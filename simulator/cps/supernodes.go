package cps

import (
	"fmt"
	"math"
)

//SuperNodeMovement controls all the movement of the super nodes
//	as well as some methods only defined for super nodes
type SuperNodeParent interface {
	NodeParent
	Tick(p *Params, r *RegionParams)

	PathMove(p *Params)
	CentMove(p *Params)

	UpdateLoc()

	AddRoutePoint(Coord, *Params, *RegionParams)
	UpdatePath(p *Params, r *RegionParams)
	Route(grid [][]*Square, c1 Coord, c2 Coord, list[]Coord) []Coord

	IncSquareMoved(int)
	IncAllPoints()

	GetRoutePath() []Coord
	GetRoutePoints() []Coord
	GetNumDest() int
	GetCenter() Coord
	GetSquaresMoved() int
	GetPointsVisited() int
	GetId() int
	GetAvgResponseTime() float64
	GetAllPoints() []Coord
	GetSuperNodeType() int

	SetNumDest(int)
	SetRoutePath([]Coord)
	SetRoutePoints([]Coord)
}

//Super nodes travel through the gird based of a Route
//	dictated by a grid path algorithm
//They contain attributes not contained by other nodes
//	to control their movement through the grid
type Supern struct {
	*NodeImpl
	X_speed int
	Y_speed int

	XRadius int
	YRadius int

	NumDestinations int

	RoutePoints []Coord
	RoutePath   []Coord

	Center Coord

	SquaresMoved int

	PointsVisited   int
	TotResponseTime int
	AvgResponseTime float64

	SuperNodeType int

	AllPoints []Coord
}

//toString for super nodes
func (n Supern) String() string {
	return fmt.Sprintf("x: %v y: %v RoutePoints: %v Path: %v",
		n.X, n.Y, n.RoutePoints, n.RoutePath)
}

//Moves the super node along the path determined by the RoutePath list
//It also maintains the RoutePoints and RoutePath lists, deleting elements
//	once they have been visited by the super node
func (n *Supern) PathMove(p *Params) {

	//Increase the time of the points that are being travelled to
	for p, _ := range n.RoutePoints {
		n.RoutePoints[p].Time += 1
	}
	//The index of the point the super node is moving to
	//When the superNodeSpeed is greater than 1, the removal_index is used
	//	in the remove_range function
	removal_index := -1

	//Moves the super node to the next Coord in the RoutePath
	//If there are enough Coords for the super node to move its full speed it will
	//If there are not enough Coords in the RoutePath list, the super node will move
	//	as far as possible
	if len(n.RoutePath) >= p.SuperNodeSpeed {
		n.X = n.RoutePath[p.SuperNodeSpeed-1].X
		n.Y = n.RoutePath[p.SuperNodeSpeed-1].Y

		//Saves the value of the index to remove the Coords
		removal_index = p.SuperNodeSpeed
	} else {
		n.X = n.RoutePath[len(n.RoutePath)-1].X
		n.Y = n.RoutePath[len(n.RoutePath)-1].Y

		//Saves the value of the index to remove the Coords
		removal_index = len(n.RoutePath)
	}
	//The first element in the RoutePoints list is always the current location
	//	of the super node and is so updated
	n.UpdateLoc()

	//A boolean flag used for super nodes of type 2
	//Squares moved towards the center do not count towards their squaresTravelled total
	countSquares := true

	//Loops through the RoutePath at the points to be removed
	//If one of those points is a point of interest in the RoutePoints list it is removed
	for i := 0; i < removal_index; i++ {
		if len(n.RoutePoints) > 1 {
			if (n.RoutePoints[1].X == n.RoutePath[i].X) && (n.RoutePoints[1].Y == n.RoutePath[i].Y) {
				//Increases the number of points visited by the super node
				//Calculates the totResponseTime and avgResponseTime
				n.PointsVisited++
				n.TotResponseTime += n.RoutePoints[1].Time
				n.AvgResponseTime = float64(n.TotResponseTime / n.PointsVisited)

				//If a super node of type 2 is moving towards its center than the squares it
				//	moves should not count towards it total
				if (n.RoutePoints[1].X == n.Center.X) && (n.RoutePoints[1].Y == n.Center.Y) && (p.SuperNodeType == 2) {
					countSquares = false
				}
				//It is then removed from the RoutePoints list
				n.SetRoutePoints(n.RoutePoints[:1+copy(n.RoutePoints[1:], n.RoutePoints[2:])])

				//Since this point has been visited and removed the numDestinations is
				//	also decreased
				n.NumDestinations--
			}
		}
	}
	//The range of points travelled by the super node are then removed
	n.SetRoutePath(n.RoutePath[:0+copy(n.RoutePath[0:], n.RoutePath[removal_index:])])

	//Increases the amount of squares moved by the super node
	//For super nodes of type 2, their movement back to the center position is not counted
	//	as squares travelled since this movement is idle behavior
	if countSquares {
		n.IncSquareMoved(removal_index)
	} else {
		countSquares = true
	}
}

