/*
-logNodes=false
-logPosition=true
-logGrid=false
-logEnergy=false
-regionRouting=true
-noEnergy=true
-csvMove=true
-zipFiles=true
-windRegionPath=hull_fine_bomb_9x9.txt
-inputFileName=Scenario_3.txt
-imageFileName=marathon_street_map.png
-stimFileName=circle_0.txt
-outRoutingStatsName=routingStats.txt
-iterations=11000
-superNodes=false
-doOptimize=false
-movementPath=C:/Users/patrick/Downloads/marathon_movement/marathon2_2000_3.scb
-totalNodes=2000
-detectionThreshold=5
-detectionDistance=6
-sittingStopThreshold=5
-negativeSittingStopThreshold=-10
-GridCapacityPercentage=0.900000
-naturalLoss=0.005000
-sensorSamplingLoss=0.001000
-GPSSamplingLoss=0.005000
-serverSamplingLoss=0.010000
-SamplingLossBTCM=0.000100
-SamplingLossWifiCM=0.001000
-SamplingLoss4GCM=0.005000
-SamplingLossAccelCM=0.001000
-thresholdBatteryToHave=30
-thresholdBatteryToUse=10
-movementSamplingSpeed=20
-movementSamplingPeriod=20
-maxBufferCapacity=25
-sensorSamplingPeriod=1000
-GPSSamplingPeriod=1000
-serverSamplingPeriod=1000
-nodeStoredSamples=10
-GridStoredSamples=10
-errorMultiplier=0.60000
-numSuperNodes=4
-RecalibrationThreshold=3
-StandardDeviationThreshold=1.700000
-SuperNodeSpeed=3
-SquareRowCM=60
-SquareColCM=320
-validationThreshold=2
-serverRecal=true
-driftExplorer=true
-commandBomb=false
-fineSensorPath=C:/Users/patrick/Dropbox/Patrick/udel/SUMMER2019/GitSimulator/fine_bomb9x9.csv
-csvSensor=false
-OutputFileName=C:/Users/patrick/Downloads/testFolder/
-detectionWindow=59
-moveSize=4000
*/

/*
-logNodes=false
-logPosition=true
-logGrid=false
-logEnergy=false
-regionRouting=true
-noEnergy=true
-csvMove=true
-zipFiles=true
-windRegionPath=hull_fine_bomb_9x9.txt
-inputFileName=Scenario_4.txt
-imageFileName=C:/Users/patrick/Downloads/fireflyEdits/Firefly.png
-stimFileName=circle_0.txt
-outRoutingStatsName=routingStats.txt
-iterations=11000
-superNodes=false
-doOptimize=false
-movementPath=C:/Users/patrick/Downloads/fireflyEdits/firefly.scb
-totalNodes=1000
-detectionThreshold=5
-detectionDistance=6
-sittingStopThreshold=5
-negativeSittingStopThreshold=-10
-GridCapacityPercentage=0.900000
-naturalLoss=0.005000
-sensorSamplingLoss=0.001000
-GPSSamplingLoss=0.005000
-serverSamplingLoss=0.010000
-SamplingLossBTCM=0.000100
-SamplingLossWifiCM=0.001000
-SamplingLoss4GCM=0.005000
-SamplingLossAccelCM=0.001000
-thresholdBatteryToHave=30
-thresholdBatteryToUse=10
-movementSamplingSpeed=20
-movementSamplingPeriod=20
-maxBufferCapacity=25
-sensorSamplingPeriod=1000
-GPSSamplingPeriod=1000
-serverSamplingPeriod=1000
-nodeStoredSamples=10
-GridStoredSamples=10
-errorMultiplier=0.60000
-numSuperNodes=4
-RecalibrationThreshold=3
-StandardDeviationThreshold=1.700000
-SuperNodeSpeed=3
-SquareRowCM=372
-SquareColCM=288
-validationThreshold=2
-serverRecal=true
-driftExplorer=true
-commandBomb=false
-fineSensorPath=C:/Users/patrick/Dropbox/Patrick/udel/SUMMER2019/GitSimulator/fine_bomb9x9.csv
-csvSensor=false
-OutputFileName=C:/Users/patrick/Downloads/testFolder/
-detectionWindow=59
-moveSize=4000
*/

/*
-logNodes=false
-logPosition=true
-logGrid=false
-logEnergy=false
-regionRouting=true
-noEnergy=true
-csvMove=true
-zipFiles=true
-windRegionPath=hull_fine_bomb_9x9.txt
-inputFileName=Scenario_3.txt
-imageFileName=marathon_street_map.png
-stimFileName=circle_0.txt
-outRoutingStatsName=routingStats.txt
-iterations=11000
-superNodes=false
-doOptimize=false
-movementPath=C:/Users/patrick/Downloads/marathon_movement/marathon2_2000_3.scb
-totalNodes=2000
-detectionThreshold=5
-detectionDistance=6
-sittingStopThreshold=5
-negativeSittingStopThreshold=-10
-GridCapacityPercentage=0.900000
-naturalLoss=0.005000
-sensorSamplingLoss=0.001000
-GPSSamplingLoss=0.005000
-serverSamplingLoss=0.010000
-SamplingLossBTCM=0.000100
-SamplingLossWifiCM=0.001000
-SamplingLoss4GCM=0.005000
-SamplingLossAccelCM=0.001000
-thresholdBatteryToHave=30
-thresholdBatteryToUse=10
-movementSamplingSpeed=20
-movementSamplingPeriod=20
-maxBufferCapacity=25
-sensorSamplingPeriod=1000
-GPSSamplingPeriod=1000
-serverSamplingPeriod=1000
-nodeStoredSamples=10
-GridStoredSamples=10
-errorMultiplier=0.60000
-numSuperNodes=4
-RecalibrationThreshold=3
-StandardDeviationThreshold=1.700000
-SuperNodeSpeed=3
-SquareRowCM=60
-SquareColCM=320
-validationThreshold=2
-serverRecal=true
-driftExplorer=false
-commandBomb=false
-fineSensorPath=C:/Users/patrick/Dropbox/Patrick/udel/SUMMER2019/GitSimulator/fine_bomb9x9.csv
-csvSensor=false
-OutputFileName=C:/Users/patrick/Downloads/testFolder/
-detectionWindow=59
-moveSize=4000
*/

