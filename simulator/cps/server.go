/*
This is the server GO file and it is a model of our server
*/
package cps

import (
	"bytes"
	"fmt"
	"math"
	"sort"
	"sync"
)

var (
	falsePositives	int
	truePositives	int
	scheduler		*Scheduler
)

//FusionCenter is the server class which contains statistical, reading, and recalibration data
type FusionCenter struct {
	P *Params

	TimeBuckets 	[][]Reading //2D array where each sub array is made of the readings at one iteration
	Mean 			[]float64
	StdDev 			[]float64
	Variance 		[]float64
	Times 			map[int]bool
	LastRecal		[]int
}

//Init initializes the values for the server
func (s *FusionCenter) Init(){
	s.TimeBuckets = make([][]Reading, s.P.Iterations_used)
	s.Mean = make([]float64, s.P.Iterations_used)
	s.StdDev = make([]float64, s.P.Iterations_used)
	s.Variance = make([]float64, s.P.Iterations_used)
	s.Times = make(map[int]bool, 0)

	falsePositives = 0
	truePositives = 0

	s.LastRecal = make([]int, s.P.NumNodes) //s.P.NumNodes
}

//Reading packages the data sent by a node
type Reading struct {
	SensorVal float64
	Xpos      int
	YPos      int
	Time      int //Time represented by iteration number
	Id        int //Node Id number
}

//MakeGrid initializes a grid of Square objects according to the size of the map
func (s FusionCenter) MakeGrid() {
	p := s.P
	p.Grid = make([][]*Square, p.SquareColCM) //this creates the p.Grid and only works if row is same size as column
	for i := range p.Grid {
		p.Grid[i] = make([]*Square, p.SquareRowCM)
	}

	for i := 0; i < p.SquareColCM; i++ {
		for j := 0; j < p.SquareRowCM; j++ {

			travelList := make([]bool, 0)
			for k := 0; k < p.NumSuperNodes; k++ {
				travelList = append(travelList, true)
			}

			p.Grid[i][j] = &Square{i, j, 0.0, 0, make([]float32, p.NumGridSamples),
				p.NumGridSamples, 0.0, 0, 0, false,
				0.0, 0.0, false, travelList, sync.Mutex{}}
		}
	}
}

//CheckDetections iterates through the grid and validates detections by nodes
func (s FusionCenter) CheckDetections(p *Params, scheduler *Scheduler) {
	r := &RegionParams{}
	//scheduler := Scheduler{}

	for x := 0; x < p.SquareColCM; x++ { //k
		for y := 0; y < p.SquareRowCM; y++ { //z
			bombSquare := p.Grid[p.B.X/p.XDiv][p.B.Y/p.YDiv]
			bs_y := float64(p.B.Y / p.YDiv)
			bs_x := float64(p.B.X / p.XDiv)
			iters := p.Iterations_used

			p.Grid[x][y].StdDev = math.Sqrt(p.Grid[x][y].GetSquareValues() / float64(p.Grid[x][y].NumNodes-1))

			//check for false negatives/positives
			if p.Grid[x][y].NumNodes > 0 && float64(p.Grid[x][y].Avg) < p.DetectionThreshold && bombSquare == p.Grid[x][y] && !p.Grid[x][y].HasDetected {
				//this is a p.Grid false negative
				fmt.Fprintln(p.DriftFile, "Grid False Negative Avg:", p.Grid[x][y].Avg, "Square Row:", y, "Square Column:", x, "Iteration:", iters)
				p.Grid[x][y].HasDetected = false
			}

			if float64(p.Grid[x][y].Avg) >= p.DetectionThreshold && (math.Abs(bs_y-float64(y)) >= 1.1 && math.Abs(bs_x-float64(x)) >= 1.1) && !p.Grid[x][y].HasDetected {
				//this is a false positive
				fmt.Fprintln(p.DriftFile, "Grid False Positive Avg:", p.Grid[x][y].Avg, "Square Row:", y, "Square Column:", x, "Iteration:", iters)
				//report to supernodes
				xLoc := (x * p.XDiv) + int(p.XDiv/2)
				yLoc := (y * p.YDiv) + int(p.YDiv/2)
				p.CenterCoord = Coord{X: xLoc, Y: yLoc}
				scheduler.AddRoutePoint(p.CenterCoord, p, r)
				p.Grid[x][y].HasDetected = true
			}

			if float64(p.Grid[x][y].Avg) >= p.DetectionThreshold && (math.Abs(bs_y-float64(y)) <= 1.1 && math.Abs(bs_x-float64(x)) <= 1.1) && !p.Grid[x][y].HasDetected {
				//this is a true positive
				fmt.Fprintln(p.DriftFile, "Grid True Positive Avg:", p.Grid[x][y].Avg, "Square Row:", y, "Square Column:", x, "Iteration:", iters)
				//report to supernodes
				xLoc := (x * p.XDiv) + int(p.XDiv/2)
				yLoc := (y * p.YDiv) + int(p.YDiv/2)
				p.CenterCoord = Coord{X: xLoc, Y: yLoc}
				scheduler.AddRoutePoint(p.CenterCoord, p, r)
				p.Grid[x][y].HasDetected = true
			}

			p.Grid[x][y].SetSquareValues(0)
			p.Grid[x][y].NumNodes = 0
		}
	}
}

