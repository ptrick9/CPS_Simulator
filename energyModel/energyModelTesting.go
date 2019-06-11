
package main

import (
	//"./cps"
	"../simulator/cps"
	"bytes"
	"flag"
	"fmt"
	"math"
	"math/rand"
	"sync"
	"time"

	"image"
	"image/png"
)

var (
	p  *cps.Params
	r  *cps.RegionParams
	err 		 error
	// End the command line variables.
)

func init() {
	image.RegisterFormat("png", "png", png.Decode, png.DecodeConfig)
}

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
	p.NumNodes = cps.GetDashedInput("numNodes", p)
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


	cps.ReadMap(p, r)

	sum := 0
	for i := range p.BoardMap {
		for j := range p.BoardMap[i] {
			sum += p.BoardMap[i][j]
		}
	}
	fmt.Println(sum)

	top_left_corner := cps.Coord{X: 0, Y: 0}
	top_right_corner := cps.Coord{X: 0, Y: 0}
	bot_left_corner := cps.Coord{X: 0, Y: 0}
	bot_right_corner := cps.Coord{X: 0, Y: 0}

	tl_min := p.Height + p.Width
	tr_max := -1
	bl_max := -1
	br_max := -1

	for x := 0; x < p.Height; x++ {
		for y := 0; y < p.Width; y++ {
			if r.Point_dict[cps.Tuple{x, y}] {
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

	starting_locs := make([]cps.Coord, 4)
	starting_locs[0] = top_left_corner
	starting_locs[1] = top_right_corner
	starting_locs[2] = bot_left_corner
	starting_locs[3] = bot_right_corner





	cps.GenerateRouting(p, r)


	//This is where the text file reading ends
	Vn := make([]float64, 1000)
	for i := 0; i < 1000; i++ {
		Vn[i] = rand.NormFloat64() * -.5
	}

	if p.NumNodes > p.NumNodeNodes {
		for i := 0; i < p.NumNodes-p.NumNodeNodes; i++ {

			//p.Npos = append(p.Npos, []int{rangeInt(1, p.MaxX), rangeInt(1, p.MaxY), 0})
			p.Npos = append(p.Npos, []int{0, 0, 0})
		}
	}

	rand.Seed(time.Now().UnixNano()) //sets random to work properly by tying to to clock
	p.ThreshHoldBatteryToHave = 30.0   //This is the threshhold battery to have for all phones



	p.Iterations_used = 0
	p.Iterations_of_event = 1000
	p.EstimatedPingsNeeded = 10200



	cps.SetupFiles(p)



	cps.SetupParameters(p)




	//Printing important information to the p.Grid log file
	//fmt.Fprintln(p.GridFile, "Grid:", p.SquareRowCM, "x", p.SquareColCM)
	//fmt.Fprintln(p.GridFile, "Total Number of Nodes:", (p.NumNodes + numSuperNodes))
	//fmt.Fprintln(p.GridFile, "Runs:", iterations_of_event)


	fmt.Println("xDiv is ", p.XDiv, " yDiv is ", p.YDiv, " square capacity is ", p.SquareCapacity)



	//The scheduler determines which supernode should pursue a point of interest
	scheduler := &cps.Scheduler{}

	//List of all the supernodes on the grid
	scheduler.SNodeList = make([]cps.SuperNodeParent, p.NumSuperNodes)

	for i := 0; i < p.NumSuperNodes; i++ {
		snode_points := make([]cps.Coord, 1)
		snode_path := make([]cps.Coord, 0)
		all_points := make([]cps.Coord, 0)

		//Defining the starting x and y values for the super node
		//This super node starts at the middle of the grid
		x_val, y_val := starting_locs[i].X, starting_locs[i].Y
		nodeCenter := cps.Coord{X: x_val, Y: y_val}

		scheduler.SNodeList[i] = &cps.Sn_zero{&cps.Supern{&cps.NodeImpl{X: x_val, Y: y_val, Id: i}, 1,
			1, p.SuperNodeRadius, p.SuperNodeRadius, 0, snode_points, snode_path,
			nodeCenter, 0, 0, 0, 0, 0, all_points}}

		//The super node's current location is always the first element in the routePoints list
		scheduler.SNodeList[i].UpdateLoc()
	}

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
				0.0, 0.0, false, travelList, make(map[cps.Tuple]*cps.NodeImpl)}
			p.Grid[i][j].NodesInSquare = make(map[cps.Tuple]*cps.NodeImpl)
		}
	}

	fmt.Println("Super Node Type", p.SuperNodeType)
	fmt.Println("Dimensions: ", p.MaxX, "x", p.MaxY)
	fmt.Println(p.SensorTimes)

	//This function initializes the super nodes in the scheduler's SNodeList
	//scheduler.MakeSuperNodes(p)

	fmt.Printf("Running Simulator iteration %d\\%v",0, p.Iterations_of_event)


	iters := 0
	p.CurrTime = 0
	for iters = 0; iters < p.Iterations_of_event && !p.FoundBomb; iters++ {

		for i := 0; i < len(p.SensorTimes); i++ {
			if p.Iterations_used == p.SensorTimes[i] {
				p.CurrTime = i
			}
		}
		fmt.Printf("Current time: %d\n", p.CurrTime)


		makeNodes()
		//fmt.Println(iterations_used)
		fmt.Printf("\rRunning Simulator iteration %d\\%v",iters, p.Iterations_of_event)
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
				//p.NodeList[i].GetReadingsCSV(p)
				p.NodeList[i].GetReadings(p)
			}(i)
		}
		wg.Wait()
		p.DriftFile.Sync()
		p.NodeFile.Sync()
		p.PositionFile.Sync()

		fmt.Fprintln(p.EnergyFile, "Amount:", len(p.NodeList))
		cps.HandleMovement(p)

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
		if (iters%10 == 0) && (iters <= 990) {
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
					fmt.Fprintln(p.DriftFile, "Grid False Negative Avg:", p.Grid[k][z].Avg, "Square Row:", k, "Square Column:", z, "Iteration:", iters)
					p.Grid[k][z].HasDetected = false
				}

				if float64(p.Grid[k][z].Avg) >= p.DetectionThreshold && (math.Abs(bs_y-float64(k)) >= 1.1 && math.Abs(bs_x-float64(z)) >= 1.1) && !p.Grid[k][z].HasDetected {
					//this is a false positive
					fmt.Fprintln(p.DriftFile, "Grid False Positive Avg:", p.Grid[k][z].Avg, "Square Row:", k, "Square Column:", z, "Iteration:", iters)
					//report to supernodes
					xLoc := (z * p.XDiv) + int(p.XDiv/2)
					yLoc := (k * p.YDiv) + int(p.YDiv/2)
					p.CenterCoord = cps.Coord{X: xLoc, Y: yLoc}
					scheduler.AddRoutePoint(p.CenterCoord, p, r)
					p.Grid[k][z].HasDetected = true
				}

				if float64(p.Grid[k][z].Avg) >= p.DetectionThreshold && (math.Abs(bs_y-float64(k)) <= 1.1 && math.Abs(bs_x-float64(z)) <= 1.1) && !p.Grid[k][z].HasDetected {
					//this is a true positive
					fmt.Fprintln(p.DriftFile, "Grid True Positive Avg:", p.Grid[k][z].Avg, "Square Row:", k, "Square Column:", z, "Iteration:", iters)
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
				fmt.Fprintln(p.AttractionFile, p.Attractions[number])
			}
			fmt.Fprint(p.AttractionFile, "----------------\n")
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
	fmt.Fprintf(p.PositionFile, "Amount: %-8v\n", iters)

	if (iters < p.Iterations_of_event - 1) {
		fmt.Printf("\nFound bomb at iteration: %v \nSimulation Complete\n", iters)
	} else {
		fmt.Println("\nSimulation Complete")
	}

	for i := range p.BoolGrid {
		fmt.Fprintln(p.BoolFile, p.BoolGrid[i])
	}
}

