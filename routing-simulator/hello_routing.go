package main

import (
	"bytes"
	"flag"
	"fmt"
	"image"
	"image/png"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
	"./cps"
)

var (


	p  *cps.Params
	r  *cps.RegionParams


	err 		 error
)

func init() {
	image.RegisterFormat("png", "png", png.Decode, png.DecodeConfig)
}



func main() {
	
	p = &cps.Params{}
	r = &cps.RegionParams{}
	
	
	prnt := false
	
	rand.Seed(time.Now().UTC().UnixNano())

	getFlags()

	p.MaxX = 408
	p.MaxY = 408
	
	//squareRow = squareRowCM
	//squareCol = squareColCM

	p.XDiv = p.MaxX / p.SquareColCM
	p.YDiv = p.MaxY / p.SquareRowCM

	cps.CreateBoard(p.MaxX, p.MaxY, p)

	roadFile, err := os.Create("roadLog.txt")
	if err != nil {
		log.Fatal("Cannot create file", err)
	}
	defer roadFile.Close()

	p.PositionFile, err = os.Create("Log-simulatorOutput.txt")
	if err != nil {
		log.Fatal("Cannot create file", err)
	}
	defer p.PositionFile.Close()

	fmt.Fprintln(p.PositionFile, "Width:", p.MaxX)
	fmt.Fprintln(p.PositionFile, "Height:", p.MaxY)
	fmt.Fprintf(p.PositionFile, "Amount: %v\n", 0)
	fmt.Fprintf(p.PositionFile, "Bomb X: %v\n", 0)
	fmt.Fprintf(p.PositionFile, "Bomb Y: %v\n", 0)

	r.Point_list = make([]cps.Tuple, 0)

	r.Point_list2 = make([][]bool, 0)

	r.Point_dict = make(map[cps.Tuple]bool)

	r.Square_list = make([]cps.RoutingSquare, 0)

	r.Border_dict = make(map[int][]int)

	//stimName := stimFileNameCM
	absPath, _ := filepath.Abs(p.StimFileNameCM)
	stimData, err := ioutil.ReadFile(absPath)
	if err != nil {
		log.Fatal("Cannot create file", err)
	}

	stim_line := strings.Split(string(stimData), "\n")
	r.Stim_list = make(map[int]cps.Tuple)
	fmt.Println("STIM_FILE LOOP")
	fmt.Println(len(stim_line))
	for i := 0; i < len(stim_line)-1; i++ {
		line := strings.Split(stim_line[i], ", ")
		x, _ := strconv.Atoi(line[0])
		y, _ := strconv.Atoi(line[1])
		t, _ := strconv.Atoi(line[2])
		fmt.Printf("%d %d %d %s\n", x, y, t, line)
		r.Stim_list[t] = cps.Tuple{x, y}
		//x, _ := strconv.Atoi(line[1])
		//y, _ := strconv.Atoi(line[3][:len(line[3])-1])

		//boardMap[y][x] = -1
	}

	imgfile, err := os.Open(p.ImageFileNameCM)
	if err != nil {
		fmt.Println("image file not found!")
		fmt.Println(p.ImageFileNameCM)
		os.Exit(1)
	}

	defer imgfile.Close()

	imgCfg, _, err := image.DecodeConfig(imgfile)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	width := imgCfg.Width
	height := imgCfg.Height

	fmt.Println("Width : ", width)
	fmt.Println("Height : ", height)

	imgfile.Seek(0, 0)

	img, _, err := image.Decode(imgfile)

	for x := 0; x < height; x++ {
		r.Point_list2 = append(r.Point_list2, make([]bool, width))
	}

	for x := 0; x < height; x++ {
		for y := 0; y < width; y++ {
			rr, _, _, _ := img.At(x, y).RGBA()
			if rr != 0 {
				r.Point_list2[x][y] = true
				r.Point_dict[cps.Tuple{x, y}] = true
				if prnt {
					fmt.Printf("X: %d, Y: %d\n", x, y)
				}

			} else {
				r.Point_dict[cps.Tuple{x, y}] = false
				p.BoardMap[y][x] = -1
			}
		}
	}

	top_left_corner := cps.Coord{X: 0, Y: 0}
	top_right_corner := cps.Coord{X: 0, Y: 0}
	bot_left_corner := cps.Coord{X: 0, Y: 0}
	bot_right_corner := cps.Coord{X: 0, Y: 0}

	tl_min := height + width
	tr_max := -1
	bl_max := -1
	br_max := -1

	for x := 0; x < height; x++ {
		for y := 0; y < width; y++ {
			if r.Point_dict[cps.Tuple{x, y}] {
				if x+y < tl_min {
					tl_min = x + y
					top_left_corner.X = x
					top_left_corner.Y = y
				}
				if y-x > tr_max {
					tr_max = y - x
					top_right_corner.X = x
					top_right_corner.Y = y
				}
				if x-y > bl_max {
					bl_max = x - y
					bot_left_corner.X = x
					bot_left_corner.Y = y
				}
				if x+y > br_max {
					br_max = x + y
					bot_right_corner.X = x
					bot_right_corner.Y = y
				}
			}
		}
	}

	fmt.Printf("TL: %v, TR %v, BL %v, BR %v\n", top_left_corner, top_right_corner, bot_left_corner, bot_right_corner)

	starting_locs := make([]cps.Coord, 4)
	starting_locs[0] = top_left_corner
	starting_locs[1] = top_right_corner
	starting_locs[2] = bot_left_corner
	starting_locs[3] = bot_right_corner

	statsFile, err := os.Create(p.OutRoutingStatsNameCM)
	if err != nil {
		log.Fatal("Cannot create file", err)
	}
	defer statsFile.Close()

	//New routing initialization
	if p.RegionRouting {

		fmt.Println("Beginning Region Routing")

		/*
			for x := 0; x < 200; x ++ {
				for y := 0; y < 200; y++ {
					if r.Point_dict[cps.Tuple{x, y}] {
						fmt.Print("_")
					} else {
						fmt.Print("X")
					}
				}
				fmt.Println()
			}
			fmt.Println(img.At(203, 26).RGBA())*/

		id_counter := 0
		done := false
		//for len(r.Point_list) != 0 {
		for !done {
			top_left := cps.Tuple{-1, -1}
			fmt.Println("starting")
			for x := 0; x < height; x++ {
				for y := 0; y < width; y++ {
					//fmt.Printf("X: %d, Y: %d, v: %d", x, y, r.Point_list2[x][y])
					if r.Point_list2[x][y] {
						top_left = cps.Tuple{x, y}
						break
					}
				}
				if (top_left != cps.Tuple{-1, -1}) {
					break
				}
			}
			fmt.Printf("working %d %d\n", top_left.X, top_left.Y)
			if (top_left == cps.Tuple{-1, -1}) {
				done = true
				break
			}
			//top_left := r.Point_list[0]
			temp := cps.Tuple{top_left.X, top_left.Y}

			for r.Point_dict[cps.Tuple{temp.X + 1, temp.Y}] {
				temp.X += 1
			}

			collide := false
			y_test := cps.Tuple{top_left.X, top_left.Y}

			for !collide {
				y_test.Y += 1
				if prnt {
					fmt.Println(y_test.Y)
				}

				for x_val := top_left.X; x_val < temp.X; x_val++ {
					if !r.Point_dict[cps.Tuple{x_val, y_test.Y}] {
						collide = true
					}
				}
			}

			bottom_right := cps.Tuple{temp.X, y_test.Y - 1}

			fmt.Println(top_left.X, bottom_right.X, top_left.Y, bottom_right.Y)

			new_square := cps.RoutingSquare{top_left.X, bottom_right.X, top_left.Y, bottom_right.Y, true, id_counter, make([]cps.Tuple, 0)}
			id_counter++
			fmt.Println("start_r_square")
			cps.RemoveRoutingSquare(new_square, r)
			fmt.Println("end_r_square")
			r.Square_list = append(r.Square_list, new_square)
		}

		length := len(r.Square_list)
		for y, _ := range r.Square_list {
			square := r.Square_list[y]
			r.Square_list[y].Routers = make([]cps.Tuple, length)

			for z := y + 1; z < len(r.Square_list); z++ {
				new_square := r.Square_list[z]

				if new_square.X1 >= square.X1 && new_square.X2 <= square.X2 {
					if new_square.Y1 == square.Y2+1 {
						r.Border_dict[y] = append(r.Border_dict[y], z)
						r.Border_dict[z] = append(r.Border_dict[z], y)

					} else if new_square.Y2 == square.Y1-1 {
						r.Border_dict[y] = append(r.Border_dict[y], z)
						r.Border_dict[z] = append(r.Border_dict[z], y)
					}
				} else if new_square.Y1 >= square.Y1 && new_square.Y2 <= square.Y2 {
					if new_square.X1 == square.X2+1 {
						r.Border_dict[y] = append(r.Border_dict[y], z)
						r.Border_dict[z] = append(r.Border_dict[z], y)

					} else if new_square.X2 == square.X1-1 {
						r.Border_dict[y] = append(r.Border_dict[y], z)
						r.Border_dict[z] = append(r.Border_dict[z], y)
					}
				}
				if square.X1 >= new_square.X1 && square.X2 <= new_square.X2 {
					if square.Y1 == new_square.Y2+1 {
						r.Border_dict[y] = append(r.Border_dict[y], z)
						r.Border_dict[z] = append(r.Border_dict[z], y)

					} else if square.Y2 == new_square.Y1-1 {
						r.Border_dict[y] = append(r.Border_dict[y], z)
						r.Border_dict[z] = append(r.Border_dict[z], y)
					}
				} else if square.Y1 >= new_square.Y1 && square.Y2 <= new_square.Y2 {
					if square.X1 == new_square.X2+1 {
						r.Border_dict[y] = append(r.Border_dict[y], z)
						r.Border_dict[z] = append(r.Border_dict[z], y)

					} else if square.X2 == new_square.X1-1 {
						r.Border_dict[y] = append(r.Border_dict[y], z)
						r.Border_dict[z] = append(r.Border_dict[z], y)
					}
				}
			}
		}
		fmt.Println(r.Border_dict)

		//Cutting takes place in this loop
		for true {
			rebuilt := false

			for i := 0; i < len(r.Square_list) && !rebuilt; i++ {

				for _, n := range r.Border_dict[i] {

					s_rat := cps.Side_ratio(r.Square_list[i], r.Square_list[n])
					if s_rat > 0.6 {
						new_squares := cps.Single_cut(r.Square_list[i], r.Square_list[n])

						s1 := r.Square_list[n]
						s2 := r.Square_list[i]

						cps.Square_list_remove(s1, r)
						cps.Square_list_remove(s2, r)

						r.Square_list = append(r.Square_list, new_squares...)

						cps.Rebuild(r.Square_list, r)

						rebuilt = true

						break
					}

					a_rat := cps.Area_ratio(r.Square_list[i], r.Square_list[n])
					if a_rat > 0.6 {
						new_squares := cps.Double_cut(r.Square_list[i], r.Square_list[n])

						s1 := r.Square_list[n]
						s2 := r.Square_list[i]

						cps.Square_list_remove(s1, r)
						cps.Square_list_remove(s2, r)

						new_squares[2].Id_num = len(r.Square_list)

						r.Square_list = append(r.Square_list, new_squares...)

						cps.Rebuild(r.Square_list, r)

						rebuilt = true

						break
					}
				}
			}

			if !rebuilt {
				break
			}
		}

		

		r.Node_tables = make([]map[cps.Tuple]float64, len(r.Square_list))

		for key, values := range r.Border_dict {
			if key < len(r.Square_list) {
				r.Node_tables[key] = make(map[cps.Tuple]float64)
				if len(values) > 1 {
					for n := 0; n < len(values); n++ {
						next := n + 1
						for next < len(values) {
							node_a := r.Border_dict[key][n]
							node_b := r.Border_dict[key][next]

							p1 := r.Square_list[key].Routers[node_a]
							p2 := r.Square_list[key].Routers[node_b]

							r.Node_tables[key][cps.Tuple{node_a, node_b}] = cps.Dist(p1, p2)
							r.Node_tables[key][cps.Tuple{node_b, node_a}] = cps.Dist(p1, p2)

							next += 1
						}
					}
				}
			}
		}

		routingFile, err := os.Create(p.OutRoutingNameCM)
		if err != nil {
			log.Fatal("Cannot create file", err)
		}
		defer routingFile.Close()

		//The scheduler determines which supernode should pursue a point of interest
		scheduler := &cps.Scheduler{}

		//List of all the supernodes on the grid
		scheduler.SNodeList = make([]cps.SuperNodeParent, p.NumSuperNodes)

		for i := 0; i < p.NumSuperNodes; i++ {
			snode_points := make([]cps.Coord, 1)
			snode_path := make([]cps.Coord, 0)
			all_points := make([]cps.Coord, 0)

			//Defining the starting x and y values for the super node
			//This super node starts at the middle of the grid
			x_val, y_val := starting_locs[i].X, starting_locs[i].Y
			nodeCenter := cps.Coord{X: x_val, Y: y_val}

			scheduler.SNodeList[i] = &cps.Sn_zero{&cps.Supern{&cps.NodeImpl{X: x_val, Y: y_val, Id: i}, 1,
				1, p.SuperNodeRadius, p.SuperNodeRadius, 0, snode_points, snode_path,
				nodeCenter, 0, 0, 0, 0, 0, all_points}}

			//The super node's current location is always the first element in the routePoints list
			scheduler.SNodeList[i].UpdateLoc()
		}

		fmt.Printf("Iteration %d/%v", 0, p.Iterations_of_event)
		for i := 0; i < p.Iterations_of_event; i++ {
			fmt.Printf("\rIteration %d/%v", i, p.Iterations_of_event)

			if t, ok := r.Stim_list[i]; ok {
				start := time.Now()
				scheduler.AddRoutePoint(cps.Coord{X: t.X, Y: t.Y}, p, r)
				elapsed := time.Since(start)

				fmt.Fprint(statsFile, "Region Routing Elapsed: ", elapsed, "\n")

				fmt.Printf("\nAdding %d %d %d\n", i, t.X, t.Y)
			}

			for _, s := range scheduler.SNodeList {
				points_len := len(s.GetRoutePoints())
				response_time := -1
				if points_len > 1 {
					response_time = s.GetRoutePoints()[1].Time
				}

				s.Tick(p, r)

				if points_len > len(s.GetRoutePoints()) {
					fmt.Fprint(statsFile, "Response Time: ", response_time, "\n")
				}

				fmt.Fprint(routingFile, s)
				p := printPoints(s)
				fmt.Fprint(routingFile, " UnvisitedPoints: ")
				fmt.Fprintln(routingFile, p.String())
			}
		}
		for _, s := range scheduler.SNodeList {
			fmt.Println("\nSQUARES MOVED", s.GetSquaresMoved())
			fmt.Println("\nPOINTS VISITED", s.GetPointsVisited())

			fmt.Fprintf(statsFile, "Squares Moved: %d %d\n", s.GetId(), s.GetSquaresMoved())
			fmt.Fprintf(statsFile, "Points Visited: %d %d\n", s.GetId(), s.GetPointsVisited())
		}

	}
	//End new routing initialization

	//aStar routing
	if p.AStarRouting {

		fmt.Println("Beginning aStar Routing")

		routingFile, err := os.Create(p.OutRoutingNameCM)
		if err != nil {
			log.Fatal("Cannot create file", err)
		}
		defer routingFile.Close()

		//The scheduler determines which supernode should pursue a point of interest
		scheduler := &cps.Scheduler{}

		//List of all the supernodes on the grid
		scheduler.SNodeList = make([]cps.SuperNodeParent, p.NumSuperNodes)

		for i := 0; i < p.NumSuperNodes; i++ {
			snode_points := make([]cps.Coord, 1)
			snode_path := make([]cps.Coord, 0)
			all_points := make([]cps.Coord, 0)

			//Defining the starting x and y values for the super node
			//This super node starts at the middle of the grid
			x_val, y_val := starting_locs[i].X, starting_locs[i].Y
			nodeCenter := cps.Coord{X: x_val, Y: y_val}

			scheduler.SNodeList[i] = &cps.Sn_zero{&cps.Supern{&cps.NodeImpl{X: x_val, Y: y_val, Id: i}, 1,
				1, p.SuperNodeRadius, p.SuperNodeRadius, 0, snode_points, snode_path,
				nodeCenter, 0, 0, 0, 0, 0, all_points}}

			//The super node's current location is always the first element in the routePoints list
			scheduler.SNodeList[i].UpdateLoc()
		}

		fmt.Print("POINTS", scheduler.SNodeList[0].GetPointsVisited())

		fmt.Printf("Iteration %d/%v", 0, p.Iterations_of_event)
		for i := 0; i < p.Iterations_of_event; i++ {
			fmt.Printf("\rIteration %d/%v", i, p.Iterations_of_event)

			if t, ok := r.Stim_list[i]; ok {
				//scheduler.addRoutePoint(cps.Coord{X: t.X, Y: t.Y})
				start := time.Now()
				scheduler.AddRoutePoint(cps.Coord{X: t.X, Y: t.Y}, p, r)
				elapsed := time.Since(start)

				fmt.Fprint(statsFile, "aStar Routing Elapsed ", elapsed, "\n")

				fmt.Printf("\nAdding %d %d %d\n", i, t.X, t.Y)
			}

			for _, s := range scheduler.SNodeList {
				points_len := len(s.GetRoutePoints())
				response_time := -1
				if points_len > 1 {
					response_time = s.GetRoutePoints()[1].Time
				}

				s.Tick(p, r)

				if points_len > len(s.GetRoutePoints()) {
					fmt.Fprint(statsFile, "Response Time: ", response_time, "\n")
				}

				fmt.Fprint(routingFile, s)
				p := printPoints(s)
				fmt.Fprint(routingFile, " UnvisitedPoints: ")
				fmt.Fprintln(routingFile, p.String())
			}
		}
		for _, s := range scheduler.SNodeList {
			fmt.Println("\nSQUARES MOVED", s.GetSquaresMoved())
			fmt.Println("\nPOINTS VISITED", s.GetPointsVisited())

			fmt.Fprintf(statsFile, "Distance Traveled: %d %d\n", s.GetId(), s.GetSquaresMoved())
			fmt.Fprintf(statsFile, "Points Visited: %d %d\n", s.GetId(), s.GetPointsVisited())
		}
	}
}

