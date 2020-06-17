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

	allReadings 	   [1000]float64
	calibrateTimes 	   []int
	calibrateReading   []float64

	BatteryOverTime	   map[int]float32

	TotalPacketsSent    int
	TotalBytesSent		int
	IsClusterMember		bool
	IsClusterHead		bool
	Recalibrated 		bool

	NodeClusterParams *ClusterMemberParams
	ReadingsBuffer		[]Reading
	TimeLastSentReadings	int
	TimeLastAccel		int
	LastMoveTime		int
	Velocity			float64
	SampleRate			float64

	LastReading			 float64
	ReadingPercentChange float64
	LastAccel			 float64
	MovePercentChange	 float64
	MovementModifier	 float64
	SensorModifier		 float64

	SampleRateSensor	 float64
	SampleRateMovement   float64
	SampleRateBattery    float64

	ScheduledEvent		 *Event
	BatteryPercent		 float64

	LastNSampleRates	 []float64
	CurTree				 *Quadtree
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
func (node *NodeImpl) Row(div int) int {
	return int(node.Y) / div
}

//Returns the y Index of the square in which the specific
//	node currently resides
func (node *NodeImpl) Col(div int) int {
	return int(node.X) / div
}

func (node *NodeImpl) InBounds(p *Params) bool {
	return int(node.X) < p.Width && int(node.X) >= 0 && int(node.Y) < p.Height && node.Y >= 0
}

func (node *NodeImpl) TurnValid(x, y int, p *Params) bool {
	return x < p.Width && x >= 0 && y < p.Height && y >= 0
}



func (node *NodeImpl) ADCReading(raw float32) int {

	level := (raw - node.P.ADCOffset)/ node.P.ADCWidth

	if level > node.P.MaxADC {
		level = node.P.MaxADC
	} else if level < 0 {
		level = 0
	}

	return int(level)
}

func (node NodeImpl) String() string {
	//return fmt.Sprintf("x: %v y: %v Id: %v battery: %v sensor checked: %v sensor checks: %v GPS checked: %v GPS checks: %v server checked: %v server checks: %v buffer: %v ", int(curNode.X), curNode.Y, curNode.Id, curNode.Battery, curNode.HasCheckedSensor, curNode.TotalChecksSensor, curNode.HasCheckedGPS, curNode.TotalChecksGPS, curNode.HasCheckedServer, curNode.TotalChecksServer,curNode.BufferI)
	//return fmt.Sprintf("x: %v y: %v valid: %v", int(curNode.X), int(curNode.Y), curNode.Valid)
	//return fmt.Sprintf("battery: %v sensor checked: %v GPS checked: %v ", int(curNode.Battery), curNode.HasCheckedSensor, curNode.HasCheckedGPS)
	return fmt.Sprintf("battery: %v sensor checked: %v GPS checked: %v ", int(node.Battery), true, true)

}

func (c Coord) String() string {
	return fmt.Sprintf("{%v %v %v}", c.X, c.Y, c.Time)
}

func (c Coord) Equals(c2 Coord) bool {
	return c.X == c2.X && c.Y == c2.Y
}

