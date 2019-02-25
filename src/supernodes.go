package main

import (
	"fmt"
	"math"
)

//SuperNodeMovement controls all the movement of the super nodes
//	as well as some methods only defined for super nodes
type SuperNodeParent interface {
	NodeParent
	tick()

	pathMove()
	centMove()

	updateLoc()

	addRoutePoint(Coord)
	updatePath()
	route([][]*Square, Coord, Coord, []Coord) []Coord

	incSquareMoved(int)
	incAllPoints()

	getRoutePath() []Coord
	getRoutePoints() []Coord
	getNumDest() int
	getCenter() Coord
	getSquaresMoved() int
	getPointsVisited() int
	getId() int
	getAvgResponseTime() float64
	getAllPoints() []Coord
	getSuperNodeType() int

	setNumDest(int)
	setRoutePath([]Coord)
	setRoutePoints([]Coord)
}

//Super nodes travel through the gird based of a route
//	dictated by a grid path algorithm
//They contain attributes not contained by other nodes
//	to control their movement through the grid
type supern struct {
	*NodeImpl
	x_speed int
	y_speed int

	xRadius int
	yRadius int

	numDestinations int

	routePoints []Coord
	routePath   []Coord

	center Coord

	squaresMoved int

	pointsVisited   int
	totResponseTime int
	avgResponseTime float64

	superNodeType int

	allPoints []Coord
}

//toString for super nodes
func (n supern) String() string {
	return fmt.Sprintf("x: %v y: %v RoutePoints: %v Path: %v",
		n.x, n.y, n.routePoints, n.routePath)
}

//Moves the super node along the path determined by the routePath list
//It also maintains the routePoints and routePath lists, deleting elements
//	once they have been visited by the super node
func (n *supern) pathMove() {

	//Increase the time of the points that are being travelled to
	for p, _ := range n.routePoints {
		n.routePoints[p].time += 1
	}
	//The index of the point the super node is moving to
	//When the superNodeSpeed is greater than 1, the removal_index is used
	//	in the remove_range function
	removal_index := -1

	//Moves the super node to the next Coord in the routePath
	//If there are enough Coords for the super node to move its full speed it will
	//If there are not enough Coords in the routePath list, the super node will move
	//	as far as possible
	if len(n.routePath) >= superNodeSpeed {
		n.x = n.routePath[superNodeSpeed-1].x
		n.y = n.routePath[superNodeSpeed-1].y

		//Saves the value of the index to remove the Coords
		removal_index = superNodeSpeed
	} else {
		n.x = n.routePath[len(n.routePath)-1].x
		n.y = n.routePath[len(n.routePath)-1].y

		//Saves the value of the index to remove the Coords
		removal_index = len(n.routePath)
	}
	//The first element in the routePoints list is always the current location
	//	of the super node and is so updated
	n.updateLoc()

	//A boolean flag used for super nodes of type 2
	//Squares moved towards the center do not count towards their squaresTravelled total
	countSquares := true

	//Loops through the routePath at the points to be removed
	//If one of those points is a point of interest in the routePoints list it is removed
	for i := 0; i < removal_index; i++ {
		if len(n.routePoints) > 1 {
			if (n.routePoints[1].x == n.routePath[i].x) && (n.routePoints[1].y == n.routePath[i].y) {
				//Increases the number of points visited by the super node
				//Calculates the totResponseTime and avgResponseTime
				n.pointsVisited++
				n.totResponseTime += n.routePoints[1].time
				n.avgResponseTime = float64(n.totResponseTime / n.pointsVisited)

				//If a super node of type 2 is moving towards its center than the squares it
				//	moves should not count towards it total
				if (n.routePoints[1].x == n.center.x) && (n.routePoints[1].y == n.center.y) && (superNodeType == 2) {
					countSquares = false
				}
				//It is then removed from the routePoints list
				n.setRoutePoints(n.routePoints[:1+copy(n.routePoints[1:], n.routePoints[2:])])

				//Since this point has been visited and removed the numDestinations is
				//	also decreased
				n.numDestinations--
			}
		}
	}
	//The range of points travelled by the super node are then removed
	n.setRoutePath(n.routePath[:0+copy(n.routePath[0:], n.routePath[removal_index:])])

	//Increases the amount of squares moved by the super node
	//For super nodes of type 2, their movement back to the center position is not counted
	//	as squares travelled since this movement is idle behavior
	if countSquares {
		n.incSquareMoved(removal_index)
	} else {
		countSquares = true
	}
}

