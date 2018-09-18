package main

import "fmt"

/*
This is the server GO file and it is a model of our server
*/

//This is the server's data structure for a phone (or node)
type phoneFile struct {
	id int //This is the phone's unique id
	xPos []int //These are the saved x pos of the phone
	yPos []int //These are the saved y pos of the phone
	val []int //These are the saved values of the phone
	time []int //These are the saved times of the GPS/sensor readings
	bufferSizes []int //these are the saved buffer sizes when info was dumped to server
	speeds []int //these are the saved accelerometer based of the phone

	refined [][][][]int //x,y,val,time for all time
}

//The server is merely al list of phone files for now
type server struct {
	//p [numNodes]phoneFile
	p [200]phoneFile
}
//This is for later when the server becomes more advanced
type serverThink interface {
}

//This is the server absorbing data from the nodes and writing it to its phone files
func getData(s *server,xPos []int, yPos []int, val []int, time []int, id int, buffer int) () {
	//s.p[id].xPos = append(s.p[id].xPos,xPos ...)
	//s.p[id].yPos = append(s.p[id].yPos,yPos...)
	//s.p[id].val = append(s.p[id].val,val...)
	//s.p[id].time = append(s.p[id].time,time...)
	//s.p[id].bufferSizes = append(s.p[id].bufferSizes,buffer)
}

//This is a debugging function to be removed later
func (s server) String() {
	fmt.Println("Length of string",int(len(s.p))," ")
}

//This refines the phone files to fill in the gaps between where the server did not check the GPS or sensor
func reifne( p *phoneFile) (bool) {
	//This fills the positions
	if (len(p.yPos) == len(p.time)) == (len(p.yPos) == len(p.val)) {
		inbetween := 0
		open := false
		for i := 0; i < len(p.time); i++ {
			if p.xPos[i] == -1 && p.yPos[i] == -1 {
				inbetween += 1
			}
			if p.xPos[i] != -1 && p.yPos[i] != -1 && open == true && inbetween > 0 {
				diviserX := (p.xPos[i] - p.xPos[i-inbetween-1])/(inbetween+1)
				diviserY := (p.yPos[i] - p.yPos[i-inbetween-1])/(inbetween+1)
				for x := 0; x < inbetween; x++ {
					p.xPos[i-inbetween+x] = diviserX + p.xPos[i-inbetween+x-1]
					p.yPos[i-inbetween+x] = diviserY + p.yPos[i-inbetween+x-1]
				}
				inbetween = 0
			} else if p.xPos[i] != -1 && p.yPos[i] != -1 && open == false {
				open = true
				inbetween = 0
			}
		}
		inbetween = 0
		open = false
		//This fills the values
		for i := 1; i < len(p.time); i++ {
			if p.val[i] == -1 {
				inbetween += 1
			}
			if p.val[i] != -1 && p.val[i-1] == -1 && inbetween > 0 && open == true {
				diviserV := (p.val[i] - p.val[i-inbetween-1])/(inbetween+1)
				for x:= 0; x < inbetween; x++ {
					p.val[i-inbetween+x] = diviserV + p.val[i-inbetween+x-1]
				}
				inbetween = 0
			} else if p.val[i] != -1 && open == false {
				open = true
				inbetween = 0
			}
		}
		return true
	} else {
		return false
	}
}