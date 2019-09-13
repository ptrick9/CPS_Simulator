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
	//"CPS_Simulator/simulator/cps"
	"./cps"
	"bytes"
	"container/heap"

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
	p.Server = cps.FusionCenter{p, r, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil}

	p.Tau1 = 3500
	p.Tau2 = 9000
	p.FoundBomb = false

	rand.Seed(time.Now().UTC().UnixNano())

	//getFlags()
	fmt.Fprintf(p.RunParamFile,"Starting file\n")
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

	//Defaults to false

	if p.CommBomb {
		p.BombX = p.BombXCM
		p.BombY = p.BombYCM
		p.B = &cps.Bomb{X: p.BombX, Y: p.BombY}
	}

	if p.RandomBomb {
		reg := cps.RandomInt(0, len(r.Square_list))
		xval := cps.RandomInt(r.Square_list[reg].X1, r.Square_list[reg].X2 + 1)
		yval := cps.RandomInt(r.Square_list[reg].Y2, r.Square_list[reg].Y1 + 1)
		p.BombX = xval
		p.BombY = yval
		p.B = &cps.Bomb{X: p.BombX, Y: p.BombY}
	}
	fmt.Printf("Bomb location: %v, %v\n", p.BombX, p.BombY)

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
	p.Server.MakeNodeData()

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




	p.CurrentTime = 0
	for len(p.Events) > 0 && p.CurrentTime < 1000*p.Iterations_of_event && !p.FoundBomb{
		event := heap.Pop(&p.Events).(*cps.Event)
		//fmt.Println(event)
		//fmt.Println(p.CurrentNodes)
		p.CurrentTime = event.Time
		if event.Node != nil {
			if event.Instruction == cps.SENSE {

				if p.CurrentTime/1000 < p.NumNodeMovements-5 {
					if (p.CSVMovement) {
						event.Node.MoveCSV(p)
					} else {
						event.Node.MoveNormal(p)
					}
				}
				if (p.DriftExplorer) {
					event.Node.GetSensor()
				} else {
					if (p.CSVSensor) {
						event.Node.GetReadingsCSV()
					} else {
						event.Node.GetReadings()
					}
				}

			} else if event.Instruction == cps.MOVE {
				if(p.CSVMovement) {
					event.Node.MoveCSV(p)
				} else {
					event.Node.MoveNormal(p)
				}
				if p.CurrentTime/1000 < p.NumNodeMovements-5 {
					p.Events.Push(&cps.Event{event.Node, cps.MOVE, p.CurrentTime + 100, 0})
				}
			}
		} else {
			if event.Instruction == cps.POSITION {
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
				p.Iterations_used += 1
				p.Events.Push(&cps.Event{nil, cps.POSITION, p.CurrentTime + 1000, 0})

			} else if event.Instruction == cps.SERVER {
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
			} else if event.Instruction == cps.DRIFTLOG {
				if p.DriftExplorer || !p.DriftExplorer {
					cps.DriftHist(p)
					p.Events.Push(&cps.Event{nil, cps.DRIFTLOG, p.CurrentTime + 1000, 0})
				}
			}  else if event.Instruction == cps.CLEANUPREADINGS {
				p.Server.CleanupReadings()
				p.Events.Push(&cps.Event{nil,cps.CLEANUPREADINGS, p.CurrentTime + 1000, 0})
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


	PrintNodeBatteryOverTimeFast(p)

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

		if p.DriftExplorer {
			p.DriftExploreFile.Close()
		}

		output := p.OutputFileNameCM + ".zip"
		if err := cps.ZipFiles(output, p.Files); err != nil {
			panic(err)
		}
		fmt.Println("Zipped File:", output)

		for _, file := range(p.Files) {

			var err = os.Remove(file)
			if err != nil {
				fmt.Println(err)
			}
		}

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

	fmt.Fprint(p.BatteryFile, "Time,")
	for i := range p.NodeList{
		n := p.NodeList[i]
		fmt.Fprint(p.BatteryFile, "Node",n.GetID(),",")
	}
	fmt.Fprint(p.BatteryFile, "\n")

	for t:=0; t<p.Iterations_of_event; t++{
		fmt.Fprint(p.BatteryFile, t, ",")
		for i := range p.NodeList{
			n := p.NodeList[i]
			fmt.Fprint(p.BatteryFile, n.BatteryOverTime[t],",")
		}
		fmt.Fprint(p.BatteryFile, "\n")
	}
	p.BatteryFile.Sync()
}
func PrintNodeBatteryOverTimeFast(p * cps.Params)  {
	var buffer bytes.Buffer
	buffer.WriteString("Time,")
	for i := range p.NodeList{
		n := p.NodeList[i]
		buffer.WriteString(fmt.Sprintf("Node %v,",n.GetID()))
	}
	buffer.WriteString("\n")

	for t:=0; t<p.CurrentTime/1000; t++{
		buffer.WriteString(fmt.Sprintf("%v,", t))
		for i := range p.NodeList{
			n := p.NodeList[i]
			buffer.WriteString(fmt.Sprintf("%v,", n.BatteryOverTime[t]))
		}
		buffer.WriteString("\n")
	}
	fmt.Fprintf(p.BatteryFile, buffer.String())
}