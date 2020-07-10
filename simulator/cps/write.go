//CPS contains the simulator
package cps

import (
	"archive/zip"
	"bytes"
	"encoding/csv"
	"flag"
	"io"
	"math/rand"
	"strings"
	//"bufio"
	"fmt"
	"image"
	"io/ioutil"
	"math"
	"os"
	"regexp"
	"strconv"
	//"Time"
	"log"
)

var (
	//iterations_used int = 0

	// variables for making maps
	bufferCurrent = [][]int{{2, 2}, {0, 0}} // points currently being worked with
	bufferFuture  = [][]int{{}}             // point to be worked with
	starter       = 1                       // This is the destination number
	/*p.BoardMap      = [][]int{                // This is the map with all the position variables for path finding
	{0, 0, 0, 0, 0},
	{0, 0, 0, 0, 0},
	{0, 0, 0, 0, 0},
	{0, 0, 0, 0, 0},
	{0, 0, 0, 0, 0}}*/

	wallPoints = [][]int{{1, 1}, {1, 2}, {1, 3}, {2, 1}}
	// end variables for making maps

	/*
		npos    [][]int // node positions
		wpos    [][]int // wall positions
		spos    [][]int // super node positions
		ppos    [][]int // super node points of interest positions
		poikpos [][]int // points of interest kinetic
		poispos [][]int // points of interest static
	*/

	//numNodeNodes               int
	//numWallNodes               int
	//numPoints                  int
	//numPointsOfInterestKinetic int
	//numPointsOfInterestStatic  int

	//fileName = "Log1_in.txt"

	makeBoardMapFile = true
)

//GetDashedInput reads numbers from a file which contain values
//for elements like the number of nodes.
func GetDashedInput(s string, p *Params) int {
	b := ReadFromFile(p.FileName)
	r := regexp.MustCompile(string(s + "-[0-9]+"))
	w := r.FindAllString(string(b), 1)
	r = regexp.MustCompile("[0-9]+")
	w = r.FindAllString(w[0], 10)
	s1, err := strconv.ParseInt(w[len(w)-1], 10, 32)
	Check(err)
	return int(s1)
}

//getString reads the input file and extracts the data from the specified category
//headExp is a regular expression to denote the label and dataExp is the form the data will take
func getString(p *Params, bytes []byte, headExp string, dataExp string) ([][]int, []string){
	regex := regexp.MustCompile(headExp)
	text := regex.FindAllString(string(bytes), 10)
	regex = regexp.MustCompile("[0-9]+")
	text = regex.FindAllString(text[0], 10)
	size, err := strconv.ParseInt(text[0], 10, 32)
	Check(err)
	regex = regexp.MustCompile(dataExp)
	fai := regex.FindAllIndex(bytes, int(size))
	text = regex.FindAllString(string(bytes), int(size))
	return fai, text
}

//GetListedInput reads from the Scenario file to determine the node number
//the number and location of attractions, and other relevant data
//that is stored in that text file.
func GetListedInput(p *Params) {
	var temp []byte
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(dir)
	fileBytes := ReadFromFile(p.FileName)

	fai, text := getString(p, fileBytes, "N: [0-9]+", "x:[0-9]+, y:[0-9]+, t:[0-9]+")
	if len(fai) > 0 {
		temp = fileBytes[fai[len(fai)-1][1]:]
	} else {
		temp = fileBytes
	}
	FillInts(text, 0, p)

	fai, text = getString(p, temp, "W: [0-9]+", "x:[0-9]+, y:[0-9]+")
	if len(fai) > 0 {
		temp = temp[fai[len(fai)-1][1]:]
	}
	FillInts(text, 1, p)

	fai, text = getString(p, temp, "S: [0-9]+", "x:[0-9]+, y:[0-9]+")
	if len(fai) > 0 {
		temp = temp[fai[len(fai)-1][1]:]
	}
	FillInts(text, 2, p)

	fai, text = getString(p, temp, "POIS: [0-9]+", "x:[0-9]+, y:[0-9]+, ti:[0-9]+, to:[0-9]+")
	FillInts(text, 5, p)

	p.CurrentNodes = len(p.NodeEntryTimes)
	p.NumWallNodes = len(p.Wpos)
	//numPoints = len(ppos)
	//fmt.Println(p.NumNodeNodes, p.NumWallNodes, p.NumPointsOfInterestStatic)
}

//FillInts reads the scenario file for the integer values containing the locations of
//nodes, attractions, and their entry times
func FillInts(s []string, place int, p *Params) {
	if place == 0 {
		for i := 0; i < len(s); i++ {
			r := regexp.MustCompile("x:[0-9]+")
			X := r.FindAllString(s[i], 1)
			r = regexp.MustCompile("[0-9]+")
			x := r.FindAllString(X[0], 1)
			s1, err := strconv.ParseInt(x[0], 10, 32)
			Check(err)
			ap := []int{int(s1), 0, 0}
			p.NodeEntryTimes = append(p.NodeEntryTimes, ap)

			r = regexp.MustCompile("y:[0-9]+")
			X = r.FindAllString(s[i], 1)
			r = regexp.MustCompile("[0-9]+")
			x = r.FindAllString(X[0], 1)
			s1, err = strconv.ParseInt(x[0], 10, 32)
			Check(err)
			p.NodeEntryTimes[i][1] = int(s1)

			r = regexp.MustCompile("t:[0-9]+")
			X = r.FindAllString(s[i], 1)
			r = regexp.MustCompile("[0-9]+")
			x = r.FindAllString(X[0], 1)
			s1, err = strconv.ParseInt(x[0], 10, 32)
			Check(err)
			p.NodeEntryTimes[i][2] = int(s1)
		}
	} else if place == 1 {
		for i := 0; i < len(s); i++ {
			r := regexp.MustCompile("x:[0-9]+")
			X := r.FindAllString(s[i], 1)
			r = regexp.MustCompile("[0-9]+")
			x := r.FindAllString(X[0], 1)
			s1, err := strconv.ParseInt(x[0], 10, 32)
			Check(err)
			ap := []int{int(s1), 0, 0}
			p.Wpos = append(p.Wpos, ap)

			r = regexp.MustCompile("y:[0-9]+")
			X = r.FindAllString(s[i], 1)
			r = regexp.MustCompile("[0-9]+")
			x = r.FindAllString(X[0], 1)
			s1, err = strconv.ParseInt(x[0], 10, 32)
			Check(err)
			p.Wpos[i][1] = int(s1)

		}
	} else if place == 2 {
		for i := 0; i < len(s); i++ {
			r := regexp.MustCompile("x:[0-9]+")
			X := r.FindAllString(s[i], 1)
			r = regexp.MustCompile("[0-9]+")
			x := r.FindAllString(X[0], 1)
			s1, err := strconv.ParseInt(x[0], 10, 32)
			Check(err)
			ap := []int{int(s1), 0, 0}
			p.Wpos = append(p.Spos, ap)

			r = regexp.MustCompile("y:[0-9]+")
			X = r.FindAllString(s[i], 1)
			r = regexp.MustCompile("[0-9]+")
			x = r.FindAllString(X[0], 1)
			s1, err = strconv.ParseInt(x[0], 10, 32)
			Check(err)
			p.Wpos[i][1] = int(s1)

			r = regexp.MustCompile("t:[0-9]+")
			X = r.FindAllString(s[i], 1)
			r = regexp.MustCompile("[0-9]+")
			x = r.FindAllString(X[0], 1)
			s1, err = strconv.ParseInt(x[0], 10, 32)
			Check(err)
			p.Spos[i][2] = int(s1)
		}
	} else if place == 3 {
		for i := 0; i < len(s); i++ {
			r := regexp.MustCompile("x:[0-9]+")
			X := r.FindAllString(s[i], 1)
			r = regexp.MustCompile("[0-9]+")
			x := r.FindAllString(X[0], 1)
			s1, err := strconv.ParseInt(x[0], 10, 32)
			Check(err)
			ap := []int{int(s1), 0, 0}
			p.Ppos = append(p.Ppos, ap)

			r = regexp.MustCompile("y:[0-9]+")
			X = r.FindAllString(s[i], 1)
			r = regexp.MustCompile("[0-9]+")
			x = r.FindAllString(X[0], 1)
			s1, err = strconv.ParseInt(x[0], 10, 32)
			Check(err)
			p.Ppos[i][1] = int(s1)

			r = regexp.MustCompile("t:[0-9]+")
			X = r.FindAllString(s[i], 1)
			r = regexp.MustCompile("[0-9]+")
			x = r.FindAllString(X[0], 1)
			s1, err = strconv.ParseInt(x[0], 10, 32)
			Check(err)
			p.Ppos[i][2] = int(s1)
		}
	} else if place == 5 {
		for i := 0; i < len(s); i++ {
			r := regexp.MustCompile("x:[0-9]+")
			X := r.FindAllString(s[i], 1)
			r = regexp.MustCompile("[0-9]+")
			x := r.FindAllString(X[0], 1)
			s1, err := strconv.ParseInt(x[0], 10, 32)
			Check(err)
			ap := []int{int(s1), 0, 0, 0}
			p.Poispos = append(p.Poispos, ap)

			r = regexp.MustCompile("y:[0-9]+")
			X = r.FindAllString(s[i], 1)
			r = regexp.MustCompile("[0-9]+")
			x = r.FindAllString(X[0], 1)
			s1, err = strconv.ParseInt(x[0], 10, 32)
			Check(err)
			p.Poispos[i][1] = int(s1)

			r = regexp.MustCompile("ti:[0-9]+")
			X = r.FindAllString(s[i], 1)
			r = regexp.MustCompile("[0-9]+")
			x = r.FindAllString(X[0], 1)
			s1, err = strconv.ParseInt(x[0], 10, 32)
			Check(err)
			p.Poispos[i][2] = int(s1)

			r = regexp.MustCompile("to:[0-9]+")
			X = r.FindAllString(s[i], 1)
			r = regexp.MustCompile("[0-9]+")
			x = r.FindAllString(X[0], 1)
			s1, err = strconv.ParseInt(x[0], 10, 32)
			Check(err)
			p.Poispos[i][3] = int(s1)
		}
	}
}

