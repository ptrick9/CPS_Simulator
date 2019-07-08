package cps

import (
	"math"
)

//Super node of type 1
//The super node with the minimum distance routing algorithm
type Sn_one struct {
	*Supern
}

//This function is called every tick
//It adds points to the super node's RoutePoints and RoutePath lists
//	and moves it along the path
func (n *Sn_one) Tick() {
	//If there are still squares to visit the super node should move along
	//	the path
	if len(n.GetRoutePoints()) > 1 {
		//If a new point of interest has been added the current path needs
		//	to be updated
		if n.GetNumDest() < len(n.GetRoutePoints())-1 {
			n.UpdatePath()
		}
		//The super node should then follow the current path through the grid
		if len(n.GetRoutePath()) > 0 {
			n.PathMove()
		}
	}
}

//updatePath is called when the number of destinations is less
//	than the number of points in the RoutePoints list
//This means a new point of interest was added and needs to be
//	visited
//This function adds points to the RoutePath list in order for
//	for the super node to follow, it also inserts points into the
//	RoutePath if a new point appears and qualifies for insertion
func (n *Sn_one) UpdatePath() {
	//If NumDestinations is 0 then there are currently no points for the super
	//	node to visit
	//That means that all the points in the nodePoints list can be visited without
	//	needing to alter the current RoutePath
	if n.NumDestinations == 0 {

		//Creates a list of distances between points in the RoutePoints list
		//The list contains Path types which contain the start and end Index
		//	of the path in RoutePoints and the distance between them
		dist_list := make([]Path, 0)

		//Calculates the distance between each possible connection of points
		for x := 0; x < (len(n.RoutePoints))-1; x++ {
			for y := x + 1; y < (len(n.RoutePoints)); y++ {
				sqrt := math.Pow(float64((n.RoutePoints[x].X)-(n.RoutePoints[y].X)), 2.0) + math.Pow(float64((n.RoutePoints[x].Y)-(n.RoutePoints[y].Y)), 2.0)
				dist_list = append(dist_list, Path{x, y, math.Sqrt(sqrt)})
			}
		}

		//Finds the shortest distance in the dist_list
		//Also saves the start and end points of the shortest path
		end := 0

		//Makes a copy of the RoutePoints list to edit and rearrange, maintaining
		//	the original order of the RoutePoints list to iterate through
		cpy := make([]Coord, len(n.RoutePoints))
		cpy[0] = n.RoutePoints[0]

		//Loops through the list of Paths between points of interest
		//Finds the shortest distance between where the super node's current starting
		//Saves the destination Index to be the startinf point for the next iteration
		for len(dist_list) > 0 {
			start := end
			lowdist := float64(n.P.MaxX * n.P.MaxY) //1000.0
			for _, d := range dist_list {
				if d.X == start {
					if d.Dist < lowdist {
						lowdist = d.Dist
						end = d.Y
					}
				} else if d.Y == start {
					if d.Dist < lowdist {
						lowdist = d.Dist
						end = d.X
					}
				}
			}
			//Adds the points of the path to the RoutePath list using the
			//	Route function
			n.RoutePath = append(n.RoutePath, AStar(n.RoutePoints[start], n.RoutePoints[end], n.P)...)

			//Once a point of interest's path is added to the RoutePath
			//	the number of destination can be increased
			n.NumDestinations++

			//The end point of the created path is placed in the correct location in the
			//	RoutePoints list
			cpy[n.NumDestinations] = n.RoutePoints[end]

			//Removes every Path that contains the point that was just departed
			//This is to prevent the Path from never looping back to
			//	the previous point
			w := 0
			for w < len(dist_list) {
				d := dist_list[w]
				if (d.X == start) || (d.Y == start) {
					dist_list = Remove_index(dist_list, w)
				} else {
					w++
				}
			}
		}
		copy(n.RoutePoints, cpy)

		//If NumDestinations is not 0 then the current RoutePath needs to be altered
	} else {
		//Loops through the points added to the RoutePoints list that are not accounted
		//	for by the super nodes NumDestinations
		for i := n.NumDestinations + 1; i < len(n.RoutePoints); i++ {

			notAdded := true
			ind := 0
			oldNumDest := n.NumDestinations

			//This loop goes from the beginning of the RoutePoints list
			//	until the end of the currently accounted for points
			for (ind < oldNumDest) && notAdded {

				//These distances are the distances from the currenlty selected point to the
				//	next point in the path, the distance from the newly added point to the next
				//	point in the path, and the distance from the currently selected point to the
				//	new point
				dist1 := math.
					Sqrt(math.
						Pow(float64((n.RoutePoints[ind].X)-(n.RoutePoints[ind+1].X)), 2.0) + math.
						Pow(float64((n.RoutePoints[ind].Y)-(n.RoutePoints[ind+1].Y)), 2.0))
				dist2 := math.
					Sqrt(math.
						Pow(float64((n.RoutePoints[i].X)-(n.RoutePoints[ind+1].X)), 2.0) + math.
						Pow(float64((n.RoutePoints[i].Y)-(n.RoutePoints[ind+1].Y)), 2.0))
				dist3 := math.
					Sqrt(math.
						Pow(float64((n.RoutePoints[ind].X)-(n.RoutePoints[i].X)), 2.0) + math.
						Pow(float64((n.RoutePoints[ind].Y)-(n.RoutePoints[i].Y)), 2.0))

				//If the new point is closer to the current point than the next point AND
				//	is closer to the next point than the current point, it is added to the
				//	RoutePoints list in between them and the path connecting the current point
				//	to the new point and new point to the next point is added in the place of
				//	the path between the current point and the next point
				if (dist2 < dist1) && (dist3 < dist1) {

					notAdded = false

					start_ind := 0
					end_ind := 0
					start_found := false
					end_found := false

					//This for loop finds the location of the points in the RoutePoints
					//	list in the RoutePath list
					//If the new point is being placed in between the current location
					//	and the first destination point, the start_ind must be manually
					//	set because the current location will never be found in the
					//	RoutePath list
					if ind == 0 {
						start_found = true
					}
					for j := 0; j < len(n.RoutePath) && !end_found; j++ {
						if (n.RoutePath[j].X == n.RoutePoints[ind].X) && (n.RoutePath[j].Y == n.RoutePoints[ind].Y) && !start_found {
							start_ind = j
							start_found = true
						} else if (n.RoutePath[j].X == n.RoutePoints[ind+1].X) && (n.RoutePath[j].Y == n.RoutePoints[ind+1].Y) && start_found {
							end_ind = j
							end_found = true
						}
					}

					//The Route function is called twice and stored in a new array
					//This separate array appends the two paths to make one path from
					//	each RoutePath point to the new point
					arr := make([]Coord, 0)
					arr = AStar(n.RoutePoints[ind], n.RoutePoints[i], n.P)
					arr = append(arr, AStar(n.RoutePoints[i], n.RoutePoints[ind+1], n.P)...)

					if ind != 0 {
						start_ind += 1
					}

					//Removes the current path between two points
					n.RoutePath = Remove_range(n.RoutePath, start_ind, end_ind)

					//Replaces that path with a a path that connects the two points
					//	by going through the newly added point
					n.RoutePath = Insert_array(n.RoutePath, arr, start_ind)

					//The new point is inserted into the RoutePoints list in between
					//	the two points
					n.RoutePoints = Remove_and_insert(n.RoutePoints, i, ind+1)

					//Since a new point has been added to the RoutePath, the number
					//	of destinations is increased
					n.NumDestinations++

					//Otherwise, the next point becomes the current point and its next becomes
					//	the new next
				} else {
					ind++
				}
			}
			//If the new point is not added anywhere in between the currently Routed points
			//	it is added at the end of the current path
			if notAdded {
				n.RoutePath = append(n.RoutePath, AStar(n.RoutePoints[i-1], n.RoutePoints[i], n.P)...)
				n.NumDestinations++
			}
		}
	}
}

//Adds a RoutePoint to the super node's RoutePoints
func (n *Sn_one) AddRoutePoint(c Coord) {
	n.RoutePoints = append(n.RoutePoints, c)
}
