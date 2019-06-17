/*
This is the server GO file and it is a model of our server
*/
package cps

import (
	"fmt"
	"math"
	"sort"
)

/*var (
	s.TimeBuckets [1000][]float64 //2D array where each sub array is the sensor readings at one iteration
	s.Mean [1000]float64
	s.StdDev [1000]float64
	s.Variance [1000]float64
)*/

type FusionCenter struct {
	P *Params

	TimeBuckets 	[][]Reading //2D array where each sub array is made of the readings at one iteration
	Mean 			[]float64
	StdDev 			[]float64
	Variance 		[]float64
	Times 			map[int]bool
}

func (s *FusionCenter) Init(){
	s.TimeBuckets = make([][]Reading, s.P.Iterations_used)
	s.Mean = make([]float64, s.P.Iterations_used)
	s.StdDev = make([]float64, s.P.Iterations_used)
	s.Variance = make([]float64, s.P.Iterations_used)
	s.Times = make(map[int]bool, 0)
}

//Data sent by node
type Reading struct {
	SensorVal float64
	Xpos      int
	YPos      int
	Time      int //Time represented by iteration number
	Id        int //Node Id number
	//StdDevFromMean	float64
}

func (s FusionCenter) GetSquareAverage(tile *Square) float32 {
	return tile.Avg
}

