package main

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
	"time"
	"log"
	"io/ioutil"
	"path/filepath"
	"strings"
	"strconv"
	"bytes"
)

var (
	squareRow        int
	squareCol        int
	numStoredSamples int
	maxX             int
	maxY             int

	threshHoldBatteryToHave  float32
	totalPercentBatteryToUse float32
	iterations_used          int
	iterations_of_event      int
	b                        *bomb

	//How the griid is divided into rows and columns
	xDiv int
	yDiv int

	recalibrate    bool
	squareCapacity int
	boolGrid       [][]bool
	grid           [][]*Square

	detectionThreshold float64 //Value of sensor reading which determines a "detection" of a bomb

	// These are the command line variables.
	negativeSittingStopThresholdCM int     // This is a negative number for the sitting to be set to when map is reset
	sittingStopThresholdCM         int     // This is the threshold for the longest time a node can sit before no longer moving
	gridCapacityPercentageCM       float64 // This is the percent of a subgrid that can be filled with nodes, between 0.0 and 1.0
	errorModifierCM				   float64 // Multiplier for error model
	outputFileNameCM               string  // This is the prefix of the output text file
	inputFileNameCM                string  // This must be the name of the input text file with ".txt"
	naturalLossCM                  float64 // This can be any number n: 0 < n < .1
	sensorSamplingLossCM           float64 // This can be any number n: 0 < n < .1
	GPSSamplingLossCM              float64 // This can be any number n: 0 < n < GPSSamplingLossCM < .1
	serverSamplingLossCM           float64 // This can be any number n: 0 < n < serverSamplingLossCM < .1
	thresholdBatteryToHaveCM       int     // This can be any number n: 0 < n < 50
	thresholdBatteryToUseCM        int     // This can be any number n: 0 < n < 20 < 100-thresholdBatteryToHaveCM
	movementSamplingSpeedCM        int     // This can be any number n: 0 < n < 100
	movementSamplingPeriodCM       int     // This can be any int number n: 1 <= n <= 100
	maxBufferCapacityCM            int     // This can be aby int number n: 10 <= n <= 100
	energyModelCM                  string  // This can be "custom", "2StageServer", or other string will result in dynamic
	noEnergyModelCM				   bool    // If set to true, all energy model values ignored
	sensorSamplingPeriodCM         int     // This can be any int n: 1 <= n <= 100
	GPSSamplingPeriodCM            int     // This can be any int n: 1 <= n < sensorSamplingPeriodCM <=  100
	serverSamplingPeriodCM         int     // This can be any int n: 1 <= n < GPSSamplingPeriodCM <= 100
	numStoredSamplesCM             int     // This can be any int n: 5 <= n <= 25
	gridStoredSamplesCM            int     // This can be any int n: 5 <= n <= 25
	detectionThresholdCM           float64 //This is whatever value 1-1000 we determine should constitute a "detection" of a bomb
	positionPrintCM                bool    //This is either true or false for whether to print positions to log file
	energyPrintCM                  bool    //This is either true or false for whether to print energy info to log file
	nodesPrintCM                   bool    //This is either true or false for whether to print node readings/averages to log file
	gridPrintCM                    bool    //This is either true or false for whether to print grid readings to log file
	squareRowCM                    int     //This is an int 1 through maxX representing how many rows of squares there are
	squareColCM                    int     //This is an int 1 through maxY representing how many columns of squares there are

	numSuperNodes  int
	superNodeType  int
	superNodeSpeed int
	doOptimize     bool
	//superNodeVariation int
	superNodeRadius int

	center Coord

	positionPrint bool
	nodesPrint    bool

	driftFile    *os.File
	nodeFile     *os.File
	positionFile *os.File

	// End the command line variables.
)

const Tau1 = 10
const Tau2 = 500

type Tuple struct{
	x, y int
}