//Tick is performed every iteration to move supernodes and check possible detections
func (srv FusionCenter) Tick() {
	p := srv.P
	r := &RegionParams{}
	optimize := false

	for _, s := range scheduler.SNodeList {
		//Saves the current length of the super node's list of routePoints
		//If a routePoint is reached by a super node the scheduler should
		// 	reorganize the paths
		length := len(s.GetRoutePoints())

		//The super node executes it's per iteration code
		s.Tick(p, r)

		//Compares the path lengths to decide if optimization is needed
		//Optimization will only be done if he optimization requirements are met
		//	AND if the simulator is currently in a mode that requests optimization

		if length != len(s.GetRoutePoints()) {
			bombSquare := p.Grid[p.B.X/p.XDiv][p.B.Y/p.YDiv]
			sSquare := p.Grid[s.GetX()/p.XDiv][s.GetY()/p.YDiv]
			p.Grid[s.GetX()/p.XDiv][s.GetY()/p.YDiv].HasDetected = false

			bdist := float32(math.Pow(float64(math.Pow(float64(math.Abs(float64(s.GetX())-float64(p.B.X))), 2)+math.Pow(float64(math.Abs(float64(s.GetY())-float64(p.B.Y))), 2)), .5))

			if bombSquare == sSquare || bdist < 8.0 {
				p.FoundBomb = true
			} else {
				sSquare.Reset()
			}

		}

		if length != len(s.GetRoutePoints()) {
			optimize = p.DoOptimize // true &&
		}

		//Writes the super node information to a file
		fmt.Fprint(p.RoutingFile, s)
		pp := srv.printPoints(s)
		fmt.Fprint(p.RoutingFile, " UnvisitedPoints: ")
		fmt.Fprintln(p.RoutingFile, pp.String())
	}

	//Executes the optimization code if the optimize flag is true
	if optimize {
		//The scheduler optimizes the paths of each super node
		scheduler.Optimize(p, r)
		//Resets the optimize flag
		optimize = false
	}
	srv.CheckDetections(p, scheduler)

}

//printPoints saves the Coords in the allPoints list into a buffer to print to the file
func (srv FusionCenter) printPoints(s SuperNodeParent) bytes.Buffer {
	var buffer bytes.Buffer
	buffer.WriteString((fmt.Sprintf("[")))
	for ind, i := range s.GetAllPoints() {
		buffer.WriteString(i.String())

		if ind != len(s.GetAllPoints())-1 {
			buffer.WriteString(" ")
		}
	}
	buffer.WriteString((fmt.Sprintf("]")))
	return buffer
}