// Returns the char number associated with a byte
func GetIntFromByte(a byte) int {
	if a <= 57 && a >= 48 {
		return int(a - 48)
	} else {
		return -1
	}
}

// Returns the string character of a byte
func GetLetterFromByte(a byte) string {
	return string([]byte{a})
}

// Clears file then writes message
func WriteToFile(name string, message string) {
	d1 := []byte(message)
	err := ioutil.WriteFile(name, append(ReadFromFile(name), d1...), 0644)
	Check(err)
}

// Reads entire file to array of bytes
func ReadFromFile(name string) (b []byte) {
	b, err := ioutil.ReadFile(name)
	Check(err)
	return
}

// Creates a file file with specific name
func CreateFile(name string) {
	file, err := os.Create(name) // creates text file
	Check(err)                   // Checks if text file is created properly
	file.Close()                 // closes the file at the end
}

// Checks an error
func Check(e error) {
	if e != nil {
		fmt.Println(e)
	}
}

// Creates p.BoardMap
func CreateBoard(x int, y int, p *Params) {
	p.BoardMap = [][]int{}
	for i := 0; i < x; i++ {
		p.BoardMap = append(p.BoardMap, []int{})
		for j := 0; j < y; j++ {
			p.BoardMap[i] = append(p.BoardMap[i], 0)
		}
	}
}

//HandleMovement adjusts BoolGrid when nodes move around the map
func HandleMovement(p *Params) {
	for j := 0; j < len(p.NodeList); j++ {

		oldX, oldY := p.NodeList[j].GetLoc()
		p.BoolGrid[int(oldX)][int(oldY)] = false //set the old spot false since the node will now move away

		//move the node to its new location
		p.NodeList[j].Move(p)

		//set the new location in the boolean field to true
		newX, newY := p.NodeList[j].GetLoc()
		p.BoolGrid[int(newX)][int(newY)] = true

		//sync the QuadTree
		if p.ClusteringOn {
			p.NodeTree.NodeMovement(p.NodeList[j])
		}

		//writes the node information to the file
		if p.EnergyPrint {
			fmt.Fprintln(p.EnergyFile, p.NodeList[j])
		}

		//Add the node into its new Square's p.TotalNodes
		//If the node hasn't left the square, that Square's p.TotalNodes will
		//remain the same after these calculations
	}
}

//HandleMovementCSV does the same as HandleMovement
func HandleMovementCSV(p *Params) {
	time := p.Iterations_used
	for j := 0; j < len(p.NodeList); j++ {

		if p.NodeList[j].Valid {
			oldX, oldY := p.NodeList[j].GetLoc()
			p.BoolGrid[int(oldX)][int(oldY)] = false //set the old spot false since the node will now move away
		}
		//move the node to its new location
		//p.NodeList[j].Move(p)

		id := p.NodeList[j].GetID()
		p.NodeList[j].X = float32(p.NodeMovements[id][time].X)
		p.NodeList[j].Y = float32(p.NodeMovements[id][time].Y)


		//set the new location in the boolean field to true
		newX, newY := p.NodeList[j].GetLoc()

		//sync the QuadTree
		if p.ClusteringOn {
			p.NodeTree.NodeMovement(p.NodeList[j])
		}

		if p.NodeList[j].InBounds(p) {
			p.NodeList[j].Valid = true
		} else {
			p.NodeList[j].Valid = false
		}
		if p.NodeList[j].Valid {
			p.BoolGrid[int(newX)][int(newY)] = true
		}

		//writes the node information to the file
		if p.EnergyPrint {
			fmt.Fprintln(p.EnergyFile, p.NodeList[j])
		}

		//Add the node into its new Square's p.TotalNodes
		//If the node hasn't left the square, that Square's p.TotalNodes will
		//remain the same after these calculations
	}
}

//InitializeNodeParameters sets all the defaulted node values to semi-random values
func InitializeNodeParameters(p *Params, nodeNum int) *NodeImpl{

	var initHistory = make([]float32, p.NumStoredSamples)

	//initialize nodes to invalid starting point as starting point will be selected after initialization
	curNode := NodeImpl{P: p, X: -1, Y: -1, Id: len(p.NodeList), SampleHistory: initHistory, Concentration: 0, Cascade: nodeNum}

	//values to determine coefficients
	curNode.SetS0(rand.Float64()*0.005 + 0.33)
	curNode.SetS1(rand.Float64()*0.005 + 0.33)
	curNode.SetS2(rand.Float64()*0.005 + 0.33)
	//values to determine error in coefficients
	s0, s1, s2 := curNode.GetCoefficients()
	curNode.SetE0(rand.Float64() * 0.02 * p.ErrorModifierCM * s0)
	curNode.SetE1(rand.Float64() * 0.02 * p.ErrorModifierCM * s1)
	curNode.SetE2(rand.Float64() * 0.02 * p.ErrorModifierCM * s2)
	//Values to determine error in exponents
	curNode.SetET1(p.Tau1 * rand.Float64() * p.ErrorModifierCM * 0.05)
	curNode.SetET2(p.Tau1 * rand.Float64() * p.ErrorModifierCM * 0.05)

	//set node Time and initial sensitivity
	curNode.NodeTime = 0
	curNode.InitialSensitivity = s0 + (s1)*math.Exp(-float64(curNode.NodeTime)/p.Tau1) + (s2)*math.Exp(-float64(curNode.NodeTime)/p.Tau2)
	curNode.Sensitivity = curNode.InitialSensitivity

	// Initialize New Battery Model Variables
	curNode.CurrentBatteryLevel = int(float64(p.BatteryCapacity) * RandomBatteryLevel(p.AverageBatteryLevel))
	curNode.InitialBatteryLevel = curNode.CurrentBatteryLevel
	curNode.SamplingPeriod		= p.SamplingPeriodMS
	return &curNode
}

func SetupCSVNodes(p *Params) {
	for i := 0; i < p.TotalNodes; i++ {
		newNode := InitializeNodeParameters(p, i)

		newNode.X = float32(p.NodeMovements[i][0].X)
		newNode.Y = float32(p.NodeMovements[i][0].Y)

		if newNode.InBounds(p) {
			newNode.Valid = true
			//fmt.Printf("Valid NODE %v %v %v\n", newNode.Id, newNode.X, newNode.Y)
			p.BoolGrid[int(newNode.X)][int(newNode.Y)] = true
		} else {
			newNode.Valid = false
		}

		p.NodeList = append(p.NodeList, newNode)
		p.AliveNodes[newNode] = true
		p.CurrentNodes += 1

		if p.ClusteringOn {
			newNode.IsClusterHead = false
			newNode.IsClusterMember = false
			newNode.NodeClusterParams = &ClusterMemberParams{}
			newNode.OutOfRange = false
			p.NodeTree.Insert(newNode)
			p.ClusterNetwork.ClearClusterParams(newNode)
			//newNode.TimeLastSentReadings = p.CurrentTime
			//newNode.ReadingsBuffer = []Reading{}
		}

		newNode.AccelerometerSpeed = []float32{}
		//newNode.TimeLastAccel = p.CurrentTime
		//newNode.LastMoveTime = p.CurrentTime

		p.Events.Push(&Event{newNode, SENSE, 0, 0})
		p.Events.Push(&Event{newNode, MOVE, 0, 0})

	}

}
//SetupRandomNodes creates random nodes and appends them to the node list
func SetupRandomNodes(p *Params) {
	for i := 0; i < len(p.NodeEntryTimes); i++ {
		if p.Iterations_used == p.NodeEntryTimes[i][2] {

			newNode := InitializeNodeParameters(p, i)

			xx := rangeInt(1, p.MaxX)
			yy := rangeInt(1, p.MaxY)
			for p.BoolGrid[xx][yy] == true {
				xx = rangeInt(1, p.MaxX)
				yy = rangeInt(1, p.MaxY)
			}

			newNode.X = float32(xx)
			newNode.Y = float32(yy)

			newNode.Valid = true
			p.BoolGrid[xx][yy] = true

			p.NodeList = append(p.NodeList, newNode)
			p.CurrentNodes += 1

		}
	}
}


// Fills the walls into the board based on the wall positions extrapolated from the file
func FillInWallsToBoard(p *Params) {
	for i := 0; i < len(p.Wpos); i++ {
		p.BoardMap[p.Wpos[i][0]][p.Wpos[i][1]] = -1
	}
}

// Fills the points of interest into the current buffer

