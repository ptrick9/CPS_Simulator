package cps

import (
	"fmt"
	"math"
	"math/rand"
	"sort"
)

//Global variables used in battery loss dynamics
var (
	naturalLoss float32
)

//The NodeParent interface is inherited by all node types
type NodeParent interface {
	Distance(b Bomb) float32        //distance to bomb in the form of the node's reading
	Row(div int) int                //Row of node
	Col(div int) int                //Column of node
	GetSpeed() []float32            //History of accelerometer based speeds of node
	//BatteryLossDynamic()   //Battery loss based of ratios of battery usage
	//BatteryLossDynamic1()  //2 stage buffer battery loss
	UpdateHistory(newValue float32) //updates history of node's samples
	IncrementTotalSamples()         //increments total number of samples node has taken
	GetAvg() float32                //returns average of node's past samples
	IncrementNumResets()            //increments the number of times a node has been reset
	SetConcentration(conc float64)  //sets the concentration of a node
	GeoDist(b Bomb) float32         //returns distance from bomb (rather than reading of node)
	GetID() int                     //returns ID of node
	GetLoc() (x, y float32)             //returns x and y values of node

	//following functions set drifting parameters of nodes
	SetS0(s0 float64)
	SetS1(s1 float64)
	SetS2(s2 float64)
	SetE0(e0 float64)
	SetE1(e1 float64)
	SetE2(e2 float64)
	SetET1(et1 float64)
	SetET2(et2 float64)
	GetParams() (float64, float64, float64, float64, float64, float64, float64, float64) //returns all of the above parameters
	GetCoefficients() (float64, float64, float64)                                        //returns some of the above parameters
	GetX() int                                                                           //returns x position of node
	GetY() int                                                                           //returns y position of node
}

//NodeImpl is a struct that implements all the methods listed
//	above in NodeParent
type NodeImpl struct {
	P 								*Params
	Id                              int      //Id of node
	OldX                            int      // for movement
	OldY                            int      // for movement
	Sitting                         int      // for movement
	X                               float32      //x pos of node
	Y                               float32      //y pos of node
	Battery                         float32  //battery of node
	BatteryLossScalar               float32  //natural incremental battery loss of node
	BatteryLossSensor				float32  //sensor based battery loss of node
	BatteryLossGPS		            float32  //GPS based battery loss of node
	BatteryLossServer				float32  //server communication based battery loss of node

	BatteryLossBT					float32
	BatteryLossWifi					float32
	BatteryLoss4G					float32
	BatteryLossAccelerometer		float32

	ToggleCheckIterator             int      //node's personal iterator mostly for cascading pings
	//HasCheckedSensor                bool     //did the node just ping the sensor?
	TotalChecksSensor               int      //total sensor pings of node
	//HasCheckedGPS                   bool     //did the node just ping the GPS?
	TotalChecksGPS                  int      //total GPS pings of node
	//HasCheckedServer                bool     //did the node just communicate with the server?
	TotalChecksServer               int      //how many times did the node communicate with the server?
	PingPeriod                      float32  //This is the aggregate ping period used in some ping rate determining algorithms
	SensorPingPeriod                float32  //This is the ping period for the sensor
	GPSPingPeriod                   float32  //This is the ping period for the GPS
	ServerPingPeriod                float32  //This is the ping period for the server
	Pings                           float32  //This is an aggregate pings used in some ping rate determining algorithms
	SensorPings                     float32  //This is the total sensor pings to be made
	GPSPings                        float32  //This is the total GPS pings to be made
	ServerPings                     float32  //This is the total server pings to be made
	Cascade                         int      //This cascades the pings of the nodes
	BufferI                         int      //This is to keep track of the node's buffer size
	XPos                            [100]float32 //x pos buffer of node
	YPos                            [100]float32 //y pos buffer of node
	Value                           [100]int //Value buffer of node
	AccelerometerSpeedServer        [100]int //Accelerometer speed history of node
	Time                            [100]int //This keeps track of when specific pings are made
	//speedGPSPeriod int //This is a special period for speed based GPS pings but it is not used and may never be
	AccelerometerPosition [2][3]int //This is the accelerometer model of node
	AccelerometerSpeed    []float32 //History of accelerometer speeds recorded
	InverseSensor         float32   //Algorithm place holder declared here for speed
	InverseGPS            float32   //Algorithm place holder declared here for speed
	InverseServer         float32   //Algorithm place holder declared here for speed
	SampleHistory         []float32 //a history of the node's readings
	Avg                   float32   //weighted average of the node's most recent readings
	TotalSamples          int       //total number of samples taken by a node
	SpeedWeight           float32   //weight given to averaging of node's samples, based on node's speed
	NumResets             int       //number of times a node has had to reset due to drifting
	Concentration         float64   //used to determine reading of node
	SpeedGPSPeriod        int

	Current  int
	Previous int
	Diffx    int
	Diffy    int
	Speed    float32

	//The following values are all various drifting parameters of the node
	NewX               int
	NewY               int
	S0                 float64
	S1                 float64
	S2                 float64
	E0                 float64
	E1                 float64
	E2                 float64
	ET1                float64
	ET2                float64
	NodeTime           int
	Sensitivity        float64
	InitialSensitivity float64
	Valid 			   bool

	//allReadings 	   [1000]float64
	calibrateTimes 	   []int
	calibrateReading   []float64

	BatteryOverTime	   map[int]float32

	TotalPacketsSent    int
	TotalBytesSent		int

	IsClusterMember		bool
	IsClusterHead		bool

	NodeBounds			*Bounds //the node's representative bounds object
	NodeClusterParams	*ClusterMemberParams
	Recalibrated 		bool
	CHPenalty			float64
	ReadingsBuffer		[]Reading
	TimeLastSentReadings	int
}

