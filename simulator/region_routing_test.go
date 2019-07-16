/*
-inputFileName=Scenario_3.txt
-imageFileName=routingCases.png
-logPosition=true
-logGrid=true
-logEnergy=true
-logNodes=false
-noEnergy=true
-sensorPath=smoothed_marathon.csv
-SquareRowCM=60
-SquareColCM=320
-csvMove=true
-movementPath=marathon_2k.txt
-iterations=1000
-csvSensor=false
-detectionThreshold=5
-numSuperNodes=4
*/
package main

import (
	"CPS_Simulator/simulator/cps"
	"flag"
	"fmt"
	"testing"
)

func TestRegions(t *testing.T) {
	p := &cps.Params{}
	r := &cps.RegionParams{}
	p.Server = cps.FusionCenter{p, r, nil, nil, nil, nil, nil, nil, nil}

	getFlagsTest(p)

	p.FileName = p.InputFileNameCM
	cps.GetListedInput(p)
	p.MaxX = cps.GetDashedInput("maxX", p)
	p.MaxY = cps.GetDashedInput("maxY", p)
	p.TotalNodes = cps.GetDashedInput("numNodes", p)
	p.BombX = cps.GetDashedInput("bombX", p)
	p.BombY = cps.GetDashedInput("bombY", p)

	p.B = &cps.Bomb{X: p.BombX, Y: p.BombY}

	p.NumStoredSamples = p.NumStoredSamplesCM
	p.NumGridSamples = p.GridStoredSamplesCM
	p.DetectionThreshold = p.DetectionThresholdCM
	p.PositionPrint = p.PositionPrintCM
	p.GridPrint = p.GridPrintCM
	p.NodesPrint = p.NodesPrintCM
	p.SquareRowCM = p.SquareRowCM
	p.SquareColCM = p.SquareColCM

	cps.MakeBoolGrid(p)
	p.Server.Init()
	cps.ReadMap(p, r)

	cps.SetupParameters(p)

	p.Server.MakeGrid()
	p.Server.MakeSuperNodes()

	cps.GenerateRouting(p, r)
	cps.FlipSquares(p, r)

	p.Iterations_used = 0
	p.Iterations_of_event = p.IterationsCM
	p.EstimatedPingsNeeded = 10200

	/*RegionVisual, err := os.Create(p.OutputFileNameCM + "-regionVisual.txt")
	if err != nil {
		log.Fatal("Cannot create file", err)
	}

	var buffer bytes.Buffer
	for i:= range r.Square_list {
		buffer.WriteString(fmt.Sprintf("X:%v Y:%v X:%v Y:%v\n", r.Square_list[i].X1, r.Square_list[i].Y1, r.Square_list[i].X2, r.Square_list[i].Y2))
	}
	fmt.Fprintf(RegionVisual, buffer.String())*/

	p.Server.Sch.SNodeList[0].SetLoc(cps.Coord{X:7, Y:26})
	p.Server.Sch.SNodeList[1].SetLoc(cps.Coord{X:17, Y:23})
	p.Server.Sch.SNodeList[2].SetLoc(cps.Coord{X:20, Y:14})
	p.Server.Sch.SNodeList[3].SetLoc(cps.Coord{X:31, Y:24})

	p.Server.Sch.SNodeList[4].SetLoc(cps.Coord{X:45, Y:7})
	p.Server.Sch.SNodeList[5].SetLoc(cps.Coord{X:45, Y:22})
	p.Server.Sch.SNodeList[6].SetLoc(cps.Coord{X:56, Y:26})
	p.Server.Sch.SNodeList[7].SetLoc(cps.Coord{X:32, Y:2})

	p.Server.Sch.AddRoutePoint(cps.Coord{X: 4, Y: 16})
	p.Server.Sch.AddRoutePoint(cps.Coord{X: 12, Y: 13})
	p.Server.Sch.AddRoutePoint(cps.Coord{X: 27, Y: 10})
	p.Server.Sch.AddRoutePoint(cps.Coord{X: 32, Y: 14})

	p.Server.Sch.AddRoutePoint(cps.Coord{X: 37, Y: 11})
	p.Server.Sch.AddRoutePoint(cps.Coord{X: 52, Y: 15})
	p.Server.Sch.AddRoutePoint(cps.Coord{X: 62, Y: 15})
	p.Server.Sch.AddRoutePoint(cps.Coord{X: 32, Y: 10})
	for i:= 0; i < 20; i++ {
		p.Server.Tick()
	}

}