func FillInBufferCurrent(p *Params) {
	bufferCurrent = [][]int{}
	for i := 0; i < len(p.Poispos); i++ {
		if p.Iterations_used >= p.Poispos[i][2] && p.Iterations_used < p.Poispos[i][3] {
			bufferCurrent = append(bufferCurrent, []int{p.Poispos[i][0], p.Poispos[i][1]})
			//fmt.Println("1ho- ", iterations_used, "2ho", bufferCurrent)
		}
	}
}

// Fills the points of interest to the board
func FillPointsToBoard(p *Params) {
	for i := 0; i < len(bufferCurrent); i++ {
		//fmt.Println(bufferCurrent)
		p.BoardMap[bufferCurrent[i][0]][bufferCurrent[i][1]] = starter
	}
}

//FillInMap fills BoardMap with the appropriate values
func FillInMap(p *Params) {
	/*start := Time.Now()

	defer func() {
		elapsed := Time.Since(start)
		//fmt.Println("Board Map took", elapsed)
	}()*/

	for len(bufferFuture) > 0 {
		bufferFuture = [][]int{}
		for i := 0; i < len(bufferCurrent); i++ {
			// empty buffer future
			//Check above
			//fmt.Println(len(p.BoardMap[1]),i)
			if bufferCurrent[i][0]-1 < len(p.BoardMap) &&
				bufferCurrent[i][1] < len(p.BoardMap[1]) &&
				bufferCurrent[i][0]-1 >= 0 &&
				bufferCurrent[i][1] >= 0 &&
				p.BoardMap[bufferCurrent[i][0]-1][bufferCurrent[i][1]] == 0 {

				p.BoardMap[bufferCurrent[i][0]-1][bufferCurrent[i][1]] = starter + 1
				bufferFuture = append(bufferFuture, []int{bufferCurrent[i][0] - 1, bufferCurrent[i][1]})
			}
			//Check below
			if bufferCurrent[i][0]+1 < len(p.BoardMap) &&
				bufferCurrent[i][1] < len(p.BoardMap[1]) && // p.BoardMap[1] to
				bufferCurrent[i][0]+1 >= 0 &&
				bufferCurrent[i][1] >= 0 &&
				p.BoardMap[bufferCurrent[i][0]+1][bufferCurrent[i][1]] == 0 {

				p.BoardMap[bufferCurrent[i][0]+1][bufferCurrent[i][1]] = starter + 1
				bufferFuture = append(bufferFuture, []int{bufferCurrent[i][0] + 1, bufferCurrent[i][1]})
			}
			// to the right
			if bufferCurrent[i][0] < len(p.BoardMap) &&
				bufferCurrent[i][1]+1 < len(p.BoardMap[1]) &&
				bufferCurrent[i][0] >= 0 &&
				bufferCurrent[i][1]+1 >= 0 &&
				p.BoardMap[bufferCurrent[i][0]][bufferCurrent[i][1]+1] == 0 {

				p.BoardMap[bufferCurrent[i][0]][bufferCurrent[i][1]+1] = starter + 1
				bufferFuture = append(bufferFuture, []int{bufferCurrent[i][0], bufferCurrent[i][1] + 1})
			}
			// Check to the left
			if bufferCurrent[i][0] < len(p.BoardMap) &&
				bufferCurrent[i][1]-1 < len(p.BoardMap[1]) &&
				bufferCurrent[i][0] >= 0 &&
				bufferCurrent[i][1]-1 >= 0 &&
				p.BoardMap[bufferCurrent[i][0]][bufferCurrent[i][1]-1] == 0 {

				p.BoardMap[bufferCurrent[i][0]][bufferCurrent[i][1]-1] = starter + 1
				bufferFuture = append(bufferFuture, []int{bufferCurrent[i][0], bufferCurrent[i][1] - 1})
			}
		}
		bufferCurrent = [][]int{}
		bufferCurrent = append(bufferCurrent, bufferFuture...)
		starter += 1
	}
	starter = 1
	bufferFuture = [][]int{{}}
}

//MakeBoolGrid generates BoolGrid and initializes all of its values to false
func MakeBoolGrid(p *Params) {
	p.BoolGrid = make([][]bool, p.MaxX)
	for i := range p.BoolGrid {
		p.BoolGrid[i] = make([]bool, p.MaxY)
	}
	//Initializing the boolean field with values of false
	for i := 0; i < p.MaxX; i++ {
		for j := 0; j < p.MaxY; j++ {
			p.BoolGrid[i][j] = false
		}
	}
}

//ReadMap takes the proper values for the map and writes the walls
//and nodes to the map
func ReadMap(p *Params, r *RegionParams) {

	CreateBoard(p.MaxX, p.MaxY, p)

	r.Point_list = make([]Tuple, 0)

	r.Point_list2 = make([][]bool, 0)

	r.Point_dict = make(map[Tuple]bool)

	r.Square_list = make([]RoutingSquare, 0)

	r.Border_dict = make(map[int][]int)

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

	p.Width = imgCfg.Width
	p.Height = imgCfg.Height

	fmt.Println("Width : ", p.Width)
	fmt.Println("Height : ", p.Height)

	imgfile.Seek(0, 0)

	img, _, err := image.Decode(imgfile)

	for x := 0; x < p.Width; x++ {
		r.Point_list2 = append(r.Point_list2, make([]bool, p.Height))
	}

	for x := 0; x < p.Width; x++ {
		for y := 0; y < p.Height; y++ {
			rr, _, _, _ := img.At(x, y).RGBA()
			/*rr, gg, bb, _ := img.At(1599, 90).RGBA()
			fmt.Printf("r: %d %d %d\n", rr, gg, bb)
			rr, gg, bb, _ = img.At(1599, 89).RGBA()
			fmt.Printf("r: %d %d %d\n", rr, gg, bb)*/
			if rr >= 60000 {
				r.Point_list2[x][y] = true
				r.Point_dict[Tuple{x, y}] = true

			} else {
				r.Point_dict[Tuple{x, y}] = false
				p.BoardMap[x][y] = -1
				temp := make([]int, 2)
				temp[0] = x
				temp[1] = y
				p.Wpos = append(p.Wpos, temp)
				p.BoolGrid[x][y] = true
			}
		}
	}

	CreateBoard(p.MaxX, p.MaxY, p)
	FillInWallsToBoard(p)
	FillInBufferCurrent(p)
	FillPointsToBoard(p)
	FillInMap(p)

}

