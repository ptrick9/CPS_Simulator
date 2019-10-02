package cps

import "math"

//This function takes a start and an end Coord and returns a list of Coords
//	that connect them without going through walls
func AStar(a Coord, b Coord, p *Params) []Coord {
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
		lowScore := p.MaxX * p.MaxY
		lowInd := -1
		for s, _ := range openList {
			if openList[s].Score < lowScore {
				lowScore = openList[s].Score
				lowInd = s
			}
		}

		//The Coord with the lowest score is removed from the openList
		//	and added to the closedList
		curr := openList[lowInd]
		openList = Remove_coord_index(openList, lowInd)
		closedList = append(closedList, curr)

		//If the destination is in the closedList, meaning the current
		//	Coord is the destination, the algorithm is completed
		contains, _ := b.In(closedList)

		if contains {
			//The destination has been reached
			break
		}

		//Saves the parent Coord for checking in the getWalkable function
		//The parent Coord need not be checked since it is definitely not a viable
		//	path for the supernode since it has already been there
		parent_coord := Coord{X: -1, Y: -1}
		if !curr.Equals(a) {
			parent_coord = *(curr.Parent)
		}

		//adjacent is a Coord list of all the Coords that are adjacent to
		//	the current Coord
		//Random Walls and previously travelled Coords will not appear in this list
		adjacent := curr.GetWalkable(b, closedList, parent_coord, p)

		//This loops through all the adjacent Coords to the current Coord
		//As long as the Coords ae not in the closedList, they are added to
		//	the openList
		//If they are already in the openList their score is updated to the lowest
		//	of the two scores
		for j, _ := range adjacent {
			inList, _ := adjacent[j].In(closedList)
			if !inList {
				inOpenList, o := adjacent[j].In(openList)
				if !inOpenList {
					openList = append(openList, adjacent[j])
				} else {
					if adjacent[j].Score < openList[o].Score {
						openList[o].Score = adjacent[j].Score
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

	parent := path[0].Parent

	for parent != nil {
		path = append(path, *parent)
		parent = parent.Parent
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
func (a *Coord) GetWalkable(b Coord, closedList []Coord, p Coord, pp *Params) []Coord {
	walkableList := make([]Coord, 0)

	//Boundary conditions that check whether AStar coordinates
	//	are walls, previously travelled AStar coordinates or
	//	outside the grid
	if a.X+1 < pp.MaxX {
		newA := Coord{a, a.X + 1, a.Y, 0, 0, 0, 0}
		inList, _ := newA.In(closedList)
		if pp.BoardMap[a.Y][a.X] != -1 && !inList && !newA.Equals(p) {
			newA.SetScore(b, pp)
			walkableList = append(walkableList, newA)
		}
	}
	if a.X-1 >= 0 {
		newA := Coord{a, a.X - 1, a.Y, 0, 0, 0, 0}
		inList, _ := newA.In(closedList)
		if pp.BoardMap[a.Y][a.X] != -1 && !inList && !newA.Equals(p) {
			newA.SetScore(b, pp)
			walkableList = append(walkableList, newA)
		}
	}
	if a.Y+1 < pp.MaxY {
		newA := Coord{a, a.X, a.Y + 1, 0, 0, 0, 0}
		inList, _ := newA.In(closedList)
		if pp.BoardMap[a.Y][a.X] != -1 && !inList && !newA.Equals(p) {
			newA.SetScore(b, pp)
			walkableList = append(walkableList, newA)
		}
	}
	if a.Y-1 >= 0 {
		newA := Coord{a, a.X, a.Y - 1, 0, 0, 0, 0}
		inList, _ := newA.In(closedList)
		if pp.BoardMap[a.Y][a.X] != -1 && !inList && !newA.Equals(p) {
			newA.SetScore(b, pp)
			walkableList = append(walkableList, newA)
		}
	}
	return walkableList
}

//setScore takes a Coord and the destination and
//	calculates the new score for this Coord
func (a *Coord) SetScore(b Coord, p *Params) {
	//The G Value is the amount of squares travelled since the
	//	starting point, therefore its just the parent's G Value + 1
	a.G = a.Parent.G + 1

	//The H Value is the actual distance between the current AStar
	//	coordinate and the destination AStar coordinate
	h := math.Sqrt(math.Pow(float64(a.X-b.X), 2.0) + math.Pow(float64(a.Y-b.Y), 2.0))
	a.H = int(h)

	//The score is the two values summed
	a.Score = a.G + a.H
}

//This function returns whether or not a Coord
//	is in a list of Coords and where in the list it is
func (a *Coord) In(list []Coord) (bool, int) {
	for i, _ := range list {
		if (list[i].X == a.X) && (list[i].Y == a.Y) {
			return true, i
		}
	}
	return false, -1
}
