package cps

import (
	"math"
)

type AdHocNetwork struct {
	ClusterHeads        []*NodeImpl
	SingularNodes       []*NodeImpl
	FullReclusters      int //Counts the number of full reclusters that occur in a simulation
	LocalReclusters     int //Counts the number of local reclusters that occur in a simulation
	PotentialReclusters int //Counts the number of reclusters that would have occurred if there was a check every second
	AverageClusterSize  int //The current average cluster size is added to this every iteration (when cluster print is enabled)
	AverageNumClusters  int //The current number of clusters is added to this every iteration
	CSJoins             int //Counts the number of times a new node hello leads to a node joining an existing cluster
	CSSolos             int //Counts the number of times a new node hello leads to the creation of a new cluster
	LRHeads             int //Counts the number of cluster heads created by local reclusters
	LostReadings        int //Counts the number of readings that were lost because the cluster head died before sending
	NextClusterNum      int //For testing, may remove later
	TotalWaits          int //Counts the number of times a node waited instead of initiating a cluster search
}

type Cluster struct {
	ClusterHead    *NodeImpl
	ClusterMembers []*NodeImpl
	ClusterNum     int //integer identifier of cluster
}

type ClusterMemberParams struct {
	CurrentCluster  *Cluster
	RecvMsgs        []*HelloMsg
	ThisNodeHello   *HelloMsg
	AttemptedToJoin []*NodeImpl
}

type HelloMsg struct {
	Sender      *NodeImpl //pointer to Node sending the Hello Msg
	NodeCHScore float64 //score for determining how suitable a node is to be a clusterhead
}

type ClusterNode struct { //made for testing, only has parameters that the cluster needs to know from a NodeImpl
	NodeBounds        *Bounds
	Battery           float32
	ClusterHead       *NodeImpl
	NodeClusterParams *ClusterMemberParams
	IsClusterHead     bool
	IsClusterMember   bool
}

//Computes the cluster score (higher the score the better chance a node becomes a cluster head)
func (node *NodeImpl) ComputeClusterScore(p *Params, numWithinDist int, ) float64 {
	degree := math.Min(float64(numWithinDist), float64(p.ClusterMaxThreshold))
	//battery := node.GetBatteryPercentage() * float64(p.ClusterMaxThreshold) //Multiplying by threshold ensures that battery and degree have the same maximum value

	//weighted sum, 60% from degree (# of nodes within distance), 40% from its battery life
	// penalty used to increase a nodes chance at staying a clusterhead
	//score := p.DegreeWeight*degree + p.BatteryWeight*battery
	score := degree * node.GetBatteryPercentage()
	//if node.IsClusterHead {
	return score
	//} else {
	//	return score * p.Penalty
	//}
}

//Generates Hello Message for node to form/maintain clusters
func (node *NodeImpl) GenerateHello(score float64) {
	message := &HelloMsg{ Sender: node, NodeCHScore: score}
	node.NodeClusterParams.ThisNodeHello = message
}

