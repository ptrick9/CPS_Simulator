package main

import (
	"bytes"
	"flag"
	"fmt"
	"log"
	"math"
	"math/rand"
	"os"
	"sync"
	"time"
)

var (
	squareRow        int
	squareCol        int
	numNodes         int
	numStoredSamples int
	numGridSamples   int
	maxX             int
	maxY             int
	bombX            int
	bombY            int

	wallNodeList []wallNodes
	nodeList     []NodeImpl

	//s              server
	batteryCharges []float32
	batteryLosses  []float32

	batteryLossesCheckingSensorScalar []float32
	batteryLossesCheckingGPSScalar    []float32
	batteryLossesCheckingServerScalar []float32
	//checkingBoard [][]float32

	threshHoldBatteryToHave  float32
	totalPercentBatteryToUse float32
	iterations_used          int
	iterations_of_event      int
	estimatedPingsNeeded     int
	b                        *bomb

	//How the griid is divided into rows and columns
	xDiv int
	yDiv int

	recalibrate    bool
	squareCapacity int
	boolGrid       [][]bool
	numAtt         int
	attractions    []*attraction
	grid           [][]*Square
	bombSquare     *Square
	xLoc           int
	yLoc           int

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

	centerCoord Coord

	center Coord

	positionPrint bool
	energyPrint   bool
	nodesPrint    bool
	gridPrint     bool

	driftFile    *os.File
	nodeFile     *os.File
	positionFile *os.File
	foundBomb    bool
	err 		 error

	// End the command line variables.
)

const Tau1 = 10
const Tau2 = 500