//The route function adds Coords to the routePath of the specific
//	super node
//It recursively finds the Square with the lest number of nodes between
//	the two Coords
//Once the Square with the lowest numNodes is found the route
//	function is called again between the beginning node/the lowest
//	node and the lowest node/the end node
//Eventually this adds Coords from the beginning to the end travelling
//	along the least populated nodes along the way
func (n supern) route(grid [][]*Square, c1, c2 Coord, list []Coord) []Coord {
	lowNum := 100
	lowCoord := Coord{x: -1, y: -1}

	newx1 := 0
	newx2 := 0
	newy1 := 0
	newy2 := 0

	//If the two Coords are the same the recursion is broken
	if (c1.y == c2.y) && (c1.x == c2.x) {
		ret_list := make([]Coord, 0)
		return ret_list
		//Otherwise the bounds of the search need to be defined correctly
	} else {
		if c1.x > c2.x {
			newx1 = c2.x
			newx2 = c1.x + 1
		} else {
			newx1 = c1.x
			newx2 = c2.x + 1
		}
		if c1.y > c2.y {
			newy1 = c2.y
			newy2 = c1.y + 1
		} else {
			newy1 = c1.y
			newy2 = c2.y + 1
		}
	}

	//This double for loop find the Square with the lowest numNodes
	//	in between the two Coords
	for i := newy1; i < newy2; i++ {
		for j := newx1; j < newx2; j++ {
			if !(i == c1.y && j == c1.x) && !(i == c2.y && j == c2.x) {
				if grid[i][j].numNodes < lowNum {
					lowNum = grid[i][j].numNodes
					lowCoord = Coord{x: j, y: i}
				}
			}
		}
	}

	//If a Square has been found, the function is called recursively and that
	//	Square's coordinates are added to the routePath
	if (lowCoord.x != -1) && (lowCoord.y != -1) {
		list = n.route(grid, c1, lowCoord, list)
		list = append(list, lowCoord)
		list = n.route(grid, lowCoord, c2, list)
	}

	//Otherwise this recursion loop ends with this function call
	return list
}

//Moves the super node back towards it's assigned center point
//Instead of constructing an entire path to follow and then scrapping
//	it when a new point of interest arrives, centMove adds only the
//	point of the path and then calls path move
//This allows the super node to have empty routePoints and routePath
//	when a new point of interest is added
func (n *supern) centMove() {
	arr := make([]Coord, 0)
	arr = aStar(n.routePoints[0], n.center)
	arr = append(arr, n.center)
	n.routePath = append(n.routePath, arr...)

	n.pathMove()
}

//Updates the location of the super node within the gird
//Determines what square it's in not the exact x,y
func (n *supern) updateLoc() {
	n.routePoints[0].x = n.x
	n.routePoints[0].y = n.y
}

//This function helps keep track of the amount of squares the
//	super node is travelling, this is used to test the effectiveness
//	of the routing algorithms
func (n *supern) incSquareMoved(num int) {
	n.squaresMoved += num
}

//Increments the time of all points not currently being visited
//	byt the super node
func (n *supern) incAllPoints() {
	for p, _ := range n.allPoints {
		n.allPoints[p].time++
	}
}