func main() {
	rand.Seed(time.Now().UTC().UnixNano())

	getFlags()

	maxX = 100
	maxY = 100
	squareRow = squareRowCM
	squareCol = squareColCM

	xDiv = maxX / squareCol
	yDiv = maxY / squareRow

	createBoard(maxX, maxY)

	roadFile, err := os.Create("roadLog.txt")
	if err != nil {
		log.Fatal("Cannot create file", err)
	}
	defer roadFile.Close()

	positionFile, err = os.Create("Log-simulatorOutput.txt")
	if err != nil {
		log.Fatal("Cannot create file", err)
	}
	defer positionFile.Close()

	fmt.Fprintln(positionFile, "Width:", maxX)
	fmt.Fprintln(positionFile, "Height:", maxY)
	fmt.Fprintf(positionFile, "Amount: %v\n", 0)
	fmt.Fprintf(positionFile, "Bomb x: %v\n", 0)
	fmt.Fprintf(positionFile, "Bomb y: %v\n", 0)

	doWalls := true
	if (doWalls) {

		//wallString := strconv.Itoa(49)
		routingName := "Testing Walls Output/Log-path-wall-maze.txt"
		wallName := "../CPS_Simulator/Testing Walls/out_initial_wall_maze.txt"

		absPath, _ := filepath.Abs(wallName)
		wallData, err := ioutil.ReadFile(absPath)
		if err != nil {
			log.Fatal("Cannot create file", err)
		}
		wallsWithHeader := string(wallData)

		walls := strings.Split(wallsWithHeader, "\n")

		walls = walls[3 : len(walls)-1]
		for i := 0; i < len(walls); i++ {
			line := strings.Split(walls[i], " ")
			x, _ := strconv.Atoi(line[1])
			y, _ := strconv.Atoi(line[3][:len(line[3])-1])

			boardMap[y][x] = -1
		}

		routingFile, err := os.Create(routingName)
		if err != nil {
			log.Fatal("Cannot create file", err)
		}
		defer routingFile.Close()

		//The scheduler determines which supernode should pursue a point of interest
		scheduler := &Scheduler{}

		//List of all the supernodes on the grid
		scheduler.sNodeList = make([]SuperNodeParent, numSuperNodes)

		for i := 0; i < numSuperNodes; i++ {
			snode_points := make([]Coord, 1)
			snode_path := make([]Coord, 0)
			all_points := make([]Coord, 0)

			//Defining the starting x and y values for the super node
			//This super node starts at the middle of the grid
			nodeCenter := Coord{x: 2, y: 2}
			x_val, y_val := nodeCenter.x, nodeCenter.y

			scheduler.sNodeList[i] = &sn_zero{&supern{&NodeImpl{x: x_val, y: y_val, id: i}, 1,
				1, superNodeRadius, superNodeRadius, 0, snode_points, snode_path,
				nodeCenter, 0, 0, 0, 0, 0, all_points}}

			//The super node's current location is always the first element in the routePoints list
			scheduler.sNodeList[i].updateLoc()
		}
		maxLength := -1

		for _, s := range scheduler.sNodeList {
			s.addRoutePoint(Coord{x: maxX - 3, y: 2})

			s.tick()

			//Writes the super node information to a file
			fmt.Fprint(routingFile, s)
			p := printPoints(s)
			fmt.Fprint(routingFile, " UnvisitedPoints: ")
			fmt.Fprintln(routingFile, p.String())

			routeLength := len(s.getRoutePath())
			if routeLength > maxLength {
				maxLength = routeLength
			}
		}

		fmt.Printf("Iteration %d/%v", 0, maxLength)
		for i := 0; i < (maxLength + 1); i++ {
			fmt.Printf("\rIteration %d/%v", i, maxLength)
			for _, s := range scheduler.sNodeList {
				s.tick()

				//Writes the super node information to a file
				fmt.Fprint(routingFile, s)
				p := printPoints(s)
				fmt.Fprint(routingFile, " UnvisitedPoints: ")
				fmt.Fprintln(routingFile, p.String())
			}
		}
	}
}

