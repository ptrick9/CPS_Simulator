package main

import (
	"../simulator/cps"
	"math"
	"math/rand"
	"sync"
	"testing"
)

func TestInit(t *testing.T) {
	p := cps.Params{}
	p.Server = cps.FusionCenter{&p, nil, nil, nil, nil, nil}
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
		0.0, 0.0, false, travelList, sync.Mutex{}}

	testAvg := srv.GetSquareAverage(&squ)
	if testAvg != 0.0 {
		t.Errorf("Wrong value for square average, got: %v, want %v", testAvg, 0.0)
	}
}

func TestUpdateSquareAvg(t *testing.T) {
	p := cps.Params{}
	p.YDiv = 1
	p.XDiv = 1
	p.NumGridSamples = 1
	p.Server = cps.FusionCenter{&p, nil, nil, nil, nil, nil}
	srv := p.Server
	rd := cps.Reading{10,0,0,0,0}
	travelList := make([]bool, 0)
	travelList = append(travelList, false)

	squ := cps.Square{0, 0, 0.0, 1, make([]float32, p.NumGridSamples),
		p.NumGridSamples, 0.0, 0, 0, false,
		0.0, 0.0, false, travelList, sync.Mutex{}}

	p.Grid = append(p.Grid, []*cps.Square{ &squ })

	srv.UpdateSquareAvg(rd)
	avg := srv.GetSquareAverage(&squ)

	if avg != 10.0 {
		t.Errorf("Square average reading incorrectly updated, got: %v, wanted %v", avg, 10.0)
	}


}

func TestUpdateSquareNumNodes(t *testing.T) {
	p := cps.Params{}
	p.Server = cps.FusionCenter{&p, nil, nil, nil, nil, nil}
	srv := p.Server
	p.NodeList = make([]cps.NodeImpl, 2)
	p.NumNodes = 2
	p.YDiv = 2
	p.XDiv = 2
	makeNodesForTest(&p)

	test_squares := make([]*cps.Square, 4)
	travelList := make([]bool, 0)
	travelList = append(travelList, false)

	//Create 4 squares
	test_squares[0] = &cps.Square{0, 0, 0.0, 0, make([]float32, p.NumGridSamples),
		p.NumGridSamples, 0.0, 0, 0, false,
		0.0, 0.0, false, travelList, sync.Mutex{}}

	test_squares[1] = &cps.Square{1, 0, 0.0, 0, make([]float32, p.NumGridSamples),
		p.NumGridSamples, 0.0, 0, 0, false,
		0.0, 0.0, false, travelList, sync.Mutex{}}

	test_squares[2] = &cps.Square{0, 1, 0.0, 0, make([]float32, p.NumGridSamples),
		p.NumGridSamples, 0.0, 0, 0, false,
		0.0, 0.0, false, travelList, sync.Mutex{}}

	test_squares[3] = &cps.Square{1, 1, 0.0, 0, make([]float32, p.NumGridSamples),
		p.NumGridSamples, 0.0, 0, 0, false,
		0.0, 0.0, false, travelList, sync.Mutex{}}

	p.Grid = append(p.Grid, test_squares[0:2])
	p.Grid = append(p.Grid, test_squares[2:])

	srv.UpdateSquareNumNodes()

	if p.Grid[0][0].ActualNumNodes != 2 {
		t.Errorf("Number of nodes in square (%v,%v) updated incorrectly, got: %v, wanted %v", 0, 0, p.Grid[0][0].ActualNumNodes, 2)
	}
	if p.Grid[0][1].ActualNumNodes != 0 {
		t.Errorf("Number of nodes in square (%v,%v) updated incorrectly, got: %v, wanted %v", 0, 1, p.Grid[0][1].ActualNumNodes, 0)
	}
	if p.Grid[1][0].ActualNumNodes != 0 {
		t.Errorf("Number of nodes in square (%v,%v) updated incorrectly, got: %v, wanted %v", 1, 0, p.Grid[1][0].ActualNumNodes, 0)
	}
	if p.Grid[1][1].ActualNumNodes != 0 {
		t.Errorf("Number of nodes in square (%v,%v) updated incorrectly, got: %v, wanted %v", 1, 1, p.Grid[1][1].ActualNumNodes, 0)
	}


}

func TestSend(t *testing.T) {
	p := cps.Params{}
	p.Server = cps.FusionCenter{&p, nil, nil, nil, nil, nil}
	srv := p.Server
	rd := cps.Reading{0,0,0,0,0}
	p.XDiv = 1
	p.YDiv = 1
	travelList := make([]bool, 0)
	travelList = append(travelList, false)
	p.NumGridSamples = 1

	squ := cps.Square{0, 0, 0.0, 1, make([]float32, p.NumGridSamples),
		p.NumGridSamples, 0.0, 0, 0, false,
		0.0, 0.0, false, travelList, sync.Mutex{}}

	p.Grid = append(p.Grid, []*cps.Square{ &squ })

	srv.Send(rd)

	if srv.Times[0] != true {
		t.Errorf("Time 0 not included in packet, got %v, wanted true", srv.Times[0])
	}
	if srv.TimeBuckets[0][0] != 0 {
		t.Errorf("TimeBuckets not updated properly, got %v, wanted %v", srv.TimeBuckets[0][0], 0)
	}

}

func TestCalcStats(t *testing.T) {
	p := cps.Params{}
	p.Server = cps.FusionCenter{&p, nil, nil, nil, nil, nil}
	srv := p.Server
	srv.Times = make(map[int]bool)
	srv.Times[0] = true
	srv.Times[1] = true


	srv.TimeBuckets = [][]float64{[]float64{1, 2, 3, 4}, []float64{2, 4, 6, 8}}
	srv.CalcStats()

	expectedMean := []float64{2.5, 5.0}
	expectedStdDev := []float64{1.118033988749895, 2.236067977}
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
	for i := 0; i < len(p.Npos); i++ {

		if p.Iterations_used == p.Npos[i][2] {

			var initHistory = make([]float32, p.NumStoredSamples)

			xx := cps.RangeInt(1, p.MaxX)
			yy := cps.RangeInt(1, p.MaxY)
			for p.BoolGrid[xx][yy] == true {
				xx = cps.RangeInt(1, p.MaxX)
				yy = cps.RangeInt(1, p.MaxY)
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

