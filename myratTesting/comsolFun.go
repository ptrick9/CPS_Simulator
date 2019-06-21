/*
-inputFileName=Scenario_3.txt
-imageFileName=marathon_street_map.png
-logPosition=true
-logGrid=true
-logEnergy=true
-logNodes=false
-noEnergy=true
-sensorPath=C:/Users/patrick/Dropbox/Patrick/udel/SUMMER2019/GitSimulator/smoothed_marathon.csv
-SquareRowCM=60
-SquareColCM=320
*/

package main

import (
	//"./cps"
	"../simulator/cps"
	"bytes"
	"flag"
	"fmt"
	"log"
	"math"
	"math/rand"
	"os"
	"runtime"
	"runtime/pprof"
	"time"

	"image"
	"image/png"
)

var (
	p *cps.Params
	r *cps.RegionParams

	err error

	// End the command line variables.

)

func init() {
	image.RegisterFormat("png", "png", png.Decode, png.DecodeConfig)
}

func main() {

	p = &cps.Params{}
	r = &cps.RegionParams{}

	p.Server = cps.FusionCenter{p, nil, nil, nil, nil, nil, nil}

	p.Tau1 = 10
	p.Tau2 = 500
	p.FoundBomb = false

	rand.Seed(time.Now().UTC().UnixNano())

	getFlags()

	if p.CPUProfile != "" {
		f, err := os.Create(p.CPUProfile)
		if err != nil {
			log.Fatal("could not create CPU profile: ", err)
		}
		defer f.Close()
		if err := pprof.StartCPUProfile(f); err != nil {
			log.Fatal("could not start CPU profile: ", err)
		}
		defer pprof.StopCPUProfile()
	}

	//fileName = p.InputFileNameCM
	p.FileName = p.InputFileNameCM
	cps.GetListedInput(p)
	//os.Exit(0)

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

	cps.MakeBoolGrid(p)

	p.Server.Init()
	cps.ReadMap(p, r)

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
	p.ThreshHoldBatteryToHave = 30.0 //This is the threshold battery to have for all phones

	p.Iterations_used = 0
	p.Iterations_of_event = 200
	p.EstimatedPingsNeeded = 10200

	cps.SetupFiles(p)

	cps.SetupParameters(p)

	//Printing important information to the p.Grid log file
	//fmt.Fprintln(p.GridFile, "Grid:", p.SquareRowCM, "x", p.SquareColCM)
	//fmt.Fprintln(p.GridFile, "Total Number of Nodes:", (p.NumNodes + numSuperNodes))
	//fmt.Fprintln(p.GridFile, "Runs:", iterations_of_event)

	fmt.Println("xDiv is ", p.XDiv, " yDiv is ", p.YDiv, " square capacity is ", p.SquareCapacity)

	p.Server.MakeSuperNodes()

	p.B = &cps.Bomb{X: p.BombX, Y: p.BombY}

	p.WallNodeList = make([]cps.WallNodes, p.NumWallNodes)

	p.NodeList = make([]cps.NodeImpl, 0)

	for i := 0; i < p.NumWallNodes; i++ {
		p.WallNodeList[i] = cps.WallNodes{Node: &cps.NodeImpl{X: p.Wpos[i][0], Y: p.Wpos[i][1]}}
	}

	p.Server.MakeGrid()

	fmt.Println("Super Node Type", p.SuperNodeType)
	fmt.Println("Dimensions: ", p.MaxX, "x", p.MaxY)
	//fmt.Println(p.SensorTimes)

	//This function initializes the super nodes in the scheduler's SNodeList
	//scheduler.MakeSuperNodes(p)

	fmt.Printf("Running Simulator iteration %d\\%v", 0, p.Iterations_of_event)

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
		fmt.Printf("\rRunning Simulator iteration %d\\%v", iters, p.Iterations_of_event)
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

		//is square thread safe
		//var wg sync.WaitGroup
		//wg.Add(len(p.NodeList))
		fmt.Fprintln(p.MoveReadingsFile, "T=", p.Iterations_used)
		for i := 0; i < len(p.NodeList); i++ {
			//go func(i int) {
			//	defer wg.Done()
				if !p.NoEnergyModelCM {
					p.NodeList[i].BatteryLossMostDynamic(p)
				} else {
					p.NodeList[i].HasCheckedSensor = true
					p.NodeList[i].Sitting = 0
				}
				p.NodeList[i].GetReadingsCSV(p)
				p.NodeList[i].GetReadings(p)
			//}(i)
		}
		//wg.Wait()
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
		//optimize := false

		p.Server.Tick()

		//Adding random points that the supernodes must visit
		if (iters%10 == 0) && (iters <= 990) {
			//fmt.Println(p.SuperNodeType)
			//fmt.Println(p.SuperNodeVariation)
			//scheduler.addRoutePoint(Coord{nil, rangeInt(0, p.MaxX), ranpositionPrintgeInt(0, p.MaxY), 0, 0, 0, 0})
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
		p.Server.CalcStats()
	}

	p.PositionFile.Seek(0, 0)
	fmt.Fprintln(p.PositionFile, "Width:", p.MaxX)
	fmt.Fprintln(p.PositionFile, "Height:", p.MaxY)
	fmt.Fprintf(p.PositionFile, "Amount: %-8v\n", iters)

	if iters < p.Iterations_of_event-1 {
		fmt.Printf("\nFound bomb at iteration: %v \nSimulation Complete\n", iters)
	} else {
		fmt.Println("\nSimulation Complete")
	}

	for i := range p.BoolGrid {
		fmt.Fprintln(p.BoolFile, p.BoolGrid[i])
	}

	if p.MemProfile != "" {
		f, err := os.Create(p.MemProfile)
		if err != nil {
			log.Fatal("could not create memory profile: ", err)
		}
		defer f.Close()
		runtime.GC() // get up-to-date statistics
		if err := pprof.WriteHeapProfile(f); err != nil {
			log.Fatal("could not write memory profile: ", err)
		}
	}
	p.Server.PrintStatsFile()
}


