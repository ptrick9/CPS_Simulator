

# General Description
The bomb simulator is a program designed to simulate the detection of a bomb utilizing cheap chemical sensors carried by every member of the event. Specifically, the simulator attempts to emulate realistic detector equations for sensor drifting, energy usage, and communication overheads. The simulator utilizes a general flow which will be listed and described below.
## Flow
1. Parse User Input
   * [Arena Map](#arena-map)
   * [Energy Models](#energy-models)
   * [Super Node Types](#super-node-types)
   * [Sensor Drifting Parameters](#sensor-drifting-parameters)
   * [Node Sampling Parameters](#node-sampling-parameters)
   * [Logging Parameters](#logging-parameters)
2. Build and Populate Arena for simulation
3. Execute Simulator with given parameters
4. Output Log Files with Details and Results of Simulation

### [Quick Link To Tutorial](#tutorial)


## General Detection
The simulator runs with the given parameters until it detects a bomb. In order to detect a bomb, a few things must happen. A node (person) must walk near enough to a bomb for it to create a high reading on the node's sensor. The node's individual energy model must also decide to take a sample within the time period that the node is within the detection radius of the bomb. At some point in the future, again, as determined by the energy model, the node will send its data to the server. 

Upon receiving the data from the node, the server will place the readings into a map of the arena, allowing for spatial correlation of multiple readings. In order to better correlate readings and decrease the number of false positives caused by sensor error, the server averages readings over areas. These areas are represented as large (user settable) squares in which the most recent readings in that area are averaged together. Once a square has a high enough average to trigger a detection, a supernode (a highly accurate node with greater sensor accuracy and detection distance) is routed to that area to ensure that there is no false alarm. 

# User Input
The simulator is designed with the philosophy of being as reconfigurable as possible. While this leads to a lengthy configuration string required for a launch, it provides the greatest ability for customization and reconfiguration as new ideas or parameters come to light. 

## Arena Map 
The arena map is an input file which specifies multiple parameters for the simulator. Notably, the file provides 6 values at the top of the file

```
numNodes-1000
superNodeType-1
maxX-500
maxY-100
bombX-250
bombY-22
```

* numNodes
   * The number of total nodes to be utilized in the simulation
* superNodeType
   * The routing model utilized by super nodes. There are many models for this which will be discussed later.
* maxX
   * Maximum X dimension for the arena
* maxY
   * Maximum Y dimension for the arena
* bombX
   * X Coordinate location of the bomb
* bombY
   * Y Coordinate location of the bomb


After these beginning arguments, is another four sections of inputs. 

```
N: 2
x:62, y:28, t:55

x:62, y:15, t:12

W: 2
x:46, y:3

x:46, y:9

S: 1
x:67, y:35, t:56

POIS: 3
x:20, y:20, ti:0, to:500

x:480, y:20, ti:0, to:500

x:250, y:20, ti:500, to:2000
```

These lines represent variables within the simulation, specifically:
* N represents the nodes section
* W represents the walls section
* S represeents the supernode section
* POIS represents points of interest for the nodes 

Each section follows a similar format. First, we specify the section title and number of items contained, for example: `N: 2` specifies that there are 2 entries in the nodes section. 

The format of the entries is as follows
* Node entries
   * X coordinate, Y coordinate, time at which node enters simulation
* Wall entries
   * X coordinate, Y coordinate
* Super Node Entries
   * X coordinate, Y coordinate, time at which super node enters simulation
* Points of Interest Entries
   * X coordiante, Y coordinate, time at which POI enters simulation, time at which POI is removed from simulation

----

## Energy Models 
The energy model has many different parameters that control it's main function. To be most efficient, the GPS, sensor, and server are all sampled at different times and adaptively based upon movement speed, battery level, and number of stored samples waiting to be transmitted. The command line arguments that control the battery model are listed and described below:

* naturalLoss
    * Loss over every iteration of the simulator, essentially battery loss per second. 
    * 0 < naturalLoss < .1
* sensorSamplingLoss
    * Loss from sampling the sensor. 
    * 0 < sensorSamplingLoss < .1
* GPSSamplingLoss
    * Loss from sampling the GPS unit. 
    * 0 < GPSSamplingLoss < .1
* serverSamplingLoss
    * Loss from sending data to the server. 
    * 0 < serverSapmlingLoss < .1
* thresholdBatteryToHave
    * Minimum amount of battery to have at the end of simulation. If simulator detects that the battery will fall beneath this level, it stops sampling and sending information to the server. 
    * 0 < thresholdBatteryToHave < 50
* thresholdBatteryToUse
    * The amount of battery to dedicate to bomb detection. 
    * 0 < thresholdBatteryToUse < 20
* movementSamplingSpeed
    * Speed at which a mode must be moving in order for it to sample ahead of time.
    * 0 < movementSamplingSpeed < 100
* movementSamplingPeriod
    * How often to sample the node GPS. 
    * 0 < movementSamplingPeriod < 100
* maxBufferCapacity
    * Number of samples to hold before we send to the server. 
    * 10 < maxBufferCapacity < 100
* sensorSamplingPeriod
    * How often to sample the sensor
    * 0 < sensorSamplingPeriod < 100
* GPSSamplingPeriod
    * How often to sample the GPS
    * 0 < GPSSamplingPeriod < 100
* serverSamplingPeriod
    * How often to send samples to the server
    * 0 < serverSamplingPeriod < 100

----

## Super Node Types

There are 7 types of supernodes currently available. Supernodes have 2 major parameters, whether they can be optimized or not, and whether they require a square area to work correctly. 

A super node that can be optimized is one that upon reaching a destination re-examines the list of destinations to visit next and optimizes the path. 

Some super node routing strategies break the area into regions and give supernodes responsibility over a specific region. Due to the nature of these algorithms, many of them require a square to work at the moment. 

Quick Links to the different types:
* [Type 0](###-Super-Node-Type-0)
* [Type 1](###-Super-Node-Type-1)
* [Type 2](###-Super-Node-Type-2)
* [Type 3](###-Super-Node-Type-3)
* [Type 4](###-Super-Node-Type-4)
* [Type 5](###-Super-Node-Type-5)
* [Type 6](###-Super-Node-Type-6)
* [Type 7](###-Super-Node-Type-7)

----

### Super Node Type 0


Number of Super Nodes: NO RESTRICTIONS  
Can Be Optimized: NO  
Square Grid Only: NO

Super nodes of type 0 are the simplest super node. Their behavior is the least sophisticated of all the super nodes. When the scheduler is notified of a point of interest it immediately adds it to the nearest super node. That super node adds the new point of interest to its current routePath.

----

### Super Node Type 1
Number of Super Nodes: NO RESTRICTIONS  
Can Be Optimized: YES  
Square Grid Only: NO

Super nodes of type 1 are an improved version of super nodes of type 0. When the scheduler is notified of a point of interest it finds the super node with the smallest total nodeDist value. This value is the sum of the length of the super node's current routePath and the distance from the new point of interest to the closest point in the routePath. If the super node's routePath is empty, the distance from the new point of interest to the super node is used as the nodeDist.

When a new point of interest is added to a super node of type 1 it does not immediately add it to the end of the routePath. To ensure the super node is visiting all of its routePoints in the minimum distance required, it finds the place in the current routePath that demands the least distance to divert the path to the new point of interest. As a last resort it adds the point of interest to the end of the routePath.

The algorithm of super nodes of type 1 can be optimized. By setting the ```doOptimize``` flag to ```true``` the behavior of the super nodes can be changed. Once a super node reaches a point of interest all the current points being visited by the super nodes are reorganized to more efficiently.

----

### Super Node Type 2
Number of Super Nodes: 4  
Can Be Optimized: NO  
Square Grid Only: YES  

Super nodes of type 2 follow the same algorithm that super nodes of type 1 except the area they cover is restricted. Instead of covering the entirety of the grid these super nodes only visit points of interest inside circles in the corners of the grid. If a line was drawn from the corner of the grid to the center, the center of these circles would be the midpoint of that line with the radii extending from the center of the circle to the corner of the grid and to the center of the grid. 

The benefit of restricting the area each super node covers is that it minimizes the total number of points of interest each super node needs to visit. These circles also overlap slightly, allowing for super nodes that are particularly busy to be aided by bordering super nodes that can visit points inside these areas.

----

### Super Node Type 3
Number of Super Nodes: 4  
Can Be Optimized: NO  
Square Grid Only: YES  

Super nodes of type 3 operate exactly the same as super nodes of type 2 except the circles are centered at the midpoint of each side of the grid. These circles have a radius equal to half the length/width of the square grid. The overlapping areas of these circles is larger than those of super nodes of type 2. The increased total area, as well as the increased area of overlap, serve to allow busy super nodes to be aided even more easily by other super nodes.

----

### Super Node Type 4
Number of Super Nodes: 4  
Can Be Optimized: NO  
Square Grid Only: YES  

Super nodes of type 4 also operate exactly the same as super nodes of type 2 and 3. The circles the super nodes of type 4 are restricted inside of are centered inside the corners of the grid. These circle's radii are the length of the diagonal from the corner of the grid to the center of the grid. Now no super node has any area where only it is the only super node that can reach a point of interest. The area of overlap sometimes contains all four super nodes. The increased overlapping area attempts to prevent any one super node from being completely overwhelmed; however the increased total area any one super node can travel adds to this potential strain.

----

### Super Node Type 5
Number of Super Nodes: 1, 2, 4  
Can Be Optimized: NO  
Square Grid Only: NO  

Super nodes of type 5 used the same routing algorithm as super nodes of type 1 but instead of being restricted to the inside of a circle like other super node types, these super nodes are restricted inside of regions with no overlap. These regions can be the size of the entire grid, halves of the grid, or quarters of the grid. When the scheduler is notified of a new point of interest it adds it to the super node who's region the point falls inside. 

----

### Super Node Type 6
Number of Super Nodes: 1, 2, 4  
Can Be Optimized: NO  
Square Grid Only: NO  

Super nodes of type 6 are restricted into the same exact regions as super nodes of type 5. The routing algorithms, however, are entirely different. Super nodes of type 6 begin centered inside their respective regions. Every iteration each super node divides its respective region into four quadrants. Looping through the quadrants the super node finds the quadrant with the most points of interest. The super node then splits the quadrant along the diagonal from the center to corner of the quadrant. Starting with the triangle that contains the most points of interest the super node plots a path from the closest point to the super node all the way to the furthest point from the super node. It then plots a path back towards the center going through all the points in the
other triangle, from furthest from the super node to closest. 

Once all the points of interest have been visited the super node moves back to the center. Once a quadrant has been visited it cannot revisit that quadrant until all other quadrants have been visited or all the other quadrants are empty.

----

### Super Node Type 7
Number of Super Nodes: 1, 2, 4  
Can Be Optimized: NO  
Square Grid Only: NO  

Super nodes of type 7 are exactly the same as super nodes of type 6 except intead of finding the quadrant with the most points of interest inside of it, it finds the quadrant with the oldest super node inside of it. All the same restrictions and behaviors apply.

----

## Sensor Drifting Parameters 

Sensor drifting is handled by the following equations:



<p align="center"> 
 <a href="https://www.codecogs.com/eqnedit.php?latex=$$&space;Concentration&space;=&space;\frac{1000}{\sqrt[3]{\frac{Distance}{.2}}}&space;$$\\" target="_blank"><img src="https://latex.codecogs.com/png.latex?$$&space;Concentration&space;=&space;\frac{1000}{\sqrt[3]{\frac{Distance}{.2}}}&space;$$\\" title="$$ Concentration = \frac{1000}{\sqrt[3]{\frac{Distance}{.2}}} $$\\" /></a>
 </p>



<p align="center"> 
<a href="https://www.codecogs.com/eqnedit.php?latex=$$&space;Sensitivity&space;=&space;(S0&plus;E0)&space;&plus;&space;(S1&plus;E1)e^{\frac{-t(i)}{\tau1&space;&plus;&space;ET1}}&space;&plus;&space;(S2&plus;E2)e^{\frac{-t(i)}{\tau2&space;&plus;&space;ET2}}&space;$$" target="_blank"><img src="https://latex.codecogs.com/png.latex?$$&space;Sensitivity&space;=&space;(S0&plus;E0)&space;&plus;&space;(S1&plus;E1)e^{\frac{-t(i)}{\tau1&space;&plus;&space;ET1}}&space;&plus;&space;(S2&plus;E2)e^{\frac{-t(i)}{\tau2&space;&plus;&space;ET2}}&space;$$" title="$$ Sensitivity = (S0+E0) + (S1+E1)e^{\frac{-t(i)}{\tau1 + ET1}} + (S2+E2)e^{\frac{-t(i)}{\tau2 + ET2}} $$" /></a>
</p>
 

<p align="center"> 
<a href="https://www.codecogs.com/eqnedit.php?latex=$$&space;SensitivityEstimate&space;=&space;S0&space;&plus;&space;S1e^{\frac{-t(i)}{\tau1}}&space;&plus;&space;S2e^{\frac{-t(i)}{\tau2}}&space;$$" target="_blank"><img src="https://latex.codecogs.com/png.latex?$$&space;SensitivityEstimate&space;=&space;S0&space;&plus;&space;S1e^{\frac{-t(i)}{\tau1}}&space;&plus;&space;S2e^{\frac{-t(i)}{\tau2}}&space;$$" title="$$ SensitivityEstimate = S0 + S1e^{\frac{-t(i)}{\tau1}} + S2e^{\frac{-t(i)}{\tau2}} $$" /></a>
</p>


<p align="center"> 
<a href="https://www.codecogs.com/eqnedit.php?latex=$$&space;MeasurementNoise&space;=&space;gaussian*0.5&space;&plus;&space;Concentration*Sensitivity&space;$$" target="_blank"><img src="https://latex.codecogs.com/png.latex?$$&space;MeasurementNoise&space;=&space;gaussian*0.5&space;&plus;&space;Concentration*Sensitivity&space;$$" title="$$ MeasurementNoise = gaussian*0.5 + Concentration*Sensitivity $$" /></a>
</p>


<p align="center"> 
<a href="https://www.codecogs.com/eqnedit.php?latex=$$&space;MeasurementExtimation&space;=&space;\frac{MeasurementNoise}{SensitivityEstimate}&space;$$" target="_blank"><img src="https://latex.codecogs.com/png.latex?$$&space;MeasurementExtimation&space;=&space;\frac{MeasurementNoise}{SensitivityEstimate}&space;$$" title="$$ MeasurementExtimation = \frac{MeasurementNoise}{SensitivityEstimate} $$" /></a>
</p>

 Currently, <a href="https://www.codecogs.com/eqnedit.php?latex=\inline&space;$\tau1=10$" target="_blank"><img src="https://latex.codecogs.com/png.latex?\inline&space;$\tau1=10$" title="$\tau1=10$" /></a> and <a href="https://www.codecogs.com/eqnedit.php?latex=\inline&space;$\tau2=500$" target="_blank"><img src="https://latex.codecogs.com/png.latex?\inline&space;$\tau2=500$" title="$\tau2=500$" /></a>. S0, S1, S2 are all sensor specific and randomly chosen from a random distribution such that $rand()*0.2+0.1$. E0, E1, E2 are chosen as
 

<p align="center"> 
<a href="https://www.codecogs.com/eqnedit.php?latex=$$&space;E0&space;=&space;rand()*0.1*S0&space;$$" target="_blank"><img src="https://latex.codecogs.com/png.latex?$$&space;E0&space;=&space;rand()*0.1*S0&space;$$" title="$$ E0 = rand()*0.1*S0 $$" /></a>
</p>



 <p align="center"> 
<a href="https://www.codecogs.com/eqnedit.php?latex=$$&space;E1&space;=&space;rand()*0.1*S1&space;$$" target="_blank"><img src="https://latex.codecogs.com/png.latex?$$&space;E1&space;=&space;rand()*0.1*S1&space;$$" title="$$ E1 = rand()*0.1*S1 $$" /></a>
</p>
 

<p align="center"> 
<a href="https://www.codecogs.com/eqnedit.php?latex=$$&space;E1&space;=&space;rand()*0.1*S1&space;$$" target="_blank"><img src="https://latex.codecogs.com/png.latex?$$&space;E2&space;=&space;rand()*0.1*S2&space;$$" title="$$ E2 = rand()*0.1*S2 $$" /></a>
</p>

 and ET1 and ET2 are chosen as


 <p align="center"> 
<a href="https://www.codecogs.com/eqnedit.php?latex=$$&space;ET1&space;=&space;\tau1*rand()*0.05&space;$$" target="_blank"><img src="https://latex.codecogs.com/png.latex?$$&space;ET1&space;=&space;\tau1*rand()*0.05&space;$$" title="$$ ET1 = \tau1*rand()*0.05 $$" /></a>
</p>

 <p align="center"> 
<a href="https://www.codecogs.com/eqnedit.php?latex=$$&space;ET2&space;=&space;\tau2*rand()*0.05&space;$$" target="_blank"><img src="https://latex.codecogs.com/png.latex?$$&space;ET2&space;=&space;\tau2*rand()*0.05&space;$$" title="$$ ET2 = \tau2*rand()*0.05 $$" /></a>
</p>

In all cases, each value is chosen and remembered per node, meaning that each node has a slightly different set of parameters and error values that are inherent to it. Time is kept track of per node, with time 0 being when the node is initialized and whenever the node is recalibrated. Currently, none of these parameters are exposed via the command line as they were specified as unchanging. 

Recalibration is triggered by two conditions. In the first case, when the sensitivity has fallen to 50% of the initial sensitivity, the the node is recalibrated. Furthermore, when a node's reading is more than two standard deviations away from the average for the area it is reporting in, the node is told to recalibrate by the server. A recalibration resets the node's sensitivity back to its initial sensitivity. 

----

## Node Sampling Parameters 

Node sampling is handled mainly by the energy models as described in [Energy Models](#energy-models). Once data is sampled it is stored in a few places.

Each node keeps a running average of its own readings. This value is implemented as a weighted moving average in which readings are rated equally if a node is sitting still, or more heavily as a node begins to move. This average value is utilized for debugging and diagnostics.

Each grid square also keeps a running average of readings within the grid in order to both trigger recalibrations if node readings are unreasonable and to trigger a supernode visit if readings are above a detection threshold.

There are a number of parameters that can be set from the command line in order to better control sampling.

* gridStoredSamples
   * The number of samples that will be utilized in the running average for each grid square. This can be viewed as a FIFO queue of chosen size in which all items in the queue are averaged.
   * 5 <= gridStoredSamples <= 25

* numStoredSamples
    * The number of samples considered in the moving average of each node
    * 5 <= numStoredSamples <= 25

* detectionThreshold
    * Threshold for triggering a supernode to be routed to the grid position. Increasing this threshold decreases false positives, but can lead to many false negatives if the grid size is much larger than the detection radius of the individual nodes. 
    * 1 <= detectionThreshold <= 1000

* squareRow
    * Utilized as a divisor for the Y coordinate size and dictates the number rows that the entire arena will be divided into. Linked closely with squareCol, which determines the number of columns in the arena. Individual square dimensions can be calculated as $\frac{maxX}{squareCol}$ x $\frac{maxX}{squareRow}$

* squareCol
    * Utilized as a divisor for the X coordinate size and dictates the number columns that the entire arena will be divided into. Linked closely with squareRow, which determines the number of rows in the arena. Individual square dimensions can be calculated as $\frac{maxX}{squareCol}$ x $\frac{maxX}{squareRow}$

----

## Logging Parameters 

Logging for each run of the simulator is controlled by command line parameters. Logging can be extremely verbose or limited to just important statistics such as detections.

* logPosition
    * Determines whether to print each nodes position at every update cycle. This can be utilized for viewing with the included BombDetection.jar allowing for viewing of crowd movement over time.
    * true or false

* logEnergy
    * Prints energy usage for each node after every update cycle. Useful for debugging energy models.
    * true or false

* logNodes
    * Print node readings and averages after every update cycle.
    * true or false

* logGrid
    * Print grid readings and averages after every update cycle
    * true or false

* outputFileName
    * Filename prefix for output files. For example, the drifting log file is named outputFileName_drift.txt

The most important information is printed in the drift file. This file contains information as to the types of detections, when, and where they occur. Possible detections are:
### Grid 
* Grid False Negative
    * The bomb is located in this grid square, but the nodes that have taken readings in this grid square were either too far away from the bomb or had drifted too much to detect.

* Grid False Positive
    * The nodes that have taken readings in this square have drifted too much, or the bomb is located close to a neighboring grid edge, leading to a high reading in this grid square that is incorrect.

* Grid True Positive
    * The detection in this square is above the detection threshold and the bomb is in this square

### Drift/Energy

* Drifting False Negative
    * The sensor has drifted far enough that the value is larger than the deteciton threshold, but the bomb is not near enough to be detected

* Energy False Negative
    * The sensor has not yet drifted far enough to read a false negative if it were sampled **and** the energy model didn't allow the node to take a sample at this time step. Thus leading to a missed detection or false negative.

* Drifting and Energy False Negative
    * The sensor has drifted far enough that it wouldn't have detected the bomb that was in range, and the energy model didn't allow the sensor to sample, meaning that even if the node had sampled it still wouldn't have detected the bomb.

* Drifting False Positive
    * The sensor has drifted far enought that it has incorrectly detected a bomb reading. The sensor must not be in range of the bomb for this to occur.

----

# Tutorial

The simulator project comes with a bash script that builds and runs a test project. In order for everything to compile correctly, some dependencies are required:

* [Golang](https://golang.org/dl/) must be installed 
    * Choose most recent version for your system, install, and restart any terminal
* [Java](http://www.oracle.com/technetwork/java/javase/downloads/jdk10-downloads-4416644.html) w/ JDK must be installed in order to run viewer. Accept agreement and install
    * If on LINUX, run ``` sudo apt-get install openjdk openjfx ``` which will install the required packages. If you are running Ubuntu 18.04, this does not seem to work. 
* [Python3](https://www.python.org/downloads/) with numpy, pandas, matplotlib, and jupyter must be installed if you want to use the included statistics processing. 
    * 