/*func (adhoc *AdHocNetwork) SendHelloMessage(curNode *NodeImpl, p *Params) {
	withinDist := p.NodeTree.WithinRadius(p.NodeBTRange, curNode, []*NodeImpl{})
	curNode.DrainBatteryBluetooth()	//Broadcasting first hello message, no score
	for j := 0; j < len(withinDist); j++ {
		withinDist[j].DrainBatteryBluetooth() //Every node in bluetooth range receives the first hello message, allowing them to count how many neighbors they have
	}

	if curNode.IsAlive() { //If the node survives sending the first message
		curNode.GenerateHello(curNode.ComputeClusterScore(p, len(withinDist)))
		curNode.DrainBatteryBluetooth() //Broadcasting second hello message with score

		for j := 0; j < len(withinDist); j++ {
			if withinDist[j].IsAlive() { //If the node survives receiving the first message
				withinDist[j].DrainBatteryBluetooth() //Every node in bluetooth range receives a second hello message, the one including curNode's score
				withinDist[j].NodeClusterParams.RecvMsgs = append(withinDist[j].NodeClusterParams.RecvMsgs, curNode.NodeClusterParams.ThisNodeHello)
			}
		}
	}
}

func (adhoc *AdHocNetwork) SendLocalHello(curNode *NodeImpl, members []*NodeImpl, p *Params) {
	withinDist := p.NodeTree.WithinRadius(p.NodeBTRange, curNode, []*NodeImpl{})
	membersWithinDist := SimpleIntersect(withinDist, members)

	curNode.DrainBatteryBluetooth() //Broadcasting first hello message, no score

	for j := 0; j < len(withinDist); j++ {
		withinDist[j].DrainBatteryBluetooth() //Every node in bluetooth range receives the first hello message, allowing them to count how many neighbors they have
	}

	if curNode.IsAlive() {
		curNode.GenerateHello(curNode.ComputeClusterScore(p, len(membersWithinDist))) //Score is computed using members in range because only members respond to the initial hello
		curNode.DrainBatteryBluetooth()                                               //Broadcasting second hello message with score
		for j := 0; j < len(withinDist); j++ {
			withinDist[j].DrainBatteryBluetooth() //Every node in bluetooth range then receives a second hello message, the one including curNode's score
		}

		for j := 0; j < len(membersWithinDist); j++ {
			//Of nodes in bluetooth range, only members of the old cluster add the hello message to their received messages
			membersWithinDist[j].NodeClusterParams.RecvMsgs = append(membersWithinDist[j].NodeClusterParams.RecvMsgs, curNode.NodeClusterParams.ThisNodeHello)
		}
	}
}*/

/*
Used for sending hello messages in both local and global reclusters

curNode - The node sending the hello message
members - The nodes participating in the recluster. nil if global
p		- The parameters of the simulation
*/
func (adhoc *AdHocNetwork) SendHello(curNode *NodeImpl, members []*NodeImpl, p *Params) {
	withinDist := p.NodeTree.WithinRadius(p.NodeBTRange, curNode, []*NodeImpl{})
	membersWithinDist := withinDist
	if members != nil {
		membersWithinDist = SimpleIntersect(withinDist, members)
	}

	curNode.DrainBatteryBluetooth(&curNode.P.Server.ReclusterBTCounter) //Broadcasting first hello message, no score

	for j := 0; j < len(withinDist); j++ {
		withinDist[j].DrainBatteryBluetooth(&curNode.P.Server.ReclusterBTCounter) //Every node in bluetooth range receives the first hello message, allowing them to count how many neighbors they have
	}

	if curNode.IsAlive() { //If the node survived sending the first message
		curNode.GenerateHello(curNode.ComputeClusterScore(p, len(membersWithinDist))) //Score is computed using members in range because only members respond to the initial hello
		curNode.DrainBatteryBluetooth(&curNode.P.Server.ReclusterBTCounter)                                               //Broadcasting second hello message with score
		for j := 0; j < len(withinDist); j++ {
			if withinDist[j].IsAlive() { //If the node survived receiving the first message
				withinDist[j].DrainBatteryBluetooth(&curNode.P.Server.ReclusterBTCounter) //Every node in bluetooth range then receives a second hello message, the one including curNode's score
			}
		}

		for j := 0; j < len(membersWithinDist); j++ {
			if withinDist[j].IsAlive() { //If the node survived receiving the first message
				//Of nodes in bluetooth range, only members add the hello message to their received messages
				membersWithinDist[j].NodeClusterParams.RecvMsgs = append(membersWithinDist[j].NodeClusterParams.RecvMsgs, curNode.NodeClusterParams.ThisNodeHello)
			}
		}
	}
}

