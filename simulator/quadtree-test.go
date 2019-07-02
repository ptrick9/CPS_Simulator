package main

import (
	"../simulator/cps"
	"fmt"
	"math"
	"math/rand"
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

	squareDim := 8.0
	size := 10
	qt := cps.Quadtree{
		Bounds: cps.Bounds{
			X:      0,
			Y:      0,
			Width:  squareDim,
			Height: squareDim,
		},
		MaxObjects: 4,
		MaxLevels:  12,
		Level:      0,
		Objects:    make([]*cps.Bounds, 0),
		ParentTree: nil,
		SubTrees:   make([]*cps.Quadtree, 0),
	}

	nodes := make([]*cps.Bounds, size) //random 10,000 nodes
	for i:=0; i<size; i++{
		nodeX := rand.Float64() * squareDim
		nodeY := rand.Float64() * squareDim
		nodes[i] = &cps.Bounds{
			X:       nodeX,
			Y:       nodeY,
			Width:   0,
			Height:  0,
			CurTree: &qt,
		}
	}

	for i:=0; i<size; i++{
		qt.Insert(nodes[i])
	}

	searchRadius := 5.0
	iterativeResults := make([]int, size)
	treeResults := make([]int, size)

	for i:=0; i<size; i++{
		treeResults[i] = len(nodes[i].CurTree.WithinDistance(searchRadius, nodes[i], []cps.Bounds{}, true))
	}

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
				}
			}
		}
	}

	for i:=0; i<size; i++{
		fmt.Printf("%d %d %d ", i, iterativeResults[i], treeResults[i])
		if(iterativeResults[i] == treeResults[i]){
			fmt.Print("\n")
			continue
		} else{
			fmt.Print("Not equal")
			//break
		}
		fmt.Print("\n")
	}
	testInd := 3 //size-1

	withinDist := []cps.Bounds{}

	qt.PrintTree("")
	fmt.Println()
	fmt.Println(nodes[testInd])
	fmt.Println(nodes[testInd].CurTree)
	withinDist = nodes[testInd].CurTree.WithinDistance(searchRadius, nodes[testInd], withinDist, true)
	withinDistNode0 := len(withinDist)
	fmt.Println(withinDist)
	fmt.Println(withinDistNode0)
	fmt.Println()
	for i:=0; i<len(nodes); i++{
		fmt.Println(*(nodes[i]))
	}
}