//MakeNodes creates all of the nodes in the simulator and stores them into NodeList
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

			p.BoolGrid[xx][yy] = true
		}
	}
}


//getFlags defines and sets commandline values
func getFlags() {
	//p = cps.Params{}

	flag.StringVar(&p.CPUProfile, "cpuprofile", "", "write cpu profile to `file`")
	flag.StringVar(&p.MemProfile, "memprofile", "", "write memory profile to `file`")

	//fmt.Println(os.Args[1:], "\nhmmm? \n ") //C:\Users\Nick\Desktop\comand line experiments\src
	flag.IntVar(&p.NegativeSittingStopThresholdCM, "negativeSittingStopThreshold", -10,
		"Negative number sitting is set to when board map is reset")
	flag.IntVar(&p.SittingStopThresholdCM, "sittingStopThreshold", 5,
		"How long it takes for a node to stay seated")
	flag.Float64Var(&p.GridCapacityPercentageCM, "GridCapacityPercentage", .9,
		"Percent the sub-Grid can be filled")
	flag.StringVar(&p.InputFileNameCM, "inputFileName", "Log1_in.txt",
		"Name of the input text file")
	flag.StringVar(&p.SensorPath, "sensorPath", "Circle_2D.csv", "Sensor Reading Inputs")
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
	flag.Float64Var(&p.DetectionThresholdCM, "detectionThreshold", 10000.0, //11180.0,
		"Value where if a node gets this reading or higher, it will trigger a detection")
	flag.Float64Var(&p.ErrorModifierCM, "errorMultiplier", 1.0,
		"Multiplier for error values in system")

	//Range 1, 2, or 4
	//Currently works for only a few numbers, can be easily expanded but is not currently dynamic
	flag.IntVar(&p.NumSuperNodes, "numSuperNodes", 4, "the number of super nodes in the simulator")
	flag.Float64Var(&p.CalibrationThresholdCM, "Recalibration Threshold", 3.0, "Value over grid average to recalibrate node")
	flag.Float64Var(&p.StdDevThresholdCM, "StandardDeviationThreshold", 1.7, "Detection Threshold based on standard deviations from mean")

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
	flag.IntVar(&p.SquareRowCM, "SquareRowCM", 50, "Number of rows of p.Grid squares, 1 through p.MaxX")
	flag.IntVar(&p.SquareColCM, "SquareColCM", 50, "Number of columns of p.Grid squares, 1 through p.MaxY")

	flag.StringVar(&p.ImageFileNameCM, "imageFileName", "circle_justWalls_x4.png", "Name of the input text file")
	flag.StringVar(&p.StimFileNameCM, "stimFileName", "circle_0.txt", "Name of the stimulus text file")
	flag.StringVar(&p.OutRoutingNameCM, "outRoutingName", "log.txt", "Name of the stimulus text file")
	flag.StringVar(&p.OutRoutingStatsNameCM, "outRoutingStatsName", "routingStats.txt", "Name of the output file for stats")

	flag.BoolVar(&p.RegionRouting, "regionRouting", true, "True if you want to use the new routing algorithm with regions and cutting")

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


//rangeInt returns a random number between max and min
func rangeInt(min, max int) int {
	return rand.Intn(max-min) + min
}

//printGrid saves the current measurements of each Square into a buffer to print into the file
func printGrid(g [][]*cps.Square) bytes.Buffer {
	var buffer bytes.Buffer
	/*for i := range g {
		for _, x := range g[i] {
			buffer.WriteString(fmt.Sprintf("%.2f\t", x.Avg))
		}
		buffer.WriteString(fmt.Sprintf("\n"))
	}
	return buffer*/
	for y := 0; y < len(g[0]); y++ {
		for x:=0; x < len(g); x++ {
			buffer.WriteString(fmt.Sprintf("%.2f\t", g[x][y].Avg))
		}
		buffer.WriteString(fmt.Sprintf("\n"))
	}
	return buffer
}

//printGridNodes saves the current p.NumNodes of each Square into a buffer to print to the file
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

//printSuperStats writes supernode data to a buffer
func printSuperStats(SNodeList []cps.SuperNodeParent) bytes.Buffer {
	var buffer bytes.Buffer
	for _, i := range SNodeList {
		buffer.WriteString(fmt.Sprintf("SuperNode: %d\t", i.GetId()))
		buffer.WriteString(fmt.Sprintf("SquaresMoved: %d\t", i.GetSquaresMoved()))
		buffer.WriteString(fmt.Sprintf("AvgResponseTime: %.2f\t", i.GetAvgResponseTime()))
	}
	return buffer
}