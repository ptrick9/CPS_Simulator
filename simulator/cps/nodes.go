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
	BatteryLossDynamic(p *Params)   //Battery loss based of ratios of battery usage
	BatteryLossDynamic1(p *Params)  //2 stage buffer battery loss
	UpdateHistory(newValue float32) //updates history of node's samples
	IncrementTotalSamples()         //increments total number of samples node has taken
	GetAvg() float32                //returns average of node's past samples
	IncrementNumResets()            //increments the number of times a node has been reset
	SetConcentration(conc float64)  //sets the concentration of a node
	GeoDist(b Bomb) float32         //returns distance from bomb (rather than reading of node)
	GetID() int                     //returns ID of node
	GetLoc() (x, y int)             //returns x and y values of node

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
	Id                              int      //id of node
	OldX                            int      // for movement
	OldY                            int      // for movement
	Sitting                         int      // for movement
	X                               int      //x pos of node
	Y                               int      //y pos of node
	Battery                         float32  //battery of node
	BatteryLossScalar               float32  //natural incremental battery loss of node
	BatteryLossCheckingSensorScalar float32  //sensor based battery loss of node
	BatteryLossGPSScalar            float32  //GPS based battery loss of node
	BatteryLossCheckingServerScalar float32  //server communication based battery loss of node
	ToggleCheckIterator             int      //node's personal iterator mostly for cascading pings
	HasCheckedSensor                bool     //did the node just ping the sensor?
	TotalChecksSensor               int      //total sensor pings of node
	HasCheckedGPS                   bool     //did the node just ping the GPS?
	TotalChecksGPS                  int      //total GPS pings of node
	HasCheckedServer                bool     //did the node just communicate with the server?
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
	XPos                            [100]int //x pos buffer of node
	YPos                            [100]int //y pos buffer of node
	Value                           [100]int //value buffer of node
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

//Returns the x index of the square in which the specific
//	node currently resides
func (n *NodeImpl) Row(div int) int {
	return n.Y / div
}

//Returns the y index of the square in which the specific
//	node currently resides
func (n *NodeImpl) Col(div int) int {
	return n.X / div
}

//Returns a float representing the detection of the bomb
//	by the specific node depending on distance
func (n *NodeImpl) Distance(b Bomb) float32 {
	dist := float32(math.Pow(float64(math.Abs(float64(n.X)-float64(b.X))), 2) + math.Pow(float64(math.Abs(float64(n.Y)-float64(b.Y))), 2))

	if dist == 0 {
		return 1000
	} else {
		//return float32(1000 / (math.Pow((float64(dist)/0.2)*0.25,1.5)))
		return float32(math.Pow(1000/float64(dist), 1.5))
	}
}

// These are the toString methods for battery levels
func (n Bn) String() string { // extra extra string statements
	return fmt.Sprintf("x: %v y: %v Xspeed: %v Yspeed: %v id: %v battery: %v sensor checked: %v sensor checks: %v GPS checked: %v GPS checks: %v server checked: %v server checks: %v buffer: %v ", n.X, n.Y, n.X_speed, n.Y_speed, n.Id, n.Battery, n.HasCheckedSensor, n.TotalChecksSensor, n.HasCheckedGPS, n.TotalChecksGPS, n.HasCheckedServer, n.TotalChecksServer, n.BufferI)
}

func (n Wn) String() string {
	return fmt.Sprintf("x: %v y: %v speed: %v dir: %v id: %v battery: %v sensor checked: %v sensor checks: %v GPS checked: %v GPS checks: %v server checked: %v server checks: %v buffer: %v ", n.X, n.Y, n.Speed, n.Dir, n.Id, n.Battery, n.HasCheckedSensor, n.TotalChecksSensor, n.HasCheckedGPS, n.TotalChecksGPS, n.HasCheckedServer, n.TotalChecksServer, n.BufferI)
}

func (n Rn) String() string {
	return fmt.Sprintf("x: %v y: %v id: %v battery: %v sensor checked: %v sensor checks: %v GPS checked: %v GPS checks: %v server checked: %v server checks: %v buffer: %v ", n.X, n.Y, n.Id, n.Battery, n.HasCheckedSensor, n.TotalChecksSensor, n.HasCheckedGPS, n.TotalChecksGPS, n.HasCheckedServer, n.TotalChecksServer, n.BufferI)
} // end extra extra string statements

func (n NodeImpl) String() string {
	//return fmt.Sprintf("x: %v y: %v id: %v battery: %v sensor checked: %v sensor checks: %v GPS checked: %v GPS checks: %v server checked: %v server checks: %v buffer: %v ", n.X, n.Y, n.Id, n.Battery, n.HasCheckedSensor, n.TotalChecksSensor, n.HasCheckedGPS, n.TotalChecksGPS, n.HasCheckedServer, n.TotalChecksServer,n.BufferI)
	return fmt.Sprintf("battery: %v sensor checked: %v GPS checked: %v ", int(n.Battery), n.HasCheckedSensor, n.HasCheckedGPS)

}

func (c Coord) String() string {
	return fmt.Sprintf("{%v %v %v}", c.X, c.Y, c.Time)
}

func (c Coord) Equals(c2 Coord) bool {
	return c.X == c2.X && c.Y == c2.Y
}

