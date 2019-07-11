package main

import (
	"./cps"
	"flag"
	"fmt"
	"math"
	"math/rand"
	"reflect"
	"sync"
	"testing"
)

func TestInit(t *testing.T) {
	p := cps.Params{}
	r := cps.RegionParams{}
	p.Server = cps.FusionCenter{&p, &r, nil, nil, nil, nil, nil, nil, nil}
	srv := p.Server
	srv.Init()

	if len(srv.TimeBuckets) != 0 {
		t.Errorf("Initialization failed, TimeBuckets not initialized to empty list")
	}
	if len(srv.Mean) != 0{
		t.Errorf("Initialization failed, Mean not initialized empty list")
	}
	if len(srv.StdDev) != 0{
		t.Errorf("Initialization failed, StdDev not initialized empty list")
	}
	if len(srv.Variance) != 0{
		t.Errorf("Initialization failed, Variance not initialized empty list")
	}
	if len(srv.Times) != 0{
		t.Errorf("Initialization failed, Times not initialized empty list")
	}

}

func TestGetSquareAverage(t *testing.T) {
	p := cps.Params{}
	srv := p.Server
	travelList := make([]bool, 0)
	travelList = append(travelList, false)

	squ := cps.Square{0, 0, 0.0, 0, make([]float32, p.NumGridSamples),
		p.NumGridSamples, 0.0, 0, 0, false,
		0.0, 0.0, false, travelList, map[cps.Tuple]*cps.NodeImpl{},sync.Mutex{},0}

	testAvg := srv.GetSquareAverage(&squ)
	if testAvg != 0.0 {
		t.Errorf("Wrong value for square average, got: %v, want %v", testAvg, 0.0)
	}
}

func TestUpdateSquareAvg(t *testing.T) {
	p := cps.Params{}
	r := cps.RegionParams{}
	p.YDiv = 1
	p.XDiv = 1
	p.NumGridSamples = 1
	p.Server = cps.FusionCenter{&p, &r, nil, nil, nil, nil, nil, nil, nil}
	srv := p.Server
	rd := cps.Reading{10,0,0,0,0}
	travelList := make([]bool, 0)
	travelList = append(travelList, false)

	squ := cps.Square{0, 0, 0.0, 1, make([]float32, p.NumGridSamples),
		p.NumGridSamples, 0.0, 0, 0, false,
		0.0, 0.0, false, travelList, map[cps.Tuple]*cps.NodeImpl{},sync.Mutex{},0}

	p.Grid = append(p.Grid, []*cps.Square{ &squ })

	srv.UpdateSquareAvg(rd)
	avg := srv.GetSquareAverage(&squ)

	if avg != 10.0 {
		t.Errorf("Square average reading incorrectly updated, got: %v, wanted %v", avg, 10.0)
	}


}

func TestUpdateSquareNumNodes(t *testing.T) {
	p := cps.Params{}
	r := cps.RegionParams{}
	p.Server = cps.FusionCenter{&p, &r, nil, nil, nil, nil, nil, nil,nil}
	srv := p.Server
	p.NodeList = make([]cps.NodeImpl, 2)
	p.TotalNodes = 2
	p.YDiv = 2
	p.XDiv = 2

	cps.SetupRandomNodes(&p)
	p.NodeList[0].Valid = true
	p.NodeList[1].Valid = true
	p.NodeList[1].X = 2
	p.NodeList[1].Y = 2

	test_squares := make([]*cps.Square, 4)
	travelList := make([]bool, 0)
	travelList = append(travelList, false)

	//Create 4 squares
	test_squares[0] = &cps.Square{0, 0, 0.0, 0, make([]float32, p.NumGridSamples),
		p.NumGridSamples, 0.0, 0, 0, false,
		0.0, 0.0, false, travelList, map[cps.Tuple]*cps.NodeImpl{},sync.Mutex{},0}

	test_squares[1] = &cps.Square{1, 0, 0.0, 0, make([]float32, p.NumGridSamples),
		p.NumGridSamples, 0.0, 0, 0, false,
		0.0, 0.0, false, travelList, map[cps.Tuple]*cps.NodeImpl{},sync.Mutex{},0}

	test_squares[2] = &cps.Square{0, 1, 0.0, 0, make([]float32, p.NumGridSamples),
		p.NumGridSamples, 0.0, 0, 0, false,
		0.0, 0.0, false, travelList, map[cps.Tuple]*cps.NodeImpl{},sync.Mutex{},0}

	test_squares[3] = &cps.Square{1, 1, 0.0, 0, make([]float32, p.NumGridSamples),
		p.NumGridSamples, 0.0, 0, 0, false,
		0.0, 0.0, false, travelList, map[cps.Tuple]*cps.NodeImpl{},sync.Mutex{},0}

	p.Grid = append(p.Grid, test_squares[0:2])
	p.Grid = append(p.Grid, test_squares[2:])

	srv.UpdateSquareNumNodes()

	if p.Grid[0][0].ActualNumNodes != 1 {
		t.Errorf("Number of nodes in square (%v,%v) updated incorrectly, got: %v, wanted %v", 0, 0, p.Grid[0][0].ActualNumNodes, 1)
	}
	if p.Grid[0][1].ActualNumNodes != 0 {
		t.Errorf("Number of nodes in square (%v,%v) updated incorrectly, got: %v, wanted %v", 0, 1, p.Grid[0][1].ActualNumNodes, 0)
	}
	if p.Grid[1][0].ActualNumNodes != 0 {
		t.Errorf("Number of nodes in square (%v,%v) updated incorrectly, got: %v, wanted %v", 1, 0, p.Grid[1][0].ActualNumNodes, 0)
	}
	if p.Grid[1][1].ActualNumNodes != 1 {
		t.Errorf("Number of nodes in square (%v,%v) updated incorrectly, got: %v, wanted %v", 1, 1, p.Grid[1][1].ActualNumNodes, 1)
	}


}