//The Route function adds Coords to the RoutePath of the specific
//	super node
//It recursively finds the Square with the least number of nodes between
//	the two Coords
//Once the Square with the lowest numNodes is found the Route
//	function is called again between the beginning node/the lowest
//	node and the lowest node/the end node
//Eventually this adds Coords from the beginning to the end travelling
//	along the least populated nodes along the way
func (n Supern) Route(grid [][]*Square, c1 Coord, c2 Coord, list []Coord) []Coord {
	lowNum := 100
	lowCoord := Coord{X: -1, Y: -1}

	newx1 := 0
	newx2 := 0
	newy1 := 0
	newy2 := 0

	//If the two Coords are the same the recursion is broken
	if (c1.Y == c2.Y) && (c1.X == c2.X) {
		ret_list := make([]Coord, 0)
		return ret_list
		//Otherwise the bounds of the search need to be defined correctly
	} else {
		if c1.X > c2.X {
			newx1 = c2.X
			newx2 = c1.X + 1
		} else {
			newx1 = c1.X
			newx2 = c2.X + 1
		}
		if c1.Y > c2.Y {
			newy1 = c2.Y
			newy2 = c1.Y + 1
		} else {
			newy1 = c1.Y
			newy2 = c2.Y + 1
		}
	}

	//This double for loop find the Square with the lowest numNodes
	//	in between the two Coords
	for i := newy1; i < newy2; i++ {
		for j := newx1; j < newx2; j++ {
			if !(i == c1.Y && j == c1.X) && !(i == c2.Y && j == c2.X) {
				if grid[i][j].NumNodes < lowNum {
					lowNum = grid[i][j].NumNodes
					lowCoord = Coord{X: j, Y: i}
				}
			}
		}
	}

	//If a Square has been found, the function is called recursively and that
	//	Square's coordinates are added to the RoutePath
	if (lowCoord.X != -1) && (lowCoord.Y != -1) {
		list = n.Route(grid, c1, lowCoord, list)
		list = append(list, lowCoord)
		list = n.Route(grid, lowCoord, c2, list)
	}

	//Otherwise this recursion loop ends with this function call
	return list
}

//Moves the super node back towards it's assigned center point
//Instead of constructing an entire path to follow and then scrapping
//	it when a new point of interest arrives, centMove adds only the
//	point of the path and then calls path move
//This allows the super node to have empty RoutePoints and RoutePath
//	when a new point of interest is added
func (n *Supern) CentMove(p *Params) {
	arr := make([]Coord, 0)
	arr = AStar(n.RoutePoints[0], n.Center, p)
	arr = append(arr, n.Center)
	n.RoutePath = append(n.RoutePath, arr...)

	n.PathMove(p)
}

//Updates the location of the super node within the gird
//Determines what square it's in not the exact x,y
func (n *Supern) UpdateLoc() {
	n.RoutePoints[0].X = n.X
	n.RoutePoints[0].Y = n.Y
}

//This function helps keep track of the amount of squares the
//	super node is travelling, this is used to test the effectiveness
//	of the routing algorithms
func (n *Supern) IncSquareMoved(num int) {
	n.SquaresMoved += num
}

//Increments the time of all points not currently being visited
//	byt the super node
func (n *Supern) IncAllPoints() {
	for p, _ := range n.AllPoints {
		n.AllPoints[p].Time++
	}
}