//This function allows the simulator to create a roadMap of the grid
//Every Coord in the grid is given an integer value corresponding to the
//	number of times the Coord is used by all paths
//The function first generates two random Coords on each half of the grid
//It then finds the path between those Coords
//It then increments the integer value of each Coord in the path by one
//This is done an amount of time to generate a conclusive distribution of paths
//	across the gird
//THe resulting roadMap is outputted to the file, first with the max number
//	if times a Coord is visited and then each Coord's integer value
func makeRoads(roadFile *os.File){
	//This map has Tuples as keys and integers as values
	//The Tuples represent the Coord in the grid and the integer represents
	//	the amount of times the Coord is visited by all paths
	roadMap := make(map[Tuple]int)

	//The max is kept track of the be outputted at the beginning of the
	//road output file
	//This is used to determine the gradient of color by the Viewer when
	//	displaying the roads
	max := 0

	aStarIterations := 100

	fmt.Printf("Running ASTAR iteration %d/%v",0, aStarIterations)
	for i:= 0; i < aStarIterations; i++{
		//Two Coords are randomly generated
		a := Coord{nil, rangeInt(0, maxX), rangeInt(0, maxY), 0, 0, 0, 0}
		b := Coord{nil, rangeInt(0, maxX), rangeInt(0, maxY), 0, 0, 0, 0}

		//The Coords' x and y positions are randomly updated to be on either the
		//	left and right side, or top and bottom of the grid
		//This is done to ensure the paths between these Coords crosses a large
		//	section of the grid
		if i %2 == 0{
			a.x = rangeInt(0, maxX/2)
			b.x = rangeInt(maxX/2, maxX)
		}else{
			a.y = rangeInt(0, maxY/2)
			b.y = rangeInt(maxY/2, maxY)
		}
		fmt.Printf("\rRunning ASTAR iteration %d/%v",i, aStarIterations)
		//The aStar path between these two Coords is created
		//Each Coord in this path is looped through and the integer value corresponding
		//	to that Coord is incremented by one
		for _, r := range aStar(a, b) {
			pos := Tuple{r.x, r.y}
			roadMap[pos]++
			if roadMap[pos] > max {
				max = roadMap[pos]
			}
		}
	}
	fmt.Fprintln(roadFile, "max", max)

	//This loops through the roadMap and outputs the integer value for every Coord
	//	in the grid
	for i := 0; i < maxX; i++ {
		for j := 0; j < maxY; j++ {
			fmt.Println("\rOutputting to roadLog: Coord", j, i)
			if boardMap[i][j] == -1 {
				fmt.Fprintln(roadFile, i, j, -1)
			}else{
				fmt.Fprintln(roadFile, i, j, roadMap[Tuple{j,i}])
			}
		}
	}
}

//Saves the Coords in the allPoints list into a buffer to
//	print to the file
func printPoints(s SuperNodeParent) bytes.Buffer {
	var buffer bytes.Buffer
	buffer.WriteString((fmt.Sprintf("[")))
	for ind, i := range s.getAllPoints() {
		buffer.WriteString(i.String())

		if ind != len(s.getAllPoints())-1 {
			buffer.WriteString(" ")
		}
	}
	buffer.WriteString((fmt.Sprintf("]")))
	return buffer
}

