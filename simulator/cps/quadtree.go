package cps

import (
	"fmt"
	"math"
)

// Quadtree - The quadtree data structure
//based off https://github.com/JamesMilnerUK/quadtree-go/blob/master/quadtree.go

type Quadtree struct {
	Bounds     Bounds
	MaxObjects int // Maximum objects a node can hold before splitting into 4 SubTrees
	MaxLevels  int // Total max levels inside root Quadtree
	Level      int // Depth level, required for SubTrees
	Objects    []Bounds
	ParentTree *Quadtree
	SubTrees   []Quadtree
	Total      int
}

// Bounds - A bounding box with a x,y origin and width and height
type Bounds struct {
	X      	float64
	Y      	float64
	Width  	float64
	Height 	float64
	CurTree	*Quadtree //tree that this object is in
	//curNode *NodeImpl
}

//IsPoint - Checks if a bounds object is a point or not (has no width or height)
func (b *Bounds) IsPoint() bool {

	if b.Width == 0 && b.Height == 0 {
		return true
	}

	return false

}

// Intersects - Checks if a Bounds object intersects with another Bounds
func (b *Bounds) Intersects(a Bounds) bool {

	aMaxX := a.X + a.Width
	aMaxY := a.Y + a.Height
	bMaxX := b.X + b.Width
	bMaxY := b.Y + b.Height

	// a is left of b
	if aMaxX < b.X {
		return false
	}

	// a is right of b
	if a.X > bMaxX {
		return false
	}

	// a is above b
	if aMaxY < b.Y {
		return false
	}

	// a is below b
	if a.Y > bMaxY {
		return false
	}

	// The two overlap
	return true

}

// TotalSubTrees - Retrieve the total number of sub-Quadtrees in a Quadtree
func (qt *Quadtree) TotalSubTrees() int {

	total := 0

	if len(qt.SubTrees) > 0 {
		for i := 0; i < len(qt.SubTrees); i++ {
			total += 1
			total += qt.SubTrees[i].TotalSubTrees()
		}
	}

	return total

}

// split - split the node into 4 SubTrees
func (qt *Quadtree) split() {

	if len(qt.SubTrees) == 4 {
		return
	}

	nextLevel := qt.Level + 1
	subWidth := qt.Bounds.Width / 2
	subHeight := qt.Bounds.Height / 2
	x := qt.Bounds.X
	y := qt.Bounds.Y

	//top right node (0)
	b0 := Bounds{
		X:      x + subWidth,
		Y:      y,
		Width:  subWidth,
		Height: subHeight,
	}
	sub0 :=  Quadtree{
		Bounds: b0,
		MaxObjects: qt.MaxObjects,
		MaxLevels:  qt.MaxLevels,
		Level:      nextLevel,
		Objects:    make([]Bounds, 0),
		ParentTree: qt,
		SubTrees:   make([]Quadtree, 0, 4),
	}
	sub0.Bounds.CurTree = &sub0
	qt.SubTrees = append(qt.SubTrees, sub0)

	//qt.SubTrees = append(qt.SubTrees, Quadtree{
	//	Bounds: Bounds{
	//		X:      x + subWidth,
	//		Y:      y,
	//		Width:  subWidth,
	//		Height: subHeight,
	//	},
	//	MaxObjects: qt.MaxObjects,
	//	MaxLevels:  qt.MaxLevels,
	//	Level:      nextLevel,
	//	Objects:    make([]Bounds, 0),
	//	ParentTree: qt,
	//	SubTrees:   make([]Quadtree, 0, 4),
	//})

	//top left node (1)
	b1 := Bounds{
		X:      x,
		Y:      y,
		Width:  subWidth,
		Height: subHeight,
	}
	sub1 :=  Quadtree{
		Bounds: b1,
		MaxObjects: qt.MaxObjects,
		MaxLevels:  qt.MaxLevels,
		Level:      nextLevel,
		Objects:    make([]Bounds, 0),
		ParentTree: qt,
		SubTrees:   make([]Quadtree, 0, 4),
	}
	sub1.Bounds.CurTree = &sub1
	qt.SubTrees = append(qt.SubTrees, sub1)

	//qt.SubTrees = append(qt.SubTrees, Quadtree{
	//	Bounds: Bounds{
	//		X:      x,
	//		Y:      y,
	//		Width:  subWidth,
	//		Height: subHeight,
	//	},
	//	MaxObjects: qt.MaxObjects,
	//	MaxLevels:  qt.MaxLevels,
	//	Level:      nextLevel,
	//	Objects:    make([]Bounds, 0),
	//	ParentTree: qt,
	//	SubTrees:   make([]Quadtree, 0, 4),
	//})

	//bottom left node (2)
	b2 := Bounds{
		X:      x,
		Y:      y + subHeight,
		Width:  subWidth,
		Height: subHeight,
	}
	sub2 :=  Quadtree{
		Bounds: b2,
		MaxObjects: qt.MaxObjects,
		MaxLevels:  qt.MaxLevels,
		Level:      nextLevel,
		Objects:    make([]Bounds, 0),
		ParentTree: qt,
		SubTrees:   make([]Quadtree, 0, 4),
	}
	sub2.Bounds.CurTree = &sub2
	qt.SubTrees = append(qt.SubTrees, sub2)

	//qt.SubTrees = append(qt.SubTrees, Quadtree{
	//	Bounds: Bounds{
	//		X:      x,
	//		Y:      y + subHeight,
	//		Width:  subWidth,
	//		Height: subHeight,
	//	},
	//	MaxObjects: qt.MaxObjects,
	//	MaxLevels:  qt.MaxLevels,
	//	Level:      nextLevel,
	//	Objects:    make([]Bounds, 0),
	//	ParentTree: qt,
	//	SubTrees:   make([]Quadtree, 0, 4),
	//})

	//bottom right node (3)
	b3 := Bounds{
		X:      x + subWidth,
		Y:      y + subHeight,
		Width:  subWidth,
		Height: subHeight,
	}
	sub3 :=  Quadtree{
		Bounds: b3,
		MaxObjects: qt.MaxObjects,
		MaxLevels:  qt.MaxLevels,
		Level:      nextLevel,
		Objects:    make([]Bounds, 0),
		ParentTree: qt,
		SubTrees:   make([]Quadtree, 0, 4),
	}
	sub3.Bounds.CurTree = &sub3
	qt.SubTrees = append(qt.SubTrees, sub3)


	//qt.SubTrees = append(qt.SubTrees, Quadtree{
	//	Bounds: Bounds{
	//		X:      x + subWidth,
	//		Y:      y + subHeight,
	//		Width:  subWidth,
	//		Height: subHeight,
	//	},
	//	MaxObjects: qt.MaxObjects,
	//	MaxLevels:  qt.MaxLevels,
	//	Level:      nextLevel,
	//	Objects:    make([]Bounds, 0),
	//	ParentTree: qt,
	//	SubTrees:   make([]Quadtree, 0, 4),
	//})

}