//Initializes the regionList for the super node
func makeRegionList(sNodeNum int) []Region {
	r := make([]Region, 4)
	r0 := make([]Coord, 0)
	r1 := make([]Coord, 0)
	r2 := make([]Coord, 0)
	r3 := make([]Coord, 0)

	if numSuperNodes == 1 {
		r[0] = Region{Coord{x: maxX / 2, y: maxY / 2}, Coord{x: 0, y: 0}, r0}
		r[1] = Region{Coord{x: maxX / 2, y: maxY / 2}, Coord{x: maxX, y: 0}, r1}
		r[2] = Region{Coord{x: maxX / 2, y: maxY / 2}, Coord{x: 0, y: maxY}, r2}
		r[3] = Region{Coord{x: maxX / 2, y: maxY / 2}, Coord{x: maxX, y: maxY}, r3}
	} else if numSuperNodes == 2 {
		r[0] = Region{Coord{x: (3 - (2 * sNodeNum)) * (maxX / 4), y: maxY / 2},
			Coord{x: (maxX / 2) - ((maxX / 2) * sNodeNum), y: 0}, r0}
		r[1] = Region{Coord{x: (3 - (2 * sNodeNum)) * (maxX / 4), y: maxY / 2},
			Coord{x: maxX - ((maxX / 2) * sNodeNum), y: 0}, r1}
		r[2] = Region{Coord{x: (3 - (2 * sNodeNum)) * (maxX / 4), y: maxY / 2},
			Coord{x: (maxX / 2) - ((maxX / 2) * sNodeNum), y: maxY}, r2}
		r[3] = Region{Coord{x: (3 - (2 * sNodeNum)) * (maxX / 4), y: maxY / 2},
			Coord{x: maxX - ((maxX / 2) * sNodeNum), y: maxY}, r3}
	} else if numSuperNodes == 4 {
		if sNodeNum == 0 {
			r[0] = Region{Coord{x: maxX / 4, y: maxY / 4}, Coord{x: 0, y: 0}, r0}
			r[1] = Region{Coord{x: maxX / 4, y: maxY / 4}, Coord{x: maxX / 2, y: 0}, r1}
			r[2] = Region{Coord{x: maxX / 4, y: maxY / 4}, Coord{x: 0, y: maxY / 2}, r2}
			r[3] = Region{Coord{x: maxX / 4, y: maxY / 4}, Coord{x: maxX / 2, y: maxY / 2}, r3}
		} else if sNodeNum == 1 {
			r[0] = Region{Coord{x: 3 * (maxX / 4), y: maxY / 4}, Coord{x: maxX / 2, y: 0}, r0}
			r[1] = Region{Coord{x: 3 * (maxX / 4), y: maxY / 4}, Coord{x: maxX, y: 0}, r1}
			r[2] = Region{Coord{x: 3 * (maxX / 4), y: maxY / 4}, Coord{x: maxX / 2, y: maxY / 2}, r2}
			r[3] = Region{Coord{x: 3 * (maxX / 4), y: maxY / 4}, Coord{x: maxX, y: maxY / 2}, r3}
		} else if sNodeNum == 2 {
			r[0] = Region{Coord{x: maxX / 4, y: 3 * (maxY / 4)}, Coord{x: 0, y: maxY / 2}, r0}
			r[1] = Region{Coord{x: maxX / 4, y: 3 * (maxY / 4)}, Coord{x: maxX / 2, y: maxY / 2}, r1}
			r[2] = Region{Coord{x: maxX / 4, y: 3 * (maxY / 4)}, Coord{x: 0, y: maxY}, r2}
			r[3] = Region{Coord{x: maxX / 4, y: 3 * (maxY / 4)}, Coord{x: maxX / 2, y: maxY}, r3}
		} else if sNodeNum == 3 {
			r[0] = Region{Coord{x: 3 * (maxX / 4), y: 3 * (maxY / 4)}, Coord{x: maxX / 2, y: maxY / 2}, r0}
			r[1] = Region{Coord{x: 3 * (maxX / 4), y: 3 * (maxY / 4)}, Coord{x: maxX, y: maxY / 2}, r1}
			r[2] = Region{Coord{x: 3 * (maxX / 4), y: 3 * (maxY / 4)}, Coord{x: maxX / 2, y: maxY}, r2}
			r[3] = Region{Coord{x: 3 * (maxX / 4), y: 3 * (maxY / 4)}, Coord{x: maxX, y: maxY}, r3}
		}
	}
	return r
}

//This function determines the center position for the super nodes
//The center depends on the number of super nodes in the simulation
func makeCenter1(sNodeNum int) (Coord, int, int) {
	nodeCenter := center
	x_val := maxX / 2
	y_val := maxY / 2
	//If there is only one super node it should be place in the center of
	//	the gird
	if numSuperNodes != 1 {
		//Determining the angle at which to separate the super nodes
		//This value is multiplied by the current number of this super node
		//For example if there are 3 super nodes, they should be separated by
		//	120 degress around the center point of the grid
		angle := (2 * math.Pi) / float64(numSuperNodes) * float64(sNodeNum)

		//Determining how far from the center this super node should be
		x_dist := math.Cos(angle)
		y_dist := math.Sin(angle)

		//Initializing the x and y position of the super node that
		//	correspond to that center
		x_val = center.x + int(x_dist*float64(maxX/4))
		y_val = center.y + int(y_dist*float64(maxY/4))

		//Creating the Coord
		nodeCenter = Coord{x: x_val, y: y_val}
	}

	return nodeCenter, x_val, y_val
}

//This function determines the center position for the super nodes' circles
//This currently works for only 4 super nodes as this is a special test for a
// unique version of super nodes of type 1
//This version has the super nodes' circles positioned in the four corners
func makeCenter1_corners(sNodeNum int) (Coord, int, int, int, int) {

	//Radius of super nodes of type 3
	xRad := int(float64(maxX)/2.83) + 1
	yRad := int(float64(maxY)/2.83) + 1

	nodeCenter := center
	x_val := maxX / 4
	y_val := maxY / 4

	if numSuperNodes == 4 {
		if sNodeNum == 0 {
			nodeCenter = Coord{x: center.x - x_val, y: center.y - y_val}
		} else if sNodeNum == 1 {
			nodeCenter = Coord{x: center.x + x_val, y: center.y - y_val}
		} else if sNodeNum == 2 {
			nodeCenter = Coord{x: center.x - x_val, y: center.y + y_val}
		} else if sNodeNum == 3 {
			nodeCenter = Coord{x: center.x + x_val, y: center.y + y_val}
		}
		x_val = nodeCenter.x
		y_val = nodeCenter.y
	} else {
		fmt.Println("ONLY USE THIS FUNCTION WITH 4 SUPER NODES FOR NOW")
	}
	return nodeCenter, x_val, y_val, xRad, yRad
}

