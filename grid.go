package main

import (
	"fmt"
	"math"
)

//Square type contains the average of the readings in the Square,
//	the number of total entries that have been taken in that
//	specific Square, the list of values that are used to calculate
//	the average, the maximum number of entries that list can
//	contain, the total of those entries before they are averaged,
//	and the number of nodes inside of the specific Square
type Square struct {
	x int
	y int

	avg      float32
	numEntry int
	values   []float32
	maxEntry int
	tot      float32

	numNodes        int
	actualNumNodes  int
	pointOfInterest bool

	stdDev       float64
	squareValues float64
	hasDetected  bool

	canBeTravelledTo []bool
}

//Adds a node to this Square, increasing its numNodes
func (s *Square) addNode() {
	s.numNodes += 1
}

//Removes a node from this Square, decreasing its numNodes
func (s *Square) removeNode() {
	s.numNodes -= 1
}

//reset the square to prevent repetitive false positives
func (s *Square) reset() {
	for i := 0; i < s.maxEntry; i++ {
		s.values[i] = 0.0
	}
}

//This function takes a measurement from a node inside that
//	Square and adds its value to the value list and calculates
//	the new average
func (s *Square) takeMeasurement(x float32) {
	if s.numEntry < s.maxEntry {
		s.tot += x
		s.values[(s.numEntry)%s.maxEntry] = x
	} else {
		s.tot -= s.values[(s.numEntry)%s.maxEntry]
		s.tot += x
		s.values[(s.numEntry)%s.maxEntry] = x
	}
	s.numEntry += 1
	s.avg = s.tot / float32(math.Min(float64(s.numEntry), float64(s.maxEntry)))
}

//toString method for Squares
func (s *Square) String() string {
	return fmt.Sprintf("avg: %v num: %v vals: %v max: %v nodes: %v", s.avg, s.numEntry, s.values, s.maxEntry, s.numNodes)
}

//getter function for squareValues
func (s *Square) getSquareValues() float64 {
	return s.squareValues
}

//setter function for squareValues
func (s *Square) setSquareValues(squareVals float64) {
	s.squareValues = squareVals
}
