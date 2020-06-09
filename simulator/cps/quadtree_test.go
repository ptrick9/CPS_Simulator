package cps

import (
	"reflect"
	"testing"
)

func TestBounds_Intersects(t *testing.T) {
	type fields struct {
		X       float64
		Y       float64
		Width   float64
		Height  float64
		CurTree *Quadtree
	}
	type args struct {
		a Bounds
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &Bounds{
				X:       tt.fields.X,
				Y:       tt.fields.Y,
				Width:   tt.fields.Width,
				Height:  tt.fields.Height,
				CurTree: tt.fields.CurTree,
			}
			if got := b.Intersects(tt.args.a); got != tt.want {
				t.Errorf("Intersects() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBounds_IsPoint(t *testing.T) {
	type fields struct {
		X       float64
		Y       float64
		Width   float64
		Height  float64
		CurTree *Quadtree
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &Bounds{
				X:       tt.fields.X,
				Y:       tt.fields.Y,
				Width:   tt.fields.Width,
				Height:  tt.fields.Height,
				CurTree: tt.fields.CurTree,
			}
			if got := b.IsPoint(); got != tt.want {
				t.Errorf("IsPoint() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetSearchBounds(t *testing.T) {
	type args struct {
		node   *NodeImpl
		radius float64
	}
	tests := []struct {
		name string
		args args
		want Bounds
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetSearchBounds(tt.args.node, tt.args.radius); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetSearchBounds() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestQuadtree_BringNodesUp(t *testing.T) {
	type fields struct {
		Bounds     Bounds
		MaxObjects int
		MaxLevels  int
		Level      int
		Objects    []*NodeImpl
		ParentTree *Quadtree
		SubTrees   []*Quadtree
	}
	tests := []struct {
		name   string
		fields fields
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			qt := &Quadtree{
				Bounds:     tt.fields.Bounds,
				MaxObjects: tt.fields.MaxObjects,
				MaxLevels:  tt.fields.MaxLevels,
				Level:      tt.fields.Level,
				Objects:    tt.fields.Objects,
				ParentTree: tt.fields.ParentTree,
				SubTrees:   tt.fields.SubTrees,
			}
		})
	}
}

func TestQuadtree_CleanUp(t *testing.T) {
	type fields struct {
		Bounds     Bounds
		MaxObjects int
		MaxLevels  int
		Level      int
		Objects    []*NodeImpl
		ParentTree *Quadtree
		SubTrees   []*Quadtree
	}
	tests := []struct {
		name   string
		fields fields
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			qt := &Quadtree{
				Bounds:     tt.fields.Bounds,
				MaxObjects: tt.fields.MaxObjects,
				MaxLevels:  tt.fields.MaxLevels,
				Level:      tt.fields.Level,
				Objects:    tt.fields.Objects,
				ParentTree: tt.fields.ParentTree,
				SubTrees:   tt.fields.SubTrees,
			}
		})
	}
}

func TestQuadtree_Clear(t *testing.T) {
	type fields struct {
		Bounds     Bounds
		MaxObjects int
		MaxLevels  int
		Level      int
		Objects    []*NodeImpl
		ParentTree *Quadtree
		SubTrees   []*Quadtree
	}
	tests := []struct {
		name   string
		fields fields
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			qt := &Quadtree{
				Bounds:     tt.fields.Bounds,
				MaxObjects: tt.fields.MaxObjects,
				MaxLevels:  tt.fields.MaxLevels,
				Level:      tt.fields.Level,
				Objects:    tt.fields.Objects,
				ParentTree: tt.fields.ParentTree,
				SubTrees:   tt.fields.SubTrees,
			}
		})
	}
}

func TestQuadtree_Insert(t *testing.T) {
	type fields struct {
		Bounds     Bounds
		MaxObjects int
		MaxLevels  int
		Level      int
		Objects    []*NodeImpl
		ParentTree *Quadtree
		SubTrees   []*Quadtree
	}
	type args struct {
		node *NodeImpl
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			qt := &Quadtree{
				Bounds:     tt.fields.Bounds,
				MaxObjects: tt.fields.MaxObjects,
				MaxLevels:  tt.fields.MaxLevels,
				Level:      tt.fields.Level,
				Objects:    tt.fields.Objects,
				ParentTree: tt.fields.ParentTree,
				SubTrees:   tt.fields.SubTrees,
			}
		})
	}
}

func TestQuadtree_NodeMovement(t *testing.T) {
	type fields struct {
		Bounds     Bounds
		MaxObjects int
		MaxLevels  int
		Level      int
		Objects    []*NodeImpl
		ParentTree *Quadtree
		SubTrees   []*Quadtree
	}
	type args struct {
		movingNode *NodeImpl
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			qt := &Quadtree{
				Bounds:     tt.fields.Bounds,
				MaxObjects: tt.fields.MaxObjects,
				MaxLevels:  tt.fields.MaxLevels,
				Level:      tt.fields.Level,
				Objects:    tt.fields.Objects,
				ParentTree: tt.fields.ParentTree,
				SubTrees:   tt.fields.SubTrees,
			}
		})
	}
}

