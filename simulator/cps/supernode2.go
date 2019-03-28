package cps

import (
	"math"
)

//This interface is implemented by super node of type 2
type SuperNode2Impl interface {
	calc_triangle_tot(Region) (int, int)
	add_route_point([][]Square, Coord)
	withinRegion(Region)
}

//Super node of type 2
//The super node with the regional routing algorithm
type sn_two struct {
	*supern

	regionList []Region
	usedRegionList []Region
}

//The struct region defines a rectangular area inside the
//	grid
type Region struct {
	center, edge Coord
	points []Coord
}

//This function is called every tick
//It adds points to the super nodes routePoints and routePath lists
//	and moves it along those paths
func (n* sn_two) tick() {
	//Increments the time of all the points in the allPoints list
	n.incAllPoints()

	//If there is a path for the super node to follow it should follow its path
	//Otherwise it should do one of two things:
	//		Move back to the center
	//		OR if already in the center, plot a new path
	if len(n.routePath) > 0 {
		n.pathMove()
	}else {
		if (n.x == n.center.x) && (n.y == n.center.y) {
			n.updatePath()
		}else {
			n.centMove()
		}
	}
}

//This function updates the path for the super node to follow
//updatePath is only called when the super node is in the center of the grid
//	with no more points to visit
//It searches the four quadrants of the grid and finds the one with the most points
//	of interest inside of it or the one with the oldest point of interest inside it,
//	it then plots a route that takes the super node to the edge of that region and
//	back, intercepting points along the way
//Once a region has been visited it will not be visited again until all of the
//	region have been visited or all the other regions are currently empty
func (n* sn_two) updatePath(){

	//If all the regions have been traversed then all regions may be looked at again
	if len(n.regionList) <= 0{
		n.regionList = n.usedRegionList
		n.usedRegionList = make([]Region, 0)
	}

	//This is the index of the region with the most points in it
	//This region is then used to plot a path for the super node
	high_ind := n.oldestPointRegion()

	//If none of the remaining regions contain any points, all the regions
	//	are looked at again
	if high_ind == -1 {
		n.regionList = append(n.regionList, n.usedRegionList...)
		n.usedRegionList = make([]Region, 0)

		//Otherwise a path is created through the region found to have the greatest
		//	number of points of interest
	}else {
		//Adds the points inside the Region to that Region's points list
		n.addToRegion(high_ind)

		//Initializes several lists to store the return parameters from the function
		//	calc_triangle_tot
		//These lists will be used to plot a path out and back through the region
		lower_list := make([]Coord, 0)
		upper_list := make([]Coord, 0)
		start := 0

		//Finds the number of points in the lower and upper halves of the region
		//This allows the super node to visit points on the way out and back towards the
		//	center, saving time on the return trip
		lower_list, upper_list = n.calc_triangle_tot(n.regionList[high_ind])

		start_list := make([]Coord, 0)
		end_list := make([]Coord, 0)

		//Ensures the start_list is the largest of the two lists, end_list becomes the other
		if len(lower_list) > len(upper_list) {
			start_list = lower_list
			end_list = upper_list
		}else{
			start_list = upper_list
			end_list = lower_list
		}

		//Initializes a list to store lengths of Paths between possible points to visit
		dist_list := make([]Path, 0)

		//Calculates the distance between each point and the center
		counter := 0
		for _, c := range start_list {
			sqrt := math.Pow(float64((c.x)-(n.regionList[high_ind].center.x)), 2.0) + math.Pow(float64((c.y)-(n.regionList[high_ind].center.y)), 2.0)
			dist_list = append(dist_list, Path{0, counter, math.Sqrt(sqrt)})
			counter++
		}

		//This loop plots the path from the super node's current location to the
		//	points it must visit in this part of the region
		for len(dist_list) > 0 {
			lowind := 0
			lowdist := float64(maxX*maxY);//1000.0
			for _, d := range dist_list {
				if (d.dist < lowdist) {
					lowdist = d.dist
					lowind = d.y
				}
			}
			//Adds the selected point to the routePoints list
			n.routePoints = append(n.routePoints, start_list[lowind])

			//Finds the path to that selected point, then adds that path to the routePath
			//	list, as well as adding the point itself
			n.routePath = append(n.routePath, aStar(n.routePoints[start], n.routePoints[start+1])...)

			//Since the super node must visit this point the numDestinations is increased
			n.numDestinations++

			//The selected node is removed from the dist_list allowing the super node to
			//	travel to the remaining points in this half of the region
			w := 0
			for w < len(dist_list) {
				d := dist_list[w]
				if (d.y == lowind) {
					dist_list = remove_index(dist_list, w)
				} else {
					w++
				}
			}
			//The index to add things into the routePoints list is incremented each iteration
			start++
		}

		//If the other half of the region does not have any points the super
		//	node is directed back to its center point
		if (len(end_list) > 0) {
			dist_list = make([]Path, 0)

			//Calculates the distance between each point and the center
			counter = 0
			for _, c := range end_list {
				sqrt := math.Pow(float64((c.x)-(n.regionList[high_ind].center.x)), 2.0) + math.Pow(float64((c.y)-(n.regionList[high_ind].center.y)), 2.0)
				dist_list = append(dist_list, Path{0, counter, math.Sqrt(sqrt)})
				counter++
			}

			//Finds the point furthest from the center in the other half of the region
			//This point will be connected to the path on the other side of the region
			longdist := -1.0
			longind := -1
			for _, d := range dist_list {
				if (d.dist > longdist) {
					longdist = d.dist
					longind = d.y
				}
			}

			//Adds the point to the routePoints list
			n.routePoints = append(n.routePoints, end_list[longind])

			//Increases the numDestinations
			n.numDestinations++

			//Adds the path between that point and the last point of the previous path
			//	to the routePath list
			n.routePath = append(n.routePath, aStar(n.routePoints[start], n.routePoints[start+1])...)

			//The index to add things into the routePoints list is incremented each iteration
			start++

			//The selected node is removed from the dist_list allowing the super node to
			//	travel to the remaining points in this half of the region
			w := 0
			for w < len(dist_list) {
				d := dist_list[w]
				if (d.y == longind) {
					dist_list = remove_index(dist_list, w)
				} else {
					w++
				}
			}

			//This loop plots the path from the super node's current location to the
			//	points it must visit in this part of the region
			for len(dist_list) > 0 {
				longind := 0
				longdist := -1.0
				for _, d := range dist_list {
					if (d.dist > longdist) {
						longdist = d.dist
						longind = d.y
					}
				}

				//Adds the selected point to the routePoints list
				n.routePoints = append(n.routePoints, end_list[longind])


				//Finds the path to that selected point, then adds that path to the routePath
				//	list, as well as adding the point itself
				n.routePath = append(n.routePath,aStar(n.routePoints[start], n.routePoints[start+1])...)

				//Since the super node must visit this point the numDestinations is increased
				n.numDestinations++

				//The selected node is removed from the dist_list allowing the super node to
				//	travel to the remaining points in this half of the region
				w := 0
				for w < len(dist_list) {
					d := dist_list[w]
					if (d.y == longind) {
						dist_list = remove_index(dist_list, w)
					} else {
						w++
					}
				}
				//The index to add things into the routePoints list is incremented each iteration
				start++
			}
		}

		//The Region's points list must be emptied or else it's length will still include
		//	previously visited points
		n.regionList[high_ind].points = make([]Coord, 0)

		//Adding the currently visited region to the usedRegionList
		//This ensures it won't be visited right away again
		n.usedRegionList = append(n.usedRegionList, n.regionList[high_ind])

		//Removes the currently visited region from the regionList
		n.regionList = remove_region(n.regionList, high_ind)
	}
}

