package main

import "./cps"

func main() {
	tree := cps.Quadtree{
		Bounds: cps.Bounds{
			X:      0,
			Y:      0,
			Width:  100,
			Height: 100,
		},
		MaxObjects: 1,
		MaxLevels:  50,
		Objects:    []*cps.NodeImpl{},
		SubTrees:   []*cps.Quadtree{},
	}

	nodes := []*cps.NodeImpl{
		{
			X: 20,
			Y: 2,
		},

		{
			X: 55,
			Y: 5,
		},

		{
			X: 27,
			Y: 5,
		},

		{
			X: 55,
			Y: 60,
		},

		{
			X: 40,
			Y: 2,
		},

		{
			X: 40,
			Y: 15,
		},

		{
			X: 40,
			Y: 20,
		},

		{
			X: 41,
			Y: 16,
		},
	}

	for i, node := range nodes {
		tree.Insert(node)
		node.Id = i
		node.NodeClusterParams = &cps.ClusterMemberParams{RecvMsgs: []*cps.HelloMsg{}}
		node.IsClusterMember = false
	}

	network := cps.AdHocNetwork{Threshold: 8}

	for _, node := range nodes {
		network.SendHelloMessage(8, node, &tree)
	}

	for _, node := range nodes {
		network.ElectClusterHead(node)
	}

	//for _, head := range network.ClusterHeads {
	//	network.FormClusters(head)
	//}

	//network.FinalizeClusters()

	tree.PrintTree()
}
