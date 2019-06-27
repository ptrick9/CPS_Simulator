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
		SubTrees:   make([]cps.Quadtree, 0),
	}

	// Insert some boxes
	qt.Insert(cps.Bounds{
		X:      1,
		Y:      1,
		Width:  10,
		Height: 10,
	})
	qt.Insert(cps.Bounds{
		X:      5,
		Y:      5,
		Width:  10,
		Height: 10,
	})
	qt.Insert(cps.Bounds{
		X:      10,
		Y:      10,
		Width:  10,
		Height: 10,
	})
	qt.Insert(cps.Bounds{
		X:      3,
		Y:      3,
		Width:  7,
		Height: 14,
	})
	qt.Insert(cps.Bounds{
		X:      15,
		Y:      15,
		Width:  10,
		Height: 10,
	})
	qt.Insert(cps.Bounds{
		X:      2,
		Y:      8,
		Width:  3,
		Height: 3,
	})


	////Get all the intersections
	//intersections := qt.RetrieveIntersections(cps.Bounds{
	//	X:      4,
	//	Y:      4,
	//	Width:  3,
	//	Height: 3,
	//})

	intersections := qt.RetrieveIntersections(cps.Bounds{10,12,10,10})

	fmt.Println(intersections)
	fmt.Println(qt.Objects)

	qt.PrintTree("")


	fmt.Println(qt.Retrieve(cps.Bounds{14, 14, 10, 10}))

	// Clear the Quadtree
	qt.Clear()
}
