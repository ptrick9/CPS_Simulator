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
)

//FusionCenter is the server class which contains statistical, reading, and recalibration data
type FusionCenter struct {
	P *Params
	R *RegionParams

	TimeBuckets 	[][]Reading //2D array where each sub array is made of the readings at one iteration
	Mean 			[]float64
	StdDev 			[]float64
	Variance 		[]float64
	Times 			map[int]bool
	LastRecal		[]int
	Sch		*Scheduler
	Readings		map[Key][]Reading
	CheckedIds		[]int
}

//Init initializes the values for the server
func (s *FusionCenter) Init(){
	s.TimeBuckets = make([][]Reading, s.P.Iterations_used)
	s.Mean = make([]float64, s.P.Iterations_of_event)
	s.StdDev = make([]float64, s.P.Iterations_of_event)
	s.Variance = make([]float64, s.P.Iterations_of_event)
	s.Times = make(map[int]bool, 0)

	falsePositives = 0
	truePositives = 0

	s.LastRecal = make([]int, s.P.TotalNodes) //s.P.TotalNodes
	s.Sch = &Scheduler{s.P, s.R, nil}

	s.Readings = make(map[Key][]Reading)
	s.CheckedIds = make([]int, 0)
}

//Reading packages the data sent by a node
type   Reading struct {
	SensorVal float64
	Xpos      float32
	YPos      float32
	Time      int //Time represented by iteration number
	Id        int //Node Id number
}

//Key for dictionary of sensor readings
type Key struct {
	X 		int
	Y		int
	Time 	int
}

//MakeGrid initializes a grid of Square objects according to the size of the map
func (s FusionCenter) MakeGrid() {
	navigable := true
	s.P.Grid = make([][]*Square, s.P.SquareColCM) //this creates the p.Grid and only works if row is same size as column
	for i := range s.P.Grid {
		s.P.Grid[i] = make([]*Square, s.P.SquareRowCM)
	}

	for i := 0; i < s.P.SquareColCM; i++ {
		for j := 0; j < s.P.SquareRowCM; j++ {

			travelList := make([]bool, 0)
			for k := 0; k < s.P.NumSuperNodes; k++ {
				travelList = append(travelList, true)
			}
			//xLoc := (i * s.P.XDiv) + int(s.P.XDiv/2)
			//yLoc := (j * s.P.YDiv) + int(s.P.YDiv/2)
			xLoc := i * s.P.XDiv
			yLoc := j * s.P.YDiv
			navigable = true
			for x:= xLoc; x < xLoc + s.P.XDiv; x++ {
				for y := yLoc; y < yLoc + s.P.YDiv; y++ {
					//fmt.Printf("X:%v, Y:%v, Region:%v\n", x, y, RegionContaining(Tuple{x, y}, s.R))
					if RegionContaining(Tuple{x, y}, s.R) == -1 {
						navigable = false
					}
				}
			}

			s.P.Grid[i][j] = &Square{i, j, 0.0, 0, make([]float32, s.P.NumGridSamples),
				s.P.NumGridSamples, 0.0, 0, 0, false,
				0.0, 0.0, false, travelList, map[Tuple]*NodeImpl{}, sync.Mutex{}, 0, navigable, false}
		}
	}
}