//NodeMovement controls the movement of all the normal nodes
//It inherits all the methods and attributes from NodeParent
//	and NodeImpl
type NodeMovement interface {
	NodeParent
	Move(p *Params)
}

//Bouncing nodes bound around the grid
type Bn struct {
	*NodeImpl
	X_speed int
	Y_speed int
}

//Wall nodes go in a straight line from top/bottom or
//	side/side
type Wn struct {
	*NodeImpl
	Speed int
	Dir   int
}

//Random nodes get assigned a random x, y velocity every
//	move update
type Rn struct {
	*NodeImpl
}

type WallNodes struct {
	Node *NodeImpl
}

//Coord is a struct that contains x and y coordinates of
//	a square in the grid
//This struct is used by the super node type to create its
//	route through the grid
type Coord struct {
	Parent      *Coord
	X, Y        int
	Time        int
	G, H, Score int
}

//Path is a struct that contains an x and y integer and
//	a float for distance
//This struct is used when calculating the distance between
//	points of interest on the grid during super node route
//	scheduling
type Path struct {
	X, Y int
	Dist float64
}

//Returns the x Index of the square in which the specific
//	node currently resides
func (curNode *NodeImpl) Row(div int) int {
	return int(curNode.Y) / div
}

//Returns the y Index of the square in which the specific
//	node currently resides
func (curNode *NodeImpl) Col(div int) int {
	return int(curNode.X) / div
}

func (curNode *NodeImpl) InBounds(p *Params) bool {
	if int(curNode.X) < curNode.P.Width && int(curNode.X) >= 0 {
		if int(curNode.Y) < curNode.P.Height && curNode.Y >= 0 {
			return true
		}
	}
	return false
}

func (curNode *NodeImpl) TurnValid(x, y int, p *Params) bool {
	if x < curNode.P.Width && x >= 0 {
		if y < curNode.P.Height && y >= 0 {
			return true
		}
	}
	return false
}



func (curNode *NodeImpl) ADCReading(raw float32) int {

	level := (raw - curNode.P.ADCOffset)/curNode.P.ADCWidth

	if level > curNode.P.MaxADC {
		level = curNode.P.MaxADC
	} else if level < 0 {
		level = 0
	}

	return int(level)
}

func (curNode NodeImpl) String() string {
	//return fmt.Sprintf("x: %v y: %v Id: %v battery: %v sensor checked: %v sensor checks: %v GPS checked: %v GPS checks: %v server checked: %v server checks: %v buffer: %v ", int(curNode.X), curNode.Y, curNode.Id, curNode.Battery, curNode.HasCheckedSensor, curNode.TotalChecksSensor, curNode.HasCheckedGPS, curNode.TotalChecksGPS, curNode.HasCheckedServer, curNode.TotalChecksServer,curNode.BufferI)
	//return fmt.Sprintf("x: %v y: %v valid: %v", int(curNode.X), int(curNode.Y), curNode.Valid)
	//return fmt.Sprintf("battery: %v sensor checked: %v GPS checked: %v ", int(curNode.Battery), curNode.HasCheckedSensor, curNode.HasCheckedGPS)
	return fmt.Sprintf("battery: %v sensor checked: %v GPS checked: %v ", int(curNode.Battery), true, true)

}

func (c Coord) String() string {
	return fmt.Sprintf("{%v %v %v}", c.X, c.Y, c.Time)
}

func (c Coord) Equals(c2 Coord) bool {
	return c.X == c2.X && c.Y == c2.Y
}

