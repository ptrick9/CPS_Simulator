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
			newPath := GetPath(n.RoutePoints[i], n.RoutePoints[i+1], n.R, n.P)
			n.RoutePath = append(n.RoutePath, newPath...)
			//fmt.Printf("Route Path:%v\n", n.RoutePath)
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

func (n *Sn_zero) AddRoutePointUrgent(c Coord) {
	n.RoutePoints = []Coord{ Coord{X: n.X, Y: n.Y}, c}
	n.RoutePath = make([]Coord, 0)
	n.NumDestinations = 0
	n.UpdatePath()
}

//Adds a routePoint to the super node's routePoints
func (n *Sn_zero) AddRoutePoint(c Coord) {
	n.RoutePoints = append(n.RoutePoints, c)
	if len(n.RoutePoints) > 1 {
		newRoute := make([]Coord, 0)
		destinations := make([][]float64, len(n.RoutePoints))
		for i:= range destinations {
			destinations[i] = make([]float64, len(n.RoutePoints))
		}
		for i:= range n.RoutePoints {
			for j:= range n.RoutePoints {
				if i == j || j == 0{
					destinations[i][j] = 1000000.0
				} else {
					destinations[i][j] = Dist(Tuple{n.RoutePoints[i].X,n.RoutePoints[i].Y}, Tuple{n.RoutePoints[j].X,n.RoutePoints[j].Y})
				}
			}
		}
		min := 999999.9
		minIndex := -1
		i:=0
		for len(newRoute) < len(n.RoutePoints) {
			if i == 0 {
				newRoute = append(newRoute, n.RoutePoints[0])
			}
			min = 999999.9
			for j:= 0; j < len(destinations); j++ {
				if destinations[i][j] < min{
					min = destinations[i][j]
					minIndex = j
				}
			}
			//fmt.Printf("Min: %v, MinIndex: %v, Row:%v\n", min, minIndex, destinations[minIndex])

			newRoute = append(newRoute, n.RoutePoints[minIndex])
			i = minIndex
			for k := range destinations {
				destinations[k][minIndex] = 1000000.0
			}

			//Check destinations array
			/*for a:= range destinations {
				for b:= range destinations {
					if destinations[a][b] == 1000000.0 {
						fmt.Print(" INFINITY ")
					} else {
						fmt.Printf(" %08.3f ",destinations[a][b])
					}
				}
				fmt.Println()
			}*/

		}

		//fmt.Printf("RoutePoints:%v\nNewRoute:   %v\n", n.RoutePoints, newRoute)
		n.RoutePoints = newRoute
		n.RoutePath = make([]Coord, 0)
		n.NumDestinations = 0
	}
	n.UpdatePath()
}