func (adhoc *AdHocNetwork) ClusterMovement(node *NodeImpl, p *Params) {
	if node.Valid && node.IsAlive() {
		 if !node.IsClusterHead && p.CurrentTime >= p.InitClusterTime {
			if node.NodeClusterParams.CurrentCluster.ClusterHead == nil {
				adhoc.ClusterSearch(node, p)
			} else if !node.IsWithinRange(node.NodeClusterParams.CurrentCluster.ClusterHead, p.NodeBTRange) {
				/*	Node knows it is within range of its cluster head if it receives confirmation after sending its
					reading. This is the cost of sending the reading to the out-of-range cluster head. The cost is
					normally handled in node.go's SendToClusterHead, when the head is in range. */
				node.DrainBatteryBluetooth(&node.P.Server.ReadingBTCounter)
				if node.IsAlive() { //If the node survived sending the message to the out-of-range cluster head
					adhoc.ClusterSearch(node, p)
				}
			} else {
				node.Wait = 0
			}
		} else if !p.GlobalRecluster && node.IsClusterHead && len(node.NodeClusterParams.CurrentCluster.ClusterMembers) <= p.ClusterMinThreshold {
			adhoc.ClusterSearch(node, p)
		} else {
			node.Wait = 0
		}
	}
}

func (adhoc *AdHocNetwork) ClusterSearch(node *NodeImpl, p *Params) {
	if node.Wait < p.ClusterSearchThreshold {
		node.Wait++
		adhoc.TotalWaits++
	} else {
		node.Wait = 0
		toJoin := node.FindNearbyHead(p, &p.Server.ClusterSearchBTCounter)

		if toJoin != nil {
			adhoc.ClearClusterParams(node)
			node.Join(toJoin.NodeClusterParams.CurrentCluster)
			adhoc.CSJoins++
		} else {
			adhoc.ClearClusterParams(node)
			adhoc.ElectClusterHead(node, p)
			adhoc.CSSolos++
		}
	}
}

func (node *NodeImpl) HasMaxNodeScore(p *Params) *NodeImpl {
	maxNode := node
	maxScore := node.NodeClusterParams.ThisNodeHello.NodeCHScore //* p.Penalty
	for i := 0; i < len(node.NodeClusterParams.RecvMsgs); i++ {
		sender := node.NodeClusterParams.RecvMsgs[i].Sender
		score := node.NodeClusterParams.RecvMsgs[i].NodeCHScore
		//if !sender.IsClusterHead {
		//	score *= p.Penalty
		//}
		//do not consider nodes already with a clusterhead
		//if received a message from a node who does not have a cluster head
		if sender != nil && !sender.IsClusterMember && len(sender.NodeClusterParams.CurrentCluster.ClusterMembers) < p.ClusterMaxThreshold {
			//if their score higher than current node score
			if score > maxScore {
				maxScore = score
				maxNode = sender
			} else if math.Abs(score-maxScore) < 0.1 {
				if sender.Id < maxNode.Id {
					maxNode = sender
				}
			}
		}
	}

	return maxNode
}

func (node *NodeImpl) Join(cluster *Cluster) {
	node.IsClusterMember = true
	node.IsClusterHead = false
	node.NodeClusterParams.CurrentCluster = cluster
	cluster.ClusterMembers = append(cluster.ClusterMembers, node)
}

//Prints a representation of a node
/*func (node *NodeImpl) PrintClusterNode() {
	fmt.Print("{")
	fmt.Print(node.X)
	fmt.Print(",")
	fmt.Print(node.Y)
	fmt.Print(" ")
	fmt.Print(node.GetBatteryPercentage())
	fmt.Print(" ")
	//fmt.Print(curNode.ClusterHead)
	//fmt.Print(" ")
	fmt.Print(node.NodeClusterParams.CurrentCluster)
	fmt.Print(" ")
	fmt.Print(node.IsClusterHead)
	fmt.Print("}")
	fmt.Println()
}*/

