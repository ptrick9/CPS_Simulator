package main

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
	distance(b bomb) float32        //distance to bomb in the form of the node's reading
	row(div int) int                //Row of node
	col(div int) int                //Column of node
	getSpeed() []float32            //History of accelerometer based speeds of node
	batteryLossDynamic()            //Battery loss based of ratios of battery usage
	batteryLossDynamic1()           //2 stage buffer battery loss
	updateHistory(newValue float32) //updates history of node's samples
	incrementTotalSamples()         //increments total number of samples node has taken
	getAvg() float32                //returns average of node's past samples
	incrementNumResets()            //increments the number of times a node has been reset
	setConcentration(conc float64)  //sets the concentration of a node
	geoDist(b bomb) float32         //returns distance from bomb (rather than reading of node)
	getID() int                     //returns ID of node
	getLoc() (x, y int)             //returns x and y values of node
	//following functions set drifting parameters of nodes
	setS0(s0 float64)
	setS1(s1 float64)
	setS2(s2 float64)
	setE0(e0 float64)
	setE1(e1 float64)
	setE2(e2 float64)
	setET1(et1 float64)
	setET2(et2 float64)
	getParams() (float64, float64, float64, float64, float64, float64, float64, float64) //returns all of the above parameters
	getCoefficients() (float64, float64, float64)                                        //returns some of the above parameters
	getX() int                                                                           //returns x position of node
	getY() int                                                                           //returns y position of node
}

//NodeImpl is a struct that implements all the methods listed
//	above in NodeParent
type NodeImpl struct {
	id                              int      //id of node
	oldX                            int      // for movement
	oldY                            int      // for movement
	sitting                         int      // for movement
	x                               int      //x pos of node
	y                               int      //y pos of node
	battery                         float32  //battery of node
	batteryLossScalar               float32  //natural incremental battery loss of node
	batteryLossCheckingSensorScalar float32  //sensor based battery loss of node
	batteryLossGPSScalar            float32  //GPS based battery loss of node
	batteryLossCheckingServerScalar float32  //server communication based battery loss of node
	toggleCheckIterator             int      //node's personal iterator mostly for cascading pings
	hasCheckedSensor                bool     //did the node just ping the sensor?
	totalChecksSensor               int      //total sensor pings of node
	hasCheckedGPS                   bool     //did the node just ping the GPS?
	totalChecksGPS                  int      //total GPS pings of node
	hasCheckedServer                bool     //did the node just communicate with the server?
	totalChecksServer               int      //how many times did the node communicate with the server?
	pingPeriod                      float32  //This is the aggregate ping period used in some ping rate determining algorithms
	sensorPingPeriod                float32  //This is the ping period for the sensor
	GPSPingPeriod                   float32  //This is the ping period for the GPS
	serverPingPeriod                float32  //This is the ping period for the server
	pings                           float32  //This is an aggregate pings used in some ping rate determining algorithms
	sensorPings                     float32  //This is the total sensor pings to be made
	GPSPings                        float32  //This is the total GPS pings to be made
	serverPings                     float32  //This is the total server pings to be made
	cascade                         int      //This cascades the pings of the nodes
	bufferI                         int      //This is to keep track of the node's buffer size
	xPos                            [100]int //x pos buffer of node
	yPos                            [100]int //y pos buffer of node
	value                           [100]int //value buffer of node
	accelerometerSpeedServer        [100]int //Accelerometer speed history of node
	time                            [100]int //This keeps track of when specific pings are made
	//speedGPSPeriod int //This is a special period for speed based GPS pings but it is not used and may never be
	accelerometerPosition [2][3]int //This is the accelerometer model of node
	accelerometerSpeed    []float32 //History of accelerometer speeds recorded
	inverseSensor         float32   //Algorithm place holder declared here for speed
	inverseGPS            float32   //Algorithm place holder declared here for speed
	inverseServer         float32   //Algorithm place holder declared here for speed
	sampleHistory         []float32 //a history of the node's readings
	avg                   float32   //weighted average of the node's most recent readings
	totalSamples          int       //total number of samples taken by a node
	speedWeight           float32   //weight given to averaging of node's samples, based on node's speed
	numResets             int       //number of times a node has had to reset due to drifting
	concentration         float64   //used to determine reading of node
	speedGPSPeriod        int

	current  int
	previous int
	diffx    int
	diffy    int
	speed    float32

	//The following values are all various drifting parameters of the node
	newX               int
	newY               int
	S0                 float64
	S1                 float64
	S2                 float64
	E0                 float64
	E1                 float64
	E2                 float64
	ET1                float64
	ET2                float64
	nodeTime           int
	sensitivity        float64
	initialSensitivity float64
}

//NodeMovement controls the movement of all the normal nodes
//It inherits all the methods and attributes from NodeParent
//	and NodeImpl
type NodeMovement interface {
	NodeParent
	move()
}

//Bouncing nodes bound around the grid
type bn struct {
	*NodeImpl
	x_speed int
	y_speed int
}

//Wall nodes go in a straight line from top/bottom or
//	side/side
type wn struct {
	*NodeImpl
	speed int
	dir   int
}

//Random nodes get assigned a random x, y velocity every
//	move update
type rn struct {
	*NodeImpl
}

type wallNodes struct {
	node *NodeImpl
}

//Coord is a struct that contains x and y coordinates of
//	a square in the grid
//This struct is used by the super node type to create its
//	route through the grid
type Coord struct {
	parent      *Coord
	x, y        int
	time        int
	g, h, score int
}

//Path is a struct that contains an x and y integer and
//	a float for distance
//This struct is used when calculating the distance between
//	points of interest on the grid during super node route
//	scheduling
type Path struct {
	x, y int
	dist float64
}

//Returns the x index of the square in which the specific
//	node currently resides
func (n *NodeImpl) row(div int) int {
	return n.y / div
}

//Returns the y index of the square in which the specific
//	node currently resides
func (n *NodeImpl) col(div int) int {
	return n.x / div
}

