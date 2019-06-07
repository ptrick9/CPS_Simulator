package cps

import (
	"testing"
		"math"
	)


func TestDist(t *testing.T) {
	n := &NodeImpl{x: 3, y: 3, id: 0}
	b := bomb{1,1}

	dist := n.distance(b)
	partialDist := (float32(math.Pow(float64(math.Abs(float64(n.x)-float64(b.x))),2) + math.Pow(float64(math.Abs(float64(n.y)-float64(b.y))),2)))
	expectedDist := float32(1000 / (math.Pow((float64(partialDist)/0.2)*0.25,1.5)))
	if dist != expectedDist {
		t.Errorf("Distance was incorrect, got: %f, want: %f", dist, expectedDist)
	}

	n.x = 1
	n.y = 1
	dist = n.distance(b)
	expectedDist = 1000
	if dist != expectedDist {
		t.Errorf("Distance was incorrect, got: %f, want %f", dist, expectedDist)
	}

}

func TestRowCol(t *testing.T) {
	n := &NodeImpl{id:0,x:15,y:15}
	xDiv = 5
	yDiv = 5
	row := n.row(yDiv)
	col := n.col(xDiv)
	if row != 3 {
		t.Errorf("Row was Incorrect, got: %d, want: %d",row,3)
	}
	if col != 3 {
		t.Errorf("Column was Incorrect, got: %d, want: %d",col,3)
	}
}

func TestUpdateHistory(t *testing.T) {
	numStoredSamples = 4
	n := &NodeImpl{sampleHistory:make([]float32,numStoredSamples),totalSamples:100}
	n.updateHistory(4.0)
	newHist := []float32{4.0,0.0,0.0,0.0}
	var areEqual bool
	areEqual = true
	for i := range n.sampleHistory {
		if n.sampleHistory[i] != newHist[i] {
			areEqual = false
		}
	}
	if !areEqual {
		t.Errorf("updateHistory was incorect, got: %v, wanted: %v",n.sampleHistory,
			newHist)
	}
}

func TestGeoDist(t *testing.T) {
	n := &NodeImpl{x:5,y:17}
	b := bomb{x:30,y:12}

	dist := n.geoDist(b)
	expectedDist := 650
	if dist != float32(expectedDist) {
		t.Errorf("geoDist was incorrect, got: %v, wanted:%v",dist,expectedDist)
	}
}

func TestRecalibrate(t *testing.T) {
	n := &NodeImpl{initialSensitivity:0.82,sensitivity:0.56,nodeTime:34}
	n.recalibrate()
	if n.sensitivity != 0.82 {
		t.Errorf("Recalibrate didn't fix sensitivity, should be: %f, is: %f",
			n.initialSensitivity,n.sensitivity)
	}
	if n.nodeTime != 0.0 {
		t.Errorf("Recalibrate didn't fix node time, should be: %d, is: %d",
			0,n.nodeTime)
	}
}

func TestSquare(t *testing.T) {
	s := Square{numEntry:5,maxEntry:10,tot:0,values:make([]float32,11),avg:0.0}
	s.takeMeasurement(0.5)
	if s.values[5] != 0.5 {
		t.Errorf("Square is incorrect. (values)")
	}
	if s.tot != 0.5 {
		t.Errorf("Square is incorrect. (total)")
	}
	if s.numEntry != 6 {
		t.Errorf("Square is incorrect. (numEntry)")
	}
	if s.avg != (0.5 / 6) {
		t.Errorf("Square is incorrect. (avg)")
	}

	s.numEntry = 11
	s.values = make([]float32,11)
	s.tot = 0
	s.avg = 0.0
	s.takeMeasurement(0.5)
	if s.tot != 0.5 {
		t.Errorf("Square is incorrect (tot)")
	}
	if s.values[1] != 0.5 {
		t.Errorf("Square is incorrect (values)")
	}
	if s.numEntry != 12 {
		t.Errorf("Square is incorrect (numEntry)")
	}
	if s.avg != 0.05 {
		t.Errorf("Square is incorrect (avg)")
	}
}