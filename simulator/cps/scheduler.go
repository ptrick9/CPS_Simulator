package cps

import (
	"math"
)

//Then scheduler type contains the list of all the super nodes in the simulation
//It is responsible for adding the points of interest to the super nodes, this
//	process differs on the type of super node
type Scheduler struct {
	P *Params
	R *RegionParams
	SNodeList []SuperNodeParent
	//SNodeList []Supern
}

//This function is called every Time a super node reaches a point of interest along
//	its routePath
//The scheduler removes all the nodes' points of interest and redistributes them again
//	attempting to make the routing process more efficient
func (s *Scheduler) Optimize() {
	//pp := s.P
	//r := s.R
	points := make([]Coord, 0)

	//Loops through the super node list appending every points to the newly created points list
	//Then removes the range of points from each super node's routePoints list
	for sn, _ := range s.SNodeList {
		for i := 1; i < len(s.SNodeList[sn].GetRoutePoints()); i++ {
			points = append(points, s.SNodeList[sn].GetRoutePoints()[i])
		}

		//Removes every point from the routePoints list except the first point which represents
		//	the location of the super node
		s.SNodeList[sn].SetRoutePoints(Remove_range(s.SNodeList[sn].GetRoutePoints(),
			1, len(s.SNodeList[sn].GetRoutePoints())-1))

		//Resets the routePath of the super node
		path := make([]Coord, 0)
		s.SNodeList[sn].SetRoutePath(path)

		//Resets the numDestinations to 0
		s.SNodeList[sn].SetNumDest(0)
	}

	//Adds the points back into the simulator
	for p, _ := range points {
		s.AddRoutePoint(points[p])
	}
}

//This function determines which super node adding method should be called
func (s *Scheduler) AddRoutePoint(c Coord) {
	if s.P.SuperNodeType == 0 {
		s.AddRoutePoint0(c)
	} else if s.P.SuperNodeType == 1 {
		s.AddRoutePoint1(c)
	} else if s.P.SuperNodeType == 2 || s.P.SuperNodeType == 3 || s.P.SuperNodeType == 4 {
		s.AddRoutePoint1_circle(c)
	} else if s.P.SuperNodeType == 5 {
		s.AddRoutePoint1_regions(c)
	} else if s.P.SuperNodeType == 6 || s.P.SuperNodeType == 7 {
		s.AddRoutePoint2(c)
	}
}

//Adds a point of interest to a super node of type 0
//Since super node 0 operates on the default scheduling algorithm the
//	scheduler adds the new point of interest to the super node who's
//	final destination is closest to the point
func (s *Scheduler) AddRoutePoint0(c Coord) {
	dist := 100000.0
	nodeDist := 100000.0
	closestNode := -1

	//Finds the super node closest to the newly added point
	for n, _ := range s.SNodeList {
		length := len(s.SNodeList[n].GetRoutePath())
		reg := RegionContaining(Tuple{s.SNodeList[n].GetX(), s.SNodeList[n].GetY()}, s.R)
		//fmt.Printf("\nSuper Node X:%v, Y:%v\n", s.SNodeList[n].GetX(), s.SNodeList[n].GetY())
		//fmt.Printf("Starting region: %v\n", reg)
		if ValidPath(reg, c, s.R) {
			//fmt.Println("Valid path found!")
			if length != 0 {
				nodeDist = math.Sqrt(math.
					Pow(float64(s.SNodeList[n].GetRoutePath()[length-1].X-c.X), 2.0) + math.
					Pow(float64(s.SNodeList[n].GetRoutePath()[length-1].Y-c.Y), 2.0))

				nodeDist += float64(length)
			} else {
				nodeDist = math.Sqrt(math.
					Pow(float64(s.SNodeList[n].GetX()-c.X), 2.0) + math.
					Pow(float64(s.SNodeList[n].GetY()-c.Y), 2.0))
			}
			if nodeDist < dist {
				dist = nodeDist
				closestNode = n
			}
		}
	}
	//fmt.Printf("\nClosest Node: %v\n", closestNode)
	//Tells that super node to add that point
	//fmt.Printf("\nAdded route point %v\n", c)
	s.SNodeList[closestNode].AddRoutePoint(c)
}