//MakeSuperNodes initializes the supernodes to the corners of the map
func (s FusionCenter) MakeSuperNodes() {
	p := s.P
	r := RegionParams{}

	top_left_corner := Coord{X: 0, Y: 0}
	top_right_corner := Coord{X: 0, Y: 0}
	bot_left_corner := Coord{X: 0, Y: 0}
	bot_right_corner := Coord{X: 0, Y: 0}

	tl_min := p.Height + p.Width
	tr_max := -1
	bl_max := -1
	br_max := -1

	for x := 0; x < p.Width; x++ {
		for y := 0; y < p.Height; y++ {
			if r.Point_dict[Tuple{x, y}] {
				if x+y < tl_min {
					tl_min = x + y
					top_left_corner.X = x
					top_left_corner.Y = y
				}
				if y-x > tr_max {
					tr_max = y - x
					top_right_corner.X = x
					top_right_corner.Y = y
				}
				if x-y > bl_max {
					bl_max = x - y
					bot_left_corner.X = x
					bot_left_corner.Y = y
				}
				if x+y > br_max {
					br_max = x + y
					bot_right_corner.X = x
					bot_right_corner.Y = y
				}
			}
		}
	}

	fmt.Printf("TL: %v, TR %v, BL %v, BR %v\n", top_left_corner, top_right_corner, bot_left_corner, bot_right_corner)

	starting_locs := make([]Coord, 4)
	starting_locs[0] = top_left_corner
	starting_locs[1] = top_right_corner
	starting_locs[2] = bot_left_corner
	starting_locs[3] = bot_right_corner

	//The scheduler determines which supernode should pursue a point of interest
	scheduler = &Scheduler{}

	//List of all the supernodes on the grid
	scheduler.SNodeList = make([]SuperNodeParent, p.NumSuperNodes)

	for i := 0; i < p.NumSuperNodes; i++ {
		snode_points := make([]Coord, 1)
		snode_path := make([]Coord, 0)
		all_points := make([]Coord, 0)

		//Defining the starting x and y values for the super node
		//This super node starts at the middle of the grid
		x_val, y_val := starting_locs[i].X, starting_locs[i].Y
		nodeCenter := Coord{X: x_val, Y: y_val}

		scheduler.SNodeList[i] = &Sn_zero{&Supern{&NodeImpl{X: x_val, Y: y_val, Id: i}, 1,
			1, p.SuperNodeRadius, p.SuperNodeRadius, 0, snode_points, snode_path,
			nodeCenter, 0, 0, 0, 0, 0, all_points}}

		//The super node's current location is always the first element in the routePoints list
		scheduler.SNodeList[i].UpdateLoc()
	}
}

//GetSquareAverage grabs and returns the average of a particular Square
func (s FusionCenter) GetSquareAverage(tile *Square) float32 {
	return tile.Avg
}

//UpdateSquareAvg takes a node reading and updates the parameters in the Square the node took the reading in
func (s FusionCenter) UpdateSquareAvg(rd Reading) {
	tile := s.P.Grid[rd.Xpos/s.P.XDiv][rd.YPos/s.P.YDiv]
	tile.TakeMeasurement(float32(rd.SensorVal))
}

//UpdateSquareNumNodes searches the node list and updates the number of nodes in each square
func (s FusionCenter) UpdateSquareNumNodes() {
	var node NodeImpl

	//Clear number of nodes for each square
	for i:=0; i < len(s.P.Grid); i++ {
		for j:=0; j < len(s.P.Grid[i]); j++ {
			s.P.Grid[i][j].ActualNumNodes = 0
		}
	}

	//Count number of nodes in each square
	for i:=0; i < s.P.NumNodes; i++ {
		node = s.P.NodeList[i]
		s.P.Grid[node.X/s.P.XDiv][node.Y/s.P.YDiv].ActualNumNodes += 1
	}
}

//Send is called by a node to deliver a reading to the server.
// Statistics are calculated each time data is received
func (s *FusionCenter) Send(n *NodeImpl, rd Reading) {
	//fmt.Printf("Sending to server:\nTime: %v, ID: %v, X: %v, Y: %v, Sensor Value: %v\n", rd.Time, rd.Id, rd.Xpos, rd.YPos, rd.SensorVal)
	s.Times = make(map[int]bool, 0)
	if s.Times[rd.Time] {

	} else {
		s.Times[rd.Time] = true
	}

	if len(s.TimeBuckets) <= rd.Time {
		s.TimeBuckets = append(s.TimeBuckets, make([]Reading,0))
	}
	currBucket := (s.TimeBuckets)[rd.Time]
	if len(currBucket) != 0 { //currBucket != nil
		(s.TimeBuckets)[rd.Time] = append(currBucket, rd)
	} else {
		(s.TimeBuckets)[rd.Time] = append((s.TimeBuckets)[rd.Time], rd) //s.TimeBuckets[rd.Time] = []float64{rd.sensorVal}
	}

	s.UpdateSquareAvg(rd)
	tile := s.P.Grid[rd.Xpos/s.P.XDiv][rd.YPos/s.P.YDiv]
	tile.SquareValues += math.Pow(float64(rd.SensorVal-float64(tile.Avg)), 2)
	if rd.SensorVal > (float64(s.GetSquareAverage(s.P.Grid[rd.Xpos/s.P.XDiv][rd.YPos/s.P.YDiv])) + s.P.CalibrationThresholdCM){ //Check if x over grid avg
		n.Recalibrate()
		s.LastRecal[n.Id] = s.P.Iterations_used
		//fmt.Println(s.LastRecal)
	}
}

