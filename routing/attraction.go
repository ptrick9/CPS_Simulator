package main

import "fmt"

//Attraction type contains the locations of the attraction point
type attraction struct {
	x int
	y int
}

//function to "teleport" an attraction
func (a *attraction) move(x, y int) {
	a.x = x
	a.y = y
}

func (a attraction) String() string {
	return fmt.Sprintf("X: %v, Y: %v", a.x, a.y)
}
