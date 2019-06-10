package cps

import (
	"encoding/csv"
	"strings"

	//"bufio"
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"regexp"
	"strconv"
	"image"
	//"time"
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

	//b []byte
	b2 []byte
	b3 []byte
	b4 []byte
	b5 []byte
	b6 []byte
	b7 []byte

	//numNodeNodes               int
	//numWallNodes               int
	//numPoints                  int
	//numPointsOfInterestKinetic int
	//numPointsOfInterestStatic  int

	//fileName = "Log1_in.txt"

	makeBoardMapFile = true
	NodePositionMap			map[Tuple]*NodeImpl
)

//func main() {
//
//	getListedInput()
//
//	squareRow := getDashedInput("squareRow")
//	squareCol:= getDashedInput("squareCol")
//	numNodes:= getDashedInput("numNodes")
//	numStoredSamples:= getDashedInput("numStoredSamples")
//	Tau1:= getDashedInput("Tau1")
//	Tau2:= getDashedInput("Tau2")
//	superNodeType:= getDashedInput("superNodeType")
//	maxX:= getDashedInput("maxX")
//	maxY:= getDashedInput("maxY")
//	bombX:= getDashedInput("bombX")
//	bombY:= getDashedInput("bombY")
//	Tester:= getDashedInput("Tester")
//
//	fmt.Println(squareRow,
//		squareCol,
//		numNodes,
//		numStoredSamples,
//		Tau1,
//		Tau2,
//		superNodeType,
//		maxX,
//		maxY,
//		bombX,
//		bombY,
//		Tester)
//
//	createBoard(maxX,maxY)
//
//	fillInWallsToBoard()
//
//	fillInBufferCurrent()
//
//	fillPointsToBoard()
//
//	fillInMap()
//
//	writeBordMapToFile()
//
//}

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

