//CPS contains the simulator
package cps

import (
	"encoding/csv"
	"flag"
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
		p.BoolGrid[int(oldX)][int(oldY)] = false //set the old spot false since the node will now move away

		//move the node to its new location
		p.NodeList[j].Move(p)

		//set the new location in the boolean field to true
		newX, newY := p.NodeList[j].GetLoc()
		p.BoolGrid[int(newX)][int(newY)] = true

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
	curNode := NodeImpl{P: p, X: -1, Y: -1, Id: len(p.NodeList), SampleHistory: initHistory, Concentration: 0,
		Cascade: nodeNum, Battery: p.BatteryCharges[nodeNum], BatteryLossScalar: p.BatteryLosses[nodeNum],
		BatteryLossSensor:        p.BatteryLossesSensor[nodeNum],
		BatteryLossGPS:           p.BatteryLossesGPS[nodeNum],
		BatteryLossServer:        p.BatteryLossesServer[nodeNum],
		BatteryLoss4G:            p.BatteryLosses4G[nodeNum],
		BatteryLossAccelerometer: p.BatteryLossesAccelerometer[nodeNum],
		BatteryLossWifi:		  p.BatteryLossesWiFi[nodeNum],
		BatteryLossBT:			  p.BatteryLossesBT[nodeNum]}

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

	//set node Time and initial sensitivity
	curNode.NodeTime = 0
	curNode.InitialSensitivity = s0 + (s1)*math.Exp(-float64(curNode.NodeTime)/p.Tau1) + (s2)*math.Exp(-float64(curNode.NodeTime)/p.Tau2)
	curNode.Sensitivity = curNode.InitialSensitivity

	return &curNode
}

