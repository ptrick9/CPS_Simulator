package main

import (
	"../simulator/cps"
	"testing"
)


//func TestDist(t *testing.T) {
//	n := &cps.NodeImpl{X: 3, Y: 3, Id: 0}
//	b := cps.Bomb{1,1}
//
//	dist := n.Distance(b)
//	partialDist := (float32(math.Pow(float64(math.Abs(float64(n.X)-float64(b.X))),2) + math.Pow(float64(math.Abs(float64(n.Y)-float64(b.Y))),2)))
//	expectedDist := float32(1000 / (math.Pow((float64(partialDist)/0.2)*0.25,1.5)))
//	if dist != expectedDist {
//		t.Errorf("Distance was incorrect, got: %f, want: %f", dist, expectedDist)
//	}
//
//	n.X = 1
//	n.Y = 1
//	dist = n.Distance(b)
//	expectedDist = 1000
//	if dist != expectedDist {
//		t.Errorf("Distance was incorrect, got: %f, want %f", dist, expectedDist)
//	}
//
//}

func TestRowCol(t *testing.T) {
	n := &cps.NodeImpl{Id:0,X:15,Y:15}
	var xDiv = 5
	var yDiv = 5
	row := n.Row(yDiv)
	col := n.Col(xDiv)
	if row != 3 {
		t.Errorf("Row was Incorrect, got: %d, want: %d",row,3)
	}
	if col != 3 {
		t.Errorf("Column was Incorrect, got: %d, want: %d",col,3)
	}
}

func TestUpdateHistory(t *testing.T) {
	var numStoredSamples = 4
	n := &cps.NodeImpl{SampleHistory:make([]float32,numStoredSamples),TotalSamples:100}
	n.UpdateHistory(4.0)
	newHist := []float32{4.0,0.0,0.0,0.0}
	var areEqual bool
	areEqual = true
	for i := range n.SampleHistory {
		if n.SampleHistory[i] != newHist[i] {
			areEqual = false
		}
	}
	if !areEqual {
		t.Errorf("updateHistory was incorect, got: %v, wanted: %v",n.SampleHistory,
			newHist)
	}
}

func TestGeoDist(t *testing.T) {
	n := &cps.NodeImpl{X:5,Y:17}
	b := cps.Bomb{X:30,Y:12}

	dist := n.GeoDist(b)
	expectedDist := 650
	if dist != float32(expectedDist) {
		t.Errorf("geoDist was incorrect, got: %v, wanted:%v",dist,expectedDist)
	}
}

func TestRecalibrate(t *testing.T) {
	n := &cps.NodeImpl{InitialSensitivity:0.82,Sensitivity:0.56,NodeTime:34}
	n.Recalibrate()
	if n.Sensitivity != 0.82 {
		t.Errorf("Recalibrate didn't fix sensitivity, should be: %f, is: %f",
			n.InitialSensitivity,n.Sensitivity)
	}
	if n.NodeTime != 0.0 {
		t.Errorf("Recalibrate didn't fix node time, should be: %d, is: %d",
			0,n.NodeTime)
	}
}

func TestSquare(t *testing.T) {
	s := cps.Square{NumEntry:5,MaxEntry:10,Tot:0,Values:make([]float32,11),Avg:0.0}
	s.TakeMeasurement(0.5)
	if s.Values[5] != 0.5 {
		t.Errorf("Square is incorrect. (values)")
	}
	if s.Tot != 0.5 {
		t.Errorf("Square is incorrect. (total)")
	}
	if s.NumEntry != 6 {
		t.Errorf("Square is incorrect. (numEntry)")
	}
	if s.Avg != (0.5 / 6) {
		t.Errorf("Square is incorrect. (avg)")
	}

	s.NumEntry = 11
	s.Values = make([]float32,11)
	s.Tot = 0
	s.Avg = 0.0
	s.TakeMeasurement(0.5)
	if s.Tot != 0.5 {
		t.Errorf("Square is incorrect (tot)")
	}
	if s.Values[1] != 0.5 {
		t.Errorf("Square is incorrect (values)")
	}
	if s.NumEntry != 12 {
		t.Errorf("Square is incorrect (numEntry)")
	}
	if s.Avg != 0.05 {
		t.Errorf("Square is incorrect (avg)")
	}
}

func TestDecrementAccel(t *testing.T){
	n := &cps.NodeImpl{Id:0,X:15,Y:15}
	oldBattery := n.Battery
	n.DecrementPowerAccel()
	newBattery := n.Battery
	if(newBattery != oldBattery - oldBattery*n.BatteryLossAccelerometer){
		t.Errorf("DecrementAccel Failed: NewBattery is: %f, NewBattery Should be: %f", newBattery, oldBattery - oldBattery*n.BatteryLossAccelerometer)
	}
}

func TestDecrementBT(t *testing.T){
	n := &cps.NodeImpl{Id:0,X:15,Y:15}
	oldBattery := n.Battery
	n.DecrementPowerBT(10)
	newBattery := n.Battery
	if(newBattery != oldBattery - oldBattery*n.BatteryLossBT){
		t.Errorf("DecrementBT Failed: NewBattery is: %f, NewBattery Should be: %f", newBattery, oldBattery - oldBattery*n.BatteryLossBT)
	}
}

func TestDecrementWifi(t *testing.T){
	n := &cps.NodeImpl{Id:0,X:15,Y:15}
	oldBattery := n.Battery
	n.DecrementPowerWifi(10)
	newBattery := n.Battery
	if(newBattery != oldBattery - oldBattery*n.BatteryLossWifi){
		t.Errorf("DecrementWifi Failed: NewBattery is: %f, NewBattery Should be: %f", newBattery, oldBattery - oldBattery*n.BatteryLossWifi)
	}
}

func TestDecrement4G(t *testing.T){
	n := &cps.NodeImpl{Id:0,X:15,Y:15}
	oldBattery := n.Battery
	n.DecrementPower4G(10)
	newBattery := n.Battery
	if(newBattery != oldBattery - oldBattery*n.BatteryLossBT){
		t.Errorf("Decrement4G Failed: NewBattery is: %f, NewBattery Should be: %f", newBattery, oldBattery - oldBattery*n.BatteryLoss4G)
	}
}

func TestDecrementGPS(t *testing.T){
	n := &cps.NodeImpl{Id:0,X:15,Y:15}
	oldBattery := n.Battery
	n.DecrementPowerGPS()
	newBattery := n.Battery
	if(newBattery != oldBattery - oldBattery*n.BatteryLossGPS){
		t.Errorf("DecrementGPS Failed: NewBattery is: %f, NewBattery Should be: %f", newBattery, oldBattery - oldBattery*n.BatteryLossGPS)
	}
}

func TestDecrementSensor(t *testing.T){
	n := &cps.NodeImpl{Id:0,X:15,Y:15}
	oldBattery := n.Battery
	n.DecrementPowerSensor()
	newBattery := n.Battery
	if(newBattery != oldBattery - oldBattery*n.BatteryLossSensor){
		t.Errorf("DecrementSenor Failed: NewBattery is: %f, NewBattery Should be: %f", newBattery, oldBattery - oldBattery*n.BatteryLossSensor)
	}
}