//SetupFiles initilizes all of the output files to be used for the simulator
func SetupFiles(p *Params) {
	fmt.Printf("Building Output Files\n")
	dummy, err := os.Create("dummyFile.txt")
	if err != nil {
		log.Fatal("cannot create file", err)
	}
	dummy.Close()
	os.Remove("dummyFile.txt")
	p.PositionFile, err = os.Create(p.OutputFileNameCM + "-simulatorOutput.txt")
	if err != nil {
		log.Fatal("Cannot create file", err)
	}
	p.Files = append(p.Files, p.OutputFileNameCM + "-simulatorOutput.txt")

	//Print parameters to position file
	fmt.Fprintln(p.PositionFile, "Image:", p.ImageFileNameCM)
	fmt.Fprintln(p.PositionFile, "Width:", p.MaxX)
	fmt.Fprintln(p.PositionFile, "Height:", p.MaxY)
	fmt.Fprintf(p.PositionFile, "Amount: %-8v\n", p.Iterations_of_event)
	fmt.Fprintf(p.PositionFile, "Bomb x: %v\n", p.BombX)
	fmt.Fprintf(p.PositionFile, "Bomb y: %v\n", p.BombY)

	p.DriftFile, err = os.Create(p.OutputFileNameCM + "-drift.txt")
	if err != nil {
		log.Fatal("Cannot create file", err)
	}
	p.Files = append(p.Files, p.OutputFileNameCM + "-drift.txt")

	//Printing parameters to driftFile
	fmt.Fprintln(p.DriftFile, "Number of Nodes:", p.TotalNodes)
	fmt.Fprintln(p.DriftFile, "Rows:", p.SquareRowCM)
	fmt.Fprintln(p.DriftFile, "Columns:", p.SquareColCM)
	fmt.Fprintln(p.DriftFile, "Samples Stored by Node:", p.NumStoredSamples)
	fmt.Fprintln(p.DriftFile, "Samples Stored by Grid:", p.NumGridSamples)
	fmt.Fprintln(p.DriftFile, "Width:", p.MaxX)
	fmt.Fprintln(p.DriftFile, "Height:", p.MaxY)
	fmt.Fprintln(p.DriftFile, "Bomb x:", p.BombX)
	fmt.Fprintln(p.DriftFile, "Bomb y:", p.BombY)
	fmt.Fprintln(p.DriftFile, "Iterations:", p.Iterations_of_event)
	fmt.Fprintln(p.DriftFile, "Size of Square:", p.XDiv, "x", p.YDiv)
	fmt.Fprintln(p.DriftFile, "Detection Threshold:", p.DetectionThreshold)
	fmt.Fprintln(p.DriftFile, "Input File Name:", p.InputFileNameCM)
	fmt.Fprintln(p.DriftFile, "Output File Name:", p.OutputFileNameCM)
	fmt.Fprintln(p.DriftFile, "Printing Position:", p.PositionPrint, "\nPrinting Energy:", p.EnergyPrint, "\nPrinting Nodes:", p.NodesPrint)
	fmt.Fprintln(p.DriftFile, "Super Nodes:", p.NumSuperNodes, "\nSuper Node Type:", p.SuperNodeType, "\nSuper Node Speed:", p.SuperNodeSpeed, "\nSuper Node Radius:", p.SuperNodeRadius)
	fmt.Fprintln(p.DriftFile, "Error Multiplier:", p.ErrorModifierCM)
	fmt.Fprintln(p.DriftFile, "--------------------")

	p.GridFile, err = os.Create(p.OutputFileNameCM + "-grid.txt")
	if err != nil {
		log.Fatal("Cannot create file", err)
	}
	p.Files = append(p.Files, p.OutputFileNameCM + "-grid.txt")

	p.OutputLog, err = os.Create(p.OutputFileNameCM + "-OutputLog.txt")
	if err != nil {
		log.Fatal("Cannot create file", err)
	}
	p.Files = append(p.Files, p.OutputFileNameCM + "-OutputLog.txt")



	//Write parameters to gridFile
	p.GridHeight = int(math.Ceil(float64(p.MaxY)/float64(p.SquareRowCM)))
	p.GridWidth = int(math.Ceil(float64(p.MaxX)/float64(p.SquareColCM)))
	fmt.Fprintln(p.GridFile, "Width:", p.GridWidth)
	fmt.Fprintln(p.GridFile, "Height:", p.GridHeight)


	p.EnergyFile, err = os.Create(p.OutputFileNameCM + "-node.txt")
	if err != nil {
		log.Fatal("Cannot create file", err)
	}
	p.Files = append(p.Files, p.OutputFileNameCM + "-node.txt")

	p.BatteryFile, err = os.Create(p.OutputFileNameCM + "-batteryusage.txt")
	if err != nil {
		log.Fatal("Cannot create file", err)
	}
	p.Files = append(p.Files, p.OutputFileNameCM + "-batteryusage.txt")

	p.RunParamFile, err = os.Create(p.OutputFileNameCM + "-parameters.txt")
	if err != nil {
		log.Fatal("Cannot create file", err)
	}
	p.Files = append(p.Files, p.OutputFileNameCM + "-parameters.txt")

	p.RoutingFile, err = os.Create(p.OutputFileNameCM + "-path.txt")
	if err != nil {
		log.Fatal("Cannot create file", err)
	}
	p.Files = append(p.Files, p.OutputFileNameCM + "-path.txt")

	p.MoveReadingsFile, err = os.Create(p.OutputFileNameCM + "-movementReadings.txt")
	if err != nil {
		log.Fatal("Cannot create file", err)
	}
	p.Files = append(p.Files, p.OutputFileNameCM + "-movementReadings.txt")

	p.ServerFile, err = os.Create(p.OutputFileNameCM + "-server.txt")
	if err != nil {
		log.Fatal("Cannot create file", err)
	}
	p.Files = append(p.Files, p.OutputFileNameCM + "-server.txt")

	p.DetectionFile, err = os.Create(p.OutputFileNameCM + "-detection.txt")
	if err != nil {
		log.Fatal("Cannot create file", err)
	}
	p.Files = append(p.Files, p.OutputFileNameCM + "-detection.txt")

	if p.DriftExplorer || !p.DriftExplorer {
		p.DriftExploreFile, err = os.Create(p.OutputFileNameCM + "-driftExplore.txt")
		if err != nil {
			log.Fatal("Cannot create file", err)
		}
		p.Files = append(p.Files, p.OutputFileNameCM+"-driftExplore.txt")
	}

	//defer p.ServerFile.Close()
	p.NodeDataFile, err = os.Create(p.OutputFileNameCM + "-nodeData.txt")
	if err != nil {
		log.Fatal("Cannot create file", err)
	}
	p.Files = append(p.Files, p.OutputFileNameCM+"-nodeData.txt")

	p.DistanceFile, err = os.Create(p.OutputFileNameCM + "-distance.txt")
	if err != nil {
		log.Fatal("Cannot create file", err)
	}
	p.Files = append(p.Files, p.OutputFileNameCM+"-distance.txt")

	if p.ClusteringOn {
		p.ClusterStatsFile, err = os.Create(p.OutputFileNameCM + "-clusterStats.txt")
		if err != nil {
			log.Fatal("Cannot create file", err)
		}
		p.Files = append(p.Files, p.OutputFileNameCM+"-clusterStats.txt")

		p.ClusterFile, err = os.Create(p.OutputFileNameCM + "-clusters.txt")
		if err != nil {
			log.Fatal("Cannot create file", err)
		}
		p.Files = append(p.Files, p.OutputFileNameCM+"-clusters.txt")

		p.ClusterDebugFile, err = os.Create(p.OutputFileNameCM + "-clusterDebug.txt")
		if err != nil {
			log.Fatal("Cannot create file", err)
		}
		p.Files = append(p.Files, p.OutputFileNameCM+"-clusterDebug.txt")
	}

	fmt.Println(p.Files)

}

func SetupParameters(p *Params, r *RegionParams) {

	//p.XDiv = p.MaxX / p.SquareColCM
	//p.YDiv = p.MaxY / p.SquareRowCM
	p.XDiv = p.SquareColCM
	p.YDiv = p.SquareRowCM

	//The capacity for a square should be equal to the area of the square
	//So we take the side length (xDiv) and square it
	p.SquareCapacity = int(math.Pow(float64(p.XDiv), 2))

	//Center of the p.Grid
	p.Center.X = p.MaxX / 2
	p.Center.Y = p.MaxY / 2

	p.Attractions = make([]*Attraction, p.NumAtt)

	if p.CSVSensor {
		readSensorCSV(p, r)
		readFineSensorCSV(p)
	} else {

		readFineSensorCSV(p)

		/*
		p.MaxRaw = 1000
		//p.EdgeRaw = 36
		p.EdgeRaw = RawConcentration(float32(p.DetectionDistance/2))
		fmt.Println("Raw:", p.EdgeRaw)
		p.MaxADC = 4095
		p.EdgeADC = 3
		p.ADCWidth = p.MaxRaw/p.MaxADC
		p.ADCOffset = p.EdgeRaw - p.EdgeADC * p.ADCWidth
		fmt.Printf("%v %v\n", p.ADCWidth, p.ADCOffset)*/
	}

	if p.CSVMovement {
		//readMovementCSV(p)
		readInitialMovementsCSV(p)
	}


}

func InterpolateFloat (start float32, end float32, portion float32) float32{
	return (float32(end-start) * portion + float32(start))
}

func CalculateADCSetting(reading float64, x, y, time int, p *Params) {
	fmt.Println("\n", reading, time, x, y)
	total := float32(0.0)
	counted := float32(0.0)
	for i := -1; i < 2; i++ {
		for j := -1; j < 2; j++ {
			if i != 0 || j != 0 {
				if i+x > 0 && i+x < p.Width {
					if j+y > 0 && j+y < p.Height {
						part := InterpolateFloat(float32(p.SensorReadings[x][y][time]),  float32(p.SensorReadings[i + x][j + y][time]), .2/float32(Dist(Tuple{x, y}, Tuple{x+i, y+j})))
						total += part
						//fmt.Println(i, j, float32(p.SensorReadings[x][y][time]),  float32(p.SensorReadings[i + x][j + y][time]), .2/float32(Dist(Tuple{x, y}, Tuple{x+i, y+j})), part)
						counted += 1
					}
				}
			}
		}
	}

	p.MaxRaw = total/counted
	fmt.Printf("Raw: %v\n", p.MaxRaw)
	straight_increment := int(math.Ceil(p.DetectionDistance))
	//fmt.Println(straight_increment)
	diag_increment := int(math.Ceil(p.DetectionDistance/math.Sqrt2))
	//fmt.Println(diag_increment)

	points := [8][2]int{
		{straight_increment, 0},			//right
		{diag_increment, diag_increment}, 	//upper right
		{0, straight_increment},			//up
		{-diag_increment, diag_increment}, 	//upper left
		{-straight_increment, 0},			//left
		{-diag_increment, -diag_increment}, //lower left
		{0, -straight_increment},			//down
		{diag_increment, -diag_increment}}	//lower right

	total = 0
	counted = 0
	for _, point := range (points) {
		if point[0]+x > 0 && point[0]+x < p.Width {
			if point[1]+y > 0 && point[1]+y < p.Height {
				//if legal then interpolate the correct value based on distance
				part := InterpolateFloat(float32(p.SensorReadings[x][y][time]), float32(p.SensorReadings[x+point[0]][y+point[1]][time]), float32(p.DetectionDistance/(Dist(Tuple{x, y}, Tuple{x + point[0], y + point[1]}))))
				total += part
				//fmt.Println(float32(p.SensorReadings[x][y][time]), float32(p.SensorReadings[x + point[0]][y + point[1]][time]), float32(p.DetectionDistance/(Dist(Tuple{x, y}, Tuple{x + point[0], y + point[1]}))), part)
				counted += 1
			}
		}
	}

	fmt.Println(p.MaxRaw)
	//p.EdgeRaw = 36
	p.EdgeRaw = total/counted
	fmt.Println(p.EdgeRaw)
	p.MaxADC = 4095
	p.EdgeADC = 3
	p.ADCWidth = p.MaxRaw/p.MaxADC
	p.ADCOffset = p.EdgeRaw - p.EdgeADC * p.ADCWidth

	fmt.Println("Edge Raw:", p.EdgeRaw)
	fmt.Printf("%v %v\n", p.ADCWidth, p.ADCOffset)

	fmt.Printf("\n")
}

