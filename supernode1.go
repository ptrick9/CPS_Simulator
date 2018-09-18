package main

import (
	"math"
)

//Super node of type 1
//The super node with the minimum distance routing algorithm
type sn_one struct {
	*supern
}

//This function is called every tick
//It adds points to the super node's routePoints and routePath lists
//	and moves it along the path
func (n* sn_one) tick() {
	//If there are still squares to visit the super node should move along
	//	the path
	if len(n.getRoutePoints()) > 1 {
		//If a new point of interest has been added the current path needs
		//	to be updated
		if n.getNumDest() < len(n.getRoutePoints())-1 {
			n.updatePath()
		}
		//The super node should then follow the current path through the grid
		if len(n.getRoutePath()) > 0 {
			n.pathMove()
		}
	}
}

//updatePath is called when the number of destinations is less
//	than the number of points in the routePoints list
//This means a new point of interest was added and needs to be
//	visited
//This function adds points to the routePath list in order for
//	for the super node to follow, it also inserts points into the
//	routePath if a new point appears and qualifies for insertion
func (n* sn_one) updatePath(){
	//If numDestinations is 0 then there are currently no points for the super
	//	node to visit
	//That means that all the points in the nodePoints list can be visited without
	//	needing to alter the current routePath
	if n.numDestinations == 0 {

		//Creates a list of distances between points in the routePoints list
		//The list contains Path types which contain the start and end index
		//	of the path in routePoints and the distance between them
		dist_list := make([]Path, 0)

		//Calculates the distance between each possible connection of points
		for x := 0; x < (len(n.routePoints))-1; x++ {
			for y := x + 1; y < (len(n.routePoints)); y++ {
				sqrt := math.Pow(float64((n.routePoints[x].x)-(n.routePoints[y].x)), 2.0) + math.Pow(float64((n.routePoints[x].y)-(n.routePoints[y].y)), 2.0)
				dist_list = append(dist_list, Path{x, y, math.Sqrt(sqrt)})
			}
		}

		//Finds the shortest distance in the dist_list
		//Also saves the start and end points of the shortest path
		end := 0

		//Makes a copy of the routePoints list to edit and rearrange, maintaining
		//	the original order of the routePoints list to iterate through
		cpy := make([]Coord, len(n.routePoints))
		cpy[0] = n.routePoints[0]

		//Loops through the list of Paths between points of interest
		//Finds the shortest distance between where the super node's current starting
		//Saves the destination index to be the startinf point for the next iteration
		for len(dist_list) > 0 {
			start := end
			lowdist := float64(maxX*maxY);//1000.0
			for _, d := range dist_list {
				if (d.x == start) {
					if (d.dist < lowdist) {
						lowdist = d.dist
						end = d.y
					}
				} else if (d.y == start) {
					if (d.dist < lowdist) {
						lowdist = d.dist
						end = d.x
					}
				}
			}
			//Adds the points of the path to the routePath list using the
			//	route function
			n.routePath = append(n.routePath, aStar(n.routePoints[start], n.routePoints[end])...)

			//Once a point of interest's path is added to the routePath
			//	the number of destination can be increased
			n.numDestinations++

			//The end point of the created path is placed in the correct location in the
			//	routePoints list
			cpy[n.numDestinations] = n.routePoints[end]

			//Removes every Path that contains the point that was just departed
			//This is to prevent the Path from never looping back to
			//	the previous point
			w := 0
			for w < len(dist_list) {
				d := dist_list[w]
				if (d.x == start) || (d.y == start) {
					dist_list = remove_index(dist_list, w)
				} else {
					w++
				}
			}
		}
		copy(n.routePoints, cpy)

		//If numDestinations is not 0 then the current routePath needs to be altered
	}else {
		//Loops through the points added to the routePoints list that are not accounted
		//	for by the super nodes numDestinations
		for i := n.numDestinations + 1; i < len(n.routePoints); i++ {

			notAdded := true
			ind := 0
			oldNumDest := n.numDestinations

			//This loop goes from the beginning of the routePoints list
			//	until the end of the currently accounted for points
			for (ind < oldNumDest) && notAdded {

				//These distances are the distances from the currenlty selected point to the
				//	next point in the path, the distance from the newly added point to the next
				//	point in the path, and the distance from the currently selected point to the
				//	new point
				dist1 := math.
					Sqrt(math.
						Pow(float64((n.routePoints[ind].x)-(n.routePoints[ind+1].x)), 2.0) + math.
						Pow(float64((n.routePoints[ind].y)-(n.routePoints[ind+1].y)), 2.0))
				dist2 := math.
					Sqrt(math.
						Pow(float64((n.routePoints[i].x)-(n.routePoints[ind+1].x)), 2.0) + math.
						Pow(float64((n.routePoints[i].y)-(n.routePoints[ind+1].y)), 2.0))
				dist3 := math.
					Sqrt(math.
						Pow(float64((n.routePoints[ind].x)-(n.routePoints[i].x)), 2.0) + math.
						Pow(float64((n.routePoints[ind].y)-(n.routePoints[i].y)), 2.0))

				//If the new point is closer to the current point than the next point AND
				//	is closer to the next point than the current point, it is added to the
				//	routePoints list in between them and the path connecting the current point
				//	to the new point and new point to the next point is added in the place of
				//	the path between the current point and the next point
				if (dist2 < dist1) && (dist3 < dist1) {

					notAdded = false

					start_ind := 0
					end_ind := 0
					start_found := false
					end_found := false

					//This for loop finds the location of the points in the routePoints
					//	list in the routePath list
					//If the new point is being placed in between the current location
					//	and the first destination point, the start_ind must be manually
					//	set because the current location will never be found in the
					//	routePath list
					if ind == 0 {
						start_found = true
					}
					for j := 0; (j < len(n.routePath) && !end_found); j++ {
						if (n.routePath[j].x == n.routePoints[ind].x) && (n.routePath[j].y == n.routePoints[ind].y) && !start_found{
							start_ind = j
							start_found = true
						}else if (n.routePath[j].x == n.routePoints[ind+1].x) && (n.routePath[j].y == n.routePoints[ind+1].y) && start_found {
							end_ind = j
							end_found = true
						}
					}

					//The route function is called twice and stored in a new array
					//This separate array appends the two paths to make one path from
					//	each routePath point to the new point
					arr := make([]Coord, 0)
					arr = aStar(n.routePoints[ind], n.routePoints[i])
					arr = append(arr, aStar(n.routePoints[i], n.routePoints[ind+1])...)

					if ind != 0{
						start_ind += 1
					}

					//Removes the current path between two points
					n.routePath = remove_range(n.routePath, start_ind, end_ind)

					//Replaces that path with a a path that connects the two points
					//	by going through the newly added point
					n.routePath = insert_array(n.routePath, arr, start_ind)

					//The new point is inserted into the routePoints list in between
					//	the two points
					n.routePoints = remove_and_insert(n.routePoints, i, ind+1)

					//Since a new point has been added to the routePath, the number
					//	of destinations is increased
					n.numDestinations++


					//Otherwise, the next point becomes the current point and its next becomes
					//	the new next
				}else{
					ind++
				}
			}
			//If the new point is not added anywhere in between the currently routed points
			//	it is added at the end of the current path
			if notAdded {
				n.routePath = append(n.routePath, aStar(n.routePoints[i-1], n.routePoints[i])...)
				n.numDestinations++
			}
		}
	}
}

//Adds a routePoint to the super node's routePoints
func (n* sn_one) addRoutePoint(c Coord) {
	n.routePoints = append(n.routePoints, c)
}