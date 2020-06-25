package cps

import (
	"math"
)

type AdHocNetwork struct {
	ClusterHeads  []*NodeImpl
	SingularNodes []*NodeImpl
	TotalMsgs      int //used to counts total messages sent/received in one iteration
	FullReclusters int //Counts the number of full reclusters that occur in a simulation
	PotentialReclusters int //Counts the number of reclusters that would have occured if there was a check every second
	NextClusterNum int //For testing, may remove later
}

type Cluster struct {
	ClusterHead    *NodeImpl
	ClusterMembers []*NodeImpl
	ClusterNetwork *AdHocNetwork
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
	battery := node.GetBatteryPercentage() * float64(p.ClusterMaxThreshold) //Multiplying by threshold ensures that battery and degree have the same maximum value

	//weighted sum, 60% from degree (# of nodes within distance), 40% from its battery life
	// penalty used to increase a nodes chance at staying a clusterhead
	score := p.DegreeWeight*degree + p.BatteryWeight*battery
	if node.IsClusterHead {
		return score
	} else {
		return score * p.Penalty
	}
}

//Generates Hello Message for node to form/maintain clusters. Returns message as a string
func (node *NodeImpl) GenerateHello(score float64) {
	message := &HelloMsg{ Sender: node, NodeCHScore: score}
	node.NodeClusterParams.ThisNodeHello = message
}

func (adhoc *AdHocNetwork) SendHelloMessage(curNode *NodeImpl, p *Params) {
	withinDist := p.NodeTree.WithinRadius(p.NodeBTRange, curNode, []*NodeImpl{})
	numWithinDist := len(withinDist)

	curNode.GenerateHello(curNode.ComputeClusterScore(p, numWithinDist))
	curNode.DrainBatteryBluetooth()	//Broadcasting first hello message, no score
	curNode.DrainBatteryBluetooth()	//Broadcasting second hello message with score

	//var buffer bytes.Buffer
	for j := 0; j < numWithinDist; j++ {
		withinDist[j].NodeClusterParams.RecvMsgs = append(withinDist[j].NodeClusterParams.RecvMsgs, curNode.NodeClusterParams.ThisNodeHello)
		//buffer.WriteString(fmt.Sprintf("SenderId=%v\tRecieverId=%v\tSenderCHS=%v\n",curNode.Id,withinDist[j].CurNode.Id,curNode.NodeClusterParams.ThisNodeHello.NodeCHScore))
		withinDist[j].DrainBatteryBluetooth()	//Every node in bluetooth range receives the first hello message, allowing them to count how many neighbors they have
		withinDist[j].DrainBatteryBluetooth()	//Every node in bluetooth range then receives a second hello message, the one including curNode's score
		adhoc.TotalMsgs++
	}
	//fmt.Fprintf(curNode.P.ClusterMessages,buffer.String())
}

func (adhoc *AdHocNetwork) ClusterMovement(node *NodeImpl, p *Params) {
	if node.Valid {
		if !node.Alive {
			node.CurTree.RemoveAndClean(node)
			if node.IsClusterHead {
				adhoc.DissolveCluster(node)
			} else {
				adhoc.ClearClusterParams(node)
			}
		} else if (!node.IsClusterHead && (node.NodeClusterParams.CurrentCluster.ClusterHead == nil || !node.IsWithinRange(node.NodeClusterParams.CurrentCluster.ClusterHead, p.NodeBTRange))) || len(node.NodeClusterParams.CurrentCluster.ClusterMembers) <= 0 {
			adhoc.NewNodeHello(node, p)
		}
	}
}

func (adhoc *AdHocNetwork) NewNodeHello(node *NodeImpl, p *Params) {
	withinDist := p.NodeTree.WithinRadius(p.NodeBTRange, node, []*NodeImpl{})
	numWithinDist := len(withinDist)

	node.DrainBatteryBluetooth()	//Broadcasting hello message

	for i := 0; i < numWithinDist; i++ {
		withinDist[i].DrainBatteryBluetooth()	//Every node in bluetooth range receives the hello message
		if withinDist[i].IsClusterHead && len(withinDist[i].NodeClusterParams.CurrentCluster.ClusterMembers) < p.ClusterMaxThreshold {
			withinDist[i].DrainBatteryBluetooth() //Every cluster head with room for this node sends a reply
		}
	}

	for i := 0; i < numWithinDist; i++ {
		if withinDist[i].IsClusterHead && len(withinDist[i].NodeClusterParams.CurrentCluster.ClusterMembers) < p.ClusterMaxThreshold {
			//withinDist[j].GenerateHello(withinDist[j].ComputeClusterScore(p, len(p.NodeTree.WithinRadius(p.NodeBTRange, withinDist[j], []*NodeImpl{}))))
			adhoc.ClearClusterParams(node)
			node.Join(withinDist[i].NodeClusterParams.CurrentCluster)
			//adhoc.Joins++
			return
		}
	}

	adhoc.ClearClusterParams(node)
	adhoc.ElectClusterHead(node, p)
	//fmt.Fprintf(curNode.P.ClusterMessages,buffer.String())
}

func (node *NodeImpl) HasMaxNodeScore(p *Params) *NodeImpl {
	maxNode := node
	maxScore := node.NodeClusterParams.ThisNodeHello.NodeCHScore
	for i := 0; i < len(node.NodeClusterParams.RecvMsgs); i++ {
		sender := node.NodeClusterParams.RecvMsgs[i].Sender
		//do not consider nodes already with a clusterhead
		//if received a message from a node who does not have a cluster head
		if !sender.IsClusterMember && len(sender.NodeClusterParams.CurrentCluster.ClusterMembers) < p.ClusterMaxThreshold {
			//if their score higher than current node score
			if node.NodeClusterParams.RecvMsgs[i].NodeCHScore > maxScore {
				maxScore = node.NodeClusterParams.RecvMsgs[i].NodeCHScore
				maxNode = sender
			} else if math.Abs(node.NodeClusterParams.RecvMsgs[i].NodeCHScore-maxScore) < 0.1 {
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
	fmt.Print(node.Battery)
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
	for _, member := range node.NodeClusterParams.CurrentCluster.ClusterMembers {
		adhoc.ClearClusterParams(member)
	}
	adhoc.ClearClusterParams(node)
}

func (adhoc *AdHocNetwork) ResetClusters(p *Params) {
	for i := 0; i < len(p.AliveList); i++ {
		p.AliveList[i].IsClusterHead = false
		adhoc.ClearClusterParams(p.AliveList[i])
		if !p.AliveList[i].Alive {
			p.AliveList = append(p.AliveList[:i], p.AliveList[i+1:]...)
			i--
		}
	}
	adhoc.ClusterHeads = []*NodeImpl{}
	adhoc.SingularNodes = []*NodeImpl{}
	adhoc.TotalMsgs = 0
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
		maxNode.IsClusterMember = false
		adhoc.ClusterHeads = append(adhoc.ClusterHeads, maxNode)
		maxNode.NodeClusterParams.CurrentCluster = &Cluster{maxNode, []*NodeImpl{}, adhoc, adhoc.NextClusterNum}
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
	for _, node := range p.AliveList {
		node.DrainBatteryWifi()	//Server sends message to all nodes that reclustering is happening
		if node.Valid {
			adhoc.SendHelloMessage(node, p)
		}
	}
	for _, node := range p.AliveList {
		if !node.IsClusterHead && node.Valid {
			adhoc.ElectClusterHead(node, p)
		}
	}
}