func GetListedInput(p *Params) {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(dir)
	b := ReadFromFile(p.FileName)
	r := regexp.MustCompile("N: [0-9]+")
	w := r.FindAllString(string(b), 10)
	r = regexp.MustCompile("[0-9]+")
	w = r.FindAllString(w[0], 10)
	s, err := strconv.ParseInt(w[0], 10, 32)
	Check(err)
	r = regexp.MustCompile("x:[0-9]+, y:[0-9]+, t:[0-9]+")
	fai := r.FindAllIndex(b, int(s))
	w = r.FindAllString(string(b), int(s))
	if len(fai) > 0 {
		b2 = b[fai[len(fai)-1][1]:]
	} else {
		b2 = b
	}
	//fmt.Println(w)
	FillInts(w, 0, p)
	//fmt.Println(npos)

	r = regexp.MustCompile("W: [0-9]+")
	w = r.FindAllString(string(b2), 10)
	r = regexp.MustCompile("[0-9]+")
	w = r.FindAllString(w[0], 10)
	s, err = strconv.ParseInt(w[0], 10, 32)
	r = regexp.MustCompile("x:[0-9]+, y:[0-9]+")
	Check(err)
	fai = r.FindAllIndex(b2, int(s))
	w = r.FindAllString(string(b2), int(s))
	if len(fai) > 0 {
		b3 = b2[fai[len(fai)-1][1]:]
	} else {
		b3 = b2
	}
	//fmt.Println(w)
	FillInts(w, 1, p)
	//fmt.Println(wpos)

	r = regexp.MustCompile("S: [0-9]+")
	w = r.FindAllString(string(b3), 10)
	r = regexp.MustCompile("[0-9]+")
	w = r.FindAllString(w[0], 10)
	s, err = strconv.ParseInt(w[0], 10, 32)
	r = regexp.MustCompile("x:[0-9]+, y:[0-9]+, t:[0-9]+")
	Check(err)
	fai = r.FindAllIndex(b3, int(s))
	w = r.FindAllString(string(b3), int(s))
	if len(fai) > 0 {
		b4 = b3[fai[len(fai)-1][1]:]
	} else {
		b4 = b3
	}
	//fmt.Println(w)
	FillInts(w, 2, p)
	//fmt.Println(spos)

	/*r = regexp.MustCompile("P: [0-9]+")
	w = r.FindAllString(string(b4), 10)
	r = regexp.MustCompile("[0-9]+")
	w = r.FindAllString(w[0], 10)
	s, err = strconv.ParseInt(w[0], 10, 32)
	r = regexp.MustCompile("x:[0-9]+, y:[0-9]+, t:[0-9]+")
	Check(err)
	fai = r.FindAllIndex(b4, int(s))
	w = r.FindAllString(string(b4), int(s))
	if len(fai) > 0 {
		b5 = b4[fai[len(fai)-1][1]:]
	} else {
		b5 = b4
	}
	fmt.Println(w)
	fillInts(w, 3)
	fmt.Println(ppos)*/

	b5 = b4

	r = regexp.MustCompile("POIS: [0-9]+")
	w = r.FindAllString(string(b5), 10)
	r = regexp.MustCompile("[0-9]+")
	w = r.FindAllString(w[0], 10)
	s, err = strconv.ParseInt(w[0], 10, 32)
	r = regexp.MustCompile("x:[0-9]+, y:[0-9]+, ti:[0-9]+, to:[0-9]+")
	Check(err)
	fai = r.FindAllIndex(b5, int(s))
	w = r.FindAllString(string(b5), int(s))
	if len(fai) > 0 {
		b6 = b5[fai[len(fai)-1][1]:]
	} else {
		b6 = b5
	}
	//fmt.Println(w)
	FillInts(w, 5, p)
	//fmt.Println(poispos)

	p.NumNodeNodes = len(p.Npos)
	p.NumWallNodes = len(p.Wpos)
	//numPoints = len(ppos)
	p.NumPointsOfInterestStatic = len(p.Poispos)
	//fmt.Println(numNodeNodes, numWallNodes, numPoints, numPointsOfInterestKinetic, numPointsOfInterestStatic)
}

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
			p.Npos = append(p.Npos, ap)

			r = regexp.MustCompile("y:[0-9]+")
			X = r.FindAllString(s[i], 1)
			r = regexp.MustCompile("[0-9]+")
			x = r.FindAllString(X[0], 1)
			s1, err = strconv.ParseInt(x[0], 10, 32)
			Check(err)
			p.Npos[i][1] = int(s1)

			r = regexp.MustCompile("t:[0-9]+")
			X = r.FindAllString(s[i], 1)
			r = regexp.MustCompile("[0-9]+")
			x = r.FindAllString(X[0], 1)
			s1, err = strconv.ParseInt(x[0], 10, 32)
			Check(err)
			p.Npos[i][2] = int(s1)
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
	for i := 0; i < y; i++ {
		p.BoardMap = append(p.BoardMap, []int{})
		for j := 0; j < x; j++ {
			p.BoardMap[i] = append(p.BoardMap[i], 0)
		}
	}
}

func HandleMovement(p *Params) {
	for j := 0; j < len(p.NodeList); j++ {

		oldX, oldY := p.NodeList[j].GetLoc()
		p.BoolGrid[oldY][oldX] = false //set the old spot false since the node will now move away

		//move the node to its new location
		p.NodeList[j].Move(p)

		//set the new location in the boolean field to true
		newX, newY := p.NodeList[j].GetLoc()
		p.BoolGrid[newY][newX] = true

		//writes the node information to the file
		if p.EnergyPrint {
			fmt.Fprintln(p.EnergyFile, p.NodeList[j])
		}

		//Add the node into its new Square's p.NumNodes
		//If the node hasn't left the square, that Square's p.NumNodes will
		//remain the same after these calculations
	}
}


// Fills the walls into the board based on the wall positions extrapolated from the file
func FillInWallsToBoard(p *Params) {
	for i := 0; i < len(p.Wpos); i++ {
		p.BoardMap[p.Wpos[i][1]][p.Wpos[i][0]] = -1
	}
}



// Fills the points of interest into the current buffer