package main

import (
	//"CPS_Simulator/simulator/cps"
	"./cps"
	"bytes"
	"container/heap"
	"math"
	//"CPS_Simulator/simulator/cps"
	"fmt"
	"log"
	"math/rand"
	"os"
	"runtime"
	"runtime/pprof"
	"time"

	"image"
	"image/png"
)

var (
	p *cps.Params
	r *cps.RegionParams

	err error

	// End the command line variables.

)

func init() {
	image.RegisterFormat("png", "png", png.Decode, png.DecodeConfig)
}

func main() {

	Events := make(cps.PriorityQueue, 0)

	heap.Init(&Events)

	p = &cps.Params{}
	r = &cps.RegionParams{}

	p.Events = Events
	p.Server = cps.FusionCenter{P: p, R: r}

	p.Tau1 = 3500
	p.Tau2 = 9000
	p.FoundBomb = false

	rand.Seed(time.Now().UTC().UnixNano())

	//getFlags()
	fmt.Fprintf(p.RunParamFile, "Starting file\n")
	cps.GetFlags(p)

	fmt.Println("Getting Wind regions...")
	cps.ReadWindRegion(p)
	fmt.Println("Done!")

	if p.CPUProfile != "" {
		f, err := os.Create(p.CPUProfile)
		if err != nil {
			log.Fatal("could not create CPU profile: ", err)
		}
		defer f.Close()
		if err := pprof.StartCPUProfile(f); err != nil {
			log.Fatal("could not start CPU profile: ", err)
		}
		defer pprof.StopCPUProfile()
	}

	p.FileName = p.InputFileNameCM
	cps.GetListedInput(p)

	//p.SquareRowCM = getDashedInput("p.SquareRowCM")
	//p.SquareColCM = getDashedInput("p.SquareColCM")
	//p.TotalNodes = cps.GetDashedInput("numNodes", p)    //Now a command line argument
	//numStoredSamples = getDashedInput("numStoredSamples")
	p.MaxX = cps.GetDashedInput("maxX", p)
	p.MaxY = cps.GetDashedInput("maxY", p)
	p.BombX = cps.GetDashedInput("bombX", p)
	p.BombY = cps.GetDashedInput("bombY", p)

	p.B = &cps.Bomb{X: p.BombX, Y: p.BombY}

	//numAtt = getDashedInput("numAtt")

	/*resultFile, err := os.Create("result.txt")
	if err != nil {
		log.Fatal("Cannot create file", err)
	}
	defer resultFile.Close()*/

	p.NumStoredSamples = p.NumStoredSamplesCM
	p.NumGridSamples = p.GridStoredSamplesCM
	p.DetectionThreshold = p.DetectionThresholdCM
	//p.SquareRowCM = p.SquareRowCM
	//p.SquareColCM = p.SquareColCM
	p.MinDistance = 1000

	//Initializers
	cps.MakeBoolGrid(p)
	p.Server.Init()
	cps.ReadMap(p, r)
	if p.SuperNodes {
		p.Server.MakeSuperNodes()
	}
	//cps.GenerateRouting(p, r)

	//cps.FlipSquares(p, r)

	//Defaults to false

	if p.CommBomb {
		p.BombX = p.BombXCM
		p.BombY = p.BombYCM
		p.B = &cps.Bomb{X: p.BombX, Y: p.BombY}
	}

	if p.RandomBomb {
		reg := cps.RandomInt(0, len(r.Square_list))
		xval := cps.RandomInt(r.Square_list[reg].X1, r.Square_list[reg].X2+1)
		yval := cps.RandomInt(r.Square_list[reg].Y2, r.Square_list[reg].Y1+1)
		p.BombX = xval
		p.BombY = yval
		p.B = &cps.Bomb{X: p.BombX, Y: p.BombY}
	}
	fmt.Printf("Bomb location: %v, %v\n", p.BombX, p.BombY)

	p.NodeTree = &cps.Quadtree{
		Bounds: cps.Bounds{
			X:      0,
			Y:      0,
			Width:  float64(p.MaxX),
			Height: float64(p.MaxY),
		},
		MaxObjects: 1,
		MaxLevels:  15,
		Level:      0,
		Objects:    make([]*cps.NodeImpl, 0),
		ParentTree: nil,
		SubTrees:   make([]*cps.Quadtree, 0),
	}

	p.ClusterNetwork = &cps.AdHocNetwork{
		ClusterHeads: []*cps.NodeImpl{},
		FullReclusters: 0,
	}

	//This is where the text file reading ends
	Vn := make([]float64, 1000)
	for i := 0; i < 1000; i++ {
		Vn[i] = rand.NormFloat64() * -.5
	}

	if p.TotalNodes > p.CurrentNodes {
		for i := 0; i < p.TotalNodes-p.CurrentNodes; i++ {

			//p.NodeEntryTimes = append(p.NodeEntryTimes, []int{rangeInt(1, p.MaxX), rangeInt(1, p.MaxY), 0})
			p.NodeEntryTimes = append(p.NodeEntryTimes, []int{0, 0, 0})
		}
	}

	rand.Seed(time.Now().UnixNano()) //sets random to work properly by tying to to clock

	p.Iterations_used = 0
	p.Iterations_of_event = p.IterationsCM

	cps.SetupFiles(p)
	cps.SetupParameters(p, r)

	//Printing important information to the p.Grid log file
	//fmt.Fprintln(p.GridFile, "Grid:", p.SquareRowCM, "x", p.SquareColCM)
	//fmt.Fprintln(p.GridFile, "Total Number of Nodes:", (p.TotalNodes + numSuperNodes))
	//fmt.Fprintln(p.GridFile, "Runs:", iterations_of_event)

	fmt.Println("xDiv is ", p.XDiv, " yDiv is ", p.YDiv, " square capacity is ", p.SquareCapacity)

	p.WallNodeList = make([]cps.WallNodes, p.NumWallNodes)

	p.NodeList = make([]*cps.NodeImpl, 0)

	for i := 0; i < p.NumWallNodes; i++ {
		p.WallNodeList[i] = cps.WallNodes{Node: &cps.NodeImpl{X: float32(p.Wpos[i][0]), Y: float32(p.Wpos[i][1])}}
	}

	p.Server.MakeGrid()

	fmt.Println("Super Node Type", p.SuperNodeType)
	fmt.Println("Dimensions: ", p.MaxX, "x", p.MaxY)
	fmt.Printf("Running Simulator iteration %d\\%v", 0, p.Iterations_of_event)

	p.TimeStep = 0

	cps.WriteFlags(p)

	p.AliveNodes = make(map[*cps.NodeImpl]bool)
	p.Server.Clusters = make(map[*cps.NodeImpl]*cps.Cluster)
	p.Server.ClusterHeadsOf = make(map[*cps.NodeImpl][]*cps.NodeImpl)
	p.Server.AloneNodes = make(map[*cps.NodeImpl]int)

	if p.CSVMovement {
		cps.SetupCSVNodes(p)
	} else {
		cps.SetupRandomNodes(p)
	}
	p.Server.MakeNodeData()

	testMove:=false			//USES REGULAR MOVEMENT FOR ADAPTIVE SAMPLING

	//p.Events.Push(&cps.Event{&p.NodeList[0], "sense", 0, 0})
	p.Events.Push(&cps.Event{nil, cps.POSITION, 999, 0})
	if p.EnergyPrint {
		p.Events.Push(&cps.Event{nil, cps.ENERGYPRINT, 999, 0})
	}

	p.Events.Push(&cps.Event{nil, cps.SERVER, 999, 0})
	p.Events.Push(&cps.Event{nil, cps.GRID, 999, 0})
	p.Events.Push(&cps.Event{nil, cps.TIME, -1, 0})
	p.Events.Push(&cps.Event{nil, cps.GARBAGECOLLECT, 999, 0})
	p.Events.Push(&cps.Event{nil, cps.DRIFTLOG, 999, 0})
	p.Events.Push(&cps.Event{nil, cps.CLEANUPREADINGS, (p.ReadingHistorySize + 1) * 1000, 0})
	//p.Events.Push(&cps.Event{nil, cps.VALIDNODES, 999, 0})
	p.Events.Push(&cps.Event{nil, cps.LOADMOVE, (p.MovementSize - 2) * 1000, 0})
	p.Events.Push(&cps.Event{nil, cps.CLEANUPREADINGS, (p.ReadingHistorySize + 1) * 1000, 0})
	p.Events.Push(&cps.Event{nil, cps.SERVERSTATS, 1000, 0})
	if p.ClusteringOn {
		p.InitClusterTime *= 1000
		if p.InitClusterTime > 0 {
			p.Events.Push(&cps.Event{nil, cps.INITCLUSTER, p.InitClusterTime, 0})
		} else if p.GlobalRecluster == 2 {
			p.Events.Push(&cps.Event{nil, cps.FULLRECLUSTER, 5000, 0})
		}
		if p.ClusterPrint {
			fmt.Fprintf(p.ClusterFile, "Number of clusters, average cluster size, global reclusters, potential global reclusters, local reclusters, clusters above member threshold, clusters below member threshold, aliveValidNodes, Cluster Searches, CSJoins, CSSolos, Waits, Lost Readings\n")
			p.Events.Push(&cps.Event{nil, cps.CLUSTERPRINT, 999, 0})
		}
	}
	if p.BatteryPrint {
		fmt.Fprintf(p.BatteryFile, "percent alive, average battery level, min, max, samples, wifi, bluetooth, BTRecluster, BTClusterSearch, BTReadings\n")
		p.Events.Push(&cps.Event{nil, cps.BATTERYPRINT, 1000, 0})
	}
	p.CurrentTime = 0
	for len(p.Events) > 0 && p.CurrentTime < 1000*p.Iterations_of_event && !p.FoundBomb {
		event := heap.Pop(&p.Events).(*cps.Event)
		//fmt.Println(event)
		//fmt.Println(p.CurrentNodes)
		p.CurrentTime = event.Time
		switch event.Instruction {
		case cps.SENSE:
			if testMove {
				if p.CurrentTime/1000 < p.NumNodeMovements-5 {
					if p.CSVMovement {
						p.IsSense = true
						event.Node.MoveCSV(p)
						p.IsSense=false
					} else {
						event.Node.MoveNormal(p)
					}
				}
			} else {
				p.IsSense=false
			}
			event.Node.TimeLastSensed = p.CurrentTime
			event.Node.DrainBatterySample()
			event.Node.ScheduleNextSense()
			if p.DriftExplorer { //no sensor csv, just checking FP
				event.Node.GetSensor()
			} else {
				if p.CSVSensor { //if we have a big CSV file of the entire event
					event.Node.GetReadingsCSV()
				} else {
					event.Node.GetReadings() //if we have no big file, just the small 'FINE' csv file
				}
			}
		case cps.MOVE:
            p.IsSense = false
			if p.CSVMovement {
				event.Node.MoveCSV(p)
			} else {
				event.Node.MoveNormal(p)
			}
			if p.ClusteringOn && event.Node.Valid && event.Node.IsAlive() {
				p.NodeTree.NodeMovement(event.Node)
				event.Node.UpdateOutOfRange(p)
			}
			if p.CurrentTime/1000 < p.NumNodeMovements-5 {
				p.Events.Push(&cps.Event{event.Node, cps.MOVE, p.CurrentTime + 100, 0})
			}
		//case cps.INITCLUSTER:
		//	p.ClusterNetwork.FullRecluster(p)
		//	if p.GlobalRecluster > 0 {
		//		p.Events.Push(&cps.Event{nil, cps.FULLRECLUSTER, p.CurrentTime + p.ReclusterPeriod*1000, 0})
		//	}
		//case cps.FULLRECLUSTER:
		//	p.ClusterNetwork.FullRecluster(p)
		//	p.Events.Push(&cps.Event{nil, cps.FULLRECLUSTER, p.CurrentTime + p.ReclusterPeriod * 1000, 0})
		case cps.POSITION: //runs through all valid nodes and prints their location
			//fmt.Printf("Current Time: %v \n", p.CurrentTime)
			var avBuffer bytes.Buffer
			validCount := 0
			aliveCount := 0
			for i := 0; i < len(p.NodeList); i++ {
				if p.NodeList[i].Valid {
					validCount++
				}
			}
			avBuffer.WriteString(fmt.Sprintf("Valid Nodes:%v,Alive Nodes:%v\n", validCount, aliveCount))
			fmt.Fprintf(p.AliveValidNodes, avBuffer.String())

			if p.PositionPrint {
				amount := 0
				for i := 0; i < len(p.NodeList); i++ {
					//fmt.Printf("%v\n", p.NodeList[i].Valid)
					if p.NodeList[i].Valid {
						amount += 1
					}
				}
				fmt.Fprintln(p.PositionFile, "t= ", int(p.CurrentTime/1000), " amount= ", amount)
				var buffer bytes.Buffer
				for i := 0; i < len(p.NodeList); i++ {

					if p.NodeList[i].Valid {
						buffer.WriteString(fmt.Sprintf("ID: %v x: %v y: %v\n", p.NodeList[i].GetID(), int(p.NodeList[i].GetX()), int(p.NodeList[i].GetY())))
						//fmt.Fprintln(p.PositionFile, "ID:", p.NodeList[i].GetID(), "x:", int(p.NodeList[i].GetX()), "y:", int(p.NodeList[i].GetY()))
					}
				}
				fmt.Fprint(p.PositionFile, buffer.String())
			}
			fmt.Printf("\rRunning Simulator iteration %d\\%v", int(p.CurrentTime/1000), p.Iterations_of_event)
			p.Iterations_used += 1
			p.Events.Push(&cps.Event{nil, cps.POSITION, p.CurrentTime + 1000, 0})

		case cps.CLEANUPREADINGS:
			p.Server.CleanupReadings()
			p.Events.Push(&cps.Event{nil, cps.CLEANUPREADINGS, p.CurrentTime + 1000, 0})
		case cps.SERVERSTATS:
			//p.Server.CalcStats()
			p.Events.Push(&cps.Event{nil, cps.SERVERSTATS, p.CurrentTime + 1000, 0})

		case cps.SERVER: //scheduling for super nodes
			if !p.SuperNodes {
				fmt.Fprintln(p.RoutingFile, "Amount:", 0)
			} else {
				fmt.Fprintln(p.RoutingFile, "Amount:", p.NumSuperNodes)
			}
			p.Server.Tick()
			p.Events.Push(&cps.Event{nil, cps.SERVER, p.CurrentTime + 1000, 0})

		case cps.TIME:
			current := int(p.CurrentTime / 1000)
			for i := 0; i < len(p.SensorTimes); i++ {
				if current == p.SensorTimes[i] {
					p.TimeStep = i
					break
				}
			}
			if p.TimeStep+1 < len(p.SensorTimes) {
				p.Events.Push(&cps.Event{nil, cps.TIME, p.SensorTimes[p.TimeStep+1] * 1000, 0})
			}
			//fmt.Printf("\nSetting timestep to %v at %v next event at %v\n", p.SensorTimes[p.TimeStep], p.CurrentTime, p.SensorTimes[p.TimeStep+1]*1000)
		case cps.ENERGYPRINT:
			if p.EnergyPrint {
				fmt.Fprintln(p.EnergyFile,
					"Amount:", len(p.NodeList),
					"Samples:", p.Server.SamplesCounter,
					"Wifi:", p.Server.WifiCounter,
					"Bluetooth:", p.Server.BluetoothCounter,
					"Recluster:", p.Server.ReclusterBTCounter,
					"Cluster Search:", p.Server.ClusterSearchBTCounter,
					"Reading:", p.Server.ReadingBTCounter)
				var buffer bytes.Buffer
				//for i := 0; i < len(p.NodeList); i++ {
				//	//p.NodeList[i].BatteryOverTime[p.CurrentTime/1000] = p.NodeList[i].Battery
				//	buffer.WriteString(fmt.Sprintf("%v\n", p.NodeList[i]))
				//}
				fmt.Fprintf(p.EnergyFile, buffer.String())
			}
			p.Events.Push(&cps.Event{nil, cps.ENERGYPRINT, p.CurrentTime + 1000, 0})
		case cps.GRID:
			if p.GridPrint {
				//x := printGrid(p.Grid)
				printGrid(p, p.Grid)

				p.Events.Push(&cps.Event{nil, cps.GRID, p.CurrentTime + 1000, 0})
				fmt.Fprint(p.GridFile, "----------------\n")

			}
		case cps.GARBAGECOLLECT:
			runtime.GC()
			if p.MemProfile != "" {
				s := fmt.Sprintf("%s-%v", p.MemProfile, p.CurrentTime/1000)
				f, err := os.Create(s)
				if err != nil {
					log.Fatal("could not create memory profile: ", err)
				}
				defer f.Close()
				runtime.GC() // get up-to-date statistics
				if err := pprof.WriteHeapProfile(f); err != nil {
					log.Fatal("could not write memory profile: ", err)
				}
			}
			p.Events.Push(&cps.Event{nil, cps.GARBAGECOLLECT, p.CurrentTime + 100000, 0})
		case cps.DRIFTLOG:
			if p.DriftExplorer || !p.DriftExplorer {
				cps.DriftHist(p)
				p.Events.Push(&cps.Event{nil, cps.DRIFTLOG, p.CurrentTime + 1000, 0})
			}
		case cps.VALIDNODES:
			val := 0
			for _, n := range p.NodeList {
				if n.Valid {
					val += 1
				}
			}
			fmt.Printf("Valid: %v\n", val)
			p.Events.Push(&cps.Event{nil, cps.VALIDNODES, p.CurrentTime + 1000*100, 0})
		case cps.LOADMOVE:
			cps.PartialReadMovementCSV(p)
			p.Events.Push(&cps.Event{nil, cps.LOADMOVE, (p.MovementOffset + p.MovementSize - 2) * 1000, 0})
		case cps.CLUSTERPRINT:
			var clusterStatsBuffer bytes.Buffer
			var clusterDebugBuffer bytes.Buffer

			if p.ClusterDebug {

				clusterHeadCount := 0
				clusterMemberCount := 0
				clusterDebugBuffer.WriteString(fmt.Sprint(""))
				for node := range p.AliveNodes {
					if node.IsClusterHead {
						clusterHeadCount++
					} else if node.IsClusterMember {
						clusterMemberCount++
					}
				}
				clusterDebugBuffer.WriteString(fmt.Sprintf("Iteration: %v\tlen(p.ClusterNetwork.ClusterHeads): %v\tClusterHeads: %v\tClusterMembers: %v\n", p.CurrentTime/1000, len(p.ClusterNetwork.ClusterHeads), clusterHeadCount, clusterMemberCount))

				outOfRange := 0
				outOfRangeAndSensed := 0
				for node := range p.AliveNodes {
					if node.IsClusterHead {
						for mem := range node.ClusterMembers {
							if !(node.IsWithinRange(mem, p.NodeBTRange)) {
								if mem.TimeLastSensed > mem.TimeMovedOutOfRange {
									xDist := node.X - mem.X
									yDist := node.Y - mem.Y
									radDist := math.Sqrt(float64(xDist*xDist) + float64(yDist*yDist))
									clusterDebugBuffer.WriteString(fmt.Sprintf("\tCluster Member Out of Range: Member:{ID=%v, Coord(%v,%v)} Cluster:{CH_ID=%v, Coord(%v,%v),Size=%v} Dist: %.4f, Valid: %v\n",
										mem.Id, mem.X, mem.Y,
										node.Id, node.X, node.Y, len(node.ClusterMembers), radDist, node.Valid))
									outOfRangeAndSensed++
								}
								outOfRange++
							}
						}

						for j := 0; j < len(p.ClusterNetwork.ClusterHeads); j++ {
							for mem := range p.ClusterNetwork.ClusterHeads[j].ClusterMembers {
								if node == mem {
									clusterDebugBuffer.WriteString(fmt.Sprintf("\tCluster Head {CH_ID: %v, Size=%v} is cluster member of {CH_ID: %v, Size=%v}\n", node.Id, len(node.ClusterMembers), p.ClusterNetwork.ClusterHeads[j].Id, len(p.ClusterNetwork.ClusterHeads[j].ClusterMembers)))
								}
							}
						}
					} else if node.IsClusterMember {
						clusterCount := 0
						clusterHeadIsHead := node.ClusterHead.IsClusterHead
						clusterHeadInHeads := false
						for j := 0; j < len(p.ClusterNetwork.ClusterHeads); j++ {
							if p.ClusterNetwork.ClusterHeads[j] == node.ClusterHead {
								clusterHeadInHeads = true
							}
							for mem := range p.ClusterNetwork.ClusterHeads[j].ClusterMembers {
								if node == mem {
									clusterCount++
								}
							}
						}
						if clusterCount != 1 {
							clusterDebugBuffer.WriteString(fmt.Sprintf("\tNode ID=%v is cluster member of %v clusters. Cluster head in ClusterHeads: %v. Cluster head is head: %v\n", node.Id, clusterCount, clusterHeadInHeads, clusterHeadIsHead))
						}
					} else {
						if node.Valid {
							clusterDebugBuffer.WriteString(fmt.Sprintf("\tNode ID=%v is niether member nor head.\n", node.Id))
						}
					}
				}
				//clusterStatsBuffer.WriteString(fmt.Sprintf("Out of range and sensed: %v/%v\n", outOfRangeAndSensed, outOfRange))
				//clusterStatsBuffer.WriteString("\n")



				//fmt.Fprintf(p.ClusterFile, clusterStatsBuffer.String())
			}

			for i := 0; i < len(p.ClusterNetwork.ClusterHeads); i++ {
				//if len(p.ClusterNetwork.ClusterHeads[i].ClusterMembers) > 0 {
				//	clusterStatsBuffer.WriteString(fmt.Sprintf("%v: [", p.ClusterNetwork.ClusterHeads[i].Id))
				//	for j := 0; j < len(p.ClusterNetwork.ClusterHeads[i].ClusterMembers); j++ {
				//		clusterStatsBuffer.WriteString(fmt.Sprintf("%v", p.ClusterNetwork.ClusterHeads[i].ClusterMembers[j].Id))
				//		if j+1 != len(p.ClusterNetwork.ClusterHeads[i].ClusterMembers) {
				//			clusterStatsBuffer.WriteString(fmt.Sprintf(", "))
				//		}
				//	}
				//	clusterStatsBuffer.WriteString(fmt.Sprintf("]\n"))
				//}

				clusterStatsBuffer.WriteString(fmt.Sprintf("%v", len(p.ClusterNetwork.ClusterHeads[i].ClusterMembers)))
				if i+1 != len(p.ClusterNetwork.ClusterHeads) {
					clusterStatsBuffer.WriteString(fmt.Sprintf(","))
				}
			}
			clusterStatsBuffer.WriteString(fmt.Sprintln(""))

			fmt.Fprintf(p.ClusterStatsFile, clusterStatsBuffer.String())
			fmt.Fprintf(p.ClusterDebugFile, clusterDebugBuffer.String())

			clustersAboveThresh := len(p.ClusterNetwork.ClusterHeads)
			clustersBelowThresh := 0
			nodesBelowThresh := 0
			nodesInClusters := 0
			smallClustersNumWithinDist := 0
			for i := 0; i < len(p.ClusterNetwork.ClusterHeads); i++ {
				if len(p.ClusterNetwork.ClusterHeads[i].ClusterMembers) <= p.ClusterMinThreshold {
					if p.ReportBTAverages {
						smallClustersNumWithinDist += len(p.NodeTree.WithinRadius(p.NodeBTRange, p.ClusterNetwork.ClusterHeads[i], []*cps.NodeImpl{}))
					}
					nodesBelowThresh += len(p.ClusterNetwork.ClusterHeads[i].ClusterMembers) + 1
					clustersAboveThresh--
					clustersBelowThresh++
				}
				nodesInClusters += len(p.ClusterNetwork.ClusterHeads[i].ClusterMembers) + 1 //Plus one for the cluster head
			}
			if float64(nodesBelowThresh)/float64(len(p.AliveNodes)) > p.ReclusterThreshold {
				p.ClusterNetwork.PotentialReclusters++
			}

			average := 0
			if len(p.ClusterNetwork.ClusterHeads) > 0 {
				average = nodesInClusters / len(p.ClusterNetwork.ClusterHeads)
				p.ClusterNetwork.AverageClusterSize += average
			}

			aliveValidNodes := 0
			//countedCHs := 0
			//countedCMs := 0
			//unaccounted := 0
			//numWithinDist := 0
			for node := range p.AliveNodes {
				if node.Valid {
					//if node.IsClusterHead {
					//	countedCHs++
					//} else if node.IsClusterMember {
					//	countedCMs++
					//} else {
					//	unaccounted++
					//}
					aliveValidNodes++
					//if p.ReportBTAverages {
					//	numWithinDist += len(p.NodeTree.WithinRadius(p.NodeBTRange, node, []*cps.NodeImpl{}))
					//}
				}
			}
			//clusterStatsBuffer.WriteString(fmt.Sprintf("Iteration: %v\n", p.CurrentTime/1000))
			//clusterStatsBuffer.WriteString(fmt.Sprintf("Clusters above member threshold: %v\n", clustersAboveThresh))
			//clusterStatsBuffer.WriteString(fmt.Sprintf("Clusters below member threshold: %v\n", clustersBelowThresh))
			//clusterStatsBuffer.WriteString(fmt.Sprintf("Alive, valid nodes: %v\n", aliveValidNodes))
			//clusterStatsBuffer.WriteString(fmt.Sprintf("Total clustered nodes: %v\n", nodesInClusters))
			//clusterStatsBuffer.WriteString(fmt.Sprintf("Counted heads: %v\n", countedCHs))
			//clusterStatsBuffer.WriteString(fmt.Sprintf("Counted members: %v\n", countedCMs))
			//clusterStatsBuffer.WriteString(fmt.Sprintf("Unaccounted nodes: %v\n", unaccounted))
			//clusterStatsBuffer.WriteString(fmt.Sprintf("Average cluster size: %v\n", average))
			//if p.ReportBTAverages {
			//	if aliveValidNodes > 0 {
			//		clusterStatsBuffer.WriteString(fmt.Sprintf("Average number of nodes in BT range: %v\n", numWithinDist/aliveValidNodes))
			//	} else {
			//		clusterStatsBuffer.WriteString(fmt.Sprintf("Average number of nodes in BT range: 0\n"))
			//	}
			//	if clustersBelowThresh > 0 {
			//		clusterStatsBuffer.WriteString(fmt.Sprintf("Average number of nodes in BT range of small cluster heads: %v\n", smallClustersNumWithinDist/clustersBelowThresh))
			//	} else {
			//		clusterStatsBuffer.WriteString(fmt.Sprintf("Average number of nodes in BT range of small cluster heads: 0\n"))
			//	}
			//}

			//Number of clusters, average cluster size, global reclusters, potential global reclusters, local reclusters, clusters above member threshold, clusters below member threshold, aliveValidNodes, Cluster Searches, CSJoins, CSSolos, Waits, Lost Readings
			fmt.Fprintf(p.ClusterFile, "%v,%v,%v,%v,%v,%v,%v,%v,%v,%v,%v,%v,%v\n", len(p.ClusterNetwork.ClusterHeads), average, p.ClusterNetwork.FullReclusters, p.ClusterNetwork.PotentialReclusters, p.ClusterNetwork.LocalReclusters, clustersAboveThresh, clustersBelowThresh, aliveValidNodes, p.ClusterNetwork.CSJoins + p.ClusterNetwork.CSSolos, p.ClusterNetwork.CSJoins, p.ClusterNetwork.CSSolos, p.ClusterNetwork.TotalWaits, p.ClusterNetwork.LostReadings)

			p.ClusterNetwork.AverageNumClusters += len(p.ClusterNetwork.ClusterHeads)

			p.Events.Push(&cps.Event{nil, cps.CLUSTERPRINT, p.CurrentTime + 1000, 0})
		case cps.BATTERYPRINT:
			averageBattery := 0.0
			lowBattery := p.NodeList[0].GetBatteryPercentage()
			highBattery := lowBattery
			for i := 1; i < len(p.NodeList); i++ {
				averageBattery += p.NodeList[i].GetBatteryPercentage()
				if p.NodeList[i].GetBatteryPercentage() < lowBattery {
					lowBattery = p.NodeList[i].GetBatteryPercentage()
				} else if p.NodeList[i].GetBatteryPercentage() > highBattery {
					highBattery = p.NodeList[i].GetBatteryPercentage()
				}
			}
			averageBattery = averageBattery / float64(len(p.NodeList))
			//percent alive, average battery level, min, max, samples, wifi, bluetooth, BTRecluster, BTClusterSearch, BTReadings
			fmt.Fprintf(p.BatteryFile, "%v,%v,%v,%v,%v,%v,%v,%v,%v,%v\n", float64(len(p.AliveNodes))/float64(p.TotalNodes), averageBattery, lowBattery, highBattery, p.Server.SamplesCounter, p.Server.WifiCounter, p.Server.BluetoothCounter, p.Server.ReclusterBTCounter, p.Server.ClusterSearchBTCounter, p.Server.ReadingBTCounter)
			p.Events.Push(&cps.Event{nil, cps.BATTERYPRINT, p.CurrentTime + 1000, 0})
		}
	}

	if p.MemProfile != "" {
		f, err := os.Create(p.MemProfile)
		if err != nil {
			log.Fatal("could not create memory profile: ", err)
		}
		defer f.Close()
		runtime.GC() // get up-to-date statistics
		if err := pprof.WriteHeapProfile(f); err != nil {
			log.Fatal("could not write memory profile: ", err)
		}
	}

	if p.GridPrint {
		printGrid(p, p.Grid)

		amount := 0
		for i := 0; i < len(p.NodeList); i++ {
			//fmt.Printf("%v\n", p.NodeList[i].Valid)
			if p.NodeList[i].Valid {
				amount += 1
			}
		}
		fmt.Fprintln(p.PositionFile, "t= ", int(p.CurrentTime/1000), " amount= ", amount)
		var buffer bytes.Buffer
		for i := 0; i < len(p.NodeList); i++ {

			if p.NodeList[i].Valid {
				buffer.WriteString(fmt.Sprintf("ID: %v x: %v y: %v\n", p.NodeList[i].GetID(), int(p.NodeList[i].GetX()), int(p.NodeList[i].GetY())))
				//fmt.Fprintln(p.PositionFile, "ID:", p.NodeList[i].GetID(), "x:", int(p.NodeList[i].GetX()), "y:", int(p.NodeList[i].GetY()))
			}
		}
		fmt.Fprint(p.PositionFile, buffer.String())

		if !p.SuperNodes {
			fmt.Fprintln(p.RoutingFile, "Amount:", 0)
		} else {
			fmt.Fprintln(p.RoutingFile, "Amount:", p.NumSuperNodes)
		}
		fmt.Fprintln(p.EnergyFile, "Amount:", len(p.NodeList)) //big time waster
		//if p.EnergyPrint {
		//	var buffer bytes.Buffer
		//	for i := 0; i < p.CurrentNodes; i++ {
		//		buffer.WriteString(fmt.Sprintf("%v\n", p.NodeList[i]))
		//	}
		//	fmt.Fprintf(p.EnergyFile, buffer.String())
		//}
	}

	if p.ClusteringOn && p.ClusterPrint {
		fmt.Fprintln(p.ClusterStatsFile, "New nodes joining existing clusters:", p.ClusterNetwork.CSJoins)
		fmt.Fprintln(p.ClusterStatsFile, "New nodes creating a new cluster:", p.ClusterNetwork.CSSolos)
		fmt.Fprintln(p.ClusterStatsFile, "Average number of clusters:", p.ClusterNetwork.AverageNumClusters / (p.CurrentTime / 1000))
		fmt.Fprintln(p.ClusterStatsFile, "Overall average cluster size:", p.ClusterNetwork.AverageClusterSize / (p.CurrentTime / 1000))
		fmt.Fprintln(p.ClusterStatsFile, "Local Reclusters:", p.ClusterNetwork.LocalReclusters)
		fmt.Fprintln(p.ClusterStatsFile, "Cluster heads created by local reclusters:", p.ClusterNetwork.LRHeads)
		fmt.Fprintln(p.ClusterStatsFile, "Potential Full Reclusters:", p.ClusterNetwork.PotentialReclusters)
		fmt.Fprintln(p.ClusterStatsFile, "Actual Full Reclusters:", p.ClusterNetwork.FullReclusters)
		fmt.Fprintln(p.ClusterStatsFile, "Lost Readings:", p.ClusterNetwork.LostReadings, "/", p.Server.SamplesCounter)
	}

	p.PositionFile.Seek(0, 0)
	fmt.Fprintln(p.PositionFile, "Image:", p.ImageFileNameCM)
	fmt.Fprintln(p.PositionFile, "Width:", p.MaxX)
	fmt.Fprintln(p.PositionFile, "Height:", p.MaxY)
	fmt.Fprintf(p.PositionFile, "Amount: %-8v\n", int(p.CurrentTime/1000)+1)

	if p.Iterations_used < p.Iterations_of_event-1 {
		fmt.Printf("\nFound bomb at iteration: %v \nSimulation Complete\n", int(p.CurrentTime/1000))
	} else {
		fmt.Println("\nSimulation Complete")
	}

	p.Server.GoThroughSquares()
	if p.OutputPrint {
		for k, v := range p.Server.SquareTime { //k= key v=value
			fmt.Fprintln(p.OutputLog, "Square and max", k, v.MaxDelta)
		}
	}

	p.Server.PrintBatteryStats()
	fmt.Println("Total Samples",p.Server.SamplesCounter)
	fmt.Println("Total Adaptations",p.TotalAdaptations)

	if p.ZipFiles {
		p.MoveReadingsFile.Close()
		p.DriftFile.Close()
		p.NodeFile.Close()
		p.PositionFile.Close()
		p.GridFile.Close()
		p.EnergyFile.Close()
		p.RoutingFile.Close()
		p.ServerFile.Close()
		p.DetectionFile.Close()
		p.BatteryFile.Close()
		p.RunParamFile.Close()
		p.NodeDataFile.Close()
		p.DistanceFile.Close()

		if p.DriftExplorer {
			p.DriftExploreFile.Close()
		}

		output := p.OutputFileNameCM + ".zip"
		if err := cps.ZipFiles(output, p.Files); err != nil {
			panic(err)
		}
		fmt.Println("Zipped File:", output)

		for _, file := range p.Files {

			var err = os.Remove(file)
			if err != nil {
				fmt.Println(err)
			}
		}

	}

	//p.Server.PrintStatsFile()

}