func main() {

	foundBomb = false

	rand.Seed(time.Now().UTC().UnixNano())

	getFlags()

	fileName = inputFileNameCM

	getListedInput()

	//squareRow = getDashedInput("squareRow")
	//squareCol = getDashedInput("squareCol")
	numNodes = getDashedInput("numNodes")
	//numStoredSamples = getDashedInput("numStoredSamples")
	maxX = getDashedInput("maxX")
	maxY = getDashedInput("maxY")
	bombX = getDashedInput("bombX")
	bombY = getDashedInput("bombY")
	//numAtt = getDashedInput("numAtt")

	/*resultFile, err := os.Create("result.txt")
	if err != nil {
		log.Fatal("Cannot create file", err)
	}
	defer resultFile.Close()*/

	numStoredSamples = numStoredSamplesCM
	numGridSamples = gridStoredSamplesCM
	detectionThreshold = detectionThresholdCM
	positionPrint = positionPrintCM
	gridPrint = gridPrintCM
	energyPrint = energyPrintCM
	nodesPrint = nodesPrintCM
	squareRow = squareRowCM
	squareCol = squareColCM

	createBoard(maxX, maxY)

	fillInWallsToBoard()

	//This is where the text file reading ends
	Vn := make([]float64, 1000)
	for i := 0; i < 1000; i++ {
		Vn[i] = rand.NormFloat64() * -.5
	}

	if numNodes > numNodeNodes {
		for i := 0; i < numNodes-numNodeNodes; i++ {
			npos = append(npos, []int{rangeInt(1, maxX), rangeInt(1, maxY), 0})
		}
	}

	rand.Seed(time.Now().UnixNano()) //sets random to work properly by tying to to clock
	threshHoldBatteryToHave = 30.0   //This is the threshhold battery to have for all phones
	totalPercentBatteryToUse = float32(thresholdBatteryToUseCM)

	iterations_used = 0
	iterations_of_event = 5000
	estimatedPingsNeeded = 10200
	batteryCharges = getLinearBatteryValues(len(npos))
	batteryLosses = getLinearBatteryLossConstant(len(npos), float32(naturalLossCM))
	batteryLossesCheckingSensorScalar = getLinearBatteryLossConstant(len(npos), float32(sensorSamplingLossCM))
	batteryLossesCheckingGPSScalar = getLinearBatteryLossConstant(len(npos), float32(GPSSamplingLossCM))
	batteryLossesCheckingServerScalar = getLinearBatteryLossConstant(len(npos), float32(serverSamplingLossCM))
	attractions = make([]*attraction, numAtt)

	positionFile, err = os.Create(outputFileNameCM + "-simulatorOutput.txt")
	if err != nil {
		log.Fatal("Cannot create file", err)
	}
	defer positionFile.Close()

	driftFile, err = os.Create(outputFileNameCM + "-drift.txt")
	if err != nil {
		log.Fatal("Cannot create file", err)
	}
	defer driftFile.Close()

	gridFile, err := os.Create(outputFileNameCM + "-grid.txt")
	if err != nil {
		log.Fatal("Cannot create file", err)
	}
	defer gridFile.Close()

	nodeFile, err = os.Create(outputFileNameCM + "-node_reading.txt")
	if err != nil {
		log.Fatal("Cannot create file", err)
	}
	defer nodeFile.Close()

	energyFile, err := os.Create(outputFileNameCM + "-node.txt")
	if err != nil {
		log.Fatal("Cannot create file", err)
	}
	defer energyFile.Close()

	routingFile, err := os.Create(outputFileNameCM + "-path.txt")
	if err != nil {
		log.Fatal("Cannot create file", err)
	}
	defer routingFile.Close()

	boolFile, err := os.Create(outputFileNameCM + "-bool.txt")
	if err != nil {
		log.Fatal("Cannot create file", err)
	}
	defer boolFile.Close()

	attractionFile, err := os.Create(outputFileNameCM + "-attraction.txt")
	if err != nil {
		log.Fatal("Cannot create file", err)
	}
	defer attractionFile.Close()

	xDiv = maxX / squareCol
	yDiv = maxY / squareRow

	fmt.Fprintln(positionFile, "Width:", maxX)
	fmt.Fprintln(positionFile, "Height:", maxY)
	fmt.Fprintf(positionFile, "Amount: %-8v\n", iterations_of_event)
	fmt.Fprintf(positionFile, "Bomb x: %v\n", bombX)
	fmt.Fprintf(positionFile, "Bomb y: %v\n", bombY)

	//The capacity for a square should be equal to the area of the square
	//So we take the side length (xDiv) and square it
	squareCapacity = int(math.Pow(float64(xDiv), 2))
	fmt.Println("xDiv is ", xDiv, " yDiv is ", yDiv, " square capacity is ", squareCapacity)

	//Center of the grid
	center.x = maxX / 2
	center.y = maxY / 2

	//Printing important information to the grid log file
	//fmt.Fprintln(gridFile, "Grid:", squareRow, "x", squareCol)
	//fmt.Fprintln(gridFile, "Total Number of Nodes:", (numNodes + numSuperNodes))
	//fmt.Fprintln(gridFile, "Runs:", iterations_of_event)

	fmt.Fprintln(gridFile,"Width:", squareCol)
	fmt.Fprintln(gridFile,"Height:", squareRow)

	//Printing parameters to driftFile
	fmt.Fprintln(driftFile, "Number of Nodes:", numNodes)
	fmt.Fprintln(driftFile, "Rows:", squareRow)
	fmt.Fprintln(driftFile, "Columns:", squareCol)
	fmt.Fprintln(driftFile, "Samples Stored by Node:", numStoredSamples)
	fmt.Fprintln(driftFile, "Samples Stored by Grid:", numGridSamples)
	fmt.Fprintln(driftFile, "Width:", maxX)
	fmt.Fprintln(driftFile, "Height:", maxY)
	fmt.Fprintln(driftFile, "Bomb x:", bombX)
	fmt.Fprintln(driftFile, "Bomb y:", bombY)
	fmt.Fprintln(driftFile, "Iterations:", iterations_of_event)
	fmt.Fprintln(driftFile, "Size of Square:", xDiv, "x", yDiv)
	fmt.Fprintln(driftFile, "Detection Threshold:", detectionThreshold)
	fmt.Fprintln(driftFile, "Input File Name:", inputFileNameCM)
	fmt.Fprintln(driftFile, "Output File Name:", outputFileNameCM)
	fmt.Fprintln(driftFile, "Battery Natural Loss:", naturalLossCM)
	fmt.Fprintln(driftFile, "Sensor Loss:", sensorSamplingLossCM, "\nGPS Loss:", GPSSamplingLossCM, "\nServer Loss:", serverSamplingLossCM)
	fmt.Fprintln(driftFile, "Printing Position:", positionPrint, "\nPrinting Energy:", energyPrint, "\nPrinting Nodes:", nodesPrint)
	fmt.Fprintln(driftFile, "Super Nodes:", numSuperNodes, "\nSuper Node Type:", superNodeType, "\nSuper Node Speed:", superNodeSpeed, "\nSuper Node Radius:", superNodeRadius)
	fmt.Fprintln(driftFile, "Error Multiplier:", errorModifierCM)
	fmt.Fprintln(driftFile, "--------------------")
	//Initializing the size of the boolean field representing coordinates
	boolGrid = make([][]bool, maxY)
	for i := range boolGrid {
		boolGrid[i] = make([]bool, maxX)
	}
	//Initializing the boolean field with values of false
	for i := 0; i < maxY; i++ {
		for j := 0; j < maxX; j++ {
			boolGrid[i][j] = false
		}
	}

	//The scheduler determines which supernode should pursue a point of interest
	scheduler := &Scheduler{}

	//List of all the supernodes on the grid
	//Currently only one for testing
	scheduler.sNodeList = make([]SuperNodeParent, numSuperNodes)

	b = &bomb{x: bombX, y: bombY}

	wallNodeList = make([]wallNodes, numWallNodes)

	nodeList = make([]NodeImpl, 0)

	for i := 0; i < numWallNodes; i++ {
		wallNodeList[i] = wallNodes{node: &NodeImpl{x: wpos[i][0], y: wpos[i][1]}}
	}

	grid = make([][]*Square, squareRow) //this creates the grid and only works if row is same size as column
	for i := range grid {
		grid[i] = make([]*Square, squareCol)
	}

	for i := 0; i < squareRow; i++ {
		for j := 0; j < squareCol; j++ {

			travelList := make([]bool, 0)
			for k := 0; k < numSuperNodes; k++ {
				travelList = append(travelList, true)
			}

			grid[i][j] = &Square{i, j, 0.0, 0, make([]float32, numGridSamples),
				numGridSamples, 0.0, 0, 0, false,
				0.0, 0.0, false, travelList}
		}
	}

	fmt.Println("Super Node Type", superNodeType)
	fmt.Println("Dimensions: ", maxX, "x", maxY)

	//This function initializes the super nodes in the scheduler's sNodeList
	scheduler.makeSuperNodes()

	fmt.Printf("Running Simulator iteration %d\\%v",0, iterations_of_event)

	i := 0
	for i = 0; i < iterations_of_event && !foundBomb; i++ {

		makeNodes()
		//fmt.Println(iterations_used)
		fmt.Printf("\rRunning Simulator iteration %d\\%v",i, iterations_of_event)
		if positionPrint {
			fmt.Fprintln(positionFile, "t= ", iterations_used, " amount= ", len(nodeList))
		}
		for i := 0; i < len(poispos); i++ {
			if iterations_used == poispos[i][2] || iterations_used == poispos[i][3] {
				for i := 0; i < len(nodeList); i++ {
					nodeList[i].sitting = negativeSittingStopThresholdCM
				}
				createBoard(maxX, maxY)
				fillInWallsToBoard()
				fillInBufferCurrent()
				fillPointsToBoard()
				fillInMap()
				//writeBordMapToFile()
				i = len(poispos)
			}
		}

		//start := time.Now()
		var wg sync.WaitGroup
		wg.Add(len(nodeList))
		for i := 0; i < len(nodeList); i++ {
			go func(i int) {
				defer wg.Done()
				if !noEnergyModelCM {
					nodeList[i].batteryLossMostDynamic()
				} else {
					nodeList[i].hasCheckedSensor = true
					nodeList[i].sitting = 0
				}
				nodeList[i].getReadings()
			}(i)
		}
		wg.Wait()
		driftFile.Sync()
		nodeFile.Sync()
		positionFile.Sync()

		fmt.Fprintln(energyFile, "Amount:", len(nodeList))
		for j := 0; j < len(nodeList); j++ {

			oldX, oldY := nodeList[j].getLoc()
			boolGrid[oldY][oldX] = false //set the old spot false since the node will now move away

			//move the node to its new location
			nodeList[j].move()

			//set the new location in the boolean field to true
			newX, newY := nodeList[j].getLoc()
			boolGrid[newY][newX] = true

			//writes the node information to the file
			if energyPrint {
				fmt.Fprintln(energyFile, nodeList[j])
			}

			//Add the node into its new Square's numNodes
			//If the node hasn't left the square, that Square's numNodes will
			//remain the same after these calculations
		}

		fmt.Fprintln(routingFile, "Amount:", numSuperNodes)

		//Alerts the scheduler to redraw the paths of super nodes as efficiently
		// as possible
		//This should optimize the distances the super nodes have to travel as the
		//	longer the simulator runs the more inefficient the paths can become
		optimize := false

		//Loops through each super node and calls their tick function
		//The tick function does all the node maintenance specific to that
		//	type of super node including: updating the routePath, adding points
		// 	of interest to the super node, and moving the super node
		for _, s := range scheduler.sNodeList {
			//Saves the current length of the super node's list of routePoints
			//If a routePoint is reached by a super node the scheduler should
			// 	reorganize the paths
			length := len(s.getRoutePoints())

			//The super node executes it's per iteration code
			s.tick()

			//Compares the path lengths to decide if optimization is needed
			//Optimization will only be done if he optimization requirements are met
			//	AND if the simulator is currently in a mode that requests optimization

			if length != len(s.getRoutePoints()) {
				bombSquare := grid[b.y/yDiv][b.x/xDiv]
				sSquare := grid[s.getY()/yDiv][s.getX()/xDiv]
				grid[s.getY()/yDiv][s.getX()/xDiv].hasDetected = false

				bdist := float32(math.Pow(float64(math.Pow(float64(math.Abs(float64(s.getX())-float64(b.x))), 2)+math.Pow(float64(math.Abs(float64(s.getY())-float64(b.y))), 2)), .5))

				if bombSquare == sSquare || bdist < 8.0 {
					foundBomb = true
				} else {
					sSquare.reset()
				}

			}

			if length != len(s.getRoutePoints()) {
				optimize = doOptimize // true &&
			}

			//Writes the super node information to a file
			fmt.Fprint(routingFile, s)
			p := printPoints(s)
			fmt.Fprint(routingFile, " UnvisitedPoints: ")
			fmt.Fprintln(routingFile, p.String())
		}

		//Executes the optimization code if the optimize flag is true
		if optimize {
			//The scheduler optimizes the paths of each super node
			scheduler.optimize()
			//Resets the optimize flag
			optimize = false
		}

		//Adding random points that the supernodes must visit
		if (i%10 == 0) && (i <= 990) {
			//fmt.Println(superNodeType)
			//fmt.Println(superNodeVariation)
			//scheduler.addRoutePoint(Coord{nil, rangeInt(0, maxX), ranpositionPrintgeInt(0, maxY), 0, 0, 0, 0})
		}

		//Loop over every square in the grid once again
		for k := 0; k < squareRow; k++ {
			for z := 0; z < squareCol; z++ {
				bombSquare := grid[b.y/yDiv][b.x/xDiv]
				bs_y := float64(b.y / yDiv)
				bs_x := float64(b.x / xDiv)

				grid[k][z].stdDev = math.Sqrt(grid[k][z].getSquareValues() / float64(grid[k][z].numNodes-1))

				//check for false negatives/positives
				if grid[k][z].numNodes > 0 && float64(grid[k][z].avg) < detectionThreshold && bombSquare == grid[k][z] && !grid[k][z].hasDetected {
					//this is a grid false negative
					fmt.Fprintln(driftFile, "Grid False Negative Avg:", grid[k][z].avg, "Square Row:", k, "Square Column:", z, "Iteration:", i)
					grid[k][z].hasDetected = true
				}

				if float64(grid[k][z].avg) >= detectionThreshold && (math.Abs(bs_y-float64(k)) >= 1.1 && math.Abs(bs_x-float64(z)) >= 1.1) && !grid[k][z].hasDetected {
					//this is a false positive
					fmt.Fprintln(driftFile, "Grid False Positive Avg:", grid[k][z].avg, "Square Row:", k, "Square Column:", z, "Iteration:", i)
					//report to supernodes
					xLoc := (z * xDiv) + int(xDiv/2)
					yLoc := (k * yDiv) + int(yDiv/2)
					centerCoord = Coord{x: xLoc, y: yLoc}
					scheduler.addRoutePoint(centerCoord)
					grid[k][z].hasDetected = true
				}

				if float64(grid[k][z].avg) >= detectionThreshold && (math.Abs(bs_y-float64(k)) <= 1.1 && math.Abs(bs_x-float64(z)) <= 1.1) && !grid[k][z].hasDetected {
					//this is a true positive
					fmt.Fprintln(driftFile, "Grid True Positive Avg:", grid[k][z].avg, "Square Row:", k, "Square Column:", z, "Iteration:", i)
					//report to supernodes
					xLoc := (z * xDiv) + int(xDiv/2)
					yLoc := (k * yDiv) + int(yDiv/2)
					centerCoord = Coord{x: xLoc, y: yLoc}
					scheduler.addRoutePoint(centerCoord)
					grid[k][z].hasDetected = true
				}

				grid[k][z].setSquareValues(0)
				grid[k][z].numNodes = 0
			}
		}

		//printing to log files
		if gridPrint {
			x := printGrid(grid)
			for number := range attractions {
				fmt.Fprintln(attractionFile, attractions[number])
			}
			fmt.Fprint(attractionFile, "----------------\n")
			fmt.Fprintln(gridFile, x.String())
		}
		fmt.Fprint(driftFile, "----------------\n")
		if energyPrint {
			//fmt.Fprint(energyFile, "----------------\n")
		}
		fmt.Fprint(gridFile, "----------------\n")
		if nodesPrint {
			fmt.Fprint(nodeFile, "----------------\n")
		}

		iterations_used++
	}

	positionFile.Seek(0, 0)
	fmt.Fprintln(positionFile, "Width:", maxX)
	fmt.Fprintln(positionFile, "Height:", maxY)
	fmt.Fprintf(positionFile, "Amount: %-8v\n", i)

	if (i < iterations_of_event - 1) {
		fmt.Printf("\nFound bomb at iteration: %v \nSimulation Complete\n", i)
	} else {
		fmt.Println("\nSimulation Complete")
	}

	for i := range boolGrid {
		fmt.Fprintln(boolFile, boolGrid[i])
	}


}