func TestQuadtree_PrintTree(t *testing.T) {
	type fields struct {
		Bounds     Bounds
		MaxObjects int
		MaxLevels  int
		Level      int
		Objects    []*NodeImpl
		ParentTree *Quadtree
		SubTrees   []*Quadtree
	}
	type args struct {
		tab string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			qt := &Quadtree{
				Bounds:     tt.fields.Bounds,
				MaxObjects: tt.fields.MaxObjects,
				MaxLevels:  tt.fields.MaxLevels,
				Level:      tt.fields.Level,
				Objects:    tt.fields.Objects,
				ParentTree: tt.fields.ParentTree,
				SubTrees:   tt.fields.SubTrees,
			}
		})
	}
}

func TestQuadtree_Remove(t *testing.T) {
	type fields struct {
		Bounds     Bounds
		MaxObjects int
		MaxLevels  int
		Level      int
		Objects    []*NodeImpl
		ParentTree *Quadtree
		SubTrees   []*Quadtree
	}
	type args struct {
		node *NodeImpl
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			qt := &Quadtree{
				Bounds:     tt.fields.Bounds,
				MaxObjects: tt.fields.MaxObjects,
				MaxLevels:  tt.fields.MaxLevels,
				Level:      tt.fields.Level,
				Objects:    tt.fields.Objects,
				ParentTree: tt.fields.ParentTree,
				SubTrees:   tt.fields.SubTrees,
			}
			if got := qt.Remove(tt.args.node); got != tt.want {
				t.Errorf("Remove() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestQuadtree_TotalNodes(t *testing.T) {
	type fields struct {
		Bounds     Bounds
		MaxObjects int
		MaxLevels  int
		Level      int
		Objects    []*NodeImpl
		ParentTree *Quadtree
		SubTrees   []*Quadtree
	}
	tests := []struct {
		name   string
		fields fields
		want   int
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			qt := &Quadtree{
				Bounds:     tt.fields.Bounds,
				MaxObjects: tt.fields.MaxObjects,
				MaxLevels:  tt.fields.MaxLevels,
				Level:      tt.fields.Level,
				Objects:    tt.fields.Objects,
				ParentTree: tt.fields.ParentTree,
				SubTrees:   tt.fields.SubTrees,
			}
			if got := qt.TotalNodes(); got != tt.want {
				t.Errorf("TotalNodes() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestQuadtree_TotalSubTrees(t *testing.T) {
	type fields struct {
		Bounds     Bounds
		MaxObjects int
		MaxLevels  int
		Level      int
		Objects    []*NodeImpl
		ParentTree *Quadtree
		SubTrees   []*Quadtree
	}
	tests := []struct {
		name   string
		fields fields
		want   int
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			qt := &Quadtree{
				Bounds:     tt.fields.Bounds,
				MaxObjects: tt.fields.MaxObjects,
				MaxLevels:  tt.fields.MaxLevels,
				Level:      tt.fields.Level,
				Objects:    tt.fields.Objects,
				ParentTree: tt.fields.ParentTree,
				SubTrees:   tt.fields.SubTrees,
			}
			if got := qt.TotalSubTrees(); got != tt.want {
				t.Errorf("TotalSubTrees() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestQuadtree_WithinRadius(t *testing.T) {
	type fields struct {
		Bounds     Bounds
		MaxObjects int
		MaxLevels  int
		Level      int
		Objects    []*NodeImpl
		ParentTree *Quadtree
		SubTrees   []*Quadtree
	}
	type args struct {
		radius       float64
		center       *NodeImpl
		searchBounds Bounds
		withinDist   []*NodeImpl
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   []*NodeImpl
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			qt := &Quadtree{
				Bounds:     tt.fields.Bounds,
				MaxObjects: tt.fields.MaxObjects,
				MaxLevels:  tt.fields.MaxLevels,
				Level:      tt.fields.Level,
				Objects:    tt.fields.Objects,
				ParentTree: tt.fields.ParentTree,
				SubTrees:   tt.fields.SubTrees,
			}
			if got := qt.WithinRadius(tt.args.radius, tt.args.center, tt.args.searchBounds, tt.args.withinDist); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithinRadius() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestQuadtree_getIndex(t *testing.T) {
	type fields struct {
		Bounds     Bounds
		MaxObjects int
		MaxLevels  int
		Level      int
		Objects    []*NodeImpl
		ParentTree *Quadtree
		SubTrees   []*Quadtree
	}
	type args struct {
		node NodeImpl
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   int
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			qt := &Quadtree{
				Bounds:     tt.fields.Bounds,
				MaxObjects: tt.fields.MaxObjects,
				MaxLevels:  tt.fields.MaxLevels,
				Level:      tt.fields.Level,
				Objects:    tt.fields.Objects,
				ParentTree: tt.fields.ParentTree,
				SubTrees:   tt.fields.SubTrees,
			}
			if got := qt.getIndex(tt.args.node); got != tt.want {
				t.Errorf("getIndex() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestQuadtree_split(t *testing.T) {
	type fields struct {
		Bounds     Bounds
		MaxObjects int
		MaxLevels  int
		Level      int
		Objects    []*NodeImpl
		ParentTree *Quadtree
		SubTrees   []*Quadtree
	}
	tests := []struct {
		name   string
		fields fields
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			qt := &Quadtree{
				Bounds:     tt.fields.Bounds,
				MaxObjects: tt.fields.MaxObjects,
				MaxLevels:  tt.fields.MaxLevels,
				Level:      tt.fields.Level,
				Objects:    tt.fields.Objects,
				ParentTree: tt.fields.ParentTree,
				SubTrees:   tt.fields.SubTrees,
			}
		})
	}
}