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
)

var (
	squareRow        int
	squareCol        int
	numStoredSamples int
	maxX             int
	maxY             int

	threshHoldBatteryToHave  float32
	totalPercentBatteryToUse float32
	iterations_used          int
	iterations_of_event      int
	b                        *bomb

	//How the griid is divided into rows and columns
	xDiv int
	yDiv int

	recalibrate    bool
	squareCapacity int
	boolGrid       [][]bool
	grid           [][]*Square

	detectionThreshold float64 //Value of sensor reading which determines a "detection" of a bomb

	// These are the command line variables.
	negativeSittingStopThresholdCM int     // This is a negative number for the sitting to be set to when map is reset
	sittingStopThresholdCM         int     // This is the threshold for the longest time a node can sit before no longer moving
	gridCapacityPercentageCM       float64 // This is the percent of a subgrid that can be filled with nodes, between 0.0 and 1.0
	errorModifierCM                float64 // Multiplier for error model
	outputFileNameCM               string  // This is the prefix of the output text file
	inputFileNameCM                string  // This must be the name of the input text file with ".txt"
	naturalLossCM                  float64 // This can be any number n: 0 < n < .1
	sensorSamplingLossCM           float64 // This can be any number n: 0 < n < .1
	GPSSamplingLossCM              float64 // This can be any number n: 0 < n < GPSSamplingLossCM < .1
	serverSamplingLossCM           float64 // This can be any number n: 0 < n < serverSamplingLossCM < .1
	thresholdBatteryToHaveCM       int     // This can be any number n: 0 < n < 50
	thresholdBatteryToUseCM        int     // This can be any number n: 0 < n < 20 < 100-thresholdBatteryToHaveCM
	movementSamplingSpeedCM        int     // This can be any number n: 0 < n < 100
	movementSamplingPeriodCM       int     // This can be any int number n: 1 <= n <= 100
	maxBufferCapacityCM            int     // This can be aby int number n: 10 <= n <= 100
	energyModelCM                  string  // This can be "custom", "2StageServer", or other string will result in dynamic
	noEnergyModelCM                bool    // If set to true, all energy model values ignored
	sensorSamplingPeriodCM         int     // This can be any int n: 1 <= n <= 100
	GPSSamplingPeriodCM            int     // This can be any int n: 1 <= n < sensorSamplingPeriodCM <=  100
	serverSamplingPeriodCM         int     // This can be any int n: 1 <= n < GPSSamplingPeriodCM <= 100
	numStoredSamplesCM             int     // This can be any int n: 5 <= n <= 25
	gridStoredSamplesCM            int     // This can be any int n: 5 <= n <= 25
	detectionThresholdCM           float64 //This is whatever value 1-1000 we determine should constitute a "detection" of a bomb
	positionPrintCM                bool    //This is either true or false for whether to print positions to log file
	energyPrintCM                  bool    //This is either true or false for whether to print energy info to log file
	nodesPrintCM                   bool    //This is either true or false for whether to print node readings/averages to log file
	gridPrintCM                    bool    //This is either true or false for whether to print grid readings to log file
	squareRowCM                    int     //This is an int 1 through maxX representing how many rows of squares there are
	squareColCM                    int     //This is an int 1 through maxY representing how many columns of squares there are

	numSuperNodes  int
	superNodeType  int
	superNodeSpeed int
	doOptimize     bool
	//superNodeVariation int
	superNodeRadius int

	center Coord

	positionPrint bool
	nodesPrint    bool

	regionRouting bool
	aStarRouting  bool

	imageFileNameCM       string // This must be the name of the wall image file with ".png"
	stimFileNameCM        string // This must be the name of the stim file with ".txt"
	outRoutingNameCM      string // This is the name of the output routing file with ".txt"
	outRoutingStatsNameCM string // This is the name of the output stats file with ".txt"

	driftFile    *os.File
	nodeFile     *os.File
	positionFile *os.File
	statsFile    *os.File

	point_list []Tuple

	point_list2 [][]bool

	point_dict map[Tuple]bool

	square_list []RoutingSquare

	border_dict map[int][]int

	node_tables []map[Tuple]float64

	possible_paths [][]int

	stim_list map[int]Tuple

	prnt bool = false

	// End the command line variables.
)