func (adhoc *AdHocNetwork) ClearClusterParams(node *NodeImpl) {
	if node.IsClusterHead {
		node.IsClusterHead = false
		adhoc.DissolveCluster(node)
		return
	} else if node.NodeClusterParams.CurrentCluster != nil {
		for i := range node.NodeClusterParams.CurrentCluster.ClusterMembers {
			if node == node.NodeClusterParams.CurrentCluster.ClusterMembers[i] {
				node.NodeClusterParams.CurrentCluster.ClusterMembers = append(node.NodeClusterParams.CurrentCluster.ClusterMembers[:i], node.NodeClusterParams.CurrentCluster.ClusterMembers[i+1:]...)
				break
			}
		}
	}
	node.NodeClusterParams.CurrentCluster = &Cluster{}
	node.NodeClusterParams.CurrentCluster.ClusterMembers = []*NodeImpl{}

	node.NodeClusterParams.RecvMsgs = []*HelloMsg{}
	node.NodeClusterParams.ThisNodeHello = &HelloMsg{}

	node.IsClusterMember = false
}

func (adhoc *AdHocNetwork) DissolveCluster(node *NodeImpl) {
	//Assume node is cluster head
	for i := range adhoc.ClusterHeads {
		if node == adhoc.ClusterHeads[i] {
			adhoc.ClusterHeads = append(adhoc.ClusterHeads[:i], adhoc.ClusterHeads[i+1:]...)
			break
		}
	}
	for i := 0; len(node.NodeClusterParams.CurrentCluster.ClusterMembers) > 0; {
		member := node.NodeClusterParams.CurrentCluster.ClusterMembers[i]
		adhoc.ClearClusterParams(member)
		if node.P.LocalRecluster <= 0 {
			/*	Assuming local reclustering is disabled, the members only know that this cluster has dissolved because
				they will later try to send their reading to the cluster head over bluetooth and will receive no
				confirmation. That bluetooth cost has to be simulated now. */
			member.DrainBatteryBluetooth(&node.P.Server.ReclusterBTCounter)
		} else {
			/*	Assuming local reclustering is enabled, the members will know that this cluster has dissolved because
				they will be told so by the server */
			member.DrainBatteryWifi()
		}
	}
	adhoc.ClearClusterParams(node)
}

func (adhoc *AdHocNetwork) ResetClusters(p *Params) {
	for node := range p.AliveNodes {
		node.IsClusterHead = false
		p.ClusterNetwork.ClearClusterParams(node)
	}
	adhoc.ClusterHeads = []*NodeImpl{}
	adhoc.SingularNodes = []*NodeImpl{}
}

//sorts messages by distance to the node: 0th = closest, nth = farthest
/*func (node *NodeImpl) SortMessages() {

	distances := []float64{}

	msgs := node.NodeClusterParams.RecvMsgs

	for i := 0; i < len(msgs); i++ {
		xDist := node.X - msgs[i].Sender.X
		yDist := node.Y - msgs[i].Sender.Y
		dist := math.Sqrt(float64(xDist*xDist) + float64(yDist*yDist))
		distances = append(distances, dist)
	}

	//TODO efficient sorting?
	for i := 0; i < len(msgs); i++ {
		for j := 0; j < len(msgs)-i-1; j++ {
			if distances[j] > distances[j+1] {
				helloTemp := msgs[j]
				msgs[j] = msgs[j+1]
				msgs[j+1] = helloTemp

				distTemp := distances[j]
				distances[j] = distances[j+1]
				distances[j+1] = distTemp
			}
		}
	}
	k := 0
	for k < len(distances) && distances[k] <= 8 {
		k++
	}
	if k < len(distances) {
		node.NodeClusterParams.RecvMsgs = node.NodeClusterParams.RecvMsgs[:k]
	}
}*/

