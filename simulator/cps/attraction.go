package cps

import "fmt"

//Attraction type contains the locations of the attraction point
type Attraction struct {
	X int
	Y int
}

//function to "teleport" an attraction
func (a *Attraction) Move(x, y int) {
	a.X = x
	a.Y = y
}

func (a Attraction) String() string {
	return fmt.Sprintf("X: %v, Y: %v", a.X, a.Y)
}
