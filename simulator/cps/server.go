/*
This is the server GO file and it is a model of our server
*/
package cps

import (
	"fmt"
	"math"
)

/*var (
	s.TimeBuckets [1000][]float64 //2D array where each sub array is the sensor readings at one iteration
	s.Mean [1000]float64
	s.StdDev [1000]float64
	s.Variance [1000]float64
)*/

type FusionCenter struct {
	P *Params

	TimeBuckets 	[][]float64 //2D array where each sub array is the sensor readings at one iteration
	Mean 			[]float64
	StdDev 			[]float64
	Variance 		[]float64
	//TimesInPacket mapset.Set//Set of times in received packet
	Times 			map[int]bool
}

func (s *FusionCenter) Init(){
	s.TimeBuckets = make([][]float64, s.P.Iterations_used)
	s.Mean = make([]float64, s.P.Iterations_used)
	s.StdDev = make([]float64, s.P.Iterations_used)
	s.Variance = make([]float64, s.P.Iterations_used)

	//s.TimesInPacket = mapset.NewSet()
	s.Times = make(map[int]bool, 0)
}

//Data sent by node
type Reading struct {
	sensorVal 	float64
	xPos 		int
	yPos 		int
	time 		int //Time represented by iteration number
	id 			int //Node id number
}

func (s FusionCenter) GetSquareAverage(tile *Square) float32 {
	return tile.Avg
}

func (s FusionCenter) UpdateSquareAvg(rd Reading) {
	//var curNode NodeImpl
	tile := s.P.Grid[rd.yPos/s.P.YDiv][rd.xPos/s.P.XDiv]
	tile.TakeMeasurement(float32(rd.sensorVal))

}

//Searches node list and updates the number of nodes in each square
func (s FusionCenter) UpdateSquareNumNodes() {
	var node NodeImpl

	//Clear number of nodes for each square
	for i:=0; i < len(s.P.Grid); i++ {
		for j:=0; j < len(s.P.Grid[i]); j++ {
			s.P.Grid[i][j].ActualNumNodes = 0
		}
	}

	//Count number of nodes in each square
	for i:=0; i < s.P.NumNodes; i++ {
		node = s.P.NodeList[i]
		s.P.Grid[node.Y/s.P.YDiv][node.X/s.P.XDiv].ActualNumNodes += 1
	}

	//Debugging: Check if total nodes after update is equal to total expected nodes
	/*totalNodes := 0
	for i:=0; i < len(s.P.Grid); i++ {
		for j:=0; j < len(s.P.Grid[0]); j++ {
			totalNodes += s.P.Grid[i][j].ActualNumNodes
		}
	}
	if totalNodes != s.P.NumNodes {
		fmt.Printf("Error, number of nodes do not match! Found %v nodes out of %v\n", totalNodes, s.P.NumNodeNodes)
	}*/
}

//Node calls this to send data to server. Statistics are calculated each time data is received
func (s *FusionCenter) Send(rd Reading) {
	s.Times = make(map[int]bool, 0)
	if s.Times[rd.time] {

	} else {
		s.Times[rd.time] = true
	}
	//fmt.Printf("Sending to server:\nTime: %v, ID: %v, X: %v, Y: %v, Sensor Value: %v\n", rd.time, rd.id, rd.xPos, rd.yPos, rd.sensorVal)
	if len(s.TimeBuckets) <= rd.time {
		s.TimeBuckets = append(s.TimeBuckets, make([]float64,0))
	}
	currBucket := (s.TimeBuckets)[rd.time]
	if len(currBucket) != 0 { //currBucket != nil
		(s.TimeBuckets)[rd.time] = append(currBucket, rd.sensorVal)
	} else {
		(s.TimeBuckets)[rd.time] = append((s.TimeBuckets)[rd.time], rd.sensorVal) //s.TimeBuckets[rd.time] = []float64{rd.sensorVal}
	}

	s.UpdateSquareAvg(rd)
	tile := s.P.Grid[rd.yPos/s.P.YDiv][rd.xPos/s.P.XDiv]
	tile.SquareValues += math.Pow(float64(rd.sensorVal-float64(tile.Avg)), 2)
	//s.CalcStats()
}