func (node *NodeImpl) Move(p *Params) {
	if node.Sitting <= p.SittingStopThresholdCM {
		node.OldX = int(node.X) / p.XDiv
		node.OldY = int(node.Y) / p.YDiv

		var potentialSpots []GridSpot

		//only add the ones that are valid to move to into the list
		if int(node.Y)-1 >= 0 &&
			int(node.X) >= 0 &&
			int(node.X) < p.Width &&
			int(node.Y)-1 < p.Height &&

			p.BoardMap[int(node.X)][int(node.Y)-1] != -1 &&
			p.BoolGrid[int(node.X)][int(node.Y)-1] == false { // &&
			//curp.BoardMap[int(curNode.X)][curNode.Y-1] <= curp.BoardMap[int(curNode.X)][curNode.Y] {

			up := GridSpot{int(node.X), int(node.Y) - 1, p.BoardMap[int(node.X)][int(node.Y)-1]}
			potentialSpots = append(potentialSpots, up)
		}
		if int(node.X)+1 < p.Width &&
			int(node.X)+1 >= 0 &&
			int(node.Y) < p.Height &&
			node.Y >= 0 &&

			p.BoardMap[int(node.X)+1][int(node.Y)] != -1 &&
			p.BoolGrid[int(node.X)+1][int(node.Y)] == false { // &&
			//curp.BoardMap[int(curNode.X)+1][curNode.Y] <= curp.BoardMap[int(curNode.X)][curNode.Y] {

			right := GridSpot{int(node.X) + 1, int(node.Y), p.BoardMap[int(node.X)+1][int(node.Y)]}
			potentialSpots = append(potentialSpots, right)
		}
		if int(node.Y)+1 < p.Height &&
			node.Y+1 >= 0 &&
			int(node.X) < p.Width &&
			int(node.X) >= 0 &&

			p.BoardMap[int(node.X)][int(node.Y)+1] != -1 &&
			p.BoolGrid[int(node.X)][int(node.Y)+1] == false { //&&
			//curp.BoardMap[int(curNode.X)][curNode.Y+1] <= curp.BoardMap[int(curNode.X)][curNode.Y] {

			down := GridSpot{int(node.X), int(node.Y) + 1, p.BoardMap[int(node.X)][int(node.Y)+1]}
			potentialSpots = append(potentialSpots, down)
		}
		if int(node.X)-1 >= 0 &&
			int(node.X)-1 < p.Width &&
			node.Y >= 0 &&
			int(node.Y) < p.Height &&

			p.BoardMap[int(node.X)-1][int(node.Y)] != -1 &&
			p.BoolGrid[int(node.X)-1][int(node.Y)] == false { // &&
			//curp.BoardMap[int(curNode.X)-1][curNode.Y] <= curp.BoardMap[int(curNode.X)][curNode.Y] {

			left := GridSpot{int(node.X) - 1, int(node.Y), p.BoardMap[int(node.X)-1][int(node.Y)]}
			potentialSpots = append(potentialSpots, left)
		}

		sort.Sort(byRandom(potentialSpots))
		sort.Sort(byRandom(potentialSpots))
		sort.Sort(byRandom(potentialSpots))
		sort.Sort(byValue(potentialSpots))

		/*for i := 0; i < len(potentialSpots); i++ {
			if curp.Grid[potentialSpots[i].Y/curp.YDiv][potentialSpots[i].X/curp.XDiv].ActualNumNodes <= curp.SquareCapacity {
				int(curNode.X) = potentialSpots[i].X
				curNode.Y = potentialSpots[i].Y
				break
			}
		}*/

		//If there are no potential spots, do not move
		if len(potentialSpots) > 0 {
			node.X = float32(potentialSpots[0].X)
			node.Y = float32(potentialSpots[0].Y)
		}

		//Change number of nodes in square
		/*if int(curNode.X)/curp.XDiv != curNode.OldX || curNode.Y/curp.YDiv != curNode.OldY {
			curp.Grid[curNode.Y/curp.YDiv][int(curNode.X)/curp.XDiv].ActualNumNodes = curp.Grid[curNode.Y/curp.YDiv][int(curNode.X)/curp.XDiv].ActualNumNodes + 1
			curp.Grid[curNode.OldY][curNode.OldX].ActualNumNodes = curp.Grid[curNode.OldY][curNode.OldX].ActualNumNodes - 1
		}*/

		//curp.Server.UpdateSquareNumNodes()
		if node.Diffx == 0 && node.Diffy == 0 || node.Sitting < 0 {
			node.Sitting = node.Sitting + 1
		} else {
			node.Sitting = 0
		}
	}
}