func (adhoc *AdHocNetwork) ElectClusterHead(curNode *NodeImpl, p *Params) {
	maxNode := curNode.HasMaxNodeScore(p)

	if !maxNode.IsClusterHead {
		maxNode.IsClusterHead = true
		maxNode.TimeBecameClusterHead = p.CurrentTime
		maxNode.BatteryBecameClusterHead = maxNode.GetBatteryPercentage()
		maxNode.IsClusterMember = false
		adhoc.ClusterHeads = append(adhoc.ClusterHeads, maxNode)
		maxNode.NodeClusterParams.CurrentCluster = &Cluster{maxNode, []*NodeImpl{}, adhoc.NextClusterNum}
		adhoc.NextClusterNum++
	}
	if curNode != maxNode {
		maxNode.NodeClusterParams.CurrentCluster.ClusterMembers = append(maxNode.NodeClusterParams.CurrentCluster.ClusterMembers, curNode)
		curNode.NodeClusterParams.CurrentCluster = maxNode.NodeClusterParams.CurrentCluster
		curNode.IsClusterMember = true
	}
}

//Assumed to be called by cluster heads
/*func (adhoc *AdHocNetwork) FormClusters(clusterHead *NodeImpl, p *Params) {

	msgs := clusterHead.NodeClusterParams.RecvMsgs
	if clusterHead.NodeClusterParams.CurrentCluster == nil {
		clusterHead.NodeClusterParams.CurrentCluster = &Cluster{clusterHead, []*NodeImpl{}, adhoc, adhoc.NextClusterNum}
		adhoc.NextClusterNum++
		adhoc.ClusterHeads = append(adhoc.ClusterHeads, clusterHead)
	}

	for i := 0; i < len(msgs) && len(clusterHead.NodeClusterParams.CurrentCluster.ClusterMembers) < p.ClusterMaxThreshold; i++ {
		if !msgs[i].Sender.IsClusterHead && !msgs[i].Sender.IsClusterMember {
			//if(clusterHead.IsWithinRange(clusterHead.NodeClusterParams.RecvMsgs[i].Sender,8)){
			clusterHead.NodeClusterParams.CurrentCluster.ClusterMembers = append(clusterHead.NodeClusterParams.CurrentCluster.ClusterMembers, msgs[i].Sender)

			msgs[i].Sender.IsClusterMember = true
			msgs[i].Sender.NodeClusterParams.CurrentCluster = clusterHead.NodeClusterParams.CurrentCluster

			clusterHead.DrainBatteryBluetooth()
			clusterHead.NodeClusterParams.RecvMsgs[i].Sender.DrainBatteryBluetooth()

			//}
		}
	}

	for i := 0; i < len(msgs); i++ {
		if msgs[i].Sender.IsClusterHead || msgs[i].Sender.IsClusterMember {
			//if i < len(msgs) {
			msgs = append(msgs[:i], msgs[i+1:]...)
			//}
		}
	}

	if len(msgs) > 0 {
		for i := 0; i < len(msgs); i++ {
			//if msgs[i].Sender.Battery > msgs[i].Sender.P.ThreshHoldBatteryToHave {

			adhoc.ElectClusterHead(msgs[i].Sender, p)

			//Add to SingularNodes if it has
			exists := false
			for j := 0; j < len(adhoc.SingularNodes); j++ {
				if adhoc.SingularNodes[j] == msgs[i].Sender {
					exists = true
					break
				}
			}

			if !exists {
				adhoc.SingularNodes = append(adhoc.SingularNodes, msgs[i].Sender)
			}
			//}
		}
	}
}*/

//sorts clusterheads by distance to the current node
/*func (adhoc *AdHocNetwork) SortClusterHeads(curNode *NodeImpl) (viableOptions []*NodeImpl) {

	distances := []float64{}
	viableOptions = []*NodeImpl{}

	for j := 0; j < len(adhoc.ClusterHeads); j++ {
		xDist := curNode.X - adhoc.ClusterHeads[j].X
		yDist := curNode.Y - adhoc.ClusterHeads[j].Y
		dist := math.Sqrt(float64(xDist*xDist) + float64(yDist*yDist))
		distances = append(distances, dist)
		viableOptions = append(viableOptions, adhoc.ClusterHeads[j])
	}

	for i := 0; i < len(viableOptions); i++ {
		for j := 0; j < len(viableOptions)-i-1; j++ {
			if distances[j] > distances[j+1] {
				chTemp := viableOptions[j]
				viableOptions[j] = viableOptions[j+1]
				viableOptions[j+1] = chTemp

				distTemp := distances[j]
				distances[j] = distances[j+1]
				distances[j+1] = distTemp
			}
		}
	}
	k := 0
	for k < len(distances) && distances[k] <= 8 {
		k++
	}
	if k < len(distances) {
		viableOptions = viableOptions[:k]
	}

	//fmt.Println(distances)
	return viableOptions
}*/