func makeNodes() {
	for i := 0; i < len(npos); i++ {

		if iterations_used == npos[i][2] {

			var initHistory = make([]float32, numStoredSamples)

			nodeList = append(nodeList, NodeImpl{x: npos[i][0], y: npos[i][1], id: len(nodeList), sampleHistory: initHistory, concentration: 0,
				cascade: i, battery: batteryCharges[i], batteryLossScalar: batteryLosses[i],
				batteryLossCheckingSensorScalar: batteryLossesCheckingSensorScalar[i],
				batteryLossGPSScalar:            batteryLossesCheckingGPSScalar[i],
				batteryLossCheckingServerScalar: batteryLossesCheckingServerScalar[i]})

			nodeList[len(nodeList)-1].setConcentration(((1000) / (math.Pow((float64(nodeList[len(nodeList)-1].geoDist(*b))/0.2)*0.25, 1.5))))

			curNode := nodeList[len(nodeList)-1] //variable to keep track of current node being added

			//values to determine coefficients
			curNode.setS0(rand.Float64()*0.2 + 0.1)
			curNode.setS1(rand.Float64()*0.2 + 0.1)
			curNode.setS2(rand.Float64()*0.2 + 0.1)
			//values to determine error in coefficients
			s0, s1, s2 := curNode.getCoefficients()
			curNode.setE0(rand.Float64() * 0.1 * errorModifierCM * s0)
			curNode.setE1(rand.Float64() * 0.1 * errorModifierCM * s1)
			curNode.setE2(rand.Float64() * 0.1 * errorModifierCM * s2)
			//Values to determine error in exponents
			curNode.setET1(Tau1 * rand.Float64() * errorModifierCM * 0.05)
			curNode.setET2(Tau1 * rand.Float64() * errorModifierCM * 0.05)

			//set node time and initial sensitivity
			curNode.nodeTime = 0
			curNode.initialSensitivity = s0 + (s1)*math.Exp(-float64(curNode.nodeTime)/Tau1) + (s2)*math.Exp(-float64(curNode.nodeTime)/Tau2)
			curNode.sensitivity = curNode.initialSensitivity

			nodeList[len(nodeList)-1] = curNode

			boolGrid[npos[i][1]][npos[i][0]] = true
		}
	}
}

