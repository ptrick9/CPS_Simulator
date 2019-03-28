package cps

import "fmt"

//Super node of type 1
//The super node with the minimum distance routing algorithm
type sn_zero struct {
	*supern
}

//This function is called every tick
//It adds points to the super node's routePoints and routePath
//	lists and moves the super node if necessary
func (n *sn_zero) tick() {
	//If there are new points that are not currently destinations the route
	//	needs to be updated
	/*if len(n.getRoutePoints()) > n.getNumDest() {
		n.updatePath()
	}*/
	//If there are points left in the route the super node must move
	//	along the path
	if len(n.getRoutePath()) > 0 {
		n.pathMove()
	}
}

//This super node just adds the newest point at the end of the current path
func (n *sn_zero) updatePath() {

	if regionRouting {
		//Loops through all the points in the routePoints list that are currently not
		//	destinations
		for i := n.numDestinations; i < len(n.routePoints)-1; i++ {

			//Adds the points of the path to the routePath list using the
			//	route function
			newPath := getPath(n.routePoints[i], n.routePoints[i+1])
			n.routePath = append(n.routePath, newPath...)

			//Once a point of interest's path is added to the routePath
			//	the number of destination can be increased
			n.numDestinations++
		}
	} else {
		//Loops through all the points in the routePoints list that are currently not
		//	destinations
		for i := n.numDestinations; i < len(n.routePoints)-1; i++ {

			//Adds the points of the path to the routePath list using the
			//	route function
			fmt.Println("\ncalling aStar")
			newPath := aStar(n.routePoints[i], n.routePoints[i+1])
			n.routePath = append(n.routePath, newPath...)

			//Once a point of interest's path is added to the routePath
			//	the number of destination can be increased
			n.numDestinations++
		}
	}
}

//Adds a routePoint to the super node's routePoints
func (n *sn_zero) addRoutePoint(c Coord) {
	n.routePoints = append(n.routePoints, c)
	n.updatePath()
}