/*func (adhoc *AdHocNetwork) FinalizeClusters(p *Params) {
	//TODO clean this up. This code block should not be needed
	for i := 0; i < len(p.NodeList); i++ {
		//Nodes marked as members but not in a cluster added to SingularNodes
		if p.NodeList[i].IsClusterMember && !p.NodeList[i].IsClusterHead {
			if p.NodeList[i].NodeClusterParams.CurrentCluster == nil {
			} else {
				found := false
				j := 0
				for j < len(adhoc.ClusterHeads) {
					if adhoc.ClusterHeads[j] == p.NodeList[i].NodeClusterParams.CurrentCluster.ClusterHead {
						break
					} else {
						j++
					}
				}
				if j < len(adhoc.ClusterHeads) {
					for k := 0; k < len(adhoc.ClusterHeads[j].NodeClusterParams.CurrentCluster.ClusterMembers) && !found; k++ {
						if adhoc.ClusterHeads[j].NodeClusterParams.CurrentCluster.ClusterMembers[k] == p.NodeList[i] {
							found = true
						}
					}
				}
				if !found {
					p.NodeList[i].IsClusterMember = false
					adhoc.SingularNodes = append(adhoc.SingularNodes, p.NodeList[i])
				}
			}
		}
	}

	//Finds viable options for singular nodes (join another cluster or form own)
	for i := 0; i < len(adhoc.SingularNodes); i++ {
		if adhoc.SingularNodes[i].IsClusterHead {
			adhoc.FormClusters(adhoc.SingularNodes[i], p)
		} else if !adhoc.SingularNodes[i].IsClusterMember {
			viableOptions := adhoc.SortClusterHeads(adhoc.SingularNodes[i])

			k := 0
			joined := false
			atj := []*NodeImpl{}
			for !joined && k < len(viableOptions) {
				//fmt.Printf("\tViableOption Total: %d\n",viableOptions[0].NodeClusterParams.CurrentCluster.Total)
				if len(viableOptions[k].NodeClusterParams.CurrentCluster.ClusterMembers) < p.ClusterMaxThreshold {
					clusterHead := viableOptions[k]

					clusterHead.NodeClusterParams.CurrentCluster.ClusterMembers = append(clusterHead.NodeClusterParams.CurrentCluster.ClusterMembers, adhoc.SingularNodes[i])

					adhoc.SingularNodes[i].IsClusterMember = true
					adhoc.SingularNodes[i].NodeClusterParams.CurrentCluster = clusterHead.NodeClusterParams.CurrentCluster
					joined = true

					adhoc.SingularNodes[i].DrainBatteryBluetooth()
					clusterHead.DrainBatteryBluetooth()
				} else {
					atj = append(atj, viableOptions[k])
				}
				k++
			}
			if k == len(viableOptions) {
				adhoc.SingularNodes[i].NodeClusterParams.CurrentCluster = &Cluster{}
				adhoc.SingularNodes[i].NodeClusterParams.CurrentCluster.ClusterHead = adhoc.SingularNodes[i]
				adhoc.SingularNodes[i].NodeClusterParams.CurrentCluster.ClusterMembers = []*NodeImpl{}
				adhoc.SingularNodes[i].NodeClusterParams.CurrentCluster.ClusterNetwork = adhoc
				adhoc.SingularNodes[i].IsClusterHead = true
				adhoc.SingularNodes[i].IsClusterMember = false

				adhoc.ClusterHeads = append(adhoc.ClusterHeads, adhoc.SingularNodes[i])
				adhoc.SingularNodes[i].NodeClusterParams.AttemptedToJoin = append(adhoc.SingularNodes[i].NodeClusterParams.AttemptedToJoin, atj...)

				adhoc.SingularNodes[i].DrainBatteryBluetooth()
			}
		}
	}

	//If no longer a singular node remove from array and decrease count
	for i := 0; i < len(adhoc.SingularNodes); i++ {
		if adhoc.SingularNodes[i].IsClusterHead || adhoc.SingularNodes[i].IsClusterMember { // && adhoc.SingularNodes[i].NodeClusterParams.CurrentCluster!=nil && adhoc.SingularNodes[i].NodeClusterParams.CurrentCluster.Total>0){
			adhoc.SingularNodes = append(adhoc.SingularNodes[:i], adhoc.SingularNodes[i+1:]...)
		}
	}

	//if a cluster head is part of a cluster remove from that cluster
	for i := 0; i < len(adhoc.ClusterHeads); i++ {
		for j := 0; j < len(adhoc.ClusterHeads); j++ {
			for k := 0; k < len(adhoc.ClusterHeads[j].NodeClusterParams.CurrentCluster.ClusterMembers); k++ {
				if adhoc.ClusterHeads[j].NodeClusterParams.CurrentCluster.ClusterMembers[k] == adhoc.ClusterHeads[i] {
					adhoc.ClusterHeads[j].NodeClusterParams.CurrentCluster.ClusterMembers = append(adhoc.ClusterHeads[j].NodeClusterParams.CurrentCluster.ClusterMembers[:k], adhoc.ClusterHeads[j].NodeClusterParams.CurrentCluster.ClusterMembers[k+1:]...)
				}
			}
		}
	}
}*/