func CalculateFineADCSetting(reading float64, x, y float32, time int, p *Params) {
	straight_increment := float32(.2)  //this is actually .1m away from the bomb by nature of the computation that occurs in the interpolation
	diag_increment := float32(math.Sqrt(math.Pow(float64(straight_increment), 2.0)*2.0))


	total := float32(0.0)

	closePoints := [8][2]float32{
		{straight_increment, 0},			//right
		{diag_increment, diag_increment}, 	//upper right
		{0, straight_increment},			//up
		{-diag_increment, diag_increment}, 	//upper left
		{-straight_increment, 0},			//left
		{-diag_increment, -diag_increment}, //lower left
		{0, -straight_increment},			//down
		{diag_increment, -diag_increment}}	//lower right

	for i := 0; i < 8; i++ {
		//fmt.Printf("%v %v\n", x+closePoints[i][0], y+closePoints[i][1])
		total += float32(interpolateReading(x + closePoints[i][0], y + closePoints[i][1], time*1000, time, true, p))
	}

	p.MaxRaw = total/8.0


	straight_increment = float32((p.DetectionDistance))
	//fmt.Printf("straight: %v\n", straight_increment)
	diag_increment = float32(straight_increment/math.Sqrt2)
	//fmt.Printf("diag: %v\n", diag_increment)




	points := [8][2]float32{
		{straight_increment, 0},			//right
		{diag_increment, diag_increment}, 	//upper right
		{0, straight_increment},			//up
		{-diag_increment, diag_increment}, 	//upper left
		{-straight_increment, 0},			//left
		{-diag_increment, -diag_increment}, //lower left
		{0, -straight_increment},			//down
		{diag_increment, -diag_increment}}	//lower right

	total = 0.0
	for i := 0; i < 8; i++ {
		total += float32(interpolateReading(x + points[i][0], y + points[i][1], time*1000, time, true, p))
	}


	p.EdgeRaw = total/8.0
	//fmt.Println(p.EdgeRaw)
	p.MaxADC = 4095
	p.EdgeADC = 3
	p.MaxRaw = 12285 //////CHANGE
	if p.DriftExplorer {
		p.ADCWidth = (p.MaxRaw) / p.MaxADC
	} else {
		p.ADCWidth = (p.MaxRaw) / p.MaxADC
	}
	p.ADCOffset = p.EdgeRaw - p.EdgeADC * p.ADCWidth

	//fmt.Printf("\n")

	fmt.Println("Max Raw:", p.MaxRaw)
	fmt.Println("Edge Raw:", p.EdgeRaw)
	fmt.Printf("%v %v\n", p.ADCWidth, p.ADCOffset)

}




//Helper to count lines and learn progress
func lineCounter(fileName string) (int, error) {

	r, _ := os.Open(fileName)
	buf := make([]byte, 32*1024)
	count := 0
	lineSep := []byte{'\n'}

	for {
		c, err := r.Read(buf)
		count += bytes.Count(buf[:c], lineSep)

		switch {
		case err == io.EOF:
			return count, nil

		case err != nil:
			return count, err
		}
	}
}

func ReadWindRegion(p *Params) {
	in, err := os.Open(p.WindRegionPath)
	if err != nil {
		println("error opening file", err)
	}
	defer in.Close()

	record, err := ioutil.ReadAll(in)
	line := strings.Split(string(record), "\r\n")

	data := make([][]string, 0)
	for i := range line {
		data = append(data, strings.Split(line[i], ","))
	}
	//fmt.Println(data[1])
	count := 0
	x := int64(0)
	y := int64(0)
	p.WindRegion = make([][]Coord, len(data))
	for t := range data {
		for i:= 0; i < len(data[t]); i++ {
			if i % 2 == 0 {
				x,_ = strconv.ParseInt(data[t][i], 10, 32)
				count ++
			} else {
				y,_ = strconv.ParseInt(data[t][i], 10, 32)
				count ++
			}
			if count == 2 {
				p.WindRegion[t] = append(p.WindRegion[t], Coord{X: int(x), Y: int(y)})
				count = 0
			}
		}
	}

	//fmt.Println(len(p.WindRegion))
	//fmt.Println(p.WindRegion[1])
}

//readSensorCSV reads the sensor values from a file
func readSensorCSV(p *Params, region *RegionParams) {

	in, err := os.Open(p.SensorPath)
	if err != nil {
		println("error opening file")
	}
	defer in.Close()
	fmt.Printf("Reading Sensor Files\n")

	lines, _ := lineCounter(p.SensorPath)

	//fmt.Println(lines)
	r := csv.NewReader(in)
	r.ReuseRecord = true

	r.FieldsPerRecord = -1

	record, err := r.Read()
	//record, err := r.ReadAll()

	reg, _ := regexp.Compile("([0-9]+)")
	times := reg.FindAllStringSubmatch(strings.Join(record, " "), -1)

	p.SensorTimes = make([]int, len(times))
	for i := range times {
		p.SensorTimes[i], _ = strconv.Atoi(times[i][1])
	}
	p.MaxTimeStep = len(times)

	numSamples := len(record) - 3

	//record, err = r.Read();
	//record, err = r.Read();

	//p.FineScale, _ = strconv.Atoi(record[0])

	fmt.Println(numSamples)

	p.SensorReadings = make([][][]float64, p.Width)
	for i := range p.SensorReadings {
		p.SensorReadings[i] = make([][]float64, p.Height)
		/*for j := range p.SensorReadings[i] {

		}*/
	}

	i := 1
	fmt.Printf("Sensor CSV Processing\n")

	for {
		record, err = r.Read()
		if err != nil {
			if err == io.EOF {
				break
			}
		}
		p.FineScale, _ = strconv.Atoi(record[0])
		x, _ := strconv.ParseInt(record[1], 10, 32);
		y, _ := strconv.ParseInt(record[2], 10, 32);

		if(SoftRegionContaining(Tuple{int(x),int(y)}, region) >= 0) {

			p.SensorReadings[x][y] = make([]float64, numSamples)
			for k := range p.SensorReadings[x][y] {
				p.SensorReadings[x][y][k] = 0
			}

			j := 3
			//fmt.Println((len(record)))
			for j < len(record) {
				read1, _ := strconv.ParseFloat(record[j], 32);
				if read1 < 0 {
					read1 = 0
				}
				p.SensorReadings[x][y][j-3] = read1

				j += 1

			}



		}
		i++
		if (i%1000 == 0) {
			prog := int(float32(i) / float32(lines) * 100)
			fmt.Printf("\rProgress [%s%s] %d ", strings.Repeat("=", prog), strings.Repeat(".", 100-prog), prog)
		}
	}

	//CalculateADCSetting(maxReading, maxLocX, maxLocY, maxTime, p)
	//fmt.Println(p.BombX, p.BombY)
	//CalculateADCSetting(p.SensorReadings[p.B.X][p.B.Y][10], p.B.X, p.B.Y, 10, p)
	fmt.Println("")
}



//readSensorCSV reads the sensor values from a file
func readFineSensorCSV(p *Params) {

	in, err := os.Open(p.FineSensorPath)
	if err != nil {
		println("error opening file")
	}
	defer in.Close()

	fmt.Printf("Reading Fine Sensor File\n")

	lines, _ := lineCounter(p.FineSensorPath)

	r := csv.NewReader(in)
	r.FieldsPerRecord = -1
	record, err := r.Read()

	reg, _ := regexp.Compile("([0-9]+)")
	times := reg.FindAllStringSubmatch(strings.Join(record, " "), -1)

	p.SensorTimes = make([]int, len(times))
	for i := range times {
		p.SensorTimes[i], _ = strconv.Atoi(times[i][1])
	}
	p.MaxTimeStep = len(times)

	numSamples := len(record) - 3


	record, _ = r.Read()

	p.FineScale, _ = strconv.Atoi(record[0])
	p.FineWidth = int(math.Sqrt(float64(lines-1)))
	p.FineHeight = int(math.Sqrt(float64(lines-1)))

	fmt.Printf("%v %v\n", p.FineWidth, p.FineHeight)


	p.FineSensorReadings = make([][][]float64, p.FineWidth)
	for i := range p.FineSensorReadings {
		p.FineSensorReadings[i] = make([][]float64, p.FineHeight)
		for j := range p.FineSensorReadings[i] {
			p.FineSensorReadings[i][j] = make([]float64, numSamples)
			for k := range p.FineSensorReadings[i][j] {
				p.FineSensorReadings[i][j][k] = 0
			}
		}
	}

	//in.Seek(0,0)

	//r.Read()

	i := 1
	fmt.Printf("Fine Sensor CSV Processing\n")


	for  {
		record, err = r.Read()
		if err != nil {
			if err == io.EOF {
				break
			}
		}

		x, _ := strconv.ParseInt(record[1], 10, 32);
		y, _ := strconv.ParseInt(record[2], 10, 32);
		//fmt.Printf("%d %d\n", x, y)

		j := 3

		for j < len(record) {
			read1, _ := strconv.ParseFloat(record[j], 32);
			//fmt.Printf("%d %d %v\n", x, y, read1)
			if read1 < 0 {
				read1 = 0
			}
			p.FineSensorReadings[x][y][j-3] = read1
			/*if(x == 25 && y == 25) {
				fmt.Printf("%d %d %v\n", x, y, p.FineSensorReadings[x][y][j-3])
			}*/
			j += 1
		}
		i++

		if(i % 100 == 0) {
			prog := int(float32(i)/float32(lines)*100)
			fmt.Printf("\rProgress [%s%s] %d ", strings.Repeat("=", prog), strings.Repeat(".", 100-prog), prog)
		}
	}

	//CalculateADCSetting(maxReading, maxLocX, maxLocY, maxTime, p)
	//fmt.Println(p.BombX, p.BombY)
	fmt.Println()
	//CalculateFineADCSetting(p.FineSensorReadings[p.FineWidth/2][p.FineHeight/2][0], p.FineWidth/2, p.FineHeight/2, 400, p)
	CalculateFineADCSetting(p.FineSensorReadings[p.FineWidth/2][p.FineHeight/2][0], float32(p.B.X), float32(p.B.Y), 120, p)
}

