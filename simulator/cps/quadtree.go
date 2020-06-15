package cps

import (
	"fmt"
	"strconv"
)

// Quadtree - The quadtree data structure
//based off https://github.com/JamesMilnerUK/quadtree-go/blob/master/quadtree.go

type Quadtree struct {
	Bounds     Bounds
	MaxObjects int // Maximum objects a node can hold before splitting into 4 SubTrees
	MaxLevels  int // Total max levels inside root Quadtree
	Level      int // Depth level, required for SubTrees
	Objects    []*NodeImpl
	ParentTree *Quadtree
	SubTrees   []*Quadtree
	//Total      int
}

// Bounds - A bounding box with a x,y origin and width and height
type Bounds struct {
	X       float64
	Y       float64
	Width   float64
	Height  float64
	CurTree *Quadtree //tree whose boundaries this object represents
	//CurNode *NodeImpl //pointer to current Node it represents
	//CurTestNode *ClusterNode
}

//IsPoint - Checks if a bounds object is a point or not (has no width or height)
//func (b *Bounds) IsPoint() bool {
//
//	return b.Width == 0 && b.Height == 0
//
//}

// Intersects - Checks if a Bounds object intersects with another Bounds
func (b *Bounds) Intersects(a Bounds) bool {

	aMaxX := a.X + a.Width
	aMaxY := a.Y + a.Height
	bMaxX := b.X + b.Width
	bMaxY := b.Y + b.Height

	//		 a left of b	a right of b	a above b	   a below b
	return !(aMaxX < b.X || a.X > bMaxX || aMaxY < b.Y || a.Y > bMaxY)

}

func (qt *Quadtree) NodeMovement(movingNode *NodeImpl) {
	movedDown := float64(movingNode.Y) < movingNode.CurTree.Bounds.Y
	movedUp := float64(movingNode.Y) > movingNode.CurTree.Bounds.Y+movingNode.CurTree.Bounds.Height
	movedRight := float64(movingNode.X) < movingNode.CurTree.Bounds.X
	movedLeft := float64(movingNode.X) > movingNode.CurTree.Bounds.X+movingNode.CurTree.Bounds.Width

	//if moved out in any direction, reinsert in tree
	if movedDown || movedUp || movedRight || movedLeft {
		qt.RemoveAndClean(movingNode)
		qt.Insert(movingNode)
		//fmt.Printf("INSERT IN TREE: Moving Nodeto (%.4f,%.4f)\n", movingNode.X,movingNode.Y)
	}
}

// TotalSubTrees - Retrieve the total number of sub-Quadtrees in a Quadtree
func (qt *Quadtree) TotalSubTrees() int {

	total := len(qt.SubTrees)

	if total > 0 {
		for i := 0; i < len(qt.SubTrees); i++ {
			total += qt.SubTrees[i].TotalSubTrees()
		}
	}

	return total

}

