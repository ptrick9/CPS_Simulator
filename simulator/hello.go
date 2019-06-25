package main

import (
	"CPS_Simulator/simulator/cps"
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
	p *cps.Params
	r *cps.RegionParams

	err error

	// End the command line variables.
)

func main() {

	p = &cps.Params{}
	r = &cps.RegionParams{}

	p.Tau1 = 10
	p.Tau2 = 500
	p.FoundBomb = false

	rand.Seed(time.Now().UTC().UnixNano())

	getFlags()

	//fileName = p.InputFileNameCM
	p.FileName = p.InputFileNameCM
	cps.GetListedInput(p)

	//p.SquareRowCM = getDashedInput("p.SquareRowCM")
	//p.SquareColCM = getDashedInput("p.SquareColCM")
	p.TotalNodes = cps.GetDashedInput("numNodes", p)
	//numStoredSamples = getDashedInput("numStoredSamples")
	p.MaxX = cps.GetDashedInput("maxX", p)
	p.MaxY = cps.GetDashedInput("maxY", p)
	p.BombX = cps.GetDashedInput("bombX", p)
	p.BombY = cps.GetDashedInput("bombY", p)
	//numAtt = getDashedInput("numAtt")

	/*resultFile, err := os.Create("result.txt")
	if err != nil {
		log.Fatal("Cannot create file", err)
	}
	defer resultFile.Close()*/

	p.NumStoredSamples = p.NumStoredSamplesCM
	p.NumGridSamples = p.GridStoredSamplesCM
	p.DetectionThreshold = p.DetectionThresholdCM
	p.PositionPrint = p.PositionPrintCM
	p.GridPrint = p.GridPrintCM
	p.EnergyPrint = p.EnergyPrintCM
	p.NodesPrint = p.NodesPrintCM
	p.SquareRowCM = p.SquareRowCM
	p.SquareColCM = p.SquareColCM

	cps.CreateBoard(p.MaxX, p.MaxY, p)

	cps.FillInWallsToBoard(p)

	//This is where the text file reading ends
	Vn := make([]float64, 1000)
	for i := 0; i < 1000; i++ {
		Vn[i] = rand.NormFloat64() * -.5
	}

	if p.TotalNodes > p.CurrentNodes {
		for i := 0; i < p.TotalNodes-p.CurrentNodes; i++ {
			p.NodeEntryTimes = append(p.NodeEntryTimes, []int{rangeInt(1, p.MaxX), rangeInt(1, p.MaxY), 0})
		}
	}

	rand.Seed(time.Now().UnixNano()) //sets random to work properly by tying to to clock
	p.ThreshHoldBatteryToHave = 30.0 //This is the threshhold battery to have for all phones
	p.TotalPercentBatteryToUse = float32(p.ThresholdBatteryToUseCM)

	p.Iterations_used = 0
	p.Iterations_of_event = 5000
	p.EstimatedPingsNeeded = 10200
	p.BatteryCharges = cps.GetLinearBatteryValues(len(p.NodeEntryTimes))
	p.BatteryLosses = cps.GetLinearBatteryLossConstant(len(p.NodeEntryTimes), float32(p.NaturalLossCM))
	p.BatteryLossesCheckingSensorScalar = cps.GetLinearBatteryLossConstant(len(p.NodeEntryTimes), float32(p.SensorSamplingLossCM))
	p.BatteryLossesCheckingGPSScalar = cps.GetLinearBatteryLossConstant(len(p.NodeEntryTimes), float32(p.GPSSamplingLossCM))
	p.BatteryLossesCheckingServerScalar = cps.GetLinearBatteryLossConstant(len(p.NodeEntryTimes), float32(p.ServerSamplingLossCM))
	p.Attractions = make([]*cps.Attraction, p.NumAtt)

	p.PositionFile, err = os.Create(p.OutputFileNameCM + "-simulatorOutput.txt")
	if err != nil {
		log.Fatal("Cannot create file", err)
	}
	defer p.PositionFile.Close()

	p.DriftFile, err = os.Create(p.OutputFileNameCM + "-drift.txt")
	if err != nil {
		log.Fatal("Cannot create file", err)
	}
	defer p.DriftFile.Close()

	p.GridFile, err = os.Create(p.OutputFileNameCM + "-p.Grid.txt")
	if err != nil {
		log.Fatal("Cannot create file", err)
	}
	defer p.GridFile.Close()

	p.NodeFile, err = os.Create(p.OutputFileNameCM + "-node_reading.txt")
	if err != nil {
		log.Fatal("Cannot create file", err)
	}
	defer p.NodeFile.Close()

	p.EnergyFile, err = os.Create(p.OutputFileNameCM + "-node.txt")
	if err != nil {
		log.Fatal("Cannot create file", err)
	}
	defer p.EnergyFile.Close()

	p.RoutingFile, err = os.Create(p.OutputFileNameCM + "-path.txt")
	if err != nil {
		log.Fatal("Cannot create file", err)
	}
	defer p.RoutingFile.Close()

	boolFile, err := os.Create(p.OutputFileNameCM + "-bool.txt")
	if err != nil {
		log.Fatal("Cannot create file", err)
	}
	defer boolFile.Close()

	attractionFile, err := os.Create(p.OutputFileNameCM + "-attraction.txt")
	if err != nil {
		log.Fatal("Cannot create file", err)
	}
	defer attractionFile.Close()

	p.XDiv = p.MaxX / p.SquareColCM
	p.YDiv = p.MaxY / p.SquareRowCM

	fmt.Fprintln(p.PositionFile, "Width:", p.MaxX)
	fmt.Fprintln(p.PositionFile, "Height:", p.MaxY)
	fmt.Fprintf(p.PositionFile, "Amount: %-8v\n", p.Iterations_of_event)
	fmt.Fprintf(p.PositionFile, "Bomb x: %v\n", p.BombX)
	fmt.Fprintf(p.PositionFile, "Bomb y: %v\n", p.BombY)

	//The capacity for a square should be equal to the area of the square
	//So we take the side length (xDiv) and square it
	p.SquareCapacity = int(math.Pow(float64(p.XDiv), 2))
	fmt.Println("xDiv is ", p.XDiv, " yDiv is ", p.YDiv, " square capacity is ", p.SquareCapacity)

	//Center of the p.Grid
	p.Center.X = p.MaxX / 2
	p.Center.Y = p.MaxY / 2

	//Printing important information to the p.Grid log file
	//fmt.Fprintln(p.GridFile, "Grid:", p.SquareRowCM, "x", p.SquareColCM)
	//fmt.Fprintln(p.GridFile, "Total Number of Nodes:", (p.TotalNodes + numSuperNodes))
	//fmt.Fprintln(p.GridFile, "Runs:", iterations_of_event)

	fmt.Fprintln(p.GridFile, "Width:", p.SquareColCM)
	fmt.Fprintln(p.GridFile, "Height:", p.SquareRowCM)

	//Printing parameters to driftFile
	fmt.Fprintln(p.DriftFile, "Number of Nodes:", p.TotalNodes)
	fmt.Fprintln(p.DriftFile, "Rows:", p.SquareRowCM)
	fmt.Fprintln(p.DriftFile, "Columns:", p.SquareColCM)
	fmt.Fprintln(p.DriftFile, "Samples Stored by Node:", p.NumStoredSamples)
	fmt.Fprintln(p.DriftFile, "Samples Stored by Grid:", p.NumGridSamples)
	fmt.Fprintln(p.DriftFile, "Width:", p.MaxX)
	fmt.Fprintln(p.DriftFile, "Height:", p.MaxY)
	fmt.Fprintln(p.DriftFile, "Bomb x:", p.BombX)
	fmt.Fprintln(p.DriftFile, "Bomb y:", p.BombY)
	fmt.Fprintln(p.DriftFile, "Iterations:", p.Iterations_of_event)
	fmt.Fprintln(p.DriftFile, "Size of Square:", p.XDiv, "x", p.YDiv)
	fmt.Fprintln(p.DriftFile, "Detection Threshold:", p.DetectionThreshold)
	fmt.Fprintln(p.DriftFile, "Input File Name:", p.InputFileNameCM)
	fmt.Fprintln(p.DriftFile, "Output File Name:", p.OutputFileNameCM)
	fmt.Fprintln(p.DriftFile, "Battery Natural Loss:", p.NaturalLossCM)
	fmt.Fprintln(p.DriftFile, "Sensor Loss:", p.SensorSamplingLossCM, "\nGPS Loss:", p.GPSSamplingLossCM, "\nServer Loss:", p.ServerSamplingLossCM)
	fmt.Fprintln(p.DriftFile, "Printing Position:", p.PositionPrint, "\nPrinting Energy:", p.EnergyPrint, "\nPrinting Nodes:", p.NodesPrint)
	fmt.Fprintln(p.DriftFile, "Super Nodes:", p.NumSuperNodes, "\nSuper Node Type:", p.SuperNodeType, "\nSuper Node Speed:", p.SuperNodeSpeed, "\nSuper Node Radius:", p.SuperNodeRadius)
	fmt.Fprintln(p.DriftFile, "Error Multiplier:", p.ErrorModifierCM)
	fmt.Fprintln(p.DriftFile, "--------------------")
	//Initializing the size of the boolean field representing coordinates
	p.BoolGrid = make([][]bool, p.MaxY)
	for i := range p.BoolGrid {
		p.BoolGrid[i] = make([]bool, p.MaxX)
	}
	//Initializing the boolean field with values of false
	for i := 0; i < p.MaxY; i++ {
		for j := 0; j < p.MaxX; j++ {
			p.BoolGrid[i][j] = false
		}
	}

	//The scheduler determines which supernode should pursue a point of interest
	scheduler := &cps.Scheduler{}

	//List of all the supernodes on the p.Grid
	//Currently only one for testing
	scheduler.SNodeList = make([]cps.SuperNodeParent, p.NumSuperNodes)

	p.B = &cps.Bomb{X: p.BombX, Y: p.BombY}

	p.WallNodeList = make([]cps.WallNodes, p.NumWallNodes)

	p.NodeList = make([]cps.NodeImpl, 0)

	for i := 0; i < p.NumWallNodes; i++ {
		p.WallNodeList[i] = cps.WallNodes{Node: &cps.NodeImpl{X: p.Wpos[i][0], Y: p.Wpos[i][1]}}
	}

	p.Grid = make([][]*cps.Square, p.SquareRowCM) //this creates the p.Grid and only works if row is same size as column
	for i := range p.Grid {
		p.Grid[i] = make([]*cps.Square, p.SquareColCM)
	}

	for i := 0; i < p.SquareRowCM; i++ {
		for j := 0; j < p.SquareColCM; j++ {

			travelList := make([]bool, 0)
			for k := 0; k < p.NumSuperNodes; k++ {
				travelList = append(travelList, true)
			}

			p.Grid[i][j] = &cps.Square{i, j, 0.0, 0, make([]float32, p.NumGridSamples),
				p.NumGridSamples, 0.0, 0, 0, false,
				0.0, 0.0, false, travelList, sync.Mutex{}}
		}
	}

	fmt.Println("Super Node Type", p.SuperNodeType)
	fmt.Println("Dimensions: ", p.MaxX, "x", p.MaxY)

	//This function initializes the super nodes in the scheduler's SNodeList
	scheduler.MakeSuperNodes(p)

	fmt.Printf("Running Simulator iteration %d\\%v", 0, p.Iterations_of_event)

	i := 0
	for i = 0; i < p.Iterations_of_event && !p.FoundBomb; i++ {

		makeNodes()
		//fmt.Println(iterations_used)
		fmt.Printf("\rRunning Simulator iteration %d\\%v", i, p.Iterations_of_event)
		if p.PositionPrint {
			fmt.Fprintln(p.PositionFile, "t= ", p.Iterations_used, " amount= ", len(p.NodeList))
		}
		for i := 0; i < len(p.Poispos); i++ {
			if p.Iterations_used == p.Poispos[i][2] || p.Iterations_used == p.Poispos[i][3] {
				for i := 0; i < len(p.NodeList); i++ {
					p.NodeList[i].Sitting = p.NegativeSittingStopThresholdCM
				}
				cps.CreateBoard(p.MaxX, p.MaxY, p)
				cps.FillInWallsToBoard(p)
				cps.FillInBufferCurrent(p)
				cps.FillPointsToBoard(p)
				cps.FillInMap(p)
				//writeBordMapToFile()
				i = len(p.Poispos)
			}
		}

		//start := time.Now()
		var wg sync.WaitGroup
		wg.Add(len(p.NodeList))
		for i := 0; i < len(p.NodeList); i++ {
			go func(i int) {
				defer wg.Done()
				if !p.NoEnergyModelCM {
					p.NodeList[i].BatteryLossMostDynamic(p)
				} else {
					p.NodeList[i].HasCheckedSensor = true
					p.NodeList[i].Sitting = 0
				}
				p.NodeList[i].GetReadings(p)
			}(i)
		}
		wg.Wait()
		p.DriftFile.Sync()
		p.NodeFile.Sync()
		p.PositionFile.Sync()

		fmt.Fprintln(p.EnergyFile, "Amount:", len(p.NodeList))
		for j := 0; j < len(p.NodeList); j++ {

			oldX, oldY := p.NodeList[j].GetLoc()
			p.BoolGrid[oldY][oldX] = false //set the old spot false since the node will now move away

			//move the node to its new location
			p.NodeList[j].Move(p)

			//set the new location in the boolean field to true
			newX, newY := p.NodeList[j].GetLoc()
			p.BoolGrid[newY][newX] = true

			//writes the node information to the file
			if p.EnergyPrint {
				fmt.Fprintln(p.EnergyFile, p.NodeList[j])
			}

			//Add the node into its new Square's p.TotalNodes
			//If the node hasn't left the square, that Square's p.TotalNodes will
			//remain the same after these calculations
		}

		fmt.Fprintln(p.RoutingFile, "Amount:", p.NumSuperNodes)

		//Alerts the scheduler to redraw the paths of super nodes as efficiently
		// as possible
		//This should optimize the distances the super nodes have to travel as the
		//	longer the simulator runs the more inefficient the paths can become
		optimize := false

		//Loops through each super node and calls their tick function
		//The tick function does all the node maintenance specific to that
		//	type of super node including: updating the routePath, adding points
		// 	of interest to the super node, and moving the super node
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
				bombSquare := p.Grid[p.B.Y/p.YDiv][p.B.X/p.XDiv]
				sSquare := p.Grid[s.GetY()/p.YDiv][s.GetX()/p.XDiv]
				p.Grid[s.GetY()/p.YDiv][s.GetX()/p.XDiv].HasDetected = false

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
			pp := printPoints(s)
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

		//Adding random points that the supernodes must visit
		if (i%10 == 0) && (i <= 990) {
			//fmt.Println(p.SuperNodeType)
			//fmt.Println(p.SuperNodeVariation)
			//scheduler.addRoutePoint(Coord{nil, rangeInt(0, p.MaxX), ranpositionPrintgeInt(0, p.MaxY), 0, 0, 0, 0})
		}

		//Loop over every square in the p.Grid once again
		for k := 0; k < p.SquareRowCM; k++ {
			for z := 0; z < p.SquareColCM; z++ {
				bombSquare := p.Grid[p.B.Y/p.YDiv][p.B.X/p.XDiv]
				bs_y := float64(p.B.Y / p.YDiv)
				bs_x := float64(p.B.X / p.XDiv)

				p.Grid[k][z].StdDev = math.Sqrt(p.Grid[k][z].GetSquareValues() / float64(p.Grid[k][z].NumNodes-1))

				//check for false negatives/positives
				if p.Grid[k][z].NumNodes > 0 && float64(p.Grid[k][z].Avg) < p.DetectionThreshold && bombSquare == p.Grid[k][z] && !p.Grid[k][z].HasDetected {
					//this is a p.Grid false negative
					fmt.Fprintln(p.DriftFile, "Grid False Negative Avg:", p.Grid[k][z].Avg, "Square Row:", k, "Square Column:", z, "Iteration:", i)
					p.Grid[k][z].HasDetected = true
				}

				if float64(p.Grid[k][z].Avg) >= p.DetectionThreshold && (math.Abs(bs_y-float64(k)) >= 1.1 && math.Abs(bs_x-float64(z)) >= 1.1) && !p.Grid[k][z].HasDetected {
					//this is a false positive
					fmt.Fprintln(p.DriftFile, "Grid False Positive Avg:", p.Grid[k][z].Avg, "Square Row:", k, "Square Column:", z, "Iteration:", i)
					//report to supernodes
					xLoc := (z * p.XDiv) + int(p.XDiv/2)
					yLoc := (k * p.YDiv) + int(p.YDiv/2)
					p.CenterCoord = cps.Coord{X: xLoc, Y: yLoc}
					scheduler.AddRoutePoint(p.CenterCoord, p, r)
					p.Grid[k][z].HasDetected = true
				}

				if float64(p.Grid[k][z].Avg) >= p.DetectionThreshold && (math.Abs(bs_y-float64(k)) <= 1.1 && math.Abs(bs_x-float64(z)) <= 1.1) && !p.Grid[k][z].HasDetected {
					//this is a true positive
					fmt.Fprintln(p.DriftFile, "Grid True Positive Avg:", p.Grid[k][z].Avg, "Square Row:", k, "Square Column:", z, "Iteration:", i)
					//report to supernodes
					xLoc := (z * p.XDiv) + int(p.XDiv/2)
					yLoc := (k * p.YDiv) + int(p.YDiv/2)
					p.CenterCoord = cps.Coord{X: xLoc, Y: yLoc}
					scheduler.AddRoutePoint(p.CenterCoord, p, r)
					p.Grid[k][z].HasDetected = true
				}

				p.Grid[k][z].SetSquareValues(0)
				p.Grid[k][z].NumNodes = 0
			}
		}

		//printing to log files
		if p.GridPrint {
			x := printGrid(p.Grid)
			for number := range p.Attractions {
				fmt.Fprintln(attractionFile, p.Attractions[number])
			}
			fmt.Fprint(attractionFile, "----------------\n")
			fmt.Fprintln(p.GridFile, x.String())
		}
		fmt.Fprint(p.DriftFile, "----------------\n")
		if p.EnergyPrint {
			//fmt.Fprint(energyFile, "----------------\n")
		}
		fmt.Fprint(p.GridFile, "----------------\n")
		if p.NodesPrint {
			fmt.Fprint(p.NodeFile, "----------------\n")
		}

		p.Iterations_used++
	}

	p.PositionFile.Seek(0, 0)
	fmt.Fprintln(p.PositionFile, "Width:", p.MaxX)
	fmt.Fprintln(p.PositionFile, "Height:", p.MaxY)
	fmt.Fprintf(p.PositionFile, "Amount: %-8v\n", i)

	if i < p.Iterations_of_event-1 {
		fmt.Printf("\nFound bomb at iteration: %v \nSimulation Complete\n", i)
	} else {
		fmt.Println("\nSimulation Complete")
	}

	for i := range p.BoolGrid {
		fmt.Fprintln(boolFile, p.BoolGrid[i])
	}

}