//Returns a float representing the detection of the bomb
//	by the specific node depending on distance
func (n *NodeImpl) distance(b bomb) float32 {
	dist := float32(math.Pow(float64(math.Abs(float64(n.x)-float64(b.x))), 2) + math.Pow(float64(math.Abs(float64(n.y)-float64(b.y))), 2))

	if dist == 0 {
		return 1000
	} else {
		//return float32(1000 / (math.Pow((float64(dist)/0.2)*0.25,1.5)))
		return float32(math.Pow(1000/float64(dist), 1.5))
	}
}

// These are the toString methods for battery levels
func (n bn) String() string { // extra extra string statements
	return fmt.Sprintf("x: %v y: %v Xspeed: %v Yspeed: %v id: %v battery: %v sensor checked: %v sensor checks: %v GPS checked: %v GPS checks: %v server checked: %v server checks: %v buffer: %v ", n.x, n.y, n.x_speed, n.y_speed, n.id, n.battery, n.hasCheckedSensor, n.totalChecksSensor, n.hasCheckedGPS, n.totalChecksGPS, n.hasCheckedServer, n.totalChecksServer, n.bufferI)
}

func (n wn) String() string {
	return fmt.Sprintf("x: %v y: %v speed: %v dir: %v id: %v battery: %v sensor checked: %v sensor checks: %v GPS checked: %v GPS checks: %v server checked: %v server checks: %v buffer: %v ", n.x, n.y, n.speed, n.dir, n.id, n.battery, n.hasCheckedSensor, n.totalChecksSensor, n.hasCheckedGPS, n.totalChecksGPS, n.hasCheckedServer, n.totalChecksServer, n.bufferI)
}

func (n rn) String() string {
	return fmt.Sprintf("x: %v y: %v id: %v battery: %v sensor checked: %v sensor checks: %v GPS checked: %v GPS checks: %v server checked: %v server checks: %v buffer: %v ", n.x, n.y, n.id, n.battery, n.hasCheckedSensor, n.totalChecksSensor, n.hasCheckedGPS, n.totalChecksGPS, n.hasCheckedServer, n.totalChecksServer, n.bufferI)
} // end extra extra string statements

func (n NodeImpl) String() string {
	//return fmt.Sprintf("x: %v y: %v id: %v battery: %v sensor checked: %v sensor checks: %v GPS checked: %v GPS checks: %v server checked: %v server checks: %v buffer: %v ", n.x, n.y, n.id, n.battery, n.hasCheckedSensor, n.totalChecksSensor, n.hasCheckedGPS, n.totalChecksGPS, n.hasCheckedServer, n.totalChecksServer,n.bufferI)
	return fmt.Sprintf("battery: %v sensor checked: %v GPS checked: %v ", int(n.battery), n.hasCheckedSensor, n.hasCheckedGPS)

}

func (c Coord) String() string {
	return fmt.Sprintf("{%v %v %v}", c.x, c.y, c.time)
}

func (c Coord) equals(c2 Coord) bool {
	return c.x == c2.x && c.y == c2.y
}

func (n *NodeImpl) move() {
	if n.sitting <= sittingStopThresholdCM {
		n.oldX = n.x / xDiv
		n.oldY = n.y / yDiv

		var potentialSpots []gridSpot

		//only add the ones that are valid to move to into the list
		if n.y-1 >= 0 &&
			n.x >= 0 &&
			n.x < len(boardMap[n.y-1]) &&
			n.y-1 < len(boardMap) &&

			boardMap[n.y-1][n.x] != -1 &&
			boolGrid[n.y-1][n.x] == false { // &&
			//boardMap[n.x][n.y-1] <= boardMap[n.x][n.y] {

			up := gridSpot{n.x, n.y - 1, boardMap[n.y-1][n.x]}
			potentialSpots = append(potentialSpots, up)
		}
		if n.x+1 < len(boardMap[n.y]) &&
			n.x+1 >= 0 &&
			n.y < len(boardMap) &&
			n.y >= 0 &&

			boardMap[n.y][n.x+1] != -1 &&
			boolGrid[n.y][n.x+1] == false { // &&
			//boardMap[n.x+1][n.y] <= boardMap[n.x][n.y] {

			right := gridSpot{n.x + 1, n.y, boardMap[n.y][n.x+1]}
			potentialSpots = append(potentialSpots, right)
		}
		if n.y+1 < len(boardMap) &&
			n.y+1 >= 0 &&
			n.x < len(boardMap[n.y+1]) &&
			n.x >= 0 &&

			boardMap[n.y+1][n.x] != -1 &&
			boolGrid[n.y+1][n.x] == false { //&&
			//boardMap[n.x][n.y+1] <= boardMap[n.x][n.y] {

			down := gridSpot{n.x, n.y + 1, boardMap[n.y+1][n.x]}
			potentialSpots = append(potentialSpots, down)
		}
		if n.x-1 >= 0 &&
			n.x-1 < len(boardMap[n.y]) &&
			n.y >= 0 &&
			n.y < len(boardMap) &&

			boardMap[n.y][n.x-1] != -1 &&
			boolGrid[n.y][n.x-1] == false { // &&
			//boardMap[n.x-1][n.y] <= boardMap[n.x][n.y] {

			left := gridSpot{n.x - 1, n.y, boardMap[n.y][n.x-1]}
			potentialSpots = append(potentialSpots, left)
		}

		sort.Sort(byRandom(potentialSpots))
		sort.Sort(byRandom(potentialSpots))
		sort.Sort(byRandom(potentialSpots))
		sort.Sort(byValue(potentialSpots))

		for i := 0; i < len(potentialSpots); i++ {
			if grid[potentialSpots[i].y/yDiv][potentialSpots[i].x/xDiv].actualNumNodes <= squareCapacity {
				n.x = potentialSpots[i].x
				n.y = potentialSpots[i].y
				break
			}
		}

		if n.x/xDiv != n.oldX || n.y/yDiv != n.oldY {
			grid[n.y/yDiv][n.x/xDiv].actualNumNodes = grid[n.y/yDiv][n.x/xDiv].actualNumNodes + 1
			grid[n.oldY][n.oldX].actualNumNodes = grid[n.oldY][n.oldX].actualNumNodes - 1
		}
		if n.diffx == 0 && n.diffy == 0 || n.sitting < 0 {
			n.sitting = n.sitting + 1
		} else {
			n.sitting = 0
		}
	}
}

