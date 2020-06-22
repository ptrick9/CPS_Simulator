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
		{"",
			fields {
			250.0,
			250.0,
			500.0,
			500.0,
			nil,
			},
		args {
			Bounds {
				50.0,
				50.0,
				300.0,
				300.0,
				nil,
			},
		},
		true,
		},
		{"",
			fields {
				250.0,
				250.0,
				500.0,
				500.0,
				nil,
			},
			args {
				Bounds {
					50.0,
					50.0,
					300.0,
					199.0,
					nil,
				},
			},
			false,
		},
		{"",
			fields {
				250.0,
				250.0,
				500.0,
				500.0,
				nil,
			},
			args {
				Bounds {
					50.0,
					50.0,
					300.0,
					199.0,
					nil,
				},
			},
			false,
		},
		{"",
			fields{
				250.0,
				250.0,
				500.0,
				500.0,
				nil,
			},
			args{
				Bounds{
					50.0,
					50.0,
					201.0,
					201.0,
					nil,
				},
			},
			true,
		},
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
		{
			"",
			args {
				&NodeImpl {
					X: 5,
					Y: 5,
				},
				2,
			},
			Bounds {
				3,
				3,
				4,
				4,
				nil,
			},
		},
		{
			"",
			args {
				&NodeImpl {
					X: 20,
					Y: 300,
				},
				15,
			},
			Bounds {
				5,
				285,
				30,
				30,
				nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetSearchBounds(tt.args.node, tt.args.radius); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetSearchBounds() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestQuadtree_PrintTree (t *testing.T) {
	top := Quadtree {
		Bounds {
			X: 0,
			Y: 0,
			Width: 100,
			Height: 100,
		},
		1,
		50,
		0,
		[]*NodeImpl{},
		nil,
		[]*Quadtree{},
	}

	top.Insert(&NodeImpl{
		X: 20,
		Y: 2,
	})

	top.Insert(&NodeImpl{
		X: 55,
		Y: 5,
	})

	top.Insert(&NodeImpl{
		X: 27,
		Y: 5,
	})

	top.Insert(&NodeImpl{
		X: 55,
		Y: 60,
	})

	top.Insert(&NodeImpl{
		X: 40,
		Y: 2,
	})

	top.Insert(&NodeImpl{
		X: 40,
		Y: 15,
	})

	top.Insert(&NodeImpl{
		X: 40,
		Y: 20,
	})

	top.PrintTree()
}

/*func TestQuadtree_BringNodesUp(t *testing.T) {
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

/*func TestQuadtree_CleanUp(t *testing.T) {
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
}*/