// getIndex - Determine which quadrant the object belongs to (0-3)
func (qt *Quadtree) getIndex(pRect Bounds) int {

	index := -1 // index of the subnode (0-3), or -1 if pRect cannot completely fit within a subnode and is part of the parent node

	verticalMidpoint := qt.Bounds.X + (qt.Bounds.Width / 2)
	horizontalMidpoint := qt.Bounds.Y + (qt.Bounds.Height / 2)

	//pRect can completely fit within the top quadrants
	topQuadrant := (pRect.Y < horizontalMidpoint) && (pRect.Y+pRect.Height < horizontalMidpoint)

	//pRect can completely fit within the bottom quadrants
	bottomQuadrant := (pRect.Y > horizontalMidpoint)

	//pRect can completely fit within the left quadrants
	if (pRect.X < verticalMidpoint) && (pRect.X+pRect.Width < verticalMidpoint) {

		if topQuadrant {
			index = 1
		} else if bottomQuadrant {
			index = 2
		}

	} else if pRect.X > verticalMidpoint {
		//pRect can completely fit within the right quadrants

		if topQuadrant {
			index = 0
		} else if bottomQuadrant {
			index = 3
		}

	}

	return index

}

// Insert - Insert the object into the tree. If the tree exceeds the capacity,
// it will split and add all objects to their corresponding SubTrees.
func (qt *Quadtree) Insert(pRect * Bounds) {
	pRect.CurTree = qt
	qt.Total++

	i := 0
	var index int

	// If we have SubTrees within the Quadtree
	if (len(qt.SubTrees) > 0) {

		index = qt.getIndex(*pRect)

		if index != -1 {
			pRect.CurTree.ParentTree = qt
			pRect.CurTree = &qt.SubTrees[index]
			qt.SubTrees[index].Insert(pRect)
			return
		}
	}
	// If we don't SubTrees within the Quadtree
	qt.Objects = append(qt.Objects, *pRect)

	// If total objects is greater than max objects and level is less than max levels
	if (len(qt.Objects) > qt.MaxObjects) && (qt.Level < qt.MaxLevels) {

		// split if we don't already have SubTrees
		if len(qt.SubTrees) > 0 == false {
			qt.split()
		}

		// Add all objects to there corresponding SubTrees
		for i < len(qt.Objects) {

			index = qt.getIndex(qt.Objects[i])

			if index != -1 {
				splice := qt.Objects[i]                                  // Get the object out of the slice
				qt.Objects = append(qt.Objects[:i], qt.Objects[i+1:]...) // Remove the object from the slice
				qt.SubTrees[index].Insert(&splice)
			} else {

				i++

			}

		}

	}



}