//printGrid saves the current measurements of each Square into a buffer to print into the file
func printGrid(p *cps.Params, g [][]*cps.Square) {
	var buffer bytes.Buffer
	for y := 0; y < p.GridHeight; y++ {
		for x := 0; x < p.GridWidth; x++ {
			buffer.WriteString(fmt.Sprintf("%.2f\t", g[x][y].Avg))
		}
		buffer.WriteString("\n")
	}
	buffer.WriteString("\n")
	fmt.Fprintf(p.GridFile, buffer.String())
}

//printGridNodes saves the current p.NumNodes of each Square into a buffer to print to the file
func printGridNodes(g [][]*cps.Square) bytes.Buffer {
	var buffer bytes.Buffer
	for i, _ := range g {
		for _, x := range g[i] {
			buffer.WriteString(fmt.Sprintf("%d\t", x.NumNodes))
		}
		buffer.WriteString(fmt.Sprintf("\n"))
	}
	return buffer
}

//printSuperStats writes supernode data to a buffer
func printSuperStats(SNodeList []cps.SuperNodeParent) bytes.Buffer {
	var buffer bytes.Buffer
	for _, i := range SNodeList {
		buffer.WriteString(fmt.Sprintf("SuperNode: %d\t", i.GetId()))
		buffer.WriteString(fmt.Sprintf("SquaresMoved: %d\t", i.GetSquaresMoved()))
		buffer.WriteString(fmt.Sprintf("AvgResponseTime: %.2f\t", i.GetAvgResponseTime()))
	}
	return buffer
}