/*
func (scheduler *Scheduler) makeSuperNodes() {
	for i := 0; i < numSuperNodes; i++ {
		snode_points := make([]Coord, 1)
		snode_path := make([]Coord, 0)
		all_points := make([]Coord, 0)

		if superNodeType == 0 {

			//Defining the starting x and y values for the super node
			//This super node starts at the middle of the grid
			nodeCenter, x_val, y_val := makeCenter1(i)

			scheduler.sNodeList[i] = &sn_zero{&supern{&NodeImpl{x: x_val, y: y_val, id: numNodes + i}, 1,
				1, superNodeRadius, 0, snode_points, snode_path,
				nodeCenter, 0, 0,
				0, 0, 0, all_points}}
		} else if superNodeType == 1 {
			nodeCenter := Coord{}
			x_val := 0
			y_val := 0

			//makeCenter creates the Coord that represents the super node's center
			if superNodeVariation == 0 {
				nodeCenter, x_val, y_val = makeCenter1(i)
			} else if superNodeVariation == 1 || superNodeVariation == 4 {
				nodeCenter, x_val, y_val = makeCenter1_corners(i)
			} else if superNodeVariation == 2 {
				nodeCenter, x_val, y_val = makeCenter1_sides(i)
			} else if superNodeVariation == 3 {
				nodeCenter, x_val, y_val = makeCenter1_largeCorners(i)
			}

			scheduler.sNodeList[i] = &sn_one{&supern{&NodeImpl{x: x_val, y: y_val, id: numNodes + i}, 1,
				1, superNodeRadius, 0, snode_points, snode_path,
				nodeCenter, 0, 0,
				0, 0, 1, all_points}}
		} else if superNodeType == 2 {
			//makeRegionList initializes the regionList for this super node
			r_list := makeRegionList(i)

			//makeCenter creates the Coord that represents the super node's center
			nodeCenter, x_val, y_val := makeCenter2(i, r_list)

			//The useRegionList is just initialized to an empty list
			ur_list := make([]Region, 0)

			scheduler.sNodeList[i] = &sn_two{&supern{&NodeImpl{id: numNodes + i, x: x_val, y: y_val}, 1,
				1, superNodeRadius, 0, snode_points, snode_path,
				nodeCenter, 0, 0, 0, 0, 1,
				all_points}, r_list, ur_list}
		}

		//The super node's current location is always the first element in the routePoints list
		scheduler.sNodeList[i].updateLoc()
	}

}*/