func (curNode *NodeImpl) Move(p *Params) {
	if curNode.Sitting <= curNode.P.SittingStopThresholdCM {
		curNode.OldX = int(curNode.X) / curNode.P.XDiv
		curNode.OldY = int(curNode.Y) / curNode.P.YDiv

		var potentialSpots []GridSpot

		//only add the ones that are valid to move to into the list
		if int(curNode.Y)-1 >= 0 &&
			int(curNode.X) >= 0 &&
			int(curNode.X) < curNode.P.Width &&
			int(curNode.Y)-1 < curNode.P.Height &&

			curNode.P.BoardMap[int(curNode.X)][int(curNode.Y)-1] != -1 &&
			curNode.P.BoolGrid[int(curNode.X)][int(curNode.Y)-1] == false { // &&
			//curNode.P.BoardMap[int(curNode.X)][curNode.Y-1] <= curNode.P.BoardMap[int(curNode.X)][curNode.Y] {

			up := GridSpot{int(curNode.X), int(curNode.Y) - 1, curNode.P.BoardMap[int(curNode.X)][int(curNode.Y)-1]}
			potentialSpots = append(potentialSpots, up)
		}
		if int(curNode.X)+1 < curNode.P.Width &&
			int(curNode.X)+1 >= 0 &&
			int(curNode.Y) < curNode.P.Height &&
			curNode.Y >= 0 &&

			curNode.P.BoardMap[int(curNode.X)+1][int(curNode.Y)] != -1 &&
			curNode.P.BoolGrid[int(curNode.X)+1][int(curNode.Y)] == false { // &&
			//curNode.P.BoardMap[int(curNode.X)+1][curNode.Y] <= curNode.P.BoardMap[int(curNode.X)][curNode.Y] {

			right := GridSpot{int(curNode.X) + 1, int(curNode.Y), curNode.P.BoardMap[int(curNode.X)+1][int(curNode.Y)]}
			potentialSpots = append(potentialSpots, right)
		}
		if int(curNode.Y)+1 < curNode.P.Height &&
			curNode.Y+1 >= 0 &&
			int(curNode.X) < curNode.P.Width &&
			int(curNode.X) >= 0 &&

			curNode.P.BoardMap[int(curNode.X)][int(curNode.Y)+1] != -1 &&
			curNode.P.BoolGrid[int(curNode.X)][int(curNode.Y)+1] == false { //&&
			//curNode.P.BoardMap[int(curNode.X)][curNode.Y+1] <= curNode.P.BoardMap[int(curNode.X)][curNode.Y] {

			down := GridSpot{int(curNode.X), int(curNode.Y) + 1, curNode.P.BoardMap[int(curNode.X)][int(curNode.Y)+1]}
			potentialSpots = append(potentialSpots, down)
		}
		if int(curNode.X)-1 >= 0 &&
			int(curNode.X)-1 < curNode.P.Width &&
			curNode.Y >= 0 &&
			int(curNode.Y) < curNode.P.Height &&

			curNode.P.BoardMap[int(curNode.X)-1][int(curNode.Y)] != -1 &&
			curNode.P.BoolGrid[int(curNode.X)-1][int(curNode.Y)] == false { // &&
			//curNode.P.BoardMap[int(curNode.X)-1][curNode.Y] <= curNode.P.BoardMap[int(curNode.X)][curNode.Y] {

			left := GridSpot{int(curNode.X) - 1, int(curNode.Y), curNode.P.BoardMap[int(curNode.X)-1][int(curNode.Y)]}
			potentialSpots = append(potentialSpots, left)
		}

		sort.Sort(byRandom(potentialSpots))
		sort.Sort(byRandom(potentialSpots))
		sort.Sort(byRandom(potentialSpots))
		sort.Sort(byValue(potentialSpots))

		/*for i := 0; i < len(potentialSpots); i++ {
			if curNode.P.Grid[potentialSpots[i].Y/curNode.P.YDiv][potentialSpots[i].X/curNode.P.XDiv].ActualNumNodes <= curNode.P.SquareCapacity {
				int(curNode.X) = potentialSpots[i].X
				curNode.Y = potentialSpots[i].Y
				break
			}
		}*/

		//If there are no potential spots, do not move
		if len(potentialSpots) > 0 {
			curNode.X = float32(potentialSpots[0].X)
			curNode.Y = float32(potentialSpots[0].Y)
		}

		//Change number of nodes in square
		/*if int(curNode.X)/curNode.P.XDiv != curNode.OldX || curNode.Y/curNode.P.YDiv != curNode.OldY {
			curNode.P.Grid[curNode.Y/curNode.P.YDiv][int(curNode.X)/curNode.P.XDiv].ActualNumNodes = curNode.P.Grid[curNode.Y/curNode.P.YDiv][int(curNode.X)/curNode.P.XDiv].ActualNumNodes + 1
			curNode.P.Grid[curNode.OldY][curNode.OldX].ActualNumNodes = curNode.P.Grid[curNode.OldY][curNode.OldX].ActualNumNodes - 1
		}*/

		//curNode.P.Server.UpdateSquareNumNodes()
		if curNode.Diffx == 0 && curNode.Diffy == 0 || curNode.Sitting < 0 {
			curNode.Sitting = curNode.Sitting + 1
		} else {
			curNode.Sitting = 0
		}
	}
}

func (curNode *NodeImpl) Recalibrate() {
	curNode.Sensitivity = curNode.InitialSensitivity
	curNode.NodeTime = 0
	//fmt.Printf("Node %v recalibrated!\curNode", curNode.Id)
	curNode.Recalibrated = true
}

//Returns the arr with the element at Index curNode removed
func Remove_index(arr []Path, curNode int) []Path {
	return arr[:curNode+copy(arr[curNode:], arr[curNode+1:])]
}

//Returns the array with the range of elements from Index
//	a to b removed
func Remove_range(arr []Coord, a, b int) []Coord {
	if b > a {
		temp := b
		b = a
		a = temp
	}
	if len(arr) == 0 {
		if b+1 >= len(arr) {
			return arr[:a]
		} else {
			return append(arr[:a], arr[b+1:]...)
		}
	} else {
		new_arr := make([]Coord, 0)
		return new_arr
	}
}

//Returns the array with the specified array inserted inside at
//	Index curNode
func Insert_array(arr1 []Coord, arr2 []Coord, curNode int) []Coord {
	if len(arr1) == 0 {
		return arr2
	} else {
		return append(arr1[:curNode], append(arr2, arr1[curNode:]...)...)
	}
}

//Returns the array with the element at ind1 moved to ind2
//ind1 must always be further in the array than ind2
func Remove_and_insert(arr []Coord, ind1, ind2 int) []Coord {
	arr1 := make([]Coord, 0)
	c := arr[ind1]
	arr = arr[:ind1+copy(arr[ind1:], arr[ind1+1:])]
	arr1 = append(arr1, c)
	return append(arr[:ind2], append(arr1, arr[ind2:]...)...)
}

func (curNode *NodeImpl) LogBatteryPower(t int){
	//fmt.Println("entered function")
	//t should be p.TimeStep
	if(curNode.BatteryOverTime == nil){
		curNode.BatteryOverTime = map[int]float32{}
	}
	curNode.BatteryOverTime[t] = curNode.Battery;
	//used to test the log file writing and the python processing code
	//if(curNode.Id%4==0){
	//	curNode.DecrementPowerSensor()
	//	curNode.DecrementPower4G(100)
	//}
	//if(curNode.Id%3==0){
	//	curNode.DecrementPowerSensor()
	//}
}


//decrement battery due to transmitting/receiving over BlueTooth
func (curNode *NodeImpl) DecrementPowerBT(){
	curNode.Battery = curNode.Battery - float32(curNode.P.SamplingLossBTCM)
}

//decrement battery due to transmitting/receiving over WiFi
func (curNode *NodeImpl) DecrementPowerWifi(packet int){
	curNode.Battery = curNode.Battery - float32(curNode.P.SamplingLossWifiCM) *curNode.Battery
}

//decrement battery due to transmitting/receiving over 4G
func (curNode *NodeImpl) DecrementPower4G(){
	curNode.Battery = curNode.Battery - float32(curNode.P.SamplingLoss4GCM)
}