func makeNodes() {
	for i := 0; i < len(p.NodeEntryTimes); i++ {

		if p.Iterations_used == p.NodeEntryTimes[i][2] {

			var initHistory = make([]float32, p.NumStoredSamples)

			p.NodeList = append(p.NodeList, cps.NodeImpl{X: p.NodeEntryTimes[i][0], Y: p.NodeEntryTimes[i][1], Id: len(p.NodeList), SampleHistory: initHistory, Concentration: 0,
				Cascade: i, Battery: p.BatteryCharges[i], BatteryLossScalar: p.BatteryLosses[i],
				BatteryLossCheckingSensorScalar: p.BatteryLossesCheckingSensorScalar[i],
				BatteryLossGPSScalar:            p.BatteryLossesCheckingGPSScalar[i],
				BatteryLossCheckingServerScalar: p.BatteryLossesCheckingServerScalar[i]})

			p.NodeList[len(p.NodeList)-1].SetConcentration(((1000) / (math.Pow((float64(p.NodeList[len(p.NodeList)-1].GeoDist(*p.B))/0.2)*0.25, 1.5))))

			curNode := p.NodeList[len(p.NodeList)-1] //variable to keep track of current node being added

			//values to determine coefficients
			curNode.SetS0(rand.Float64()*0.2 + 0.1)
			curNode.SetS1(rand.Float64()*0.2 + 0.1)
			curNode.SetS2(rand.Float64()*0.2 + 0.1)
			//values to determine error in coefficients
			s0, s1, s2 := curNode.GetCoefficients()
			curNode.SetE0(rand.Float64() * 0.1 * p.ErrorModifierCM * s0)
			curNode.SetE1(rand.Float64() * 0.1 * p.ErrorModifierCM * s1)
			curNode.SetE2(rand.Float64() * 0.1 * p.ErrorModifierCM * s2)
			//Values to determine error in exponents
			curNode.SetET1(p.Tau1 * rand.Float64() * p.ErrorModifierCM * 0.05)
			curNode.SetET2(p.Tau1 * rand.Float64() * p.ErrorModifierCM * 0.05)

			//set node time and initial sensitivity
			curNode.NodeTime = 0
			curNode.InitialSensitivity = s0 + (s1)*math.Exp(-float64(curNode.NodeTime)/p.Tau1) + (s2)*math.Exp(-float64(curNode.NodeTime)/p.Tau2)
			curNode.Sensitivity = curNode.InitialSensitivity

			p.NodeList[len(p.NodeList)-1] = curNode

			p.BoolGrid[p.NodeEntryTimes[i][1]][p.NodeEntryTimes[i][0]] = true
		}
	}
}