// Total - Retrieve the total number of nodes in a quadtree
func (qt *Quadtree) TotalNodes() int {

	// if tree has no subTrees, return the number of objects in the tree
	if len(qt.SubTrees) <= 0 {
		return len(qt.Objects)
	}

	// if the tree has subtrees, sum and return the number of objects in each subtree
	total := 0
	for i := 0; i < len(qt.SubTrees); i++ {
		total += qt.SubTrees[i].TotalNodes()
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

	//initialize all the subtrees
	for i := 0; i < 4; i++ {
		qt.SubTrees = append(qt.SubTrees, &Quadtree{
			Bounds: Bounds{
				Width:  subWidth,
				Height: subHeight,
			},
			MaxObjects: qt.MaxObjects,
			MaxLevels:  qt.MaxLevels,
			Level:      nextLevel,
			Objects:    make([]*NodeImpl, 0),
			ParentTree: qt,
			SubTrees:   make([]*Quadtree, 0, 4),
		})
		qt.SubTrees[i].Bounds.CurTree = qt.SubTrees[i]
	}

	//top right node (0)
	qt.SubTrees[0].Bounds.X, qt.SubTrees[0].Bounds.Y = x+subWidth, y

	//top left node (1)
	qt.SubTrees[1].Bounds.X, qt.SubTrees[1].Bounds.Y = x, y

	//bottom left node (2)
	qt.SubTrees[2].Bounds.X, qt.SubTrees[2].Bounds.Y = x, y+subHeight

	//bottom right node (3)
	qt.SubTrees[3].Bounds.X, qt.SubTrees[3].Bounds.Y = x+subWidth, y+subHeight

}

// getIndex - Determine which quadrant the object belongs to (0-3)
func (qt *Quadtree) getIndex(node NodeImpl) int {

	index := -1 // index of the subnode (0-3), or -1 if pRect cannot completely fit within a subnode and is part of the parent node

	verticalMidpoint := qt.Bounds.X + (qt.Bounds.Width / 2)
	horizontalMidpoint := qt.Bounds.Y + (qt.Bounds.Height / 2)

	x := float64(node.X)
	y := float64(node.Y)

	//pRect can completely fit within the top quadrants
	//topQuadrant := (pRect.Y < horizontalMidpoint) && (pRect.Y+pRect.Height < horizontalMidpoint)

	//pRect can completely fit within the bottom quadrants
	//bottomQuadrant := (pRect.Y > horizontalMidpoint)

	//pRect can completely fit within the left quadrants
	//if (pRect.X < verticalMidpoint) && (pRect.X+pRect.Width < verticalMidpoint) {
	//println("get index:")
	//print("x: ")
	//println(int(x))
	//print("y: ")
	//println(int(y))
	//print("vert: ")
	//println(int(verticalMidpoint))
	//print("hori: ")
	//println(int(horizontalMidpoint))
	if x < verticalMidpoint { //pRect fits within the left quadrants

		if y < horizontalMidpoint {
			index = 1 //top-left
		} else {
			index = 2 //bottom-left
		}

	} else { //pRect fits within the right quadrants

		if y < horizontalMidpoint {
			index = 0 //top-right
		} else {
			index = 3 //bottom-right
		}

	}

	return index

}

// Insert - Insert the object into the tree. If the tree exceeds its capacity,
// it will split and add all objects to its corresponding SubTrees.
func (qt *Quadtree) Insert(node *NodeImpl) {
	//print("X: ")
	//println(int(node.X))
	//print("Y: ")
	//println(int(node.Y))
	//pRect.CurTree = qt
	//qt.Total++

	i := 0
	var index int

	// If we have SubTrees within the Quadtree
	//print("insert: ")
	//println(len(qt.SubTrees))
	if len(qt.SubTrees) > 0 {

		index = qt.getIndex(*node)

		//if index != -1 {
		//node.CurTree = qt.SubTrees[index]
		//node.CurTree.ParentTree = qt
		//print("index: ")
		//println(index)
		qt.SubTrees[index].Insert(node)
		return
		//}
	}
	// If we don't have SubTrees within the Quadtree
	qt.Objects = append(qt.Objects, node)

	// If total objects is greater than max objects and level is less than max levels
	//println(len(qt.Objects) > qt.MaxObjects)
	if len(qt.Objects) > qt.MaxObjects && qt.Level < qt.MaxLevels {

		// split if we don't already have SubTrees
		//if len(qt.SubTrees) > 0 == false {
		qt.split()
		//}

		// Add all objects to their corresponding SubTrees
		for i < len(qt.Objects) {

			index = qt.getIndex(*qt.Objects[i])

			//if index != -1 {
			splice := qt.Objects[i]                                  // Get the object out of the slice
			//println(len(qt.Objects))
			qt.Objects = append(qt.Objects[:i], qt.Objects[i+1:]...) // Remove the object from the slice
			//println(len(qt.Objects))
			//print("index: ")
			//println(index)
			splice.CurTree = qt.SubTrees[index]
			//splice.CurTree.ParentTree = qt
			//splice.CurTree.Objects = append(splice.CurTree.Objects, splice)
			qt.SubTrees[index].Insert(splice)
			//} else {
			//	i++
			//}
		}
	}
}

// Retrieve - Return all objects that could collide with the given object
//func (qt *Quadtree) Retrieve(pRect Bounds) []*Bounds {
//
//	index := qt.getIndex(pRect)
//
//	// Array with all detected objects
//	returnObjects := qt.Objects
//
//	//if we have SubTrees ...
//	if len(qt.SubTrees) > 0 {
//
//		//if pRect fits into a subnode ..
//		if index != -1 {
//
//			returnObjects = append(returnObjects, qt.SubTrees[index].Retrieve(pRect)...)
//
//		} else {
//
//			//if pRect does not fit into a subnode, check it against all SubTrees
//			for i := 0; i < len(qt.SubTrees); i++ {
//				returnObjects = append(returnObjects, qt.SubTrees[i].Retrieve(pRect)...)
//			}
//
//		}
//	}
//
//	return returnObjects
//
//}

// RetrievePoints - Return all points that collide
//func (qt *Quadtree) RetrievePoints(find Bounds) []Bounds {
//
//	var foundPoints []Bounds
//	potentials := qt.Retrieve(find)
//	for o := 0; o < len(potentials); o++ {
//
//		// X and Ys are the same and it has no Width and Height (Point)
//		xyMatch := potentials[o].X == find.X && potentials[o].Y == find.Y
//		if xyMatch && potentials[o].IsPoint() {
//			foundPoints = append(foundPoints, find)
//		}
//	}
//
//	return foundPoints
//
//}

// RetrieveIntersections - Bring back all the bounds in a Quadtree that intersect with a provided bounds
//func (qt *Quadtree) RetrieveIntersections(find Bounds) []Bounds {
//
//	var foundIntersections []Bounds
//
//	potentials := qt.Retrieve(find)
//	for o := 0; o < len(potentials); o++ {
//		if potentials[o].Intersects(find) {
//			foundIntersections = append(foundIntersections, *potentials[o])
//		}
//	}
//
//	return foundIntersections
//
//}

//Clear - Clear the Quadtree
func (qt *Quadtree) Clear() {

	qt.Objects = []*NodeImpl{}

	if len(qt.SubTrees)-1 > 0 {
		for i := 0; i < len(qt.SubTrees); i++ {
			qt.SubTrees[i].Clear()
		}
	}

	qt.SubTrees = []*Quadtree{}
	//qt.Total = 0

}

//PrintTree - Prints the Tree, its SubTrees, and all objects in the subtree in a clean manner
//			- helps see the hierarchy of the tree
//func (qt *Quadtree) PrintTree(tab string) {
//	var recursiveTab = tab
//	for i := 0; i < len(qt.SubTrees); i++ {
//		fmt.Printf("%sSubtree %d: ", tab, i)
//		if qt.SubTrees != nil && qt.SubTrees[i] != nil {
//			if qt.SubTrees[i].SubTrees != nil {
//				if len(qt.SubTrees[i].SubTrees) > 0 {
//					fmt.Print(qt.SubTrees[i].Bounds)
//					fmt.Print(" ")
//					fmt.Print(qt.SubTrees[i].Objects)
//					fmt.Print()
//					fmt.Print(qt.SubTrees[i].TotalNodes())
//					fmt.Println()
//					recursiveTab = tab + "\t"
//					qt.SubTrees[i].PrintTree(recursiveTab)
//				} else {
//					fmt.Print(qt.SubTrees[i])
//					fmt.Println()
//				}
//			}
//		}
//	}
//}

//PrintTree - Prints a visualization of the tree in the console
func (qt *Quadtree) PrintTree() {
	grid := [][]string{}

	top := []string{}
	bottom := make([]string, int(qt.Bounds.Width))
	for i := 0; i < int(qt.Bounds.Width); i++ {
		top = append(top, "-")
	}
	copy(bottom, top)

	row := []string{}
	row = append(row, "|")
	for i := 1; i < int(qt.Bounds.Width) - 1; i++ {
		row = append(row, " ")
	}
	row = append(row, "|")

	grid = append(grid, top)
	for i := 1; i < int(qt.Bounds.Height) - 1; i++ {
		grid = append(grid, make([]string, int(qt.Bounds.Width)))
		copy(grid[i], row)
	}
	grid = append(grid, bottom)

	if len(qt.SubTrees) > 0 {
		qt.PrintHelper(grid)
	} else {
		for _, node := range qt.Objects {
			if node.NodeClusterParams != nil && node.NodeClusterParams.CurrentCluster != nil {
				grid[int(node.Y)][int(node.X)] = strconv.Itoa(node.NodeClusterParams.CurrentCluster.ClusterNum)
			} else {
				grid[int(node.Y)][int(node.X)] = "x"
			}
		}
	}

	for i := 0; i < int(qt.Bounds.Height); i++ {
		for j := 0; j < int(qt.Bounds.Width); j++ {
			fmt.Print(grid[i][j])
		}
		fmt.Println()
	}
}

//PrintHelper - Helps PrintTree create the visualization by dealing with SubTrees
func (qt *Quadtree) PrintHelper(grid [][]string) {
	xStart := int(qt.Bounds.X)
	xMid := int(qt.Bounds.X + qt.Bounds.Width/2)
	xEnd := int(qt.Bounds.X + qt.Bounds.Width)
	yStart := int(qt.Bounds.Y)
	yMid := int(qt.Bounds.Y + qt.Bounds.Height/2)
	yEnd := int(qt.Bounds.Y + qt.Bounds.Height)
	//print("xStart: ")
	//println(xStart)
	//print("xMid: ")
	//println(xMid)
	//print("xEnd: ")
	//println(xEnd)
	//print("yStart: ")
	//println(yStart)
	//print("yMid: ")
	//println(yMid)
	//print("yEnd: ")
	//println(yEnd)
	for i := yStart + 1; i < yEnd; i++ {
		grid[i][xMid] = "|"
	}

	for i := xStart; i < xEnd; i++ {
		grid[yMid][i] = "-"
	}

	for _, st := range qt.SubTrees {
		if len(st.SubTrees) > 0 {
			st.PrintHelper(grid)
		} else {
			for _, node := range st.Objects {
				if node.NodeClusterParams != nil && node.NodeClusterParams.CurrentCluster != nil {
					grid[int(node.Y)][int(node.X)] = strconv.Itoa(node.NodeClusterParams.CurrentCluster.ClusterNum)
				} else {
					grid[int(node.Y)][int(node.X)] = "x"
				}
			}
		}
	}
}

//Remove - removes a node (bounds) from the tree, DOES NOT reconfigure the tree
func (qt *Quadtree) Remove(node *NodeImpl) bool {

	//remove from Objects in Current Tree
	for i := 0; i < len(node.CurTree.Objects); i++ {
		if node.CurTree.Objects[i] == node {
			node.CurTree.Objects = append(node.CurTree.Objects[:i], node.CurTree.Objects[i+1:]...) //remove from objects
			node.CurTree = nil
			return true
		}
	}
	return false
}

//CleanEntireTree - cleans tree from top to bottom, touching every subtree
//func (qt *Quadtree) CleanEntireTree() {
//	total := qt.TotalNodes()
//	if len(qt.SubTrees) > 0 {
//		if total > qt.MaxObjects {
//			//cleanup each subtree
//			for i := 0; i < len(qt.SubTrees); i++ {
//				qt.SubTrees[i].CleanEntireTree()
//			}
//		} else if total > 0 {
//			//if parent holds nodes, but less than or equal to the max allowed at a single level, move objects to this level (subtrees are redundant)
//			qt.BringNodesUp()
//			qt.SubTrees = []*Quadtree{}
//		} else {
//			//all objects were removed from the subtrees, so subtrees can be removed
//			qt.SubTrees = []*Quadtree{}
//		}
//	}
//}

//CleanUp - Cleans tree upwards from given quadTree, should be called after remove
func (qt * Quadtree) CleanUp(){

	if qt.TotalNodes() <= qt.MaxObjects {
		//if tree holds nodes, but less than or equal to the max allowed at a single level, move objects to this level (subtrees are redundant)
		qt.BringNodesUp()
		qt.SubTrees = []*Quadtree{}
		if qt.ParentTree != nil {
			qt.ParentTree.CleanUp()
		}
	}

}

//Removes a node from the tree and cleans the tree
func (qt * Quadtree) RemoveAndClean(node *NodeImpl){
	parent := node.CurTree.ParentTree
	qt.Remove(node)
	parent.CleanUp()
}

//BringNodesUp - Brings nodes from SubTrees of qt to qt's Objects array
func (qt *Quadtree) BringNodesUp() {

	for i := 0; i < len(qt.SubTrees); i++ {
		if len(qt.SubTrees[i].SubTrees) > 0 {
			qt.SubTrees[i].BringNodesUp()
		}
		for j := 0; j < len(qt.SubTrees[i].Objects); j++ {
			qt.Objects = append(qt.Objects, qt.SubTrees[i].Objects[j])
		}
	}

}

//Creates Search-bounds for a Node/Bounds b with a given search radius
//output of this function is used as input for WithinRadius()
func GetSearchBounds(node *NodeImpl, radius float64) Bounds {
	//returns Bounds that is a square of 2r x 2r, with a center at the Node (Bounds b)

	searchBounds := Bounds{
		X:       float64(node.X) - radius,
		Y:       float64(node.Y) - radius,
		Width:   2 * radius,
		Height:  2 * radius,
	}
	return searchBounds
}

//Traverses the tree finding all points that intersect with the searchBounds and also are within in a radius distance of center point
func (qt *Quadtree) WithinRadius(radius float64, center *NodeImpl, withinDist []*NodeImpl) []*NodeImpl {

	searchBounds := GetSearchBounds(center, radius)

	//First traverse through subtrees. If there are any subtrees there are no objects to check in the current tree
	if qt.SubTrees != nil && len(qt.SubTrees) > 0 {
		for i := 0; i < len(qt.SubTrees); i++ {
			//Only move down a level if the tree bounds intersect with the search bounds
			if qt.SubTrees[i].Bounds.Intersects(searchBounds) {
				withinDist = qt.SubTrees[i].WithinRadius(radius, center, withinDist)
			}
		}
		//Traverse objects in the subtree.
	} else {
		if qt.Objects != nil && len(qt.Objects) > 0 {
			for i := 0; i < len(qt.Objects); i++ {
				if qt.Objects[i] != center {
					//actual search area is a circle, search bounds is a square
					//its very likely all points in the bounds are also in the area, but not always true
					//so check with the distance formula

					//if the distance is less than the search radius then it is in the search area, thus add to array
					if center.IsWithinRange(qt.Objects[i], radius) {
						withinDist = append(withinDist, qt.Objects[i])
					}
				}
			}
		}
	}
	return withinDist
}