//decrement battery due to sampling Accelerometer
func (curNode *NodeImpl) DecrementPowerAccel(){
	curNode.Battery = curNode.Battery - float32(curNode.P.SamplingLossAccelCM) *curNode.Battery
}

//decrement battery due to transmitting/receiving GPS
func (curNode *NodeImpl) DecrementPowerGPS(){
	curNode.Battery = curNode.Battery - float32(curNode.P.SamplingLossGPSCM)
}

//decrement battery due to using GPS
func (curNode *NodeImpl) DecrementPowerSensor(){
	curNode.Battery = curNode.Battery - float32(curNode.P.SamplingLossSensorCM)
}


/* updateHistory shifts all values in the sample history slice to the right and adds the Value at the beginning
Therefore, each Time a node takes a sample in main, it also adds this sample to the beginning of the sample history.
Each sample is only stored until ln more samples have been taken (this variable is in hello.go)
*/
func (curNode *NodeImpl) UpdateHistory(newValue float32) {

	//loop through the sample history slice in reverse order, excluding 0th Index
	for i := len(curNode.SampleHistory) - 1; i > 0; i-- {
		curNode.SampleHistory[i] = curNode.SampleHistory[i-1] //set the current Index equal to the Value of the previous Index
	}

	curNode.SampleHistory[0] = newValue //set 0th Index to new measured Value

	/* Now calculate the weighted average of the sample history. Note that if a node is stationary, all values
	averaged over are weighted equally. The faster the node is moving, the less the older values are worth when
	calculating the average, because in that case we want the average to more closely reflect the newer values
	*/
	var sum float32
	var numSamples int //variable for number of samples to average over

	var decreaseRatio = curNode.SpeedWeight / 100.0

	if curNode.TotalSamples > len(curNode.SampleHistory) { //if the node has taken more than x total samples
		numSamples = len(curNode.SampleHistory) //we only average over the x most recent ones
	} else { //if it doesn't have x samples taken yet
		numSamples = curNode.TotalSamples //we only average over the number of samples it's taken
	}

	for i := 0; i < numSamples; i++ {
		if curNode.SampleHistory[i] != 0 {
			//weight the values of the sampleHistory when added to the sum variable based on the speed, so older values are weighted less
			sum += curNode.SampleHistory[i] - ((decreaseRatio) * float32(i))
		} else {
			sum += 0
		}
	}
	curNode.Avg = sum / float32(numSamples)
}

func (curNode *NodeImpl) getDriftSlope() (float32, float32){
	var r float32
	var slope float32

	var sum float32
	var yAvg float32 = 0.0
	squareSumX := 0.0
	squareSumY := 0.0

	var xSum float32
	var ySum float32
	var xySum float32
	var xSqrSum float32
	//size := float32(len(curNode.SampleHistory))

	for i:= range curNode.SampleHistory {
		ySum += float32(i)
	}
	yAvg = ySum / float32(len(curNode.SampleHistory))
	for i := range curNode.SampleHistory {
		sum += (curNode.SampleHistory[i] - curNode.Avg) * (float32(i) - yAvg)
		squareSumX += math.Pow( float64(curNode.SampleHistory[i] - curNode.Avg), 2)
		squareSumY += math.Pow( float64(i - 1), 2)

		xSum += curNode.SampleHistory[i]
		xySum += curNode.SampleHistory[i] * float32(i)
		xSqrSum += float32(math.Pow(float64(curNode.SampleHistory[i]), 2))
	}
	r = sum / float32(math.Sqrt(squareSumX * squareSumY))
	//slope = ( (size * xySum) - (xSum * ySum) ) / ( (size * float32(xSqrSum)) - float32(math.Pow(float64(xSum), 2)) )
	//slope = sum / float32(squareSumX)
	if r > 1 || r < -1 {
		fmt.Printf("Bad r Value! Got %v\n", r)
	}
	return r, slope
}

/* this function increments a node's total number of samples by 1
it's called whenever the node takes a new sample */
func (curNode *NodeImpl) IncrementTotalSamples() {
	curNode.TotalSamples++
}

//getter function for average
func (curNode *NodeImpl) GetAvg() float32 {
	return curNode.Avg
}

//increases numResets field
func (curNode *NodeImpl) IncrementNumResets() {
	curNode.NumResets++
}

//setter function for concentration field
func (curNode *NodeImpl) SetConcentration(conc float64) {
	curNode.Concentration = conc
}

//getter function for ID field
func (curNode *NodeImpl) GetID() int {
	return curNode.Id
}

//getter function for x and y locations
func (curNode *NodeImpl) GetLoc() (float32, float32) {
	return curNode.X, curNode.Y
}

//setter function for S0
func (curNode *NodeImpl) SetS0(s0 float64) {
	curNode.S0 = s0
}

//setter function for S1
func (curNode *NodeImpl) SetS1(s1 float64) {
	curNode.S1 = s1
}

//setter function for S2
func (curNode *NodeImpl) SetS2(s2 float64) {
	curNode.S2 = s2
}

//setter function for E0
func (curNode *NodeImpl) SetE0(e0 float64) {
	curNode.E0 = e0
}

//setter function for E1
func (curNode *NodeImpl) SetE1(e1 float64) {
	curNode.E1 = e1
}

//setter function for E2
func (curNode *NodeImpl) SetE2(e2 float64) {
	curNode.E2 = e2
}

//setter function for ET1
func (curNode *NodeImpl) SetET1(et1 float64) {
	curNode.ET1 = et1
}

//setter function for ET2
func (curNode *NodeImpl) SetET2(et2 float64) {
	curNode.ET2 = et2
}

//getter function for all parameters
func (curNode *NodeImpl) GetParams() (float64, float64, float64, float64, float64, float64, float64, float64) {
	return curNode.S0, curNode.S1, curNode.S2, curNode.E0, curNode.E1, curNode.E2, curNode.ET1, curNode.ET2
}