//This function allows the simulator to create a roadMap of the grid
//Every cps.Coord in the grid is given an integer value corresponding to the
//	number of times the cps.Coord is used by all paths
//The function first generates two random cps.Coords on each half of the grid
//It then finds the path between those cps.Coords
//It then increments the integer value of each cps.Coord in the path by one
//This is done an amount of time to generate a conclusive distribution of paths
//	across the gird
//THe resulting roadMap is outputted to the file, first with the max number
//	if times a cps.Coord is visited and then each cps.Coord's integer value
func makeRoads(roadFile *os.File) {
	//This map has cps.Tuples as keys and integers as values
	//The cps.Tuples represent the cps.Coord in the grid and the integer represents
	//	the amount of times the cps.Coord is visited by all paths
	roadMap := make(map[cps.Tuple]int)

	//The max is kept track of the be outputted at the beginning of the
	//road output file
	//This is used to determine the gradient of color by the Viewer when
	//	displaying the roads
	max := 0

	aStarIterations := 100

	fmt.Printf("Running ASTAR p.Iteration %d/%v", 0, aStarIterations)
	for i := 0; i < aStarIterations; i++ {
		//Two cps.Coords are randomly generated
		a := cps.Coord{nil, rangeInt(0, p.MaxX), rangeInt(0, p.MaxY), 0, 0, 0, 0}
		b := cps.Coord{nil, rangeInt(0, p.MaxX), rangeInt(0, p.MaxY), 0, 0, 0, 0}

		//The cps.Coords' x and y positions are randomly updated to be on either the
		//	left and right side, or top and bottom of the grid
		//This is done to ensure the paths between these cps.Coords crosses a large
		//	section of the grid
		if i%2 == 0 {
			a.X = rangeInt(0, p.MaxX/2)
			b.X = rangeInt(p.MaxX/2, p.MaxX)
		} else {
			a.Y = rangeInt(0, p.MaxY/2)
			b.Y = rangeInt(p.MaxY/2, p.MaxY)
		}
		fmt.Printf("\rRunning ASTAR p.Iteration %d/%v", i, aStarIterations)
		//The aStar path between these two cps.Coords is created
		//Each cps.Coord in this path is looped through and the integer value corresponding
		//	to that cps.Coord is incremented by one
		for _, rr := range cps.AStar(a, b, p) {
			pos := cps.Tuple{rr.X, rr.Y}
			roadMap[pos]++
			if roadMap[pos] > max {
				max = roadMap[pos]
			}
		}
	}
	fmt.Fprintln(roadFile, "max", max)

	//This loops through the roadMap and outputs the integer value for every cps.Coord
	//	in the grid
	for i := 0; i < p.MaxX; i++ {
		for j := 0; j < p.MaxY; j++ {
			fmt.Println("\rOutputting to roadLog: cps.Coord", j, i)
			if p.BoardMap[i][j] == -1 {
				fmt.Fprintln(roadFile, i, j, -1)
			} else {
				fmt.Fprintln(roadFile, i, j, roadMap[cps.Tuple{j, i}])
			}
		}
	}
}

