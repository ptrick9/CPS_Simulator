//CPS contains the simulator
package cps

import (
	"encoding/csv"
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
	p.NumPointsOfInterestStatic = len(p.Poispos)
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
		p.BoolGrid[oldX][oldY] = false //set the old spot false since the node will now move away

		//move the node to its new location
		p.NodeList[j].Move(p)

		//set the new location in the boolean field to true
		newX, newY := p.NodeList[j].GetLoc()
		p.BoolGrid[newX][newY] = true

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
			p.BoolGrid[oldX][oldY] = false //set the old spot false since the node will now move away
		}
		//move the node to its new location
		//p.NodeList[j].Move(p)

		id := p.NodeList[j].GetID()
		p.NodeList[j].X = p.NodeMovements[id][time].X
		p.NodeList[j].Y = p.NodeMovements[id][time].Y


		//set the new location in the boolean field to true
		newX, newY := p.NodeList[j].GetLoc()

		if p.NodeList[j].InBounds(p) {
			p.NodeList[j].Valid = true
		} else {
			p.NodeList[j].Valid = false
		}
		if p.NodeList[j].Valid {
			p.BoolGrid[newX][newY] = true
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
	curNode := NodeImpl{P: p, X: -1, Y: -1, Id: len(p.NodeList), SampleHistory: initHistory, Concentration: 0,
		Cascade: nodeNum, Battery: p.BatteryCharges[nodeNum], BatteryLossScalar: p.BatteryLosses[nodeNum],
		BatteryLossCheckingSensorScalar: p.BatteryLossesCheckingSensorScalar[nodeNum],
		BatteryLossGPSScalar:            p.BatteryLossesCheckingGPSScalar[nodeNum],
		BatteryLossCheckingServerScalar: p.BatteryLossesCheckingServerScalar[nodeNum]}


	//values to determine coefficients
	curNode.SetS0(rand.Float64()*0.2 + 0.1)
	curNode.SetS1(rand.Float64()*0.2 + 0.1)
	curNode.SetS2(rand.Float64()*0.2 + 0.1)
	//values to determine error in coefficients
	s0, s1, s2 := curNode.GetCoefficients()
	curNode.SetE0(rand.Float64() * 0.1 * p.ErrorModifierCM * s0)
	curNode.SetE1(rand.Float64() * 0.1 * p.ErrorModifierCM * s1)
	curNode.SetE2(rand.Float64() * 0.1 * p.ErrorModifierCM * s2)
	//Values to determine error in exponents
	curNode.SetET1(p.Tau1 * rand.Float64() * p.ErrorModifierCM * 0.05)
	curNode.SetET2(p.Tau1 * rand.Float64() * p.ErrorModifierCM * 0.05)

	//set node time and initial sensitivity
	curNode.NodeTime = 0
	curNode.InitialSensitivity = s0 + (s1)*math.Exp(-float64(curNode.NodeTime)/p.Tau1) + (s2)*math.Exp(-float64(curNode.NodeTime)/p.Tau2)
	curNode.Sensitivity = curNode.InitialSensitivity

	return &curNode
}