//Calculates s.Mean, standard deviation and s.Variance
func (s *FusionCenter) CalcStats() ([]float64, []float64, []float64) {
	//fmt.Printf("Calculating stats for times: %v", s.times)
	//Calculate the mean
	sum := 0.0
	for i:= range s.Times {
		for j := 0; j < len(s.TimeBuckets[i]); j++ {
			//fmt.Printf("Bucket size: %v\n", len(s.TimeBuckets[i]))
			sum += (s.TimeBuckets)[i][j]
			//fmt.Printf("Time : %v, Elements #: %v, Value: %v\n", i, j, s.TimeBuckets[i][j])
		}
		if len(s.Mean) <= i {
			s.Mean = append(s.Mean, sum / float64( len(s.TimeBuckets[i]) ))
		} else {
			s.Mean[i] = sum / float64(len(s.TimeBuckets[i]))
		}
	}

	//Calculate the standard deviation and variance
	sum = 0.0
	for i:= range s.Times {
		for j := 0; j < len((s.TimeBuckets)[i]); j++ {
			sum += math.Pow((s.TimeBuckets)[i][j] - s.Mean[i], 2)
		}

		if len(s.Variance) <= i {
			s.Variance = append(s.Variance, sum / float64( len(s.TimeBuckets[i]) ))
		} else {
			s.Variance[i] = sum / float64(len(s.TimeBuckets[i]))
		}

		if len(s.StdDev) <= i {
			s.StdDev = append(s.StdDev, math.Sqrt(sum / float64( len((s.TimeBuckets)[i])) ))
		} else {
			s.StdDev[i] = math.Sqrt(sum / float64( len((s.TimeBuckets)[i])) )
		}

		//s.Variance[i] = sum / float64( len((s.TimeBuckets)[i]) )
		//s.StdDev[i] = math.Sqrt(sum / float64( len((s.TimeBuckets)[i])) )
	}

	//sum := 0.0
	//Calculate the Mean
	/*for i := 0; i < len(s.TimeBuckets); i++ {
		for j := 0; j < len(s.TimeBuckets[i]); j++ {
			//fmt.Printf("Bucket size: %v\n", len(s.TimeBuckets[i]))
			sum += (s.TimeBuckets)[i][j]
			//fmt.Printf("Time : %v, Elements #: %v, Value: %v\n", i, j, s.TimeBuckets[i][j])
		}
		if len(s.Mean) <= i {
			s.Mean = append(s.Mean, sum / float64( len(s.TimeBuckets[i]) ))
		} else {
			s.Mean[i] = sum / float64(len(s.TimeBuckets[i]))
		}
	}

	//Calculate Standard Deviation and Variance
	sum = 0.0
	for i:= 0; i < len(s.TimeBuckets); i++ {
		for j := 0; j < len((s.TimeBuckets)[i]); j++ {
			sum += math.Pow((s.TimeBuckets)[i][j] - s.Mean[i], 2)
		}

		if len(s.Variance) <= i {
			s.Variance = append(s.Variance, sum / float64( len(s.TimeBuckets[i]) ))
		} else {
			s.Variance[i] = sum / float64(len(s.TimeBuckets[i]))
		}

		if len(s.StdDev) <= i {
			s.StdDev = append(s.StdDev, math.Sqrt(sum / float64( len((s.TimeBuckets)[i])) ))
		} else {
			s.StdDev[i] = math.Sqrt(sum / float64( len((s.TimeBuckets)[i])) )
		}

		//s.Variance[i] = sum / float64( len((s.TimeBuckets)[i]) )
		//s.StdDev[i] = math.Sqrt(sum / float64( len((s.TimeBuckets)[i])) )
	}*/
	return s.Mean, s.StdDev, s.Variance
}

func (s FusionCenter) PrintStats() {
	for i:= 0; i < s.P.Iterations_used; i++ {
		fmt.Printf("Time: %v, Mean: %v, Std Deviation: %v, Variance: %v\n", i, s.Mean[i], s.StdDev[i], s.Variance[i])
	}
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
//	if (len(s.P.yPos) == len(p.time)) == (len(p.yPos) == len(p.val)) {
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