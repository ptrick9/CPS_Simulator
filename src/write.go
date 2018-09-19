package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strconv"
	//"time"
)

var (
	//iterations_used int = 0

	// variables for making maps
	bufferCurrent = [][]int{{2, 2}, {0, 0}} // points currently being worked with
	bufferFuture  = [][]int{{}}             // point to be worked with
	starter       = 1                       // This is the destination number
	boardMap      = [][]int{                // This is the map with all the position variables for path finding
		{0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0}}

	wallPoints = [][]int{{1, 1}, {1, 2}, {1, 3}, {2, 1}}
	// end variables for making maps

	npos    [][]int // node positions
	wpos    [][]int // wall positions
	spos    [][]int // super node positions
	ppos    [][]int // super node points of interest positions
	poikpos [][]int // points of interest kinetic
	poispos [][]int // points of interest static

	//b []byte
	b2 []byte
	b3 []byte
	b4 []byte
	b5 []byte
	b6 []byte
	b7 []byte

	numNodeNodes               int
	numWallNodes               int
	numPoints                  int
	numPointsOfInterestKinetic int
	numPointsOfInterestStatic  int

	fileName = "Log1_in.txt"

	makeBoardMapFile = true
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

func getDashedInput(s string) int {
	b := readFromFile(fileName)
	r := regexp.MustCompile(string(s + "-[0-9]+"))
	w := r.FindAllString(string(b), 1)
	r = regexp.MustCompile("[0-9]+")
	w = r.FindAllString(w[0], 10)
	s1, err := strconv.ParseInt(w[len(w)-1], 10, 32)
	check(err)
	return int(s1)
}

func getListedInput() {
	b := readFromFile(fileName)
	r := regexp.MustCompile("N: [0-9]+")
	w := r.FindAllString(string(b), 10)
	r = regexp.MustCompile("[0-9]+")
	w = r.FindAllString(w[0], 10)
	s, err := strconv.ParseInt(w[0], 10, 32)
	check(err)
	r = regexp.MustCompile("x:[0-9]+, y:[0-9]+, t:[0-9]+")
	fai := r.FindAllIndex(b, int(s))
	w = r.FindAllString(string(b), int(s))
	if len(fai) > 0 {
		b2 = b[fai[len(fai)-1][1]:]
	} else {
		b2 = b
	}
	//fmt.Println(w)
	fillInts(w, 0)
	//fmt.Println(npos)

	r = regexp.MustCompile("W: [0-9]+")
	w = r.FindAllString(string(b2), 10)
	r = regexp.MustCompile("[0-9]+")
	w = r.FindAllString(w[0], 10)
	s, err = strconv.ParseInt(w[0], 10, 32)
	r = regexp.MustCompile("x:[0-9]+, y:[0-9]+")
	check(err)
	fai = r.FindAllIndex(b2, int(s))
	w = r.FindAllString(string(b2), int(s))
	if len(fai) > 0 {
		b3 = b2[fai[len(fai)-1][1]:]
	} else {
		b3 = b2
	}
	//fmt.Println(w)
	fillInts(w, 1)
	//fmt.Println(wpos)

	r = regexp.MustCompile("S: [0-9]+")
	w = r.FindAllString(string(b3), 10)
	r = regexp.MustCompile("[0-9]+")
	w = r.FindAllString(w[0], 10)
	s, err = strconv.ParseInt(w[0], 10, 32)
	r = regexp.MustCompile("x:[0-9]+, y:[0-9]+, t:[0-9]+")
	check(err)
	fai = r.FindAllIndex(b3, int(s))
	w = r.FindAllString(string(b3), int(s))
	if len(fai) > 0 {
		b4 = b3[fai[len(fai)-1][1]:]
	} else {
		b4 = b3
	}
	//fmt.Println(w)
	fillInts(w, 2)
	//fmt.Println(spos)

	/*r = regexp.MustCompile("P: [0-9]+")
	w = r.FindAllString(string(b4), 10)
	r = regexp.MustCompile("[0-9]+")
	w = r.FindAllString(w[0], 10)
	s, err = strconv.ParseInt(w[0], 10, 32)
	r = regexp.MustCompile("x:[0-9]+, y:[0-9]+, t:[0-9]+")
	check(err)
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
	check(err)
	fai = r.FindAllIndex(b5, int(s))
	w = r.FindAllString(string(b5), int(s))
	if len(fai) > 0 {
		b6 = b5[fai[len(fai)-1][1]:]
	} else {
		b6 = b5
	}
	//fmt.Println(w)
	fillInts(w, 5)
	//fmt.Println(poispos)

	numNodeNodes = len(npos)
	numWallNodes = len(wpos)
	//numPoints = len(ppos)
	numPointsOfInterestStatic = len(poispos)
	//fmt.Println(numNodeNodes, numWallNodes, numPoints, numPointsOfInterestKinetic, numPointsOfInterestStatic)
}

func fillInts(s []string, place int) {
	if place == 0 {
		for i := 0; i < len(s); i++ {
			r := regexp.MustCompile("x:[0-9]+")
			X := r.FindAllString(s[i], 1)
			r = regexp.MustCompile("[0-9]+")
			x := r.FindAllString(X[0], 1)
			s1, err := strconv.ParseInt(x[0], 10, 32)
			check(err)
			ap := []int{int(s1), 0, 0}
			npos = append(npos, ap)

			r = regexp.MustCompile("y:[0-9]+")
			X = r.FindAllString(s[i], 1)
			r = regexp.MustCompile("[0-9]+")
			x = r.FindAllString(X[0], 1)
			s1, err = strconv.ParseInt(x[0], 10, 32)
			check(err)
			npos[i][1] = int(s1)

			r = regexp.MustCompile("t:[0-9]+")
			X = r.FindAllString(s[i], 1)
			r = regexp.MustCompile("[0-9]+")
			x = r.FindAllString(X[0], 1)
			s1, err = strconv.ParseInt(x[0], 10, 32)
			check(err)
			npos[i][2] = int(s1)
		}
	} else if place == 1 {
		for i := 0; i < len(s); i++ {
			r := regexp.MustCompile("x:[0-9]+")
			X := r.FindAllString(s[i], 1)
			r = regexp.MustCompile("[0-9]+")
			x := r.FindAllString(X[0], 1)
			s1, err := strconv.ParseInt(x[0], 10, 32)
			check(err)
			ap := []int{int(s1), 0, 0}
			wpos = append(wpos, ap)

			r = regexp.MustCompile("y:[0-9]+")
			X = r.FindAllString(s[i], 1)
			r = regexp.MustCompile("[0-9]+")
			x = r.FindAllString(X[0], 1)
			s1, err = strconv.ParseInt(x[0], 10, 32)
			check(err)
			wpos[i][1] = int(s1)

		}
	} else if place == 2 {
		for i := 0; i < len(s); i++ {
			r := regexp.MustCompile("x:[0-9]+")
			X := r.FindAllString(s[i], 1)
			r = regexp.MustCompile("[0-9]+")
			x := r.FindAllString(X[0], 1)
			s1, err := strconv.ParseInt(x[0], 10, 32)
			check(err)
			ap := []int{int(s1), 0, 0}
			spos = append(spos, ap)

			r = regexp.MustCompile("y:[0-9]+")
			X = r.FindAllString(s[i], 1)
			r = regexp.MustCompile("[0-9]+")
			x = r.FindAllString(X[0], 1)
			s1, err = strconv.ParseInt(x[0], 10, 32)
			check(err)
			spos[i][1] = int(s1)

			r = regexp.MustCompile("t:[0-9]+")
			X = r.FindAllString(s[i], 1)
			r = regexp.MustCompile("[0-9]+")
			x = r.FindAllString(X[0], 1)
			s1, err = strconv.ParseInt(x[0], 10, 32)
			check(err)
			spos[i][2] = int(s1)
		}
	} else if place == 3 {
		for i := 0; i < len(s); i++ {
			r := regexp.MustCompile("x:[0-9]+")
			X := r.FindAllString(s[i], 1)
			r = regexp.MustCompile("[0-9]+")
			x := r.FindAllString(X[0], 1)
			s1, err := strconv.ParseInt(x[0], 10, 32)
			check(err)
			ap := []int{int(s1), 0, 0}
			ppos = append(ppos, ap)

			r = regexp.MustCompile("y:[0-9]+")
			X = r.FindAllString(s[i], 1)
			r = regexp.MustCompile("[0-9]+")
			x = r.FindAllString(X[0], 1)
			s1, err = strconv.ParseInt(x[0], 10, 32)
			check(err)
			ppos[i][1] = int(s1)

			r = regexp.MustCompile("t:[0-9]+")
			X = r.FindAllString(s[i], 1)
			r = regexp.MustCompile("[0-9]+")
			x = r.FindAllString(X[0], 1)
			s1, err = strconv.ParseInt(x[0], 10, 32)
			check(err)
			ppos[i][2] = int(s1)
		}
	} else if place == 5 {
		for i := 0; i < len(s); i++ {
			r := regexp.MustCompile("x:[0-9]+")
			X := r.FindAllString(s[i], 1)
			r = regexp.MustCompile("[0-9]+")
			x := r.FindAllString(X[0], 1)
			s1, err := strconv.ParseInt(x[0], 10, 32)
			check(err)
			ap := []int{int(s1), 0, 0, 0}
			poispos = append(poispos, ap)

			r = regexp.MustCompile("y:[0-9]+")
			X = r.FindAllString(s[i], 1)
			r = regexp.MustCompile("[0-9]+")
			x = r.FindAllString(X[0], 1)
			s1, err = strconv.ParseInt(x[0], 10, 32)
			check(err)
			poispos[i][1] = int(s1)

			r = regexp.MustCompile("ti:[0-9]+")
			X = r.FindAllString(s[i], 1)
			r = regexp.MustCompile("[0-9]+")
			x = r.FindAllString(X[0], 1)
			s1, err = strconv.ParseInt(x[0], 10, 32)
			check(err)
			poispos[i][2] = int(s1)

			r = regexp.MustCompile("to:[0-9]+")
			X = r.FindAllString(s[i], 1)
			r = regexp.MustCompile("[0-9]+")
			x = r.FindAllString(X[0], 1)
			s1, err = strconv.ParseInt(x[0], 10, 32)
			check(err)
			poispos[i][3] = int(s1)
		}
	}
}

// Returns the char number associated with a byte
func getIntFromByte(a byte) int {
	if a <= 57 && a >= 48 {
		return int(a - 48)
	} else {
		return -1
	}
}

// Returns the string character of a byte
func getLetterFromByte(a byte) string {
	return string([]byte{a})
}

// Clears file then writes message
func writeToFile(name string, message string) {
	d1 := []byte(message)
	err := ioutil.WriteFile(name, append(readFromFile(name), d1...), 0644)
	check(err)
}

// Reads entire file to array of bytes
func readFromFile(name string) (b []byte) {
	b, err := ioutil.ReadFile(name)
	check(err)
	return
}

// Creates a file file with specific name
func createFile(name string) {
	file, err := os.Create(name) // creates text file
	check(err)                   // checks if text file is created properly
	file.Close()                 // closes the file at the end
}

// Checks an error
func check(e error) {
	if e != nil {
		fmt.Println(e)
	}
}

// Creates boardMap
func createBoard(x int, y int) {
	boardMap = [][]int{}
	for i := 0; i < y; i++ {
		boardMap = append(boardMap, []int{})
		for j := 0; j < x; j++ {
			boardMap[i] = append(boardMap[i], 0)
		}
	}
}

// Fills the walls into the board based on the wall positions extrapolated from the file
func fillInWallsToBoard() {
	for i := 0; i < len(wpos); i++ {
		boardMap[wpos[i][1]][wpos[i][0]] = -1
	}
}

// Fills the points of interest into the current buffer

func fillInBufferCurrent() {
	bufferCurrent = [][]int{}
	for i := 0; i < len(poispos); i++ {
		if iterations_used >= poispos[i][2] && iterations_used < poispos[i][3] {
			bufferCurrent = append(bufferCurrent, []int{poispos[i][1], poispos[i][0]})
			//fmt.Println("1ho- ", iterations_used, "2ho", bufferCurrent)
		}
	}
}

// Fills the points of interest to the board
func fillPointsToBoard() {
	for i := 0; i < len(bufferCurrent); i++ {
		//fmt.Println(bufferCurrent)
		boardMap[bufferCurrent[i][0]][bufferCurrent[i][1]] = starter
	}
}

// Fills in board map with the path finding values
func fillInMap1() {
	/*start := time.Now()

	defer func() {
		elapsed := time.Since(start)
		//fmt.Println("Board Map took", elapsed)
	}()*/

	for len(bufferFuture) > 0 {
		bufferFuture = [][]int{}
		for i := 0; i < len(bufferCurrent); i++ {
			checkLeft(i)
			checkRight(i)
			checkUp(i)
			checkDown(i)
		}
		bufferCurrent = [][]int{}
		bufferCurrent = append(bufferCurrent, bufferFuture...)
		starter += 1
	}
	starter = 1
	bufferFuture = [][]int{{}}
}

func checkLeft(i int) {
	defer func() {
		recover()
	}()
	if
	//bufferCurrent[i][0]-1 < len(boardMap) &&
	//bufferCurrent[i][1] < len(boardMap[1]) &&
	//bufferCurrent[i][0]-1 >= 0 &&
	//bufferCurrent[i][1] >= 0 &&
	boardMap[bufferCurrent[i][0]-1][bufferCurrent[i][1]] == 0 {

		boardMap[bufferCurrent[i][0]-1][bufferCurrent[i][1]] = starter + 1
		bufferFuture = append(bufferFuture, []int{bufferCurrent[i][0] - 1, bufferCurrent[i][1]})
	}
}

func checkRight(i int) {
	defer func() {
		recover()
	}()
	if
	//bufferCurrent[i][0]+1 < len(boardMap) &&
	//bufferCurrent[i][1] < len(boardMap[1]) && // boardMap[1] to
	//bufferCurrent[i][0]+1 >= 0 &&
	//bufferCurrent[i][1] >= 0 &&
	boardMap[bufferCurrent[i][0]+1][bufferCurrent[i][1]] == 0 {

		boardMap[bufferCurrent[i][0]+1][bufferCurrent[i][1]] = starter + 1
		bufferFuture = append(bufferFuture, []int{bufferCurrent[i][0] + 1, bufferCurrent[i][1]})
	}
}
func checkUp(i int) {
	defer func() {
		recover()
	}()
	if
	//bufferCurrent[i][0] < len(boardMap) &&
	//bufferCurrent[i][1]+1 < len(boardMap[1]) &&
	//bufferCurrent[i][0] >= 0 &&
	//bufferCurrent[i][1]+1 >= 0 &&
	boardMap[bufferCurrent[i][0]][bufferCurrent[i][1]+1] == 0 {

		boardMap[bufferCurrent[i][0]][bufferCurrent[i][1]+1] = starter + 1
		bufferFuture = append(bufferFuture, []int{bufferCurrent[i][0], bufferCurrent[i][1] + 1})
	}
}
func checkDown(i int) {
	defer func() {
		recover()
	}()
	if
	//bufferCurrent[i][0] < len(boardMap) &&
	//bufferCurrent[i][1]-1 < len(boardMap[1]) &&
	//bufferCurrent[i][0] >= 0 &&
	//bufferCurrent[i][1]-1 >= 0 &&
	boardMap[bufferCurrent[i][0]][bufferCurrent[i][1]-1] == 0 {

		boardMap[bufferCurrent[i][0]][bufferCurrent[i][1]-1] = starter + 1
		bufferFuture = append(bufferFuture, []int{bufferCurrent[i][0], bufferCurrent[i][1] - 1})
	}
}

func fillInMap() {
	/*start := time.Now()

	defer func() {
		elapsed := time.Since(start)
		//fmt.Println("Board Map took", elapsed)
	}()*/

	for len(bufferFuture) > 0 {
		bufferFuture = [][]int{}
		for i := 0; i < len(bufferCurrent); i++ {
			// empty buffer future
			//check above
			//fmt.Println(len(boardMap[1]),i)
			if bufferCurrent[i][0]-1 < len(boardMap) &&
				bufferCurrent[i][1] < len(boardMap[1]) &&
				bufferCurrent[i][0]-1 >= 0 &&
				bufferCurrent[i][1] >= 0 &&
				boardMap[bufferCurrent[i][0]-1][bufferCurrent[i][1]] == 0 {

				boardMap[bufferCurrent[i][0]-1][bufferCurrent[i][1]] = starter + 1
				bufferFuture = append(bufferFuture, []int{bufferCurrent[i][0] - 1, bufferCurrent[i][1]})
			}
			//check below
			if bufferCurrent[i][0]+1 < len(boardMap) &&
				bufferCurrent[i][1] < len(boardMap[1]) && // boardMap[1] to
				bufferCurrent[i][0]+1 >= 0 &&
				bufferCurrent[i][1] >= 0 &&
				boardMap[bufferCurrent[i][0]+1][bufferCurrent[i][1]] == 0 {

				boardMap[bufferCurrent[i][0]+1][bufferCurrent[i][1]] = starter + 1
				bufferFuture = append(bufferFuture, []int{bufferCurrent[i][0] + 1, bufferCurrent[i][1]})
			}
			// to the rite
			if bufferCurrent[i][0] < len(boardMap) &&
				bufferCurrent[i][1]+1 < len(boardMap[1]) &&
				bufferCurrent[i][0] >= 0 &&
				bufferCurrent[i][1]+1 >= 0 &&
				boardMap[bufferCurrent[i][0]][bufferCurrent[i][1]+1] == 0 {

				boardMap[bufferCurrent[i][0]][bufferCurrent[i][1]+1] = starter + 1
				bufferFuture = append(bufferFuture, []int{bufferCurrent[i][0], bufferCurrent[i][1] + 1})
			}
			// check to the left
			if bufferCurrent[i][0] < len(boardMap) &&
				bufferCurrent[i][1]-1 < len(boardMap[1]) &&
				bufferCurrent[i][0] >= 0 &&
				bufferCurrent[i][1]-1 >= 0 &&
				boardMap[bufferCurrent[i][0]][bufferCurrent[i][1]-1] == 0 {

				boardMap[bufferCurrent[i][0]][bufferCurrent[i][1]-1] = starter + 1
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

// This prints the board map to the terminal
//func printBoardMap(){
//	for i:= 0; i < len(boardMap); i++{
//		fmt.Println(boardMap[i])
//	}
//}

var fileBoard, errBoard = os.Create("boardMap.txt")

// This prints board Map to a txt file.
func writeBordMapToFile2() {
	/*start := time.Now()
	defer func() {
		elapsed := time.Since(start)
		fmt.Println("Printing Board Map took", elapsed)
	}()*/
	check(errBoard)
	var s = ""
	s = s + "\nt=" + strconv.Itoa(iterations_used) + "\n\n"
	for i := 0; i < len(boardMap); i++ {
		for j := 0; j < len(boardMap[i]); j++ {
			s = s + strconv.Itoa(boardMap[i][j]) + " "
		}
		s = s + "\n"
	}
	//s = s + "\nt=" + strconv.Itoa(iterations_used) + "\n\n"
	n3, err := fileBoard.WriteString(s)
	check(err)
	fmt.Printf("wrote %d bytes\n", n3)
}

func writeBordMapToFile() {
	//start := time.Now()
	check(errBoard)
	w := bufio.NewWriter(fileBoard)
	w.WriteString("\nt=" + strconv.Itoa(iterations_used) + "\n\n")
	for i := 0; i < len(boardMap); i++ {
		for j := 0; j < len(boardMap[i]); j++ {
			w.WriteString(strconv.Itoa(boardMap[i][j]) + " ")
		}
		w.WriteString("\n")
	}
	w.Flush()
	//elapsed := time.Since(start)
	//fmt.Println("Printing Board Map took", elapsed)
}