func init() {
	image.RegisterFormat("png", "png", png.Decode, png.DecodeConfig)
}

const Tau1 = 10
const Tau2 = 500

type Tuple struct {
	x, y int
}

func main() {
	rand.Seed(time.Now().UTC().UnixNano())

	getFlags()

	maxX = 408
	maxY = 408
	squareRow = squareRowCM
	squareCol = squareColCM

	xDiv = maxX / squareCol
	yDiv = maxY / squareRow

	createBoard(maxX, maxY)

	roadFile, err := os.Create("roadLog.txt")
	if err != nil {
		log.Fatal("Cannot create file", err)
	}
	defer roadFile.Close()

	positionFile, err = os.Create("Log-simulatorOutput.txt")
	if err != nil {
		log.Fatal("Cannot create file", err)
	}
	defer positionFile.Close()

	fmt.Fprintln(positionFile, "Width:", maxX)
	fmt.Fprintln(positionFile, "Height:", maxY)
	fmt.Fprintf(positionFile, "Amount: %v\n", 0)
	fmt.Fprintf(positionFile, "Bomb x: %v\n", 0)
	fmt.Fprintf(positionFile, "Bomb y: %v\n", 0)

	point_list = make([]Tuple, 0)

	point_list2 = make([][]bool, 0)

	point_dict = make(map[Tuple]bool)

	square_list = make([]RoutingSquare, 0)

	border_dict = make(map[int][]int)

	stimName := stimFileNameCM
	absPath, _ := filepath.Abs(stimName)
	stimData, err := ioutil.ReadFile(absPath)
	if err != nil {
		log.Fatal("Cannot create file", err)
	}

	stim_line := strings.Split(string(stimData), "\n")
	stim_list = make(map[int]Tuple)
	fmt.Println("STIM_FILE LOOP")
	fmt.Println(len(stim_line))
	for i := 0; i < len(stim_line)-1; i++ {
		line := strings.Split(stim_line[i], ", ")
		x, _ := strconv.Atoi(line[0])
		y, _ := strconv.Atoi(line[1])
		t, _ := strconv.Atoi(line[2])
		fmt.Printf("%d %d %d %s\n", x, y, t, line)
		stim_list[t] = Tuple{x, y}
		//x, _ := strconv.Atoi(line[1])
		//y, _ := strconv.Atoi(line[3][:len(line[3])-1])

		//boardMap[y][x] = -1
	}

	imgfile, err := os.Open(imageFileNameCM)
	if err != nil {
		fmt.Println("image file not found!")
		fmt.Println(imageFileNameCM)
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
		point_list2 = append(point_list2, make([]bool, width))
	}

	for x := 0; x < height; x++ {
		for y := 0; y < width; y++ {
			r, _, _, _ := img.At(x, y).RGBA()
			if r != 0 {
				point_list2[x][y] = true
				point_dict[Tuple{x, y}] = true
				if prnt {
					fmt.Printf("X: %d, Y: %d\n", x, y)
				}

			} else {
				point_dict[Tuple{x, y}] = false
				boardMap[y][x] = -1
			}
		}
	}

	top_left_corner := Coord{x: 0, y: 0}
	top_right_corner := Coord{x: 0, y: 0}
	bot_left_corner := Coord{x: 0, y: 0}
	bot_right_corner := Coord{x: 0, y: 0}

	tl_min := height + width
	tr_max := -1
	bl_max := -1
	br_max := -1

	for x := 0; x < height; x++ {
		for y := 0; y < width; y++ {
			if point_dict[Tuple{x, y}] {
				if x+y < tl_min {
					tl_min = x + y
					top_left_corner.x = x
					top_left_corner.y = y
				}
				if y-x > tr_max {
					tr_max = y - x
					top_right_corner.x = x
					top_right_corner.y = y
				}
				if x-y > bl_max {
					bl_max = x - y
					bot_left_corner.x = x
					bot_left_corner.y = y
				}
				if x+y > br_max {
					br_max = x + y
					bot_right_corner.x = x
					bot_right_corner.y = y
				}
			}
		}
	}

	fmt.Printf("TL: %v, TR %v, BL %v, BR %v\n", top_left_corner, top_right_corner, bot_left_corner, bot_right_corner)

	starting_locs := make([]Coord, 4)
	starting_locs[0] = top_left_corner
	starting_locs[1] = top_right_corner
	starting_locs[2] = bot_left_corner
	starting_locs[3] = bot_right_corner

	statsFile, err := os.Create(outRoutingStatsNameCM)
	if err != nil {
		log.Fatal("Cannot create file", err)
	}
	defer statsFile.Close()

	//New routing initialization
	if regionRouting {

		fmt.Println("Beginning Region Routing")

		/*
			for x := 0; x < 200; x ++ {
				for y := 0; y < 200; y++ {
					if point_dict[Tuple{x, y}] {
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
		//for len(point_list) != 0 {
		for !done {
			top_left := Tuple{-1, -1}
			fmt.Println("starting")
			for x := 0; x < height; x++ {
				for y := 0; y < width; y++ {
					//fmt.Printf("X: %d, Y: %d, v: %d", x, y, point_list2[x][y])
					if point_list2[x][y] {
						top_left = Tuple{x, y}
						break
					}
				}
				if (top_left != Tuple{-1, -1}) {
					break
				}
			}
			fmt.Printf("working %d %d\n", top_left.x, top_left.y)
			if (top_left == Tuple{-1, -1}) {
				done = true
				break
			}
			//top_left := point_list[0]
			temp := Tuple{top_left.x, top_left.y}

			for point_dict[Tuple{temp.x + 1, temp.y}] {
				temp.x += 1
			}

			collide := false
			y_test := Tuple{top_left.x, top_left.y}

			for !collide {
				y_test.y += 1
				if prnt {
					fmt.Println(y_test.y)
				}

				for x_val := top_left.x; x_val < temp.x; x_val++ {
					if !point_dict[Tuple{x_val, y_test.y}] {
						collide = true
					}
				}
			}

			bottom_right := Tuple{temp.x, y_test.y - 1}

			fmt.Println(top_left.x, bottom_right.x, top_left.y, bottom_right.y)

			new_square := RoutingSquare{top_left.x, bottom_right.x, top_left.y, bottom_right.y, true, id_counter, make([]Tuple, 0)}
			id_counter++
			fmt.Println("start_r_square")
			removeRoutingSquare(new_square)
			fmt.Println("end_r_square")
			square_list = append(square_list, new_square)
		}

		length := len(square_list)
		for y, _ := range square_list {
			square := square_list[y]
			square_list[y].routers = make([]Tuple, length)

			for z := y + 1; z < len(square_list); z++ {
				new_square := square_list[z]

				if new_square.x1 >= square.x1 && new_square.x2 <= square.x2 {
					if new_square.y1 == square.y2+1 {
						border_dict[y] = append(border_dict[y], z)
						border_dict[z] = append(border_dict[z], y)

					} else if new_square.y2 == square.y1-1 {
						border_dict[y] = append(border_dict[y], z)
						border_dict[z] = append(border_dict[z], y)
					}
				} else if new_square.y1 >= square.y1 && new_square.y2 <= square.y2 {
					if new_square.x1 == square.x2+1 {
						border_dict[y] = append(border_dict[y], z)
						border_dict[z] = append(border_dict[z], y)

					} else if new_square.x2 == square.x1-1 {
						border_dict[y] = append(border_dict[y], z)
						border_dict[z] = append(border_dict[z], y)
					}
				}
				if square.x1 >= new_square.x1 && square.x2 <= new_square.x2 {
					if square.y1 == new_square.y2+1 {
						border_dict[y] = append(border_dict[y], z)
						border_dict[z] = append(border_dict[z], y)

					} else if square.y2 == new_square.y1-1 {
						border_dict[y] = append(border_dict[y], z)
						border_dict[z] = append(border_dict[z], y)
					}
				} else if square.y1 >= new_square.y1 && square.y2 <= new_square.y2 {
					if square.x1 == new_square.x2+1 {
						border_dict[y] = append(border_dict[y], z)
						border_dict[z] = append(border_dict[z], y)

					} else if square.x2 == new_square.x1-1 {
						border_dict[y] = append(border_dict[y], z)
						border_dict[z] = append(border_dict[z], y)
					}
				}
			}
		}
		fmt.Println(border_dict)

		//Cutting takes place in this loop
		for true {
			rebuilt := false

			for i := 0; i < len(square_list) && !rebuilt; i++ {

				for _, n := range border_dict[i] {

					s_rat := side_ratio(square_list[i], square_list[n])
					if s_rat > 0.6 {
						new_squares := single_cut(square_list[i], square_list[n])

						s1 := square_list[n]
						s2 := square_list[i]

						square_list_remove(s1)
						square_list_remove(s2)

						square_list = append(square_list, new_squares...)

						rebuild(square_list)

						rebuilt = true

						break
					}

					a_rat := area_ratio(square_list[i], square_list[n])
					if a_rat > 0.6 {
						new_squares := double_cut(square_list[i], square_list[n])

						s1 := square_list[n]
						s2 := square_list[i]

						square_list_remove(s1)
						square_list_remove(s2)

						new_squares[2].id_num = len(square_list)

						square_list = append(square_list, new_squares...)

						rebuild(square_list)

						rebuilt = true

						break
					}
				}
			}

			if !rebuilt {
				break
			}
		}

		/*
			type Changeable interface {
				Set(x, y int, c color.Color)
			}

			if cimg, ok := img.(Changeable); ok {

				for x := 0; x < len(square_list); x++ {
					r := uint8(rand.Intn(255))
					g := uint8(rand.Intn(255))
					b := uint8(rand.Intn(255))

					fmt.Println(square_list[x], r, g, b)
					for xx := square_list[x].x1; xx <= square_list[x].x2; xx++ {
						for yy := square_list[x].y1; yy <= square_list[x].y2; yy++ {
							cimg.Set(xx, yy, color.RGBA{r, g, b, 255})
						}
					}
				}
				cimg.Set(199, 14, color.RGBA{255, 255, 255, 255})
				cimg.Set(200, 14, color.RGBA{255, 255, 255, 255})
				//cimg.Set(193, 341, color.RGBA{255, 255, 255, 255})
				//cimg.Set(341, 193, color.RGBA{255, 255, 255, 255})

				ff, err := os.Create("img.png")
				if err != nil {
					panic(err)
				}
				defer ff.Close()
				png.Encode(ff, img)

			} else {
				fmt.Println("can't edit image :(")
			}
		*/

		node_tables = make([]map[Tuple]float64, len(square_list))

		for key, values := range border_dict {
			if key < len(square_list) {
				node_tables[key] = make(map[Tuple]float64)
				if len(values) > 1 {
					for n := 0; n < len(values); n++ {
						next := n + 1
						for next < len(values) {
							node_a := border_dict[key][n]
							node_b := border_dict[key][next]

							p1 := square_list[key].routers[node_a]
							p2 := square_list[key].routers[node_b]

							node_tables[key][Tuple{node_a, node_b}] = dist(p1, p2)
							node_tables[key][Tuple{node_b, node_a}] = dist(p1, p2)

							next += 1
						}
					}
				}
			}
		}

		routingFile, err := os.Create(outRoutingNameCM)
		if err != nil {
			log.Fatal("Cannot create file", err)
		}
		defer routingFile.Close()

		//The scheduler determines which supernode should pursue a point of interest
		scheduler := &Scheduler{}

		//List of all the supernodes on the grid
		scheduler.sNodeList = make([]SuperNodeParent, numSuperNodes)

		for i := 0; i < numSuperNodes; i++ {
			snode_points := make([]Coord, 1)
			snode_path := make([]Coord, 0)
			all_points := make([]Coord, 0)

			//Defining the starting x and y values for the super node
			//This super node starts at the middle of the grid
			x_val, y_val := starting_locs[i].x, starting_locs[i].y
			nodeCenter := Coord{x: x_val, y: y_val}

			scheduler.sNodeList[i] = &sn_zero{&supern{&NodeImpl{x: x_val, y: y_val, id: i}, 1,
				1, superNodeRadius, superNodeRadius, 0, snode_points, snode_path,
				nodeCenter, 0, 0, 0, 0, 0, all_points}}

			//The super node's current location is always the first element in the routePoints list
			scheduler.sNodeList[i].updateLoc()
		}

		fmt.Printf("Iteration %d/%v", 0, iterations_of_event)
		for i := 0; i < iterations_of_event; i++ {
			fmt.Printf("\rIteration %d/%v", i, iterations_of_event)

			if t, ok := stim_list[i]; ok {
				start := time.Now()
				scheduler.addRoutePoint(Coord{x: t.x, y: t.y})
				elapsed := time.Since(start)

				fmt.Fprint(statsFile, "Region Routing Elapsed: ", elapsed, "\n")

				fmt.Printf("\nAdding %d %d %d\n", i, t.x, t.y)
			}

			for _, s := range scheduler.sNodeList {
				points_len := len(s.getRoutePoints())
				response_time := -1
				if points_len > 1 {
					response_time = s.getRoutePoints()[1].time
				}

				s.tick()

				if points_len > len(s.getRoutePoints()) {
					fmt.Fprint(statsFile, "Response Time: ", response_time, "\n")
				}

				fmt.Fprint(routingFile, s)
				p := printPoints(s)
				fmt.Fprint(routingFile, " UnvisitedPoints: ")
				fmt.Fprintln(routingFile, p.String())
			}
		}
		for _, s := range scheduler.sNodeList {
			fmt.Println("\nSQUARES MOVED", s.getSquaresMoved())
			fmt.Println("\nPOINTS VISITED", s.getPointsVisited())

			fmt.Fprintf(statsFile, "Squares Moved: %d %d\n", s.getId(), s.getSquaresMoved())
			fmt.Fprintf(statsFile, "Points Visited: %d %d\n", s.getId(), s.getPointsVisited())
		}

	}
	//End new routing initialization

	//aStar routing
	if aStarRouting {

		fmt.Println("Beginning aStar Routing")

		routingFile, err := os.Create(outRoutingNameCM)
		if err != nil {
			log.Fatal("Cannot create file", err)
		}
		defer routingFile.Close()

		//The scheduler determines which supernode should pursue a point of interest
		scheduler := &Scheduler{}

		//List of all the supernodes on the grid
		scheduler.sNodeList = make([]SuperNodeParent, numSuperNodes)

		for i := 0; i < numSuperNodes; i++ {
			snode_points := make([]Coord, 1)
			snode_path := make([]Coord, 0)
			all_points := make([]Coord, 0)

			//Defining the starting x and y values for the super node
			//This super node starts at the middle of the grid
			x_val, y_val := starting_locs[i].x, starting_locs[i].y
			nodeCenter := Coord{x: x_val, y: y_val}

			scheduler.sNodeList[i] = &sn_zero{&supern{&NodeImpl{x: x_val, y: y_val, id: i}, 1,
				1, superNodeRadius, superNodeRadius, 0, snode_points, snode_path,
				nodeCenter, 0, 0, 0, 0, 0, all_points}}

			//The super node's current location is always the first element in the routePoints list
			scheduler.sNodeList[i].updateLoc()
		}

		fmt.Print("POINTS", scheduler.sNodeList[0].getPointsVisited())

		fmt.Printf("Iteration %d/%v", 0, iterations_of_event)
		for i := 0; i < iterations_of_event; i++ {
			fmt.Printf("\rIteration %d/%v", i, iterations_of_event)

			if t, ok := stim_list[i]; ok {
				//scheduler.addRoutePoint(Coord{x: t.x, y: t.y})
				start := time.Now()
				scheduler.addRoutePoint(Coord{x: t.x, y: t.y})
				elapsed := time.Since(start)

				fmt.Fprint(statsFile, "aStar Routing Elapsed ", elapsed, "\n")

				fmt.Printf("\nAdding %d %d %d\n", i, t.x, t.y)
			}

			for _, s := range scheduler.sNodeList {
				points_len := len(s.getRoutePoints())
				response_time := -1
				if points_len > 1 {
					response_time = s.getRoutePoints()[1].time
				}

				s.tick()

				if points_len > len(s.getRoutePoints()) {
					fmt.Fprint(statsFile, "Response Time: ", response_time, "\n")
				}

				fmt.Fprint(routingFile, s)
				p := printPoints(s)
				fmt.Fprint(routingFile, " UnvisitedPoints: ")
				fmt.Fprintln(routingFile, p.String())
			}
		}
		for _, s := range scheduler.sNodeList {
			fmt.Println("\nSQUARES MOVED", s.getSquaresMoved())
			fmt.Println("\nPOINTS VISITED", s.getPointsVisited())

			fmt.Fprintf(statsFile, "Distance Traveled: %d %d\n", s.getId(), s.getSquaresMoved())
			fmt.Fprintf(statsFile, "Points Visited: %d %d\n", s.getId(), s.getPointsVisited())
		}
	}
}