func (n *NodeImpl) recalibrate() {
	n.sensitivity = n.initialSensitivity
	n.nodeTime = 0
}

//Moves the bouncing node
func (n *bn) move() {
	//Boundary conditions
	if n.x+n.x_speed < maxX && n.x+n.x_speed >= 0 {
		n.x = n.x + n.x_speed
	} else {
		if n.x+n.x_speed >= maxX {
			n.x = n.x - (n.x_speed - (maxX - 1 - n.x))
			n.x_speed = n.x_speed * -1
		} else {
			n.x = (n.x_speed + n.x) * -1
			n.x_speed = n.x_speed * -1
		}
	}
	if n.y+n.y_speed < maxY && n.y+n.y_speed >= 0 {
		n.y = n.y + n.y_speed
	} else {
		if n.y+n.y_speed >= maxY {
			n.y = n.y - (n.y_speed - (maxY - 1 - n.y))
			n.y_speed = n.y_speed * -1
		} else {
			n.y = (n.y_speed + n.y) * -1
			n.y_speed = n.y_speed * -1
		}
	}
}

//Moves the wall nodes
func (n *wn) move() {
	if n.dir == 0 { //x-axis
		//Boundary conditions
		if n.x+n.speed < maxX && n.x+n.speed >= 0 {
			n.x = n.x + n.speed
		} else {
			if n.x+n.speed >= maxX {
				n.x = n.x - (n.speed - (maxX - 1 - n.x))
				n.speed = n.speed * -1
			} else {
				n.x = (n.speed + n.x) * -1
				n.speed = n.speed * -1
			}
		}
	} else { //y-axis
		if n.y+n.speed < maxY && n.y+n.speed >= 0 {
			n.y = n.y + n.speed
		} else {
			if n.y+n.speed >= maxY {
				n.y = n.y - (n.speed - (maxY - 1 - n.y))
				n.speed = n.speed * -1
			} else {
				n.y = (n.speed + n.y) * -1
				n.speed = n.speed * -1
			}
		}
	}
}

//Moves the random nodes
func (n *rn) move() {
	x_speed := rangeInt(-3, 3)
	y_speed := rangeInt(-3, 3)

	//Boundary conditions
	if n.x+x_speed < maxX && n.x+x_speed >= 0 {
		n.x = n.x + x_speed
	} else {
		if n.x+x_speed >= maxX {
			n.x = n.x - (x_speed - (maxX - 1 - n.x))
		} else {
			n.x = (x_speed + n.x) * -1
		}
	}
	if n.y+y_speed < maxY && n.y+y_speed >= 0 {
		n.y = n.y + y_speed
	} else {
		if n.y+y_speed >= maxY {
			n.y = n.y - (y_speed - (maxY - 1 - n.y))
		} else {
			n.y = (y_speed + n.y) * -1
		}
	}
}

//Returns the arr with the element at index n removed
func remove_index(arr []Path, n int) []Path {
	return arr[:n+copy(arr[n:], arr[n+1:])]
}