//Initializes the regionList for the super node
func MakeRegionList(sNodeNum int, p *Params) []Region {
	r := make([]Region, 4)
	r0 := make([]Coord, 0)
	r1 := make([]Coord, 0)
	r2 := make([]Coord, 0)
	r3 := make([]Coord, 0)

	if p.NumSuperNodes == 1 {
		r[0] = Region{Coord{X: p.MaxX / 2, Y: p.MaxY / 2}, Coord{X: 0, Y: 0}, r0}
		r[1] = Region{Coord{X: p.MaxX / 2, Y: p.MaxY / 2}, Coord{X: p.MaxX, Y: 0}, r1}
		r[2] = Region{Coord{X: p.MaxX / 2, Y: p.MaxY / 2}, Coord{X: 0, Y: p.MaxY}, r2}
		r[3] = Region{Coord{X: p.MaxX / 2, Y: p.MaxY / 2}, Coord{X: p.MaxX, Y: p.MaxY}, r3}
	} else if p.NumSuperNodes == 2 {
		r[0] = Region{Coord{X: (3 - (2 * sNodeNum)) * (p.MaxX / 4), Y: p.MaxY / 2},
			Coord{X: (p.MaxX / 2) - ((p.MaxX / 2) * sNodeNum), Y: 0}, r0}
		r[1] = Region{Coord{X: (3 - (2 * sNodeNum)) * (p.MaxX / 4), Y: p.MaxY / 2},
			Coord{X: p.MaxX - ((p.MaxX / 2) * sNodeNum), Y: 0}, r1}
		r[2] = Region{Coord{X: (3 - (2 * sNodeNum)) * (p.MaxX / 4), Y: p.MaxY / 2},
			Coord{X: (p.MaxX / 2) - ((p.MaxX / 2) * sNodeNum), Y: p.MaxY}, r2}
		r[3] = Region{Coord{X: (3 - (2 * sNodeNum)) * (p.MaxX / 4), Y: p.MaxY / 2},
			Coord{X: p.MaxX - ((p.MaxX / 2) * sNodeNum), Y: p.MaxY}, r3}
	} else if p.NumSuperNodes == 4 {
		if sNodeNum == 0 {
			r[0] = Region{Coord{X: p.MaxX / 4, Y: p.MaxY / 4}, Coord{X: 0, Y: 0}, r0}
			r[1] = Region{Coord{X: p.MaxX / 4, Y: p.MaxY / 4}, Coord{X: p.MaxX / 2, Y: 0}, r1}
			r[2] = Region{Coord{X: p.MaxX / 4, Y: p.MaxY / 4}, Coord{X: 0, Y: p.MaxY / 2}, r2}
			r[3] = Region{Coord{X: p.MaxX / 4, Y: p.MaxY / 4}, Coord{X: p.MaxX / 2, Y: p.MaxY / 2}, r3}
		} else if sNodeNum == 1 {
			r[0] = Region{Coord{X: 3 * (p.MaxX / 4), Y: p.MaxY / 4}, Coord{X: p.MaxX / 2, Y: 0}, r0}
			r[1] = Region{Coord{X: 3 * (p.MaxX / 4), Y: p.MaxY / 4}, Coord{X: p.MaxX, Y: 0}, r1}
			r[2] = Region{Coord{X: 3 * (p.MaxX / 4), Y: p.MaxY / 4}, Coord{X: p.MaxX / 2, Y: p.MaxY / 2}, r2}
			r[3] = Region{Coord{X: 3 * (p.MaxX / 4), Y: p.MaxY / 4}, Coord{X: p.MaxX, Y: p.MaxY / 2}, r3}
		} else if sNodeNum == 2 {
			r[0] = Region{Coord{X: p.MaxX / 4, Y: 3 * (p.MaxY / 4)}, Coord{X: 0, Y: p.MaxY / 2}, r0}
			r[1] = Region{Coord{X: p.MaxX / 4, Y: 3 * (p.MaxY / 4)}, Coord{X: p.MaxX / 2, Y: p.MaxY / 2}, r1}
			r[2] = Region{Coord{X: p.MaxX / 4, Y: 3 * (p.MaxY / 4)}, Coord{X: 0, Y: p.MaxY}, r2}
			r[3] = Region{Coord{X: p.MaxX / 4, Y: 3 * (p.MaxY / 4)}, Coord{X: p.MaxX / 2, Y: p.MaxY}, r3}
		} else if sNodeNum == 3 {
			r[0] = Region{Coord{X: 3 * (p.MaxX / 4), Y: 3 * (p.MaxY / 4)}, Coord{X: p.MaxX / 2, Y: p.MaxY / 2}, r0}
			r[1] = Region{Coord{X: 3 * (p.MaxX / 4), Y: 3 * (p.MaxY / 4)}, Coord{X: p.MaxX, Y: p.MaxY / 2}, r1}
			r[2] = Region{Coord{X: 3 * (p.MaxX / 4), Y: 3 * (p.MaxY / 4)}, Coord{X: p.MaxX / 2, Y: p.MaxY}, r2}
			r[3] = Region{Coord{X: 3 * (p.MaxX / 4), Y: 3 * (p.MaxY / 4)}, Coord{X: p.MaxX, Y: p.MaxY}, r3}
		}
	}
	return r
}