func (node *NodeImpl) IsWithinRange(node2 *NodeImpl, searchRange float64) bool {
	xDist := node.X - node2.X
	yDist := node.Y - node2.Y
	radDist := math.Sqrt(float64(xDist*xDist) + float64(yDist*yDist))

	return radDist <= searchRange
}

func (adhoc *AdHocNetwork) FullRecluster(p *Params) {
	adhoc.FullReclusters++
	adhoc.ResetClusters(p)
	for node := range p.AliveNodes {
		node.DrainBatteryWifi()	//Server sends message to all nodes that reclustering is happening
		if node.Valid {
			adhoc.SendHello(node, nil, p)
		}
	}
	for node := range p.AliveNodes {
		if !node.IsClusterHead && node.Valid {
			adhoc.ElectClusterHead(node, p)
		}
	}
}

func (adhoc *AdHocNetwork) LocalRecluster(head *NodeImpl, members []*NodeImpl, p *Params) {
	adhoc.LocalReclusters++
	adhoc.ClearClusterParams(head)
	head.DrainBatteryWifi() //Head informs the server it is dying
	for i := 0; i < len(members); i++ {
		if members[i].IsAlive() {
			if p.LocalRecluster == 1 {
				toJoin := members[i].FindNearbyHead(p, &p.Server.ReclusterBTCounter)
				if toJoin != nil {
					members[i].Join(toJoin.NodeClusterParams.CurrentCluster)
					members = append(members[:i], members[i+1:]...)
					i--
					adhoc.CSJoins++
					continue
				}
			}
			adhoc.SendHello(members[i], members, p)
		}
	}
	for _, node := range members {
		if !node.IsClusterHead && node.IsAlive() {
			adhoc.ElectClusterHead(node, p)
		}
	}
	//For logging, not function
	for _, node := range members {
		if node.IsClusterHead {
			adhoc.LRHeads++
		}
	}
}