func getFlags() {
	//fmt.Println(os.Args[1:], "\nhmmm? \n ") //C:\Users\Nick\Desktop\comand line experiments\src
	flag.IntVar(&negativeSittingStopThresholdCM, "negativeSittingStopThreshold", -10,
		"Negative number sitting is set to when board map is reset")
	flag.IntVar(&sittingStopThresholdCM, "sittingStopThreshold", 5,
		"How long it takes for a node to stay seated")
	flag.Float64Var(&gridCapacityPercentageCM, "gridCapacityPercentage", .9,
		"Percent the sub-grid can be filled")
	flag.StringVar(&inputFileNameCM, "inputFileName", "Log1_in.txt",
		"Name of the input text file")
	flag.StringVar(&outputFileNameCM, "outputFileName", "Log",
		"Name of the output text file prefix")
	flag.Float64Var(&naturalLossCM, "naturalLoss", .005,
		"battery loss due to natural causes")
	flag.Float64Var(&sensorSamplingLossCM, "sensorSamplingLoss", .001,
		"battery loss due to sensor sampling")
	flag.Float64Var(&GPSSamplingLossCM, "GPSSamplingLoss", .005,
		"battery loss due to GPS sampling")
	flag.Float64Var(&serverSamplingLossCM, "serverSamplingLoss", .01,
		"battery loss due to server sampling")
	flag.IntVar(&thresholdBatteryToHaveCM, "thresholdBatteryToHave", 30,
		"Threshold battery phones should have")
	flag.IntVar(&thresholdBatteryToUseCM, "thresholdBatteryToUse", 10,
		"Threshold of battery phones should consume from all forms of sampling")
	flag.IntVar(&movementSamplingSpeedCM, "movementSamplingSpeed", 20,
		"the threshold of speed to increase sampling rate")
	flag.IntVar(&movementSamplingPeriodCM, "movementSamplingPeriod", 1,
		"the threshold of speed to increase sampling rate")
	flag.IntVar(&maxBufferCapacityCM, "maxBufferCapacity", 25,
		"maximum capacity for the buffer before it sends data to the server")
	flag.StringVar(&energyModelCM, "energyModel", "variable",
		"this determines the energy loss model that will be used")
	flag.BoolVar(&noEnergyModelCM, "noEnergy", false,
		"Whether or not to ignore energy model for simulation")
	flag.IntVar(&sensorSamplingPeriodCM, "sensorSamplingPeriod", 1000,
		"rate of the sensor sampling period when custom energy model is chosen")
	flag.IntVar(&GPSSamplingPeriodCM, "GPSSamplingPeriod", 1000,
		"rate of the GridGPS sampling period when custom energy model is chosen")
	flag.IntVar(&serverSamplingPeriodCM, "serverSamplingPeriod", 1000,
		"rate of the server sampling period when custom energy model is chosen")
	flag.IntVar(&numStoredSamplesCM, "nodeStoredSamples", 10,
		"number of samples stored by individual nodes for averaging")
	flag.IntVar(&gridStoredSamplesCM, "gridStoredSamples", 10,
		"number of samples stored by grid squares for averaging")
	flag.Float64Var(&detectionThresholdCM, "detectionThreshold", 11180.0,
		"Value where if a node gets this reading or higher, it will trigger a detection")
	flag.Float64Var(&errorModifierCM, "errorMultiplier", 1.0,
		"Multiplier for error values in system")
	//Range 1, 2, or 4
	//Currently works for only a few numbers, can be easily expanded but is not currently dynamic
	flag.IntVar(&numSuperNodes, "numSuperNodes", 1, "the number of super nodes in the simulator")

	//Range: 0-2
	//0: default routing algorithm, points added onto the end of the path and routed to in that order
	//flag.IntVar(&superNodeType, "superNodeType", 0, "the type of super node used in the simulator")
	//better descriptions incoming
	//Range: 0-6
	//	0: default routing algorithm, points added onto the end of the path and routed to in that order
	//	1: sophisticated routing algorithm, begin in center, routed anywhere
	//	2: sophisticated routing algorithm, begin inside circles located in the corners, only routed inside circle
	//	3: sophisticated routing algorithm, begin inside circles located on the sides, only routed inside circle
	//	4: sophisticated routing algorithm, being inside large circles located in the corners, only routed inside circle
	//	5: sophisticated routing algorithm, begin inside regions, only routed inside region
	//	6: regional return trip routing algorithm, routed inside region based on most points
	//	7: regional return trip routing algorithm, routed inside region based on oldest point
	flag.IntVar(&superNodeType, "superNodeType", 6, "the type of super node used in the simulator")

	//Range: 0-...
	//Theoretically could be as high as possible
	//Realistically should remain around 10
	flag.IntVar(&superNodeSpeed, "superNodeSpeed", 1, "the speed of the super node")

	//Range: true/false
	//Tells the simulator whether or not to optimize the path of all the super nodes
	//Only works when multiple super nodes are active in the same area
	flag.BoolVar(&doOptimize, "doOptimize", false, "whether or not to optimize the simulator")

	//Range: 0-4
	//	0: begin in center, routed anywhere
	//	1: begin inside circles located in the corners, only routed inside circle
	//	2: begin inside circles located on the sides, only routed inside circle
	//	3: being inside large circles located in the corners, only routed inside circle
	//	4: begin inside regions, only routed inside region
	//Only used for super nodes of type 1
	//flag.IntVar(&superNodeVariation, "superNodeVariation", 3, "super nodes of type 1 have different variations")

	flag.BoolVar(&positionPrintCM, "logPosition", false, "Whether you want to write position info to a log file")
	flag.BoolVar(&gridPrintCM, "logGrid", false, "Whether you want to write grid info to a log file")
	flag.BoolVar(&energyPrintCM, "logEnergy", false, "Whether you want to write energy into to a log file")
	flag.BoolVar(&nodesPrintCM, "logNodes", false, "Whether you want to write node readings to a log file")
	flag.IntVar(&squareRowCM, "squareRow", 100, "Number of rows of grid squares, 1 through maxX")
	flag.IntVar(&squareColCM, "squareCol", 100, "Number of columns of grid squares, 1 through maxY")

	flag.Parse()
}

func rangeInt(min, max int) int { //returns a random number between max and min
	return rand.Intn(max-min) + min
}