//Saves the cps.Coords in the allPoints list into a buffer to
//	print to the file
func printPoints(s cps.SuperNodeParent) bytes.Buffer {
	var buffer bytes.Buffer
	buffer.WriteString((fmt.Sprintf("[")))
	for ind, i := range s.GetAllPoints() {
		buffer.WriteString(i.String())

		if ind != len(s.GetAllPoints())-1 {
			buffer.WriteString(" ")
		}
	}
	buffer.WriteString((fmt.Sprintf("]")))
	return buffer
}

func getFlags() {
	flag.IntVar(&p.Iterations_of_event, "p.Iterations_of_event", 1000, "how many times the simulation will run")

	//fmt.Println(os.Args[1:], "\nhmmm? \n ") //C:\Users\Nick\Desktop\comand line experiments\src
	flag.IntVar(&p.NegativeSittingStopThresholdCM, "negativeSittingStopThreshold", -10,
		"Negative number sitting is set to when board map is reset")
	flag.IntVar(&p.SittingStopThresholdCM, "sittingStopThreshold", 5,
		"How long it takes for a node to stay seated")
	flag.Float64Var(&p.GridCapacityPercentageCM, "gridCapacityPercentage", .9,
		"Percent the sub-grid can be filled")
	flag.StringVar(&p.InputFileNameCM, "inputFileName", "Log1_in.txt",
		"Name of the input text file")
	flag.StringVar(&p.OutputFileNameCM, "outputFileName", "Log",
		"Name of the output text file prefix")
	flag.Float64Var(&p.NaturalLossCM, "naturalLoss", .005,
		"battery loss due to natural causes")
	flag.Float64Var(&p.SensorSamplingLossCM, "sensorSamplingLoss", .001,
		"battery loss due to sensor sampling")
	flag.Float64Var(&p.GPSSamplingLossCM, "GPSSamplingLoss", .005,
		"battery loss due to GPS sampling")
	flag.Float64Var(&p.ServerSamplingLossCM, "serverSamplingLoss", .01,
		"battery loss due to server sampling")
	flag.IntVar(&p.ThresholdBatteryToHaveCM, "thresholdBatteryToHave", 30,
		"Threshold battery phones should have")
	flag.IntVar(&p.ThresholdBatteryToUseCM, "thresholdBatteryToUse", 10,
		"Threshold of battery phones should consume from all forms of sampling")
	flag.IntVar(&p.MovementSamplingSpeedCM, "movementSamplingSpeed", 20,
		"the threshold of speed to increase sampling rate")
	flag.IntVar(&p.MovementSamplingPeriodCM, "movementSamplingPeriod", 1,
		"the threshold of speed to increase sampling rate")
	flag.IntVar(&p.MaxBufferCapacityCM, "maxBufferCapacity", 25,
		"maximum capacity for the buffer before it sends data to the server")
	flag.StringVar(&p.EnergyModelCM, "energyModel", "variable",
		"this determines the energy loss model that will be used")
	flag.BoolVar(&p.NoEnergyModelCM, "noEnergy", false,
		"Whether or not to ignore energy model for simulation")
	flag.IntVar(&p.SensorSamplingPeriodCM, "sensorSamplingPeriod", 1000,
		"rate of the sensor sampling period when custom energy model is chosen")
	flag.IntVar(&p.GPSSamplingPeriodCM, "GPSSamplingPeriod", 1000,
		"rate of the GridGPS sampling period when custom energy model is chosen")
	flag.IntVar(&p.ServerSamplingPeriodCM, "serverSamplingPeriod", 1000,
		"rate of the server sampling period when custom energy model is chosen")
	flag.IntVar(&p.NumStoredSamplesCM, "nodeStoredSamples", 10,
		"number of samples stored by individual nodes for averaging")
	flag.IntVar(&p.GridStoredSamplesCM, "gridStoredSamples", 10,
		"number of samples stored by grid squares for averaging")
	flag.Float64Var(&p.DetectionThresholdCM, "detectionThreshold", 11180.0,
		"Value where if a node gets this reading or higher, it will trigger a detection")
	flag.Float64Var(&p.ErrorModifierCM, "errorMultiplier", 1.0,
		"Multiplier for error values in system")
	//Range 1, 2, or 4
	//Currently works for only a few numbers, can be easily expanded but is not currently dynamic
	flag.IntVar(&p.NumSuperNodes, "numSuperNodes", 1, "the number of super nodes in the simulator")

	//Range: 0-2
	//0: default routing algorithm, points added onto the end of the path and routed to in that order
	//flag.IntVar(&superNodeType, "superNodeType", 0, "the type of super node used in the simulator")
	//better descriptions incoming
	//Range: 0-6
	//	0: default routing algorithm, points added onto the end of the path and routed to in that order
	//	1: sophisticated routing algorithm, begin in center, routed anywhere
	//	2: sophisticated routing algorithm, begin inside circles located in the corners, only routed inside circle
	//	3: sophisticated routing algorithm, begin inside circles located on the sides, only routed inside circle
	//	4: sophisticated routing algorithm, being inside large circles located in the corners, only routed inside circle
	//	5: sophisticated routing algorithm, begin inside regions, only routed inside region
	//	6: regional return trip routing algorithm, routed inside region based on most points
	//	7: regional return trip routing algorithm, routed inside region based on oldest point
	flag.IntVar(&p.SuperNodeType, "superNodeType", 0, "the type of super node used in the simulator")

	//Range: 0-...
	//Theoretically could be as high as possible
	//Realistically should remain around 10
	flag.IntVar(&p.SuperNodeSpeed, "superNodeSpeed", 1, "the speed of the super node")

	//Range: true/false
	//Tells the simulator whether or not to optimize the path of all the super nodes
	//Only works when multiple super nodes are active in the same area
	flag.BoolVar(&p.DoOptimize, "doOptimize", false, "whether or not to optimize the simulator")

	//Range: 0-4
	//	0: begin in center, routed anywhere
	//	1: begin inside circles located in the corners, only routed inside circle
	//	2: begin inside circles located on the sides, only routed inside circle
	//	3: being inside large circles located in the corners, only routed inside circle
	//	4: begin inside regions, only routed inside region
	//Only used for super nodes of type 1
	//flag.IntVar(&superNodeVariation, "superNodeVariation", 3, "super nodes of type 1 have different variations")

	flag.BoolVar(&p.PositionPrintCM, "logPosition", false, "Whether you want to write position info to a log file")
	flag.BoolVar(&p.GridPrintCM, "logGrid", false, "Whether you want to write grid info to a log file")
	flag.BoolVar(&p.EnergyPrintCM, "logEnergy", false, "Whether you want to write energy into to a log file")
	flag.BoolVar(&p.NodesPrintCM, "logNodes", false, "Whether you want to write node readings to a log file")
	flag.IntVar(&p.SquareRowCM, "squareRow", 100, "Number of rows of grid squares, 1 through p.MaxX")
	flag.IntVar(&p.SquareColCM, "squareCol", 100, "Number of columns of grid squares, 1 through p.MaxY")

	flag.BoolVar(&p.RegionRouting, "regionRouting", false, "True if you want to use the new routing algorithm with regions and cutting")
	flag.BoolVar(&p.AStarRouting, "aStarRouting", false, "True if you want to use the old routing algorithm with aStar")

	flag.StringVar(&p.ImageFileNameCM, "imageFileName", "circle_justWalls_x4.png", "Name of the input text file")
	flag.StringVar(&p.StimFileNameCM, "stimFileName", "circle_0.txt", "Name of the stimulus text file")
	flag.StringVar(&p.OutRoutingNameCM, "outRoutingName", "log.txt", "Name of the stimulus text file")
	flag.StringVar(&p.OutRoutingStatsNameCM, "outRoutingStatsName", "routingStats.txt", "Name of the output file for stats")

	flag.Parse()
}

func rangeInt(min, max int) int { //returns a random number between max and min
	return rand.Intn(max-min) + min
}