func PartialReadMovementCSV(p *Params) {
	fmt.Println(p.MovementPath)
	in, err := os.Open(p.MovementPath)
	if err != nil {
		println("error opening file")
	}
	defer in.Close()

	r := csv.NewReader(in)
	r.FieldsPerRecord = -1
	//record, err := r.Read()
	r.ReuseRecord = true

	timeSteps, _ := lineCounter(p.MovementPath)


	p.MovementOffset = p.MovementSize + p.MovementOffset - 2 //go back one

	record, err := r.Read()
	i := 0
	for i < p.MovementOffset {  //-1 for initial read
		record, err = r.Read()
		i += 1
	}




	time := p.MovementOffset
	fmt.Printf("offset: %v\n", time)
	fmt.Printf("Movement CSV Processing %d TimeSteps for %d Nodes  %d\n", timeSteps, len(record), p.TotalNodes)
	for time < timeSteps && time < (p.MovementOffset + p.MovementSize) {
		iter := 0

		//fmt.Printf("%v\n", time)
		for iter < len(record)-1 && iter/2 < p.TotalNodes {

			x, _ := strconv.ParseInt(record[iter], 10, 32);
			y, _ := strconv.ParseInt(record[iter+1], 10, 32);

			p.NodeMovements[iter/2][time-p.MovementOffset] = Tuple{int(x), int(y)}

			iter += 2
		}
		time++

		if(time % 10 == 0) {
			prog := int(float32(time)/float32(timeSteps)*100)
			fmt.Printf("\rProgress [%s%s] %d ", strings.Repeat("=", prog), strings.Repeat(".", 100-prog), prog)
		}

		record, err = r.Read()
		if err != nil {
			if err == io.EOF {
				break
			} else{
				break
			}
		}
	}

	fmt.Printf("\n")
}



//readMovementCSV reads the movement parameters from a file
func readInitialMovementsCSV(p *Params) {
	fmt.Println(p.MovementPath)
	in, err := os.Open(p.MovementPath)
	if err != nil {
		println("error opening file")
	}
	defer in.Close()

	r := csv.NewReader(in)
	r.FieldsPerRecord = -1
	//record, err := r.Read()
	r.ReuseRecord = true


	timeSteps, _ := lineCounter(p.MovementPath)

	p.NumNodeMovements = timeSteps

	p.NodeMovements = make([][]Tuple, p.TotalNodes)
	for i := range p.NodeMovements {
		p.NodeMovements[i] = make([]Tuple, p.MovementSize)
	}

	record, err := r.Read()

	p.MovementOffset = 0
	time := 0
	fmt.Printf("Movement CSV Processing %d TimeSteps for %d Nodes  %d\n", timeSteps, len(record), p.TotalNodes)
	for time < timeSteps && time < (p.MovementOffset + p.MovementSize) {
		iter := 0

		//fmt.Printf("in %v\n", time)
		for iter < len(record)-1 && iter/2 < p.TotalNodes {

			x, _ := strconv.ParseInt(record[iter], 10, 32);
			y, _ := strconv.ParseInt(record[iter+1], 10, 32);

			p.NodeMovements[iter/2][time] = Tuple{int(x), int(y)}

			iter += 2
		}
		time++

		if(time % 10 == 0) {
			prog := int(float32(time)/float32(timeSteps)*100)
			fmt.Printf("\rProgress [%s%s] %d ", strings.Repeat("=", prog), strings.Repeat(".", 100-prog), prog)
		}

		record, err = r.Read()
		if err != nil {
			if err == io.EOF {
				break
			} else{
				break
			}
		}
	}

	fmt.Printf("\n")
}


func RangeInt(min, max int) int { //returns a random number between max and min
	return rand.Intn(max-min) + min
}


func FlipSquares(p *Params, r *RegionParams) {
	//tmp := 0
	for i:= range(r.Square_list) {
		//tmp = r.Square_list[i].Y1
		r.Square_list[i].Y1 = p.Height - r.Square_list[i].Y1 - 1
		r.Square_list[i].Y2 = p.Height - r.Square_list[i].Y2 - 1
		for j:= range r.Square_list[i].Routers {
			r.Square_list[i].Routers[j].Y = p.Height - r.Square_list[i].Routers[j].Y - 1
		}
	}
}