//CalcStats calculates the mean, standard deviation and variance of the entire area at one time
func (s *FusionCenter) CalcStats() ([]float64, []float64, []float64) {
	//fmt.Printf("Calculating stats for times: %v", s.times)
	s.UpdateSquareNumNodes()

	//Calculate the mean
	sum := 0.0
	StdDevFromMean := 0.0
	for i:= range s.Times {
		for j := 0; j < len(s.TimeBuckets[i]); j++ {
			//fmt.Printf("Bucket size: %v\n", len(s.TimeBuckets[i]))
			sum += (s.TimeBuckets)[i][j].SensorVal
			//fmt.Printf("Time : %v, Elements #: %v, Value: %v\n", i, j, s.TimeBuckets[i][j])
		}
		if len(s.Mean) <= i {
			s.Mean = append(s.Mean, sum / float64( len(s.TimeBuckets[i]) ))
		} else {
			s.Mean[i] = sum / float64(len(s.TimeBuckets[i]))
		}
		sum = 0
	}

	//Calculate the standard deviation and variance
	sum = 0.0
	for i:= range s.Times {
		for j := 0; j < len((s.TimeBuckets)[i]); j++ {
			sum += math.Pow((s.TimeBuckets)[i][j].SensorVal - s.Mean[i], 2)
		}

		if len(s.Variance) <= i {
			s.Variance = append(s.Variance, sum / float64( len(s.TimeBuckets[i]) ))
		} else {
			s.Variance[i] = sum / float64(len(s.TimeBuckets[i]))
		}

		if len(s.StdDev) <= i {
			s.StdDev = append(s.StdDev, math.Sqrt(sum / float64( len((s.TimeBuckets)[i])) ))
		} else {
			s.StdDev[i] = math.Sqrt(sum / float64( len((s.TimeBuckets)[i])) )
		}

		//Determine how many std deviations data is away from mean
		for j:= range s.TimeBuckets[i] {
			StdDevFromMean = (s.TimeBuckets[i][j].SensorVal - s.Mean[i]) / s.StdDev[i]
			if StdDevFromMean > s.P.StdDevThresholdCM || StdDevFromMean < (-1.0 * s.P.StdDevThresholdCM){ //4
				//fmt.Printf("Potential detection by node %v at X:%v, Y:%v with reading %v\n", s.TimeBuckets[i][j].Id, s.TimeBuckets[i][j].Xpos, s.TimeBuckets[i][j].YPos, s.TimeBuckets[i][j].SensorVal)
				fmt.Fprintf(s.P.DetectionFile,"Potential detection by node %v at X:%v, Y:%v with reading %v\n", s.TimeBuckets[i][j].Id, s.TimeBuckets[i][j].Xpos, s.TimeBuckets[i][j].YPos, s.TimeBuckets[i][j].SensorVal)
				dist := math.Pow(float64(math.Abs(float64(s.TimeBuckets[i][j].Xpos)-float64(s.P.B.X))), 2) + math.Pow(float64(math.Abs(float64(s.TimeBuckets[i][j].YPos)-float64(s.P.B.X))), 2)
				if dist > s.P.DetectionThresholdCM {
					fmt.Fprintf(s.P.DetectionFile,"False positive!\n")
					falsePositives++
				} else {
						truePositives++
						fmt.Fprintf(s.P.DetectionFile,"True Positive\n")
				}

			}
		}
		sum = 0
	}
	return s.Mean, s.StdDev, s.Variance
}

//getMedian gets the median from a data set and returns it
func (s FusionCenter) getMedian(arr []float64) float64{
	sort.Float64s(arr)
	size := 0.0
	median := 0.0
	size = float64(len(arr))
	//index := 0
	//Check if even
	if int(size) % 2 == 0 {
		median = (arr[int(size / 2.0)] + arr[int(size / 2.0 + 1)] ) / 2
	} else {
		median = arr[int(size / 2.0 + 0.5)]
	}
	return median
}

//PrintStats prints the mean, standard deviation, and variance for the whole map at every iteration
func (s FusionCenter) PrintStats() {
	for i:= 0; i < s.P.Iterations_used; i++ {
		fmt.Printf("Time: %v, Mean: %v, Std Deviation: %v, Variance: %v\n", i, s.Mean[i], s.StdDev[i], s.Variance[i])
	}
}