//CheckDetections iterates through the grid and validates detections by nodes
func (s FusionCenter) CheckDetections() {
	for x := 0; x < s.P.SquareColCM; x++ {
		for y := 0; y < s.P.SquareRowCM; y++ {
			bombSquare := s.P.Grid[s.P.B.X/s.P.XDiv][s.P.B.Y/s.P.YDiv]
			bs_y := float64(s.P.B.Y / s.P.YDiv)
			bs_x := float64(s.P.B.X / s.P.XDiv)
			iters := s.P.Iterations_used

			s.P.Grid[x][y].StdDev = math.Sqrt(s.P.Grid[x][y].GetSquareValues() / float64(s.P.Grid[x][y].NumNodes-1))

			//check for false negatives/positives
			if s.P.Grid[x][y].NumNodes > 0 && float64(s.P.Grid[x][y].Avg) < s.P.DetectionThreshold && bombSquare == s.P.Grid[x][y] && !s.P.Grid[x][y].HasDetected {
				//this is a s.P.Grid false negative
				fmt.Fprintln(s.P.DriftFile, "Grid False Negative Avg:", s.P.Grid[x][y].Avg, "Square Row:", y, "Square Column:", x, "Iteration:", iters)
				s.P.Grid[x][y].HasDetected = false
			}

			if float64(s.P.Grid[x][y].Avg) >= s.P.DetectionThreshold && (math.Abs(bs_y-float64(y)) >= 1.1 && math.Abs(bs_x-float64(x)) >= 1.1) && !s.P.Grid[x][y].HasDetected {
				//this is a false positive
				fmt.Fprintln(s.P.DriftFile, "Grid False Positive Avg:", s.P.Grid[x][y].Avg, "Square Row:", y, "Square Column:", x, "Iteration:", iters)
				//report to supernodes
				xLoc := (x * s.P.XDiv) + int(s.P.XDiv/2)
				yLoc := (y * s.P.YDiv) + int(s.P.YDiv/2)
				s.P.CenterCoord = Coord{X: xLoc, Y: yLoc}
				if s.P.SuperNodes {
					s.Sch.AddRoutePoint(s.P.CenterCoord)
				}
				s.P.Grid[x][y].HasDetected = true
			}

			if float64(s.P.Grid[x][y].Avg) >= s.P.DetectionThreshold && (math.Abs(bs_y-float64(y)) <= 1.1 && math.Abs(bs_x-float64(x)) <= 1.1) && !s.P.Grid[x][y].HasDetected {
				//this is a true positive
				fmt.Fprintln(s.P.DriftFile, "Grid True Positive Avg:", s.P.Grid[x][y].Avg, "Square Row:", y, "Square Column:", x, "Iteration:", iters)
				//report to supernodes
				xLoc := (x * s.P.XDiv) + int(s.P.XDiv/2)
				yLoc := (y * s.P.YDiv) + int(s.P.YDiv/2)
				s.P.CenterCoord = Coord{X: xLoc, Y: yLoc}
				if s.P.SuperNodes {
					s.Sch.AddRoutePoint(s.P.CenterCoord)
				}
				s.P.Grid[x][y].HasDetected = true
			}

			s.P.Grid[x][y].SetSquareValues(0)
			s.P.Grid[x][y].NumNodes = 0
		}
	}
}