func getFlags() {
	//p = cps.Params{}

	//fmt.Println(os.Args[1:], "\nhmmm? \n ") //C:\Users\Nick\Desktop\comand line experiments\src
	flag.IntVar(&p.NegativeSittingStopThresholdCM, "negativeSittingStopThreshold", -10,
		"Negative number sitting is set to when board map is reset")
	flag.IntVar(&p.SittingStopThresholdCM, "sittingStopThreshold", 5,
		"How long it takes for a node to stay seated")
	flag.Float64Var(&p.GridCapacityPercentageCM, "p.GridCapacityPercentage", .9,
		"Percent the sub-p.Grid can be filled")
	flag.StringVar(&p.InputFileNameCM, "inputFileName", "Log1_in.txt",
		"Name of the input text file")
	flag.StringVar(&p.OutputFileNameCM, "p.OutputFileName", "Log",
		"Name of the output text file prefix")
	flag.Float64Var(&p.NaturalLossCM, "naturalLoss", .005,
		"battery loss due to natural causes")
	flag.Float64Var(&p.SensorSamplingLossCM, "sensorSamplingLoss", .001,
		"battery loss due to sensor sampling")
	flag.Float64Var(&p.GPSSamplingLossCM, "GPSSamplingLoss", .005,
		"battery loss due to GPS sampling")
	flag.Float64Var(&p.ServerSamplingLossCM, "serverSamplingLoss", .01,
		"battery loss due to server sampling")
	flag.IntVar(&p.ThresholdBatteryToHaveCM, "thresholdBatteryToHave", 30,
		"Threshold battery phones should have")
	flag.IntVar(&p.ThresholdBatteryToUseCM, "thresholdBatteryToUse", 10,
		"Threshold of battery phones should consume from all forms of sampling")
	flag.IntVar(&p.MovementSamplingSpeedCM, "movementSamplingSpeed", 20,
		"the threshold of speed to increase sampling rate")
	flag.IntVar(&p.MovementSamplingPeriodCM, "movementSamplingPeriod", 1,
		"the threshold of speed to increase sampling rate")
	flag.IntVar(&p.MaxBufferCapacityCM, "maxBufferCapacity", 25,
		"maximum capacity for the buffer before it sends data to the server")
	flag.StringVar(&p.EnergyModelCM, "energyModel", "variable",
		"this determines the energy loss model that will be used")
	flag.BoolVar(&p.NoEnergyModelCM, "noEnergy", false,
		"Whether or not to ignore energy model for simulation")
	flag.IntVar(&p.SensorSamplingPeriodCM, "sensorSamplingPeriod", 1000,
		"rate of the sensor sampling period when custom energy model is chosen")
	flag.IntVar(&p.GPSSamplingPeriodCM, "GPSSamplingPeriod", 1000,
		"rate of the GridGPS sampling period when custom energy model is chosen")
	flag.IntVar(&p.ServerSamplingPeriodCM, "serverSamplingPeriod", 1000,
		"rate of the server sampling period when custom energy model is chosen")
	flag.IntVar(&p.NumStoredSamplesCM, "nodeStoredSamples", 10,
		"number of samples stored by individual nodes for averaging")
	flag.IntVar(&p.GridStoredSamplesCM, "p.GridStoredSamples", 10,
		"number of samples stored by p.Grid squares for averaging")
	flag.Float64Var(&p.DetectionThresholdCM, "detectionThreshold", 11180.0,
		"Value where if a node gets this reading or higher, it will trigger a detection")
	flag.Float64Var(&p.ErrorModifierCM, "errorMultiplier", 1.0,
		"Multiplier for error values in system")
	//Range 1, 2, or 4
	//Currently works for only a few numbers, can be easily expanded but is not currently dynamic
	flag.IntVar(&p.NumSuperNodes, "numSuperNodes", 4, "the number of super nodes in the simulator")

	//Range: 0-2
	//0: default routing algorithm, points added onto the end of the path and routed to in that order
	//flag.IntVar(&p.SuperNodeType, "p.SuperNodeType", 0, "the type of super node used in the simulator")
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
	flag.IntVar(&p.SuperNodeType, "p.SuperNodeType", 6, "the type of super node used in the simulator")

	//Range: 0-...
	//Theoretically could be as high as possible
	//Realistically should remain around 10
	flag.IntVar(&p.SuperNodeSpeed, "p.SuperNodeSpeed", 3, "the speed of the super node")

	//Range: true/false
	//Tells the simulator whether or not to optimize the path of all the super nodes
	//Only works when multiple super nodes are active in the same area
	flag.BoolVar(&p.DoOptimize, "doOptimize", false, "whether or not to optimize the simulator")

	//Range: 0-4
	//	0: begin in center, routed anywhere
	//	1: begin inside circles located in the corners, only routed inside circle
	//	2: begin inside circles located on the sides, only routed inside circle
	//	3: being inside large circles located in the corners, only routed inside circle
	//	4: begin inside regions, only routed inside region
	//Only used for super nodes of type 1
	//flag.IntVar(&p.SuperNodeVariation, "p.SuperNodeVariation", 3, "super nodes of type 1 have different variations")

	flag.BoolVar(&p.PositionPrintCM, "logPosition", false, "Whether you want to write position info to a log file")
	flag.BoolVar(&p.GridPrintCM, "logGrid", false, "Whether you want to write p.Grid info to a log file")
	flag.BoolVar(&p.EnergyPrintCM, "logEnergy", false, "Whether you want to write energy into to a log file")
	flag.BoolVar(&p.NodesPrintCM, "logNodes", false, "Whether you want to write node readings to a log file")
	flag.IntVar(&p.SquareRowCM, "p.SquareRowCM", 100, "Number of rows of p.Grid squares, 1 through p.MaxX")
	flag.IntVar(&p.SquareColCM, "p.SquareColCM", 100, "Number of columns of p.Grid squares, 1 through p.MaxY")

	flag.Parse()
	fmt.Println("Natural Loss: ", p.NaturalLossCM)
	fmt.Println("Sensor Sampling Loss: ", p.SensorSamplingLossCM)
	fmt.Println("GPS sampling loss: ", p.GPSSamplingLossCM)
	fmt.Println("Server sampling loss", p.ServerSamplingLossCM)
	fmt.Println("Threshold Battery to use: ", p.ThresholdBatteryToUseCM)
	fmt.Println("Threshold battery to have: ", p.ThresholdBatteryToHaveCM)
	fmt.Println("Moving speed for incresed sampling: ", p.MovementSamplingSpeedCM)
	fmt.Println("Period of extra sampling due to high speed: ", p.MovementSamplingPeriodCM)
	fmt.Println("Maximum size of buffer posible: ", p.MaxBufferCapacityCM)
	fmt.Println("Energy model type:", p.EnergyModelCM)
	fmt.Println("Sensor Sampling Period:", p.SensorSamplingPeriodCM)
	fmt.Println("GPS Sampling Period:", p.GPSSamplingPeriodCM)
	fmt.Println("Server Sampling Period:", p.ServerSamplingPeriodCM)
	fmt.Println("Number of Node Stored Samples:", p.NumStoredSamplesCM)
	fmt.Println("Number of Grid Stored Samples:", p.GridStoredSamplesCM)
	fmt.Println("Detection Threshold:", p.DetectionThresholdCM)

	//fmt.Println("tail:", flag.Args())
}