//Adds a point of interest to a super node of type 1
//This is a more complicated and sophisticated version of the super node 0
//	adding function
//This adds the distance from the newly added point to the closest spot on the
//	super node's path to the length of the super node's current path
//This prioritizes the proximity of the super nodes in determining which
//	one will travel to the newly added point, but also adds in the distance that
//	super node is currently attempting to travel
func (s *Scheduler) AddRoutePoint1(c Coord) {
	dist := 10000.0
	closestNode := -1
	nodeDist := 0.0

	//Loops through the list of super nodes to find the optimal super node to
	//	travel to the newly added point
	//The newly added point will be visited by the super node who can reach that
	//	point in the least distance
	for n, _ := range s.SNodeList {

		//If the super node's routePath is not zero, meaning it is currently
		//	travelling to point(s) of interest, it's closest distance is found
		//This finds the closest distance between a Coord on the super node's
		//	routePath and the newly added point
		//Otherwise the distance from the super node's current location is considered
		if len(s.SNodeList[n].GetRoutePath()) > 0 {
			nodeDist = ClosestDist(c, s.SNodeList[n].GetRoutePath())

			//The length of the current path is added to ensure no super nodes
			//	are forced to go to points of interest if they have an extremely
			//	long routePath
			nodeDist += float64(len(s.SNodeList[n].GetRoutePath()))
		} else {
			nodeDist = math.Sqrt(math.
				Pow(float64(s.SNodeList[n].GetRoutePoints()[0].X-c.X), 2.0) + math.
				Pow(float64(s.SNodeList[n].GetRoutePoints()[0].Y-c.Y), 2.0))
		}

		//Saves the smallest dist Value to return
		if nodeDist < dist {
			dist = nodeDist
			closestNode = n
		}
	}
	//Tells that super node to add that point to the decided super node
	s.SNodeList[closestNode].AddRoutePoint(c)
}

//This is a variation on the super node 1 adding function
//This restricts the super node to a circular region that covers an area of the
//	entire grid
func (s *Scheduler) AddRoutePoint1_circle(c Coord) {
	circleNode := -1

	//Loops through the list of super nodes to find the optimal super node to visit
	//	the newly added point
	for n, _ := range s.SNodeList {
		//Calculates the distance from the newly added point to center of each super
		//	node's circle
		nodeDist := int(math.Sqrt(math.
			Pow(float64(s.SNodeList[n].GetCenter().X-c.X), 2.0) + math.
			Pow(float64(s.SNodeList[n].GetCenter().Y-c.Y), 2.0)))

		//If the point of interest is inside the super node's circle it is added, unless
		//	the point has been claimed by another super node
		//If another super node is currently chosen to visit the newly added point the
		//	length of the super nodes' routePaths are compared
		if nodeDist <= s.P.SuperNodeRadius {
			if circleNode != -1 {
				if ClosestDist(c, s.SNodeList[n].GetRoutePath()) < ClosestDist(c, s.SNodeList[circleNode].GetRoutePath()) {
					circleNode = n
				}
			} else {
				circleNode = n
			}
		}
	}
	//Tells that super node to add that point to the decided super node
	s.SNodeList[circleNode].AddRoutePoint(c)
}

//This is a variation on the super node 1 adding function
//This restricts the super node to a quadrant of the grid that only it covers
func (s *Scheduler) AddRoutePoint1_regions(c Coord) {
	//Boundary conditions
	if c.X < (s.P.MaxX / 2) {
		if c.Y < (s.P.MaxY / 2) {
			s.SNodeList[0].AddRoutePoint(c)
		} else {
			s.SNodeList[2].AddRoutePoint(c)
		}
	} else {
		if c.Y < (s.P.MaxY / 2) {
			s.SNodeList[1].AddRoutePoint(c)
		} else {
			s.SNodeList[3].AddRoutePoint(c)
		}
	}
}

//Adds a point of interest to a super node of type 2
//Super nodes of type 2 schedule their routes within regions so this
//	function add the point to the super node whose center is closest
//	to the point
func (s *Scheduler) AddRoutePoint2(c Coord) {
	dist := 1000.0
	closestNode := -1

	//Finds the super node whose center is closest to the newly added point
	for n, _ := range s.SNodeList {
		nodeDist := math.Sqrt(math.Pow(float64(s.SNodeList[n].GetCenter().X-c.X), 2.0) + math.Pow(float64(s.SNodeList[n].GetCenter().Y-c.Y), 2.0))
		if nodeDist < dist {
			dist = nodeDist
			closestNode = n
		}
	}

	//Tells that super node to add that point
	s.SNodeList[closestNode].AddRoutePoint(c)
}