func GetFlags(p *Params) {
	//p = cps.Params{}

	flag.StringVar(&p.CPUProfile, "cpuprofile", "", "write cpu profile to `file`")

	flag.StringVar(&p.MemProfile, "memprofile", "", "write memory profile to `file`")

	//fmt.Println(os.Args[1:], "\nhmmm? \n ") //C:\Users\Nick\Desktop\comand line experiments\src

	flag.IntVar(&p.BombXCM, "bombX", 0, "X location of bomb")
	flag.IntVar(&p.BombYCM, "bombY", 0, "Y location of bomb")
	flag.BoolVar(&p.CommBomb, "commandBomb", false, "Whether to use command line for bomb coords")

	flag.IntVar(&p.SittingStopThresholdCM, "sittingStopThreshold", 5,
		"How long it takes for a node to stay seated")

	flag.StringVar(&p.InputFileNameCM, "inputFileName", "Log1_in.txt",
		"Name of the input text file")

	flag.StringVar(&p.ValidationType, "validationType", "square", "Type of validation: square or validators")

	flag.BoolVar(&p.RecalReject, "recalReject", false, "reject if hasn't been recalibrated recently")

	flag.StringVar(&p.SensorPath, "sensorPath", "Circle_2D.csv", "Sensor Reading Inputs")

	flag.StringVar(&p.FineSensorPath, "fineSensorPath", "Circle_2D.csv", "Sensor Reading Inputs")

	flag.StringVar(&p.MovementPath, "movementPath", "Circle_2D.csv", "Movement Inputs")

	flag.StringVar(&p.OutputFileNameCM, "OutputFileName", "Log",
		"Name of the output text file prefix")

	//flag.IntVar(&p.CMSensingTime, "cmSensingTime,",2, "seconds a cluster member will sense/record readings before sending to cluster head")
	//flag.IntVar(&p.CHSensingTime, "chSensingTime,",4, "seconds a cluster head will sense//collect from CM/record readings before sending to server")
	//flag.IntVar(&p.MaxCMReadingBufferSize, "maxCMReadingBufferSize,",10, "max readings buffer size of a cluster member. CM must send to CH when buffer is this size")
	//flag.IntVar(&p.MaxCHReadingBufferSize, "maxCHReadingBufferSize,",100, "max readings buffer size of a cluster head. CH must send to server when buffer is this size")
	flag.IntVar(&p.MaxBufferCapacityCM, "maxBufferCapacity", 25,
		"maximum capacity for the buffer before it sends data to the server")
	flag.IntVar(&p.NumStoredSamplesCM, "nodeStoredSamples", 10,
		"number of samples stored by individual nodes for averaging")

	flag.IntVar(&p.GridStoredSamplesCM, "GridStoredSamples", 10,
		"number of samples stored by p.Grid squares for averaging")

	flag.Float64Var(&p.DetectionThresholdCM, "detectionThreshold", 4000.0, //11180.0,
		"Value where if a node gets this reading or higher, it will trigger a detection")

	flag.Float64Var(&p.ErrorModifierCM, "errorMultiplier", 1.0,
		"Multiplier for error values in system")

	flag.BoolVar(&p.CSVSensor, "csvSensor", true, "Read Sensor Values from CSV")

	flag.BoolVar(&p.CSVMovement, "csvMove", true, "Read Movements from CSV")

	flag.BoolVar(&p.SuperNodes, "superNodes", true, "Enable SuperNodes")

	flag.IntVar(&p.IterationsCM, "iterations", 200, "Read Movements from CSV")

	//Range 1, 2, or 4
	//Currently works for only a few numbers, can be easily expanded but is not currently dynamic
	flag.IntVar(&p.NumSuperNodes, "numSuperNodes", 4, "the number of super nodes in the simulator")

	flag.Float64Var(&p.CalibrationThresholdCM, "RecalibrationThreshold", 3.0, "Value over grid average to recalibrate node")

	flag.Float64Var(&p.StdDevThresholdCM, "StandardDeviationThreshold", 1.7, "Detection Threshold based on standard deviations from mean")

	flag.Float64Var(&p.DetectionDistance, "detectionDistance", 6.0, "Detection Distance")

	flag.IntVar(&p.CounterThreshold,"CounterThreshold",3,"Threshold to decrease sampling rate")

	flag.Float64Var(&p.MaxMoveMeters,"MaxMoveMeters",2,"maxMoveMeters")



	//Range: 0-2
	//0: default routing algorithm, points added onto the end of the path and routed to in that order
	//flag.IntVar(&p.SuperNodeType, "p.SuperNodeType", 0, "the type of super node used in the simulator")
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
	flag.IntVar(&p.SuperNodeType, "SuperNodeType", 0, "the type of super node used in the simulator")

	//Range: 0-...
	//Theoretically could be as high as possible
	//Realistically should remain around 10
	flag.IntVar(&p.SuperNodeSpeed, "SuperNodeSpeed", 3, "the speed of the super node")


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
	//flag.IntVar(&p.SuperNodeVariation, "p.SuperNodeVariation", 3, "super nodes of type 1 have different variations")

	flag.BoolVar(&p.OutputPrint, "logOutput", true, "Whether you want to write OutputLog file")

	flag.BoolVar(&p.PositionPrint, "logPosition", false, "Whether you want to write position info to a log file")

	flag.BoolVar(&p.GridPrint, "logGrid", false, "Whether you want to write p.Grid info to a log file")

	flag.BoolVar(&p.EnergyPrint, "logEnergy", false, "Whether you want to write energy into to a log file")

	flag.BoolVar(&p.BatteryPrint, "logBattery", false, "Whether you want to write battery info to a log file")

	flag.BoolVar(&p.NodesPrint, "logNodes", false, "Whether you want to write node readings to a log file")

	flag.BoolVar(&p.ClusterPrint, "logClusters", false, "Whether you want to write cluster statistics to a log file")

	flag.BoolVar(&p.ClusterDebug, "clusterDebug", false, "Whether you want to write cluster debug information to log files")

	flag.BoolVar(&p.ReportBTAverages, "reportBTAverages", false, "Whether you want to write avg number of nodes in bluetooth range to cluster log file")

	flag.IntVar(&p.SquareRowCM, "SquareRowCM", 50, "Number of rows of p.Grid squares, 1 through p.MaxX")

	flag.IntVar(&p.SquareColCM, "SquareColCM", 50, "Number of columns of p.Grid squares, 1 through p.MaxY")


	flag.StringVar(&p.ImageFileNameCM, "imageFileName", "circle_justWalls_x4.png", "Name of the input text file")

	flag.StringVar(&p.StimFileNameCM, "stimFileName", "circle_0.txt", "Name of the stimulus text file")


	flag.StringVar(&p.OutRoutingNameCM, "outRoutingName", "log.txt", "Name of the stimulus text file")

	flag.StringVar(&p.OutRoutingStatsNameCM, "outRoutingStatsName", "routingStats.txt", "Name of the output file for stats")


	flag.BoolVar(&p.RegionRouting, "regionRouting", true, "True if you want to use the new routing algorithm with regions and cutting")

	flag.BoolVar(&p.ClusteringOn,"clusteringOn",true,"True: nodes will form clusters, False: no clusters will form")
	flag.BoolVar(&p.RedundantClustering,"redundantClustering",false,"If clusteringOn is set to true, True: nodes will join two clusters, False: clusters will form normally")
	flag.IntVar(&p.ClusterMaxThreshold, "clusterMaxThresh",8, "max number of members in a node cluster")
	flag.IntVar(&p.ClusterMinThreshold, "clusterMinThresh", 0, "max number of members in a node cluster for it to be considered 'empty'")
	flag.Float64Var(&p.NodeBTRange, "nodeBTRange",20.0,"bluetooth range of each node")
	flag.Float64Var(&p.DegreeWeight, "degreeWeight", 0.6, "The weight constant applied to the number of neighboring nodes when calculating a node's score")
	flag.Float64Var(&p.BatteryWeight, "batteryWeight", 0.4, "The weight constant applied to a node's battery when calculating a node's score")
	flag.Float64Var(&p.Penalty, "penalty", 0.8, "The penalty multiplied to a node's score when it is not already a cluster head")
	/* Global Reclustering
	0 - off
	1 - threshold-based
	2 - time-based
	*/
	flag.IntVar(&p.GlobalRecluster, "globalRecluster", 1, "Enables or disables global reclustering")
	/* Local Reclustering
	0 - off
	1 - minimal (nodes check for nearby head first)
	2 - standard
	3 or higher - expansive (nearby clusters also recluster)
	When expansive is chosen, (p.LocalRecluster - 3) will be the time in seconds that a node has been a cluster head for its cluster to join the expansive recluster
	For example, if p.LocalRecluster is set to 63, only nodes that have been cluster head for over 60 seconds will join the expansive recluster
	*/
	flag.IntVar(&p.LocalRecluster, "localRecluster", 1, "Enables or disables local reclustering")
	flag.Float64Var(&p.ReclusterThreshold, "reclusterThreshold", 0.1, "The maximum percent of clusters made up only of cluster heads before the network should fully recluster")
	flag.IntVar(&p.ReclusterPeriod, "reclusterPeriod", 30, "The period of time in seconds before the network checks if it should fully reclusters")
	flag.Float64Var(&p.SmallImprovement, "smallImprovement", 0.2, "The threshold improvement in alone nodes after a global recluster that will lead to an increased recluster threshold or period")
	flag.Float64Var(&p.LargeImprovement, "largeImprovement", 0.6, "The threshold improvement in alone nodes after a global recluster that will lead to a decreased recluster threshold or period")
	flag.Float64Var(&p.GlobalReclusterIncrement, "GRIncrement", 1.3, "The number that the recluster threshold or period will be multiplied by when it needs to be increased.")
	flag.Float64Var(&p.GlobalReclusterDecrement, "GRDecrement", 0.8, "The number that the recluster threshold or period will be multiplied by when it needs to be decreased.")
	flag.Float64Var(&p.ServerReadyThreshold, "ServerReadyThresh", 0.98, "The ratio of accounted nodes to alive nodes at which point a recluster is possible and at which point after a recluster the improvement will be calculated.")
	flag.IntVar(&p.InitClusterTime, "initClusterTime", 0, "The number of seconds to wait before clustering")
	flag.IntVar(&p.ClusterSearchThreshold, "clusterSearchThresh", 0, "The number of senses in a row required to trigger cluster search that a node either has no cluster head or is out of range of its cluster head.")
	flag.IntVar(&p.ClusterHeadTimeThreshold, "CHTimeThresh", 300, "The maximum time a can be cluster head without triggering local recluster.")
	flag.Float64Var(&p.ClusterHeadBatteryDropThreshold, "CHBatteryDropThresh", 0.3, "The maximum percent a cluster head's battery can drop before triggering a local recluster.")

	flag.StringVar(&p.WindRegionPath, "windRegionPath", "hull_testing.txt", "File containing regions formed by wind")

	flag.BoolVar(&p.DriftExplorer, "driftExplorer", false, "True if you want to JUST examine sensor drifting")

	flag.BoolVar(&p.ServerRecal, "serverRecal", true, "True if you want the server to be able to recalibrate nodes")

	flag.IntVar(&p.ReadingHistorySize, "detectionWindow", 60, "Window in which detections are kept")
	flag.IntVar(&p.ValidationThreshold, "validationThreshold", 1, "Number of detections required to validate a detection")
	flag.IntVar(&p.TotalNodes, "totalNodes", -1, "Number of Nodes")
	flag.IntVar(&p.MovementSize, "moveSize", 1000, "Number of movement records to read each load")

	flag.BoolVar(&p.RandomBomb, "randomBomb", false, "Toggles random bomb placement")
	flag.BoolVar(&p.ZipFiles, "zipFiles", false, "Toggles Zipping of output files")
	flag.IntVar(&p.DensityThreshold, "densityThreshold", 10, "Number of nodes to make a square considered dense")
	flag.IntVar(&p.SamplingPeriodMS,"SamplingPeriodMS",500,"period at which nodes sense")
	// New Battery Level Flags
	flag.IntVar(&p.BatteryCapacity, "batteryCapacity", 10000, "Max battery capacity of all nodes")
	flag.Float64Var(&p.BatteryDeadThreshold, "batteryDeadThreshold", .10, "minimum battery percentage before a node is considered dead")
	flag.Float64Var(&p.BatteryLowThreshold, "batteryLowThreshold", .15, "battery percentage to mark a node as low power")
	flag.Float64Var(&p.BatteryMediumThreshold, "batteryMediumThreshold", .25, "battery percentage to mark a node as medium power")
	flag.Float64Var(&p.BatteryHighThreshold, "batteryHighThreshold", .40, "battery percentage to mark a node as high power")
	flag.Float64Var(&p.AverageBatteryLevel, "averageBatteryLevel", 0.75, "average initial battery level to set nodes to")
	flag.Float64Var(&p.BluetoothLossPercentage, "bluetoothLossPercentage", 0.00002, "amount of battery drained each time a node uses bluetooth")
	flag.Float64Var(&p.SampleLossPercentage, "sampleLossPercentage", 0.0002, "amount of battery drained each time a node takes a sample")
	flag.Float64Var(&p.WifiLossPercentage, "wifiLossPercentage", 0.0002, "amount of battery drained each time a node uses wifi")

	flag.Parse()
	fmt.Println("Maximum size of buffer posible: ", p.MaxBufferCapacityCM)
	fmt.Println("Number of Node Stored Samples:", p.NumStoredSamplesCM)
	fmt.Println("Number of Grid Stored Samples:", p.GridStoredSamplesCM)
	fmt.Println("Detection Threshold:", p.DetectionThresholdCM)

	//fmt.Println("tail:", flag.Args())
}