//This function determines the center position for the super nodes
//The center depends on the number of super nodes in the simulation
func MakeCenter1(sNodeNum int, p *Params) (Coord, int, int) {
	nodeCenter := p.Center
	x_val := p.MaxX / 2
	y_val := p.MaxY / 2
	//If there is only one super node it should be place in the center of
	//	the gird
	if p.NumSuperNodes != 1 {
		//Determining the angle at which to separate the super nodes
		//This value is multiplied by the current number of this super node
		//For example if there are 3 super nodes, they should be separated by
		//	120 degress around the center point of the grid
		angle := (2 * math.Pi) / float64(p.NumSuperNodes) * float64(sNodeNum)

		//Determining how far from the center this super node should be
		x_dist := math.Cos(angle)
		y_dist := math.Sin(angle)

		//Initializing the x and y position of the super node that
		//	correspond to that center
		x_val = p.Center.X + int(x_dist*float64(p.MaxX/4))
		y_val = p.Center.Y + int(y_dist*float64(p.MaxY/4))

		//Creating the Coord
		nodeCenter = Coord{X: x_val, Y: y_val}
	}

	return nodeCenter, x_val, y_val
}

//This function determines the center position for the super nodes' circles
//This currently works for only 4 super nodes as this is a special test for a
// unique version of super nodes of type 1
//This version has the super nodes' circles positioned in the four corners
func MakeCenter1_corners(sNodeNum int, p *Params) (Coord, int, int, int, int) {

	//Radius of super nodes of type 3
	xRad := int(float64(p.MaxX)/2.83) + 1
	yRad := int(float64(p.MaxY)/2.83) + 1

	nodeCenter := p.Center
	x_val := p.MaxX / 4
	y_val := p.MaxY / 4

	if p.NumSuperNodes == 4 {
		if sNodeNum == 0 {
			nodeCenter = Coord{X: p.Center.X - x_val, Y: p.Center.Y - y_val}
		} else if sNodeNum == 1 {
			nodeCenter = Coord{X: p.Center.X + x_val, Y: p.Center.Y - y_val}
		} else if sNodeNum == 2 {
			nodeCenter = Coord{X: p.Center.X - x_val, Y: p.Center.Y + y_val}
		} else if sNodeNum == 3 {
			nodeCenter = Coord{X: p.Center.X + x_val, Y: p.Center.Y + y_val}
		}
		x_val = nodeCenter.X
		y_val = nodeCenter.Y
	} else {
		fmt.Println("ONLY USE THIS FUNCTION WITH 4 SUPER NODES FOR NOW")
	}
	return nodeCenter, x_val, y_val, xRad, yRad
}