// Retrieve - Return all objects that could collide with the given object
func (qt *Quadtree) Retrieve(pRect Bounds) []Bounds {

	index := qt.getIndex(pRect)

	// Array with all detected objects
	returnObjects := qt.Objects

	//if we have SubTrees ...
	if len(qt.SubTrees) > 0 {

		//if pRect fits into a subnode ..
		if index != -1 {

			returnObjects = append(returnObjects, qt.SubTrees[index].Retrieve(pRect)...)

		} else {

			//if pRect does not fit into a subnode, check it against all SubTrees
			for i := 0; i < len(qt.SubTrees); i++ {
				returnObjects = append(returnObjects, qt.SubTrees[i].Retrieve(pRect)...)
			}

		}
	}

	return returnObjects

}

// RetrievePoints - Return all points that collide
func (qt *Quadtree) RetrievePoints(find Bounds) []Bounds {

	var foundPoints []Bounds
	potentials := qt.Retrieve(find)
	for o := 0; o < len(potentials); o++ {

		// X and Ys are the same and it has no Width and Height (Point)
		xyMatch := potentials[o].X == float64(find.X) && potentials[o].Y == float64(find.Y)
		if xyMatch && potentials[o].IsPoint() {
			foundPoints = append(foundPoints, find)
		}
	}

	return foundPoints

}

// RetrieveIntersections - Bring back all the bounds in a Quadtree that intersect with a provided bounds
func (qt *Quadtree) RetrieveIntersections(find Bounds) []Bounds {

	var foundIntersections []Bounds

	potentials := qt.Retrieve(find)
	for o := 0; o < len(potentials); o++ {
		if potentials[o].Intersects(find) {
			foundIntersections = append(foundIntersections, potentials[o])
		}
	}

	return foundIntersections

}

//Clear - Clear the Quadtree
func (qt *Quadtree) Clear() {

	qt.Objects = []Bounds{}

	if len(qt.SubTrees)-1 > 0 {
		for i := 0; i < len(qt.SubTrees); i++ {
			qt.SubTrees[i].Clear()
		}
	}

	qt.SubTrees = []Quadtree{}
	qt.Total = 0

}

func (qt * Quadtree) PrintTree(tab string){
	var recursivetab = tab
	for i:=0; i<len(qt.SubTrees); i++{
		fmt.Printf("%sSubtree %d: ", tab, i)
		if(len(qt.SubTrees[i].SubTrees)>0){
			fmt.Print(qt.SubTrees[i].Bounds)
			fmt.Print(" ")
			fmt.Print(qt.SubTrees[i].Objects)
			fmt.Println()
			recursivetab = tab+"\t"
			qt.SubTrees[i].PrintTree(recursivetab)
		} else{
			fmt.Print(qt.SubTrees[i])
		}
		fmt.Println()
	}
}

func (b * Bounds) WithinDistance(radius float64) []Bounds{
	var withinDist []Bounds

	if(b.CurTree != nil) {
		if (b.CurTree.Objects != nil) {
			for i := 0; i < len(b.CurTree.Objects); i++ {
				yDist := b.Y - b.CurTree.Objects[i].Y
				xDist := b.X - b.CurTree.Objects[i].X
				radDist := math.Sqrt(yDist*yDist + xDist*xDist)
				if (radDist <= radius) {
					withinDist = append(withinDist, b.CurTree.Objects[i])
				}
			}
		}
		if (b.CurTree.SubTrees != nil) {
			for j := 0; j < len(b.CurTree.SubTrees); j++ {
				yDist := b.Y - b.CurTree.SubTrees[j].Bounds.Y
				xDist := b.X - b.CurTree.SubTrees[j].Bounds.X
				radDist := math.Sqrt(yDist*yDist + xDist*xDist)
				b.CurTree.SubTrees[j].Bounds.WithinDistance(math.Abs(radius - radDist))
			}
		}

		if (b.CurTree.ParentTree != nil) {
			yDist := b.Y - b.CurTree.ParentTree.Bounds.Y
			xDist := b.X - b.CurTree.ParentTree.Bounds.Y
			radDist := math.Sqrt(yDist*yDist + xDist*xDist)
			b.CurTree.ParentTree.Bounds.WithinDistance(math.Abs(radius - radDist))
		}
	}

	return withinDist
}