func WriteFlags(p * Params){

	var buf bytes.Buffer
	buf.WriteString(fmt.Sprintf("cpuprofile=%v\n", p.CPUProfile))
	buf.WriteString(fmt.Sprintf("memprofile=%v\n", p.MemProfile))
	buf.WriteString(fmt.Sprintf("sittingStopThreshold=%v\n", p.SittingStopThresholdCM))
	buf.WriteString(fmt.Sprintf("inputFileName=%v\n", p.InputFileNameCM))
	buf.WriteString(fmt.Sprintf("sensorPath=%v\n", p.SensorPath))
	buf.WriteString(fmt.Sprintf("fineSensorPath=%v\n", p.FineSensorPath))
	buf.WriteString(fmt.Sprintf("movementPath=%v\n", p.MovementPath))
	buf.WriteString(fmt.Sprintf("OutputFileName=%v\n", p.OutputFileNameCM))
	buf.WriteString(fmt.Sprintf("maxBufferCapacity=%v\n", p.MaxBufferCapacityCM))
	buf.WriteString(fmt.Sprintf("nodeStoredSamples=%v\n", p.NumStoredSamplesCM))
	buf.WriteString(fmt.Sprintf("GridStoredSamples=%v\n", p.GridStoredSamplesCM))
	buf.WriteString(fmt.Sprintf("detectionThreshold=%v\n", p.DetectionThresholdCM))
	buf.WriteString(fmt.Sprintf("errorMultiplier=%v\n", p.ErrorModifierCM))
	buf.WriteString(fmt.Sprintf("csvSensor=%v\n", p.CSVSensor))
	buf.WriteString(fmt.Sprintf("csvMove=%v\n", p.CSVMovement))
	buf.WriteString(fmt.Sprintf("superNodes=%v\n", p.SuperNodes))
	buf.WriteString(fmt.Sprintf("iterations=%v\n", p.IterationsCM))
	buf.WriteString(fmt.Sprintf("numSuperNodes=%v\n", p.NumSuperNodes))
	buf.WriteString(fmt.Sprintf("Recalibration Threshold=%v\n", p.CalibrationThresholdCM))
	buf.WriteString(fmt.Sprintf("StandardDeviationThreshold=%v\n", p.StdDevThresholdCM))
	buf.WriteString(fmt.Sprintf("detectionDistance=%v\n", p.DetectionDistance))
	buf.WriteString(fmt.Sprintf("inputFileName=%v\n", p.InputFileNameCM))
	buf.WriteString(fmt.Sprintf("SuperNodeSpeed=%v\n", p.SuperNodeSpeed))
	buf.WriteString(fmt.Sprintf("doOptimize=%v\n", p.DoOptimize))
	buf.WriteString(fmt.Sprintf("logOutput=%v\n", p.OutputPrint))
	buf.WriteString(fmt.Sprintf("logPosition=%v\n", p.PositionPrint))
	buf.WriteString(fmt.Sprintf("logGrid=%v\n", p.GridPrint))
	buf.WriteString(fmt.Sprintf("logEnergy=%v\n", p.EnergyPrint))
	buf.WriteString(fmt.Sprintf("logBattery=%v\n", p.BatteryPrint))
	buf.WriteString(fmt.Sprintf("logNodes=%v\n", p.NodesPrint))
	buf.WriteString(fmt.Sprintf("logClusters=%v\n", p.ClusterPrint))
	buf.WriteString(fmt.Sprintf("SquareRowCM=%v\n", p.SquareRowCM))
	buf.WriteString(fmt.Sprintf("SquareColCM=%v\n", p.SquareColCM))
	buf.WriteString(fmt.Sprintf("imageFileName=%v\n", p.ImageFileNameCM))
	buf.WriteString(fmt.Sprintf("stimFileName=%v\n", p.StimFileNameCM))
	buf.WriteString(fmt.Sprintf("outRoutingStatsName=%v\n", p.OutRoutingStatsNameCM))
	buf.WriteString(fmt.Sprintf("regionRouting=%v\n", p.RegionRouting))
	buf.WriteString(fmt.Sprintf("validationThreshold=%v\n", p.ValidationThreshold))
	buf.WriteString(fmt.Sprintf("bombX=%v\n", p.B.X))
	buf.WriteString(fmt.Sprintf("bombY=%v\n", p.B.Y))
	buf.WriteString(fmt.Sprintf("serverRecal=%v\n", p.ServerRecal))
	buf.WriteString(fmt.Sprintf("driftExplorer=%v\n", p.DriftExplorer))
	buf.WriteString(fmt.Sprintf("detectionWindow=%v\n", p.ReadingHistorySize))
	buf.WriteString(fmt.Sprintf("totalNodes=%v\n", p.TotalNodes))
	buf.WriteString(fmt.Sprintf("validaitonType=%v\n", p.ValidationType))
	buf.WriteString(fmt.Sprintf("recalReject=%v\n", p.RecalReject))
	buf.WriteString(fmt.Sprintf("clusterMaxThresh=%v\n", p.ClusterMaxThreshold))
	buf.WriteString(fmt.Sprintf("clusterMinThresh=%v\n", p.ClusterMinThreshold))
	buf.WriteString(fmt.Sprintf("nodeBTRange=%v\n", p.NodeBTRange))
	buf.WriteString(fmt.Sprintf("clusteringOn=%v\n",p.ClusteringOn))
	buf.WriteString(fmt.Sprintf("degreeWeight=%v\n",p.DegreeWeight))
	buf.WriteString(fmt.Sprintf("batteryWeight=%v\n",p.BatteryWeight))
	buf.WriteString(fmt.Sprintf("penalty=%v\n",p.Penalty))
	buf.WriteString(fmt.Sprintf("globalRecluster=%v\n",p.GlobalRecluster))
	buf.WriteString(fmt.Sprintf("localRecluster=%v\n",p.LocalRecluster))
	buf.WriteString(fmt.Sprintf("reclusterPeriod=%v\n",p.ReclusterPeriod))
	buf.WriteString(fmt.Sprintf("reclusterThreshold=%v\n",p.ReclusterThreshold))
	buf.WriteString(fmt.Sprintf("initClusterTime=%v\n",p.InitClusterTime))
	buf.WriteString(fmt.Sprintf("clusterSearchThresh=%v\n",p.ClusterSearchThreshold))
	buf.WriteString(fmt.Sprintf("CHTimeThresh=%v\n", p.ClusterHeadTimeThreshold))
	buf.WriteString(fmt.Sprintf("CHBatteryDropThresh=%v\n", p.ClusterHeadBatteryDropThreshold))
	buf.WriteString(fmt.Sprintf("batteryCapacity=%v\n",p.BatteryCapacity))
	buf.WriteString(fmt.Sprintf("bluetoothLossPercentage=%v\n",p.BluetoothLossPercentage))
	buf.WriteString(fmt.Sprintf("sampleLossPercentage=%v\n",p.SampleLossPercentage))
	buf.WriteString(fmt.Sprintf("wifiLossPercentage=%v\n",p.WifiLossPercentage))
	//buf.WriteString(fmt.Sprintf("cmSensingTime=%v\n",p.CMSensingTime))
	//buf.WriteString(fmt.Sprintf("chSensingTime=%v\n",p.CHSensingTime))
	//buf.WriteString(fmt.Sprintf("maxCMReadingBufferSize=%v\n",p.MaxCMReadingBufferSize))
	//buf.WriteString(fmt.Sprintf("maxCHReadingBufferSize=%v\n",p.MaxCHReadingBufferSize))
	fmt.Fprintf(p.RunParamFile,buf.String())
}




func ZipFiles(filename string, files []string) error {

	newZipFile, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer newZipFile.Close()

	zipWriter := zip.NewWriter(newZipFile)
	defer zipWriter.Close()

	// Add files to zip
	for _, file := range files {
		if err = AddFileToZip(zipWriter, file); err != nil {
			return err
		}
	}
	return nil
}

func AddFileToZip(zipWriter *zip.Writer, filename string) error {

	fileToZip, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer fileToZip.Close()

	// Get the file information
	info, err := fileToZip.Stat()
	if err != nil {
		return err
	}
	header, err := zip.FileInfoHeader(info)
	if err != nil {
		return err
	}
	// Using FileInfoHeader() above only uses the basename of the file. If we want
	// to preserve the folder structure we can overwrite this with the full path.
	//header.Name = //filename

	// Change to deflate to gain better compression
	// see http://golang.org/pkg/archive/zip/#pkg-constants
	header.Method = zip.Deflate

	writer, err := zipWriter.CreateHeader(header)
	if err != nil {
		return err
	}
	_, err = io.Copy(writer, fileToZip)
	return err
}

func DriftHist(p *Params) {
	meanTotal := 0.0
	varTotal := 0.0
	min := 1.0
	max := 0.0
	count := 0.0
	i := 0
	for i < len(p.NodeList) {
		n := p.NodeList[i]
		if n.Valid {
			v := n.Sensitivity / n.InitialSensitivity
			meanTotal += v
			varTotal += v * v
			if v < min {
				min = v
			}
			if v > max {
				max = v
			}
			count += 1
		}
		i += 1
	}

	mean := meanTotal / count
	meanSquare := mean * mean

	varianceSquare := (varTotal / count) - meanSquare
	variance := math.Sqrt(varianceSquare)

	fmt.Fprintf(p.DriftExploreFile, "%v %v %v %v %v %v\n", p.CurrentTime, mean, variance, min, max, count)

}