//Tick is performed every iteration to move supernodes and check possible detections
func (srv FusionCenter) Tick() {
	optimize := false
	for i := range srv.Sch.SNodeList {
		srv.P.Grid[srv.Sch.SNodeList[i].GetX() / srv.P.XDiv][srv.Sch.SNodeList[i].GetY() / srv.P.YDiv].Visited = true
	}
	/*if srv.P.Iterations_used % 60 == 0 && srv.P.Iterations_used !=0{
		DensitySquares := srv.GetLeastDenseSquares()
		leastDense := DensitySquares[0]
		for i:=0; i < len(DensitySquares); i++ {
			if DensitySquares[i].Navigable {
				leastDense = DensitySquares[i]
				fmt.Printf("\nDestination Square: X:%v, Y:%v, Navigable: %v\n", leastDense.X, leastDense.Y, leastDense.Navigable)
				break
			}
		}
		if leastDense.Navigable {
			xLoc := (leastDense.X * srv.P.XDiv) + int(srv.P.XDiv/2)
			yLoc := (leastDense.Y * srv.P.YDiv) + int(srv.P.YDiv/2)
			srv.P.CenterCoord = Coord{X: xLoc, Y: yLoc}
			fmt.Printf("Destination Coordinate: %v\n",srv.P.CenterCoord)
			fmt.Printf("Destination Region:%v\n",RegionContaining(Tuple{srv.P.CenterCoord.X, srv.P.CenterCoord.Y}, srv.R))
			srv.Sch.AddRoutePoint(srv.P.CenterCoord)
		}
	}*/

	for _, s := range srv.Sch.SNodeList {
		//Saves the current length of the super node's list of routePoints
		//If a routePoint is reached by a super node the scheduler should
		// 	reorganize the paths
		length := len(s.GetRoutePoints())

		//The super node executes it's per iteration code
		s.Tick()

		//Compares the path lengths to decide if optimization is needed
		//Optimization will only be done if he optimization requirements are met
		//	AND if the simulator is currently in a mode that requests optimization

		if length != len(s.GetRoutePoints()) {
			bombSquare := srv.P.Grid[srv.P.B.X/srv.P.XDiv][srv.P.B.Y/srv.P.YDiv]
			sSquare := srv.P.Grid[s.GetX()/srv.P.XDiv][s.GetY()/srv.P.YDiv]
			srv.P.Grid[s.GetX()/srv.P.XDiv][s.GetY()/srv.P.YDiv].HasDetected = false

			bdist := float32(math.Pow(float64(math.Pow(float64(math.Abs(float64(s.GetX())-float64(srv.P.B.X))), 2)+math.Pow(float64(math.Abs(float64(s.GetY())-float64(srv.P.B.Y))), 2)), .5))

			if bombSquare == sSquare || bdist < 8.0 {
				srv.P.FoundBomb = true
			} else {
				sSquare.Reset()
			}

		}

		if length != len(s.GetRoutePoints()) {
			optimize = srv.P.DoOptimize // true &&
		}

		//Writes the super node information to a file
		fmt.Fprint(srv.P.RoutingFile, s)
		pp := srv.printPoints(s)
		fmt.Fprint(srv.P.RoutingFile, " UnvisitedPoints: ")
		fmt.Fprintln(srv.P.RoutingFile, pp.String())
	}

	//Executes the optimization code if the optimize flag is true
	if optimize {
		//The scheduler optimizes the paths of each super node
		srv.Sch.Optimize()
		//Resets the optimize flag
		optimize = false
	}
	srv.CheckDetections()

}

//returns all of the nodes a radial distance from the current node
func (s* FusionCenter) NodesInRadius(curNode * NodeImpl, radius int)(map[Tuple]*NodeImpl) {
	var gridMaxX = s.P.MaxX;
	var gridMaxY = s.P.MaxY;

	var nodesInRadius = map[Tuple]*NodeImpl{}

	var negRadius = -1*radius;

	//iterate over the Grid by row and column
	for row := negRadius; row<=radius; row++{
		for col := negRadius; col<=radius; col++{
			//do not include current node in list of nodes in radius
			if(row == 0 && col == 0){
				continue
			}

			var testX = int(curNode.X) + col					//test X Value
			var testY = int(curNode.Y) + row					//test Y Value
			var testTup = Tuple{testX, testY}	//create Tuple from test X and Y values
			if(testX < gridMaxX && testX >= 0){			//if the testX Value is on the grid, continue
				if(testY < gridMaxY && testY >= 0){		//if the testY Value is on the grid, continue
					if(s.P.NodePositionMap[testTup] != nil){	//if the test position has a Node, continue
						nodesInRadius[testTup] = s.P.NodePositionMap[testTup]	//add the node to the nodesInRadius map
					}
				}
			}
		}
	}
	return nodesInRadius
}


