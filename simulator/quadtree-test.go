package main

import(
	"../simulator/cps"
	"fmt"
)

func main(){
	qt := cps.Quadtree{
		Bounds: cps.Bounds{
			X:      0,
			Y:      0,
			Width:  100,
			Height: 100,
		},
		MaxObjects: 4,
		MaxLevels:  4,
		Level:      0,
		Objects:    make([]cps.Bounds, 0),
		ParentTree: nil,
		SubTrees:   make([]cps.Quadtree, 0),
	}

	// Insert some boxes
	qt.Insert(&cps.Bounds{
		X:      1,
		Y:      1,
		Width:  0,
		Height: 0,
		CurTree: &qt,
	})
	qt.Insert(&cps.Bounds{
		X:      5,
		Y:      5,
		Width:  0,
		Height: 0,
		CurTree: &qt,
	})
	qt.Insert(&cps.Bounds{
		X:      10,
		Y:      10,
		Width:  0,
		Height: 0,
		CurTree: &qt,
	})
	qt.Insert(&cps.Bounds{
		X:      3,
		Y:      3,
		Width:  0,
		Height: 0,
		CurTree: &qt,
	})
	qt.Insert(&cps.Bounds{
		X:      15,
		Y:      15,
		Width:  0,
		Height: 0,
		CurTree: &qt,
	})
	qt.Insert(&cps.Bounds{
		X:      2,
		Y:      8,
		Width:  0,
		Height: 0,
		CurTree: &qt,
	})
	fmt.Println("Print 1: ")
	qt.PrintTree("")
	bounds77 := &cps.Bounds{X: 7,Y: 7,Width: 0,Height: 0, CurTree: &qt}
	qt.Insert(bounds77)
	fmt.Println("Print 2: ")
	qt.PrintTree("")


	////Get all the intersections
	//intersections := qt.RetrieveIntersections(cps.Bounds{
	//	X:      4,
	//	Y:      4,
	//	Width:  3,
	//	Height: 3,
	//})

	intersections := qt.RetrieveIntersections(cps.Bounds{10,12,0,0, &qt})

	fmt.Println(intersections)
	fmt.Println(qt.Objects)






	fmt.Println(qt.Retrieve(cps.Bounds{14, 14, 10, 10, &qt}))

	fmt.Println(bounds77.CurTree.Bounds)
	if bounds77.CurTree.ParentTree != nil{
		fmt.Println(bounds77.CurTree.ParentTree.Bounds)
	}

	within5of77 := bounds77.WithinDistance(5.0, bounds77)
	fmt.Println("Within5 of 7,7:", within5of77)
	// Clear the Quadtree
	qt.Clear()
}