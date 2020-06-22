package main

/*func main() {
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
			X: 45,
			Y: 3,
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

	//for i := 0; i < 100; i += 5 {
	//	for j := 0; j < 100; j += 3 {
	//		nodes = append(nodes, &cps.NodeImpl{X: float32(i), Y: float32(j)})
	//	}
	//}

	network := cps.AdHocNetwork{Threshold: 8, NodeTree: &tree, NodeList: []*cps.NodeImpl{}, DegreeWeight: 0.6, BatteryWeight: 0.4, Penalty: 0.7}

	for i, node := range nodes {
		tree.Insert(node)
		network.NodeList = append(network.NodeList, node)
		node.Id = i
		node.NodeClusterParams = &cps.ClusterMemberParams{RecvMsgs: []*cps.HelloMsg{}}
		//node.IsClusterMember = false
	}

	network.FullRecluster(8)
	//for _, node := range nodes {
	//	network.SendHelloMessage(8, node, &tree)
	//}
	//
	//for _, node := range nodes {
	//	network.ElectClusterHead(node)
	//}

	//for _, head := range network.ClusterHeads {
	//	network.FormClusters(head)
	//}

	//network.FinalizeClusters()

	tree.PrintTree()

	newNode := cps.NodeImpl{
		X: 52,
		Y: 3,
		NodeClusterParams: &cps.ClusterMemberParams{RecvMsgs: []*cps.HelloMsg{}},
		//Battery: 4,
	}

	tree.Insert(&newNode)
	network.NodeList = append(network.NodeList, &newNode)
	newNode.Id = 8
	newNode.NodeClusterParams = &cps.ClusterMemberParams{RecvMsgs: []*cps.HelloMsg{}}

	network.NewNodeHello(8, &newNode, &tree)
	tree.PrintTree()
}*/