func (scheduler *Scheduler) makeSuperNodes() {
	for i := 0; i < numSuperNodes; i++ {
		snode_points := make([]Coord, 1)
		snode_path := make([]Coord, 0)
		all_points := make([]Coord, 0)

		if superNodeType == 0 {

			//Defining the starting x and y values for the super node
			//This super node starts at the middle of the grid
			nodeCenter, x_val, y_val := makeCenter1(i)

			scheduler.sNodeList[i] = &sn_zero{&supern{&NodeImpl{x: x_val, y: y_val, id: i}, 1,
				1, superNodeRadius, superNodeRadius, 0, snode_points, snode_path,
				nodeCenter, 0, 0, 0, 0, 0, all_points}}
		} else if (superNodeType == 6) || (superNodeType == 7) {
			//makeRegionList initializes the regionList for this super node
			r_list := makeRegionList(i)

			//makeCenter creates the Coord that represents the super node's center
			nodeCenter, x_val, y_val := makeCenter2(i, r_list)

			//The useRegionList is just initialized to an empty list
			ur_list := make([]Region, 0)

			scheduler.sNodeList[i] = &sn_two{&supern{&NodeImpl{id: i, x: x_val, y: y_val}, 1,
				1, superNodeRadius, superNodeRadius, 0, snode_points,
				snode_path, nodeCenter, 0, 0, 0, 0,
				1, all_points}, r_list, ur_list}
		} else if (superNodeType >= 1) || (superNodeType <= 5) {
			nodeCenter := Coord{}
			x_val := 0
			y_val := 0
			xRad := 0
			yRad := 0

			//makeCenter creates the Coord that represents the super node's center
			if superNodeType == 1 {
				nodeCenter, x_val, y_val = makeCenter1(i)
			} else if superNodeType == 2 || superNodeType == 5 {
				nodeCenter, x_val, y_val, xRad, yRad = makeCenter1_corners(i)
			} else if superNodeType == 3 {
				nodeCenter, x_val, y_val, xRad, yRad = makeCenter1_sides(i)
			} else if superNodeType == 4 {
				nodeCenter, x_val, y_val, xRad, yRad = makeCenter1_largeCorners(i)
			}
			scheduler.sNodeList[i] = &sn_one{&supern{&NodeImpl{x: x_val, y: y_val, id: i}, 1,
				1, xRad, yRad, 0, snode_points, snode_path,
				nodeCenter, 0, 0,
				0, 0, 1, all_points}}
		}
		//The super node's current location is always the first element in the routePoints list
		scheduler.sNodeList[i].updateLoc()
	}
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
	flag.IntVar(&numSuperNodes, "numSuperNodes", 4, "the number of super nodes in the simulator")

	//Range: 0-2
	//0: default routing algorithm, points added onto the end of the path and routed to in that order
	//flag.IntVar(&superNodeType, "superNodeType", 0, "the type of super node used in the simulator")

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
	flag.IntVar(&superNodeSpeed, "superNodeSpeed", 3, "the speed of the super node")

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
	fmt.Println("Natural Loss: ", naturalLossCM)
	fmt.Println("Sensor Sampling Loss: ", sensorSamplingLossCM)
	fmt.Println("GPS sampling loss: ", GPSSamplingLossCM)
	fmt.Println("Server sampling loss", serverSamplingLossCM)
	fmt.Println("Threshold Battery to use: ", thresholdBatteryToUseCM)
	fmt.Println("Threshold battery to have: ", thresholdBatteryToHaveCM)
	fmt.Println("Moving speed for incresed sampling: ", movementSamplingSpeedCM)
	fmt.Println("Period of extra sampling due to high speed: ", movementSamplingPeriodCM)
	fmt.Println("Maximum size of buffer posible: ", maxBufferCapacityCM)
	fmt.Println("Energy model type:", energyModelCM)
	fmt.Println("Sensor Sampling Period:", sensorSamplingPeriodCM)
	fmt.Println("GPS Sampling Period:", GPSSamplingPeriodCM)
	fmt.Println("Server Sampling Period:", serverSamplingPeriodCM)
	fmt.Println("Number of Node Stored Samples:", numStoredSamplesCM)
	fmt.Println("Number of Grid Stored Samples:", gridStoredSamplesCM)
	fmt.Println("Detection Threshold:", detectionThresholdCM)

	//fmt.Println("tail:", flag.Args())
}