func FillInBufferCurrent(p *Params) {
	bufferCurrent = [][]int{}
	for i := 0; i < len(p.Poispos); i++ {
		if p.Iterations_used >= p.Poispos[i][2] && p.Iterations_used < p.Poispos[i][3] {
			bufferCurrent = append(bufferCurrent, []int{p.Poispos[i][1], p.Poispos[i][0]})
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

// Fills in board map with the path finding values
func FillInMap1(p *Params) {
	/*start := time.Now()

	defer func() {
		elapsed := time.Since(start)
		//fmt.Println("Board Map took", elapsed)
	}()*/

	for len(bufferFuture) > 0 {
		bufferFuture = [][]int{}
		for i := 0; i < len(bufferCurrent); i++ {
			CheckLeft(i, p)
			CheckRight(i, p)
			CheckUp(i, p)
			CheckDown(i, p)
		}
		bufferCurrent = [][]int{}
		bufferCurrent = append(bufferCurrent, bufferFuture...)
		starter += 1
	}
	starter = 1
	bufferFuture = [][]int{{}}
}

func CheckLeft(i int, p *Params) {
	defer func() {
		recover()
	}()
	if
	//bufferCurrent[i][0]-1 < len(p.BoardMap) &&
	//bufferCurrent[i][1] < len(p.BoardMap[1]) &&
	//bufferCurrent[i][0]-1 >= 0 &&
	//bufferCurrent[i][1] >= 0 &&
	p.BoardMap[bufferCurrent[i][0]-1][bufferCurrent[i][1]] == 0 {

		p.BoardMap[bufferCurrent[i][0]-1][bufferCurrent[i][1]] = starter + 1
		bufferFuture = append(bufferFuture, []int{bufferCurrent[i][0] - 1, bufferCurrent[i][1]})
	}
}

func CheckRight(i int, p *Params) {
	defer func() {
		recover()
	}()
	if
	//bufferCurrent[i][0]+1 < len(p.BoardMap) &&
	//bufferCurrent[i][1] < len(p.BoardMap[1]) && // p.BoardMap[1] to
	//bufferCurrent[i][0]+1 >= 0 &&
	//bufferCurrent[i][1] >= 0 &&
	p.BoardMap[bufferCurrent[i][0]+1][bufferCurrent[i][1]] == 0 {

		p.BoardMap[bufferCurrent[i][0]+1][bufferCurrent[i][1]] = starter + 1
		bufferFuture = append(bufferFuture, []int{bufferCurrent[i][0] + 1, bufferCurrent[i][1]})
	}
}
func CheckUp(i int, p *Params) {
	defer func() {
		recover()
	}()
	if
	//bufferCurrent[i][0] < len(p.BoardMap) &&
	//bufferCurrent[i][1]+1 < len(p.BoardMap[1]) &&
	//bufferCurrent[i][0] >= 0 &&
	//bufferCurrent[i][1]+1 >= 0 &&
	p.BoardMap[bufferCurrent[i][0]][bufferCurrent[i][1]+1] == 0 {

		p.BoardMap[bufferCurrent[i][0]][bufferCurrent[i][1]+1] = starter + 1
		bufferFuture = append(bufferFuture, []int{bufferCurrent[i][0], bufferCurrent[i][1] + 1})
	}
}
func CheckDown(i int, p *Params) {
	defer func() {
		recover()
	}()
	if
	//bufferCurrent[i][0] < len(p.BoardMap) &&
	//bufferCurrent[i][1]-1 < len(p.BoardMap[1]) &&
	//bufferCurrent[i][0] >= 0 &&
	//bufferCurrent[i][1]-1 >= 0 &&
	p.BoardMap[bufferCurrent[i][0]][bufferCurrent[i][1]-1] == 0 {

		p.BoardMap[bufferCurrent[i][0]][bufferCurrent[i][1]-1] = starter + 1
		bufferFuture = append(bufferFuture, []int{bufferCurrent[i][0], bufferCurrent[i][1] - 1})
	}
}

func FillInMap(p *Params) {
	/*start := time.Now()

	defer func() {
		elapsed := time.Since(start)
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

	for x := 0; x < p.Height; x++ {
		r.Point_list2 = append(r.Point_list2, make([]bool, p.Width))
	}

	for x := 0; x < p.Height; x++ {
		for y := 0; y < p.Width; y++ {
			rr, _, _, _ := img.At(x, y).RGBA()
			if rr != 0 {
				r.Point_list2[x][y] = true
				r.Point_dict[Tuple{x, y}] = true


			} else {
				r.Point_dict[Tuple{x, y}] = false
				p.BoardMap[y][x] = -1
				temp := make([] int, 2)
				temp[0] = x
				temp[1] = y
				p.Wpos = append(p.Wpos, temp)
				p.BoolGrid[y][x] = true
			}
		}
	}

	CreateBoard(p.MaxX, p.MaxY, p)
	FillInWallsToBoard(p)
	FillInBufferCurrent(p)
	FillPointsToBoard(p)
	FillInMap(p)

}

func SetupFiles(p *Params) {
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
	fmt.Fprintln(p.DriftFile, "Number of Nodes:", p.NumNodes)
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

	p.GridFile, err = os.Create(p.OutputFileNameCM + "-Grid.txt")
	if err != nil {
		log.Fatal("Cannot create file", err)
	}
	//defer p.GridFile.Close()

	//Write parameters to gridFile
	fmt.Fprintln(p.GridFile,"Width:", p.SquareColCM)
	fmt.Fprintln(p.GridFile,"Height:", p.SquareRowCM)

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
	//defer p.AttractionFile.Close()
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
	p.BatteryCharges = GetLinearBatteryValues(len(p.Npos))
	p.BatteryLosses = GetLinearBatteryLossConstant(len(p.Npos), float32(p.NaturalLossCM))

	//updated because of the variable renaming to BatteryLosses__ and SamplingLoss__CM
	p.BatteryLossesSensor = GetLinearBatteryLossConstant(len(p.Npos), float32(p.SamplingLossSensorCM))
	p.BatteryLossesGPS = GetLinearBatteryLossConstant(len(p.Npos), float32(p.SamplingLossGPSCM))
	p.BatteryLossesServer = GetLinearBatteryLossConstant(len(p.Npos), float32(p.SamplingLossServerCM))
	//newly added for BlueTooth, Wifi, 4G, and Accelerometer battery usage
	p.BatteryLossesBT = GetLinearBatteryLossConstant(len(p.Npos), float32(p.SamplingLossBTCM))
	p.BatteryLossesWiFi = GetLinearBatteryLossConstant(len(p.Npos), float32(p.SamplingLossWifiCM))
	p.BatteryLosses4G = GetLinearBatteryLossConstant(len(p.Npos), float32(p.SamplingLoss4GCM))
	p.BatteryLossesAccelerometer = GetLinearBatteryLossConstant(len(p.Npos), float32((p.SamplingLossAccelCM)))

	p.Attractions = make([]*Attraction, p.NumAtt)

	//readCSV(p)


}


func readCSV(p *Params) {

	in, err := os.Open(p.SensorPath)
	if err != nil {
		println("error opening file")
	}
	defer in.Close()

	r := csv.NewReader(in)
	r.FieldsPerRecord = -1
	record, err := r.ReadAll()

	reg, _ := regexp.Compile("t=([0-9]+)")
	times := (reg.FindAllStringSubmatch(strings.Join(record[8], " "), -1))

	p.SensorTimes = make([]int, len(times))
	for i := range times {
		p.SensorTimes[i], _ = strconv.Atoi(times[i][1])
	}

	/*fmt.Println(err)
	fmt.Println(len(record))*/
	big := 0.0
	topVal := 0.0
	i := 9

	numSamples := len(record[9])-2

	p.SensorReadings = make([][][] float64, p.Width)
	for i := range p.SensorReadings {
		p.SensorReadings[i] = make([][] float64, p.Height)
		for j := range p.SensorReadings[i] {
			p.SensorReadings[i][j] = make([] float64, numSamples)
			for k := range p.SensorReadings[i][j] {
				p.SensorReadings[i][j][k] = 0
			}
		}
	}

	averaged := make([][][] float64, p.Width)
	for i := range averaged {
		averaged[i] = make([][] float64, p.Height)
		for j := range averaged[i] {
			averaged[i][j] = make([] float64, numSamples)
			for k := range averaged[i][j] {
				averaged[i][j][k] = -1
			}
		}
	}

	for i < len(record) {
		//fmt.Println(record[i])
		x, err := strconv.ParseFloat(record[i][0], 32);
		/*if err == nil {
			fmt.Println(Round(x, 0.5))
		} else {
			fmt.Println(err)
		}*/
		if x > big {
			big = x
		}
		y, err := strconv.ParseFloat(record[i][1], 32);
		/*if err == nil {
			fmt.Println(Round(y, 0.5))
		}*/
		if y > big {
			big = y
		}
		j := 2
		/*fmt.Printf("%v %v\n", int(x*2), int(y*2))
		fmt.Printf("%v\n", len(record[i]))*/
		if (int(x*2) < p.Width && int(y*2) < p.Height) {
			for j < len(record[i]) {
				read1, _ := strconv.ParseFloat(record[i][j], 32);
				if err == nil {
					//fmt.Printf("%e ", read1)
				}

				//fmt.Printf("%v %v %v\n", int(x*2), int(y*2), j-2)
				p.SensorReadings[int(x*2)][int(y*2)][j-2] = read1
				if read1 > topVal {
					topVal = read1
				}

				j += 1
			}
		}
		//fmt.Println()
		fmt.Printf("\r%d\\%d", i, len(record))
		i++
	}





	fmt.Printf("\ntop: %v\n", topVal)


	cw := 7
	ch := 7
	divider := 1.0/float64(cw * ch)
	radius := cw / 2

	for k := range p.SensorReadings[0][0] {
		for i := range p.SensorReadings {
			for j := range p.SensorReadings[i] {
				total := 0.0
				for ci := radius; ci >= radius * -1; ci-- {
					for cj := radius; cj >= radius * -1; cj-- {
						if i+ci > 0 && i + ci < p.Width && j+cj > 0 && j+cj < p.Height {
							total += p.SensorReadings[i+ci][j+cj][k] * divider
						}
					}
				}
				averaged[i][j][k] = total
			}
		}
	}

	for k := range p.SensorReadings[0][0] {
		for i := range p.SensorReadings {
			for j := range p.SensorReadings[i] {
				p.SensorReadings[i][j][k] = averaged[i][j][k]
			}
		}
	}
	
}

//returns all of the nodes a radial distance from the current node
func NodesInRadius(curNode * NodeImpl, p * Params, radius int)(map[Tuple]*NodeImpl) {
	var gridMaxX = p.MaxX;
	var gridMaxY = p.MaxY;

	var nodesInRadius = map[Tuple]*NodeImpl{}

	var negRadius = -1*radius;

	//iterate over the Grid by row and column
	for row := negRadius; row<=radius; row++{
		for col := negRadius; col<=radius; col++{
			//do not include current node in list of nodes in radius
			if(row == 0 && col == 0){
				continue
			}

			var testX = curNode.X + col					//test X value
			var testY = curNode.Y + row					//test Y value
			var testTup = Tuple{testX, testY}	//create Tuple from test X and Y values
			if(testX < gridMaxX && testX >= 0){			//if the testX value is on the grid, continue
				if(testY < gridMaxY && testY >= 0){		//if the testY value is on the grid, continue
					if(NodePositionMap[testTup] != nil){	//if the test position has a Node, continue
						nodesInRadius[testTup] = NodePositionMap[testTup]	//add the node to the nodesInRadius map
					}
				}
			}
		}
	}
	return nodesInRadius
}

//returns all of the nodes dist squares away from the current node
func NodesWithinDistance(curNode * NodeImpl, p * Params, dist int)(map[Tuple]*NodeImpl){
	var gridMaxX = p.MaxX;
	var gridMaxY = p.MaxY;
	var nodesWithinDist = p.Grid[curNode.Y][curNode.X].NodesInSquare //initialize to nodes in current square
	var negDist = -1*dist;

	for row := negDist; row<=dist; row++ {
		for col := negDist; col <= dist; col++ {

			var testX = p.Grid[curNode.Y][curNode.X].X + col		//X value of test Square
			var testY = p.Grid[curNode.Y][curNode.X].Y + row		//Y value of test Square

			if(testX < gridMaxX && testX >= 0){			//if the testX value is on the grid, continue
				if(testY < gridMaxY && testY >= 0){		//if the testY value is on the grid, continue
					var testSquare =  p.Grid[testY][testX] 			//create Square from test X and Y values
					if(testSquare != nil){					//if the test Square is not null, continue
						for ind,val := range testSquare.NodesInSquare{	//iterate through nodes in square map adding each to the
							nodesWithinDist[ind] = val;					//nodes within Distance Map
						}
					}
				}
			}
		}
	}
	return nodesWithinDist
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
	start := time.Now()
	defer func() {
		elapsed := time.Since(start)
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
	//start := time.Now()
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
	//elapsed := time.Since(start)
	//fmt.Println("Printing Board Map took", elapsed)
}*/