//This function allows the simulator to create a roadMap of the grid
//Every Coord in the grid is given an integer value corresponding to the
//	number of times the Coord is used by all paths
//The function first generates two random Coords on each half of the grid
//It then finds the path between those Coords
//It then increments the integer value of each Coord in the path by one
//This is done an amount of time to generate a conclusive distribution of paths
//	across the gird
//THe resulting roadMap is outputted to the file, first with the max number
//	if times a Coord is visited and then each Coord's integer value
func makeRoads(roadFile *os.File) {
	//This map has Tuples as keys and integers as values
	//The Tuples represent the Coord in the grid and the integer represents
	//	the amount of times the Coord is visited by all paths
	roadMap := make(map[Tuple]int)

	//The max is kept track of the be outputted at the beginning of the
	//road output file
	//This is used to determine the gradient of color by the Viewer when
	//	displaying the roads
	max := 0

	aStarIterations := 100

	fmt.Printf("Running ASTAR iteration %d/%v", 0, aStarIterations)
	for i := 0; i < aStarIterations; i++ {
		//Two Coords are randomly generated
		a := Coord{nil, rangeInt(0, maxX), rangeInt(0, maxY), 0, 0, 0, 0}
		b := Coord{nil, rangeInt(0, maxX), rangeInt(0, maxY), 0, 0, 0, 0}

		//The Coords' x and y positions are randomly updated to be on either the
		//	left and right side, or top and bottom of the grid
		//This is done to ensure the paths between these Coords crosses a large
		//	section of the grid
		if i%2 == 0 {
			a.x = rangeInt(0, maxX/2)
			b.x = rangeInt(maxX/2, maxX)
		} else {
			a.y = rangeInt(0, maxY/2)
			b.y = rangeInt(maxY/2, maxY)
		}
		fmt.Printf("\rRunning ASTAR iteration %d/%v", i, aStarIterations)
		//The aStar path between these two Coords is created
		//Each Coord in this path is looped through and the integer value corresponding
		//	to that Coord is incremented by one
		for _, r := range aStar(a, b) {
			pos := Tuple{r.x, r.y}
			roadMap[pos]++
			if roadMap[pos] > max {
				max = roadMap[pos]
			}
		}
	}
	fmt.Fprintln(roadFile, "max", max)

	//This loops through the roadMap and outputs the integer value for every Coord
	//	in the grid
	for i := 0; i < maxX; i++ {
		for j := 0; j < maxY; j++ {
			fmt.Println("\rOutputting to roadLog: Coord", j, i)
			if boardMap[i][j] == -1 {
				fmt.Fprintln(roadFile, i, j, -1)
			} else {
				fmt.Fprintln(roadFile, i, j, roadMap[Tuple{j, i}])
			}
		}
	}
}