func (adhoc *AdHocNetwork) ExpansiveLocalRecluster(head *NodeImpl, members []*NodeImpl, p *Params) {
	withinDist := p.NodeTree.WithinRadius(p.NodeBTRange, head, []*NodeImpl{})
	head.IsClusterHead = false
	head.DrainBatteryBluetooth(&p.Server.ReclusterBTCounter) //Sends message to all in bluetooth range
	for i := 0; i < len(withinDist); i++ {
		withinDist[i].DrainBatteryBluetooth(&p.Server.ReclusterBTCounter) //All nodes in range receive message
		if withinDist[i].IsClusterHead && (p.CurrentTime - withinDist[i].TimeBecameClusterHead) / 1000 > (p.LocalRecluster-3) {
			withinDist[i].DrainBatteryWifi() //This cluster head informs the server that it is within range of the dying cluster head
			adhoc.ClearClusterParams(withinDist[i])
			members = append(members, withinDist[i].NodeClusterParams.CurrentCluster.ClusterMembers...)
			if withinDist[i].IsAlive() {
				members = append(members, withinDist[i])
			}
		}
	}
	adhoc.LocalRecluster(head, members, p)
}

func (node *NodeImpl) UpdateOutOfRange(p *Params) {
	if !node.IsClusterHead {
		if node.NodeClusterParams.CurrentCluster.ClusterHead != nil {
			if !node.IsWithinRange(node.NodeClusterParams.CurrentCluster.ClusterHead, p.NodeBTRange) {
				if !node.OutOfRange {
					node.OutOfRange = true
					node.TimeMovedOutOfRange = p.CurrentTime
				}
			} else {
				node.OutOfRange = false
			}
		}
	} else {
		for _, member := range node.NodeClusterParams.CurrentCluster.ClusterMembers {
			if !node.IsWithinRange(member, p.NodeBTRange) {
				if !member.OutOfRange {
					member.OutOfRange = true
					member.TimeMovedOutOfRange = p.CurrentTime
				}
			} else {
				member.OutOfRange = false
			}
		}
	}
}

func SimpleIntersect(a []*NodeImpl, b []*NodeImpl) []*NodeImpl {
	result := make([]*NodeImpl, 0)

	for i := 0; i < len(a); i++ {
		for j := 0; j < len(b); j++ {
			if a[i] == b[j] {
				result = append(result, a[i])
			}
		}
	}

	return result
}

func (node *NodeImpl) FindNearbyHead(p *Params, counter *int) *NodeImpl {
	node.DrainBatteryBluetooth(counter)	//Broadcasting hello message

	withinDist := p.NodeTree.WithinRadius(p.NodeBTRange, node, []*NodeImpl{})

	var toJoin *NodeImpl = nil
	highestBattery := p.BatteryDeadThreshold
	for i := 0; i < len(withinDist); i++ {
		withinDist[i].DrainBatteryBluetooth(counter)	//Every node in bluetooth range receives the hello message
		if withinDist[i].IsClusterHead && len(withinDist[i].NodeClusterParams.CurrentCluster.ClusterMembers) < p.ClusterMaxThreshold && withinDist[i].IsAlive() {
			withinDist[i].DrainBatteryBluetooth(counter) //Every cluster head with room for this node sends a reply
			if node.IsAlive() {
				node.DrainBatteryBluetooth(counter) //The node receives messages from every cluster head with room (unless it died)
			}
			if withinDist[i].GetBatteryPercentage() > highestBattery {
				toJoin = withinDist[i]
			}
		}
	}
	return toJoin
}

func (node *NodeImpl) ShouldLocalRecluster() bool {
	return node.P.CurrentTime - node.TimeBecameClusterHead > node.P.ClusterHeadTimeThreshold ||
		node.BatteryBecameClusterHead - node.GetBatteryPercentage() > node.P.ClusterHeadBatteryDropThreshold
}

func (node *NodeImpl) InitLocalRecluster() {
	members := make([]*NodeImpl, len(node.NodeClusterParams.CurrentCluster.ClusterMembers))
	copy(members, node.NodeClusterParams.CurrentCluster.ClusterMembers)
	if node.P.LocalRecluster < 3 {
		node.P.ClusterNetwork.LocalRecluster(node, members, node.P)
	} else {
		node.P.ClusterNetwork.ExpansiveLocalRecluster(node, members, node.P)
	}
}