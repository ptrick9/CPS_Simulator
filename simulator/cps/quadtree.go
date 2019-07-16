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
	Objects    []*Bounds
	ParentTree *Quadtree
	SubTrees   []*Quadtree
	Total      int
}

// Bounds - A bounding box with a x,y origin and width and height
type Bounds struct {
	X      	float64
	Y      	float64
	Width  	float64
	Height 	float64
	CurTree	*Quadtree 	//tree that this object is in
	CurNode *NodeImpl	//pointer to current Node it represents
	//CurTestNode *ClusterNode
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
	if aMaxX < b.X || a.X > bMaxX || aMaxY < b.Y || a.Y > bMaxY{
		return false
	}else {
		return true
	}
}

func (qt *Quadtree) NodeMovement(movingNode * Bounds){
	movedDown := movingNode.Y<movingNode.CurTree.Bounds.Y
	movedUp := movingNode.Y>movingNode.CurTree.Bounds.Y + movingNode.CurTree.Bounds.Height
	movedRight := movingNode.X < movingNode.CurTree.Bounds.X
	movedLeft := movingNode.X > movingNode.CurTree.Bounds.X + movingNode.CurTree.Bounds.Width

	//if moved out in any direction, reinsert in tree
	if(movedDown || movedUp || movedRight || movedLeft){
		qt.RemoveAndCleanup(movingNode)
		qt.Insert(movingNode)
		//fmt.Printf("INSERT IN TREE: Moving Nodeto (%.4f,%.4f)\n", movingNode.X,movingNode.Y)
	}
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
		Objects:    make([]*Bounds, 0),
		ParentTree: qt,
		SubTrees:   make([]*Quadtree, 0, 4),
	}
	sub0.Bounds.CurTree = &sub0
	qt.SubTrees = append(qt.SubTrees, &sub0)

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
		Objects:    make([]*Bounds, 0),
		ParentTree: qt,
		SubTrees:   make([]*Quadtree, 0, 4),
	}
	sub1.Bounds.CurTree = &sub1
	qt.SubTrees = append(qt.SubTrees, &sub1)


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
		Objects:    make([]*Bounds, 0),
		ParentTree: qt,
		SubTrees:   make([]*Quadtree, 0, 4),
	}
	sub2.Bounds.CurTree = &sub2
	qt.SubTrees = append(qt.SubTrees, &sub2)

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
		Objects:    make([]*Bounds, 0),
		ParentTree: qt,
		SubTrees:   make([]*Quadtree, 0, 4),
	}
	sub3.Bounds.CurTree = &sub3
	qt.SubTrees = append(qt.SubTrees, &sub3)

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
	//pRect.CurTree = qt
	qt.Total++

	i := 0
	var index int

	// If we have SubTrees within the Quadtree
	if (len(qt.SubTrees) > 0) {

		index = qt.getIndex(*pRect)

		if index != -1 {
			pRect.CurTree = qt.SubTrees[index]
			pRect.CurTree.ParentTree = qt
			qt.SubTrees[index].Insert(pRect)
			return
		}
	}
	// If we don't SubTrees within the Quadtree
	qt.Objects = append(qt.Objects, pRect)

	// If total objects is greater than max objects and level is less than max levels
	if (len(qt.Objects) > qt.MaxObjects) && (qt.Level < qt.MaxLevels) {

		// split if we don't already have SubTrees
		if len(qt.SubTrees) > 0 == false {
			qt.split()
		}

		// Add all objects to there corresponding SubTrees
		for i < len(qt.Objects) {

			index = qt.getIndex(*qt.Objects[i])

			if index != -1 {
				splice := qt.Objects[i]                                  // Get the object out of the slice
				qt.Objects = append(qt.Objects[:i], qt.Objects[i+1:]...) // Remove the object from the slice
				pRect.CurTree = qt.SubTrees[index]
				pRect.CurTree.ParentTree = qt
				splice.CurTree = qt.SubTrees[index]
				splice.CurTree.ParentTree = qt
				qt.SubTrees[index].Insert(splice)
			} else {
				i++
			}
		}
	}
}

// Retrieve - Return all objects that could collide with the given object
func (qt *Quadtree) Retrieve(pRect Bounds) []*Bounds {

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
			foundIntersections = append(foundIntersections, *potentials[o])
		}
	}

	return foundIntersections

}

//Clear - Clear the Quadtree
func (qt *Quadtree) Clear() {

	qt.Objects = []*Bounds{}

	if len(qt.SubTrees)-1 > 0 {
		for i := 0; i < len(qt.SubTrees); i++ {
			qt.SubTrees[i].Clear()
		}
	}

	qt.SubTrees = []*Quadtree{}
	qt.Total = 0

}

//PrintTree - Prints the Tree, its SubTrees, and all objects in the subtree in a clean manner
//			- helps see the hierarchy of the tree
func (qt * Quadtree) PrintTree(tab string){
	var recursivetab = tab
	for i:=0; i<len(qt.SubTrees); i++{
		fmt.Printf("%sSubtree %d: ", tab, i)
		if(qt.SubTrees!=nil && qt.SubTrees[i]!=nil) {
			if (qt.SubTrees[i].SubTrees != nil) {
				if (len(qt.SubTrees[i].SubTrees) > 0) {
					fmt.Print(qt.SubTrees[i].Bounds)
					fmt.Print(" ")
					fmt.Print(qt.SubTrees[i].Objects)
					fmt.Print()
					fmt.Print(qt.SubTrees[i].Total)
					fmt.Println()
					recursivetab = tab + "\t"
					qt.SubTrees[i].PrintTree(recursivetab)
				} else {
					fmt.Print(qt.SubTrees[i])
					fmt.Println()
				}
			}
		}
	}
}