func (node *NodeImpl) Recalibrate() {
	node.P.Server.NodeDataList[node.Id].SelfRecalTimes = append(node.P.Server.NodeDataList[node.Id].SelfRecalTimes, node.P.CurrentTime / 1000)
	node.Sensitivity = node.InitialSensitivity
	node.NodeTime = (node.P.CurrentTime/1000)
	//fmt.Fprintf(curNode.P.DriftExploreFile, "ID: %v T: %v In: %v CUR: %v NT: %v RECAL\n", curNode.Id, curNode.P.CurrentTime, curNode.InitialSensitivity, curNode.Sensitivity, curNode.NodeTime)
	//fmt.Printf("Node %v recalibrated!\curNode", curNode.Id)
	node.Recalibrated = true
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

func (node *NodeImpl) LogBatteryPower(t int){
	//fmt.Println("entered function")
	//t should be p.TimeStep
	if(node.BatteryOverTime == nil){
		node.BatteryOverTime = map[int]float32{}
	}
	node.BatteryOverTime[t] = node.Battery;
	//used to test the log file writing and the python processing code
	//if(curNode.Id%4==0){
	//	curNode.DecrementPowerSensor()
	//	curNode.DecrementPower4G(100)
	//}
	//if(curNode.Id%3==0){
	//	curNode.DecrementPowerSensor()
	//}
}

func (node *NodeImpl) SendToServer(rd Reading, tp bool){
	//int packet = num bytes in packet
	//node.TotalBytesSent += packet
	//node.TotalPacketsSent += 1

	//code to send to server goes here
	node.DecrementPower4G()
	node.P.Server.Send(node, rd, tp)
}

func (node *NodeImpl) SendToClusterHead(rd Reading, tp bool){
	//int packet = num bytes in packet
	//node.TotalBytesSent += packet
	//node.TotalPacketsSent += 1

	//code to send to cluster head goes here
	head := node.NodeClusterParams.CurrentCluster.ClusterHead

	node.DecrementPowerBT()
	head.DecrementPowerBT()
	head.DecrementPower4G()
	node.P.Server.Send(node, rd, tp)
}


//decrement battery due to transmitting/receiving over BlueTooth
func (node *NodeImpl) DecrementPowerBT(){
	node.Battery = node.Battery - node.BatteryLossBT*node.Battery
}

//decrement battery due to transmitting/receiving over WiFi
func (node *NodeImpl) DecrementPowerWifi(packet int){
	node.Battery = node.Battery - node.BatteryLossWifi
}

//decrement battery due to transmitting/receiving over 4G
func (node *NodeImpl) DecrementPower4G(){
	node.Battery = node.Battery - node.BatteryLoss4G*node.Battery
}

//decrement battery due to sampling Accelerometer
func (node *NodeImpl) DecrementPowerAccel(){
	node.Battery = node.Battery - node.BatteryLossAccelerometer*node.Battery
}

//decrement battery due to transmitting/receiving GPS
func (node *NodeImpl) DecrementPowerGPS(){
	node.Battery = node.Battery - node.BatteryLossGPS*node.Battery
}

//decrement battery due to using GPS
func (node *NodeImpl) DecrementPowerSensor(){
	node.Battery = node.Battery - node.BatteryLossSensor*node.Battery
}


/* updateHistory shifts all values in the sample history slice to the right and adds the Value at the beginning
Therefore, each Time a node takes a sample in main, it also adds this sample to the beginning of the sample history.
Each sample is only stored until ln more samples have been taken (this variable is in hello.go)
*/
func (node *NodeImpl) UpdateHistory(newValue float32) {

	//loop through the sample history slice in reverse order, excluding 0th Index
	for i := len(node.SampleHistory) - 1; i > 0; i-- {
		node.SampleHistory[i] = node.SampleHistory[i-1] //set the current Index equal to the Value of the previous Index
	}

	node.SampleHistory[0] = newValue //set 0th Index to new measured Value

	/* Now calculate the weighted average of the sample history. Note that if a node is stationary, all values
	averaged over are weighted equally. The faster the node is moving, the less the older values are worth when
	calculating the average, because in that case we want the average to more closely reflect the newer values
	*/
	var sum float32
	var numSamples int //variable for number of samples to average over

	var decreaseRatio = node.SpeedWeight / 100.0

	if node.TotalSamples > len(node.SampleHistory) { //if the node has taken more than x total samples
		numSamples = len(node.SampleHistory) //we only average over the x most recent ones
	} else { //if it doesn't have x samples taken yet
		numSamples = node.TotalSamples //we only average over the number of samples it's taken
	}

	for i := 0; i < numSamples; i++ {
		if node.SampleHistory[i] != 0 {
			//weight the values of the sampleHistory when added to the sum variable based on the speed, so older values are weighted less
			sum += node.SampleHistory[i] - ((decreaseRatio) * float32(i))
		} else {
			sum += 0
		}
	}
	node.Avg = sum / float32(numSamples)
}

func (node *NodeImpl) getDriftSlope() (float32, float32){
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

	for i:= range node.SampleHistory {
		ySum += float32(i)
	}
	yAvg = ySum / float32(len(node.SampleHistory))
	for i := range node.SampleHistory {
		sum += (node.SampleHistory[i] - node.Avg) * (float32(i) - yAvg)
		squareSumX += math.Pow( float64(node.SampleHistory[i] - node.Avg), 2)
		squareSumY += math.Pow( float64(i - 1), 2)

		xSum += node.SampleHistory[i]
		xySum += node.SampleHistory[i] * float32(i)
		xSqrSum += float32(math.Pow(float64(node.SampleHistory[i]), 2))
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
func (node *NodeImpl) IncrementTotalSamples() {
	node.TotalSamples++
}

//getter function for average
func (node *NodeImpl) GetAvg() float32 {
	return node.Avg
}

//increases numResets field
func (node *NodeImpl) IncrementNumResets() {
	node.NumResets++
}

//setter function for concentration field
func (node *NodeImpl) SetConcentration(conc float64) {
	node.Concentration = conc
}

//getter function for ID field
func (node *NodeImpl) GetID() int {
	return node.Id
}

//getter function for x and y locations
func (node *NodeImpl) GetLoc() (float32, float32) {
	return node.X, node.Y
}

func (node *NodeImpl) GetLocCoord() Coord {
	return Coord{X: int(node.X), Y: int(node.Y)}
}

func (node *NodeImpl) GetTransformedLocCoord(p *Params) Coord {
	return Coord{X: transformX(int(node.X), p), Y: transformY(int(node.Y), p)}
}

//setter function for S0
func (node *NodeImpl) SetS0(s0 float64) {
	node.S0 = s0
}

//setter function for S1
func (node *NodeImpl) SetS1(s1 float64) {
	node.S1 = s1
}

//setter function for S2
func (node *NodeImpl) SetS2(s2 float64) {
	node.S2 = s2
}

//setter function for E0
func (node *NodeImpl) SetE0(e0 float64) {
	node.E0 = e0
}

//setter function for E1
func (node *NodeImpl) SetE1(e1 float64) {
	node.E1 = e1
}

//setter function for E2
func (node *NodeImpl) SetE2(e2 float64) {
	node.E2 = e2
}

//setter function for ET1
func (node *NodeImpl) SetET1(et1 float64) {
	node.ET1 = et1
}

//setter function for ET2
func (node *NodeImpl) SetET2(et2 float64) {
	node.ET2 = et2
}

//getter function for all parameters
func (node *NodeImpl) GetParams() (float64, float64, float64, float64, float64, float64, float64, float64) {
	return node.S0, node.S1, node.S2, node.E0, node.E1, node.E2, node.ET1, node.ET2
}

//getter function for just S0 - S2 parameters
func (node *NodeImpl) GetCoefficients() (float64, float64, float64) {
	return node.S0, node.S1, node.S2
}

//getter function for x
func (node *NodeImpl) GetX() float32 {
	return node.X
}

//getter function for y
func (node *NodeImpl) GetY() float32 {
	return node.Y
}

//This is the actual upload to the server
func (node *NodeImpl) Server() {
	//getData(&s,curNode.XPos[0:curNode.BufferI],curNode.YPos[0:curNode.BufferI],curNode.Value[0:curNode.BufferI],curNode.Time[0:curNode.BufferI], curNode.Id,curNode.BufferI)
	node.BufferI = 0
}

//Returns node distance to the bomb
func (node *NodeImpl) GeoDist(b Bomb) float32 {
	//this needs to be changed
	return float32(math.Pow(float64(math.Abs(float64(node.X)-float64(b.X))), 2) + math.Pow(float64(math.Abs(float64(node.Y)-float64(b.Y))), 2))
}

//Returns array of accelerometer speeds recorded for a specific node
func (node *NodeImpl) GetSpeed() []float32 {
	return node.AccelerometerSpeed
}

//Returns a different version of the distance to the bomb
func (node *NodeImpl) GetValue() int {
	return int(math.Sqrt(math.Pow(float64(int(node.X)-node.P.B.X), 2) + math.Pow(float64(node.Y-float32(node.P.B.Y)), 2)))
}


func (node *NodeImpl) Distance(b Bomb) float32 {
	return float32(math.Sqrt(math.Pow(float64(math.Abs(float64(node.X)-float64(b.X))),2) + math.Pow(float64(math.Abs(float64(node.Y)-float64(b.Y))),2)))
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

func transformCoord(c Coord, p *Params) Coord {
	return Coord{X: transformX(c.X, p), Y: transformY(c.Y, p)}
}

func transformX(x int, p *Params) int {
	return int(float32(p.FineWidth/2) -  ((float32(p.B.X - x)/2.0)*float32(p.FineScale)))
	//return float32(p.FineWidth/2) -  ((float32(p.B.X - x)/2.0)*float32(p.FineScale))

}

func transformY(y int, p *Params) int {
	return int(float32(p.FineHeight/2) -  ((float32(p.B.Y - y)/2.0)*float32(p.FineScale)))
	//return float32(p.FineHeight/2) -  ((float32(p.B.Y - y)/2.0)*float32(p.FineScale))
}


func transformXF(x float32, p *Params) float32 {
	//return int(float32(p.FineWidth/2) -  ((float32(p.B.X - x)/2.0)*float32(p.FineScale)))
	return float32(int(p.FineWidth/2)) -  ((float32(float32(p.B.X) - x)/2.0)*float32(p.FineScale))

}

func transformYF(y float32, p *Params) float32 {
	//return int(float32(p.FineHeight/2) -  ((float32(p.B.Y - y)/2.0)*float32(p.FineScale)))
	return float32(int(p.FineHeight/2)) -  ((float32(float32(p.B.Y) - y)/2.0)*float32(p.FineScale))
}


func interpolateReading(x , y float32, time, timeStep int, fine bool, p *Params) float64{

	transX := transformXF(x, p)
	transY := transformYF(y, p)


	oldX := int(transX)
	oldY := int(transY)
	nextX := int(math.Ceil(float64(transX)))
	nextY := int(math.Ceil(float64(transY)))




	tl := 0.0
	tr := 0.0
	bl := 0.0
	br := 0.0

	if fine {

		if oldX >= p.FineWidth || oldX < 0 {
			return -1.0
		} else if oldY >= p.FineHeight || oldY < 0 {
			return -1.0
		} else if nextX >= p.FineWidth || nextX < 0 {
			return -1.0
		} else if nextY >= p.FineHeight || nextY < 0 {
			return -1.0
		}

		tl = p.FineSensorReadings[oldX][nextY][timeStep]
		tr = p.FineSensorReadings[nextX][nextY][timeStep]
		bl = p.FineSensorReadings[oldX][oldY][timeStep]
		br = p.FineSensorReadings[nextX][oldY][timeStep]
	} else {
		tl = p.SensorReadings[oldX][nextY][timeStep]
		tr = p.SensorReadings[nextX][nextY][timeStep]
		bl = p.SensorReadings[oldX][oldY][timeStep]
		br = p.SensorReadings[nextX][oldY][timeStep]
	}


	xPortion := 1.0
	if nextX != oldX {
		xPortion = float64(float32(transX - float32(oldX))/float32(float32(nextX) - float32(oldX)))
	}

	botInter := bl + xPortion * (br - bl)
	topInter := tl + xPortion * (tr - tl)

	yPortion := 1.0
	if nextY != oldY {
		yPortion = float64(float32(transY - float32(oldY))/float32(float32(nextY) - float32(oldY)))
	}


	return topInter + yPortion * (botInter - topInter)

}


func trueInterpolate(x , y float32, time, timeStep int, fine bool, p *Params) float32{

	old := interpolateReading(x, y, time, timeStep, fine, p)

	if(timeStep+1 < p.MaxTimeStep) {
		next := interpolateReading(x, y, time, timeStep+1, fine, p)

		floatTime := float32(time)/1000
		oldTime := p.SensorTimes[timeStep]
		nextTime := p.SensorTimes[timeStep+1]

		portionTime := float64((floatTime - float32(oldTime))/float32(nextTime - oldTime))

		return float32(old + portionTime * (next - old))

	}

	return float32(old)
}


//Takes cares of taking a node's readings and printing detections and stuff
func (node *NodeImpl) GetReadings() {


	if node.Valid { //Check if node should actually take readings or if it hasn't shown up yet
		newX, newY := node.GetLoc()

		//RawConc := RawConcentration(curNode.Distance(*curNode.P.B)/2) //this is the node's reported Value without error

		RawConc := 0.0


		if node.Distance(*node.P.B)/2 < float32((node.P.FineWidth/2)/node.P.FineScale) {
			RawConc = float64(trueInterpolate(newX, newY, node.P.CurrentTime, node.P.TimeStep, true, node.P))

			if RawConc == -1.0 {
				RawConc = 0.0
			}

		}

		/*if(curNode.P.RecalReject && ((curNode.P.CurrentTime/1000) - curNode.NodeTime) < 2) {

		} else {
			curNode.report(RawConc)
		}*/
		if ((node.P.CurrentTime/1000) - node.NodeTime) < 2 {
			//skip
		} else {
			node.report(RawConc)
		}

	}
	node.P.Events.Push(&Event{node, SENSE, node.P.CurrentTime + 500, 0})

}


//Takes cares of taking a node's readings and printing detections and stuff
func (node *NodeImpl) GetReadingsCSV() {

	if node.Valid { //check if node has shown up yet

		newX, newY := node.GetLoc()


		RawConcentration := 0.0
		if node.Distance(*node.P.B)/2 < float32((node.P.FineWidth/2)/node.P.FineScale) {
			//fmt.Printf("\n %v %v %v %v", curNode.P.B.X, curNode.P.B.Y, curNode.Distance(*curNode.P.B)/2, float32((curNode.P.FineWidth/2)/curNode.P.FineScale))
			RawConcentration = float64(trueInterpolate(newX, newY, node.P.CurrentTime, node.P.TimeStep, true, node.P))
			if RawConcentration == -1.0 {
				RawConcentration = float64(trueInterpolate(newX, newY, node.P.CurrentTime, node.P.TimeStep, false, node.P))

				//RawConcentration = 0.0
			}

		} else {
			RawConcentration = float64(trueInterpolate(newX, newY, node.P.CurrentTime, node.P.TimeStep, false, node.P))
		}

		/*if(curNode.P.RecalReject && ((curNode.P.CurrentTime/1000) - curNode.NodeTime) < 2) {

		} else {
			curNode.report(RawConcentration)
		}*/

		if ((node.P.CurrentTime/1000) - node.NodeTime) < 2 {
			//skip
		} else {
			node.report(RawConcentration)
		}
		//curNode.report(RawConcentration)

	}
	node.P.Events.Push(&Event{node, SENSE, node.P.CurrentTime + 500, 0})
}

func (node *NodeImpl) GetSensor() {

	if node.Valid { //check if node has shown up yet


		RawConcentration := 0.0
		//need to get the correct Time reading Value from system
		//need to verify where we read from


		/*if(curNode.P.RecalReject && ((curNode.P.CurrentTime/1000) - curNode.NodeTime) < 2) {

		} else {
			curNode.report(RawConcentration)
		}*/
		if ((node.P.CurrentTime/1000) - node.NodeTime) < 2 {
			//skip
		} else {
			node.report(RawConcentration)
		}
		//curNode.report(RawConcentration)

	}
	node.P.Events.Push(&Event{node, SENSE, node.P.CurrentTime + 500, 0})
}

func (node *NodeImpl) report(rawConc float64) {


	newX, newY := node.GetLoc()

	S0, S1, S2, E0, E1, E2, ET1, ET2 := node.GetParams()
	sError := (S0 + E0) + (S1+E1)*math.Exp(-float64(((node.P.CurrentTime/1000)-node.NodeTime))/(node.P.Tau1+ET1)) + (S2+E2)*math.Exp(-float64(((node.P.CurrentTime/1000)-node.NodeTime))/(node.P.Tau2+ET2))
	node.Sensitivity = S0 + (S1)*math.Exp(-float64(((node.P.CurrentTime/1000)-node.NodeTime))/node.P.Tau1) + (S2)*math.Exp(-float64(((node.P.CurrentTime/1000)-node.NodeTime))/node.P.Tau2)
	//sNoise := rand.NormFloat64()*float64(curNode.P.ADCWidth)*curNode.P.ErrorModifierCM + float64(rawConc)*sError
	//sNoise := rand.NormFloat64()*100*curNode.P.ErrorModifierCM + float64(rawConc)*sError
	sNoise := rand.NormFloat64()*math.Sqrt(3.0) + float64(rawConc)*sError
	errorDist := sNoise / node.Sensitivity //this is the node's actual reading with error
	clean := float64(rawConc) / node.Sensitivity


	ADCRead := float64(node.ADCReading(float32(errorDist)))
	ADCClean := float64(node.ADCReading(float32(clean)))



	d := node.Distance(*node.P.B)/2
	/*if d < 10 {
		fmt.Fprintln(curNode.P.MoveReadingsFile, "Time:", curNode.P.CurrentTime/1000, "ID:", curNode.Id, "X:", newX, "Y:",  newY, "Dist:", d, "ADCClean:", ADCClean, "ADCError:", ADCRead, "CleanSense:", clean, "Error:", errorDist, "Raw:", rawConc)
	}*/

	//increment node Time
	//curNode.NodeTime++

	//if curNode.HasCheckedSensor {
	node.IncrementTotalSamples()
	node.UpdateHistory(float32(errorDist))
	//}

	//If the reading is more than 2 standard deviations away from the grid average, then recalibrate
	//gridAverage := curNode.P.Grid[curNode.Row(curNode.P.YDiv)][curNode.Col(curNode.P.XDiv)].Avg
	//standDev := grid[curNode.Row(yDiv)][curNode.Col(xDiv)].StdDev

	//New condition added: also recalibrate when the node's sensitivity is <= 1/10 of its original sensitvity
	//New condition added: Check to make sure the sensor was pinged this iteration
	if node.Sensitivity <= node.InitialSensitivity / 2  && node.P.Iterations_used != 0 {
		fmt.Fprintf(node.P.DriftExploreFile, "ID: %v T: %v In: %v CUR: %v NT: %v RECAL\n", node.Id, node.P.CurrentTime, node.InitialSensitivity, node.Sensitivity, node.NodeTime)
		node.Recalibrate()
		node.Recalibrated = true
		node.IncrementNumResets()
	}

	//printing statements to log files, only if the sensor was pinged this iteration
	//if curNode.HasCheckedSensor && nodesPrint{
	if node.P.NodesPrint {
		if node.Recalibrated {
			fmt.Fprintln(node.P.NodeFile, "ID:", node.GetID(), "Average:", node.GetAvg(), "Reading:", rawConc, "Error Reading:", errorDist, "Recalibrated")
		} else {
			fmt.Fprintln(node.P.NodeFile, "ID:", node.GetID(), "Average:", node.GetAvg(), "Reading:", rawConc, "Error Reading:", errorDist)
		}
		//fmt.Fprintln(nodeFile, "battery:", int(curNode.Battery),)
		node.Recalibrated = false
	}


	inWind := node.P.Server.CheckFalsePosWind(node)           //true if in sensor area
	inRange := float64(d*2) < node.P.DetectionDistance        //true = out
	highConcentration := ADCClean > node.P.DetectionThreshold //true reading of the sensor
	highSensor := ADCRead > node.P.DetectionThreshold         //error model influenced reading of the sensor

	tp := false

	if inRange && highConcentration && highSensor {
		fmt.Fprintln(node.P.DetectionFile, fmt.Sprintf("TP T: %v ID: %v (%v, %v) D: %v C: %v E: %v SE: %.3f S: %.3f R: %.3f", node.P.CurrentTime, node.Id, node.X, node.Y, d, ADCClean, ADCRead, sError, node.Sensitivity, rawConc))
		tp = true
	} else if inRange && highConcentration && !highSensor {
		fmt.Fprintln(node.P.DetectionFile, fmt.Sprintf("FN Drift T: %v ID: %v (%v, %v) D: %v C: %v E: %v SE: %.3f S: %.3f R: %.3f", node.P.CurrentTime, node.Id, node.X, node.Y, d, ADCClean, ADCRead, sError, node.Sensitivity, rawConc))
	} else if inRange && !highConcentration && highSensor {
		fmt.Fprintln(node.P.DetectionFile, fmt.Sprintf("FP Drift T: %v ID: %v (%v, %v) D: %v C: %v E: %v SE: %.3f S: %.3f R: %.3f", node.P.CurrentTime, node.Id, node.X, node.Y, d, ADCClean, ADCRead, sError, node.Sensitivity, rawConc))
	} else if inRange && !highConcentration && !highSensor {
		if inWind == 1 && !node.P.CSVSensor {
			//outside bomb range and the bomb is random , this isn't a real FN
		} else if inWind == 1 && node.P.CSVSensor{
			//we are not  in the wind area, and the bomb isn't random, this is a FN due to wind
			fmt.Fprintln(node.P.DetectionFile, fmt.Sprintf("FN Wind T: %v ID: %v (%v, %v) D: %v C: %v E: %v SE: %.3f S: %.3f R: %.3f", node.P.CurrentTime, node.Id, node.X, node.Y, d, ADCClean, ADCRead, sError, node.Sensitivity, rawConc))
		} else if inWind == 0 && !node.P.CSVSensor {
			//we are in the wind zone and the bomb is random, we are
			//therefore inside the detection range but this can't happen as we would ahve a high concentration
		} else if inWind == 0 && node.P.CSVSensor {
			//we are in the wind zone and the bomb is random, so it isn't possible to get here....
			//fmt.Printf("\n %v %v %v %v %v %v\n", curNode.Id, curNode.P.TimeStep, inWind, curNode.P.CSVSensor, highSensor, highConcentration)
			fmt.Fprintln(node.P.DetectionFile, fmt.Sprintf("FN Wind T: %v ID: %v (%v, %v) D: %v C: %v E: %v SE: %.3f S: %.3f R: %.3f", node.P.CurrentTime, node.Id, node.X, node.Y, d, ADCClean, ADCRead, sError, node.Sensitivity, rawConc))
		}
	} else if !inRange && highConcentration && highSensor {
		if inWind == 0 {
			//we are in a wind zone, therefore this FP is caused by wind...not possible to have in a random
			fmt.Fprintln(node.P.DetectionFile, fmt.Sprintf("FP Wind T: %v ID: %v (%v, %v) D: %v C: %v E: %v SE: %.3f S: %.3f R: %.3f", node.P.CurrentTime, node.Id, node.X, node.Y, d, ADCClean, ADCRead, sError, node.Sensitivity, rawConc))
		}
	} else if !inRange && highConcentration && !highSensor {
		fmt.Fprintln(node.P.DetectionFile, fmt.Sprintf("FN Drift T: %v ID: %v (%v, %v) D: %v C: %v E: %v SE: %.3f S: %.3f R: %.3f", node.P.CurrentTime, node.Id, node.X, node.Y, d, ADCClean, ADCRead, sError, node.Sensitivity, rawConc))
	} else if !inRange && !highConcentration && highSensor {
		fmt.Fprintln(node.P.DetectionFile, fmt.Sprintf("FP Drift T: %v ID: %v (%v, %v) D: %v C: %v E: %v SE: %.3f S: %.3f R: %.3f", node.P.CurrentTime, node.Id, node.X, node.Y, d, ADCClean, ADCRead, sError, node.Sensitivity, rawConc))
	} else if !inRange && !highConcentration && !highSensor {
		//true negative
	}

	/*
	if ADCRead > curNode.P.DetectionThreshold && ADCClean < curNode.P.DetectionThreshold && float64(d*2) > curNode.P.DetectionDistance{
		fmt.Fprintln(curNode.P.DetectionFile, fmt.Sprintf("FP Drift T: %v ID: %v (%v, %v) D: %v C: %v E: %v SE: %.3f S: %.3f R: %.3f", curNode.P.CurrentTime, curNode.Id, curNode.X, curNode.Y, d, ADCClean, ADCRead, sError, curNode.Sensitivity, rawConc))
	} else if ADCRead < curNode.P.DetectionThreshold && ADCClean > curNode.P.DetectionThreshold && float64(d*2) < curNode.P.DetectionDistance {
		fmt.Fprintln(curNode.P.DetectionFile, fmt.Sprintf("FN Drift T: %v ID: %v (%v, %v) D: %v C: %v E: %v SE: %.3f S: %.3f R: %.3f", curNode.P.CurrentTime, curNode.Id, curNode.X, curNode.Y, d, ADCClean, ADCRead, sError, curNode.Sensitivity, rawConc))
	} else if ADCRead < curNode.P.DetectionThreshold && ADCClean < curNode.P.DetectionThreshold && float64(d*2) < curNode.P.DetectionDistance {
		fmt.Fprintln(curNode.P.DetectionFile, fmt.Sprintf("FN Wind T: %v ID: %v (%v, %v) D: %v C: %v E: %v SE: %.3f S: %.3f R: %.3f", curNode.P.CurrentTime, curNode.Id, curNode.X, curNode.Y, d, ADCClean, ADCRead, sError, curNode.Sensitivity, rawConc))
	} else if ADCRead > curNode.P.DetectionThreshold && ADCClean > curNode.P.DetectionThreshold && float64(d*2) < curNode.P.DetectionDistance {
		fmt.Fprintln(curNode.P.DetectionFile, fmt.Sprintf("TP T: %v ID: %v (%v, %v) D: %v C: %v E: %v SE: %.3f S: %.3f R: %.3f", curNode.P.CurrentTime, curNode.Id, curNode.X, curNode.Y, d, ADCClean, ADCRead, sError, curNode.Sensitivity, rawConc))
	} else if ADCRead > curNode.P.DetectionThreshold && ADCClean > curNode.P.DetectionThreshold && float64(d*2) > curNode.P.DetectionDistance {
		fmt.Fprintln(curNode.P.DetectionFile, fmt.Sprintf("FP Wind T: %v ID: %v (%v, %v) D: %v C: %v E: %v SE: %.3f S: %.3f R: %.3f", curNode.P.CurrentTime, curNode.Id, curNode.X, curNode.Y, d, ADCClean, ADCRead, sError, curNode.Sensitivity, rawConc))
	}*/

	//Receives the node's distance and calculates its running average
	//for that square
	//Only do this if the sensor was pinged this iteration

	if node.Valid {
		if node.P.ClusteringOn && node.IsClusterMember && !highSensor {
			node.SendToClusterHead(Reading{ADCRead, newX, newY, node.P.CurrentTime, node.GetID()}, tp)
		} else {
			node.SendToServer(Reading{ADCRead, newX, newY, node.P.CurrentTime, node.GetID()}, tp)
		}
	}
}

func interpolate (start int, end int, portion float32) float32{
	return float32(end-start) * portion + float32(start)
}

//HandleMovementCSV does the same as HandleMovement
func (node *NodeImpl) MoveCSV(p *Params) {
	//time := p.Iterations_used
	floatTemp := float32(p.CurrentTime)
	intTime := int(floatTemp/1000)
	portion := (floatTemp / 1000) - float32(intTime)

	id := node.GetID()

	if node.Valid {
		oldX, oldY := node.GetLoc()
		p.BoolGrid[int(oldX)][int(oldY)] = false //set the old spot false since the node will now move away

		node.X = interpolate(p.NodeMovements[id][intTime-p.MovementOffset].X, p.NodeMovements[id][intTime-p.MovementOffset+1].X, portion)
		node.Y = interpolate(p.NodeMovements[id][intTime-p.MovementOffset].Y, p.NodeMovements[id][intTime-p.MovementOffset+1].Y, portion)

		//set the new location in the boolean field to true
		newX, newY := node.GetLoc()
		//fmt.Println(oldX, oldY,newX, newY, curNode.Id, p.CurrentTime,p.NodeMovements[id][intTime].X, p.NodeMovements[id][intTime+1].X)

		if !node.InBounds(p) {
			//fmt.Println(oldX, oldY,newX, newY, curNode.Id, p.CurrentTime,p.NodeMovements[id][intTime].X, p.NodeMovements[id][intTime+1].X)
			node.Valid = false
			if p.ClusteringOn {
				p.NodeTree.RemoveAndClean(node)
			}

		} else {

			d := node.Distance(*p.B)/2
			if int(d) < p.MinDistance {
				p.MinDistance = int(d)
				fmt.Fprintf(p.DistanceFile, "ID: %v T: %v D: %v\n", node.Id, intTime, int(d))
			}

			p.BoolGrid[int(newX)][int(newY)] = true

			//sync the QuadTree
			if p.ClusteringOn {
				p.NodeTree.NodeMovement(node)
				p.ClusterNetwork.ClusterMovement(node, p.NodeBTRange)
			}
		}
	}


	if !node.Valid {
		node.Valid = node.IsValid(p.NodeMovements[id][intTime-p.MovementOffset].X, p.NodeMovements[id][intTime-p.MovementOffset].Y, p)
		node.X = float32(p.NodeMovements[id][intTime-p.MovementOffset].X)
		node.Y = float32(p.NodeMovements[id][intTime-p.MovementOffset].Y)
		if node.Valid {
			if p.DriftExplorer {
				node.NodeTime = RandomInt(-7000, 0)
			} else {
				//curNode.NodeTime = 0
				node.NodeTime = RandomInt(-7000, 0)
			}
		}
	}

}

//HandleMovement adjusts BoolGrid when nodes move around the map
func (node *NodeImpl) MoveNormal(p *Params) {

	oldX, oldY := node.GetLoc()
	p.BoolGrid[int(oldX)][int(oldY)] = false //set the old spot false since the node will now move away

	//move the node to its new location
	node.Move(p)

	//set the new location in the boolean field to true
	newX, newY := node.GetLoc()
	p.BoolGrid[int(newX)][int(newY)] = true



	//Add the node into its new Square's p.TotalNodes
	//If the node hasn't left the square, that Square's p.TotalNodes will
	//remain the same after these calculations

}

func rangeInt(min, max int) int { //returns a random number between max and min
	return rand.Intn(max-min) + min
}