//PrintStatsFile outputs statistical and detection data to log files
func (s FusionCenter) PrintStatsFile() {
	fmt.Fprintln(s.P.ServerFile, "Mean at each time:\n", s.P.Server.Mean)
	fmt.Fprintln(s.P.ServerFile, "Standard Deviations at each time:\n", s.P.Server.StdDev)
	fmt.Fprintln(s.P.ServerFile, "Variance at each time:\n", s.P.Server.Variance)
	fmt.Fprintf(s.P.DetectionFile, "Number of detections:%v\n", falsePositives + truePositives)
	fmt.Fprintf(s.P.DetectionFile, "Number of false positives:%v\n", falsePositives)
	fmt.Fprintf(s.P.DetectionFile, "Number of true positives:%v\n", truePositives)
	fmt.Fprintf(s.P.DetectionFile, "Last Recalibration times:%v\n", s.LastRecal)

}



/*
This is the OLD server GO file and it is a model of our server. It may contain compontents we want to incorporate later
*/

//This is the server's data structure for a phone (or node)
//type PhoneFile struct {
//	Id int //This is the phone's unique Id
//	Xpos []int //These are the saved x pos of the phone
//	yPos []int //These are the saved y pos of the phone
//	val []int //These are the saved values of the phone
//	Time []int //These are the saved times of the GPS/sensor readings
//	bufferSizes []int //these are the saved buffer sizes when info was dumped to server
//	speeds []int //these are the saved accelerometer based of the phone
//
//	refined [][][][]int //x,y,val,Time for all Time
//}
//
////The server is merely al list of phone files for now
//type Server struct {
//	//p [numNodes]phoneFile
//	p [200]PhoneFile
//}
////This is for later when the server becomes more advanced
//type serverThink interface {
//}
//
////This is the server absorbing data from the nodes and writing it to its phone files
//func GetData(s *Server,Xpos []int, yPos []int, val []int, Time []int, Id int, buffer int) () {
//	//s.p[Id].Xpos = append(s.p[Id].Xpos,Xpos ...)
//	//s.p[Id].yPos = append(s.p[Id].yPos,yPos...)
//	//s.p[Id].val = append(s.p[Id].val,val...)
//	//s.p[Id].Time = append(s.p[Id].Time,Time...)
//	//s.p[Id].bufferSizes = append(s.p[Id].bufferSizes,buffer)
//}
//
////This is a debugging function to be removed later
//func (s Server) String() {
//	fmt.Println("Length of string",int(len(s.p))," ")
//}
//
////This refines the phone files to fill in the gaps between where the server did not check the GPS or sensor
//func Refine( p *PhoneFile) (bool) {
//	//This fills the positions
//	if (len(s.P.yPos) == len(p.Time)) == (len(p.yPos) == len(p.val)) {
//		inbetween := 0
//		open := false
//		for i := 0; i < len(p.Time); i++ {
//			if p.Xpos[i] == -1 && p.yPos[i] == -1 {
//				inbetween += 1
//			}
//			if p.Xpos[i] != -1 && p.yPos[i] != -1 && open == true && inbetween > 0 {
//				diviserX := (p.Xpos[i] - p.Xpos[i-inbetween-1])/(inbetween+1)
//				diviserY := (p.yPos[i] - p.yPos[i-inbetween-1])/(inbetween+1)
//				for x := 0; x < inbetween; x++ {
//					p.Xpos[i-inbetween+x] = diviserX + p.Xpos[i-inbetween+x-1]
//					p.yPos[i-inbetween+x] = diviserY + p.yPos[i-inbetween+x-1]
//				}
//				inbetween = 0
//			} else if p.Xpos[i] != -1 && p.yPos[i] != -1 && open == false {
//				open = true
//				inbetween = 0
//			}
//		}
//		inbetween = 0
//		open = false
//		//This fills the values
//		for i := 1; i < len(p.Time); i++ {
//			if p.val[i] == -1 {
//				inbetween += 1
//			}
//			if p.val[i] != -1 && p.val[i-1] == -1 && inbetween > 0 && open == true {
//				diviserV := (p.val[i] - p.val[i-inbetween-1])/(inbetween+1)
//				for x:= 0; x < inbetween; x++ {
//					p.val[i-inbetween+x] = diviserV + p.val[i-inbetween+x-1]
//				}
//				inbetween = 0
//			} else if p.val[i] != -1 && open == false {
//				open = true
//				inbetween = 0
//			}
//		}
//		return true
//	} else {
//		return false
//	}
//}