func (n *NodeImpl) Move(p *Params) {
	if n.Sitting <= p.SittingStopThresholdCM {
		n.OldX = n.X / p.XDiv
		n.OldY = n.Y / p.XDiv

		var potentialSpots []GridSpot

		//only add the ones that are valid to move to into the list
		if n.Y-1 >= 0 &&
			n.X >= 0 &&
			n.X < len(p.BoardMap[n.Y-1]) &&
			n.Y-1 < len(p.BoardMap) &&

			p.BoardMap[n.Y-1][n.X] != -1 &&
			p.BoolGrid[n.Y-1][n.X] == false { // &&
			//p.BoardMap[n.X][n.Y-1] <= p.BoardMap[n.X][n.Y] {

			up := GridSpot{n.X, n.Y - 1, p.BoardMap[n.Y-1][n.X]}
			potentialSpots = append(potentialSpots, up)
		}
		if n.X+1 < len(p.BoardMap[n.Y]) &&
			n.X+1 >= 0 &&
			n.Y < len(p.BoardMap) &&
			n.Y >= 0 &&

			p.BoardMap[n.Y][n.X+1] != -1 &&
			p.BoolGrid[n.Y][n.X+1] == false { // &&
			//p.BoardMap[n.X+1][n.Y] <= p.BoardMap[n.X][n.Y] {

			right := GridSpot{n.X + 1, n.Y, p.BoardMap[n.Y][n.X+1]}
			potentialSpots = append(potentialSpots, right)
		}
		if n.Y+1 < len(p.BoardMap) &&
			n.Y+1 >= 0 &&
			n.X < len(p.BoardMap[n.Y+1]) &&
			n.X >= 0 &&

			p.BoardMap[n.Y+1][n.X] != -1 &&
			p.BoolGrid[n.Y+1][n.X] == false { //&&
			//p.BoardMap[n.X][n.Y+1] <= p.BoardMap[n.X][n.Y] {

			down := GridSpot{n.X, n.Y + 1, p.BoardMap[n.Y+1][n.X]}
			potentialSpots = append(potentialSpots, down)
		}
		if n.X-1 >= 0 &&
			n.X-1 < len(p.BoardMap[n.Y]) &&
			n.Y >= 0 &&
			n.Y < len(p.BoardMap) &&

			p.BoardMap[n.Y][n.X-1] != -1 &&
			p.BoolGrid[n.Y][n.X-1] == false { // &&
			//p.BoardMap[n.X-1][n.Y] <= p.BoardMap[n.X][n.Y] {

			left := GridSpot{n.X - 1, n.Y, p.BoardMap[n.Y][n.X-1]}
			potentialSpots = append(potentialSpots, left)
		}

		sort.Sort(byRandom(potentialSpots))
		sort.Sort(byRandom(potentialSpots))
		sort.Sort(byRandom(potentialSpots))
		sort.Sort(byValue(potentialSpots))

		/*for i := 0; i < len(potentialSpots); i++ {
			if p.Grid[potentialSpots[i].Y/p.YDiv][potentialSpots[i].X/p.XDiv].ActualNumNodes <= p.SquareCapacity {
				n.X = potentialSpots[i].X
				n.Y = potentialSpots[i].Y
				break
			}
		}*/

		//If there are no potential spots, do not move
		if len(potentialSpots) > 0 {
			n.X = potentialSpots[0].X
			n.Y = potentialSpots[0].Y
		}

		//Change number of nodes in square
		/*if n.X/p.XDiv != n.OldX || n.Y/p.YDiv != n.OldY {
			p.Grid[n.Y/p.YDiv][n.X/p.XDiv].ActualNumNodes = p.Grid[n.Y/p.YDiv][n.X/p.XDiv].ActualNumNodes + 1
			p.Grid[n.OldY][n.OldX].ActualNumNodes = p.Grid[n.OldY][n.OldX].ActualNumNodes - 1
		}*/

		p.Server.UpdateSquareNumNodes()
		if n.Diffx == 0 && n.Diffy == 0 || n.Sitting < 0 {
			n.Sitting = n.Sitting + 1
		} else {
			n.Sitting = 0
		}
	}
}

func (n *NodeImpl) Recalibrate() {
	n.Sensitivity = n.InitialSensitivity
	n.NodeTime = 0
}

//Moves the bouncing node
func (n *Bn) Move(p *Params) {
	//Boundary conditions
	if n.X+n.X_speed < p.MaxX && n.X+n.X_speed >= 0 {
		n.X = n.X + n.X_speed
	} else {
		if n.X+n.X_speed >= p.MaxX {
			n.X = n.X - (n.X_speed - (p.MaxX - 1 - n.X))
			n.X_speed = n.X_speed * -1
		} else {
			n.X = (n.X_speed + n.X) * -1
			n.X_speed = n.X_speed * -1
		}
	}
	if n.Y+n.Y_speed < p.MaxY && n.Y+n.Y_speed >= 0 {
		n.Y = n.Y + n.Y_speed
	} else {
		if n.Y+n.Y_speed >= p.MaxY {
			n.Y = n.Y - (n.Y_speed - (p.MaxY - 1 - n.Y))
			n.Y_speed = n.Y_speed * -1
		} else {
			n.Y = (n.Y_speed + n.Y) * -1
			n.Y_speed = n.Y_speed * -1
		}
	}
}

//Moves the wall nodes
func (n *Wn) Move(p *Params) {
	if n.Dir == 0 { //x-axis
		//Boundary conditions
		if n.X+n.Speed < p.MaxX && n.X+n.Speed >= 0 {
			n.X = n.X + n.Speed
		} else {
			if n.X+n.Speed >= p.MaxX {
				n.X = n.X - (n.Speed - (p.MaxX - 1 - n.X))
				n.Speed = n.Speed * -1
			} else {
				n.X = (n.Speed + n.X) * -1
				n.Speed = n.Speed * -1
			}
		}
	} else { //y-axis
		if n.Y+n.Speed < p.MaxY && n.Y+n.Speed >= 0 {
			n.Y = n.Y + n.Speed
		} else {
			if n.Y+n.Speed >= p.MaxY {
				n.Y = n.Y - (n.Speed - (p.MaxY - 1 - n.Y))
				n.Speed = n.Speed * -1
			} else {
				n.Y = (n.Speed + n.Y) * -1
				n.Speed = n.Speed * -1
			}
		}
	}
}

//Moves the random nodes
func (n *Rn) Move(p *Params) {
	x_speed := rangeInt(-3, 3)
	y_speed := rangeInt(-3, 3)

	//Boundary conditions
	if n.X+x_speed < p.MaxX && n.X+x_speed >= 0 {
		n.X = n.X + x_speed
	} else {
		if n.X+x_speed >= p.MaxX {
			n.X = n.X - (x_speed - (p.MaxX - 1 - n.X))
		} else {
			n.X = (x_speed + n.X) * -1
		}
	}
	if n.Y+y_speed < p.MaxY && n.Y+y_speed >= 0 {
		n.Y = n.Y + y_speed
	} else {
		if n.Y+y_speed >= p.MaxY {
			n.Y = n.Y - (y_speed - (p.MaxY - 1 - n.Y))
		} else {
			n.Y = (y_speed + n.Y) * -1
		}
	}
}

//Returns the arr with the element at index n removed
func Remove_index(arr []Path, n int) []Path {
	return arr[:n+copy(arr[n:], arr[n+1:])]
}