//getter function for just S0 - S2 parameters
func (curNode *NodeImpl) GetCoefficients() (float64, float64, float64) {
	return curNode.S0, curNode.S1, curNode.S2
}

//getter function for x
func (curNode *NodeImpl) GetX() float32 {
	return curNode.X
}

//getter function for y
func (curNode *NodeImpl) GetY() float32 {
	return curNode.Y
}

//This is the actual upload to the server
func (curNode *NodeImpl) Server() {
	//getData(&s,curNode.XPos[0:curNode.BufferI],curNode.YPos[0:curNode.BufferI],curNode.Value[0:curNode.BufferI],curNode.Time[0:curNode.BufferI], curNode.Id,curNode.BufferI)
	curNode.BufferI = 0
}

//Returns node distance to the bomb
func (curNode *NodeImpl) GeoDist(b Bomb) float32 {
	//this needs to be changed
	return float32(math.Pow(float64(math.Abs(float64(curNode.X)-float64(b.X))), 2) + math.Pow(float64(math.Abs(float64(curNode.Y)-float64(b.Y))), 2))
}

//Returns array of accelerometer speeds recorded for a specific node
func (curNode *NodeImpl) GetSpeed() []float32 {
	return curNode.AccelerometerSpeed
}

//Returns a different version of the distance to the bomb
func (curNode *NodeImpl) GetValue() int {
	return int(math.Sqrt(math.Pow(float64(int(curNode.X)-curNode.P.B.X), 2) + math.Pow(float64(curNode.Y-float32(curNode.P.B.Y)), 2)))
}


func (curNode *NodeImpl) Distance(b Bomb) float32 {
	return float32(math.Sqrt(math.Pow(float64(math.Abs(float64(curNode.X)-float64(b.X))),2) + math.Pow(float64(math.Abs(float64(curNode.Y)-float64(b.Y))),2)))
}

//Returns a float representing the detection of the bomb
//	by the specific node depending on distance
func RawConcentration(dist float32) float32 {
	//dist := curNode.Distance(b)
	//dist := float32(math.Pow(float64(math.Abs(float64(curNode.X)-float64(b.X))), 2) + math.Pow(float64(math.Abs(float64(curNode.Y)-float64(b.Y))), 2))

	if dist < .1 {
		return 1000
	} else {
		//reading := float32(1000.0/ math.Pow(float64(dist)/.2, 3))
		reading := float32(1000.0/ (float64(dist)/.1))
		return reading
	}
}

//Takes cares of taking a node's readings and printing detections and stuff
func (curNode *NodeImpl) GetReadings() {


	if curNode.Valid { //Check if node should actually take readings or if it hasn't shown up yet
		newX, newY := curNode.GetLoc()

		RawConc := RawConcentration(curNode.Distance(*curNode.P.B)/2) //this is the node's reported Value without error

		//need to get the correct Time reading Value from system
		//need to verify where we read from

		//Calculate error, sensitivity, and noise, as per the matlab code
		S0, S1, S2, E0, E1, E2, ET1, ET2 := curNode.GetParams()
		sError := (S0 + E0) + (S1+E1)*math.Exp(-float64(curNode.NodeTime)/(curNode.P.Tau1+ET1)) + (S2+E2)*math.Exp(-float64(curNode.NodeTime)/(curNode.P.Tau2+ET2))
		curNode.Sensitivity = S0 + (S1)*math.Exp(-float64(curNode.NodeTime)/curNode.P.Tau1) + (S2)*math.Exp(-float64(curNode.NodeTime)/curNode.P.Tau2)
		sNoise := rand.NormFloat64()*0.5*curNode.P.ErrorModifierCM + float64(RawConc)*sError

		errorDist := sNoise / curNode.Sensitivity //this is the node's actual reading with error
		clean := float64(RawConc) / curNode.Sensitivity


		ADCRead := float64(curNode.ADCReading(float32(errorDist)))
		ADCClean := float64(curNode.ADCReading(float32(clean)))

		d := curNode.Distance(*curNode.P.B)/2
		if d < 10 {
			fmt.Fprintln(curNode.P.MoveReadingsFile, "Time:", curNode.P.CurrentTime/1000, "ID:", curNode.Id, "X:", newX, "Y:",  newY, "Dist:", d, "ADCClean:", ADCClean, "ADCError:", ADCRead, "CleanSense:", clean, "Error:", errorDist, "Raw:", RawConc)
		}


		//increment node Time
		curNode.NodeTime++

		curNode.IncrementTotalSamples()
		curNode.UpdateHistory(float32(errorDist))

		//If the reading is more than 2 standard deviations away from the grid average, then recalibrate
		//gridAverage := curNode.P.Grid[curNode.Row(curNode.P.YDiv)][curNode.Col(curNode.P.XDiv)].Avg
		//standDev := grid[curNode.Row(yDiv)][curNode.Col(xDiv)].StdDev

		//New condition added: also recalibrate when the node's sensitivity is <= 1/2 of its original sensitvity
		//New condition added: Check to make sure the sensor was pinged this iteration
		if ((curNode.Sensitivity <= (curNode.InitialSensitivity / 2)) && curNode.P.Iterations_used != 0) {
			curNode.Recalibrate()
			curNode.Recalibrated = true
			curNode.IncrementNumResets()
		}

		//printing statements to log files, only if the sensor was pinged this iteration
		//if curNode.HasCheckedSensor && nodesPrint{
		if curNode.P.NodesPrint {
			if curNode.Recalibrated {
				fmt.Fprintln(curNode.P.NodeFile, "ID:", curNode.GetID(), "Average:", curNode.GetAvg(), "Reading:", RawConc, "Error Reading:", errorDist, "Recalibrated")
			} else {
				fmt.Fprintln(curNode.P.NodeFile, "ID:", curNode.GetID(), "Average:", curNode.GetAvg(), "Reading:", RawConc, "Error Reading:", errorDist)
			}
			//fmt.Fprintln(nodeFile, "battery:", int(curNode.Battery),)
			curNode.Recalibrated = false
		}


		if (ADCRead > curNode.P.DetectionThreshold && ADCClean < curNode.P.DetectionThreshold) || (ADCRead > curNode.P.DetectionThreshold && d > float32(curNode.P.DetectionDistance/2.0)){
			fmt.Fprintln(curNode.P.DetectionFile, fmt.Sprintf("FP T: %v ID: %v (%v, %v) D: %v C: %v E: %v SE: %.3f S: %.3f R: %.3f", curNode.P.CurrentTime, curNode.Id, curNode.X, curNode.Y, d, ADCClean, ADCRead, sError, curNode.Sensitivity, RawConc))
		} else if ADCRead < curNode.P.DetectionThreshold && ADCClean > curNode.P.DetectionThreshold && d < float32(curNode.P.DetectionDistance/2.0) {
			fmt.Fprintln(curNode.P.DetectionFile, fmt.Sprintf("FN T: %v ID: %v (%v, %v) D: %v C: %v E: %v SE: %.3f S: %.3f R: %.3f", curNode.P.CurrentTime, curNode.Id, curNode.X, curNode.Y, d, ADCClean, ADCRead, sError, curNode.Sensitivity, RawConc))
		} else if ADCRead > curNode.P.DetectionThreshold && ADCClean > curNode.P.DetectionThreshold {
			fmt.Fprintln(curNode.P.DetectionFile, fmt.Sprintf("TP T: %v ID: %v (%v, %v) D: %v C: %v E: %v SE: %.3f S: %.3f R: %.3f", curNode.P.CurrentTime, curNode.Id, curNode.X, curNode.Y, d, ADCClean, ADCRead, sError, curNode.Sensitivity, RawConc))
		}

		if(curNode.P.ClusteringOn){
			newReading := Reading{ADCRead, newX, newY, curNode.P.CurrentTime, curNode.GetID(), curNode.GetCHID()}
			curNode.ReadingsBuffer = append(curNode.ReadingsBuffer, newReading)
		} else{
			curNode.P.Server.Send(curNode, Reading{ADCRead, newX, newY, curNode.P.CurrentTime, curNode.GetID(), curNode.GetCHID()})
		}

	}
	curNode.P.Events.Push(&Event{curNode, SENSE, curNode.P.CurrentTime + 500, 0})


}