func (s FusionCenter) UpdateSquareAvg(rd Reading) {
	tile := s.P.Grid[rd.YPos/s.P.YDiv][rd.Xpos/s.P.XDiv]
	tile.TakeMeasurement(float32(rd.SensorVal))

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

//Node calls this to send data to server. Statistics are calculated each Time data is received
func (s *FusionCenter) Send(rd Reading) {
	//fmt.Printf("Sending to server:\nTime: %v, ID: %v, X: %v, Y: %v, Sensor Value: %v\n", rd.Time, rd.Id, rd.Xpos, rd.yPos, rd.sensorVal)
	s.Times = make(map[int]bool, 0)
	if s.Times[rd.Time] {

	} else {
		s.Times[rd.Time] = true
	}

	if len(s.TimeBuckets) <= rd.Time {
		s.TimeBuckets = append(s.TimeBuckets, make([]Reading,0))
	}
	currBucket := (s.TimeBuckets)[rd.Time]
	if len(currBucket) != 0 { //currBucket != nil
		(s.TimeBuckets)[rd.Time] = append(currBucket, rd)
	} else {
		(s.TimeBuckets)[rd.Time] = append((s.TimeBuckets)[rd.Time], rd) //s.TimeBuckets[rd.Time] = []float64{rd.sensorVal}
	}

	s.UpdateSquareAvg(rd)
	tile := s.P.Grid[rd.YPos/s.P.YDiv][rd.Xpos/s.P.XDiv]
	tile.SquareValues += math.Pow(float64(rd.SensorVal-float64(tile.Avg)), 2)
}

//Calculates s.Mean, standard deviation and s.Variance
func (s *FusionCenter) CalcStats() ([]float64, []float64, []float64) {
	//fmt.Printf("Calculating stats for times: %v", s.times)
	s.UpdateSquareNumNodes()

	//Calculate the mean
	sum := 0.0
	StdDevFromMean := 0.0
	for i:= range s.Times {
		for j := 0; j < len(s.TimeBuckets[i]); j++ {
			//fmt.Printf("Bucket size: %v\n", len(s.TimeBuckets[i]))
			sum += (s.TimeBuckets)[i][j].SensorVal
			//fmt.Printf("Time : %v, Elements #: %v, Value: %v\n", i, j, s.TimeBuckets[i][j])
		}
		if len(s.Mean) <= i {
			s.Mean = append(s.Mean, sum / float64( len(s.TimeBuckets[i]) ))
		} else {
			s.Mean[i] = sum / float64(len(s.TimeBuckets[i]))
		}
		sum = 0
	}

	//Calculate the standard deviation and variance
	sum = 0.0
	for i:= range s.Times {
		for j := 0; j < len((s.TimeBuckets)[i]); j++ {
			sum += math.Pow((s.TimeBuckets)[i][j].SensorVal - s.Mean[i], 2)
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
		//Determine how many std deviations data is away from mean
		for j:= range s.TimeBuckets[i] {
			StdDevFromMean = (s.TimeBuckets[i][j].SensorVal - s.Mean[i]) / s.StdDev[i]
			if StdDevFromMean > 4 || StdDevFromMean < -4{
				//fmt.Printf("Potential detection by node %v at X:%v, Y:%v with reading %v\n", s.TimeBuckets[i][j].Id, s.TimeBuckets[i][j].Xpos, s.TimeBuckets[i][j].YPos, s.TimeBuckets[i][j].SensorVal)
			}
		}
		sum = 0
	}
	return s.Mean, s.StdDev, s.Variance
}

//Gets median from data set and returns both the median and closest index to median
func (s FusionCenter) getMedian(arr []float64) float64{
	sort.Float64s(arr)
	size := 0.0
	median := 0.0
	size = float64(len(arr))
	//index := 0
	//Check if even
	if int(size) % 2 == 0 {
		median = (arr[int(size / 2.0)] + arr[int(size / 2.0 + 1)] ) / 2
	} else {
		median = arr[int(size / 2.0 + 0.5)]
	}
	return median
}

func (s FusionCenter) PrintStats() {
	for i:= 0; i < s.P.Iterations_used; i++ {
		fmt.Printf("Time: %v, Mean: %v, Std Deviation: %v, Variance: %v\n", i, s.Mean[i], s.StdDev[i], s.Variance[i])
	}
}

func (s FusionCenter) PrintStatsFile() {
	fmt.Fprintln(s.P.ServerFile, "Mean at each time:\n", s.P.Server.Mean)
	fmt.Fprintln(s.P.ServerFile, "Standard Deviations at each time:\n", s.P.Server.StdDev)
	fmt.Fprintln(s.P.ServerFile, "Variance at each time:\n", s.P.Server.Variance)
}



/*
This is the OLD server GO file and it is a model of our server
*/

//This is the server's data structure for a phone (or node)
//type PhoneFile struct {
//	Id int //This is the phone's unique Id
//	Xpos []int //These are the saved x pos of the phone
//	yPos []int //These are the saved y pos of the phone
//	val []int //These are the saved values of the phone
//	Time []int //These are the saved times of the GPS/sensor readings
//	bufferSizes []int //these are the saved buffer sizes when info was dumped to server
//	speeds []int //these are the saved accelerometer based of the phone
//
//	refined [][][][]int //x,y,val,Time for all Time
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
//func GetData(s *Server,Xpos []int, yPos []int, val []int, Time []int, Id int, buffer int) () {
//	//s.p[Id].Xpos = append(s.p[Id].Xpos,Xpos ...)
//	//s.p[Id].yPos = append(s.p[Id].yPos,yPos...)
//	//s.p[Id].val = append(s.p[Id].val,val...)
//	//s.p[Id].Time = append(s.p[Id].Time,Time...)
//	//s.p[Id].bufferSizes = append(s.p[Id].bufferSizes,buffer)
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
//	if (len(s.P.yPos) == len(p.Time)) == (len(p.yPos) == len(p.val)) {
//		inbetween := 0
//		open := false
//		for i := 0; i < len(p.Time); i++ {
//			if p.Xpos[i] == -1 && p.yPos[i] == -1 {
//				inbetween += 1
//			}
//			if p.Xpos[i] != -1 && p.yPos[i] != -1 && open == true && inbetween > 0 {
//				diviserX := (p.Xpos[i] - p.Xpos[i-inbetween-1])/(inbetween+1)
//				diviserY := (p.yPos[i] - p.yPos[i-inbetween-1])/(inbetween+1)
//				for x := 0; x < inbetween; x++ {
//					p.Xpos[i-inbetween+x] = diviserX + p.Xpos[i-inbetween+x-1]
//					p.yPos[i-inbetween+x] = diviserY + p.yPos[i-inbetween+x-1]
//				}
//				inbetween = 0
//			} else if p.Xpos[i] != -1 && p.yPos[i] != -1 && open == false {
//				open = true
//				inbetween = 0
//			}
//		}
//		inbetween = 0
//		open = false
//		//This fills the values
//		for i := 1; i < len(p.Time); i++ {
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