package cps

import "fmt"

//Super node of type 1
//The super node with the minimum distance routing algorithm
type Sn_zero struct {
	P *Params
	R *RegionParams
	*Supern
}

//This function is called every tick
//It adds points to the super node's routePoints and routePath
//	lists and moves the super node if necessary
func (n *Sn_zero) Tick() {
	//If there are new points that are not currently destinations the route
	//	needs to be updated
	/*if len(n.getRoutePoints()) > n.getNumDest() {
		n.updatePath()
	}*/
	//If there are points left in the route the super node must move
	//	along the path
	if len(n.GetRoutePath()) > 0 {
		n.PathMove()
	}
}

//This super node just adds the newest point at the end of the current path
func (n *Sn_zero) UpdatePath() {

	if n.P.RegionRouting {
		//Loops through all the points in the routePoints list that are currently not
		//	destinations
		for i := n.NumDestinations; i < len(n.RoutePoints)-1; i++ {

			//Adds the points of the path to the routePath list using the
			//	route function
			newPath := GetPath(n.RoutePoints[i], n.RoutePoints[i+1], n.R)
			n.RoutePath = append(n.RoutePath, newPath...)

			//Once a point of interest's path is added to the routePath
			//	the number of destination can be increased
			n.NumDestinations++
		}
	} else {
		//Loops through all the points in the routePoints list that are currently not
		//	destinations
		for i := n.NumDestinations; i < len(n.RoutePoints)-1; i++ {

			//Adds the points of the path to the routePath list using the
			//	route function
			fmt.Println("\ncalling aStar")
			newPath := AStar(n.RoutePoints[i], n.RoutePoints[i+1], n.P)
			n.RoutePath = append(n.RoutePath, newPath...)

			//Once a point of interest's path is added to the routePath
			//	the number of destination can be increased
			n.NumDestinations++
		}
	}
}

//Adds a routePoint to the super node's routePoints
func (n *Sn_zero) AddRoutePoint(c Coord) {
	n.RoutePoints = append(n.RoutePoints, c)
	n.UpdatePath()
}