//Returns the array with the range of elements from index
//	a to b removed
func remove_range(arr []Coord, a, b int) []Coord {
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
func insert_array(arr1 []Coord, arr2 []Coord, n int) []Coord {
	if len(arr1) == 0 {
		return arr2
	} else {
		return append(arr1[:n], append(arr2, arr1[n:]...)...)
	}
}

//Returns the array with the element at ind1 moved to ind2
//ind1 must always bef urther in the array than ind2
func remove_and_insert(arr []Coord, ind1, ind2 int) []Coord {
	arr1 := make([]Coord, 0)
	c := arr[ind1]
	arr = arr[:ind1+copy(arr[ind1:], arr[ind1+1:])]
	arr1 = append(arr1, c)
	return append(arr[:ind2], append(arr1, arr[ind2:]...)...)
}

// This is the battery loss function that clears the buffer at 2 different rates based on the battery percentage left
func (n *NodeImpl) batteryLossDynamic1() {
	n.hasCheckedGPS = false
	n.hasCheckedSensor = false
	n.hasCheckedServer = false

	// This is the manual buffer clearing
	if n.bufferI >= maxBufferCapacityCM {
		n.server()
	}
	// This is a generic iterator
	n.toggleCheckIterator = n.toggleCheckIterator + 1
	// These are the respective accelerometer positions
	n.current = n.toggleCheckIterator % 3
	n.previous = (n.toggleCheckIterator - 1) % 3
	if n.current == 0 {
		n.accelerometerPosition[0][0] = n.x
		n.accelerometerPosition[1][0] = n.y
	} else if n.current == 1 {
		n.accelerometerPosition[0][1] = n.x
		n.accelerometerPosition[1][1] = n.y
	} else if n.current == 2 {
		n.accelerometerPosition[0][2] = n.x
		n.accelerometerPosition[1][2] = n.y
	}
	n.diffx = n.accelerometerPosition[0][n.current] - n.accelerometerPosition[0][n.previous]
	n.diffy = n.accelerometerPosition[1][n.current] - n.accelerometerPosition[1][n.previous]
	// This is the accelerometer determined speed
	n.speed = float32(math.Sqrt(float64(n.diffx*n.diffx + n.diffy*n.diffy)))
	// This is a list of the previous accelerometer determined speeds
	n.accelerometerSpeed = append(n.accelerometerSpeed, n.speed)
	//threshHoldBatteryToHave = thres
	// This is the natural loss of the battery
	if n.battery > 0 {
		n.battery = n.battery - n.batteryLossScalar
	}
	// This is the predicted natural loss to prevent overages.
	naturalLoss = n.battery - (float32(iterations_of_event-iterations_used) * n.batteryLossScalar)
	// This is the algorithm that determines sampling rate's ratios
	n.pings = n.battery * totalPercentBatteryToUse / (n.batteryLossCheckingSensorScalar + n.batteryLossGPSScalar + n.batteryLossCheckingServerScalar) // set percentage for consumption here, also '50' if no minus
	n.inverseSensor = 1 / n.batteryLossCheckingSensorScalar
	n.inverseGPS = 1 / n.batteryLossGPSScalar
	n.inverseServer = 1 / n.batteryLossCheckingServerScalar
	n.sensorPings = n.pings * (n.inverseSensor / (n.inverseServer + n.inverseGPS + n.inverseGPS))
	n.GPSPings = n.pings * (n.inverseGPS / (n.inverseServer + n.inverseGPS + n.inverseGPS))
	n.serverPings = n.pings * (n.inverseServer / (n.inverseServer + n.inverseGPS + n.inverseGPS))

	if naturalLoss > threshHoldBatteryToHave {
		n.sensorPingPeriod = float32(iterations_of_event-iterations_used) / n.sensorPings
		if n.sensorPingPeriod < 1 {
			n.sensorPingPeriod = 1
		}
		// Checks to see if sensor is pinged
		if (n.toggleCheckIterator-n.cascade)%int(n.sensorPingPeriod) == 0 && n.battery > 1 {
			n.battery = n.battery - n.batteryLossCheckingSensorScalar
			n.totalChecksSensor = n.totalChecksSensor + 1
			n.hasCheckedSensor = true
			n.sense()
		} else {
			n.hasCheckedSensor = false
		}
		// Checks to see if GPS is pinged
		n.GPSPingPeriod = float32(iterations_of_event-iterations_used) / n.GPSPings
		if n.GPSPingPeriod < 1 {
			n.GPSPingPeriod = 1
		}
		if ((n.toggleCheckIterator-n.cascade)%int(n.GPSPingPeriod) == 0 && n.battery > 1) || (n.speed > float32(movementSamplingSpeedCM) && n.toggleCheckIterator%movementSamplingPeriodCM == 0) { // && n.toggleCheckIterator%n.speedGPSPeriod == 0
			n.battery = n.battery - n.batteryLossGPSScalar
			n.totalChecksGPS = n.totalChecksGPS + 1
			n.hasCheckedGPS = true
			n.GPS()
		} else {
			n.hasCheckedGPS = false
		}

		// This is the 2 stage buffer based on battery percentages
		if n.battery >= 75 { //100 - 75 percent
			if (n.toggleCheckIterator-n.cascade)%14 == 0 && n.battery > 1 { // 1000/70 = 14
				n.battery = n.battery - n.batteryLossCheckingServerScalar
				n.totalChecksServer = n.totalChecksServer + 1
				n.hasCheckedServer = true
				n.server()
			} else {
				n.hasCheckedServer = false
			}
		} else if n.battery >= 30 && n.battery < 75 { //70 - 30 percent
			if (n.toggleCheckIterator-n.cascade)%50 == 0 && n.battery > 1 { //1000/20 = 50
				n.battery = n.battery - n.batteryLossCheckingServerScalar
				n.totalChecksServer = n.totalChecksServer + 1
				n.hasCheckedServer = true
				n.server()
			} else {
				n.hasCheckedServer = false
			}
		} else {
			n.hasCheckedServer = false
		}
	}
}

//This is the battery loss function where the server sensor and GPS are pinged separately and by their own accord
func (n *NodeImpl) batteryLossTable() {
	n.hasCheckedGPS = false
	n.hasCheckedSensor = false
	n.hasCheckedServer = false

	// This is the buffer limit if the limit is meet the buffer is forcibly cleared
	if n.bufferI >= maxBufferCapacityCM {
		n.server()
	}
	// This iterator is generic
	n.toggleCheckIterator = n.toggleCheckIterator + 1
	// These are the nodes respective accelerometer positions
	n.current = n.toggleCheckIterator % 3
	n.previous = (n.toggleCheckIterator - 1) % 3
	// this is the accelerometer's functions
	if n.current == 0 {
		n.accelerometerPosition[0][0] = n.x
		n.accelerometerPosition[1][0] = n.y
	} else if n.current == 1 {
		n.accelerometerPosition[0][1] = n.x
		n.accelerometerPosition[1][1] = n.y
	} else if n.current == 2 {
		n.accelerometerPosition[0][2] = n.x
		n.accelerometerPosition[1][2] = n.y
	}
	n.diffx = n.accelerometerPosition[0][n.current] - n.accelerometerPosition[0][n.previous]
	n.diffy = n.accelerometerPosition[1][n.current] - n.accelerometerPosition[1][n.previous]
	// Speed as determined by accelerometer
	speed := float32(math.Sqrt(float64(n.diffx*n.diffx + n.diffy*n.diffy)))
	// This keeps track of the accelerometer values
	n.accelerometerSpeed = append(n.accelerometerSpeed, speed)
	// natural loss of the battery
	if n.battery > 0 {
		n.battery = n.battery - n.batteryLossScalar
	}
	// predicted natural loss of the battery
	naturalLoss = n.battery - (float32(iterations_of_event-iterations_used) * n.batteryLossScalar)

	// this is the battery loss based on checking the sensor, GPS, and server.
	if naturalLoss > threshHoldBatteryToHave {
		if (n.toggleCheckIterator-n.cascade)%sensorSamplingPeriodCM == 0 {
			n.battery = n.battery - n.batteryLossCheckingSensorScalar
			n.totalChecksSensor = n.totalChecksSensor + 1
			n.hasCheckedSensor = true
			n.sense()
		} else {
			n.hasCheckedSensor = false
		}
		if (n.toggleCheckIterator-n.cascade)%GPSSamplingPeriodCM == 0 || (speed > float32(movementSamplingSpeedCM) && n.toggleCheckIterator%movementSamplingPeriodCM == 0) { // && n.toggleCheckIterator%n.speedGPSPeriod == 0
			n.battery = n.battery - n.batteryLossGPSScalar
			n.totalChecksGPS = n.totalChecksGPS + 1
			n.hasCheckedGPS = true
			n.GPS()
		} else {
			n.hasCheckedGPS = false
		}

		// Check to ping server
		if (n.toggleCheckIterator-n.cascade)%serverSamplingPeriodCM == 0 {
			n.battery = n.battery - n.batteryLossCheckingServerScalar
			n.totalChecksServer = n.totalChecksServer + 1
			n.hasCheckedServer = true
			n.server()
		} else {
			n.hasCheckedServer = false
		}
	}
}

//This is the battery loss function where the server sensor and GPS are pinged separately and by their own accord
func (n *NodeImpl) batteryLossMostDynamic() {
	n.hasCheckedGPS = false
	n.hasCheckedSensor = false
	n.hasCheckedServer = false

	// This is the buffer limit if the limit is meet the buffer is forcibly cleared
	if n.bufferI >= maxBufferCapacityCM {
		n.server()
	}
	// This iterator is generic
	n.toggleCheckIterator = n.toggleCheckIterator + 1
	// These are the nodes respective accelerometer positions
	n.current = n.toggleCheckIterator % 3
	n.previous = (n.toggleCheckIterator - 1) % 3
	// this is the accelerometer's functions
	if n.current == 0 {
		n.accelerometerPosition[0][0] = n.x
		n.accelerometerPosition[1][0] = n.y
	} else if n.current == 1 {
		n.accelerometerPosition[0][1] = n.x
		n.accelerometerPosition[1][1] = n.y
	} else if n.current == 2 {
		n.accelerometerPosition[0][2] = n.x
		n.accelerometerPosition[1][2] = n.y
	}
	n.diffx = n.accelerometerPosition[0][n.current] - n.accelerometerPosition[0][n.previous]
	n.diffy = n.accelerometerPosition[1][n.current] - n.accelerometerPosition[1][n.previous]
	// Speed as determined by accelerometer
	n.speed = float32(math.Sqrt(float64(n.diffx*n.diffx + n.diffy*n.diffy)))
	// This keeps track of the accelerometer values
	n.accelerometerSpeed = append(n.accelerometerSpeed, n.speed)
	// natural loss of the battery
	if n.battery > 0 {
		n.battery = n.battery - n.batteryLossScalar
	}
	// predicted natural loss of the battery
	naturalLoss = n.battery - (float32(iterations_of_event) * n.batteryLossScalar)

	// This is the ratio algorithm used to determine the rate of pings
	n.inverseSensor = 1 / n.batteryLossCheckingSensorScalar
	n.inverseGPS = 1 / n.batteryLossGPSScalar
	n.inverseServer = 1 / n.batteryLossCheckingServerScalar

	//SensorBatteryToUse := (totalPercentBatteryToUse * (n.inverseSensor / (n.inverseServer + n.inverseGPS + n.inverseSensor)))
	//GPSBatteryToUse := (totalPercentBatteryToUse * (n.inverseGPS / (n.inverseServer + n.inverseGPS + n.inverseSensor)))
	//ServerBatteryToUse := (totalPercentBatteryToUse * (n.inverseServer / (n.inverseServer + n.inverseGPS + n.inverseSensor)))

	n.sensorPings = (totalPercentBatteryToUse * (n.inverseSensor / (n.inverseServer + n.inverseGPS + n.inverseSensor))) / n.batteryLossCheckingSensorScalar
	n.GPSPings = (totalPercentBatteryToUse * (n.inverseGPS / (n.inverseServer + n.inverseGPS + n.inverseSensor))) / n.batteryLossGPSScalar
	n.serverPings = (totalPercentBatteryToUse * (n.inverseServer / (n.inverseServer + n.inverseGPS + n.inverseSensor))) / n.batteryLossCheckingServerScalar

	// this is the battery loss based on checking the sensor, GPS, and server.
	if naturalLoss > threshHoldBatteryToHave {
		n.sensorPingPeriod = float32(iterations_of_event) / n.sensorPings //-iterations_used
		if n.sensorPingPeriod < 1 {
			n.sensorPingPeriod = 1
		}
		// Check to ping sensor
		if (n.toggleCheckIterator-n.cascade)%int(n.sensorPingPeriod) == 0 && n.battery > 1 {
			n.battery = n.battery - n.batteryLossCheckingSensorScalar
			n.totalChecksSensor = n.totalChecksSensor + 1
			n.hasCheckedSensor = true
			n.sense()
		} else {
			n.hasCheckedSensor = false
		}
		n.GPSPingPeriod = float32(iterations_of_event) / n.GPSPings //-iterations_used
		if n.GPSPingPeriod < 1 {
			n.GPSPingPeriod = 1
		}
		// Check to ping GPS, also note the extra pings made by a high speed.
		if ((n.toggleCheckIterator-n.cascade)%int(n.GPSPingPeriod) == 0 && n.battery > 1) || (n.speed > float32(movementSamplingSpeedCM) && n.toggleCheckIterator%movementSamplingPeriodCM == 0) { // && n.toggleCheckIterator%n.speedGPSPeriod == 0
			n.battery = n.battery - n.batteryLossGPSScalar
			n.totalChecksGPS = n.totalChecksGPS + 1
			n.hasCheckedGPS = true
			n.GPS()
		} else {
			n.hasCheckedGPS = false
		}
		n.serverPingPeriod = float32(iterations_of_event) / n.serverPings //-iterations_used
		if n.serverPingPeriod < 1 {
			n.serverPingPeriod = 1.1
		} else if int(n.serverPingPeriod) > iterations_of_event {
			n.serverPingPeriod = float32(iterations_of_event)
		}
		if n.toggleCheckIterator-n.cascade == 0 {
			//fmt.Println("what?")
			n.toggleCheckIterator = n.cascade + 1
		}
		// Check to ping server
		//fmt.Println(n.toggleCheckIterator,n.cascade,n.serverPingPeriod,n.id, n.serverPings,n.batteryLossCheckingServerScalar, iterations_of_event,float32(iterations_of_event),int(float32(iterations_of_event)))
		if (n.toggleCheckIterator-n.cascade)%int(n.serverPingPeriod) == 0 && n.battery > 1 {
			n.battery = n.battery - n.batteryLossCheckingServerScalar
			n.totalChecksServer = n.totalChecksServer + 1
			n.hasCheckedServer = true
			n.server()
		} else {
			n.hasCheckedServer = false
		}
	}
}

//This is the battery loss function where the server sensor and GPS are pinged separately and by their own accord
func (n *NodeImpl) batteryLossDynamic() {
	n.hasCheckedGPS = false
	n.hasCheckedSensor = false
	n.hasCheckedServer = false

	// This is the buffer limit if the limit is meet the buffer is forcibly cleared
	if n.bufferI >= maxBufferCapacityCM {
		n.server()
	}
	// This iterator is generic
	n.toggleCheckIterator = n.toggleCheckIterator + 1
	// These are the nodes respective accelerometer positions
	current := n.toggleCheckIterator % 3
	previous := (n.toggleCheckIterator - 1) % 3
	// this is the accelerometer's functions
	if current == 0 {
		n.accelerometerPosition[0][0] = n.x
		n.accelerometerPosition[1][0] = n.y
	} else if current == 1 {
		n.accelerometerPosition[0][1] = n.x
		n.accelerometerPosition[1][1] = n.y
	} else if current == 2 {
		n.accelerometerPosition[0][2] = n.x
		n.accelerometerPosition[1][2] = n.y
	}
	diffx := n.accelerometerPosition[0][current] - n.accelerometerPosition[0][previous]
	diffy := n.accelerometerPosition[1][current] - n.accelerometerPosition[1][previous]
	// Speed as determined by accelerometer
	speed := float32(math.Sqrt(float64(diffx*diffx + diffy*diffy)))
	// This keeps track of the accelerometer values
	n.accelerometerSpeed = append(n.accelerometerSpeed, speed)
	// natural loss of the battery
	if n.battery > 0 {
		n.battery = n.battery - n.batteryLossScalar
	}
	// predicted natural loss of the battery
	naturalLoss = n.battery - (float32(iterations_of_event-iterations_used) * n.batteryLossScalar)

	// This is the ratio algorithm used to determine the rate of pings
	n.pings = n.battery * .5 / (n.batteryLossCheckingSensorScalar + n.batteryLossGPSScalar + n.batteryLossCheckingServerScalar) // set percentage for consumption here, also '50' if no minus
	n.inverseSensor = 1 / n.batteryLossCheckingSensorScalar
	n.inverseGPS = 1 / n.batteryLossGPSScalar
	n.inverseServer = 1 / n.batteryLossCheckingServerScalar
	n.sensorPings = n.pings * (n.inverseSensor / (n.inverseServer + n.inverseGPS + n.inverseGPS))
	n.GPSPings = n.pings * (n.inverseGPS / (n.inverseServer + n.inverseGPS + n.inverseGPS))
	n.serverPings = n.pings * (n.inverseServer / (n.inverseServer + n.inverseGPS + n.inverseGPS))

	// this is the battery loss based on checking the sensor, GPS, and server.
	if naturalLoss > threshHoldBatteryToHave {
		n.sensorPingPeriod = float32(iterations_of_event) / n.sensorPings //-iterations_used
		if n.sensorPingPeriod < 1 {
			n.sensorPingPeriod = 1
		}
		// Check to ping sensor
		if (n.toggleCheckIterator-n.cascade)%int(n.sensorPingPeriod) == 0 && n.battery > 1 {
			n.battery = n.battery - n.batteryLossCheckingSensorScalar
			n.totalChecksSensor = n.totalChecksSensor + 1
			n.hasCheckedSensor = true
			n.sense()
		} else {
			n.hasCheckedSensor = false
		}
		n.GPSPingPeriod = float32(iterations_of_event) / n.GPSPings //-iterations_used
		if n.GPSPingPeriod < 1 {
			n.GPSPingPeriod = 1
		}
		// Check to ping GPS, also note the extra pings made by a high speed.
		if ((n.toggleCheckIterator-n.cascade)%int(n.GPSPingPeriod) == 0 && n.battery > 1) || (speed > float32(movementSamplingSpeedCM) && n.toggleCheckIterator%movementSamplingPeriodCM == 0) { // && n.toggleCheckIterator%n.speedGPSPeriod == 0
			n.battery = n.battery - n.batteryLossGPSScalar
			n.totalChecksGPS = n.totalChecksGPS + 1
			n.hasCheckedGPS = true
			n.GPS()
		} else {
			n.hasCheckedGPS = false
		}
		n.serverPingPeriod = float32(iterations_of_event) / n.serverPings //-iterations_used
		if n.serverPingPeriod < 1 {
			n.serverPingPeriod = 1
		}
		// Check to ping server
		if (n.toggleCheckIterator-n.cascade)%int(n.serverPingPeriod) == 0 && n.battery > 1 {
			n.battery = n.battery - n.batteryLossCheckingServerScalar
			n.totalChecksServer = n.totalChecksServer + 1
			n.hasCheckedServer = true
			n.server()
		} else {
			n.hasCheckedServer = false
		}
	}
}

/* updateHistory shifts all values in the sample history slice to the right and adds the value at the beginning
Therefore, each time a node takes a sample in main, it also adds this sample to the beginning of the sample history.
Each sample is only stored until ln more samples have been taken (this variable is in hello.go)
*/
func (n *NodeImpl) updateHistory(newValue float32) {

	//loop through the sample history slice in reverse order, excluding 0th index
	for i := len(n.sampleHistory) - 1; i > 0; i-- {
		n.sampleHistory[i] = n.sampleHistory[i-1] //set the current index equal to the value of the previous index
	}

	n.sampleHistory[0] = newValue //set 0th index to new measured value

	/* Now calculate the weighted average of the sample history. Note that if a node is stationary, all values
	averaged over are weighted equally. The faster the node is moving, the less the older values are worth when
	calculating the average, because in that case we want the average to more closely reflect the newer values
	*/
	var sum float32
	var numSamples int //variable for number of samples to average over

	var decreaseRatio = n.speedWeight / 100.0

	if n.totalSamples > len(n.sampleHistory) { //if the node has taken more than x total samples
		numSamples = len(n.sampleHistory) //we only average over the x most recent ones
	} else { //if it doesn't have x samples taken yet
		numSamples = n.totalSamples //we only average over the number of samples it's taken
	}

	for i := 0; i < numSamples; i++ {
		if n.sampleHistory[i] != 0 {
			//weight the values of the sampleHistory when added to the sum variable based on the speed, so older values are weighted less
			sum += n.sampleHistory[i] - ((decreaseRatio) * float32(i))
		} else {
			sum += 0
		}
	}
	n.avg = sum / float32(numSamples)
}

/* this function increments a node's total number of samples by 1
it's called whenever the node takes a new sample */
func (n *NodeImpl) incrementTotalSamples() {
	n.totalSamples++
}

//getter function for average
func (n *NodeImpl) getAvg() float32 {
	return n.avg
}

//increases numResets field
func (n *NodeImpl) incrementNumResets() {
	n.numResets++
}

//setter function for concentration field
func (n *NodeImpl) setConcentration(conc float64) {
	n.concentration = conc
}

//getter function for ID field
func (n *NodeImpl) getID() int {
	return n.id
}

//getter function for x and y locations
func (n *NodeImpl) getLoc() (int, int) {
	return n.x, n.y
}

//setter function for S0
func (n *NodeImpl) setS0(s0 float64) {
	n.S0 = s0
}

//setter function for S1
func (n *NodeImpl) setS1(s1 float64) {
	n.S1 = s1
}

//setter function for S2
func (n *NodeImpl) setS2(s2 float64) {
	n.S2 = s2
}

//setter function for E0
func (n *NodeImpl) setE0(e0 float64) {
	n.E0 = e0
}

//setter function for E1
func (n *NodeImpl) setE1(e1 float64) {
	n.E1 = e1
}

//setter function for E2
func (n *NodeImpl) setE2(e2 float64) {
	n.E2 = e2
}

//setter function for ET1
func (n *NodeImpl) setET1(et1 float64) {
	n.ET1 = et1
}

//setter function for ET2
func (n *NodeImpl) setET2(et2 float64) {
	n.ET2 = et2
}

//getter function for all parameters
func (n *NodeImpl) getParams() (float64, float64, float64, float64, float64, float64, float64, float64) {
	return n.S0, n.S1, n.S2, n.E0, n.E1, n.E2, n.ET1, n.ET2
}

//getter function for just S0 - S2 parameters
func (n *NodeImpl) getCoefficients() (float64, float64, float64) {
	return n.S0, n.S1, n.S2
}

//getter function for x
func (n *NodeImpl) getX() int {
	return n.x
}

//getter function for y
func (n *NodeImpl) getY() int {
	return n.y
}

//This is the actual pinging of the sensor
func (n *NodeImpl) sense() {
	if n.hasCheckedGPS == false {
		n.xPos[n.bufferI] = -1
		n.yPos[n.bufferI] = -1
		n.value[n.bufferI] = n.getValue()
		n.time[n.bufferI] = iterations_used
		n.bufferI = n.bufferI + 1
	} else {
		n.value[n.bufferI] = n.getValue()
	}
}

//This is the actual pinging of the GPS
func (n *NodeImpl) GPS() {
	if n.hasCheckedSensor == false {
		n.value[n.bufferI] = -1
		n.xPos[n.bufferI] = n.x
		n.yPos[n.bufferI] = n.y
		n.time[n.bufferI] = iterations_used
		n.bufferI = n.bufferI + 1
	} else {
		if n.bufferI > 0 {
			n.xPos[n.bufferI-1] = n.x
			n.yPos[n.bufferI-1] = n.y
		}
	}
}

//This is the actual upload to the server
func (n *NodeImpl) server() {
	//getData(&s,n.xPos[0:n.bufferI],n.yPos[0:n.bufferI],n.value[0:n.bufferI],n.time[0:n.bufferI], n.id,n.bufferI)
	n.bufferI = 0
}

//Returns node distance to the bomb
func (n *NodeImpl) geoDist(b bomb) float32 {
	//this needs to be changed
	return float32(math.Pow(float64(math.Abs(float64(n.x)-float64(b.x))), 2) + math.Pow(float64(math.Abs(float64(n.y)-float64(b.y))), 2))
}

//Returns array of accelerometer speeds recorded for a specific node
func (n *NodeImpl) getSpeed() []float32 {
	return n.accelerometerSpeed
}

//Returns a different version of the distance to the bomb
func (n *NodeImpl) getValue() int {
	return int(math.Sqrt(math.Pow(float64(n.x-b.x), 2) + math.Pow(float64(n.y-b.y), 2)))
}

//Takes cares of taking a node's readings and printing detections and stuff
func (curNode *NodeImpl) getReadings() {

	//driftFile.Sync()
	//nodeFile.Sync()
	//positionFile.Sync()
	//test change
	newX, newY := curNode.getLoc()

	newDist := curNode.distance(*b) //this is the node's reported value without error

	//Calculate error, sensitivity, and noise, as per the matlab code
	S0, S1, S2, E0, E1, E2, ET1, ET2 := curNode.getParams()
	sError := (S0 + E0) + (S1+E1)*math.Exp(-float64(curNode.nodeTime)/(Tau1+ET1)) + (S2+E2)*math.Exp(-float64(curNode.nodeTime)/(Tau2+ET2))
	curNode.sensitivity = S0 + (S1)*math.Exp(-float64(curNode.nodeTime)/Tau1) + (S2)*math.Exp(-float64(curNode.nodeTime)/Tau2)
	sNoise := rand.NormFloat64()*0.5*errorModifierCM + float64(newDist)*sError

	errorDist := sNoise / curNode.sensitivity //this is the node's actual reading with error

	//increment node time
	curNode.nodeTime++

	if curNode.hasCheckedSensor {
		curNode.incrementTotalSamples()
		curNode.updateHistory(float32(errorDist))
	}

	//Detection of false positives or false negatives
	if errorDist < detectionThreshold && float64(newDist) >= detectionThreshold {
		//this is a node false negative due to drifitng
		if curNode.hasCheckedSensor {
			//just drifting
			fmt.Fprintln(driftFile, "Node False Negative (drifting) ID:", curNode.id, "True Reading:", newDist, "Drifted Reading:",
				errorDist, "S0:", curNode.S0, "S1:", curNode.S1, "S2:", curNode.S2, "E0:", curNode.E0, "E1:", curNode.E1,
				"E2:", curNode.E2, "ET1:", curNode.ET1, "ET2:", curNode.ET2, "time since calibration:", curNode.nodeTime,
				"Current Time (Iter):", iterations_used, "Energy Level:", curNode.battery, "Distance from bomb:", math.Sqrt(float64(curNode.geoDist(*b))),
				"x:", curNode.x, "y:", curNode.y)
		} else {
			//both drifting and energy
			fmt.Fprintln(driftFile, "Node False Negative (both) ID:", curNode.id, "True Reading:", newDist, "Drifted Reading:",
				errorDist, "S0:", curNode.S0, "S1:", curNode.S1, "S2:", curNode.S2, "E0:", curNode.E0, "E1:", curNode.E1,
				"E2:", curNode.E2, "ET1:", curNode.ET1, "ET2:", curNode.ET2, "time since calibration:", curNode.nodeTime,
				"Current Time (Iter):", iterations_used, "Energy Level:", curNode.battery, "Distance from bomb:", math.Sqrt(float64(curNode.geoDist(*b))),
				"x:", curNode.x, "y:", curNode.y)
		}
	}

	if errorDist >= detectionThreshold && float64(newDist) >= detectionThreshold && !curNode.hasCheckedSensor {
		//false negative due solely to energy
		fmt.Fprintln(driftFile, "Node False Negative (energy) ID:", curNode.id, "True Reading:", newDist, "Drifted Reading:",
			errorDist, "S0:", curNode.S0, "S1:", curNode.S1, "S2:", curNode.S2, "E0:", curNode.E0, "E1:", curNode.E1,
			"E2:", curNode.E2, "ET1:", curNode.ET1, "ET2:", curNode.ET2, "time since calibration:", curNode.nodeTime,
			"Current Time (Iter):", iterations_used, "Energy Level:", curNode.battery, "Distance from bomb:", math.Sqrt(float64(curNode.geoDist(*b))),
			"x:", curNode.x, "y:", curNode.y)
	}

	if errorDist >= detectionThreshold && float64(newDist) < detectionThreshold {
		//this if a false positive
		//it must be due to drifting. Report relevant info to driftFile
		fmt.Fprintln(driftFile, "Node False Positive (drifting) ID:", curNode.id, "True Reading:", newDist, "Drifted Reading:",
			errorDist, "S0:", curNode.S0, "S1:", curNode.S1, "S2:", curNode.S2, "E0:", curNode.E0, "E1:", curNode.E1,
			"E2:", curNode.E2, "ET1:", curNode.ET1, "ET2:", curNode.ET2, "time since calibration:", curNode.nodeTime,
			"Current Time (Iter):", iterations_used, "Energy Level:", curNode.battery, "Distance from bomb:", math.Sqrt(float64(curNode.geoDist(*b))),
			"x:", curNode.x, "y:", curNode.y)
	}

	if errorDist >= detectionThreshold && float64(newDist) >= detectionThreshold && curNode.hasCheckedSensor {
		fmt.Fprintln(driftFile, "Node True Positive ID:", curNode.id, "True Reading:", newDist, "Drifted Reading:",
			errorDist, "S0:", curNode.S0, "S1:", curNode.S1, "S2:", curNode.S2, "E0:", curNode.E0, "E1:", curNode.E1,
			"E2:", curNode.E2, "ET1:", curNode.ET1, "ET2:", curNode.ET2, "time since calibration:", curNode.nodeTime,
			"Current Time (Iter):", iterations_used, "Energy Level:", curNode.battery, "Distance from bomb:", math.Sqrt(float64(curNode.geoDist(*b))),
			"x:", curNode.x, "y:", curNode.y)
	}

	//If the reading is more than 2 standard deviations away from the grid average, then recalibrate
	gridAverage := grid[curNode.row(yDiv)][curNode.col(xDiv)].avg
	//standDev := grid[curNode.row(yDiv)][curNode.col(xDiv)].stdDev

	//New condition added: also recalibrate when the node's sensitivity is <= 1/10 of its original sensitvity
	//New condition added: Check to make sure the sensor was pinged this iteration
	if ((curNode.sensitivity <= (curNode.initialSensitivity / 2)) && (curNode.hasCheckedSensor)) && iterations_used != 0 {
		curNode.recalibrate()
		recalibrate = true
		curNode.incrementNumResets()
	}

	//printing statements to log files, only if the sensor was pinged this iteration
	//if curNode.hasCheckedSensor && nodesPrint{
	if nodesPrint {
		if recalibrate {
			fmt.Fprintln(nodeFile, "ID:", curNode.getID(), "Average:", curNode.getAvg(), "Reading:", newDist, "Error Reading:", errorDist, "Recalibrated")
		} else {
			fmt.Fprintln(nodeFile, "ID:", curNode.getID(), "Average:", curNode.getAvg(), "Reading:", newDist, "Error Reading:", errorDist)
		}
		//fmt.Fprintln(nodeFile, "battery:", int(curNode.battery),)
	}

	if positionPrint {
		fmt.Fprintln(positionFile, "ID:", curNode.getID(), "x:", newX, "y:", newY)
	}

	recalibrate = false

	//Receives the node's distance and calculates its running average
	//for that square
	//Only do this if the sensor was pinged this iteration
	if curNode.hasCheckedSensor {
		grid[curNode.row(yDiv)][curNode.col(xDiv)].takeMeasurement(float32(errorDist))
		grid[curNode.row(yDiv)][curNode.col(xDiv)].numNodes++
		//subtract grid average from node average, square it, and add it to this variable
		grid[curNode.row(yDiv)][curNode.col(xDiv)].squareValues += (math.Pow(float64(errorDist-float64(gridAverage)), 2))
	}
}