//Remove - removes a node (bounds) from the tree, DOES NOT reconfigure the tree
func (qt * Quadtree) Remove(pRect * Bounds) * Bounds{

	//remove from Objects in Current Tree
	for i:=0; i<len(pRect.CurTree.Objects); i++{
		if(pRect.CurTree.Objects[i] == pRect){
			pRect.CurTree.Objects = append(pRect.CurTree.Objects[:i], pRect.CurTree.Objects[i+1:]...) //remove from objects
			break
		}
	}

	//update totals in current tree and all parent trees
	curTree := pRect.CurTree
	for curTree.ParentTree != nil{
		curTree.Total = curTree.Total-1
		curTree = curTree.ParentTree
	}

	return pRect
}

func (qt * Quadtree) CleanUp(){
	if(len(qt.SubTrees)==4){
		if(qt.Total > 4){ //traverse downward
			qt.SubTrees[0].CleanUp()
			qt.SubTrees[1].CleanUp()
			qt.SubTrees[2].CleanUp()
			qt.SubTrees[3].CleanUp()
		}else if(qt.Total <=4 && qt.Total>0){
			//if parent holds between 1 to 4 nodes and has 4 subtrees, move objects up one level (subtrees are redundant)
			for i:=0; i<len(qt.SubTrees); i++{
				for j:=0; j<len(qt.SubTrees[i].Objects); j++{
					qt.Objects = append(qt.Objects, qt.SubTrees[i].Objects[j])
				}
			}
			qt.SubTrees = []*Quadtree{}
		} else if (qt.Total == 0){
			//all objects were removed from the subtrees, so subtrees can be removed
			qt.SubTrees = []*Quadtree{}
		}
	}
}

//RemoveAndCleanup - Removes a node (bounds) from the tree, reconfigures the tree if neccessary
func (qt * Quadtree) RemoveAndCleanup(pRect * Bounds) *Bounds{

	for i:=0; i<len(pRect.CurTree.Objects); i++{
		if(pRect.CurTree.Objects[i] == pRect){
			pRect.CurTree.Objects = append(pRect.CurTree.Objects[:i], pRect.CurTree.Objects[i+1:]...) //remove from objects
			break
		}
	}

	//update totals
	curTree := pRect.CurTree
	for curTree.ParentTree != nil{
		curTree.Total = curTree.Total-1
		curTree = curTree.ParentTree
	}

	//if parent holds four, move all nodes up one level
	parent := pRect.CurTree.ParentTree
	if(parent != nil && parent.Total==4){
		for i:=0; i<len(parent.SubTrees); i++{
			for j:=0; j<len(parent.SubTrees[i].Objects); j++{
				parent.Objects = append(parent.Objects, parent.SubTrees[i].Objects[j])
			}
		}
		parent.SubTrees = []*Quadtree{}

		//once all objects are in parent, update objects to hold pointer to new CurTree
		for i:=0; i<len(parent.Objects); i++{
			parent.Objects[i].CurTree = parent
		}
	}


	return pRect
}

//Creates Search-bounds for a Node/Bounds b with a given search radius
//output of this function is used as input for WithinRadius()
func (b * Bounds) GetSearchBounds(radius float64) Bounds{
	//returns Bounds that is a square of 2r x 2r, with a center at the Node (Bounds b)

	searchBounds := Bounds{
		X: b.X-radius,
		Y: b.Y-radius,
		Width: 2*radius,
		Height: 2*radius,
		CurTree: b.CurTree,
	}
	return searchBounds
}

//Traverses the tree finding all points that intersect with the searchBounds and also are within in a radius distance of center point
func (qt * Quadtree) WithinRadius(radius float64, center * Bounds, searchBounds Bounds, withinDist []*Bounds) []*Bounds{

	//First traverse through subtrees. If there are any subtrees there are no objects to check in the current tree
	if(qt.SubTrees !=nil && len(qt.SubTrees) > 0){
		for i:=0; i<len(qt.SubTrees);i++{
			//Only move down a level if the tree bounds intersect with the search bounds
			if(qt.SubTrees[i].Bounds.Intersects(searchBounds)){
				withinDist = qt.SubTrees[i].WithinRadius(radius,center,searchBounds,withinDist)
			}
		}
	//Traverse objects in the subtree.
	} else{
		if(qt.Objects != nil && len(qt.Objects)>0){
			for i:=0; i<len(qt.Objects);i++{
				//Only consider if the object intersects with the search bounds
				if(qt.Objects[i].Intersects(searchBounds)){
					if(qt.Objects[i]!=center){
						//actual search area is a circle, search bounds is a square
						//its very likely all points in the bounds are also in the area, but not always true
						//so check with the distance formula

						yDist := center.Y - qt.Objects[i].Y
						xDist := center.X - qt.Objects[i].X
						radDist := math.Sqrt(yDist*yDist + xDist*xDist)

						//if the distance is less than the search radius then it is in the search area, thus add to array
						if(radDist <= radius){
							withinDist = append(withinDist,qt.Objects[i])
						}
					}
				}
			}
		}
	}
	return withinDist
}