func makeNodes() {
	for i := 0; i < len(p.Npos); i++ {

		if p.Iterations_used == p.Npos[i][2] {

			var initHistory = make([]float32, p.NumStoredSamples)

			xx := rangeInt(1, p.MaxX)
			yy := rangeInt(1, p.MaxY)
			for p.BoolGrid[xx][yy] == true {
				xx = rangeInt(1, p.MaxX)
				yy = rangeInt(1, p.MaxY)
			}

			p.NodeList = append(p.NodeList, cps.NodeImpl{X: xx, Y: yy, Id: len(p.NodeList), SampleHistory: initHistory, Concentration: 0,
				Cascade: i, Battery: p.BatteryCharges[i], BatteryLossScalar: p.BatteryLosses[i],
				BatteryLossSensor: 				 p.BatteryLossesSensor[i],
				BatteryLossGPS:			         p.BatteryLossesGPS[i],
				BatteryLossServer:				 p.BatteryLossesServer[i],
				BatteryLossBT:					 p.BatteryLossesBT[i],
				BatteryLossWifi:				 p.BatteryLossesWiFi[i],
				BatteryLoss4G:					 p.BatteryLosses4G[i],
				BatteryLossAccelerometer:		 p.BatteryLossesAccelerometer[i]})

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

			p.BoolGrid[xx][yy] = true

			p.NodePositionMap[cps.Tuple{xx,yy}] = &p.NodeList[len(p.NodeList)-1]; //add Node to the position Map
			var row = xx/p.XDiv 	//scale down based off number of squares per row
			var col = yy/p.YDiv	//scale down based off number of squares per column
			p.Grid[row][col].NodesInSquare[cps.Tuple{xx,yy}] = &p.NodeList[len(p.NodeList)-1]; //add node to NodesInSquare based off 8x8 squares
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
	flag.Float64Var(&p.GridCapacityPercentageCM, "GridCapacityPercentage", .9,
		"Percent the sub-Grid can be filled")
	flag.StringVar(&p.InputFileNameCM, "inputFileName", "Log1_in.txt",
		"Name of the input text file")
	flag.StringVar(&p.SensorPath, "Sensor Readings", "Circle_2D.csv", "Sensor Reading Inputs")
	flag.StringVar(&p.OutputFileNameCM, "p.OutputFileName", "Log",
		"Name of the output text file prefix")
	flag.Float64Var(&p.NaturalLossCM, "naturalLoss", .005,
		"battery loss due to natural causes")

	flag.Float64Var(&p.SamplingLossSensorCM, "SamplingLossSensorCM", .001,
		"battery loss due to sensor sampling")
	flag.Float64Var(&p.SamplingLossGPSCM, "SamplingLossGPSCM", .005,
		"battery loss due to GPS sampling")
	flag.Float64Var(&p.SamplingLossServerCM, "SamplingLossServerCM", .01,
		"battery loss due to server sampling")

	flag.Float64Var(&p.SamplingLossBTCM, "SamplingLossBTCM", .0001,
		"battery loss due to BlueTooth sampling")
	flag.Float64Var(&p.SamplingLossWifiCM, "SamplingLossWifiCM", .001,
		"battery loss due to WiFi sampling")
	flag.Float64Var(&p.SamplingLoss4GCM, "SamplingLoss4GCM", .005,
		"battery loss due to 4G sampling")
	flag.Float64Var(&p.SamplingLossAccelCM, "SamplingLossAccelCM", .001,
		"battery loss due to accelerometer sampling")

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
	flag.Float64Var(&p.DetectionThresholdCM, "detectionThreshold", 10000.0,//11180.0,
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
	flag.IntVar(&p.SuperNodeType, "p.SuperNodeType", 0, "the type of super node used in the simulator")

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
	flag.IntVar(&p.SquareRowCM, "p.SquareRowCM", 50, "Number of rows of p.Grid squares, 1 through p.MaxX")
	flag.IntVar(&p.SquareColCM, "p.SquareColCM", 50, "Number of columns of p.Grid squares, 1 through p.MaxY")

	flag.StringVar(&p.ImageFileNameCM, "imageFileName", "circle_justWalls_x4.png", "Name of the input text file")
	flag.StringVar(&p.StimFileNameCM, "stimFileName", "circle_0.txt", "Name of the stimulus text file")
	flag.StringVar(&p.OutRoutingNameCM, "outRoutingName", "log.txt", "Name of the stimulus text file")
	flag.StringVar(&p.OutRoutingStatsNameCM, "outRoutingStatsName", "routingStats.txt", "Name of the output file for stats")

	flag.BoolVar(&p.RegionRouting, "regionRouting", true, "True if you want to use the new routing algorithm with regions and cutting")


	flag.Parse()
	fmt.Println("Natural Loss: ", p.NaturalLossCM)
	fmt.Println("Sensor Sampling Loss: ", p.SamplingLossSensorCM)
	fmt.Println("GPS sampling loss: ", p.SamplingLossGPSCM)
	fmt.Println("Server sampling loss", p.SamplingLossServerCM)

	fmt.Println("BlueTooth sampling loss: ", p.SamplingLossBTCM)
	fmt.Println("WiFi sampling loss", p.SamplingLossWifiCM)
	fmt.Println("4G sampling loss: ", p.SamplingLoss4GCM)
	fmt.Println("Accelerometer sampling loss", p.SamplingLossAccelCM)

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

//Saves the current p.NumNodes of each Square into a buffer
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