//Saves the Coords in the allPoints list into a buffer to
//	print to the file
func printPoints(s SuperNodeParent) bytes.Buffer {
	var buffer bytes.Buffer
	buffer.WriteString((fmt.Sprintf("[")))
	for ind, i := range s.getAllPoints() {
		buffer.WriteString(i.String())

		if ind != len(s.getAllPoints())-1 {
			buffer.WriteString(" ")
		}
	}
	buffer.WriteString((fmt.Sprintf("]")))
	return buffer
}

func getFlags() {
	flag.IntVar(&iterations_of_event, "iterations_of_event", 1000, "how many times the simulation will run")

	//fmt.Println(os.Args[1:], "\nhmmm? \n ") //C:\Users\Nick\Desktop\comand line experiments\src
	flag.IntVar(&negativeSittingStopThresholdCM, "negativeSittingStopThreshold", -10,
		"Negative number sitting is set to when board map is reset")
	flag.IntVar(&sittingStopThresholdCM, "sittingStopThreshold", 5,
		"How long it takes for a node to stay seated")
	flag.Float64Var(&gridCapacityPercentageCM, "gridCapacityPercentage", .9,
		"Percent the sub-grid can be filled")
	flag.StringVar(&inputFileNameCM, "inputFileName", "Log1_in.txt",
		"Name of the input text file")
	flag.StringVar(&outputFileNameCM, "outputFileName", "Log",
		"Name of the output text file prefix")
	flag.Float64Var(&naturalLossCM, "naturalLoss", .005,
		"battery loss due to natural causes")
	flag.Float64Var(&sensorSamplingLossCM, "sensorSamplingLoss", .001,
		"battery loss due to sensor sampling")
	flag.Float64Var(&GPSSamplingLossCM, "GPSSamplingLoss", .005,
		"battery loss due to GPS sampling")
	flag.Float64Var(&serverSamplingLossCM, "serverSamplingLoss", .01,
		"battery loss due to server sampling")
	flag.IntVar(&thresholdBatteryToHaveCM, "thresholdBatteryToHave", 30,
		"Threshold battery phones should have")
	flag.IntVar(&thresholdBatteryToUseCM, "thresholdBatteryToUse", 10,
		"Threshold of battery phones should consume from all forms of sampling")
	flag.IntVar(&movementSamplingSpeedCM, "movementSamplingSpeed", 20,
		"the threshold of speed to increase sampling rate")
	flag.IntVar(&movementSamplingPeriodCM, "movementSamplingPeriod", 1,
		"the threshold of speed to increase sampling rate")
	flag.IntVar(&maxBufferCapacityCM, "maxBufferCapacity", 25,
		"maximum capacity for the buffer before it sends data to the server")
	flag.StringVar(&energyModelCM, "energyModel", "variable",
		"this determines the energy loss model that will be used")
	flag.BoolVar(&noEnergyModelCM, "noEnergy", false,
		"Whether or not to ignore energy model for simulation")
	flag.IntVar(&sensorSamplingPeriodCM, "sensorSamplingPeriod", 1000,
		"rate of the sensor sampling period when custom energy model is chosen")
	flag.IntVar(&GPSSamplingPeriodCM, "GPSSamplingPeriod", 1000,
		"rate of the GridGPS sampling period when custom energy model is chosen")
	flag.IntVar(&serverSamplingPeriodCM, "serverSamplingPeriod", 1000,
		"rate of the server sampling period when custom energy model is chosen")
	flag.IntVar(&numStoredSamplesCM, "nodeStoredSamples", 10,
		"number of samples stored by individual nodes for averaging")
	flag.IntVar(&gridStoredSamplesCM, "gridStoredSamples", 10,
		"number of samples stored by grid squares for averaging")
	flag.Float64Var(&detectionThresholdCM, "detectionThreshold", 11180.0,
		"Value where if a node gets this reading or higher, it will trigger a detection")
	flag.Float64Var(&errorModifierCM, "errorMultiplier", 1.0,
		"Multiplier for error values in system")
	//Range 1, 2, or 4
	//Currently works for only a few numbers, can be easily expanded but is not currently dynamic
	flag.IntVar(&numSuperNodes, "numSuperNodes", 1, "the number of super nodes in the simulator")

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
	flag.IntVar(&superNodeType, "superNodeType", 0, "the type of super node used in the simulator")

	//Range: 0-...
	//Theoretically could be as high as possible
	//Realistically should remain around 10
	flag.IntVar(&superNodeSpeed, "superNodeSpeed", 1, "the speed of the super node")

	//Range: true/false
	//Tells the simulator whether or not to optimize the path of all the super nodes
	//Only works when multiple super nodes are active in the same area
	flag.BoolVar(&doOptimize, "doOptimize", false, "whether or not to optimize the simulator")

	//Range: 0-4
	//	0: begin in center, routed anywhere
	//	1: begin inside circles located in the corners, only routed inside circle
	//	2: begin inside circles located on the sides, only routed inside circle
	//	3: being inside large circles located in the corners, only routed inside circle
	//	4: begin inside regions, only routed inside region
	//Only used for super nodes of type 1
	//flag.IntVar(&superNodeVariation, "superNodeVariation", 3, "super nodes of type 1 have different variations")

	flag.BoolVar(&positionPrintCM, "logPosition", false, "Whether you want to write position info to a log file")
	flag.BoolVar(&gridPrintCM, "logGrid", false, "Whether you want to write grid info to a log file")
	flag.BoolVar(&energyPrintCM, "logEnergy", false, "Whether you want to write energy into to a log file")
	flag.BoolVar(&nodesPrintCM, "logNodes", false, "Whether you want to write node readings to a log file")
	flag.IntVar(&squareRowCM, "squareRow", 100, "Number of rows of grid squares, 1 through maxX")
	flag.IntVar(&squareColCM, "squareCol", 100, "Number of columns of grid squares, 1 through maxY")

	flag.BoolVar(&regionRouting, "regionRouting", false, "True if you want to use the new routing algorithm with regions and cutting")
	flag.BoolVar(&aStarRouting, "aStarRouting", false, "True if you want to use the old routing algorithm with aStar")

	flag.StringVar(&imageFileNameCM, "imageFileName", "circle_justWalls_x4.png", "Name of the input text file")
	flag.StringVar(&stimFileNameCM, "stimFileName", "circle_0.txt", "Name of the stimulus text file")
	flag.StringVar(&outRoutingNameCM, "outRoutingName", "log.txt", "Name of the stimulus text file")
	flag.StringVar(&outRoutingStatsNameCM, "outRoutingStatsName", "routingStats.txt", "Name of the output file for stats")

	flag.Parse()
}

func rangeInt(min, max int) int { //returns a random number between max and min
	return rand.Intn(max-min) + min
}