//This function determines the center position for the super nodes' circles
//This currently works for only 4 super nodes as this is a special test for a
//	unique version of super nodes of type 1
//This version has the super nodes' circles positioned on the four sides
func makeCenter1_sides(sNodeNum int) (Coord, int, int, int, int) {

	//Radius of super nodes of type 3
	xRad := maxX / 2
	yRad := maxY / 2

	nodeCenter := center
	x_val := maxX / 2
	y_val := maxY / 2

	if numSuperNodes == 4 {
		if sNodeNum == 0 {
			nodeCenter = Coord{x: center.x - x_val, y: center.y}
		} else if sNodeNum == 1 {
			nodeCenter = Coord{x: center.x, y: center.y - y_val}
		} else if sNodeNum == 2 {
			nodeCenter = Coord{x: center.x + x_val - 1, y: center.y}
		} else if sNodeNum == 3 {
			nodeCenter = Coord{x: center.x, y: center.y + y_val - 1}
		}
		x_val = nodeCenter.x
		y_val = nodeCenter.y
	} else {
		fmt.Println("ONLY USE THIS FUNCTION WITH 4 SUPER NODES FOR NOW")
	}
	return nodeCenter, x_val, y_val, xRad, yRad
}

//This function determines the center position for the super nodes' circles
//This currently works for only 4 super nodes as this is a special test for a
//	unique version of super nodes of type 1
//This version has the super nodes' circles positioned in the four corners
//However, unlike the other circular centers the radii of these circles are larger
//	therefore the centers of the super nodes are different
func makeCenter1_largeCorners(sNodeNum int) (Coord, int, int, int, int) {

	//Radius of super nodes of type 3
	xRad := maxX
	yRad := maxY

	nodeCenter := center
	x_val := maxX / 2
	y_val := maxY / 2

	if numSuperNodes == 4 {
		if sNodeNum == 0 {
			nodeCenter = Coord{x: center.x - x_val, y: center.y - y_val}
		} else if sNodeNum == 1 {
			nodeCenter = Coord{x: center.x + x_val - 1, y: center.y - y_val}
		} else if sNodeNum == 2 {
			nodeCenter = Coord{x: center.x - x_val, y: center.y + y_val - 1}
		} else if sNodeNum == 3 {
			nodeCenter = Coord{x: center.x + x_val - 1, y: center.y + y_val - 1}
		}
		x_val = nodeCenter.x
		y_val = nodeCenter.y
	} else {
		fmt.Println("ONLY USE THIS FUNCTION WITH 4 SUPER NODES FOR NOW")
	}
	return nodeCenter, x_val, y_val, xRad, yRad
}

//This function determines the center position of super nodes of type 2
//Super nodes of type 2 are centered inside their respective regions
func makeCenter2(sNodeNum int, r_list []Region) (Coord, int, int) {
	nodeCenter := center
	x_val := maxX / 2
	y_val := maxY / 2
	//If there is only one super node it should be place in the center of
	//	the gird
	if numSuperNodes != 1 {
		//The center of the super node should be the center of the region it occupies
		nodeCenter = r_list[sNodeNum].center

		//These expressions translate the loction of the square to the location of
		//	the individual x, y location
		x_val = nodeCenter.x
		y_val = nodeCenter.y
	}
	return nodeCenter, x_val, y_val
}

//Various getters for super node attributes to be
// accessed by the SuperNodeMovement interface
func (n *supern) getRoutePath() []Coord {
	return n.routePath
}
func (n *supern) getRoutePoints() []Coord {
	return n.routePoints
}
func (n *supern) getNumDest() int {
	return n.numDestinations
}
func (n *supern) getCenter() Coord {
	return n.center
}
func (n *supern) getSquaresMoved() int {
	return n.squaresMoved
}
func (n *supern) getPointsVisited() int {
	return n.pointsVisited
}
func (n *supern) getId() int {
	return n.id
}
func (n *supern) getAvgResponseTime() float64 {
	return n.avgResponseTime
}
func (n *supern) getSuperNodeType() int {
	return n.superNodeType
}
func (n *supern) getAllPoints() []Coord {
	return n.allPoints
}
func (n *supern) getX() int {
	return n.x
}
func (n *supern) getY() int {
	return n.y
}

//Various setters for super node attributes to be
// accessed by the SuperNodeMovement interface
func (n *supern) setNumDest(d int) {
	n.numDestinations = d
}
func (n *supern) setRoutePath(c []Coord) {
	n.routePath = c
}
func (n *supern) setRoutePoints(c []Coord) {
	n.routePoints = c
}