func TestSend(t *testing.T) {
	p := cps.Params{}
	r := cps.RegionParams{}
	p.Server = cps.FusionCenter{&p, &r, nil, nil, nil, nil, nil, nil,nil}
	srv := p.Server
	//rd := cps.Reading{0,0,0,0,0}
	p.XDiv = 1
	p.YDiv = 1
	travelList := make([]bool, 0)
	travelList = append(travelList, false)
	p.NumGridSamples = 1
	rd := cps.Reading{0,0,0,0,0}


	squ := cps.Square{0, 0, 0.0, 1, make([]float32, p.NumGridSamples),
		p.NumGridSamples, 0.0, 0, 0, false,
		0.0, 0.0, false, travelList, map[cps.Tuple]*cps.NodeImpl{},sync.Mutex{},0}
	node := &cps.NodeImpl{X: 0, Y: 0, Id: len(p.NodeList), SampleHistory: []float32{}, Concentration: 0,
		Cascade: 0, Battery: 0, BatteryLossScalar: 0}

	p.Grid = append(p.Grid, []*cps.Square{ &squ })

	srv.Send(node, rd)

	if srv.Times[0] != true {
		t.Errorf("Time 0 not included in packet, got %v, wanted true", srv.Times[0])
	}
	if srv.TimeBuckets[0][0].SensorVal != 0 {
		t.Errorf("TimeBuckets not updated properly, got %v, wanted %v", srv.TimeBuckets[0][0], 0)
	}

}

func TestCalcStats(t *testing.T) {
	p := cps.Params{}
	r := cps.RegionParams{}
	//getFlagsForTest(&p)
	//cps.SetupParameters(&p)

	p.Server = cps.FusionCenter{&p, &r, nil, nil, nil, nil, nil, nil, nil}
	srv := p.Server
	srv.Init()
	srv.Times = make(map[int]bool)
	srv.Times[0] = true
	srv.Times[1] = true

	rd1 := cps.Reading{1,0,0,0,1}
	rd2 := cps.Reading{2,0,0,0,2}
	rd3 := cps.Reading{3,0,0,0,3}
	rd4 := cps.Reading{4,0,0,0,4}
	rd5 := cps.Reading{6,0,0,0,5}
	rd6 := cps.Reading{8,0,0,0,6}

	p.B = &cps.Bomb{0,0}

	//srv.TimeBuckets = [][]float64{[]float64{1, 2, 3, 4}, []float64{2, 4, 6, 8}}
	srv.TimeBuckets = [][]cps.Reading{[]cps.Reading{rd1, rd2, rd3, rd4}, []cps.Reading{rd2, rd4, rd5, rd6}}
	srv.Times = map[int]bool{0:true, 1:true}
	srv.CalcStats()

	expectedMean := []float64{2.5, 5.0}
	expectedStdDev := []float64{1.118033988749895, 2.23606797749979}
	expectedVariance := []float64{1.25, 5.0}

	for i := 0; i < 2; i++ {
		if srv.Mean[i] != expectedMean[i] {
			t.Errorf("Incorrect value for mean, got %v, wanted %v", srv.Mean[i], expectedMean[i])
		}
		if srv.StdDev[i] != expectedStdDev[i] {
			t.Errorf("Incorrect value for standard deviation, got %v, wanted %v", srv.StdDev[i], expectedStdDev[i])
		}
		if srv.Variance[i] != expectedVariance[i] {
			t.Errorf("Incorrect value for variance, got %v, wanted %v", srv.Variance[i], expectedVariance[i])
		}
	}
}


