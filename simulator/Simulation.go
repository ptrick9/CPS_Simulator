/*
-inputFileName=Scenario_3.txt
-imageFileName=marathon_street_map.png
-logPosition=true
-logGrid=true
-logEnergy=true
-logNodes=false
-noEnergy=true
-sensorPath=C:/Users/patrick/Dropbox/Patrick/udel/SUMMER2019/GitSimulator/smoothed_marathon.csv
-SquareRowCM=60
-SquareColCM=320
-csvMove=true
-movementPath=marathon_2k.txt
-iterations=1000
-csvSensor=true
-detectionThreshold=5
-superNodes=false
-detectionDistance=6
-cpuprofile=event
*/



package main

import (
	"./cps"
	"bytes"
	"container/heap"
	"math"
	"runtime"

	//"CPS_Simulator/simulator/cps"
	"fmt"
	"log"
	"math/rand"
	"os"
	"runtime/pprof"
	"time"

	"image"
	"image/png"
)

var (
	p *cps.Params
	r *cps.RegionParams

	//err error

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
	p.Server = cps.FusionCenter{p, r, nil, nil, nil, nil, nil, nil, nil, nil, nil}

	p.Tau1 = 3500
	p.Tau2 = 9000
	p.FoundBomb = false

	rand.Seed(time.Now().UTC().UnixNano())

	//getFlags()
	fmt.Fprintf(p.RunParamFile,"Starting file\n")
	cps.GetFlags(p)
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
	p.TotalNodes = cps.GetDashedInput("numNodes", p)
	//numStoredSamples = getDashedInput("numStoredSamples")
	p.MaxX = cps.GetDashedInput("maxX", p)
	p.MaxY = cps.GetDashedInput("maxY", p)
	//p.BombX = cps.GetDashedInput("bombX", p)
	//p.BombY = cps.GetDashedInput("bombY", p)

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
	p.PositionPrint = p.PositionPrintCM
	p.GridPrint = p.GridPrintCM
	p.EnergyPrint = p.EnergyPrintCM
	p.NodesPrint = p.NodesPrintCM
	p.SquareRowCM = p.SquareRowCM
	p.SquareColCM = p.SquareColCM

	//Initializers
	cps.MakeBoolGrid(p)
	p.Server.Init()
	cps.ReadMap(p, r)
	if (p.SuperNodes) {
		p.Server.MakeSuperNodes()
	}

	cps.GenerateRouting(p, r)

	cps.FlipSquares(p, r)

	p.NodeTree = &cps.Quadtree{
		Bounds: cps.Bounds{
			X:      0,
			Y:      0,
			Width:  float64(p.MaxX),
			Height: float64(p.MaxY),
		},
		MaxObjects: 4,
		MaxLevels:  4,
		Level:      0,
		Objects:    make([]*cps.Bounds, 0),
		ParentTree: nil,
		SubTrees:   make([]*cps.Quadtree, 0),
	}

	p.ClusterNetwork = &cps.AdHocNetwork{
		ClusterHeads:	[]*cps.NodeImpl{},
		TotalHeads:			0,
		Threshold:			p.ClusterThreshold,
		TotalMsgs:			0,
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
	p.ThreshHoldBatteryToHave = 30.0 //This is the threshold battery to have for all phones

	p.Iterations_used = 0
	p.Iterations_of_event = p.IterationsCM
	p.EstimatedPingsNeeded = 10200

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

	iters := 0
	p.TimeStep = 0

	cps.WriteFlags(p)

	if p.CSVMovement {
		cps.SetupCSVNodes(p)
	} else {
		cps.SetupRandomNodes(p)
	}

	//p.Events.Push(&cps.Event{&p.NodeList[0], "sense", 0, 0})
	p.Events.Push(&cps.Event{nil, cps.POSITION, 999, 0})
	p.Events.Push(&cps.Event{nil, cps.ENERGYPRINT, 999, 0})
	p.Events.Push(&cps.Event{nil, cps.SERVER, 999, 0})
	p.Events.Push(&cps.Event{nil, cps.GRID, 999, 0})
	p.Events.Push(&cps.Event{nil, cps.TIME, -1, 0})
	p.Events.Push(&cps.Event{nil, cps.GARBAGECOLLECT, 999, 0})
	p.Events.Push(&cps.Event{nil, cps.CLUSTERPRINT, 999, 0})
	p.Events.Push(&cps.Event{nil, cps.CLEANUPREADINGS, (p.ReadingHistorySize + 1) * 1000, 0})
	if(p.ClusteringOn){
		p.Events.Push(&cps.Event{nil,cps.CLUSTERLESSFORM,25,0})
	}

	p.CurrentTime = 0
	for len(p.Events) > 0 && p.CurrentTime < 1000*p.Iterations_of_event{
		//fmt.Println(p.CurFrentTime,1000*p.Iterations_of_event)
		event := heap.Pop(&p.Events).(*cps.Event)
		//fmt.Println(event)
		//fmt.Println(p.CurrentNodes)
		p.CurrentTime = event.Time
		if event.Node != nil{
			if event.Instruction == cps.SENSE {
				if event.Node.Battery > p.ThreshHoldBatteryToHave {

					if (p.CSVMovement) {
						event.Node.MoveCSV(p)
						if(event.Node.Valid){
							event.Node.DecrementPowerGPS()
						}
					} else {
						event.Node.MoveNormal(p)
						if(event.Node.Valid){
							event.Node.DecrementPowerGPS()
						}
					}

					if (p.CSVSensor) {
						event.Node.GetReadingsCSV()
					} else {
						event.Node.GetReadings()
					}
					event.Node.DecrementPowerSensor()
				}

			} else if event.Instruction == cps.MOVE {
				if event.Node.Battery > p.ThreshHoldBatteryToHave {
					if (p.CSVMovement) {
						event.Node.MoveCSV(p)
					} else {
						event.Node.MoveNormal(p)
					}

					if (event.Node.Valid) {
						p.ClusterNetwork.ClearClusterParams(event.Node)
						event.Node.DecrementPowerGPS()
					}


					//if(p.CurrentTime/1000 <= 100){
					p.Events.Push(&cps.Event{event.Node, cps.MOVE, p.CurrentTime + 100, 0})
				}
				//}

			}else if event.Instruction == cps.CLUSTERMSG {
				if event.Node.Battery > p.ThreshHoldBatteryToHave {
					if (event.Node.Valid) {
						p.ClusterNetwork.SendHelloMessage(p.NodeBTRange, event.Node)
					}
					p.Events.Push(&cps.Event{event.Node, cps.CLUSTERMSG, p.CurrentTime + 1000, 0})
				}

			} else if event.Instruction == cps.CLUSTERHEADELECT {
				if event.Node.Battery > p.ThreshHoldBatteryToHave {
					if (event.Node.Valid) {
						event.Node.SortMessages()
						p.ClusterNetwork.ElectClusterHead(event.Node)
					}
					p.Events.Push(&cps.Event{event.Node, cps.CLUSTERHEADELECT, p.CurrentTime + 1000, 0})
				}

			} else if event.Instruction == cps.CLUSTERFORM {
				if event.Node.Battery > p.ThreshHoldBatteryToHave {
					if (event.Node.Valid) {
						//p.ClusterNetwork.GenerateClusters(event.Node)
						if (event.Node.IsClusterHead) {
							p.ClusterNetwork.FormClusters(event.Node)
						}
					}
					p.Events.Push(&cps.Event{event.Node, cps.CLUSTERFORM, p.CurrentTime + 1000, 0})
				}
			} else if event.Instruction == cps.ScheduleSensor {
				event.Node.ScheduleSensing()
				p.Events.Push(&cps.Event{event.Node, cps.ScheduleSensor, p.CurrentTime + 50, 0})
			}
		} else {
			if event.Instruction == cps.POSITION {
				var avBuffer bytes.Buffer
				validCount := 0
				aliveCount := 0
				for i:=0; i<len(p.NodeList); i++{
					if(p.NodeList[i].Valid){
						validCount++
						if(p.NodeList[i].Battery>p.ThreshHoldBatteryToHave){
							aliveCount++
						}
					}
				}
				avBuffer.WriteString(fmt.Sprintf("Valid Nodes:%v,Alive Nodes:%v\n",validCount,aliveCount))
				fmt.Fprintf(p.AliveValidNodes, avBuffer.String())

				//fmt.Printf("Current Time: %v \n", p.CurrentTime)

				if p.PositionPrint {
					amount := 0
					for i := 0; i < p.CurrentNodes; i ++ {
						//fmt.Printf("%v\n", p.NodeList[i].Valid)
						if p.NodeList[i].Valid {
							amount += 1
						}
					}
					fmt.Fprintln(p.PositionFile, "t= ", int(p.CurrentTime/1000), " amount= ", amount)
					var buffer bytes.Buffer
					for i := 0; i < p.CurrentNodes; i ++ {

						if p.NodeList[i].Valid {
							buffer.WriteString(fmt.Sprintf("ID: %v x: %v y: %v\n", p.NodeList[i].GetID(), int(p.NodeList[i].GetX()), int(p.NodeList[i].GetY())))
							//fmt.Fprintln(p.PositionFile, "ID:", p.NodeList[i].GetID(), "x:", int(p.NodeList[i].GetX()), "y:", int(p.NodeList[i].GetY()))
						}
					}
					fmt.Fprint(p.PositionFile, buffer.String())
				}
				fmt.Printf("\rRunning Simulator iteration %d\\%v", int(p.CurrentTime/1000), p.Iterations_of_event)
				p.Events.Push(&cps.Event{nil, cps.POSITION, p.CurrentTime + 1000, 0})

			} else if event.Instruction == cps.CLEANUPREADINGS {
				p.Server.CleanupReadings()
				p.Events.Push(&cps.Event{nil,cps.CLEANUPREADINGS, p.CurrentTime + 1000, 0})
			}else if event.Instruction == cps.SERVER {
				if !p.SuperNodes {
					fmt.Fprintln(p.RoutingFile, "Amount:", 0)
				} else {
					fmt.Fprintln(p.RoutingFile, "Amount:", p.NumSuperNodes)
				}
				p.Server.Tick()
				p.Events.Push(&cps.Event{nil, cps.SERVER, p.CurrentTime + 1000, 0})

			} else if event.Instruction == cps.TIME {
				current := int(p.CurrentTime/1000)
				for i := 0; i < len(p.SensorTimes); i++ {
					if current == p.SensorTimes[i] {
						p.TimeStep = i
						break
					}
				}
				if (p.TimeStep+1 < len(p.SensorTimes)) {
					p.Events.Push(&cps.Event{nil, cps.TIME, p.SensorTimes[p.TimeStep+1]*1000, 0})
				}
				//fmt.Printf("\nSetting timestep to %v at %v next event at %v\n", p.SensorTimes[p.TimeStep], p.CurrentTime, p.SensorTimes[p.TimeStep+1]*1000)
			} else if event.Instruction == cps.ENERGYPRINT {
				fmt.Fprintln(p.EnergyFile, "Amount:", len(p.NodeList))  //big time waster
				if p.EnergyPrint {
					var buffer bytes.Buffer
					for i := 0; i < p.CurrentNodes; i ++ {
							p.NodeList[i].BatteryOverTime[p.CurrentTime/1000] = p.NodeList[i].Battery
							buffer.WriteString(fmt.Sprintf("%v\n", p.NodeList[i]))
					}
					fmt.Fprintf(p.EnergyFile, buffer.String())
				}
				p.Events.Push(&cps.Event{nil, cps.ENERGYPRINT, p.CurrentTime + 1000, 0})
			} else if event.Instruction == cps.GRID {
				if p.GridPrint {
					//x := printGrid(p.Grid)
					printGrid(p.Grid)

					//fmt.Fprintln(p.GridFile, x)
					p.Events.Push(&cps.Event{nil, cps.GRID, p.CurrentTime + 1000, 0})
					fmt.Fprint(p.GridFile, "----------------\n")
					//fmt.Println(p.Grid)

				}
			} else if event.Instruction == cps.GARBAGECOLLECT {
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
			} else if event.Instruction == cps.CLUSTERPRINT {
				var clusterBuffer bytes.Buffer
				var clusterStatsBuffer bytes.Buffer
				var clusterDebugBuffer bytes.Buffer

				totalHeads := p.ClusterNetwork.TotalHeads
				for i:=0; i<len(p.ClusterNetwork.ClusterHeads); i++{
					if(p.ClusterNetwork.ClusterHeads[i].NodeClusterParams.CurrentCluster.Total==0){
						totalHeads--
					}
				}
				clusterBuffer.WriteString(fmt.Sprintf("Amount: %v\n", totalHeads))
				for i:=0; i<len(p.ClusterNetwork.ClusterHeads); i++{
					if (p.ClusterNetwork.ClusterHeads[i].NodeClusterParams.CurrentCluster.Total>0){
						clusterBuffer.WriteString(fmt.Sprintf("%v: [", p.ClusterNetwork.ClusterHeads[i].Id))
						for j:=0; j<len(p.ClusterNetwork.ClusterHeads[i].NodeClusterParams.CurrentCluster.ClusterMembers); j++ {
							clusterBuffer.WriteString(fmt.Sprintf("%v", p.ClusterNetwork.ClusterHeads[i].NodeClusterParams.CurrentCluster.ClusterMembers[j].Id))
							if (j+1 != len( p.ClusterNetwork.ClusterHeads[i].NodeClusterParams.CurrentCluster.ClusterMembers)) {
								clusterBuffer.WriteString(fmt.Sprintf(", "))
							}
						}
						clusterBuffer.WriteString(fmt.Sprintf("]\n"))
					}

					clusterStatsBuffer.WriteString(fmt.Sprintf("%v", p.ClusterNetwork.ClusterHeads[i].NodeClusterParams.CurrentCluster.Total))
					if(i+1 != len(p.ClusterNetwork.ClusterHeads)){
						clusterStatsBuffer.WriteString(fmt.Sprintf(","))
					}
				}
				clusterStatsBuffer.WriteString(fmt.Sprintln(""))

				clusterHeadCount:=0
				clusterMemberCount:=0
				clusterDebugBuffer.WriteString(fmt.Sprint( "", ))
				for i:=0; i<len(p.NodeList); i++ {
					if(p.NodeList[i].IsClusterHead){
						clusterHeadCount++
					} else if(p.NodeList[i].IsClusterMember){
						clusterMemberCount++
					}
				}
				clusterDebugBuffer.WriteString(fmt.Sprintf("Iteration: %v\tlen(p.ClusterNetwork.ClusterHeads): %v\tClusterHeads: %v\tClusterMembers: %v\n",p.CurrentTime/1000,len(p.ClusterNetwork.ClusterHeads),clusterHeadCount,clusterMemberCount))

				for i:=0; i<len(p.NodeList); i++ {
					if(p.NodeList[i].IsClusterHead){
						for j:=0; j<len(p.NodeList[i].NodeClusterParams.CurrentCluster.ClusterMembers); j++{
							xDist := p.NodeList[i].X - p.NodeList[i].NodeClusterParams.CurrentCluster.ClusterMembers[j].X
							yDist := p.NodeList[i].Y - p.NodeList[i].NodeClusterParams.CurrentCluster.ClusterMembers[j].Y
							radDist := math.Sqrt(float64(xDist*xDist)+float64(yDist*yDist))
							if(!(p.NodeList[i].IsWithinRange(p.NodeList[i].NodeClusterParams.CurrentCluster.ClusterMembers[j],p.NodeBTRange))){
								clusterDebugBuffer.WriteString(fmt.Sprintf("\tCluster Member Out of Range: Member:{ID=%v, Coord(%v,%v)} Cluster:{CH_ID=%v, Coord(%v,%v),Size=%v} Dist: %.4f\n",
									p.NodeList[i].NodeClusterParams.CurrentCluster.ClusterMembers[j].Id,p.NodeList[i].NodeClusterParams.CurrentCluster.ClusterMembers[j].X,p.NodeList[i].NodeClusterParams.CurrentCluster.ClusterMembers[j].Y,
									p.NodeList[i].Id,p.NodeList[i].X,p.NodeList[i].Y,p.NodeList[i].NodeClusterParams.CurrentCluster.Total, radDist))
							}
						}

						for j:=0; j<len(p.ClusterNetwork.ClusterHeads); j++{
							for k:=0; k<len(p.ClusterNetwork.ClusterHeads[j].NodeClusterParams.CurrentCluster.ClusterMembers); k++{
								if(p.NodeList[i] == p.ClusterNetwork.ClusterHeads[j].NodeClusterParams.CurrentCluster.ClusterMembers[k]){
									clusterDebugBuffer.WriteString(fmt.Sprintf("\tCluster Head {CH_ID: %v, Size=%v} is cluster member of {CH_ID: %v, Size=%v}\n", p.NodeList[i].Id,p.NodeList[i].NodeClusterParams.CurrentCluster.Total,p.ClusterNetwork.ClusterHeads[j].Id,p.ClusterNetwork.ClusterHeads[j].NodeClusterParams.CurrentCluster.Total))
								}
							}
						}
					} else if(p.NodeList[i].IsClusterMember){
						clusterCount := 0
						for j:=0; j<len(p.ClusterNetwork.ClusterHeads); j++{
							for k:=0; k<len(p.ClusterNetwork.ClusterHeads[j].NodeClusterParams.CurrentCluster.ClusterMembers); k++{
								if(p.NodeList[i] == p.ClusterNetwork.ClusterHeads[j].NodeClusterParams.CurrentCluster.ClusterMembers[k]){
									clusterCount++
								}
							}
						}
						if(clusterCount>1){
							clusterDebugBuffer.WriteString(fmt.Sprintf("\tNode ID=%v is cluster member of %v clusters\n", p.NodeList[i].Id,clusterCount))
						}

					}
				}
				fmt.Fprintf(p.ClusterFile, clusterBuffer.String())
				fmt.Fprintf(p.ClusterStatsFile, clusterStatsBuffer.String())
				fmt.Fprintf(p.ClusterDebug, clusterDebugBuffer.String())

				fmt.Fprintf(p.ClusterMessages, "%d,%d\n",p.CurrentTime/1000,p.ClusterNetwork.TotalMsgs)

				p.Events.Push(&cps.Event{nil, cps.CLUSTERPRINT, p.CurrentTime + 1000, 0})
				p.ClusterNetwork.ResetClusters()
			} else if event.Instruction == cps.CLUSTERLESSFORM {
				p.ClusterNetwork.FinalizeClusters(p)

				//fmt.Println()
				p.Events.Push(&cps.Event{nil, cps.CLUSTERLESSFORM, p.CurrentTime + 1000, 0})
			}
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


	PrintNodeBatteryOverTime(p)

	p.PositionFile.Seek(0, 0)
	fmt.Fprintln(p.PositionFile, "Image:", p.ImageFileNameCM)
	fmt.Fprintln(p.PositionFile, "Width:", p.MaxX)
	fmt.Fprintln(p.PositionFile, "Height:", p.MaxY)
	fmt.Fprintf(p.PositionFile, "Amount: %-8v\n", int(p.CurrentTime/1000))

	if iters < p.Iterations_of_event-1 {
		fmt.Printf("\nFound bomb at iteration: %v \nSimulation Complete\n", int(p.CurrentTime/1000))
	} else {
		fmt.Println("\nSimulation Complete")
	}

	for i := range p.BoolGrid {
		fmt.Fprintln(p.BoolFile, p.BoolGrid[i])
	}


	p.Server.PrintStatsFile()

}

//printGrid saves the current measurements of each Square into a buffer to print into the file
func printGrid(g [][]*cps.Square) {
	var buffer bytes.Buffer
	for y := 0; y < len(g[0]); y++ {
		for x:=0; x < len(g); x++ {
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

func PrintNodeBatteryOverTime(p * cps.Params)  {

	//fmt.Fprint(p.BatteryFile, "Time,")
	//for i := range p.NodeList{
	//	n := p.NodeList[i]
	//	fmt.Fprintf(p.BatteryFile, "Node %d",n.Id)
	//	if(i<len(p.NodeList)){
	//		fmt.Fprint(p.BatteryFile, ",")
	//	}
	//}
	//fmt.Fprint(p.BatteryFile, "\n")

	for t:=0; t<p.Iterations_of_event; t++{
		fmt.Fprintf(p.BatteryFile, "%d,",t)
		for i := range p.NodeList{
			n := p.NodeList[i]
			fmt.Fprintf(p.BatteryFile, "%.4f",n.BatteryOverTime[t],)
			if(i<len(p.NodeList)){
				fmt.Fprint(p.BatteryFile, ",")
			}
		}
		fmt.Fprint(p.BatteryFile, "\n")
	}
	//p.BatteryFile.Sync()
}
func PrintNodeBatteryOverTimeFast(p * cps.Params)  {
	var buffer bytes.Buffer
	buffer.WriteString("Time,")
	for i := range p.NodeList{
		n := p.NodeList[i]
		buffer.WriteString(fmt.Sprintf("Node",n.GetID(),","))
	}
	buffer.WriteString("\n")

	for t:=0; t<p.Iterations_of_event; t++{
		buffer.WriteString(fmt.Sprintf("%d,", t))
		for i := range p.NodeList{
			n := p.NodeList[i]
			buffer.WriteString(fmt.Sprintf("%v,", n.BatteryOverTime[t]))
		}
		buffer.WriteString("\n")
	}
	fmt.Fprintf(p.BatteryFile, buffer.String())
}