/*
This is the server GO file and it is a model of our server
*/
package cps

import "fmt"

var (
	timeBuckets [1000][]float64 //2D array where each sub array is the sensor readings at one iteration
	mean [1000]float64
	stdDev [1000]float64
	variance [1000]float64
)

type FusionCenter struct {
	P *Params
}

type Reading struct {
	sensorVal 	float64
	xPos 		int
	yPos 		int
	time 		int //Time represented by iteration number
	id 			int //Node id number
}

func (s FusionCenter) GetGridAverage() float32 {
	return 0.0

}

func (s FusionCenter) Send(rd Reading) {
	fmt.Printf("Sending to server:\nX: %v, Y: %v, Sensor Value: %v\n", rd.xPos, rd.yPos, rd.sensorVal)
	currBucket := timeBuckets[rd.time]
	if currBucket != nil {
		timeBuckets[rd.time] = append(currBucket, rd.sensorVal)
	} else {
		timeBuckets[rd.time] = []float64{rd.sensorVal}
	}
}

func (s FusionCenter) CalcStats() ([1000]float64, [1000]float64, [1000]float64) {
	sum := 0.0
	for i := 0; i < len(timeBuckets); i++ {
		for j := 0; j < len(timeBuckets[i]); j++ {
			//fmt.Printf("Bucket size: %v\n", len(timeBuckets[i]))
			sum += timeBuckets[i][j]
			fmt.Printf("Time : %v, Elements #: %v, Value: %v\n", i, j, timeBuckets[i][j])
		}
		mean[i] = sum / float64( len(timeBuckets[i]) )
	}
	//mean = []float64{1.0}
	//stdDev = []float64{1.0}
	//variance = []float64{1.0}
	fmt.Printf("\nMean: %v", mean)
	return mean, stdDev, variance
}



/*
This is the server GO file and it is a model of our server
*/

//This is the server's data structure for a phone (or node)
//type PhoneFile struct {
//	id int //This is the phone's unique id
//	xPos []int //These are the saved x pos of the phone
//	yPos []int //These are the saved y pos of the phone
//	val []int //These are the saved values of the phone
//	time []int //These are the saved times of the GPS/sensor readings
//	bufferSizes []int //these are the saved buffer sizes when info was dumped to server
//	speeds []int //these are the saved accelerometer based of the phone
//
//	refined [][][][]int //x,y,val,time for all time
//}
//
////The server is merely al list of phone files for now
//type Server struct {
//	//p [numNodes]phoneFile
//	p [200]PhoneFile
//}
////This is for later when the server becomes more advanced
//type serverThink interface {
//}
//
////This is the server absorbing data from the nodes and writing it to its phone files
//func GetData(s *Server,xPos []int, yPos []int, val []int, time []int, id int, buffer int) () {
//	//s.p[id].xPos = append(s.p[id].xPos,xPos ...)
//	//s.p[id].yPos = append(s.p[id].yPos,yPos...)
//	//s.p[id].val = append(s.p[id].val,val...)
//	//s.p[id].time = append(s.p[id].time,time...)
//	//s.p[id].bufferSizes = append(s.p[id].bufferSizes,buffer)
//}
//
////This is a debugging function to be removed later
//func (s Server) String() {
//	fmt.Println("Length of string",int(len(s.p))," ")
//}
//
////This refines the phone files to fill in the gaps between where the server did not check the GPS or sensor
//func Refine( p *PhoneFile) (bool) {
//	//This fills the positions
//	if (len(p.yPos) == len(p.time)) == (len(p.yPos) == len(p.val)) {
//		inbetween := 0
//		open := false
//		for i := 0; i < len(p.time); i++ {
//			if p.xPos[i] == -1 && p.yPos[i] == -1 {
//				inbetween += 1
//			}
//			if p.xPos[i] != -1 && p.yPos[i] != -1 && open == true && inbetween > 0 {
//				diviserX := (p.xPos[i] - p.xPos[i-inbetween-1])/(inbetween+1)
//				diviserY := (p.yPos[i] - p.yPos[i-inbetween-1])/(inbetween+1)
//				for x := 0; x < inbetween; x++ {
//					p.xPos[i-inbetween+x] = diviserX + p.xPos[i-inbetween+x-1]
//					p.yPos[i-inbetween+x] = diviserY + p.yPos[i-inbetween+x-1]
//				}
//				inbetween = 0
//			} else if p.xPos[i] != -1 && p.yPos[i] != -1 && open == false {
//				open = true
//				inbetween = 0
//			}
//		}
//		inbetween = 0
//		open = false
//		//This fills the values
//		for i := 1; i < len(p.time); i++ {
//			if p.val[i] == -1 {
//				inbetween += 1
//			}
//			if p.val[i] != -1 && p.val[i-1] == -1 && inbetween > 0 && open == true {
//				diviserV := (p.val[i] - p.val[i-inbetween-1])/(inbetween+1)
//				for x:= 0; x < inbetween; x++ {
//					p.val[i-inbetween+x] = diviserV + p.val[i-inbetween+x-1]
//				}
//				inbetween = 0
//			} else if p.val[i] != -1 && open == false {
//				open = true
//				inbetween = 0
//			}
//		}
//		return true
//	} else {
//		return false
//	}
//}