func rangeInt(min, max int) int { //returns a random number between max and min
	return rand.Intn(max-min) + min
}

//Saves the current measurements of each Square into a
//buffer to print into the file
func printGrid(g [][]*Square) bytes.Buffer {
	var buffer bytes.Buffer
	for i, _ := range g {
		for _, x := range g[i] {
			buffer.WriteString(fmt.Sprintf("%.2f\t", x.avg))
		}
		buffer.WriteString(fmt.Sprintf("\n"))
	}
	return buffer
}

//Saves the current numNodes of each Square into a buffer
//to print to the file
func printGridNodes(g [][]*Square) bytes.Buffer {
	var buffer bytes.Buffer
	for i, _ := range g {
		for _, x := range g[i] {
			buffer.WriteString(fmt.Sprintf("%d\t", x.numNodes))
		}
		buffer.WriteString(fmt.Sprintf("\n"))
	}
	return buffer
}

func printSuperStats(sNodeList []SuperNodeParent) bytes.Buffer {
	var buffer bytes.Buffer
	for _, i := range sNodeList {
		buffer.WriteString(fmt.Sprintf("SuperNode: %d\t", i.getId()))
		buffer.WriteString(fmt.Sprintf("SquaresMoved: %d\t", i.getSquaresMoved()))
		buffer.WriteString(fmt.Sprintf("AvgResponseTime: %.2f\t", i.getAvgResponseTime()))
	}
	return buffer
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