//Returns the array with the range of elements from index
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
//	index n
func Insert_array(arr1 []Coord, arr2 []Coord, n int) []Coord {
	if len(arr1) == 0 {
		return arr2
	} else {
		return append(arr1[:n], append(arr2, arr1[n:]...)...)
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

// This is the battery loss function that clears the buffer at 2 different rates based on the battery percentage left
func (n *NodeImpl) BatteryLossDynamic1(p *Params) {
	n.HasCheckedGPS = false
	n.HasCheckedSensor = false
	n.HasCheckedServer = false

	// This is the manual buffer clearing
	if n.BufferI >= p.MaxBufferCapacityCM {
		n.Server()
	}
	// This is a generic iterator
	n.ToggleCheckIterator = n.ToggleCheckIterator + 1
	// These are the respective accelerometer positions
	n.Current = n.ToggleCheckIterator % 3
	n.Previous = (n.ToggleCheckIterator - 1) % 3
	if n.Current == 0 {
		n.AccelerometerPosition[0][0] = n.X
		n.AccelerometerPosition[1][0] = n.Y
	} else if n.Current == 1 {
		n.AccelerometerPosition[0][1] = n.X
		n.AccelerometerPosition[1][1] = n.Y
	} else if n.Current == 2 {
		n.AccelerometerPosition[0][2] = n.X
		n.AccelerometerPosition[1][2] = n.Y
	}
	n.Diffx = n.AccelerometerPosition[0][n.Current] - n.AccelerometerPosition[0][n.Previous]
	n.Diffy = n.AccelerometerPosition[1][n.Current] - n.AccelerometerPosition[1][n.Previous]
	// This is the accelerometer determined speed
	n.Speed = float32(math.Sqrt(float64(n.Diffx*n.Diffx + n.Diffy*n.Diffy)))
	// This is a list of the previous accelerometer determined speeds
	n.AccelerometerSpeed = append(n.AccelerometerSpeed, n.Speed)
	//threshHoldBatteryToHave = thres
	// This is the natural loss of the battery
	if n.Battery > 0 {
		n.Battery = n.Battery - n.BatteryLossScalar
	}
	// This is the predicted natural loss to prevent overages.
	naturalLoss = n.Battery - (float32(p.Iterations_of_event-p.Iterations_used) * n.BatteryLossScalar)
	// This is the algorithm that determines sampling rate's ratios
	n.Pings = n.Battery * p.TotalPercentBatteryToUse / (n.BatteryLossCheckingSensorScalar + n.BatteryLossGPSScalar + n.BatteryLossCheckingServerScalar) // set percentage for consumption here, also '50' if no minus
	n.InverseSensor = 1 / n.BatteryLossCheckingSensorScalar
	n.InverseGPS = 1 / n.BatteryLossGPSScalar
	n.InverseServer = 1 / n.BatteryLossCheckingServerScalar
	n.SensorPings = n.Pings * (n.InverseSensor / (n.InverseServer + n.InverseGPS + n.InverseGPS))
	n.GPSPings = n.Pings * (n.InverseGPS / (n.InverseServer + n.InverseGPS + n.InverseGPS))
	n.ServerPings = n.Pings * (n.InverseServer / (n.InverseServer + n.InverseGPS + n.InverseGPS))

	if naturalLoss > p.ThreshHoldBatteryToHave {
		n.SensorPingPeriod = float32(p.Iterations_of_event-p.Iterations_used) / n.SensorPings
		if n.SensorPingPeriod < 1 {
			n.SensorPingPeriod = 1
		}
		// Checks to see if sensor is pinged
		if (n.ToggleCheckIterator-n.Cascade)%int(n.SensorPingPeriod) == 0 && n.Battery > 1 {
			n.Battery = n.Battery - n.BatteryLossCheckingSensorScalar
			n.TotalChecksSensor = n.TotalChecksSensor + 1
			n.HasCheckedSensor = true
			n.Sense(p)
		} else {
			n.HasCheckedSensor = false
		}
		// Checks to see if GPS is pinged
		n.GPSPingPeriod = float32(p.Iterations_of_event-p.Iterations_used) / n.GPSPings
		if n.GPSPingPeriod < 1 {
			n.GPSPingPeriod = 1
		}
		if ((n.ToggleCheckIterator-n.Cascade)%int(n.GPSPingPeriod) == 0 && n.Battery > 1) || (n.Speed > float32(p.MovementSamplingSpeedCM) && n.ToggleCheckIterator%p.MovementSamplingPeriodCM == 0) { // && n.ToggleCheckIterator%n.SpeedGPSPeriod == 0
			n.Battery = n.Battery - n.BatteryLossGPSScalar
			n.TotalChecksGPS = n.TotalChecksGPS + 1
			n.HasCheckedGPS = true
			n.GPS(p)
		} else {
			n.HasCheckedGPS = false
		}

		// This is the 2 stage buffer based on battery percentages
		if n.Battery >= 75 { //100 - 75 percent
			if (n.ToggleCheckIterator-n.Cascade)%14 == 0 && n.Battery > 1 { // 1000/70 = 14
				n.Battery = n.Battery - n.BatteryLossCheckingServerScalar
				n.TotalChecksServer = n.TotalChecksServer + 1
				n.HasCheckedServer = true
				n.Server()
			} else {
				n.HasCheckedServer = false
			}
		} else if n.Battery >= 30 && n.Battery < 75 { //70 - 30 percent
			if (n.ToggleCheckIterator-n.Cascade)%50 == 0 && n.Battery > 1 { //1000/20 = 50
				n.Battery = n.Battery - n.BatteryLossCheckingServerScalar
				n.TotalChecksServer = n.TotalChecksServer + 1
				n.HasCheckedServer = true
				n.Server()
			} else {
				n.HasCheckedServer = false
			}
		} else {
			n.HasCheckedServer = false
		}
	}
}

//This is the battery loss function where the server sensor and GPS are pinged separately and by their own accord
func (n *NodeImpl) BatteryLossTable(p *Params) {
	n.HasCheckedGPS = false
	n.HasCheckedSensor = false
	n.HasCheckedServer = false

	// This is the buffer limit if the limit is meet the buffer is forcibly cleared
	if n.BufferI >= p.MaxBufferCapacityCM {
		n.Server()
	}
	// This iterator is generic
	n.ToggleCheckIterator = n.ToggleCheckIterator + 1
	// These are the nodes respective accelerometer positions
	n.Current = n.ToggleCheckIterator % 3
	n.Previous = (n.ToggleCheckIterator - 1) % 3
	// this is the accelerometer's functions
	if n.Current == 0 {
		n.AccelerometerPosition[0][0] = n.X
		n.AccelerometerPosition[1][0] = n.Y
	} else if n.Current == 1 {
		n.AccelerometerPosition[0][1] = n.X
		n.AccelerometerPosition[1][1] = n.Y
	} else if n.Current == 2 {
		n.AccelerometerPosition[0][2] = n.X
		n.AccelerometerPosition[1][2] = n.Y
	}
	n.Diffx = n.AccelerometerPosition[0][n.Current] - n.AccelerometerPosition[0][n.Previous]
	n.Diffy = n.AccelerometerPosition[1][n.Current] - n.AccelerometerPosition[1][n.Previous]
	// Speed as determined by accelerometer
	speed := float32(math.Sqrt(float64(n.Diffx*n.Diffx + n.Diffy*n.Diffy)))
	// This keeps track of the accelerometer values
	n.AccelerometerSpeed = append(n.AccelerometerSpeed, speed)
	// natural loss of the battery
	if n.Battery > 0 {
		n.Battery = n.Battery - n.BatteryLossScalar
	}
	// predicted natural loss of the battery
	naturalLoss = n.Battery - (float32(p.Iterations_of_event-p.Iterations_used) * n.BatteryLossScalar)

	// this is the battery loss based on checking the sensor, GPS, and server.
	if naturalLoss > p.ThreshHoldBatteryToHave {
		if (n.ToggleCheckIterator-n.Cascade)%p.SensorSamplingPeriodCM == 0 {
			n.Battery = n.Battery - n.BatteryLossCheckingSensorScalar
			n.TotalChecksSensor = n.TotalChecksSensor + 1
			n.HasCheckedSensor = true
			n.Sense(p)
		} else {
			n.HasCheckedSensor = false
		}
		if (n.ToggleCheckIterator-n.Cascade)%p.GPSSamplingPeriodCM == 0 || (speed > float32(p.MovementSamplingSpeedCM) && n.ToggleCheckIterator%p.MovementSamplingPeriodCM == 0) { // && n.ToggleCheckIterator%n.SpeedGPSPeriod == 0
			n.Battery = n.Battery - n.BatteryLossGPSScalar
			n.TotalChecksGPS = n.TotalChecksGPS + 1
			n.HasCheckedGPS = true
			n.GPS(p)
		} else {
			n.HasCheckedGPS = false
		}

		// Check to ping server
		if (n.ToggleCheckIterator-n.Cascade)%p.ServerSamplingPeriodCM == 0 {
			n.Battery = n.Battery - n.BatteryLossCheckingServerScalar
			n.TotalChecksServer = n.TotalChecksServer + 1
			n.HasCheckedServer = true
			n.Server()
		} else {
			n.HasCheckedServer = false
		}
	}
}

//This is the battery loss function where the server sensor and GPS are pinged separately and by their own accord
func (n *NodeImpl) BatteryLossMostDynamic(p *Params) {
	n.HasCheckedGPS = false
	n.HasCheckedSensor = false
	n.HasCheckedServer = false

	// This is the buffer limit if the limit is meet the buffer is forcibly cleared
	if n.BufferI >= p.MaxBufferCapacityCM {
		n.Server()
	}
	// This iterator is generic
	n.ToggleCheckIterator = n.ToggleCheckIterator + 1
	// These are the nodes respective accelerometer positions
	n.Current = n.ToggleCheckIterator % 3
	n.Previous = (n.ToggleCheckIterator - 1) % 3
	// this is the accelerometer's functions
	if n.Current == 0 {
		n.AccelerometerPosition[0][0] = n.X
		n.AccelerometerPosition[1][0] = n.Y
	} else if n.Current == 1 {
		n.AccelerometerPosition[0][1] = n.X
		n.AccelerometerPosition[1][1] = n.Y
	} else if n.Current == 2 {
		n.AccelerometerPosition[0][2] = n.X
		n.AccelerometerPosition[1][2] = n.Y
	}
	n.Diffx = n.AccelerometerPosition[0][n.Current] - n.AccelerometerPosition[0][n.Previous]
	n.Diffy = n.AccelerometerPosition[1][n.Current] - n.AccelerometerPosition[1][n.Previous]
	// Speed as determined by accelerometer
	n.Speed = float32(math.Sqrt(float64(n.Diffx*n.Diffx + n.Diffy*n.Diffy)))
	// This keeps track of the accelerometer values
	n.AccelerometerSpeed = append(n.AccelerometerSpeed, n.Speed)
	// natural loss of the battery
	if n.Battery > 0 {
		n.Battery = n.Battery - n.BatteryLossScalar
	}
	// predicted natural loss of the battery
	naturalLoss = n.Battery - (float32(p.Iterations_of_event) * n.BatteryLossScalar)

	// This is the ratio algorithm used to determine the rate of pings
	n.InverseSensor = 1 / n.BatteryLossCheckingSensorScalar
	n.InverseGPS = 1 / n.BatteryLossGPSScalar
	n.InverseServer = 1 / n.BatteryLossCheckingServerScalar

	//SensorBatteryToUse := (totalPercentBatteryToUse * (n.InverseSensor / (n.InverseServer + n.InverseGPS + n.InverseSensor)))
	//GPSBatteryToUse := (totalPercentBatteryToUse * (n.InverseGPS / (n.InverseServer + n.InverseGPS + n.InverseSensor)))
	//ServerBatteryToUse := (totalPercentBatteryToUse * (n.InverseServer / (n.InverseServer + n.InverseGPS + n.InverseSensor)))

	n.SensorPings = (p.TotalPercentBatteryToUse * (n.InverseSensor / (n.InverseServer + n.InverseGPS + n.InverseSensor))) / n.BatteryLossCheckingSensorScalar
	n.GPSPings = (p.TotalPercentBatteryToUse * (n.InverseGPS / (n.InverseServer + n.InverseGPS + n.InverseSensor))) / n.BatteryLossGPSScalar
	n.ServerPings = (p.TotalPercentBatteryToUse * (n.InverseServer / (n.InverseServer + n.InverseGPS + n.InverseSensor))) / n.BatteryLossCheckingServerScalar

	// this is the battery loss based on checking the sensor, GPS, and server.
	if naturalLoss > p.ThreshHoldBatteryToHave {
		n.SensorPingPeriod = float32(p.Iterations_of_event) / n.SensorPings //-iterations_used
		if n.SensorPingPeriod < 1 {
			n.SensorPingPeriod = 1
		}
		// Check to ping sensor
		if (n.ToggleCheckIterator-n.Cascade)%int(n.SensorPingPeriod) == 0 && n.Battery > 1 {
			n.Battery = n.Battery - n.BatteryLossCheckingSensorScalar
			n.TotalChecksSensor = n.TotalChecksSensor + 1
			n.HasCheckedSensor = true
			n.Sense(p)
		} else {
			n.HasCheckedSensor = false
		}
		n.GPSPingPeriod = float32(p.Iterations_of_event) / n.GPSPings //-iterations_used
		if n.GPSPingPeriod < 1 {
			n.GPSPingPeriod = 1
		}
		// Check to ping GPS, also note the extra pings made by a high speed.
		if ((n.ToggleCheckIterator-n.Cascade)%int(n.GPSPingPeriod) == 0 && n.Battery > 1) || (n.Speed > float32(p.MovementSamplingSpeedCM) && n.ToggleCheckIterator%p.MovementSamplingPeriodCM == 0) { // && n.ToggleCheckIterator%n.SpeedGPSPeriod == 0
			n.Battery = n.Battery - n.BatteryLossGPSScalar
			n.TotalChecksGPS = n.TotalChecksGPS + 1
			n.HasCheckedGPS = true
			n.GPS(p)
		} else {
			n.HasCheckedGPS = false
		}
		n.ServerPingPeriod = float32(p.Iterations_of_event) / n.ServerPings //-iterations_used
		if n.ServerPingPeriod < 1 {
			n.ServerPingPeriod = 1.1
		} else if int(n.ServerPingPeriod) > p.Iterations_of_event {
			n.ServerPingPeriod = float32(p.Iterations_of_event)
		}
		if n.ToggleCheckIterator-n.Cascade == 0 {
			n.ToggleCheckIterator = n.Cascade + 1
		}
		// Check to ping server
		//fmt.Println(n.ToggleCheckIterator,n.Cascade,n.ServerPingPeriod,n.Id, n.ServerPings,n.BatteryLossCheckingServerScalar, iterations_of_event,float32(iterations_of_event),int(float32(iterations_of_event)))
		if (n.ToggleCheckIterator-n.Cascade)%int(n.ServerPingPeriod) == 0 && n.Battery > 1 {
			n.Battery = n.Battery - n.BatteryLossCheckingServerScalar
			n.TotalChecksServer = n.TotalChecksServer + 1
			n.HasCheckedServer = true
			n.Server()
		} else {
			n.HasCheckedServer = false
		}
	}
}

//This is the battery loss function where the server sensor and GPS are pinged separately and by their own accord
func (n *NodeImpl) BatteryLossDynamic(p *Params) {
	n.HasCheckedGPS = false
	n.HasCheckedSensor = false
	n.HasCheckedServer = false

	// This is the buffer limit if the limit is meet the buffer is forcibly cleared
	if n.BufferI >= p.MaxBufferCapacityCM {
		n.Server()
	}
	// This iterator is generic
	n.ToggleCheckIterator = n.ToggleCheckIterator + 1
	// These are the nodes respective accelerometer positions
	current := n.ToggleCheckIterator % 3
	previous := (n.ToggleCheckIterator - 1) % 3
	// this is the accelerometer's functions
	if current == 0 {
		n.AccelerometerPosition[0][0] = n.X
		n.AccelerometerPosition[1][0] = n.Y
	} else if current == 1 {
		n.AccelerometerPosition[0][1] = n.X
		n.AccelerometerPosition[1][1] = n.Y
	} else if current == 2 {
		n.AccelerometerPosition[0][2] = n.X
		n.AccelerometerPosition[1][2] = n.Y
	}
	diffx := n.AccelerometerPosition[0][current] - n.AccelerometerPosition[0][previous]
	diffy := n.AccelerometerPosition[1][current] - n.AccelerometerPosition[1][previous]
	// Speed as determined by accelerometer
	speed := float32(math.Sqrt(float64(diffx*diffx + diffy*diffy)))
	// This keeps track of the accelerometer values
	n.AccelerometerSpeed = append(n.AccelerometerSpeed, speed)
	// natural loss of the battery
	if n.Battery > 0 {
		n.Battery = n.Battery - n.BatteryLossScalar
	}
	// predicted natural loss of the battery
	naturalLoss = n.Battery - (float32(p.Iterations_of_event-p.Iterations_used) * n.BatteryLossScalar)

	// This is the ratio algorithm used to determine the rate of pings
	n.Pings = n.Battery * .5 / (n.BatteryLossCheckingSensorScalar + n.BatteryLossGPSScalar + n.BatteryLossCheckingServerScalar) // set percentage for consumption here, also '50' if no minus
	n.InverseSensor = 1 / n.BatteryLossCheckingSensorScalar
	n.InverseGPS = 1 / n.BatteryLossGPSScalar
	n.InverseServer = 1 / n.BatteryLossCheckingServerScalar
	n.SensorPings = n.Pings * (n.InverseSensor / (n.InverseServer + n.InverseGPS + n.InverseGPS))
	n.GPSPings = n.Pings * (n.InverseGPS / (n.InverseServer + n.InverseGPS + n.InverseGPS))
	n.ServerPings = n.Pings * (n.InverseServer / (n.InverseServer + n.InverseGPS + n.InverseGPS))

	// this is the battery loss based on checking the sensor, GPS, and server.
	if naturalLoss > p.ThreshHoldBatteryToHave {
		n.SensorPingPeriod = float32(p.Iterations_of_event) / n.SensorPings //-iterations_used
		if n.SensorPingPeriod < 1 {
			n.SensorPingPeriod = 1
		}
		// Check to ping sensor
		if (n.ToggleCheckIterator-n.Cascade)%int(n.SensorPingPeriod) == 0 && n.Battery > 1 {
			n.Battery = n.Battery - n.BatteryLossCheckingSensorScalar
			n.TotalChecksSensor = n.TotalChecksSensor + 1
			n.HasCheckedSensor = true
			n.Sense(p)
		} else {
			n.HasCheckedSensor = false
		}
		n.GPSPingPeriod = float32(p.Iterations_of_event) / n.GPSPings //-iterations_used
		if n.GPSPingPeriod < 1 {
			n.GPSPingPeriod = 1
		}
		// Check to ping GPS, also note the extra pings made by a high speed.
		if ((n.ToggleCheckIterator-n.Cascade)%int(n.GPSPingPeriod) == 0 && n.Battery > 1) || (speed > float32(p.MovementSamplingSpeedCM) && n.ToggleCheckIterator%p.MovementSamplingPeriodCM == 0) { // && n.ToggleCheckIterator%n.SpeedGPSPeriod == 0
			n.Battery = n.Battery - n.BatteryLossGPSScalar
			n.TotalChecksGPS = n.TotalChecksGPS + 1
			n.HasCheckedGPS = true
			n.GPS(p)
		} else {
			n.HasCheckedGPS = false
		}
		n.ServerPingPeriod = float32(p.Iterations_of_event) / n.ServerPings //-iterations_used
		if n.ServerPingPeriod < 1 {
			n.ServerPingPeriod = 1
		}
		// Check to ping server
		if (n.ToggleCheckIterator-n.Cascade)%int(n.ServerPingPeriod) == 0 && n.Battery > 1 {
			n.Battery = n.Battery - n.BatteryLossCheckingServerScalar
			n.TotalChecksServer = n.TotalChecksServer + 1
			n.HasCheckedServer = true
			n.Server()
		} else {
			n.HasCheckedServer = false
		}
	}
}

/* updateHistory shifts all values in the sample history slice to the right and adds the value at the beginning
Therefore, each time a node takes a sample in main, it also adds this sample to the beginning of the sample history.
Each sample is only stored until ln more samples have been taken (this variable is in hello.go)
*/
func (n *NodeImpl) UpdateHistory(newValue float32) {

	//loop through the sample history slice in reverse order, excluding 0th index
	for i := len(n.SampleHistory) - 1; i > 0; i-- {
		n.SampleHistory[i] = n.SampleHistory[i-1] //set the current index equal to the value of the previous index
	}

	n.SampleHistory[0] = newValue //set 0th index to new measured value

	/* Now calculate the weighted average of the sample history. Note that if a node is stationary, all values
	averaged over are weighted equally. The faster the node is moving, the less the older values are worth when
	calculating the average, because in that case we want the average to more closely reflect the newer values
	*/
	var sum float32
	var numSamples int //variable for number of samples to average over

	var decreaseRatio = n.SpeedWeight / 100.0

	if n.TotalSamples > len(n.SampleHistory) { //if the node has taken more than x total samples
		numSamples = len(n.SampleHistory) //we only average over the x most recent ones
	} else { //if it doesn't have x samples taken yet
		numSamples = n.TotalSamples //we only average over the number of samples it's taken
	}

	for i := 0; i < numSamples; i++ {
		if n.SampleHistory[i] != 0 {
			//weight the values of the sampleHistory when added to the sum variable based on the speed, so older values are weighted less
			sum += n.SampleHistory[i] - ((decreaseRatio) * float32(i))
		} else {
			sum += 0
		}
	}
	n.Avg = sum / float32(numSamples)
}

/* this function increments a node's total number of samples by 1
it's called whenever the node takes a new sample */
func (n *NodeImpl) IncrementTotalSamples() {
	n.TotalSamples++
}

//getter function for average
func (n *NodeImpl) GetAvg() float32 {
	return n.Avg
}

//increases numResets field
func (n *NodeImpl) IncrementNumResets() {
	n.NumResets++
}

//setter function for concentration field
func (n *NodeImpl) SetConcentration(conc float64) {
	n.Concentration = conc
}

//getter function for ID field
func (n *NodeImpl) GetID() int {
	return n.Id
}

//getter function for x and y locations
func (n *NodeImpl) GetLoc() (int, int) {
	return n.X, n.Y
}

//setter function for S0
func (n *NodeImpl) SetS0(s0 float64) {
	n.S0 = s0
}

//setter function for S1
func (n *NodeImpl) SetS1(s1 float64) {
	n.S1 = s1
}

//setter function for S2
func (n *NodeImpl) SetS2(s2 float64) {
	n.S2 = s2
}

//setter function for E0
func (n *NodeImpl) SetE0(e0 float64) {
	n.E0 = e0
}

//setter function for E1
func (n *NodeImpl) SetE1(e1 float64) {
	n.E1 = e1
}

//setter function for E2
func (n *NodeImpl) SetE2(e2 float64) {
	n.E2 = e2
}

//setter function for ET1
func (n *NodeImpl) SetET1(et1 float64) {
	n.ET1 = et1
}

//setter function for ET2
func (n *NodeImpl) SetET2(et2 float64) {
	n.ET2 = et2
}

//getter function for all parameters
func (n *NodeImpl) GetParams() (float64, float64, float64, float64, float64, float64, float64, float64) {
	return n.S0, n.S1, n.S2, n.E0, n.E1, n.E2, n.ET1, n.ET2
}

//getter function for just S0 - S2 parameters
func (n *NodeImpl) GetCoefficients() (float64, float64, float64) {
	return n.S0, n.S1, n.S2
}

//getter function for x
func (n *NodeImpl) GetX() int {
	return n.X
}

//getter function for y
func (n *NodeImpl) GetY() int {
	return n.Y
}

//This is the actual pinging of the sensor
func (n *NodeImpl) Sense(p *Params) {
	if n.HasCheckedGPS == false {
		n.XPos[n.BufferI] = -1
		n.YPos[n.BufferI] = -1
		n.Value[n.BufferI] = n.GetValue(p)
		n.Time[n.BufferI] = p.Iterations_used
		n.BufferI = n.BufferI + 1
	} else {
		n.Value[n.BufferI] = n.GetValue(p)
	}
}

//This is the actual pinging of the GPS
func (n *NodeImpl) GPS(p *Params) {
	if n.HasCheckedSensor == false {
		n.Value[n.BufferI] = -1
		n.XPos[n.BufferI] = n.X
		n.YPos[n.BufferI] = n.Y
		n.Time[n.BufferI] = p.Iterations_used
		n.BufferI = n.BufferI + 1
	} else {
		if n.BufferI > 0 {
			n.XPos[n.BufferI-1] = n.X
			n.YPos[n.BufferI-1] = n.Y
		}
	}
}

//This is the actual upload to the server
func (n *NodeImpl) Server() {
	//getData(&s,n.XPos[0:n.BufferI],n.YPos[0:n.BufferI],n.Value[0:n.BufferI],n.Time[0:n.BufferI], n.Id,n.BufferI)
	n.BufferI = 0
}

//Returns node distance to the bomb
func (n *NodeImpl) GeoDist(b Bomb) float32 {
	//this needs to be changed
	return float32(math.Pow(float64(math.Abs(float64(n.X)-float64(b.X))), 2) + math.Pow(float64(math.Abs(float64(n.Y)-float64(b.Y))), 2))
}

//Returns array of accelerometer speeds recorded for a specific node
func (n *NodeImpl) GetSpeed() []float32 {
	return n.AccelerometerSpeed
}

//Returns a different version of the distance to the bomb
func (n *NodeImpl) GetValue(p *Params) int {
	return int(math.Sqrt(math.Pow(float64(n.X-p.B.X), 2) + math.Pow(float64(n.Y-p.B.Y), 2)))
}

//Takes cares of taking a node's readings and printing detections and stuff
func (curNode *NodeImpl) GetReadings(p *Params) {

	//driftFile.Sync()
	//nodeFile.Sync()
	//positionFile.Sync()
	//test change
	newX, newY := curNode.GetLoc()

	newDist := curNode.Distance(*p.B) //this is the node's reported value without error

	//need to get the correct time reading value from system
	//need to verify where we read from

	//Calculate error, sensitivity, and noise, as per the matlab code
	S0, S1, S2, E0, E1, E2, ET1, ET2 := curNode.GetParams()
	sError := (S0 + E0) + (S1+E1)*math.Exp(-float64(curNode.NodeTime)/(p.Tau1+ET1)) + (S2+E2)*math.Exp(-float64(curNode.NodeTime)/(p.Tau2+ET2))
	curNode.Sensitivity = S0 + (S1)*math.Exp(-float64(curNode.NodeTime)/p.Tau1) + (S2)*math.Exp(-float64(curNode.NodeTime)/p.Tau2)
	sNoise := rand.NormFloat64()*0.5*p.ErrorModifierCM + float64(newDist)*sError

	errorDist := sNoise / curNode.Sensitivity //this is the node's actual reading with error

	//increment node time
	curNode.NodeTime++

	if curNode.HasCheckedSensor {
		curNode.IncrementTotalSamples()
		curNode.UpdateHistory(float32(errorDist))
	}

	//Detection of false positives or false negatives
	if errorDist < p.DetectionThresholdCM && float64(newDist) >= p.DetectionThresholdCM {
		//this is a node false negative due to drifitng
		if curNode.HasCheckedSensor {
			//just drifting
			fmt.Fprintln(p.DriftFile, "Node False Negative (drifting) ID:", curNode.Id, "True Reading:", newDist, "Drifted Reading:",
				errorDist, "S0:", curNode.S0, "S1:", curNode.S1, "S2:", curNode.S2, "E0:", curNode.E0, "E1:", curNode.E1,
				"E2:", curNode.E2, "ET1:", curNode.ET1, "ET2:", curNode.ET2, "time since calibration:", curNode.NodeTime,
				"Current Time (Iter):", p.Iterations_used, "Energy Level:", curNode.Battery, "Distance from bomb:", math.Sqrt(float64(curNode.GeoDist(*p.B))),
				"x:", curNode.X, "y:", curNode.Y)
		} else {
			//both drifting and energy
			fmt.Fprintln(p.DriftFile, "Node False Negative (both) ID:", curNode.Id, "True Reading:", newDist, "Drifted Reading:",
				errorDist, "S0:", curNode.S0, "S1:", curNode.S1, "S2:", curNode.S2, "E0:", curNode.E0, "E1:", curNode.E1,
				"E2:", curNode.E2, "ET1:", curNode.ET1, "ET2:", curNode.ET2, "time since calibration:", curNode.NodeTime,
				"Current Time (Iter):", p.Iterations_used, "Energy Level:", curNode.Battery, "Distance from bomb:", math.Sqrt(float64(curNode.GeoDist(*p.B))),
				"x:", curNode.X, "y:", curNode.Y)
		}
	}

	if errorDist >= p.DetectionThresholdCM && float64(newDist) >= p.DetectionThresholdCM && !curNode.HasCheckedSensor {
		//false negative due solely to energy
		fmt.Fprintln(p.DriftFile, "Node False Negative (energy) ID:", curNode.Id, "True Reading:", newDist, "Drifted Reading:",
			errorDist, "S0:", curNode.S0, "S1:", curNode.S1, "S2:", curNode.S2, "E0:", curNode.E0, "E1:", curNode.E1,
			"E2:", curNode.E2, "ET1:", curNode.ET1, "ET2:", curNode.ET2, "time since calibration:", curNode.NodeTime,
			"Current Time (Iter):", p.Iterations_used, "Energy Level:", curNode.Battery, "Distance from bomb:", math.Sqrt(float64(curNode.GeoDist(*p.B))),
			"x:", curNode.X, "y:", curNode.Y)
	}

	if errorDist >= p.DetectionThresholdCM && float64(newDist) < p.DetectionThresholdCM {
		//this if a false positive
		//it must be due to drifting. Report relevant info to driftFile
		fmt.Fprintln(p.DriftFile, "Node False Positive (drifting) ID:", curNode.Id, "True Reading:", newDist, "Drifted Reading:",
			errorDist, "S0:", curNode.S0, "S1:", curNode.S1, "S2:", curNode.S2, "E0:", curNode.E0, "E1:", curNode.E1,
			"E2:", curNode.E2, "ET1:", curNode.ET1, "ET2:", curNode.ET2, "time since calibration:", curNode.NodeTime,
			"Current Time (Iter):", p.Iterations_used, "Energy Level:", curNode.Battery, "Distance from bomb:", math.Sqrt(float64(curNode.GeoDist(*p.B))),
			"x:", curNode.X, "y:", curNode.Y)
	}

	if errorDist >= p.DetectionThresholdCM && float64(newDist) >= p.DetectionThresholdCM && curNode.HasCheckedSensor {
		fmt.Fprintln(p.DriftFile, "Node True Positive ID:", curNode.Id, "True Reading:", newDist, "Drifted Reading:",
			errorDist, "S0:", curNode.S0, "S1:", curNode.S1, "S2:", curNode.S2, "E0:", curNode.E0, "E1:", curNode.E1,
			"E2:", curNode.E2, "ET1:", curNode.ET1, "ET2:", curNode.ET2, "time since calibration:", curNode.NodeTime,
			"Current Time (Iter):", p.Iterations_used, "Energy Level:", curNode.Battery, "Distance from bomb:", math.Sqrt(float64(curNode.GeoDist(*p.B))),
			"x:", curNode.X, "y:", curNode.Y)
	}

	//If the reading is more than 2 standard deviations away from the grid average, then recalibrate
	//gridAverage := p.Grid[curNode.Row(p.YDiv)][curNode.Col(p.XDiv)].Avg
	//standDev := grid[curNode.Row(yDiv)][curNode.Col(xDiv)].StdDev

	//New condition added: also recalibrate when the node's sensitivity is <= 1/10 of its original sensitvity
	//New condition added: Check to make sure the sensor was pinged this iteration
	if ((curNode.Sensitivity <= (curNode.InitialSensitivity / 2)) && (curNode.HasCheckedSensor)) && p.Iterations_used != 0 {
		curNode.Recalibrate()
		p.Recalibrate = true
		curNode.IncrementNumResets()
	}

	//printing statements to log files, only if the sensor was pinged this iteration
	//if curNode.HasCheckedSensor && nodesPrint{
	if p.NodesPrint {
		if p.Recalibrate {
			fmt.Fprintln(p.NodeFile, "ID:", curNode.GetID(), "Average:", curNode.GetAvg(), "Reading:", newDist, "Error Reading:", errorDist, "Recalibrated")
		} else {
			fmt.Fprintln(p.NodeFile, "ID:", curNode.GetID(), "Average:", curNode.GetAvg(), "Reading:", newDist, "Error Reading:", errorDist)
		}
		//fmt.Fprintln(nodeFile, "battery:", int(curNode.Battery),)
	}

	if p.PositionPrint {
		fmt.Fprintln(p.PositionFile, "ID:", curNode.GetID(), "x:", newX, "y:", newY)
	}

	p.Recalibrate = false

	//Receives the node's distance and calculates its running average
	//for that square
	//Only do this if the sensor was pinged this iteration
	if curNode.HasCheckedSensor {
		//p.Server.UpdateSquareAvg(*curNode, errorDist)
		//p.Grid[curNode.Row(p.YDiv)][curNode.Col(p.XDiv)].TakeMeasurement(float32(errorDist))
		//p.Grid[curNode.Row(p.YDiv)][curNode.Col(p.XDiv)].NumNodes++
		////subtract grid average from node average, square it, and add it to this variable
		p.Server.Send(Reading{errorDist, newX, newY, p.Iterations_used, curNode.GetID()})
		//p.Grid[curNode.Row(p.YDiv)][curNode.Col(p.XDiv)].SquareValues += math.Pow(float64(errorDist-float64(gridAverage)), 2)
	}

}

//Takes cares of taking a node's readings and printing detections and stuff
func (curNode *NodeImpl) GetReadingsCSV(p *Params) {

	//driftFile.Sync()
	//nodeFile.Sync()
	//positionFile.Sync()
	//test change
	newX, newY := curNode.GetLoc()

	//newDist := curNode.Distance(*p.B) //this is the node's reported value without error
	newDist := p.SensorReadings[newX][newY][p.CurrTime]
	//Calculate error, sensitivity, and noise, as per the matlab code
	S0, S1, S2, E0, E1, E2, ET1, ET2 := curNode.GetParams()
	sError := (S0 + E0) + (S1+E1)*math.Exp(-float64(curNode.NodeTime)/(p.Tau1+ET1)) + (S2+E2)*math.Exp(-float64(curNode.NodeTime)/(p.Tau2+ET2))
	curNode.Sensitivity = S0 + (S1)*math.Exp(-float64(curNode.NodeTime)/p.Tau1) + (S2)*math.Exp(-float64(curNode.NodeTime)/p.Tau2)
	sNoise := rand.NormFloat64()*0.5*p.ErrorModifierCM + float64(newDist)*sError

	errorDist := sNoise / curNode.Sensitivity //this is the node's actual reading with error

	//increment node time
	curNode.NodeTime++

	if curNode.HasCheckedSensor {
		curNode.IncrementTotalSamples()
		curNode.UpdateHistory(float32(errorDist))
	}

	//Detection of false positives or false negatives
	if errorDist < p.DetectionThresholdCM && float64(newDist) >= p.DetectionThresholdCM {
		//this is a node false negative due to drifitng
		if curNode.HasCheckedSensor {
			//just drifting
			fmt.Fprintln(p.DriftFile, "Node False Negative (drifting) ID:", curNode.Id, "True Reading:", newDist, "Drifted Reading:",
				errorDist, "S0:", curNode.S0, "S1:", curNode.S1, "S2:", curNode.S2, "E0:", curNode.E0, "E1:", curNode.E1,
				"E2:", curNode.E2, "ET1:", curNode.ET1, "ET2:", curNode.ET2, "time since calibration:", curNode.NodeTime,
				"Current Time (Iter):", p.Iterations_used, "Energy Level:", curNode.Battery, "Distance from bomb:", math.Sqrt(float64(curNode.GeoDist(*p.B))),
				"x:", curNode.X, "y:", curNode.Y)
		} else {
			//both drifting and energy
			fmt.Fprintln(p.DriftFile, "Node False Negative (both) ID:", curNode.Id, "True Reading:", newDist, "Drifted Reading:",
				errorDist, "S0:", curNode.S0, "S1:", curNode.S1, "S2:", curNode.S2, "E0:", curNode.E0, "E1:", curNode.E1,
				"E2:", curNode.E2, "ET1:", curNode.ET1, "ET2:", curNode.ET2, "time since calibration:", curNode.NodeTime,
				"Current Time (Iter):", p.Iterations_used, "Energy Level:", curNode.Battery, "Distance from bomb:", math.Sqrt(float64(curNode.GeoDist(*p.B))),
				"x:", curNode.X, "y:", curNode.Y)
		}
	}

	if errorDist >= p.DetectionThresholdCM && float64(newDist) >= p.DetectionThresholdCM && !curNode.HasCheckedSensor {
		//false negative due solely to energy
		fmt.Fprintln(p.DriftFile, "Node False Negative (energy) ID:", curNode.Id, "True Reading:", newDist, "Drifted Reading:",
			errorDist, "S0:", curNode.S0, "S1:", curNode.S1, "S2:", curNode.S2, "E0:", curNode.E0, "E1:", curNode.E1,
			"E2:", curNode.E2, "ET1:", curNode.ET1, "ET2:", curNode.ET2, "time since calibration:", curNode.NodeTime,
			"Current Time (Iter):", p.Iterations_used, "Energy Level:", curNode.Battery, "Distance from bomb:", math.Sqrt(float64(curNode.GeoDist(*p.B))),
			"x:", curNode.X, "y:", curNode.Y)
	}

	if errorDist >= p.DetectionThresholdCM && float64(newDist) < p.DetectionThresholdCM {
		//this if a false positive
		//it must be due to drifting. Report relevant info to driftFile
		fmt.Fprintln(p.DriftFile, "Node False Positive (drifting) ID:", curNode.Id, "True Reading:", newDist, "Drifted Reading:",
			errorDist, "S0:", curNode.S0, "S1:", curNode.S1, "S2:", curNode.S2, "E0:", curNode.E0, "E1:", curNode.E1,
			"E2:", curNode.E2, "ET1:", curNode.ET1, "ET2:", curNode.ET2, "time since calibration:", curNode.NodeTime,
			"Current Time (Iter):", p.Iterations_used, "Energy Level:", curNode.Battery, "Distance from bomb:", math.Sqrt(float64(curNode.GeoDist(*p.B))),
			"x:", curNode.X, "y:", curNode.Y)
	}

	if errorDist >= p.DetectionThresholdCM && float64(newDist) >= p.DetectionThresholdCM && curNode.HasCheckedSensor {
		fmt.Fprintln(p.DriftFile, "Node True Positive ID:", curNode.Id, "True Reading:", newDist, "Drifted Reading:",
			errorDist, "S0:", curNode.S0, "S1:", curNode.S1, "S2:", curNode.S2, "E0:", curNode.E0, "E1:", curNode.E1,
			"E2:", curNode.E2, "ET1:", curNode.ET1, "ET2:", curNode.ET2, "time since calibration:", curNode.NodeTime,
			"Current Time (Iter):", p.Iterations_used, "Energy Level:", curNode.Battery, "Distance from bomb:", math.Sqrt(float64(curNode.GeoDist(*p.B))),
			"x:", curNode.X, "y:", curNode.Y)
	}

	//If the reading is more than 2 standard deviations away from the grid average, then recalibrate
	//gridAverage := p.Grid[curNode.Row(p.YDiv)][curNode.Col(p.XDiv)].Avg
	//standDev := grid[curNode.Row(yDiv)][curNode.Col(xDiv)].StdDev

	//New condition added: also recalibrate when the node's sensitivity is <= 1/10 of its original sensitvity
	//New condition added: Check to make sure the sensor was pinged this iteration
	if ((curNode.Sensitivity <= (curNode.InitialSensitivity / 2)) && (curNode.HasCheckedSensor)) && p.Iterations_used != 0 {
		curNode.Recalibrate()
		p.Recalibrate = true
		curNode.IncrementNumResets()
	}

	//printing statements to log files, only if the sensor was pinged this iteration
	//if curNode.HasCheckedSensor && nodesPrint{
	if p.NodesPrint {
		if p.Recalibrate {
			fmt.Fprintln(p.NodeFile, "ID:", curNode.GetID(), "Average:", curNode.GetAvg(), "Reading:", newDist, "Error Reading:", errorDist, "Recalibrated")
		} else {
			fmt.Fprintln(p.NodeFile, "ID:", curNode.GetID(), "Average:", curNode.GetAvg(), "Reading:", newDist, "Error Reading:", errorDist)
		}
		//fmt.Fprintln(nodeFile, "battery:", int(curNode.Battery),)
	}

	if p.PositionPrint {
		fmt.Fprintln(p.PositionFile, "ID:", curNode.GetID(), "x:", newX, "y:", newY)
	}

	p.Recalibrate = false

	//Receives the node's distance and calculates its running average
	//for that square
	//Only do this if the sensor was pinged this iteration
	if curNode.HasCheckedSensor {
		//p.Grid[curNode.Row(p.YDiv)][curNode.Col(p.XDiv)].TakeMeasurement(float32(errorDist))
		//p.Grid[curNode.Row(p.YDiv)][curNode.Col(p.XDiv)].NumNodes++
		//subtract grid average from node average, square it, and add it to this variable
		//p.Grid[curNode.Row(p.YDiv)][curNode.Col(p.XDiv)].SquareValues += (math.Pow(float64(errorDist-float64(gridAverage)), 2))
		p.Server.Send(Reading{errorDist, newX, newY, p.Iterations_used, curNode.GetID()})
	}
}

func rangeInt(min, max int) int { //returns a random number between max and min
	return rand.Intn(max-min) + min
}
