package cps

import (
	"math"
)

//This interface is implemented by super node of type 2
type SuperNode2Impl interface {
	Calc_triangle_tot(Region) (int, int)
	Add_route_point([][]Square, Coord)
	WthinRegion(Region)
}

//Super node of type 2
//The super node with the regional routing algorithm
type Sn_two struct {
	*Supern

	RegionList     []Region
	UsedRegionList []Region
}

//The struct region defines a rectangular area inside the
//	grid
type Region struct {
	Center, Edge Coord
	Points       []Coord
}

//This function is called every tick
//It adds Points to the super nodes routePoints and routePath lists
//	and moves it along those paths
func (n *Sn_two) Tick() {
	//Increments the Time of all the Points in the allPoints list
	n.IncAllPoints()

	//If there is a path for the super node to follow it should follow its path
	//Otherwise it should do one of two things:
	//		Move back to the Center
	//		OR if already in the Center, plot a new path
	if len(n.RoutePath) > 0 {
		n.PathMove()
	} else {
		if (int(n.X) == n.Center.X) && (int(n.Y) == n.Center.Y) {
			n.UpdatePath()
		} else {
			n.CentMove()
		}
	}
}

//This function updates the path for the super node to follow
//updatePath is only called when the super node is in the Center of the grid
//	with no more Points to visit
//It searches the four quadrants of the grid and finds the one with the most Points
//	of interest inside of it or the one with the oldest point of interest inside it,
//	it then plots a route that takes the super node to the Edge of that region and
//	back, intercepting Points along the way
//Once a region has been visited it will not be visited again until all of the
//	region have been visited or all the other regions are currently empty
func (n *Sn_two) UpdatePath() {

	//If all the regions have been traversed then all regions may be looked at again
	if len(n.RegionList) <= 0 {
		n.RegionList = n.UsedRegionList
		n.UsedRegionList = make([]Region, 0)
	}

	//This is the Index of the region with the most Points in it
	//This region is then used to plot a path for the super node
	high_ind := n.OldestPointRegion()

	//If none of the remaining regions contain any Points, all the regions
	//	are looked at again
	if high_ind == -1 {
		n.RegionList = append(n.RegionList, n.UsedRegionList...)
		n.UsedRegionList = make([]Region, 0)

		//Otherwise a path is created through the region found to have the greatest
		//	number of Points of interest
	} else {
		//Adds the Points inside the Region to that Region's Points list
		n.AddToRegion(high_ind)

		//Initializes several lists to store the return parameters from the function
		//	calc_triangle_tot
		//These lists will be used to plot a path out and back through the region
		lower_list := make([]Coord, 0)
		upper_list := make([]Coord, 0)
		start := 0

		//Finds the number of Points in the lower and upper halves of the region
		//This allows the super node to visit Points on the way out and back towards the
		//	Center, saving Time on the return trip
		lower_list, upper_list = n.Calc_triangle_tot(n.RegionList[high_ind])

		start_list := make([]Coord, 0)
		end_list := make([]Coord, 0)

		//Ensures the start_list is the largest of the two lists, end_list becomes the other
		if len(lower_list) > len(upper_list) {
			start_list = lower_list
			end_list = upper_list
		} else {
			start_list = upper_list
			end_list = lower_list
		}

		//Initializes a list to store lengths of Paths between possible Points to visit
		dist_list := make([]Path, 0)

		//Calculates the distance between each point and the Center
		counter := 0
		for _, c := range start_list {
			sqrt := math.Pow(float64((c.X)-(n.RegionList[high_ind].Center.X)), 2.0) + math.Pow(float64((c.Y)-(n.RegionList[high_ind].Center.Y)), 2.0)
			dist_list = append(dist_list, Path{0, counter, math.Sqrt(sqrt)})
			counter++
		}

		//This loop plots the path from the super node's current location to the
		//	Points it must visit in this part of the region
		for len(dist_list) > 0 {
			lowind := 0
			lowdist := float64(n.P.MaxX * n.P.MaxY) //1000.0
			for _, d := range dist_list {
				if d.Dist < lowdist {
					lowdist = d.Dist
					lowind = d.Y
				}
			}
			//Adds the selected point to the routePoints list
			n.RoutePoints = append(n.RoutePoints, start_list[lowind])

			//Finds the path to that selected point, then adds that path to the routePath
			//	list, as well as adding the point itself
			n.RoutePath = append(n.RoutePath, AStar(n.RoutePoints[start], n.RoutePoints[start+1], n.P)...)

			//Since the super node must visit this point the numDestinations is increased
			n.NumDestinations++

			//The selected node is removed from the dist_list allowing the super node to
			//	travel to the remaining Points in this half of the region
			w := 0
			for w < len(dist_list) {
				d := dist_list[w]
				if d.Y == lowind {
					dist_list = Remove_index(dist_list, w)
				} else {
					w++
				}
			}
			//The Index to add things into the routePoints list is incremented each iteration
			start++
		}

		//If the other half of the region does not have any Points the super
		//	node is directed back to its Center point
		if len(end_list) > 0 {
			dist_list = make([]Path, 0)

			//Calculates the distance between each point and the Center
			counter = 0
			for _, c := range end_list {
				sqrt := math.Pow(float64((c.X)-(n.RegionList[high_ind].Center.X)), 2.0) + math.Pow(float64((c.Y)-(n.RegionList[high_ind].Center.Y)), 2.0)
				dist_list = append(dist_list, Path{0, counter, math.Sqrt(sqrt)})
				counter++
			}

			//Finds the point furthest from the Center in the other half of the region
			//This point will be connected to the path on the other side of the region
			longdist := -1.0
			longind := -1
			for _, d := range dist_list {
				if d.Dist > longdist {
					longdist = d.Dist
					longind = d.Y
				}
			}

			//Adds the point to the routePoints list
			n.RoutePoints = append(n.RoutePoints, end_list[longind])

			//Increases the numDestinations
			n.NumDestinations++

			//Adds the path between that point and the last point of the previous path
			//	to the routePath list
			n.RoutePath = append(n.RoutePath, AStar(n.RoutePoints[start], n.RoutePoints[start+1], n.P)...)

			//The Index to add things into the routePoints list is incremented each iteration
			start++

			//The selected node is removed from the dist_list allowing the super node to
			//	travel to the remaining Points in this half of the region
			w := 0
			for w < len(dist_list) {
				d := dist_list[w]
				if d.Y == longind {
					dist_list = Remove_index(dist_list, w)
				} else {
					w++
				}
			}

			//This loop plots the path from the super node's current location to the
			//	Points it must visit in this part of the region
			for len(dist_list) > 0 {
				longind := 0
				longdist := -1.0
				for _, d := range dist_list {
					if d.Dist > longdist {
						longdist = d.Dist
						longind = d.Y
					}
				}

				//Adds the selected point to the routePoints list
				n.RoutePoints = append(n.RoutePoints, end_list[longind])

				//Finds the path to that selected point, then adds that path to the routePath
				//	list, as well as adding the point itself
				n.RoutePath = append(n.RoutePath, AStar(n.RoutePoints[start], n.RoutePoints[start+1], n.P)...)

				//Since the super node must visit this point the numDestinations is increased
				n.NumDestinations++

				//The selected node is removed from the dist_list allowing the super node to
				//	travel to the remaining Points in this half of the region
				w := 0
				for w < len(dist_list) {
					d := dist_list[w]
					if d.Y == longind {
						dist_list = Remove_index(dist_list, w)
					} else {
						w++
					}
				}
				//The Index to add things into the routePoints list is incremented each iteration
				start++
			}
		}

		//The Region's Points list must be emptied or else it's length will still include
		//	previously visited Points
		n.RegionList[high_ind].Points = make([]Coord, 0)

		//Adding the currently visited region to the usedRegionList
		//This ensures it won't be visited right away again
		n.UsedRegionList = append(n.UsedRegionList, n.RegionList[high_ind])

		//Removes the currently visited region from the RegionList
		n.RegionList = Remove_region(n.RegionList, high_ind)
	}
}