func TestPathFinding(t *testing.T) {
	p := &cps.Params{}
	r := &cps.RegionParams{}
	p.Server = cps.FusionCenter{p, r, nil, nil, nil, nil, nil, nil, nil}

	getFlagsTest(p)

	//Adjust parameters for new scenario
	p.NumSuperNodes = 4
	p.ImageFileNameCM = "marathon_street_map.png"
	p.InputFileNameCM = "Scenario_Test2.txt"
	p.SquareRowCM = 50
	p.SquareColCM = 50

	p.FileName = p.InputFileNameCM
	cps.GetListedInput(p)
	p.MaxX = cps.GetDashedInput("maxX", p)
	p.MaxY = cps.GetDashedInput("maxY", p)
	p.TotalNodes = cps.GetDashedInput("numNodes", p)
	p.BombX = cps.GetDashedInput("bombX", p)
	p.BombY = cps.GetDashedInput("bombY", p)

	p.B = &cps.Bomb{X: p.BombX, Y: p.BombY}

	p.NumStoredSamples = p.NumStoredSamplesCM
	p.NumGridSamples = p.GridStoredSamplesCM
	p.DetectionThreshold = p.DetectionThresholdCM
	p.PositionPrint = p.PositionPrintCM
	p.GridPrint = p.GridPrintCM
	p.NodesPrint = p.NodesPrintCM
	p.SquareRowCM = p.SquareRowCM
	p.SquareColCM = p.SquareColCM

	cps.MakeBoolGrid(p)
	p.Server.Init()
	cps.ReadMap(p, r)

	cps.SetupParameters(p)

	p.Server.MakeGrid()
	p.Server.MakeSuperNodes()

	cps.GenerateRouting(p, r)
	cps.FlipSquares(p, r)

	p.Iterations_used = 0
	p.Iterations_of_event = p.IterationsCM
	p.EstimatedPingsNeeded = 10200

	reg := cps.RegionContaining(cps.Tuple{X: p.Server.Sch.SNodeList[0].GetX(), Y: p.Server.Sch.SNodeList[0].GetY()}, r)
	possible := 0
	found := 0
	for x:=0; x < p.Width; x++ {
		for y:=0; y < 150; y++ {
			if cps.RegionContaining(cps.Tuple{X:x, Y: y}, r) != -1 {
				possible++
				//fmt.Println(cps.ValidPath(reg, cps.Coord{X: y, Y: y}, r))
				if !cps.ValidPath(reg, cps.Coord{X: x, Y: y}, true, r) {
					//t.Errorf("No valid path found to point (%v,%v)\n", x ,y)
					fmt.Printf("No valid path found to point (%v,%v)\n", x ,y)
				} else {
					found++
				}
			}
		}
	}
	fmt.Printf("Reached %.2f%% of points(%v out of %v)\n", float64(found)/float64(possible) * 100, found, possible)
}