//This function determines the center position for the super nodes' circles
//This currently works for only 4 super nodes as this is a special test for a
//	unique version of super nodes of type 1
//This version has the super nodes' circles positioned on the four sides
func MakeCenter1_sides(sNodeNum int, p *Params) (Coord, int, int, int, int) {

	//Radius of super nodes of type 3
	xRad := p.MaxX / 2
	yRad := p.MaxY / 2

	nodeCenter := p.Center
	x_val := p.MaxX / 2
	y_val := p.MaxY / 2

	if p.NumSuperNodes == 4 {
		if sNodeNum == 0 {
			nodeCenter = Coord{X: p.Center.X - x_val, Y: p.Center.Y}
		} else if sNodeNum == 1 {
			nodeCenter = Coord{X: p.Center.X, Y: p.Center.Y - y_val}
		} else if sNodeNum == 2 {
			nodeCenter = Coord{X: p.Center.X + x_val - 1, Y: p.Center.Y}
		} else if sNodeNum == 3 {
			nodeCenter = Coord{X: p.Center.X, Y: p.Center.Y + y_val - 1}
		}
		x_val = nodeCenter.X
		y_val = nodeCenter.Y
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
func MakeCenter1_largeCorners(sNodeNum int, p *Params) (Coord, int, int, int, int) {

	//Radius of super nodes of type 3
	xRad := p.MaxX
	yRad := p.MaxY

	nodeCenter := p.Center
	x_val := p.MaxX / 2
	y_val := p.MaxY / 2

	if p.NumSuperNodes == 4 {
		if sNodeNum == 0 {
			nodeCenter = Coord{X: p.Center.X - x_val, Y: p.Center.Y - y_val}
		} else if sNodeNum == 1 {
			nodeCenter = Coord{X: p.Center.X + x_val - 1, Y: p.Center.Y - y_val}
		} else if sNodeNum == 2 {
			nodeCenter = Coord{X: p.Center.X - x_val, Y: p.Center.Y + y_val - 1}
		} else if sNodeNum == 3 {
			nodeCenter = Coord{X: p.Center.X + x_val - 1, Y: p.Center.Y + y_val - 1}
		}
		x_val = nodeCenter.X
		y_val = nodeCenter.Y
	} else {
		fmt.Println("ONLY USE THIS FUNCTION WITH 4 SUPER NODES FOR NOW")
	}
	return nodeCenter, x_val, y_val, xRad, yRad
}

//This function determines the center position of super nodes of type 2
//Super nodes of type 2 are centered inside their respective regions
func MakeCenter2(sNodeNum int, r_list []Region, p *Params) (Coord, int, int) {
	nodeCenter := p.Center
	x_val := p.MaxX / 2
	y_val := p.MaxY / 2
	//If there is only one super node it should be place in the center of
	//	the gird
	if p.NumSuperNodes != 1 {
		//The center of the super node should be the center of the region it occupies
		nodeCenter = r_list[sNodeNum].Center

		//These expressions translate the loction of the square to the location of
		//	the individual x, y location
		x_val = nodeCenter.X
		y_val = nodeCenter.Y
	}
	return nodeCenter, x_val, y_val
}

//Various getters for super node attributes to be
// accessed by the SuperNodeMovement interface
func (n *Supern) GetRoutePath() []Coord {
	return n.RoutePath
}
func (n *Supern) GetRoutePoints() []Coord {
	return n.RoutePoints
}
func (n *Supern) GetNumDest() int {
	return n.NumDestinations
}
func (n *Supern) GetCenter() Coord {
	return n.Center
}
func (n *Supern) GetSquaresMoved() int {
	return n.SquaresMoved
}
func (n *Supern) GetPointsVisited() int {
	return n.PointsVisited
}
func (n *Supern) GetId() int {
	return n.Id
}
func (n *Supern) GetAvgResponseTime() float64 {
	return n.AvgResponseTime
}
func (n *Supern) GetSuperNodeType() int {
	return n.SuperNodeType
}
func (n *Supern) GetAllPoints() []Coord {
	return n.AllPoints
}
func (n *Supern) GetX() int {
	return n.X
}
func (n *Supern) GetY() int {
	return n.Y
}

//Various setters for super node attributes to be
// accessed by the SuperNodeMovement interface
func (n *Supern) SetNumDest(d int) {
	n.NumDestinations = d
}
func (n *Supern) SetRoutePath(c []Coord) {
	n.RoutePath = c
}
func (n *Supern) SetRoutePoints(c []Coord) {
	n.RoutePoints = c
}