//This function returns the index of the Region with the most points of interest
//	inside of it
func (n* sn_two) mostPointsRegion() int{
	high_points := 0
	high_ind := -1
	for r,_ := range n.regionList {
		if n.totalInRegion(n.regionList[r]) > high_points {
			high_points = n.totalInRegion(n.regionList[r])
			high_ind = r
		}
	}
	return high_ind
}

//This function returns the index of the Region with the oldest point of interest
//	inside of it
func (n* sn_two) oldestPointRegion() int {
	point_age := -1
	high_ind := -1
	for r, _ := range n.regionList {
		if n.oldestInRegion(n.regionList[r]) > point_age {
			point_age = n.oldestInRegion(n.regionList[r])
			high_ind  = r
		}
	}
	return high_ind
}

//The function returns the arr with the element at index n removed
func remove_region(arr []Region, n int) []Region {
	return arr[:n+copy(arr[n:], arr[n+1:])]
}

//This function calculates the total number of points of interest inside
//	a regions lower and upper triangles
//The region can be split into triangles either top-right to bottom-left
//	or top-left to bottom-right
func (n *sn_two) calc_triangle_tot(r Region) ([]Coord, []Coord) {
	lower_coords := make([]Coord, 0)
	upper_coords := make([]Coord, 0)

	//SE quadrant and NW quadrant
	if ((r.center.x <= r.edge.x) && (r.center.y <= r.edge.y)) || ((r.center.x >= r.edge.x) && (r.center.y >= r.edge.y)){
		//This loop determines whether the points of interest within this Region
		//	are in the upp triangle or the lower triangle
		for p,_ := range r.points {
			if (r.points[p].x <= r.points[p].y) {
				lower_coords = append(lower_coords, r.points[p])
			}else {
				upper_coords = append(upper_coords, r.points[p])
			}
		}
		//Sw quadrant and NE quadrant
	} else {
		//This loop determines whether the points of interest within this Region
		//	are in the upp triangle or the lower triangle
		for p,_ := range r.points {
			if ((r.points[p].x + r.points[p].y) > (maxX / 2)) {
				lower_coords = append(lower_coords, r.points[p])
			}else {
				upper_coords = append(upper_coords, r.points[p])
			}
		}
	}
	return lower_coords, upper_coords
}