func rangeInt(min, max int) int { //returns a random number between max and min
	return rand.Intn(max-min) + min
}

//Saves the current measurements of each Square into a
//buffer to print into the file
func printGrid(g [][]*cps.Square) bytes.Buffer {
	var buffer bytes.Buffer
	for i, _ := range g {
		for _, x := range g[i] {
			buffer.WriteString(fmt.Sprintf("%.2f\t", x.Avg))
		}
		buffer.WriteString(fmt.Sprintf("\n"))
	}
	return buffer
}

//Saves the current p.TotalNodes of each Square into a buffer
//to print to the file
func printGridNodes(g [][]*cps.Square) bytes.Buffer {
	var buffer bytes.Buffer
	for i, _ := range g {
		for _, x := range g[i] {
			buffer.WriteString(fmt.Sprintf("%d\t", x.NumNodes))
		}
		buffer.WriteString(fmt.Sprintf("\n"))
	}
	return buffer
}

func printSuperStats(SNodeList []cps.SuperNodeParent) bytes.Buffer {
	var buffer bytes.Buffer
	for _, i := range SNodeList {
		buffer.WriteString(fmt.Sprintf("SuperNode: %d\t", i.GetId()))
		buffer.WriteString(fmt.Sprintf("SquaresMoved: %d\t", i.GetSquaresMoved()))
		buffer.WriteString(fmt.Sprintf("AvgResponseTime: %.2f\t", i.GetAvgResponseTime()))
	}
	return buffer
}

//Saves the Coords in the allPoints list into a buffer to
//	print to the file
func printPoints(s cps.SuperNodeParent) bytes.Buffer {
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
