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
}