func SetupCSVNodes(p *Params) {
	for i := 0; i < p.TotalNodes; i++ {
		newNode := InitializeNodeParameters(p, i)

		newNode.X = p.NodeMovements[i][0].X
		newNode.Y = p.NodeMovements[i][1].Y

		if newNode.InBounds(p) {
			newNode.Valid = true
			p.BoolGrid[newNode.X][newNode.Y] = true
		} else {
			newNode.Valid = false
		}

		p.NodeList = append(p.NodeList, *newNode)
		p.CurrentNodes += 1
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

			newNode.X = xx
			newNode.Y = yy

			newNode.Valid = true

			p.BoolGrid[xx][yy] = true

			p.NodeList = append(p.NodeList, *newNode)
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
	defer dummy.Close()
	p.PositionFile, err = os.Create(p.OutputFileNameCM + "-simulatorOutput.txt")
	if err != nil {
		log.Fatal("Cannot create file", err)
	}
	//defer p.PositionFile.Close()

	//Print parameters to position file
	fmt.Fprintln(p.PositionFile, "Width:", p.MaxX)
	fmt.Fprintln(p.PositionFile, "Height:", p.MaxY)
	fmt.Fprintf(p.PositionFile, "Amount: %-8v\n", p.Iterations_of_event)
	fmt.Fprintf(p.PositionFile, "Bomb x: %v\n", p.BombX)
	fmt.Fprintf(p.PositionFile, "Bomb y: %v\n", p.BombY)

	p.DriftFile, err = os.Create(p.OutputFileNameCM + "-drift.txt")
	if err != nil {
		log.Fatal("Cannot create file", err)
	}
	//defer p.DriftFile.Close()

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
	fmt.Fprintln(p.DriftFile, "Battery Natural Loss:", p.NaturalLossCM)
	fmt.Fprintln(p.DriftFile, "Sensor Loss:", p.SensorSamplingLossCM, "\nGPS Loss:", p.GPSSamplingLossCM, "\nServer Loss:", p.ServerSamplingLossCM)
	fmt.Fprintln(p.DriftFile, "Printing Position:", p.PositionPrint, "\nPrinting Energy:", p.EnergyPrint, "\nPrinting Nodes:", p.NodesPrint)
	fmt.Fprintln(p.DriftFile, "Super Nodes:", p.NumSuperNodes, "\nSuper Node Type:", p.SuperNodeType, "\nSuper Node Speed:", p.SuperNodeSpeed, "\nSuper Node Radius:", p.SuperNodeRadius)
	fmt.Fprintln(p.DriftFile, "Error Multiplier:", p.ErrorModifierCM)
	fmt.Fprintln(p.DriftFile, "--------------------")

	p.GridFile, err = os.Create(p.OutputFileNameCM + "-Grid.txt")
	if err != nil {
		log.Fatal("Cannot create file", err)
	}
	//defer p.GridFile.Close()

	//Write parameters to gridFile
	fmt.Fprintln(p.GridFile, "Width:", p.SquareColCM)
	fmt.Fprintln(p.GridFile, "Height:", p.SquareRowCM)

	p.NodeFile, err = os.Create(p.OutputFileNameCM + "-node_reading.txt")
	if err != nil {
		log.Fatal("Cannot create file", err)
	}
	//defer p.NodeFile.Close()

	p.EnergyFile, err = os.Create(p.OutputFileNameCM + "-node.txt")
	if err != nil {
		log.Fatal("Cannot create file", err)
	}
	//defer p.EnergyFile.Close()

	p.RoutingFile, err = os.Create(p.OutputFileNameCM + "-path.txt")
	if err != nil {
		log.Fatal("Cannot create file", err)
	}
	//defer p.RoutingFile.Close()

	p.BoolFile, err = os.Create(p.OutputFileNameCM + "-bool.txt")
	if err != nil {
		log.Fatal("Cannot create file", err)
	}
	//defer p.BoolFile.Close()

	p.AttractionFile, err = os.Create(p.OutputFileNameCM + "-attraction.txt")
	if err != nil {
		log.Fatal("Cannot create file", err)
	}

	p.MoveReadingsFile, err = os.Create(p.OutputFileNameCM + "-movementReadings.txt")
	if err != nil {
		log.Fatal("Cannot create file", err)
	}

	//defer p.AttractionFile.Close()

	p.ServerFile, err = os.Create(p.OutputFileNameCM + "-server.txt")
	if err != nil {
		log.Fatal("Cannot create file", err)
	}
	//defer p.ServerFile.Close()

	p.NodeTest, err = os.Create(p.OutputFileNameCM + "-nodeTest.txt")
	if err != nil {
		log.Fatal("Cannot create file", err)
	}
	//defer p.ServerFile.Close()

	p.NodeTest2, err = os.Create(p.OutputFileNameCM + "-nodeTest2.txt")
	if err != nil {
		log.Fatal("Cannot create file", err)
	}
	//defer p.ServerFile.Close()

	p.DetectionFile, err = os.Create(p.OutputFileNameCM + "-detection.txt")
	if err != nil {
		log.Fatal("Cannot create file", err)
	}
	//defer p.ServerFile.Close()
}

func SetupParameters(p *Params) {

	p.XDiv = p.MaxX / p.SquareColCM
	p.YDiv = p.MaxY / p.SquareRowCM

	//The capacity for a square should be equal to the area of the square
	//So we take the side length (xDiv) and square it
	p.SquareCapacity = int(math.Pow(float64(p.XDiv), 2))

	//Center of the p.Grid
	p.Center.X = p.MaxX / 2
	p.Center.Y = p.MaxY / 2

	p.TotalPercentBatteryToUse = float32(p.ThresholdBatteryToUseCM)
	p.BatteryCharges = GetLinearBatteryValues(len(p.NodeEntryTimes))
	p.BatteryLosses = GetLinearBatteryLossConstant(len(p.NodeEntryTimes), float32(p.NaturalLossCM))
	p.BatteryLossesCheckingSensorScalar = GetLinearBatteryLossConstant(len(p.NodeEntryTimes), float32(p.SensorSamplingLossCM))
	p.BatteryLossesCheckingGPSScalar = GetLinearBatteryLossConstant(len(p.NodeEntryTimes), float32(p.GPSSamplingLossCM))
	p.BatteryLossesCheckingServerScalar = GetLinearBatteryLossConstant(len(p.NodeEntryTimes), float32(p.ServerSamplingLossCM))
	p.Attractions = make([]*Attraction, p.NumAtt)

	if p.CSVSensor {
		readSensorCSV(p)
	}
	if p.CSVMovement {
		readMovementCSV(p)
	}


}

//readSensorCSV reads the sensor values from a file
func readSensorCSV(p *Params) {

	in, err := os.Open(p.SensorPath)
	if err != nil {
		println("error opening file")
	}
	defer in.Close()
	fmt.Printf("Reading Sensor Files\n")
	r := csv.NewReader(in)
	r.FieldsPerRecord = -1
	record, err := r.ReadAll()

	reg, _ := regexp.Compile("([0-9]+)")
	times := reg.FindAllStringSubmatch(strings.Join(record[0], " "), -1)

	p.SensorTimes = make([]int, len(times))
	for i := range times {
		p.SensorTimes[i], _ = strconv.Atoi(times[i][1])
	}

	numSamples := len(record[2]) - 2

	p.SensorReadings = make([][][]float64, p.Width)
	for i := range p.SensorReadings {
		p.SensorReadings[i] = make([][]float64, p.Height)
		for j := range p.SensorReadings[i] {
			p.SensorReadings[i][j] = make([]float64, numSamples)
			for k := range p.SensorReadings[i][j] {
				p.SensorReadings[i][j][k] = 0
			}
		}
	}



	i := 1
	fmt.Printf("Sensor CSV Processing\n")
	for i < len(record) {


		x, _ := strconv.ParseInt(record[i][0], 10, 32);
		y, _ := strconv.ParseInt(record[i][1], 10, 32);

		j := 2

		for j < len(record[i]) {
			read1, _ := strconv.ParseFloat(record[i][j], 32);
			//fmt.Printf("x:%v, y:%v, j:%v\n",x,y,j)
			p.SensorReadings[x][y][j-2] = read1
			j += 1
		}
		i++

		if(i % 1000 == 0) {
			prog := int(float32(i)/float32(len(record))*100)
			fmt.Printf("\rProgress [%s%s] %d ", strings.Repeat("=", prog), strings.Repeat(".", 100-prog), prog)
		}
	}
	fmt.Printf("\n")
}

//readMovementCSV reads the movement parameters from a file
func readMovementCSV(p *Params) {

	in, err := os.Open(p.MovementPath)
	if err != nil {
		println("error opening file")
	}
	defer in.Close()

	r := csv.NewReader(in)
	r.FieldsPerRecord = -1
	record, err := r.ReadAll()


	timeSteps := len(record)


	p.NodeMovements = make([][]Tuple, p.TotalNodes)
	for i := range p.NodeMovements {
		p.NodeMovements[i] = make([]Tuple, timeSteps)
	}


	time := 0
	fmt.Printf("Movement CSV Processing %d TimeSteps for %d Nodes  %d\n", len(record), len(record[0])/2, p.TotalNodes)
	for time < len(record) {
		iter := 0

		for iter < len(record[time])-1 {
			x, _ := strconv.ParseInt(record[time][iter], 10, 32);
			y, _ := strconv.ParseInt(record[time][iter+1], 10, 32);

			p.NodeMovements[iter/2][time] = Tuple{int(x), int(y)}
			iter += 2
		}
		time++

		if(time % 10 == 0) {
			prog := int(float32(time)/float32(len(record))*100)
			fmt.Printf("\rProgress [%s%s] %d ", strings.Repeat("=", prog), strings.Repeat(".", 100-prog), prog)
		}
	}
	fmt.Printf("\n")
}


func RangeInt(min, max int) int { //returns a random number between max and min
	return rand.Intn(max-min) + min
}

// This prints the board map to the terminal
//func printBoardMap(){
//	for i:= 0; i < len(p.BoardMap); i++{
//		fmt.Println(p.BoardMap[i])
//	}
//}

//var fileBoard, errBoard = os.Create("p.BoardMap.txt")
/*
// This prints board Map to a txt file.
func writeBordMapToFile2() {
	start := Time.Now()
	defer func() {
		elapsed := Time.Since(start)
		fmt.Println("Printing Board Map took", elapsed)
	}()
	Check(errBoard)
	var s = ""
	s = s + "\nt=" + strconv.Itoa(iterations_used) + "\n\n"
	for i := 0; i < len(p.BoardMap); i++ {
		for j := 0; j < len(p.BoardMap[i]); j++ {
			s = s + strconv.Itoa(p.BoardMap[i][j]) + " "
		}
		s = s + "\n"
	}
	//s = s + "\nt=" + strconv.Itoa(iterations_used) + "\n\n"
	n3, err := fileBoard.WriteString(s)
	Check(err)
	fmt.Printf("wrote %d bytes\n", n3)
}

func writeBordMapToFile() {
	//start := Time.Now()
	Check(errBoard)
	w := bufio.NewWriter(fileBoard)
	w.WriteString("\nt=" + strconv.Itoa(iterations_used) + "\n\n")
	for i := 0; i < len(p.BoardMap); i++ {
		for j := 0; j < len(p.BoardMap[i]); j++ {
			w.WriteString(strconv.Itoa(p.BoardMap[i][j]) + " ")
		}
		w.WriteString("\n")
	}
	w.Flush()
	//elapsed := Time.Since(start)
	//fmt.Println("Printing Board Map took", elapsed)
}*/