func TestPrintStats(t *testing.T) {
	//
}

func makeNodesForTest(p *cps.Params) {
	//p := cps.Params{}
	for i := 0; i < len(p.NodeEntryTimes); i++ {

		if p.Iterations_used == p.NodeEntryTimes[i][2] {

			var initHistory = make([]float32, p.NumStoredSamples)

			xx := cps.RangeInt(1, p.MaxX)
			yy := cps.RangeInt(1, p.MaxY)
			for p.BoolGrid[xx][yy] == true {
				xx = cps.RangeInt(1, p.MaxX)
				yy = cps.RangeInt(1, p.MaxY)
			}

			p.NodeList = append(p.NodeList, cps.NodeImpl{X: xx, Y: yy, Id: len(p.NodeList), SampleHistory: initHistory, Concentration: 0,
				Cascade: i, Battery: p.BatteryCharges[i], BatteryLossScalar: p.BatteryLosses[i]})

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

func TestMakeGrid(t *testing.T) {
	p := cps.Params{}
	r := cps.RegionParams{}
	p.SquareColCM = 5
	p.SquareRowCM = 5
	p.Server = cps.FusionCenter{&p,&r, nil,nil,nil,nil,nil,nil, nil}
	srv := p.Server

	srv.MakeGrid()
	if len(p.Grid) * len(p.Grid[0]) != 25 {
		t.Errorf("Incorrect grid size, got %v, wanted %v", len(p.Grid) * len(p.Grid[0]), 25)
	}
	for i := 0; i < 5; i++ {
		for j := 0; j < 5; j++ {
			if reflect.TypeOf(p.Grid[i][j]) != reflect.TypeOf(&cps.Square{}) {
				t.Errorf("Square not generated correctly at x: %v, y: %v, got type %v, expected type %v", i, j, reflect.TypeOf(&cps.Square{}), reflect.TypeOf(p.Grid[i][j]))
			}
		}
	}
}

func TestMakeSuperNodes(t *testing.T) {
	p := cps.Params{}
	r := cps.RegionParams{}
	p.Server = cps.FusionCenter{&p,&r, nil,nil,nil,nil,nil,nil, nil}
	srv := p.Server
	srv.Init()

	p.Height = 2
	p.Width = 2
	p.NumSuperNodes = 4

	srv.MakeSuperNodes()

	if p.NumSuperNodes != 4 {
		t.Errorf("Incorrect number of super nodes, got %v, wanted %v", p.NumSuperNodes, 4)
	}

}

func TestCheckDetections(t *testing.T) {
	p := cps.Params{}
	r := cps.RegionParams{}
	p.Server = cps.FusionCenter{&p,&r,nil,nil,nil,nil,nil, nil, nil}
	srv := p.Server
	srv.Init()
	p.YDiv = 1
	p.XDiv = 1
	getFlagsForTest(&p)
	cps.SetupParameters(&p)

	p.B = &cps.Bomb{1,1}

	//sch := srv.Sch//sch := &cps.Scheduler{}
	p.NodeList = append(p.NodeList, cps.NodeImpl{X: 1, Y: 1, Id: len(p.NodeList), SampleHistory: []float32{}, Concentration: 0,
		Cascade: 0, Battery: 0, BatteryLossScalar: 0})
	srv.Sch.SNodeList = append(srv.Sch.SNodeList, &cps.Sn_zero{&p,&r,&cps.Supern{&p,&r, &cps.NodeImpl{X: 0, Y: 0, Id: 0}, 1,
		1, p.SuperNodeRadius, p.SuperNodeRadius, 0, make([]cps.Coord, 1), make([]cps.Coord, 0),
		cps.Coord{X: 1, Y: 1}, 0, 0, 0, 0, 0, make([]cps.Coord, 0)}})

	srv.MakeGrid()

	p.Grid[0][0].Avg = 4
	srv.CheckDetections()


	if !p.Grid[0][0].HasDetected {
		t.Errorf("True positive missed!")
	}

}

func TestGetMedian(t *testing.T) {
	p := cps.Params{}
	srv := p.Server
	arr := []float64{ 1.0, 2.0, 3.0, 4.0 }
	median := srv.GetMedian(arr)
	if median != 2.5 {
		t.Errorf("Incorrect median, got %v, wanted %v", median, 2.5)
	}
	arr = []float64{ 1.0, 2.0, 3.0, 4.0 , 5.0}
	median = srv.GetMedian(arr)
	if median != 3 {
		t.Errorf("Incorrect median, got %v, wanted %v", median, 3.0)
	}
}

func getFlagsForTest(p *cps.Params) {
	flag.IntVar(&p.MaxX, "MaxX", 268, "Maximum X value")
	flag.IntVar(&p.MaxY, "MaxY", 191, "Maximum Y value")

	flag.IntVar(&p.NegativeSittingStopThresholdCM, "negativeSittingStopThreshold", -10,
		"Negative number sitting is set to when board map is reset")
	flag.IntVar(&p.SittingStopThresholdCM, "sittingStopThreshold", 5,
		"How long it takes for a node to stay seated")
	flag.Float64Var(&p.GridCapacityPercentageCM, "GridCapacityPercentage", .9,
		"Percent the sub-Grid can be filled")
	flag.StringVar(&p.InputFileNameCM, "inputFileName", "Scenario_4.txt",
		"Name of the input text file")

	flag.Float64Var(&p.NaturalLossCM, "naturalLoss", .005,
		"battery loss due to natural causes")

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


	flag.IntVar(&p.SuperNodeType, "p.SuperNodeType", 0, "the type of super node used in the simulator")

	flag.IntVar(&p.SuperNodeSpeed, "p.SuperNodeSpeed", 3, "the speed of the super node")

	flag.BoolVar(&p.DoOptimize, "doOptimize", false, "whether or not to optimize the simulator")

	flag.BoolVar(&p.PositionPrintCM, "logPosition", false, "Whether you want to write position info to a log file")
	flag.BoolVar(&p.GridPrintCM, "logGrid", false, "Whether you want to write p.Grid info to a log file")
	flag.BoolVar(&p.EnergyPrintCM, "logEnergy", false, "Whether you want to write energy into to a log file")
	flag.BoolVar(&p.NodesPrintCM, "logNodes", false, "Whether you want to write node readings to a log file")
	flag.IntVar(&p.SquareRowCM, "SquareRowCM", 1, "Number of rows of p.Grid squares, 1 through p.MaxX")
	flag.IntVar(&p.SquareColCM, "SquareColCM", 1, "Number of columns of p.Grid squares, 1 through p.MaxY")

	flag.StringVar(&p.ImageFileNameCM, "imageFileName", "testMap.png", "Name of the input text file")
	flag.StringVar(&p.StimFileNameCM, "stimFileName", "circle_0.txt", "Name of the stimulus text file")
	flag.StringVar(&p.OutRoutingNameCM, "outRoutingName", "log.txt", "Name of the stimulus text file")
	flag.StringVar(&p.OutRoutingStatsNameCM, "outRoutingStatsName", "routingStats.txt", "Name of the output file for stats")

	flag.BoolVar(&p.RegionRouting, "regionRouting", true, "True if you want to use the new routing algorithm with regions and cutting")
	flag.BoolVar(&p.CSVMovement, "CSVMovement", false, "True if you want to use the csv for node movement")
	flag.BoolVar(&p.CSVSensor, "CSVSensor", false, "True if you want to use the csv for node movement")

	flag.Parse()
	fmt.Println("Natural Loss: ", p.NaturalLossCM)
	fmt.Println("Period of extra sampling due to high speed: ", p.MovementSamplingPeriodCM)
	fmt.Println("Maximum size of buffer posible: ", p.MaxBufferCapacityCM)
	fmt.Println("Energy model type:", p.EnergyModelCM)
	fmt.Println("Sensor Sampling Period:", p.SensorSamplingPeriodCM)
	fmt.Println("GPS Sampling Period:", p.GPSSamplingPeriodCM)
	fmt.Println("Server Sampling Period:", p.ServerSamplingPeriodCM)
	fmt.Println("Number of Node Stored Samples:", p.NumStoredSamplesCM)
	fmt.Println("Number of Grid Stored Samples:", p.GridStoredSamplesCM)
	fmt.Println("Detection Threshold:", p.DetectionThresholdCM)
}

