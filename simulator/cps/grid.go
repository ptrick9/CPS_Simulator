package cps

import (
	"fmt"
	"math"
	"sync"
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
	NodesInSquare	 map[Tuple]*NodeImpl
	Lock sync.Mutex

	LastReadingTime int
	Navigable		bool
	Visited		bool

	left		*Square
	right 		*Square
	up 			*Square
	down 		*Square

	ConnectedSquares []*Square

	SuperNodeCluster int
}

//Adds a node to this Square, increasing its numNodes
func (s *Square) AddNode(n * NodeImpl) {
	s.NodesInSquare[Tuple{int(n.X),int(n.Y)}] = n;
	s.NumNodes += 1
}

//returns map of all nodes in the same Square (NodesInSquare)
func (s *Square) nearbyNodes () (map[Tuple]*NodeImpl){
	return s.NodesInSquare
}

//Removes a node from this Square, decreasing its numNodes
func (s *Square) RemoveNode(n * NodeImpl) {
	var s_temp = s.NodesInSquare[Tuple{int(n.X),int(n.Y)}]
	if s_temp == n{ //checks if node is in Square before deleting from square, avoids exception
		delete(s.NodesInSquare, Tuple{int(n.X),int(n.Y)});
		s.NumNodes -= 1
	}
}

//reset the square to prevent repetitive false positives
func (s *Square) Reset() {
	for i := 0; i < s.MaxEntry; i++ {
		s.Values[i] = 0.0
	}
}

//This function takes a measurement from a node inside that
//	Square and adds its Value to the Value list and calculates
//	the new average
func (s *Square) TakeMeasurement(x float32) {
	s.Lock.Lock()
	defer s.Lock.Unlock()
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
