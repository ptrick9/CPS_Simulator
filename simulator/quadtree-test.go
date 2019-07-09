package main

import (
	"../simulator/cps"
	"fmt"
	"math"
	"math/rand"
	"time"
)

func main(){
	//qt := cps.Quadtree{
	//	Bounds: cps.Bounds{
	//		X:      0,
	//		Y:      0,
	//		Width:  100,
	//		Height: 100,
	//	},
	//	MaxObjects: 4,
	//	MaxLevels:  4,
	//	Level:      0,
	//	Objects:    make([]cps.Bounds, 0),
	//	ParentTree: nil,
	//	SubTrees:   make([]*cps.Quadtree, 0),
	//}
	//
	//// Insert some boxes
	////qt.Insert(&cps.Bounds{
	////	X:      1,
	////	Y:      1,
	////	Width:  0,
	////	Height: 0,
	////	CurTree: &qt,
	////})
	//qt.Insert(&cps.Bounds{
	//	X:      5,
	//	Y:      5,
	//	Width:  0,
	//	Height: 0,
	//	CurTree: &qt,
	//})
	//qt.Insert(&cps.Bounds{
	//	X:      10,
	//	Y:      10,
	//	Width:  0,
	//	Height: 0,
	//	CurTree: &qt,
	//})
	//qt.Insert(&cps.Bounds{
	//	X:      3,
	//	Y:      3,
	//	Width:  0,
	//	Height: 0,
	//	CurTree: &qt,
	//})
	//qt.Insert(&cps.Bounds{
	//	X:      15,
	//	Y:      15,
	//	Width:  0,
	//	Height: 0,
	//	CurTree: &qt,
	//})
	//qt.Insert(&cps.Bounds{
	//	X:      2,
	//	Y:      8,
	//	Width:  0,
	//	Height: 0,
	//	CurTree: &qt,
	//})
	//fmt.Println("Print 1: ")
	//qt.PrintTree("")
	//bounds77 := &cps.Bounds{X: 7,Y: 7,Width: 0,Height: 0, CurTree: &qt}
	//qt.Insert(bounds77)
	//fmt.Println("Print 2: ")
	//qt.PrintTree("")
	//
	//intersections := qt.RetrieveIntersections(cps.Bounds{10,12,0,0, &qt})
	//
	//fmt.Println(intersections)
	//fmt.Println(qt.Objects)
	//
	//fmt.Println(qt.Retrieve(cps.Bounds{14, 14, 10, 10, &qt}))
	//
	//fmt.Println(bounds77.CurTree.Bounds)
	//if bounds77.CurTree.ParentTree != nil{
	//	fmt.Println(bounds77.CurTree.ParentTree.Bounds)
	//}
	////
	//withinDist := []cps.Bounds{}
	////within5of77 := bounds77.WithinDistance(5.0, bounds77, withinDist, true)
	////fmt.Println("Within5 of 7,7:", within5of77)
	//
	//qt.Remove(bounds77)
	//fmt.Println("After removing (7,7)")
	//qt.PrintTree("")
	//
	//// Clear the Quadtree
	//qt.Clear()

	treeDimX := 300.0
	treeDimY := 300.0
	size := 5000
	qt := cps.Quadtree{
		Bounds: cps.Bounds{
			X:      0,
			Y:      0,
			Width:  treeDimX,
			Height: treeDimY,
		},
		MaxObjects: 4,
		MaxLevels:  4,
		Level:      0,
		Objects:    make([]*cps.Bounds, 0),
		ParentTree: nil,
		SubTrees:   make([]*cps.Quadtree, 0),
	}

	nodes := make([]*cps.Bounds, size) //random 10,000 nodes
	for i:=0; i<size; i++{
		nodeX := rand.Float64() * treeDimX
		nodeY := rand.Float64() * treeDimY
		nodes[i] = &cps.Bounds{
			X:       nodeX,
			Y:       nodeY,
			Width:   0.0,
			Height:  0.0,
			CurTree: &qt,
		}
	}

	for i:=0; i<size; i++{
		qt.Insert(nodes[i])
	}

	searchRadius := 5.0
	iterativeResults := make([]int, size)
	treeResults := make([]int, size)

	treeStartTime := time.Now()
	for i:=0; i<size; i++{
		//treeResults[i] = len(nodes[i].CurTree.WithinDistance(searchRadius, nodes[i], []cps.Bounds{}, true))
		//treeResults[i] = len(qt.WithinDistance2(searchRadius, nodes[i], []cps.Bounds{}, true))
		treeResults[i] = len(qt.WithinRadius(searchRadius, nodes[i], nodes[i].GetSearchBounds(searchRadius), []*cps.Bounds{}))
	}
	treeEndTime := time.Since(treeStartTime)
	fmt.Print("Tree runtime: ")
	fmt.Print(treeEndTime)
	fmt.Println()
	testInd := 0 //size-1
	iterativeStartTime := time.Now()
	for i:=0; i<size; i++{
		searchingNode := nodes[i]
		for j:=0; j<size; j++{
			compareNode := nodes[j]
			if(searchingNode == compareNode){
				continue
			} else{
				difX := searchingNode.X - compareNode.X
				difY := searchingNode.Y - compareNode.Y
				radDist := math.Sqrt(difX*difX + difY*difY)
				if(radDist <= searchRadius){
					iterativeResults[i] = iterativeResults[i] + 1
					if(i==testInd){
						//fmt.Println(searchingNode, searchingNode.CurTree.ParentTree.Bounds)
						//fmt.Println(compareNode,compareNode.CurTree.Bounds,compareNode.CurTree.ParentTree.Bounds)
					}
				}
			}
		}
	}
	iterativeEndTime := time.Since(iterativeStartTime)
	fmt.Print("Iteration runtime: ")
	fmt.Print(iterativeEndTime)
	fmt.Println()

	treeTotal := 0.0
	iterativeTotal := 0.0

	resultsMatch := true
	i:=0;
	//wrong:=0
	for ; i<size; i++{
		treeTotal = treeTotal+float64(treeResults[i])
		iterativeTotal = iterativeTotal + float64(iterativeResults[i])
		if(iterativeResults[i] == treeResults[i]){
			//fmt.Print("\n")
			continue
		} else{
			//fmt.Printf("%d %d %d \n", i, iterativeResults[i], treeResults[i])
			resultsMatch = false
			//wrong++
			break
		}
	}
	fmt.Println("Done checking ")
	if(resultsMatch){
		fmt.Println("Results Match")
	} else{
		fmt.Print("Results Do NOT Match: ")
		fmt.Printf("%d %d %d \n", i, iterativeResults[i], treeResults[i])
		//fmt.Printf("Total Wrong: %d\t SuccessRate: %f\n",wrong,float64((size-wrong))/float64(size))
	}

	treeAvg := treeTotal/float64(size)
	iterativeAvg := iterativeTotal/float64(size)
	fmt.Printf("Average nodes within %f of a given node\n", searchRadius)
	fmt.Printf("Tree: %f\t Iterative: %f\n", treeAvg, iterativeAvg)

	//
	//withinDist := []*cps.Bounds{}
	//
	////qt.PrintTree("")
	//fmt.Println()
	//fmt.Println(nodes[testInd],nodes[testInd].CurTree.Bounds,nodes[testInd].CurTree.ParentTree.Bounds)
	////withinDist = qt.WithinDistance2(searchRadius, nodes[testInd], withinDist, true)
	//withinDist = qt.WithinRadius(searchRadius, nodes[testInd], nodes[testInd].GetSearchBounds(searchRadius), withinDist)
	//withinDistNode0 := len(withinDist)
	////fmt.Println(withinDist)
	//for i:=0; i<len(withinDist); i++{
	//	fmt.Print(*withinDist[i])
	//	fmt.Print("\t")
	//}
	//fmt.Println()
	//fmt.Println(withinDistNode0)
	//fmt.Println()
	//for i:=0; i<len(nodes); i++{
	//	fmt.Println(*(nodes[i]))
	//}
	fmt.Println()
	fmt.Println("Beginning remove/inserting test...")
	qt.Clear()

	qt = cps.Quadtree{
		Bounds: cps.Bounds{
			X:      0,
			Y:      0,
			Width:  16,
			Height: 16,
		},
		MaxObjects: 4,
		MaxLevels:  4,
		Level:      0,
		Objects:    make([]*cps.Bounds, 0),
		ParentTree: nil,
		SubTrees:   make([]*cps.Quadtree, 0),
	}

	nodeXValues := [7]float64{1,9,10,12.5,13,13,13}
	nodeYValues := [7]float64{1,3,2,6.5,1,1.5,3}

	for i:=0; i<7; i++{
		b := cps.Bounds{
			X:	nodeXValues[i],
			Y:	nodeYValues[i],
			Width:	0,
			Height: 0,
			CurTree: &qt,
		}
		qt.Insert(&b)
	}

	node1 := cps.Bounds{
		X:	1,
		Y:	9,
		Width:	0,
		Height: 0,
		CurTree: &qt,
	}

	node2 := cps.Bounds{
		X:	7,
		Y:	7,
		Width:	0,
		Height: 0,
		CurTree: &qt,
	}

	node3 := cps.Bounds{
		X:	15.5,
		Y:	1.5,
		Width:	0,
		Height: 0,
		CurTree: &qt,
	}

	node4 := cps.Bounds{
		X:	15,
		Y:	3,
		Width:	0,
		Height: 0,
		CurTree: &qt,
	}

	node5 := cps.Bounds{
		X:	9,
		Y:	9,
		Width:	0,
		Height: 0,
		CurTree: &qt,
	}

	node6 := cps.Bounds{
		X:	15,
		Y:	9,
		Width:	0,
		Height: 0,
		CurTree: &qt,
	}

	node7 := cps.Bounds{
		X:	15,
		Y:	15,
		Width:	0,
		Height: 0,
		CurTree: &qt,
	}

	node8 := cps.Bounds{
		X:	9,
		Y:	15,
		Width:	0,
		Height: 0,
		CurTree: &qt,
	}

	node9 := cps.Bounds{
		X:	9,
		Y:	15.1,
		Width:	0,
		Height: 0,
		CurTree: &qt,
	}

	qt.Insert(&node1) //(1,9) to (1,7)
	qt.Insert(&node2) //7,7 to 7,9
	qt.Insert(&node3) //15.5,1.5 to 15.5,3
	qt.Insert(&node4) //15,3 to 15,1.5
	qt.Insert(&node5)
	qt.Insert(&node6)
	qt.Insert(&node7)
	qt.Insert(&node8)
	qt.Insert(&node9)

	fmt.Println()
	fmt.Println("Nodes of interest:")
	fmt.Print("Node 1: ",node1)
	fmt.Printf("\t%p\n",&node1)
	fmt.Print("Node 2: ",node2)
	fmt.Printf("\t%p\n",&node2)
	fmt.Print("Node 3: ",node3)
	fmt.Printf("\t%p\n",&node3)
	fmt.Print("Node 4: ",node4)
	fmt.Printf("\t%p\n",&node4)
	fmt.Print("Node 5: ",node5)
	fmt.Printf("\t%p\n",&node5)
	fmt.Print("Node 6: ",node6)
	fmt.Printf("\t%p\n",&node6)
	fmt.Print("Node 7: ",node7)
	fmt.Printf("\t%p\n",&node7)
	fmt.Print("Node 8: ",node8)
	fmt.Printf("\t%p\n",&node8)
	fmt.Print("Node 9: ",node9)
	fmt.Printf("\t%p\n",&node9)

	fmt.Println()
	fmt.Println("Tree before removing:")
	qt.PrintTree("")

	node1.Y = 7.0
	//node2.Y = 9.0
	node3.X = 3.0
	node4.Y = 1.5

	node5.X = node5.X-8
	node6.X = node6.X-8
	node7.X = node7.X-8
	node8.X = node8.X-8
	node9.X = node9.X-8

	qt.Remove(&node1)
	//qt.Remove(&node2)
	qt.Remove(&node3)
	qt.Remove(&node4)
	qt.Remove(&node5)
	qt.Remove(&node6)
	qt.Remove(&node7)
	qt.Remove(&node8)
	qt.Remove(&node9)

	fmt.Println()
	fmt.Println("Tree after simple-remove:")
	qt.PrintTree("")

	qt.Insert(&node1)
	//qt.Insert(&node2)
	qt.Insert(&node3)
	qt.Insert(&node4)
	qt.Insert(&node5)
	qt.Insert(&node6)
	qt.Insert(&node7)
	qt.Insert(&node8)
	qt.Insert(&node9)

	fmt.Println()
	fmt.Println("Tree after moving nodes, before cleanup: ")
	qt.PrintTree("")

	qt.CleanUp()
	fmt.Println()
	fmt.Println("Tree after cleanup: ")
	qt.PrintTree("")

}