func SetupCSVNodes(p *Params) {
	for i := 0; i < p.TotalNodes; i++ {
		newNode := InitializeNodeParameters(p, i)

		newNode.X = float32(p.NodeMovements[i][0].X)
		newNode.Y = float32(p.NodeMovements[i][1].Y)

		if newNode.InBounds(p) {
			newNode.Valid = true
			p.BoolGrid[int(newNode.X)][int(newNode.Y)] = true
		} else {
			newNode.Valid = false
		}

		p.NodeList = append(p.NodeList, newNode)
		p.CurrentNodes += 1
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
	defer dummy.Close()
	p.PositionFile, err = os.Create(p.OutputFileNameCM + "-simulatorOutput.txt")
	if err != nil {
		log.Fatal("Cannot create file", err)
	}
	//defer p.PositionFile.Close()

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
	fmt.Fprintln(p.DriftFile, "Sensor Loss:", p.SamplingLossServerCM, "\nGPS Loss:", p.SamplingLossGPSCM, "\nServer Loss:", p.SamplingLossServerCM)
	fmt.Fprintln(p.DriftFile, "BlueTooth Loss:", p.SamplingLossBTCM, "\nWiFi Loss:", p.SamplingLossWifiCM, "\n4G Loss:", p.SamplingLoss4GCM, "\nAccelerometer Loss:", p.SamplingLossAccelCM)
	fmt.Fprintln(p.DriftFile, "Printing Position:", p.PositionPrint, "\nPrinting Energy:", p.EnergyPrint, "\nPrinting Nodes:", p.NodesPrint)
	fmt.Fprintln(p.DriftFile, "Super Nodes:", p.NumSuperNodes, "\nSuper Node Type:", p.SuperNodeType, "\nSuper Node Speed:", p.SuperNodeSpeed, "\nSuper Node Radius:", p.SuperNodeRadius)
	fmt.Fprintln(p.DriftFile, "Error Multiplier:", p.ErrorModifierCM)
	fmt.Fprintln(p.DriftFile, "--------------------")

	p.GridFile, err = os.Create(p.OutputFileNameCM + "-grid.txt")
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

	p.BatteryFile, err = os.Create(p.OutputFileNameCM + "-batteryusage.txt")
	if err != nil {
		log.Fatal("Cannot create file", err)
	}

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
	//p.BatteryCharges = GetLinearBatteryValues(len(p.NodeEntryTimes))
	p.BatteryCharges = GetNormDistro(p.TotalNodes, 70.0, 15) //battery values come from normal distribution
																		//mean = 70%, std = 15%
	p.BatteryLosses = GetLinearBatteryLossConstant(len(p.NodeEntryTimes), float32(p.NaturalLossCM))

	//updated because of the variable renaming to BatteryLosses__ and SamplingLoss__CM
	p.BatteryLossesSensor = GetLinearBatteryLossConstant(len(p.NodeEntryTimes), float32(p.SamplingLossSensorCM))
	p.BatteryLossesGPS = GetLinearBatteryLossConstant(len(p.NodeEntryTimes), float32(p.SamplingLossGPSCM))
	p.BatteryLossesServer = GetLinearBatteryLossConstant(len(p.NodeEntryTimes), float32(p.SamplingLossServerCM))
	//newly added for BlueTooth, Wifi, 4G, and Accelerometer battery usage
	p.BatteryLossesBT = GetLinearBatteryLossConstant(len(p.NodeEntryTimes), float32(p.SamplingLossBTCM))
	p.BatteryLossesWiFi = GetLinearBatteryLossConstant(len(p.NodeEntryTimes), float32(p.SamplingLossWifiCM))
	p.BatteryLosses4G = GetLinearBatteryLossConstant(len(p.NodeEntryTimes), float32(p.SamplingLoss4GCM))
	p.BatteryLossesAccelerometer = GetLinearBatteryLossConstant(len(p.NodeEntryTimes), float32((p.SamplingLossAccelCM)))

	p.Attractions = make([]*Attraction, p.NumAtt)

	if p.CSVSensor {
		readSensorCSV(p)
	} else {
		p.MaxRaw = 1000
		//p.EdgeRaw = 36
		p.EdgeRaw = RawConcentration(float32(p.DetectionDistance))
		fmt.Println("Raw:", p.EdgeRaw)
		p.MaxADC = 4095
		p.EdgeADC = 3
		p.ADCWidth = p.MaxRaw/p.MaxADC
		p.ADCOffset = p.EdgeRaw - p.EdgeADC * p.ADCWidth
	}

	if p.CSVMovement {
		readMovementCSV(p)
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
			if i != 0 && j != 0 {
				if i+x > 0 && i+x < p.Width {
					if j+y > 0 && j+y < p.Height {
						total += InterpolateFloat(float32(p.SensorReadings[x][y][time]),  float32(p.SensorReadings[i + x][j + y][time]), .2/float32(Dist(Tuple{x, y}, Tuple{x+i, y+j})))
						counted += 1
					}
				}
			}
		}
	}

	p.MaxRaw = total/counted


	straight_increment := int(math.Ceil(p.DetectionDistance))
	fmt.Println(straight_increment)
	diag_increment := int(math.Ceil(p.DetectionDistance/math.Sqrt2))
	fmt.Println(diag_increment)

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

	for _,point := range(points) {
		if point[0] + x > 0 && point[0] + x < p.Width {
			if point[1] + y > 0 && point[1] + y < p.Height {
				//if legal then interpolate the correct value based on distance
				total += InterpolateFloat(float32(p.SensorReadings[x][y][time]), float32(p.SensorReadings[x + point[0]][y + point[1]][time]), float32(p.DetectionDistance/(Dist(Tuple{x, y}, Tuple{x + point[0], y + point[1]}))))
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


	fmt.Printf("\n")
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
	p.MaxTimeStep = len(times)

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
	/*maxReading := 0.0
	maxLocX := 0
	maxLocY := 0
	maxTime := 0*/
	for i < len(record) {


		x, _ := strconv.ParseInt(record[i][0], 10, 32);
		y, _ := strconv.ParseInt(record[i][1], 10, 32);

		j := 2

		for j < len(record[i]) {
			read1, _ := strconv.ParseFloat(record[i][j], 32);
			if read1 < 0 {
				read1 = 0
			}
			//fmt.Printf("x:%v, y:%v, j:%v\n",x,y,j)
			p.SensorReadings[x][y][j-2] = read1

			/*if read1 > maxReading {
				maxReading = read1
				maxLocX = int(x)
				maxLocY = int(y)
				maxTime = j-2
			}*/

			j += 1

		}
		i++

		if(i % 1000 == 0) {
			prog := int(float32(i)/float32(len(record))*100)
			fmt.Printf("\rProgress [%s%s] %d ", strings.Repeat("=", prog), strings.Repeat(".", 100-prog), prog)
		}
	}



	//CalculateADCSetting(maxReading, maxLocX, maxLocY, maxTime, p)
	fmt.Println(p.BombX, p.BombY)
	CalculateADCSetting(p.SensorReadings[p.B.X][p.B.Y][10], p.B.X, p.B.Y, 10, p)




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

		for iter < len(record[time])-1 && iter/2 < p.TotalNodes {
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
	flag.IntVar(&p.NegativeSittingStopThresholdCM, "negativeSittingStopThreshold", -10,
		"Negative number sitting is set to when board map is reset")
	flag.IntVar(&p.SittingStopThresholdCM, "sittingStopThreshold", 5,
		"How long it takes for a node to stay seated")
	flag.Float64Var(&p.GridCapacityPercentageCM, "GridCapacityPercentage", .9,
		"Percent the sub-Grid can be filled")
	flag.StringVar(&p.InputFileNameCM, "inputFileName", "Log1_in.txt",
		"Name of the input text file")
	flag.StringVar(&p.SensorPath, "sensorPath", "Circle_2D.csv", "Sensor Reading Inputs")
	flag.StringVar(&p.MovementPath, "movementPath", "Circle_2D.csv", "Movement Inputs")
	flag.StringVar(&p.OutputFileNameCM, "p.OutputFileName", "Log",
		"Name of the output text file prefix")
	flag.Float64Var(&p.NaturalLossCM, "naturalLoss", .005,
		"battery loss due to natural causes")
	flag.Float64Var(&p.SamplingLossSensorCM, "sensorSamplingLoss", .001,
		"battery loss due to sensor sampling")
	flag.Float64Var(&p.SamplingLossGPSCM, "GPSSamplingLoss", .005,
		"battery loss due to GPS sampling")
	flag.Float64Var(&p.SamplingLossSensorCM, "serverSamplingLoss", .01,
		"battery loss due to server sampling")
	flag.Float64Var(&p.SamplingLossBTCM, "SamplingLossBTCM", .0001,
		"battery loss due to BlueTooth sampling")
	flag.Float64Var(&p.SamplingLossWifiCM, "SamplingLossWifiCM", .001,
		"battery loss due to WiFi sampling")
	flag.Float64Var(&p.SamplingLoss4GCM, "SamplingLoss4GCM", .005,
		"battery loss due to 4G sampling")
	flag.Float64Var(&p.SamplingLossAccelCM, "SamplingLossAccelCM", .001,
		"battery loss due to accelerometer sampling")
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
	flag.IntVar(&p.GridStoredSamplesCM, "p.GridStoredSamples", 10,
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
	flag.Float64Var(&p.CalibrationThresholdCM, "Recalibration Threshold", 3.0, "Value over grid average to recalibrate node")
	flag.Float64Var(&p.StdDevThresholdCM, "StandardDeviationThreshold", 1.7, "Detection Threshold based on standard deviations from mean")
	flag.Float64Var(&p.DetectionDistance, "detectionDistance", 6.0, "Detection Distance")
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
	flag.IntVar(&p.SuperNodeType, "p.SuperNodeType", 0, "the type of super node used in the simulator")

	//Range: 0-...
	//Theoretically could be as high as possible
	//Realistically should remain around 10
	flag.IntVar(&p.SuperNodeSpeed, "p.SuperNodeSpeed", 3, "the speed of the super node")

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

	flag.BoolVar(&p.PositionPrintCM, "logPosition", false, "Whether you want to write position info to a log file")
	flag.BoolVar(&p.GridPrintCM, "logGrid", false, "Whether you want to write p.Grid info to a log file")
	flag.BoolVar(&p.EnergyPrintCM, "logEnergy", false, "Whether you want to write energy into to a log file")
	flag.BoolVar(&p.NodesPrintCM, "logNodes", false, "Whether you want to write node readings to a log file")
	flag.IntVar(&p.SquareRowCM, "SquareRowCM", 50, "Number of rows of p.Grid squares, 1 through p.MaxX")
	flag.IntVar(&p.SquareColCM, "SquareColCM", 50, "Number of columns of p.Grid squares, 1 through p.MaxY")

	flag.StringVar(&p.ImageFileNameCM, "imageFileName", "circle_justWalls_x4.png", "Name of the input text file")
	flag.StringVar(&p.StimFileNameCM, "stimFileName", "circle_0.txt", "Name of the stimulus text file")
	flag.StringVar(&p.OutRoutingNameCM, "outRoutingName", "log.txt", "Name of the stimulus text file")
	flag.StringVar(&p.OutRoutingStatsNameCM, "outRoutingStatsName", "routingStats.txt", "Name of the output file for stats")

	flag.BoolVar(&p.RegionRouting, "regionRouting", true, "True if you want to use the new routing algorithm with regions and cutting")

	flag.Parse()
	fmt.Println("Natural Loss: ", p.NaturalLossCM)
	fmt.Println("Sensor Sampling Loss: ", p.SamplingLossSensorCM)
	fmt.Println("GPS sampling loss: ", p.SamplingLossGPSCM)
	fmt.Println("Server sampling loss", p.SamplingLossServerCM)
	fmt.Println("BlueTooth sampling loss", p.SamplingLossBTCM)
	fmt.Println("WiFi Sampling Loss", p.SamplingLossWifiCM)
	fmt.Println("4G Sampling Loss", p.SamplingLossWifiCM)
	fmt.Println("Accelerometer Sampling Loss", p.SamplingLossAccelCM)
	fmt.Println("Threshold Battery to use: ", p.ThresholdBatteryToUseCM)
	fmt.Println("Threshold battery to have: ", p.ThresholdBatteryToHaveCM)
	fmt.Println("Moving speed for incresed sampling: ", p.MovementSamplingSpeedCM)
	fmt.Println("Period of extra sampling due to high speed: ", p.MovementSamplingPeriodCM)
	fmt.Println("Maximum size of buffer posible: ", p.MaxBufferCapacityCM)
	fmt.Println("Energy model type:", p.EnergyModelCM)
	fmt.Println("Sensor Sampling Period:", p.SensorSamplingPeriodCM)
	fmt.Println("GPS Sampling Period:", p.GPSSamplingPeriodCM)
	fmt.Println("Server Sampling Period:", p.ServerSamplingPeriodCM)
	fmt.Println("Number of Node Stored Samples:", p.NumStoredSamplesCM)
	fmt.Println("Number of Grid Stored Samples:", p.GridStoredSamplesCM)
	fmt.Println("Detection Threshold:", p.DetectionThresholdCM)

	//fmt.Println("tail:", flag.Args())
}