//--XY ERROR?--
//returns all of the nodes dist squares away from the current node
func (s* FusionCenter) NodesWithinDistance(curNode * NodeImpl, dist int)(map[Tuple]*NodeImpl){
	var gridMaxX = s.P.MaxX;
	var gridMaxY = s.P.MaxY;
	var nodesWithinDist = s.P.Grid[int(curNode.Y)][int(curNode.X)].NodesInSquare //initialize to nodes in current square
	var negDist = -1*dist;

	for row := negDist; row<=dist; row++ {
		for col := negDist; col <= dist; col++ {

			var testX = s.P.Grid[int(curNode.Y)][int(curNode.X)].X + col		//X Value of test Square
			var testY = s.P.Grid[int(curNode.Y)][int(curNode.X)].Y + row		//Y Value of test Square

			if(testX < gridMaxX && testX >= 0){			//if the testX Value is on the grid, continue
				if(testY < gridMaxY && testY >= 0){		//if the testY Value is on the grid, continue
					var testSquare =  s.P.Grid[testY][testX] 			//create Square from test X and Y values
					if(testSquare != nil){					//if the test Square is not null, continue
						for ind,val := range testSquare.NodesInSquare{	//iterate through nodes in square map adding each to the
							nodesWithinDist[ind] = val;					//nodes within Distance Map
						}
					}
				}
			}
		}
	}
	return nodesWithinDist
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

	top_left_corner := Coord{X: 0, Y: 0}
	top_right_corner := Coord{X: 0, Y: 0}
	bot_left_corner := Coord{X: 0, Y: 0}
	bot_right_corner := Coord{X: 0, Y: 0}

	tl_min := s.P.Height + s.P.Width
	tr_max := -1
	bl_max := -1
	br_max := -1

	for x := 0; x < s.P.Width; x++ {
		for y := 0; y < s.P.Height; y++ {
			if s.R.Point_dict[Tuple{x, s.P.Height - y - 1}] { //
				if x+y < tl_min {
					tl_min = x + y
					top_left_corner.X = x
					top_left_corner.Y = y
				}
				if x+y > tr_max {
					tr_max = x + y
					top_right_corner.X = x
					top_right_corner.Y = y
				}
				if y-x > bl_max {
					bl_max = y - x
					bot_left_corner.X = x
					bot_left_corner.Y = y
				}
				if x-y > br_max {
					br_max = x - y
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
	//scheduler = &Scheduler{s.P, s.R, nil}

	//List of all the supernodes on the grid
	s.Sch.SNodeList = make([]SuperNodeParent, s.P.NumSuperNodes)

	for i := 0; i < s.P.NumSuperNodes; i++ {
		snode_points := make([]Coord, 1)
		snode_path := make([]Coord, 0)
		all_points := make([]Coord, 0)

		//Defining the starting x and y values for the super node
		//This super node starts at the middle of the grid
		x_val, y_val := starting_locs[i].X, starting_locs[i].Y
		nodeCenter := Coord{X: x_val, Y: y_val}

		s.Sch.SNodeList[i] = &Sn_zero{s.P, s.R,&Supern{s.P,s.R,&NodeImpl{X: float32(x_val), Y: float32(y_val), Id: i}, 1,
			1, s.P.SuperNodeRadius, s.P.SuperNodeRadius, 0, snode_points, snode_path,
			nodeCenter, 0, 0, 0, 0, 0, all_points}}

		//The super node's current location is always the first element in the routePoints list
		s.Sch.SNodeList[i].UpdateLoc()
	}
}

//GetSquareAverage grabs and returns the average of a particular Square
func (s FusionCenter) GetSquareAverage(tile *Square) float32 {
	return tile.Avg
}

//UpdateSquareAvg takes a node reading and updates the parameters in the Square the node took the reading in
func (s FusionCenter) UpdateSquareAvg(rd Reading) {
	tile := s.P.Grid[int(rd.Xpos)/s.P.XDiv][int(rd.YPos)/s.P.YDiv]
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
	for i:=0; i < len(s.P.NodeList); i++ {
		node = *s.P.NodeList[i]
		if node.Valid {
			s.P.Grid[int(node.X)/s.P.XDiv][int(node.Y)/s.P.YDiv].ActualNumNodes += 1
		}
	}
}

//Send is called by a node to deliver a reading to the server.
// Statistics are calculated each Time data is received
func (s *FusionCenter) Send(n *NodeImpl, rd Reading) {
	//fmt.Printf("Sending to server:\nTime: %v, ID: %v, X: %v, Y: %v, Sensor Value: %v\n", rd.Time, rd.Id, rd.Xpos, rd.YPos, rd.SensorVal)
	_, ok := s.Readings[Key{int(rd.Xpos / float32(s.P.XDiv)), int(rd.YPos / float32(s.P.YDiv)), rd.Time/1000}]
	if ok {
		s.Readings[Key{int(rd.Xpos / float32(s.P.XDiv)), int(rd.YPos / float32(s.P.YDiv)), rd.Time / 1000}] = append(s.Readings[Key{int(rd.Xpos / float32(s.P.XDiv)), int(rd.YPos / float32(s.P.YDiv)), rd.Time / 1000}], rd)
	} else {
		s.Readings[Key{int(rd.Xpos / float32(s.P.XDiv)), int(rd.YPos / float32(s.P.YDiv)), rd.Time / 1000}] = []Reading{rd}
	}
	s.Times = make(map[int]bool, 0)
	if s.Times[rd.Time] {

	} else {
		s.Times[rd.Time] = true
	}

	for len(s.TimeBuckets) <= rd.Time {
		s.TimeBuckets = append(s.TimeBuckets, make([]Reading,0))
	}
	currBucket := (s.TimeBuckets)[rd.Time]
	if len(currBucket) != 0 { //currBucket != nil
		(s.TimeBuckets)[rd.Time] = append(currBucket, rd)
	} else {
		(s.TimeBuckets)[rd.Time] = append((s.TimeBuckets)[rd.Time], rd) //s.TimeBuckets[rd.Time] = []float64{rd.sensorVal}
	}

	s.UpdateSquareAvg(rd)
	tile := s.P.Grid[int(rd.Xpos)/s.P.XDiv][int(rd.YPos)/s.P.YDiv]
	tile.LastReadingTime = rd.Time
	tile.SquareValues += math.Pow(float64(rd.SensorVal-float64(tile.Avg)), 2)
	if rd.SensorVal > (float64(s.GetSquareAverage(s.P.Grid[int(rd.Xpos)/s.P.XDiv][int(rd.YPos)/s.P.YDiv])) + s.P.CalibrationThresholdCM){ //Check if x over grid avg
		n.Recalibrate()
		s.LastRecal[n.Id] = s.P.Iterations_used
		//fmt.Println(s.LastRecal)
	}

	if rd.SensorVal > s.P.DetectionThreshold {
		s.CheckedIds = make([]int, 0)
		validations := 0
		for t:= (s.P.CurrentTime / 1000) - 60; t <= s.P.CurrentTime / 1000; t++ {
			for x:= int((rd.Xpos - float32(s.P.DetectionDistance)) / float32(s.P.XDiv)); x < int((rd.Xpos + float32(s.P.DetectionDistance) )/ float32(s.P.XDiv)); x++ {
				for y:= int((rd.YPos - float32(s.P.DetectionDistance)) / float32(s.P.YDiv)); y < int((rd.YPos + float32(s.P.DetectionDistance) )/ float32(s.P.YDiv)); y++ {
					for r:= range s.Readings[Key{x,y,t}] {
						if Dist(Tuple{int(s.Readings[Key{x,y,t}][r].Xpos), int(s.Readings[Key{x,y,t}][r].YPos)}, Tuple{int(rd.Xpos), int(rd.YPos)}) < s.P.DetectionDistance {
							if s.Readings[Key{x,y,t}][r].Id != rd.Id && !Is_in(s.Readings[Key{x,y,t}][r].Id, s.CheckedIds) &&
								s.Readings[Key{x,y,t}][r].SensorVal > s.P.DetectionThreshold {
								s.CheckedIds = append(s.CheckedIds, s.Readings[Key{x,y,t}][r].Id)
								validations++
							}
						}
					}
				}
			}
		}
		if validations >= s.P.ValidationThreshold {

			//fmt.Println("Valid!")
			s.P.CenterCoord = Coord{X: int(rd.Xpos), Y: int(rd.YPos)}
			if s.P.SuperNodes {
				s.Sch.AddRoutePoint(s.P.CenterCoord)
			}
			if FloatDist(Tuple32{rd.Xpos, rd.YPos}, Tuple32{float32(s.P.B.X), float32(s.P.B.Y)}) > s.P.DetectionDistance {
				fmt.Fprintln(s.P.DetectionFile, fmt.Sprintf("FP Confirmation T: %v ID: %v (%v, %v) D: %v C: %v", rd.Time, rd.Id, rd.Xpos, rd.YPos, FloatDist(Tuple32{rd.Xpos, rd.YPos}, Tuple32{float32(s.P.B.X), float32(s.P.B.Y)}) , rd.SensorVal))
			} else {
				fmt.Fprintln(s.P.DetectionFile, fmt.Sprintf("TP Confirmation T: %v ID: %v (%v, %v) D: %v C: %v", rd.Time, rd.Id, rd.Xpos, rd.YPos, FloatDist(Tuple32{rd.Xpos, rd.YPos}, Tuple32{float32(s.P.B.X), float32(s.P.B.Y)}), rd.SensorVal))
			}

		} else {
			//fmt.Println(rd)
		}
	}
}

//CalcStats calculates the mean, standard deviation and variance of the entire area at one Time
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
		for len(s.Mean) <= i {
			s.Mean = append(s.Mean, sum / float64( len(s.TimeBuckets[i]) ))
		} /*else {
			s.Mean[i] = sum / float64(len(s.TimeBuckets[i]))
		}*/
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

		for len(s.StdDev) <= i {
			s.StdDev = append(s.StdDev, math.Sqrt(sum / float64( len((s.TimeBuckets)[i])) ))
		} /*else {
			s.StdDev[i] = math.Sqrt(sum / float64( len((s.TimeBuckets)[i])) )
		}*/

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
func (s FusionCenter) GetMedian(arr []float64) float64{
	sort.Float64s(arr)
	size := 0.0
	median := 0.0
	size = float64(len(arr))
	//Index := 0
	//Check if even
	if int(size) % 2 == 0 {
		median = (arr[int(size / 2.0)] + arr[int(size / 2.0 - 1)] ) / 2
	} else {
		median = arr[int(size / 2.0 - 0.5)]
	}
	return median
}

func (s FusionCenter) GetLeastDenseSquares() Squares{
	orderedSquares := make(Squares, 0)
	for x := 0; x < len(s.P.Grid); x++ {
		for y := 0; y < len(s.P.Grid[x]); y++ {
			if !s.P.Grid[x][y].Visited {
				orderedSquares = append(orderedSquares, s.P.Grid[x][y])
			}
		}
	}
	sort.Sort(&orderedSquares)
	/*for i:= range orderedSquares {
		//total+=orderedSquares[i].ActualNumNodes
		//fmt.Printf("X:%v, Y:%v, Density:%v\n", orderedSquares[i].X, orderedSquares[i].Y, orderedSquares[i].ActualNumNodes)
	}*/

	return orderedSquares
}

type Squares []*Square

func (s *Squares) Len() int{
	return len(*s)
}

func (s *Squares) Swap(i, j int) {
	(*s)[i], (*s)[j] = (*s)[j], (*s)[i]
}

func (s *Squares) Less(i, j int) bool{
	return (*s)[i].ActualNumNodes < (*s)[j].ActualNumNodes
}

//PrintStats prints the mean, standard deviation, and variance for the whole map at every iteration
func (s FusionCenter) PrintStats() {
	for i:= 0; i < s.P.Iterations_used; i++ {
		fmt.Printf("Time: %v, Mean: %v, Std Deviation: %v, Variance: %v\n", i, s.Mean[i], s.StdDev[i], s.Variance[i])
	}
}

//PrintStatsFile outputs statistical and detection data to log files
func (s FusionCenter) PrintStatsFile() {
	fmt.Fprintln(s.P.ServerFile, "Mean at each Time:\n", s.P.Server.Mean)
	fmt.Fprintln(s.P.ServerFile, "Standard Deviations at each Time:\n", s.P.Server.StdDev)
	fmt.Fprintln(s.P.ServerFile, "Variance at each Time:\n", s.P.Server.Variance)
	fmt.Fprintf(s.P.DetectionFile, "Number of detections:%v\n", falsePositives + truePositives)
	fmt.Fprintf(s.P.DetectionFile, "Number of false positives:%v\n", falsePositives)
	fmt.Fprintf(s.P.DetectionFile, "Number of true positives:%v\n", truePositives)
	fmt.Fprintf(s.P.DetectionFile, "Last Recalibration times:%v\n", s.LastRecal)

}