//This function returns the Index of the Region with the most Points of interest
//	inside of it
func (n *Sn_two) MostPointsRegion() int {
	high_Points := 0
	high_ind := -1
	for r, _ := range n.RegionList {
		if n.TotalInRegion(n.RegionList[r]) > high_Points {
			high_Points = n.TotalInRegion(n.RegionList[r])
			high_ind = r
		}
	}
	return high_ind
}

//This function returns the Index of the Region with the oldest point of interest
//	inside of it
func (n *Sn_two) OldestPointRegion() int {
	point_age := -1
	high_ind := -1
	for r, _ := range n.RegionList {
		if n.OldestInRegion(n.RegionList[r]) > point_age {
			point_age = n.OldestInRegion(n.RegionList[r])
			high_ind = r
		}
	}
	return high_ind
}

//The function returns the arr with the element at Index n removed
func Remove_region(arr []Region, n int) []Region {
	return arr[:n+copy(arr[n:], arr[n+1:])]
}

//This function calculates the total number of Points of interest inside
//	a regions lower and upper triangles
//The region can be split into triangles either top-right to bottom-left
//	or top-left to bottom-right
func (n *Sn_two) Calc_triangle_tot(r Region) ([]Coord, []Coord) {
	lower_coords := make([]Coord, 0)
	upper_coords := make([]Coord, 0)

	//SE quadrant and NW quadrant
	if ((r.Center.X <= r.Edge.X) && (r.Center.Y <= r.Edge.Y)) || ((r.Center.X >= r.Edge.X) && (r.Center.Y >= r.Edge.Y)) {
		//This loop determines whether the Points of interest within this Region
		//	are in the upp triangle or the lower triangle
		for p, _ := range r.Points {
			if r.Points[p].X <= r.Points[p].Y {
				lower_coords = append(lower_coords, r.Points[p])
			} else {
				upper_coords = append(upper_coords, r.Points[p])
			}
		}
		//Sw quadrant and NE quadrant
	} else {
		//This loop determines whether the Points of interest within this Region
		//	are in the upp triangle or the lower triangle
		for p, _ := range r.Points {
			if (r.Points[p].X + r.Points[p].Y) > (n.P.MaxX / 2) {
				lower_coords = append(lower_coords, r.Points[p])
			} else {
				upper_coords = append(upper_coords, r.Points[p])
			}
		}
	}
	return lower_coords, upper_coords
}

