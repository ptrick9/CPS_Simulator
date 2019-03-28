package cps

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
	X int
	Y int

	Avg      float32
	NumEntry int
	Values   []float32
	MaxEntry int
	Tot      float32

	NumNodes        int
	ActualNumNodes  int
	PointOfInterest bool

	StdDev       float64
	SquareValues float64
	HasDetected  bool

	CanBeTravelledTo []bool
}

//Adds a node to this Square, increasing its numNodes
func (s *Square) AddNode() {
	s.NumNodes += 1
}

//Removes a node from this Square, decreasing its numNodes
func (s *Square) RemoveNode() {
	s.NumNodes -= 1
}

//reset the square to prevent repetitive false positives
func (s *Square) Reset() {
	for i := 0; i < s.MaxEntry; i++ {
		s.Values[i] = 0.0
	}
}

//This function takes a measurement from a node inside that
//	Square and adds its value to the value list and calculates
//	the new average
func (s *Square) TakeMeasurement(x float32) {
	if s.NumEntry < s.MaxEntry {
		s.Tot += x
		s.Values[(s.NumEntry)%s.MaxEntry] = x
	} else {
		s.Tot -= s.Values[(s.NumEntry)%s.MaxEntry]
		s.Tot += x
		s.Values[(s.NumEntry)%s.MaxEntry] = x
	}
	s.NumEntry += 1
	s.Avg = s.Tot / float32(math.Min(float64(s.NumEntry), float64(s.MaxEntry)))
}

//toString method for Squares
func (s *Square) String() string {
	return fmt.Sprintf("avg: %v num: %v vals: %v max: %v nodes: %v", s.Avg, s.NumEntry, s.Values, s.MaxEntry, s.NumNodes)
}

//getter function for squareValues
func (s *Square) GetSquareValues() float64 {
	return s.SquareValues
}

//setter function for squareValues
func (s *Square) SetSquareValues(squareVals float64) {
	s.SquareValues = squareVals
}