func getFlagsTest(p *cps.Params) {
	flag.StringVar(&p.CPUProfile, "cpuprofile", "", "write cpu profile to `file`")
	flag.StringVar(&p.MemProfile, "memprofile", "", "write memory profile to `file`")

	//fmt.Println(os.Args[1:], "\nhmmm? \n ") //C:\Users\Nick\Desktop\comand line experiments\src
	flag.IntVar(&p.NegativeSittingStopThresholdCM, "negativeSittingStopThreshold", -10,
		"Negative number sitting is set to when board map is reset")
	flag.IntVar(&p.SittingStopThresholdCM, "sittingStopThreshold", 5,
		"How long it takes for a node to stay seated")
	flag.Float64Var(&p.GridCapacityPercentageCM, "GridCapacityPercentage", .9,
		"Percent the sub-Grid can be filled")
	flag.StringVar(&p.InputFileNameCM, "inputFileName", "Scenario_Test.txt",
		"Name of the input text file")
	flag.StringVar(&p.SensorPath, "sensorPath", "Circle_2D.csv", "Sensor Reading Inputs")
	flag.StringVar(&p.MovementPath, "movementPath", "Circle_2D.csv", "Movement Inputs")
	flag.StringVar(&p.OutputFileNameCM, "p.OutputFileName", "Log",
		"Name of the output text file prefix")
	flag.Float64Var(&p.NaturalLossCM, "naturalLoss", .005,
		"battery loss due to natural causes")

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
	flag.Float64Var(&p.DetectionThresholdCM, "detectionThreshold", 4000.0, //11180.0,
		"Value where if a node gets this reading or higher, it will trigger a detection")
	flag.Float64Var(&p.ErrorModifierCM, "errorMultiplier", 1.0,
		"Multiplier for error values in system")
	flag.BoolVar(&p.CSVSensor, "csvSensor", false, "Read Sensor Values from CSV")
	flag.BoolVar(&p.CSVMovement, "csvMove", false, "Read Movements from CSV")
	flag.BoolVar(&p.SuperNodes, "superNodes", true, "Enable SuperNodes")
	flag.IntVar(&p.IterationsCM, "iterations", 200, "Read Movements from CSV")
	//Range 1, 2, or 4
	//Currently works for only a few numbers, can be easily expanded but is not currently dynamic
	flag.IntVar(&p.NumSuperNodes, "numSuperNodes", 8, "the number of super nodes in the simulator")
	flag.Float64Var(&p.CalibrationThresholdCM, "Recalibration Threshold", 3.0, "Value over grid average to recalibrate node")
	flag.Float64Var(&p.StdDevThresholdCM, "StandardDeviationThreshold", 1.7, "Detection Threshold based on standard deviations from mean")
	flag.Float64Var(&p.DetectionDistance, "detectionDistance", 6.0, "Detection Distance")

	flag.BoolVar(&p.PositionPrintCM, "logPosition", false, "Whether you want to write position info to a log file")
	flag.BoolVar(&p.GridPrintCM, "logGrid", false, "Whether you want to write p.Grid info to a log file")
	flag.BoolVar(&p.EnergyPrintCM, "logEnergy", false, "Whether you want to write energy into to a log file")
	flag.BoolVar(&p.NodesPrintCM, "logNodes", false, "Whether you want to write node readings to a log file")
	flag.IntVar(&p.SquareRowCM, "SquareRowCM", 1, "Number of rows of p.Grid squares, 1 through p.MaxX")
	flag.IntVar(&p.SquareColCM, "SquareColCM", 1, "Number of columns of p.Grid squares, 1 through p.MaxY")

	flag.StringVar(&p.ImageFileNameCM, "imageFileName", "routingCases.png", "Name of the input text file")
	flag.StringVar(&p.StimFileNameCM, "stimFileName", "circle_0.txt", "Name of the stimulus text file")
	flag.StringVar(&p.OutRoutingNameCM, "outRoutingName", "log.txt", "Name of the stimulus text file")
	flag.StringVar(&p.OutRoutingStatsNameCM, "outRoutingStatsName", "routingStats.txt", "Name of the output file for stats")

	flag.IntVar(&p.SuperNodeType, "p.SuperNodeType", 0, "the type of super node used in the simulator")
	flag.IntVar(&p.SuperNodeSpeed, "p.SuperNodeSpeed", 3, "the speed of the super node")
	flag.BoolVar(&p.DoOptimize, "doOptimize", false, "whether or not to optimize the simulator")

	flag.BoolVar(&p.RegionRouting, "regionRouting", true, "True if you want to use the new routing algorithm with regions and cutting")

}