func (curNode * NodeImpl) GetCHID()(int){
	if(curNode.NodeClusterParams != nil && curNode.NodeClusterParams.CurrentCluster != nil && curNode.NodeClusterParams.CurrentCluster.ClusterHead != nil){
		return curNode.NodeClusterParams.CurrentCluster.ClusterHead.Id
	} else{
		return -1
	}
}

func interpolateReading(x , y float32, time, timeStep int, p *Params) float32{
	oldX := int(x)
	oldY := int(y)
	nextX := int(math.Ceil(float64(x)))
	nextY := int(math.Ceil(float64(y)))

	//fmt.Printf("%v %v %v %v\n", oldX, oldY, nextX, nextY)
	//calculate reading at last 'even' position
	oldReadingA := p.SensorReadings[oldX][oldY][timeStep]
	futureReadingA := p.SensorReadings[oldX][oldY][timeStep]
	if(timeStep < p.MaxTimeStep) {
		futureReadingA = p.SensorReadings[oldX][oldY][timeStep+1]
	}

	//calculate reading at next 'even' position
	oldReadingB := p.SensorReadings[nextX][nextY][timeStep]
	futureReadingB := p.SensorReadings[nextX][nextY][timeStep]
	if(timeStep < p.MaxTimeStep) {
		futureReadingB = p.SensorReadings[nextX][nextY][timeStep+1]
	}

	//fmt.Printf("%v %v %v %v\n", oldReadingA, oldReadingB, futureReadingA, futureReadingB)
	totalDistance := float32(math.Sqrt(math.Pow(math.Abs(float64(nextX - oldX)), 2) + math.Pow(math.Abs(float64(nextY - oldY)), 2)))
	partialDist := float32((math.Sqrt(math.Pow(math.Abs(float64(float32(nextX) - x)), 2) + math.Pow(math.Abs(float64(float32(nextY) - y)), 2))))


	//determine distance we have covered between the two positions
	portionDist := partialDist / totalDistance
	if (totalDistance == 0) {
		portionDist = 1.0
	}

	//fmt.Printf("%v %v %v\n", totalDistance, partialDist, portionDist)


	oldReading := (float32(oldReadingA - oldReadingB) * portionDist + float32(oldReadingA))  //t = 0 original and next position average
	futureReading := (float32(futureReadingA - futureReadingB) * portionDist + float32(futureReadingA)) //t = 1 original and next position average
	//fmt.Printf("%v %v\n", oldReading, futureReading)

	floatTime := float32(time)/1000
	oldTime := p.SensorTimes[timeStep]
	nextTime := p.SensorTimes[timeStep]
	if(timeStep < p.MaxTimeStep) {
		nextTime = p.SensorTimes[timeStep + 1]
	}


	portionTime := (floatTime - float32(oldTime))/float32(nextTime - oldTime)
	if(oldTime == nextTime) {
		portionTime = 1.0
	}
	//fmt.Printf("%v %v %v %v\n", floatTime, oldTime, nextTime, portionTime)

	return (futureReading - oldReading)*portionTime + oldReading
}

