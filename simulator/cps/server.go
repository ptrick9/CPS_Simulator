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

	TimeBuckets  [][]Reading //2D array where each sub array is made of the readings at one iteration
	Mean         []float64
	StdDev       []float64
	Variance     []float64
	Times        map[int]bool
	LastRecal    []int
	Sch          *Scheduler
	Readings     map[Key][]Reading
	CheckedIds   []int
	NodeDataList []NodeData
	Validators   map[int]int     //stores validators...id -> time  latest time for id is stored
	NodeSquares  map[int]Tuple   //store what square a node is in
	SquarePop    map[Tuple][]int //store nodes in square
	SquareTime   map[Tuple]TimeTrack

	SamplesCounter		int // counter keeps track when a sample is taken
	BluetoothCounter	int // counter keeps track when bluetooth communication occurs

	GlobalReclusterBTCounter	int
	LocalReclusterBTCounter	int
	ReadingBTCounter	int
	ClusterSearchBTCounter	int

	WifiCounter			int // counter keeps track when wifi communication occurs


	NodeTree		*Quadtree //stores node locations in quadtree format

	Clusters	map[*NodeImpl]*Cluster
	ClusterHeadsOf	map[*NodeImpl][]*NodeImpl
	AloneNodes		map[*NodeImpl]int
	RatioBeforeRecluster	float64
	Increments	int
	Decrements	int

	Waiting		bool	//true after a global recluster until nodes accounted for is greater than server ready threshold
	NextReclusterTime	int //The next time a global recluster is scheduled
}

type Cluster struct {
	Members		map[*NodeImpl]int
	TimeFormed	int
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
	s.NodeDataList = make([]NodeData, s.P.TotalNodes)
	s.Validators = make(map[int]int)
	s.NodeSquares = make(map[int]Tuple)
	s.SquarePop = make(map[Tuple][]int)
	s.SquareTime = make(map[Tuple]TimeTrack)

	if s.P.GlobalRecluster == 2 {
		s.NextReclusterTime = int(s.P.ReclusterPeriod)
	}
}

func (s *FusionCenter) MakeNodeData() {
	for i := range s.P.NodeList {
		s.NodeDataList[i] = NodeData{i, s.P.NodeList[i].S0,s.P.NodeList[i].S1, s.P.NodeList[i].S2,
			s.P.NodeList[i].E0, s.P.NodeList[i].E1, s.P.NodeList[i].E2, s.P.NodeList[i].ET1,
			s.P.NodeList[i].ET2, []int{ 0 }, []int{ 0 }}
	}
}

