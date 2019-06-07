package main

import "math"

//This function takes a start and an end Coord and returns a list of Coords
//	that connect them without going through walls
func aStar(a Coord, b Coord) []Coord {
	//The openList contains all the points the can be travelled to
	//The closedList contains all the points that have already been
	//	travelled to
	openList := make([]Coord, 0)
	closedList := make([]Coord, 0)

	//The starting point is added to the openList to start the algorithm
	openList = append(openList, a)

	for len(openList) > 0 {

		//Finds the Coord in the openList with the lowest score, meaning
		//	it is the closest to the destination
		lowScore := maxX * maxY
		lowInd := -1
		for s, _ := range openList {
			if openList[s].score < lowScore {
				lowScore = openList[s].score
				lowInd = s
			}
		}

		//The Coord with the lowest score is removed from the openList
		//	and added to the closedList
		curr := openList[lowInd]
		openList = remove_coord_index(openList, lowInd)
		closedList = append(closedList, curr)

		//If the destination is in the closedList, meaning the current
		//	Coord is the destination, the algorithm is completed
		contains, _ := b.in(closedList)

		if contains {
			//The destination has been reached
			break
		}

		//Saves the parent Coord for checking in the getWalkable function
		//The parent Coord need not be checked since it is definitely not a viable
		//	path for the supernode since it has already been there
		parent_coord := Coord{x: -1, y: -1}
		if !curr.equals(a) {
			parent_coord = *(curr.parent)
		}

		//adjacent is a Coord list of all the Coords that are adjacent to
		//	the current Coord
		//Random Walls and previously travelled Coords will not appear in this list
		adjacent := curr.getWalkable(b, closedList, parent_coord)

		//This loops through all the adjacent Coords to the current Coord
		//As long as the Coords ae not in the closedList, they are added to
		//	the openList
		//If they are already in the openList their score is updated to the lowest
		//	of the two scores
		for j, _ := range adjacent {
			inList, _ := adjacent[j].in(closedList)
			if !inList {
				inOpenList, o := adjacent[j].in(openList)
				if !inOpenList {
					openList = append(openList, adjacent[j])
				} else {
					if adjacent[j].score < openList[o].score {
						openList[o].score = adjacent[j].score
					}
				}
			}
		}
	}

	//Starting with the destination added to the closedList, each Coord's parent
	//	is added to the path list until the parent does not exist
	//The starting Coord does not have a parent so it will end with the starting
	//	Coord
	path := make([]Coord, 0)
	path = append(path, closedList[len(closedList)-1])

	parent := path[0].parent

	for parent != nil {
		path = append(path, *parent)
		parent = parent.parent
	}

	//Reverse the path because even though it goes from source to destination,
	//	reading the list by going through each Coord's parent means the path
	//	list is ordered from the destination to the source
	//Therefore this needs to be reversed
	for i := len(path)/2 - 1; i >= 0; i-- {
		opp := len(path) - 1 - i
		path[i], path[opp] = path[opp], path[i]
	}
	return path
}

//This function returns a list of Coords that are adjacent to
//	the Coord that called this function
//It only returns Coords that can be walked to, this does not
//	include walls and previously travelled Coords
func (a *Coord) getWalkable(b Coord, closedList []Coord, p Coord) []Coord {
	walkableList := make([]Coord, 0)

	//Boundary conditions that check whether AStar coordinates
	//	are walls, previously travelled AStar coordinates or
	//	outside the grid
	if a.x+1 < maxX {
		newA := Coord{a, a.x + 1, a.y, 0, 0, 0, 0}
		inList, _ := newA.in(closedList)
		if boardMap[a.y][a.x] != -1 && !inList && !newA.equals(p) {
			newA.setScore(b)
			walkableList = append(walkableList, newA)
		}
	}
	if a.x-1 >= 0 {
		newA := Coord{a, a.x - 1, a.y, 0, 0, 0, 0}
		inList, _ := newA.in(closedList)
		if boardMap[a.y][a.x] != -1 && !inList && !newA.equals(p) {
			newA.setScore(b)
			walkableList = append(walkableList, newA)
		}
	}
	if a.y+1 < maxY {
		newA := Coord{a, a.x, a.y + 1, 0, 0, 0, 0}
		inList, _ := newA.in(closedList)
		if boardMap[a.y][a.x] != -1 && !inList && !newA.equals(p) {
			newA.setScore(b)
			walkableList = append(walkableList, newA)
		}
	}
	if a.y-1 >= 0 {
		newA := Coord{a, a.x, a.y - 1, 0, 0, 0, 0}
		inList, _ := newA.in(closedList)
		if boardMap[a.y][a.x] != -1 && !inList && !newA.equals(p) {
			newA.setScore(b)
			walkableList = append(walkableList, newA)
		}
	}
	return walkableList
}

//setScore takes a Coord and the destination and
//	calculates the new score for this Coord
func (a *Coord) setScore(b Coord) {
	//The G value is the amount of squares travelled since the
	//	starting point, therefore its just the parent's G value + 1
	a.g = a.parent.g + 1

	//The H value is the actual distance between the current AStar
	//	coordinate and the destination AStar coordinate
	h := math.Sqrt(math.Pow(float64(a.x-b.x), 2.0) + math.Pow(float64(a.y-b.y), 2.0))
	a.h = int(h)

	//The score is the two values summed
	a.score = a.g + a.h
}

//This function returns whether or not a Coord
//	is in a list of Coords and where in the list it is
func (a *Coord) in(list []Coord) (bool, int) {
	for i, _ := range list {
		if (list[i].x == a.x) && (list[i].y == a.y) {
			return true, i
		}
	}
	return false, -1
}