//Takes cares of taking a node's readings and printing detections and stuff
func (curNode *NodeImpl) GetReadingsCSV() {

	if curNode.Valid { //check if node has shown up yet

		newX, newY := curNode.GetLoc()


		RawConcentration := interpolateReading(newX, newY, curNode.P.CurrentTime, curNode.P.TimeStep, curNode.P)

		//need to get the correct Time reading Value from system
		//need to verify where we read from

		//Calculate error, sensitivity, and noise, as per the matlab code
		S0, S1, S2, E0, E1, E2, ET1, ET2 := curNode.GetParams()
		sError := (S0 + E0) + (S1+E1)*math.Exp(-float64(curNode.NodeTime)/(curNode.P.Tau1+ET1)) + (S2+E2)*math.Exp(-float64(curNode.NodeTime)/(curNode.P.Tau2+ET2))
		curNode.Sensitivity = S0 + (S1)*math.Exp(-float64(curNode.NodeTime)/curNode.P.Tau1) + (S2)*math.Exp(-float64(curNode.NodeTime)/curNode.P.Tau2)
		sNoise := rand.NormFloat64()*0.5*curNode.P.ErrorModifierCM + float64(RawConcentration)*sError

		errorDist := sNoise / curNode.Sensitivity //this is the node's actual reading with error
		clean := float64(RawConcentration) / curNode.Sensitivity


		ADCRead := float64(curNode.ADCReading(float32(errorDist)))
		ADCClean := float64(curNode.ADCReading(float32(clean)))

		d := curNode.Distance(*curNode.P.B)/2
		if d < 10 {
			fmt.Fprintln(curNode.P.MoveReadingsFile, "Time:", curNode.P.CurrentTime/1000, "ID:", curNode.Id, "X:", newX, "Y:",  newY, "Dist:", d, "ADCClean:", ADCClean, "ADCError:", ADCRead, "CleanSense:", clean, "Error:", errorDist, "Raw:", RawConcentration)
		}

		//increment node Time
		curNode.NodeTime++

		//if curNode.HasCheckedSensor {
			curNode.IncrementTotalSamples()
			curNode.UpdateHistory(float32(errorDist))
		//}
		
		//If the reading is more than 2 standard deviations away from the grid average, then recalibrate
		//gridAverage := curNode.P.Grid[curNode.Row(curNode.P.YDiv)][curNode.Col(curNode.P.XDiv)].Avg
		//standDev := grid[curNode.Row(yDiv)][curNode.Col(xDiv)].StdDev

		//New condition added: also recalibrate when the node's sensitivity is <= 1/10 of its original sensitvity
		//New condition added: Check to make sure the sensor was pinged this iteration
		if ((curNode.Sensitivity <= (curNode.InitialSensitivity / 2))  && curNode.P.Iterations_used != 0) {
			curNode.Recalibrate()
			curNode.Recalibrated = true
			curNode.IncrementNumResets()
		}

		//printing statements to log files, only if the sensor was pinged this iteration
		//if curNode.HasCheckedSensor && nodesPrint{
		if curNode.P.NodesPrint {
			if curNode.Recalibrated {
				fmt.Fprintln(curNode.P.NodeFile, "ID:", curNode.GetID(), "Average:", curNode.GetAvg(), "Reading:", RawConcentration, "Error Reading:", errorDist, "Recalibrated")
			} else {
				fmt.Fprintln(curNode.P.NodeFile, "ID:", curNode.GetID(), "Average:", curNode.GetAvg(), "Reading:", RawConcentration, "Error Reading:", errorDist)
			}
			//fmt.Fprintln(nodeFile, "battery:", int(curNode.Battery),)
			curNode.Recalibrated = false
		}

		if ADCRead > curNode.P.DetectionThreshold && ADCClean < curNode.P.DetectionThreshold{
			fmt.Fprintln(curNode.P.DetectionFile, fmt.Sprintf("FP T: %v ID: %v (%v, %v) D: %v C: %v E: %v SE: %.3f S: %.3f R: %.3f", curNode.P.CurrentTime, curNode.Id, curNode.X, curNode.Y, d, ADCClean, ADCRead, sError, curNode.Sensitivity, RawConcentration))
		} else if ADCRead < curNode.P.DetectionThreshold && ADCClean > curNode.P.DetectionThreshold {
			fmt.Fprintln(curNode.P.DetectionFile, fmt.Sprintf("FN T: %v ID: %v (%v, %v) D: %v C: %v E: %v SE: %.3f S: %.3f R: %.3f", curNode.P.CurrentTime, curNode.Id, curNode.X, curNode.Y, d, ADCClean, ADCRead, sError, curNode.Sensitivity, RawConcentration))
		} else if ADCRead > curNode.P.DetectionThreshold && ADCClean > curNode.P.DetectionThreshold {
			fmt.Fprintln(curNode.P.DetectionFile, fmt.Sprintf("TP T: %v ID: %v (%v, %v) D: %v C: %v E: %v SE: %.3f S: %.3f R: %.3f", curNode.P.CurrentTime, curNode.Id, curNode.X, curNode.Y, d, ADCClean, ADCRead, sError, curNode.Sensitivity, RawConcentration))
		}





		//Receives the node's distance and calculates its running average
		//for that square
		//Only do this if the sensor was pinged this iteration

		if curNode.Valid {
			if(curNode.P.ClusteringOn){
				if(curNode.IsClusterHead){
					if(curNode.P.CurrentTime - curNode.TimeLastSentReadings > curNode.P.CHSensingTime*1000 || len(curNode.ReadingsBuffer)>curNode.P.MaxCHReadingBufferSize){
						//Print to file
						oldest,newest := curNode.GetOldestNewestReadings()
						fmt.Fprintf(curNode.P.ClusterReadings,"%d\t %d\t %d\t Head to Server\t %d\t %d\t %d\n",curNode.P.CurrentTime/1000,curNode.TimeLastSentReadings/1000,curNode.Id,oldest/1000,newest/1000,len(curNode.ReadingsBuffer))

						//Send to Server and clear readings
						curNode.SendtoServer()
						curNode.TimeLastSentReadings = curNode.P.CurrentTime
					}
				} else if(curNode.IsClusterMember){
					if(curNode.P.CurrentTime - curNode.TimeLastSentReadings > curNode.P.CMSensingTime*1000 || len(curNode.ReadingsBuffer)>curNode.P.MaxCMReadingBufferSize){
						//Print to file
						oldest,newest := curNode.GetOldestNewestReadings()
						fmt.Fprintf(curNode.P.ClusterReadings,"%d\t %d\t %d\t Node to Head\t %d\t %d\t %d\n",curNode.P.CurrentTime/1000,curNode.TimeLastSentReadings/1000,curNode.Id,oldest/1000,newest/1000,len(curNode.ReadingsBuffer))

						//Send to Server and clear readings
						curNode.SendtoClusterHead()
						curNode.TimeLastSentReadings = curNode.P.CurrentTime
					}
				}

			newReading := Reading{ADCRead, newX, newY, curNode.P.CurrentTime, curNode.GetID(), curNode.GetCHID()}
			curNode.ReadingsBuffer = append(curNode.ReadingsBuffer, newReading)

			} else{
				curNode.P.Server.Send(curNode, Reading{ADCRead, newX, newY, curNode.P.CurrentTime, curNode.GetID(), curNode.GetCHID()})
			}
		}

	}
	curNode.P.Events.Push(&Event{curNode, SENSE, curNode.P.CurrentTime + 500, 0})
}