//Reading packages the data sent by a node
type Reading struct {
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

//Holds last time node was sampled in a square and true or false val for if it has already increased in threshold
type TimeTrack struct {
	TimeSample    int
	BeenReported  bool
	MaxDelta      int
}

type NodeData struct {
	Id 			int
	S0, S1, S2, E0, E1, E2, ET1, ET2 float64
	RecalTimes		[]int
	SelfRecalTimes  []int
}

//MakeGrid initializes a grid of Square objects according to the size of the map
func (s FusionCenter) MakeGrid() {
	navigable := true


	/*s.P.Grid = make([][]*Square, s.P.MaxX/s.P.SquareColCM + 1) //this creates the p.Grid and only works if row is same size as column
	for i := range s.P.Grid {
		s.P.Grid[i] = make([]*Square, s.P.MaxY/s.P.SquareRowCM + 1)
	}*/
	s.P.Grid = make([][]*Square, s.P.GridWidth) //this creates the p.Grid and only works if row is same size as column
	for i := range s.P.Grid {
		s.P.Grid[i] = make([]*Square, s.P.GridHeight)
	}

	fmt.Printf("squares = %v %v\n", s.P.MaxX/s.P.SquareColCM, s.P.MaxY/s.P.SquareRowCM)
	for i := 0; i < s.P.GridWidth; i++ {
		for j := 0; j < s.P.GridHeight; j++ {

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

			s.P.Grid[i][j] = &Square{i, j, 3.0, 0, make([]float32, s.P.NumGridSamples),
				s.P.NumGridSamples, 3.0, 0, 0, false,
				0.0, 0.0, false, travelList, map[Tuple]*NodeImpl{}, sync.Mutex{}, 0, navigable, false}
		}
	}
}

//CheckDetections iterates through the grid and validates detections by nodes
func (s FusionCenter) CheckDetections() {
	for x := 0; x < s.P.GridWidth; x++ {
		for y := 0; y < s.P.GridHeight; y++ {
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

func (s *FusionCenter) CleanupReadings() {
	if s.P.CurrentTime / 1000 > s.P.ReadingHistorySize{
		/*for r := range s.Readings {
			if r.Time < (rd.Time / 1000 - s.P.ReadingHistorySize) {
				delete(s.Readings, r)
			}
		}*/
		t := (s.P.CurrentTime / 1000) - s.P.ReadingHistorySize - 4
		//for t := (s.P.CurrentTime / 1000) - s.P.ReadingHistorySize - 4; t < (s.P.CurrentTime / 1000) - s.P.ReadingHistorySize - 1; t++ {
			for x := 0; x < s.P.Width/s.P.XDiv; x++ {
				for y := 0; y < s.P.Height/s.P.YDiv; y++ {
					_, ok := s.Readings[Key{x, y, t}]
					if ok {
						//fmt.Printf("Deleting time %v\n", t)
						delete(s.Readings, Key{x, y, t})
					}
				}
			}
		//}

		for nodeId, time := range s.Validators{
			if (s.P.CurrentTime/1000) - time > s.P.ReadingHistorySize {
				//this id is too old and must be deleted. We must also make a not that THIS is a false negative
				delete(s.Validators, nodeId)
				fmt.Fprintf(s.P.DetectionFile, "FN Window T: %v CT: %v ID: %v\n", time, (s.P.CurrentTime)/1000, nodeId)
			}
		}
	}
}

func pos(value int, array []int) int {
	for p, v := range array {
		if (v == value) {
			return p
		}
	}
	return -1
}

func remove(s []int, i int) []int {
	s[i] = s[0]
	return s[1:]
}


//Send is called by a node to deliver a reading to the server.
// Statistics are calculated each Time data is received
func (s *FusionCenter) Send(n *NodeImpl, rd *Reading, tp bool) {
	//fmt.Printf("Sending to server:\nTime: %v, ID: %v, X: %v, Y: %v, Sensor Value: %v\n", rd.Time, rd.Id, rd.Xpos, rd.YPos, rd.SensorVal)
	//NodeSquares 	map[int]Tuple  //store what square a node is in
	//SquarePop  		map[Tuple][]int //store nodes in square
	//SquareTime     map[Tuple][]TimeTrack stores last time a reading was taken in a square
	v, ok := s.NodeSquares[n.Id] //check if node has been recorded before
	newSquare := Tuple{int(rd.Xpos / float32(s.P.XDiv)), int(rd.YPos / float32(s.P.YDiv))}
	if ok { //node has been recorded before
		if v == newSquare { //no changes to make

		} else { //node has changed squares, need to update square pop
			elements := s.SquarePop[v]
			ind := pos(n.Id, elements)
			s.SquarePop[v] = remove(elements, ind)

			_, ok = s.SquarePop[newSquare] //check if square exists
			if ok {//exists, so just append
				s.SquarePop[newSquare] = append(s.SquarePop[newSquare], n.Id) //add id to new square
			} else {//does not exist, create
				s.SquarePop[newSquare] = []int{n.Id}
			}
			s.NodeSquares[n.Id] = newSquare //update nodes square log
		}
	} else {
		_, ok = s.SquarePop[newSquare] //check if square exists
		if ok {//exists, so just append
			s.SquarePop[newSquare] = append(s.SquarePop[newSquare], n.Id) //add id to new square
		} else {//does not exist, create
			s.SquarePop[newSquare] = []int{n.Id}
		}
		s.NodeSquares[n.Id] = newSquare //update nodes square log
	}

	s.CheckSquares(newSquare)
	s.SquareTime[newSquare] = TimeTrack{n.P.CurrentTime, false,s.SquareTime[newSquare].MaxDelta}
	// add node to correct square


	recalReject := false


	tile := s.P.Grid[int(rd.Xpos)/s.P.XDiv][int(rd.YPos)/s.P.YDiv]
	tile.LastReadingTime = rd.Time
	tile.SquareValues += math.Pow(float64(rd.SensorVal-float64(tile.Avg)), 2)
	if s.P.ServerRecal {
		if rd.SensorVal > (float64(tile.Avg) + s.P.CalibrationThresholdCM) { //Check if x over grid avg
			if s.P.RecalReject {
				if ((s.P.CurrentTime/1000)-n.NodeTime) > 200 { //hasn't been recalibrated too recently, need to reject and recal
					fmt.Fprintf(n.P.DriftExploreFile, "ID: %v T: %v In: %v CUR: %v NT: %v %v SERVRECAL\n", n.Id, n.P.CurrentTime, n.InitialSensitivity, n.Sensitivity, n.NodeTime, rd.SensorVal)
					n.Recalibrate()
					s.LastRecal[n.Id] = s.P.Iterations_used
					recalReject = true
					if(tp) {
						fmt.Fprintf(s.P.DetectionFile, "FN Recal T: %v ID: %v\n", rd.Time, rd.Id)
					}
				} else {
					recalReject = false
				}
			} else {
				fmt.Fprintf(n.P.DriftExploreFile, "ID: %v T: %v In: %v CUR: %v NT: %v %v SERVRECAL\n", n.Id, n.P.CurrentTime, n.InitialSensitivity, n.Sensitivity, n.NodeTime, rd.SensorVal)
				n.Recalibrate()
				s.LastRecal[n.Id] = s.P.Iterations_used
			}
		}
	}

	if !recalReject { //if reading shouldn't be rejected
		s.UpdateSquareAvg(*rd)

		if rd.SensorVal > s.P.DetectionThreshold {

			_, ok := s.Readings[Key{int(rd.Xpos / float32(s.P.XDiv)), int(rd.YPos / float32(s.P.YDiv)), rd.Time / 1000}]
			if ok {
				s.Readings[Key{int(rd.Xpos / float32(s.P.XDiv)), int(rd.YPos / float32(s.P.YDiv)), rd.Time / 1000}] = append(s.Readings[Key{int(rd.Xpos / float32(s.P.XDiv)), int(rd.YPos / float32(s.P.YDiv)), rd.Time / 1000}], *rd)
			} else {
				s.Readings[Key{int(rd.Xpos / float32(s.P.XDiv)), int(rd.YPos / float32(s.P.YDiv)), rd.Time / 1000}] = []Reading{*rd}
			}
			s.Times = make(map[int]bool, 0)
			if s.Times[rd.Time] {

			} else {
				s.Times[rd.Time] = true
			}
		}


	}


	if(s.P.ValidationType == "square") {
		if rd.SensorVal > s.P.DetectionThreshold {
			if tile.Avg > float32(s.P.DetectionThreshold) && tile.NumEntry > 2 {
				if s.P.RecalReject && recalReject { //true && true means we are rejecting and we should reject
					fmt.Fprintln(s.P.DetectionFile, fmt.Sprintf("Rejection T: %v ID: %v %v/%v", rd.Time, rd.Id, tile.NumEntry, len(s.SquarePop[newSquare])))
				} else { //any other case
					if FloatDist(Tuple32{rd.Xpos, rd.YPos}, Tuple32{float32(s.P.B.X), float32(s.P.B.Y)}) > s.P.DetectionDistance {
						//do nothing here
					} else {
						if !s.P.DriftExplorer && tp {
							s.P.FoundBomb = true
						}
					}
					fmt.Fprintln(s.P.DetectionFile, fmt.Sprintf("Confirmation T: %v ID: %v %v/%v %v", rd.Time, rd.Id, tile.NumEntry, len(s.SquarePop[newSquare]), s.CheckedIds))

				}
			} else {
				/*if s.P.RecalReject && recalReject {
					fmt.Fprintln(s.P.DetectionFile, fmt.Sprintf("Rejection T: %v ID: %v %v/%v", rd.Time, rd.Id, tile.NumEntry, len(s.SquarePop[newSquare])))
				} else {
					fmt.Fprintln(s.P.DetectionFile, fmt.Sprintf("Rejection T: %v ID: %v %v/%v", rd.Time, rd.Id, tile.NumEntry, len(s.SquarePop[newSquare])))
				}*/
				fmt.Fprintln(s.P.DetectionFile, fmt.Sprintf("Rejection T: %v ID: %v %v/%v", rd.Time, rd.Id, tile.NumEntry, len(s.SquarePop[newSquare])))
			}
		}
	} else if(s.P.ValidationType == "validators") {
		if rd.SensorVal > s.P.DetectionThreshold {
			s.CheckedIds = make([]int, 0)
			validations := 0
			if !recalReject {
				for t := (s.P.CurrentTime / 1000) - s.P.ReadingHistorySize; t <= s.P.CurrentTime/1000; t++ {
					for x := int((rd.Xpos - float32(s.P.DetectionDistance)) / float32(s.P.XDiv)); x < int((rd.Xpos+float32(s.P.DetectionDistance))/float32(s.P.XDiv)); x++ {
						for y := int((rd.YPos - float32(s.P.DetectionDistance)) / float32(s.P.YDiv)); y < int((rd.YPos+float32(s.P.DetectionDistance))/float32(s.P.YDiv)); y++ {
							for r := range s.Readings[Key{x, y, t}] {
								currRead := s.Readings[Key{x, y, t}][r]
								if FloatDist(Tuple32{currRead.Xpos, currRead.YPos}, Tuple32{rd.Xpos, rd.YPos}) < s.P.DetectionDistance {
									if currRead.Id != rd.Id && !Is_in(currRead.Id, s.CheckedIds) && currRead.SensorVal > s.P.DetectionThreshold {
										s.CheckedIds = append(s.CheckedIds, currRead.Id)
										if tp {
											s.Validators[rd.Id] = t
										}
										validations++
									} else if currRead.SensorVal > s.P.DetectionThreshold && tp {
										s.Validators[rd.Id] = t
									}
								}
							}
						}
					}
				}
			}
			//if validations >= s.P.ValidationThreshold {
			//offset := math.Ceil(math.Log(float64(s.P.TotalNodes)) / math.Log(10) / 2)
			//neededValidators := int(offset) + int(math.Sqrt(float64(len(s.SquarePop[newSquare]))))
			neededValidators := s.P.ValidationThreshold
			if validations >= neededValidators {
				if s.P.RecalReject && recalReject {
					fmt.Fprintln(s.P.DetectionFile, fmt.Sprintf("Rejection T: %v ID: %v %v/%v", rd.Time, rd.Id, validations, neededValidators))
				} else {
					s.P.CenterCoord = Coord{X: int(rd.Xpos), Y: int(rd.YPos)}
					if s.P.SuperNodes {
						s.Sch.AddRoutePoint(s.P.CenterCoord)
					}
					if FloatDist(Tuple32{rd.Xpos, rd.YPos}, Tuple32{float32(s.P.B.X), float32(s.P.B.Y)}) > s.P.DetectionDistance {
					} else {
						if !s.P.DriftExplorer && tp {
							s.P.FoundBomb = true
						}
					}
					fmt.Fprintln(s.P.DetectionFile, fmt.Sprintf("Confirmation T: %v ID: %v %v/%v %v", rd.Time, rd.Id, validations, neededValidators, s.CheckedIds))
				}
			} else {
				/*if s.P.RecalReject && recalReject {
					fmt.Fprintln(s.P.DetectionFile, fmt.Sprintf("Rejection T: %v ID: %v %v/%v", rd.Time, rd.Id, validations, neededValidators))
				} else {
					fmt.Fprintln(s.P.DetectionFile, fmt.Sprintf("Rejection T: %v ID: %v %v/%v", rd.Time, rd.Id, validations, neededValidators))
				}*/
				fmt.Fprintln(s.P.DetectionFile, fmt.Sprintf("Rejection T: %v ID: %v %v/%v", rd.Time, rd.Id, validations, neededValidators))
			}
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

func (s *FusionCenter) PrintBatteryStats() {

	lowestBattery := s.P.NodeList[0].GetBatteryPercentage()

	averageRemainingBattery := 0.0
	for _, node := range s.P.NodeList {
		battery := node.GetBatteryPercentage()
		averageRemainingBattery += battery
		if battery < lowestBattery {
			lowestBattery = battery
		}
	}

	fmt.Print("\nTotal Samples Taken:", s.SamplesCounter)
	fmt.Print("\nSampling Energy Consumption:", s.SamplesCounter * s.P.SampleLossAmount())
	fmt.Print("\nMinimum Remaining Battery:", lowestBattery)
	fmt.Print("\nAverage Remaining Battery:", averageRemainingBattery / float64(s.P.TotalNodes))
	fmt.Print("\nTotal Dead Nodes:", s.P.TotalNodes - len(s.P.AliveNodes), "/", s.P.TotalNodes)

	fmt.Fprintf(s.P.BatteryFile, "\nTotal Samples Taken: %v", s.SamplesCounter)
	fmt.Fprintf(s.P.BatteryFile, "\nSampling Energy Consumption: %v", s.SamplesCounter * s.P.SampleLossAmount())
	fmt.Fprintf(s.P.BatteryFile, "\nMinimum Remaining Battery: %v", lowestBattery)
	fmt.Fprintf(s.P.BatteryFile, "\nAverage Remaining Battery: %v", averageRemainingBattery / float64(s.P.TotalNodes))
	fmt.Fprintf(s.P.BatteryFile, "\nTotal Dead Nodes: %v/%v", s.P.TotalNodes - len(s.P.AliveNodes), s.P.TotalNodes)
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
	var buffer bytes.Buffer
	buffer.WriteString(fmt.Sprintf("Mean at each Time:%v\n", s.P.Server.Mean))
	buffer.WriteString(fmt.Sprintf("Standard Deviations at each Time:%v\n", s.P.Server.StdDev))
	buffer.WriteString(fmt.Sprintf("Variance at each Time:%v\n", s.P.Server.Variance))
	fmt.Fprintln(s.P.ServerFile, buffer.String())
	buffer.Reset()

	for i:= range s.NodeDataList {
		buffer.WriteString(fmt.Sprintf("ID%v,", s.NodeDataList[i].Id))
		buffer.WriteString(fmt.Sprintf("E0%v,", s.NodeDataList[i].E0))
		buffer.WriteString(fmt.Sprintf("E1%v,", s.NodeDataList[i].E1))
		buffer.WriteString(fmt.Sprintf("E2%v,", s.NodeDataList[i].E2))
		buffer.WriteString(fmt.Sprintf("ET1%v,", s.NodeDataList[i].ET1))
		buffer.WriteString(fmt.Sprintf("ET2%v,", s.NodeDataList[i].ET2))
		buffer.WriteString(fmt.Sprintf("S0%v,", s.NodeDataList[i].S0))
		buffer.WriteString(fmt.Sprintf("S1%v,", s.NodeDataList[i].S1))
		buffer.WriteString(fmt.Sprintf("S2%v,", s.NodeDataList[i].S2))
		buffer.WriteString(fmt.Sprintf("RT%v,", s.NodeDataList[i].RecalTimes))
		buffer.WriteString(fmt.Sprintf("SRT%v", s.NodeDataList[i].SelfRecalTimes))
		buffer.WriteString(fmt.Sprintln(""))
	}
	fmt.Fprintln(s.P.NodeDataFile, buffer.String())

	fmt.Fprintf(s.P.DetectionFile, "Number of detections:%v\n", falsePositives + truePositives)
	fmt.Fprintf(s.P.DetectionFile, "Number of false positives:%v\n", falsePositives)
	fmt.Fprintf(s.P.DetectionFile, "Number of true positives:%v\n", truePositives)
	fmt.Fprintf(s.P.DetectionFile, "Last Recalibration times:%v\n", s.LastRecal)

}

//Intersects returns true if line segment 'p1q1' and 'p2q2' intersect.
func Intersects(p1, q1, p2, q2 Coord) bool {
	//Orientations
	o1 := orientation(p1, q1, p2)
	o2 := orientation(p1, q1, q2)
	o3 := orientation(p2, q2, p1)
	o4 := orientation(p2, q2, q1)

	if o1 != o2 && o3 != o4 {
		return true
	}
	if o1 == 0 && onSegment(p1, p2, q1) {
		return true
	}

	// p1, q1 and q2 are colinear and q2 lies on segment p1q1
	if o2 == 0 && onSegment(p1, q2, q1) {
		return true
	}

	// p2, q2 and p1 are colinear and p1 lies on segment p2q2
	if o3 == 0 && onSegment(p2, p1, q2) {
		return true
	}

	// p2, q2 and q1 are colinear and q1 lies on segment p2q2
	if o4 == 0 && onSegment(p2, q1, q2) {
		return true
	}

	return false
}

func orientation(p, q, r Coord) int {
	val := (q.Y - p.Y) * (r.X - q.X) - (q.X - p.X) * (r.Y - q.Y)
	//colinear
	if val == 0 {
		return 0
	} else if val > 0 { //clockwise
		return 1
	} else { //counterclockwise
		return 2
	}
}

func onSegment(p, q, r Coord) bool {
	if q.X <= max(p.X, r.X) && q.X >= min(p.X, r.X) && q.Y <= max(p.Y, r.Y) && q.Y >= min(p.Y, r.Y){
		return true
	}

	return false
}

func min(num1, num2 int) int{
	if num1 < num2 {
		return num1
	} else {
		return  num2
	}
}

func max(num1, num2 int) int{
	if num1 > num2 {
		return num1
	} else {
		return  num2
	}
}

func (s *FusionCenter) CheckFalsePosWind(n *NodeImpl) int {
	sumX := 0
	sumY := 0

	if len(s.P.WindRegion[s.P.TimeStep]) == 0{
		return -1
	}
	for i:= 0; i < len(s.P.WindRegion[s.P.TimeStep]); i++ {
		sumX += transformX(s.P.WindRegion[s.P.TimeStep][i].X, s.P)  //get transformed X
		sumY += transformY(s.P.WindRegion[s.P.TimeStep][i].Y, s.P)  //get transformed Y
	}

	//calculate the center
	center := Coord{X: sumX/len(s.P.WindRegion[s.P.TimeStep]), Y: sumY/len(s.P.WindRegion[s.P.TimeStep])}
	tmp := s.P.WindRegion[s.P.TimeStep][0]
	for i:= 1; i < len(s.P.WindRegion[s.P.TimeStep]); i++ {
		//fmt.Printf("Checking if [ %v, %v ] intersects with [ %v, %v ]\n", tmp, s.P.WindRegion[s.P.TimeStep][i], center, n.GetLocCoord())
		if Intersects(transformCoord(tmp, s.P), transformCoord(s.P.WindRegion[s.P.TimeStep][i], s.P), transformCoord(center, s.P), transformCoord(n.GetLocCoord(), s.P)) {
			//fmt.Println("True!")
			return 1
		}
		tmp = s.P.WindRegion[s.P.TimeStep][i]
	}
	//fmt.Printf("Checking if [ %v, %v ] intersects with [ %v, %v ]\n", tmp, s.P.WindRegion[s.P.TimeStep][0], center, n.GetLocCoord())
	if Intersects(transformCoord(tmp, s.P), transformCoord(s.P.WindRegion[s.P.TimeStep][0], s.P), transformCoord(center, s.P), transformCoord(n.GetLocCoord(), s.P)) {
		//fmt.Println("True!")
		return 1
	}
	return 0
}

func (s *FusionCenter) UpdateClusterInfo(node *NodeImpl, rd *Reading) {
	if node.IsAlive() {
		if node.IsClusterHead {
			if len(node.StoredReadings) > 0 {
				if _, ok := s.Clusters[node]; !ok {
					s.ClearServerClusterInfo(node)
					s.Clusters[node] = &Cluster{Members: make(map[*NodeImpl]int), TimeFormed: rd.Time}
				}
				for _, reading := range node.StoredReadings {
					member := node.P.NodeList[reading.Id]
					if member.IsAlive() {
						if _, ok := s.Clusters[node].Members[member]; !ok {
							time, ok1 := s.AloneNodes[member]
							_, ok2 := s.Clusters[member]
							if (!ok1 || time < reading.Time) && (!ok2 || s.Clusters[member].TimeFormed < reading.Time) {
								delete(s.AloneNodes, member)
								if ok2 {
									for mem := range s.ClusterHeadsOf {
										s.ClusterHeadsOf[mem], _ = SearchRemove(s.ClusterHeadsOf[mem], member)
									}
									delete(s.Clusters, member)
								}
								s.Clusters[node].Members[member] = reading.Time
								s.ClusterHeadsOf[member] = append(s.ClusterHeadsOf[member], node)
								if len(s.ClusterHeadsOf[member]) > node.P.MaxClusterHeads {
									//Sort heads by oldest reading first
									sort.Slice(s.ClusterHeadsOf[member], func(i, j int) bool {
										head1 := s.ClusterHeadsOf[member][i]
										head2 := s.ClusterHeadsOf[member][j]
										time1 := s.Clusters[head1].Members[member]
										time2 := s.Clusters[head2].Members[member]
										return time1 < time2
									})
									oldestHead := s.ClusterHeadsOf[member][0]
									delete(s.Clusters[oldestHead].Members, member)
									s.ClusterHeadsOf[member] = s.ClusterHeadsOf[member][1:]
								}
							}
						} else if reading.Time > s.Clusters[node].Members[member] {
							s.Clusters[node].Members[member] = reading.Time
						}
					} else {
						s.ClearServerClusterInfo(member)
					}
				}
			} else {
				s.ClearServerClusterInfo(node)
				s.AloneNodes[node] = rd.Time
			}
		}
	} else {
		s.ClearServerClusterInfo(node)
	}
}

func (s *FusionCenter) ClearServerClusterInfo(node *NodeImpl) {
	for _, head := range s.ClusterHeadsOf[node] {
		delete(s.Clusters[head].Members, node)
	}
	if _, ok := s.Clusters[node]; ok {
		for mem := range s.Clusters[node].Members {
			s.ClusterHeadsOf[mem], _ = SearchRemove(s.ClusterHeadsOf[mem], node)
		}
	}
	delete(s.AloneNodes, node)
	delete(s.ClusterHeadsOf, node)
	delete(s.Clusters, node)
}

func (s *FusionCenter) CheckGlobalRecluster(nodesAccountedFor int) {
	ratio := 0.0
	if s.P.AloneOrClusterRatio {
		ratio = float64(len(s.AloneNodes)) / float64(nodesAccountedFor)
	} else {
		ratio = (float64(len(s.Clusters)) + float64(len(s.AloneNodes))) / float64(nodesAccountedFor)
	}
	if (s.P.GlobalRecluster == 1 && ratio > s.P.AloneThreshold) || (s.P.GlobalRecluster == 2 && s.P.CurrentTime/1000 > s.NextReclusterTime) || (s.P.GlobalRecluster == 3 && s.P.CurrentTime/1000 > s.NextReclusterTime && ratio > s.P.AloneThreshold) {
		s.RatioBeforeRecluster = ratio
		s.AloneNodes = make(map[*NodeImpl]int)
		s.Clusters = make(map[*NodeImpl]*Cluster)
		s.ClusterHeadsOf = make(map[*NodeImpl][]*NodeImpl)
		s.Waiting = true
		s.P.ClusterNetwork.FullRecluster(s.P)
		if s.P.GlobalRecluster == 3 {
			s.NextReclusterTime = s.P.CurrentTime/1000 + int(s.P.ReclusterPeriod)
		}
	}
}

func (s *FusionCenter) UpdateReclusterThresholds(nodesAccountedFor int) {
	ratio := 0.0
	if s.P.AloneOrClusterRatio {
		ratio = float64(len(s.AloneNodes)) / float64(nodesAccountedFor)
	} else {
		ratio = (float64(len(s.Clusters)) + float64(len(s.AloneNodes))) / float64(nodesAccountedFor)
	}
	improvement := (s.RatioBeforeRecluster - ratio) / s.RatioBeforeRecluster
	if improvement < s.P.SmallImprovement {
		s.Increments++
		if s.P.GlobalRecluster == 2 {
			s.P.ReclusterPeriod *= s.P.GlobalReclusterIncrement
			s.NextReclusterTime += int(s.P.ReclusterPeriod)
		} else {
			s.P.AloneThreshold *= s.P.GlobalReclusterIncrement
		}
	} else if improvement > s.P.LargeImprovement {
		s.Decrements++
		if s.P.GlobalRecluster == 2 {
			s.P.ReclusterPeriod *= s.P.GlobalReclusterDecrement
			s.NextReclusterTime += int(s.P.ReclusterPeriod)
		} else {
			s.P.AloneThreshold *= s.P.GlobalReclusterDecrement
		}
	}
}