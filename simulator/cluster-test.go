package main

import (
	"../simulator/cps"
	"fmt"
	"math"
	"math/rand"
	"time"
)

func main(){
	fmt.Println()

	rand.Seed(time.Now().UTC().UnixNano())

	w := 1600.0
	h := 300.0
	size := 5000
	transmissionRange := 33.0

	adhoc := &cps.AdHocNetwork{
		ClusterHeads:	[]*cps.NodeImpl{},
		TotalHeads:		0,
	}

	qt := cps.Quadtree{
		Bounds: cps.Bounds{
			X:      0,
			Y:      0,
			Width:  w,
			Height: h,
		},
		MaxObjects: 4,
		MaxLevels:  4,
		Level:      0,
		Objects:    make([]*cps.Bounds, 0),
		ParentTree: nil,
		SubTrees:   make([]*cps.Quadtree, 0),
	}

	nodeList := make([]*cps.NodeImpl,size)
	for i:=0; i<size;i++{
		cn := cps.NodeImpl{
			Id:	i,
			NodeBounds: &cps.Bounds{
				X:      rand.Float64()*w,
				Y:      rand.Float64()*h,
				Width:  0,
				Height: 0,
			},
			Battery: rand.Float32()*70+30,
			ClusterHead:			nil,
			NodeClusterParams: 	&cps.ClusterMemberParams{
				CurrentCluster: nil,
				RecvMsgs:	[]*cps.HelloMsg{},
				ThisNodeHello:	nil,
			},
			IsClusterHead:		false,
			IsClusterMember:	false,
			P: &cps.Params{
				NodeTree: &qt,
			},
		}
		cn.NodeBounds.CurNode = &cn

		nodeList[i] = &cn
		qt.Insert(cn.NodeBounds)
	}

	//qt.PrintTree("")

	for i:=0; i<len(nodeList); i++{
		nodeList[i].SendHelloMessage(transmissionRange)
	}

	//fmt.Println("Node locations and scores ")
	//for i:=0; i<len(nodeList); i++{
	//	fmt.Printf("Node %d: (%.4f, %.4f) %.4f  %p\n", i,nodeList[i].NodeBounds.X,nodeList[i].NodeBounds.Y,nodeList[i].NodeClusterParams.ThisNodeHello.NodeCHScore,nodeList[i].NodeBounds)
	//	//nodeList[i].PrintClusterNode()
	//}

	//fmt.Println()
	//fmt.Println("Node locations and # of received messages")
	//for i:=0; i<len(nodeList); i++{
	//	fmt.Printf("Node %d: (%.4f, %.4f) %d \n", i,nodeList[i].NodeBounds.X,nodeList[i].NodeBounds.Y,len(nodeList[i].NodeClusterParams.RecvMsgs))
	//	//nodeList[i].PrintClusterNode()
	//}

	fmt.Println()
	fmt.Println("Generating clusters...")
	for i:=0; i<len(nodeList); i++{
		adhoc.GenerateClusters(transmissionRange,nodeList[i])
	}
	fmt.Println("Done")

	//fmt.Println()
	//fmt.Println("Node#	isClusterHead")
	totalHeads := 0

	radDist := 0.0
	for i:=0; i<len(nodeList); i++{
		if(nodeList[i].IsClusterHead){
		//	fmt.Printf("Node %d: (%.4f, %.4f) %.4f\n", i,nodeList[i].NodeBounds.X,nodeList[i].NodeBounds.Y,nodeList[i].NodeClusterParams.ThisNodeHello.NodeCHScore)
		//	totalHeads++
		//	for j:=0; j<len(nodeList); j++{
		//		if(nodeList[j].IsClusterMember && nodeList[j].NodeClusterParams.CurrentCluster.ClusterHead == nodeList[i]){
		//			fmt.Printf("\tNode %d: (%.4f, %.4f) %.4f\n", j,nodeList[j].NodeBounds.X,nodeList[j].NodeBounds.Y,nodeList[j].NodeClusterParams.ThisNodeHello.NodeCHScore)
		//		}
		//	}
		}
	}

	totalHeads = adhoc.TotalHeads
	fmt.Printf("Amount: %d\n",adhoc.TotalHeads)
	for i:=0; i<len(adhoc.ClusterHeads); i++{
		fmt.Printf("%d: [", adhoc.ClusterHeads[i].Id)
		for j:=0; j<len(adhoc.ClusterHeads[i].NodeClusterParams.CurrentCluster.ClusterMembers); j++ {
			fmt.Printf("%d", adhoc.ClusterHeads[i].NodeClusterParams.CurrentCluster.ClusterMembers[j].Id)
			if (j+1 != len(adhoc.ClusterHeads[i].NodeClusterParams.CurrentCluster.ClusterMembers)) {
				fmt.Printf(", ")
			}
		}
		fmt.Printf("]\n")
	}


	dist := 0.0
	for i:=0; i<len(nodeList); i++{
		if(nodeList[i].IsClusterHead){
			for j:=0; j<len(nodeList); j++{
				if(nodeList[j].IsClusterHead && nodeList[j]!=nodeList[i]){
					difX := math.Abs(nodeList[i].NodeBounds.X - nodeList[j].NodeBounds.X)
					difY := math.Abs(nodeList[i].NodeBounds.Y - nodeList[j].NodeBounds.Y)
					dist += math.Sqrt(difX*difX + difY+difY)
				}
			}
			radDist += dist/float64((totalHeads-1))
		}
	}
	fmt.Println("Before Moving: ")
	fmt.Println("\tTotal Cluster Heads: ",totalHeads)
	fmt.Println("\tAverage Nodes per cluster: ",float32(size)/float32(totalHeads))
	fmt.Println("\tAverage Distance between clusterheads: ",radDist/float64(totalHeads))





	for i:=0; i<len(nodeList); i++{
		//Test Moving
		changeX := (w/8)- rand.Float64()*(w/4)
		changeY := (h/8) - rand.Float64()*(h/4)
		batteryDecr := rand.Float32()*10

		if(nodeList[i].NodeBounds.X+changeX<w && nodeList[i].NodeBounds.X+changeX>0){
			nodeList[i].NodeBounds.X = nodeList[i].NodeBounds.X + changeX
		}
		if(nodeList[i].NodeBounds.Y+changeY<h && nodeList[i].NodeBounds.Y+changeY>0) {
			nodeList[i].NodeBounds.Y = nodeList[i].NodeBounds.Y+changeY
		}
		qt.NodeMovement(nodeList[i].NodeBounds)
		nodeList[i].Battery = nodeList[i].Battery - batteryDecr
	}

	for i:=0; i<len(nodeList); i++{
		nodeList[i].SendHelloMessage(transmissionRange)
	}

	fmt.Println("Nodes changing position and losing battery")

//	qt.PrintTree("")

	fmt.Println("Reseting Cluster Params...")
	for i:=0; i<len(nodeList); i++ {
		adhoc.ClearClusterParams(nodeList[i])
	}

	fmt.Println()
	fmt.Println("Nodes sending hello...")
	for i:=0; i<len(nodeList); i++{
		nodeList[i].SendHelloMessage(transmissionRange)
	}

	//fmt.Println("Node locations and scores ")
	//for i:=0; i<len(nodeList); i++{
	//	fmt.Printf("Node %d: (%.4f, %.4f) %.4f  %p\n", i,nodeList[i].NodeBounds.X,nodeList[i].NodeBounds.Y,nodeList[i].NodeClusterParams.ThisNodeHello.NodeCHScore,nodeList[i].NodeBounds)
	//	//nodeList[i].PrintClusterNode()
	//}

	fmt.Println()
	fmt.Println("Generating clusters...")
	for i:=0; i<len(nodeList); i++{
		adhoc.GenerateClusters(transmissionRange,nodeList[i])
	}
	fmt.Println("Done")

	fmt.Println()
	fmt.Println("After moving:")
	totalHeads = 0
	radDist = 0.0
	for i:=0; i<len(nodeList); i++{
		if(nodeList[i].IsClusterHead){
		//	fmt.Printf("Node %d: (%.4f, %.4f) %.4f\n", i,nodeList[i].NodeBounds.X,nodeList[i].NodeBounds.Y,nodeList[i].NodeClusterParams.ThisNodeHello.NodeCHScore)
			totalHeads++
		//	for j:=0; j<len(nodeList); j++{
		//		if(nodeList[j].IsClusterMember && nodeList[j].NodeClusterParams.CurrentCluster.ClusterHead == nodeList[i]){
		//			fmt.Printf("\tNode %d: (%.4f, %.4f) %.4f\n", j,nodeList[j].NodeBounds.X,nodeList[j].NodeBounds.Y,nodeList[j].NodeClusterParams.ThisNodeHello.NodeCHScore)
		//		}
		//	}
		}
	}

	dist = 0.0
	for i:=0; i<len(nodeList); i++{
		if(nodeList[i].IsClusterHead){
			for j:=0; j<len(nodeList); j++{
				if(nodeList[j].IsClusterHead && nodeList[j]!=nodeList[i]){
					difX := math.Abs(nodeList[i].NodeBounds.X - nodeList[j].NodeBounds.X)
					difY := math.Abs(nodeList[i].NodeBounds.Y - nodeList[j].NodeBounds.Y)
					dist += math.Sqrt(difX*difX + difY+difY)
				}
			}
			radDist += dist/float64((totalHeads-1))
		}
	}


	fmt.Println("\tTotal Cluster Heads: ",totalHeads)
	fmt.Println("\tAverage Nodes per cluster: ",float32(size)/float32(totalHeads))
	fmt.Println("\tAverage Distance between clusterheads: ",radDist/float64(totalHeads))

	totalHeads = adhoc.TotalHeads
	fmt.Printf("Amount: %d\n",adhoc.TotalHeads)
	for i:=0; i<len(adhoc.ClusterHeads); i++{
		fmt.Printf("%d: [", adhoc.ClusterHeads[i].Id)
		for j:=0; j<len(adhoc.ClusterHeads[i].NodeClusterParams.CurrentCluster.ClusterMembers); j++ {
			fmt.Printf("%d", adhoc.ClusterHeads[i].NodeClusterParams.CurrentCluster.ClusterMembers[j].Id)
			if (j+1 != len(adhoc.ClusterHeads[i].NodeClusterParams.CurrentCluster.ClusterMembers)) {
				fmt.Printf(", ")
			}
		}
		fmt.Printf("]\n")
	}
}