//This function returns the distance between the specified Coord c
//	and the closest Coord in the provided list of Coords
func ClosestDist(c Coord, list []Coord) float64 {
	dist := 1000.0

	//Loops through the list to find the closest Coord
	for _, p := range list {
		newDist := math.Sqrt(math.
			Pow(float64(p.X-c.X), 2.0) + math.
			Pow(float64(p.Y-c.Y), 2.0))

		//Saves the Value of that smallest distance to return
		if newDist < dist {
			dist = newDist
		}
	}
	return dist
}

//MakeSuperNodes generates super nodes and adds them to the super node list within scheduler
func (scheduler *Scheduler) MakeSuperNodes() {
	for i := 0; i < scheduler.P.NumSuperNodes; i++ {
		snode_points := make([]Coord, 1)
		snode_path := make([]Coord, 0)
		all_points := make([]Coord, 0)

		if scheduler.P.SuperNodeType == 0 {

			//Defining the starting x and y values for the super node
			//This super node starts at the middle of the p.Grid
			nodeCenter, x_val, y_val := MakeCenter1(i, scheduler.P)

			scheduler.SNodeList[i] = &Sn_zero{scheduler.P, scheduler.R,&Supern{scheduler.P, scheduler.R,&NodeImpl{X: float32(x_val), Y: float32(y_val), Id: i}, 1,
				1, scheduler.P.SuperNodeRadius, scheduler.P.SuperNodeRadius, 0, snode_points, snode_path,
				nodeCenter, 0, 0, 0, 0, 0, all_points}}
		} else if (scheduler.P.SuperNodeType == 6) || (scheduler.P.SuperNodeType == 7) {
			//makeRegionList initializes the regionList for this super node
			r_list := MakeRegionList(i, scheduler.P)

			//makeCenter creates the Coord that represents the super node's center
			nodeCenter, x_val, y_val := MakeCenter2(i, r_list, scheduler.P)

			//The useRegionList is just initialized to an empty list
			ur_list := make([]Region, 0)

			scheduler.SNodeList[i] = &Sn_two{&Supern{scheduler.P,scheduler.R,&NodeImpl{Id: i, X: float32(x_val), Y: float32(y_val)}, 1,
				1, scheduler.P.SuperNodeRadius, scheduler.P.SuperNodeRadius, 0, snode_points,
				snode_path, nodeCenter, 0, 0, 0, 0,
				1, all_points}, r_list, ur_list}
		} else if (scheduler.P.SuperNodeType >= 1) || (scheduler.P.SuperNodeType <= 5) {
			nodeCenter := Coord{}
			x_val := 0
			y_val := 0
			xRad := 0
			yRad := 0

			//makeCenter creates the Coord that represents the super node's center
			if scheduler.P.SuperNodeType == 1 {
				nodeCenter, x_val, y_val = MakeCenter1(i, scheduler.P)
			} else if scheduler.P.SuperNodeType == 2 || scheduler.P.SuperNodeType == 5 {
				nodeCenter, x_val, y_val, xRad, yRad = MakeCenter1_corners(i, scheduler.P)
			} else if scheduler.P.SuperNodeType == 3 {
				nodeCenter, x_val, y_val, xRad, yRad = MakeCenter1_sides(i, scheduler.P)
			} else if scheduler.P.SuperNodeType == 4 {
				nodeCenter, x_val, y_val, xRad, yRad = MakeCenter1_largeCorners(i, scheduler.P)
			}
			scheduler.SNodeList[i] = &Sn_one{&Supern{scheduler.P,scheduler.R, &NodeImpl{X: float32(x_val), Y: float32(y_val), Id: i}, 1,
				1, xRad, yRad, 0, snode_points, snode_path,
				nodeCenter, 0, 0,
				0, 0, 1, all_points}}
		}
		//The super node's current location is always the first element in the routePoints list
		scheduler.SNodeList[i].UpdateLoc()
	}
}