//This function takes a Region and returns the total number of points
//	of interest inside the Region
func (n* sn_two) totalInRegion(reg Region) int {
	tot := 0
	for c, _ := range n.allPoints {
		if n.allPoints[c].isWithinRegion(reg) {
			tot++
		}
	}
	return tot
}

//This function takes a Region and returns the oldest point of interest
//	inside that Region
func (n* sn_two) oldestInRegion(reg Region) int {
	age := 0
	for c, _ := range n.allPoints {
		if n.allPoints[c].isWithinRegion(reg) {
			if n.allPoints[c].time > age {
				age = n.allPoints[c].time
			}
		}
	}
	return age
}

//This function returns true if the Coord that called it is inside
//	the specified region
func (c Coord) isWithinRegion(reg Region) bool {
	inside := false
	//Boundary conditions
	if c.x > reg.edge.x {
		if reg.center.y > reg.edge.y {
			if (c.x >= reg.edge.x) && (c.x <=reg.center.x) && (c.y >= reg.edge.y) && (c.y <= reg.center.y) {
				inside = true
			}
		}else {
			if (c.x >= reg.edge.x) && (c.x <= reg.center.x) && (c.y <= reg.edge.y) && (c.y >= reg.center.y){
				inside = true
			}
		}
	}else {
		if reg.center.y > reg.edge.y {
			if (c.x <= reg.edge.x) && (c.x >=reg.center.x) && (c.y >= reg.edge.y) && (c.y <= reg.center.y) {
				inside = true
			}
		}else {
			if (c.x <= reg.edge.x) && (c.x >= reg.center.x) && (c.y <= reg.edge.y) && (c.y >= reg.center.y){
				inside = true
			}
		}
	}
	return inside
}

//This function removes points from the allPoints list and adds them to the
//	specified region's (n.regionList[r]) points list
func (n *sn_two) addToRegion(r int) {
	c := 0
	for c < len(n.allPoints) {
		//If the Coord is within the Region it is appended to that Region's points list
		//The Coord's time is updated and then it is removed from the allPoints list
		if n.allPoints[c].isWithinRegion(n.regionList[r]) {
			n.regionList[r].points = append(n.regionList[r].points, n.allPoints[c])
			n.regionList[r].points[len(n.regionList[r].points)-1].time = n.allPoints[c].time
			n.allPoints = remove_coord_index(n.allPoints, c)
		} else {
			c++
		}
	}
}

//This function removed the Coord element at the specified index n
func remove_coord_index(arr []Coord, n int) []Coord {
	return arr[:n+copy(arr[n:], arr[n+1:])]
}

//This function adds a point of interest to the allPoints list
//This list is traversed and the point's regions are decided when
// 	updating the path
func (n *sn_two) addRoutePoint(c Coord) {
	n.allPoints = append(n.allPoints, c)
}