//This function takes a Region and returns the total number of Points
//	of interest inside the Region
func (n *Sn_two) TotalInRegion(reg Region) int {
	tot := 0
	for c, _ := range n.AllPoints {
		if n.AllPoints[c].IsWithinRegion(reg) {
			tot++
		}
	}
	return tot
}

//This function takes a Region and returns the oldest point of interest
//	inside that Region
func (n *Sn_two) OldestInRegion(reg Region) int {
	age := 0
	for c, _ := range n.AllPoints {
		if n.AllPoints[c].IsWithinRegion(reg) {
			if n.AllPoints[c].Time > age {
				age = n.AllPoints[c].Time
			}
		}
	}
	return age
}

//This function returns true if the Coord that called it is inside
//	the specified region
func (c Coord) IsWithinRegion(reg Region) bool {
	inside := false
	//Boundary conditions
	if c.X > reg.Edge.X {
		if reg.Center.Y > reg.Edge.Y {
			if (c.X >= reg.Edge.X) && (c.X <= reg.Center.X) && (c.Y >= reg.Edge.Y) && (c.Y <= reg.Center.Y) {
				inside = true
			}
		} else {
			if (c.X >= reg.Edge.X) && (c.X <= reg.Center.X) && (c.Y <= reg.Edge.Y) && (c.Y >= reg.Center.Y) {
				inside = true
			}
		}
	} else {
		if reg.Center.Y > reg.Edge.Y {
			if (c.X <= reg.Edge.X) && (c.X >= reg.Center.X) && (c.Y >= reg.Edge.Y) && (c.Y <= reg.Center.Y) {
				inside = true
			}
		} else {
			if (c.X <= reg.Edge.X) && (c.X >= reg.Center.X) && (c.Y <= reg.Edge.Y) && (c.Y >= reg.Center.Y) {
				inside = true
			}
		}
	}
	return inside
}

//This function removes Points from the allPoints list and adds them to the
//	specified region's (n.RegionList[r]) Points list
func (n *Sn_two) AddToRegion(r int) {
	c := 0
	for c < len(n.AllPoints) {
		//If the Coord is within the Region it is appended to that Region's Points list
		//The Coord's Time is updated and then it is removed from the allPoints list
		if n.AllPoints[c].IsWithinRegion(n.RegionList[r]) {
			n.RegionList[r].Points = append(n.RegionList[r].Points, n.AllPoints[c])
			n.RegionList[r].Points[len(n.RegionList[r].Points)-1].Time = n.AllPoints[c].Time
			n.AllPoints = Remove_coord_index(n.AllPoints, c)
		} else {
			c++
		}
	}
}

//This function removed the Coord element at the specified Index n
func Remove_coord_index(arr []Coord, n int) []Coord {
	return arr[:n+copy(arr[n:], arr[n+1:])]
}

//This function adds a point of interest to the allPoints list
//This list is traversed and the point's regions are decided when
// 	updating the path
func (n *Sn_two) AddRoutePoint(c Coord) {
	n.AllPoints = append(n.AllPoints, c)
}

func (n *Sn_two) AddRoutePointUrgent(c Coord) {
	n.AllPoints = append(n.AllPoints, c)
}