//Get min and max readings
func (curNode * NodeImpl) GetOldestNewestReadings()(oldest int, newest int){

	oldest = curNode.ReadingsBuffer[0].Time
	newest = curNode.ReadingsBuffer[0].Time

	for i:=0; i<len(curNode.ReadingsBuffer); i++{
		if(curNode.ReadingsBuffer[i].Time < oldest){
			oldest = curNode.ReadingsBuffer[i].Time
		}
		if (curNode.ReadingsBuffer[i].Time > newest){
			newest = curNode.ReadingsBuffer[i].Time
		}
	}

	return oldest, newest
}

//Cluster Head sends all readings to the Server
func (curNode *NodeImpl) SendtoServer(){

	//Iterates through Readings Buffer sending all readings
	for i:=0; i<len(curNode.ReadingsBuffer); i++{
		curNode.P.Server.Send(curNode, curNode.ReadingsBuffer[i])
	}

	//clears Clusterhead reading buffer
	curNode.ReadingsBuffer = []Reading{}

}

//Node sends to Cluster Head
func (curNode *NodeImpl) SendtoClusterHead(){

	//Append to clusterhead readings buffer
	clusterHead := curNode.NodeClusterParams.CurrentCluster.ClusterHead
	clusterHead.ReadingsBuffer = append(clusterHead.ReadingsBuffer, curNode.ReadingsBuffer...)

	//clear node's readings
	curNode.ReadingsBuffer = []Reading{}
}

func interpolate (start int, end int, portion float32) float32{
	return (float32(end-start) * portion + float32(start))
}

//HandleMovementCSV does the same as HandleMovement
func (curNode *NodeImpl) MoveCSV(p *Params) {
	//time := p.Iterations_used
	floatTemp := float32(p.CurrentTime)
	intTime := int(floatTemp/1000)
	portion := (floatTemp / 1000) - float32(intTime)

	id := curNode.GetID()

	if curNode.Valid {
		oldX, oldY := curNode.GetLoc()
		p.BoolGrid[int(oldX)][int(oldY)] = false //set the old spot false since the node will now move away

		curNode.X = interpolate(p.NodeMovements[id][intTime].X, p.NodeMovements[id][intTime+1].X, portion)
		curNode.Y = interpolate(p.NodeMovements[id][intTime].Y, p.NodeMovements[id][intTime+1].Y, portion)

		//set the new location in the boolean field to true
		newX, newY := curNode.GetLoc()
		p.BoolGrid[int(newX)][int(newY)] = true
		//Handle movement in the tree (remove and reinsert if neccessary)
		curNode.NodeBounds.X = float64(newX)
		curNode.NodeBounds.Y = float64(newY)
		curNode.P.NodeTree.NodeMovement(curNode.NodeBounds)
	}


	if !curNode.Valid {
		curNode.Valid = curNode.TurnValid(p.NodeMovements[id][intTime].X, p.NodeMovements[id][intTime].Y, p)
		curNode.X = float32(p.NodeMovements[id][intTime].X)
		curNode.Y = float32(p.NodeMovements[id][intTime].Y)
	}

	//fmt.Printf("%v %v %v %v %v %v %v\n", newX, newY, int(newX), int(newY), p.Width, p.Height, curNode.Valid)
	/*if curNode.InBounds(p) {
		curNode.Valid = true
	} else {
		curNode.Valid = false
	}*/

	//Add the node into its new Square's p.TotalNodes
	//If the node hasn't left the square, that Square's p.TotalNodes will
	//remain the same after these calculations
}

//HandleMovement adjusts BoolGrid when nodes move around the map
func (curNode *NodeImpl) MoveNormal(p *Params) {

	oldX, oldY := curNode.GetLoc()
	p.BoolGrid[int(oldX)][int(oldY)] = false //set the old spot false since the node will now move away

	//move the node to its new location
	curNode.Move(p)

	//set the new location in the boolean field to true
	newX, newY := curNode.GetLoc()
	p.BoolGrid[int(newX)][int(newY)] = true

	//Handle movement in the tree (remove and reinsert if neccessary)
	curNode.P.NodeTree.NodeMovement(curNode.NodeBounds)

	//Add the node into its new Square's p.TotalNodes
	//If the node hasn't left the square, that Square's p.TotalNodes will
	//remain the same after these calculations

}

func rangeInt(min, max int) int { //returns a random number between max and min
	return rand.Intn(max-min) + min
}