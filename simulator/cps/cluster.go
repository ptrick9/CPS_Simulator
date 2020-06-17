package cps

import (
	"fmt"
	"math"
)

type AdHocNetwork struct {
	ClusterHeads	[]*NodeImpl
	SingularNodes   []*NodeImpl
	TotalHeads		int
	SingularCount	int
	//Threshold		int //maximum # of nodes in a cluster
	TotalMsgs		int //used to counts total messages sent/received in one iteration
	NextClusterNum	int
}

type Cluster struct {
	ClusterHead			*NodeImpl//*NodeImpl	//id of clusterhead
	Total				int //current # of nodes in a cluster
	ClusterMembers		[]*NodeImpl
	ClusterNetwork		*AdHocNetwork
	ClusterNum			int //integer identifier of node
}

type ClusterMemberParams struct{
	CurrentCluster		*Cluster
	RecvMsgs			[]*HelloMsg
	ThisNodeHello		*HelloMsg
	AttemptedToJoin		[]*NodeImpl
}

type HelloMsg struct {
	Sender				*NodeImpl	//pointer to Node sending the Hello Msg
	ClusterHead			*NodeImpl	//pointer to of the cluster head
	//nil if not in a cluster
	//points to self if Node is a ClusterHead
	NodeCHScore			float64	//score for determining how suitable a node is to be a clusterhead
	Option				int		//0 for regular node, if a cluster head this is the # of nodes in the cluster
	BrodPeriod			float64	//broadcast period of the Sender
}

type ClusterNode struct{ //made for testing, only has parameters that the cluster needs to know from a NodeImpl
	NodeBounds			*Bounds
	Battery				float32
	ClusterHead			*NodeImpl
	NodeClusterParams 	*ClusterMemberParams
	IsClusterHead		bool
	IsClusterMember		bool
}

//Computes the cluster score (higher the score the better chance a node becomes a cluster head)
func (node * NodeImpl) ComputeClusterScore(p *Params, numWithinDist int,) float64 {
	degree := math.Min(float64(numWithinDist), float64(p.ClusterThreshold))
	battery := float64(node.Battery) * float64(p.ClusterThreshold)/100 //Multiplying by threshold/100 ensures that battery and degree have the same maximum value

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
func (node * NodeImpl)GenerateHello(score float64) {
	var option int

	//if(curNode.IsClusterHead){
	//	option = curNode.NodeClusterParams.CurrentCluster.Total
	//} else{
	option = 0
	//}

	message := &HelloMsg{
		Sender:      node,
		NodeCHScore: score,
		Option:      option,
		BrodPeriod:  0.2}
	node.NodeClusterParams.ThisNodeHello = message
}

func (adhoc * AdHocNetwork) SendHelloMessage(curNode * NodeImpl, p *Params){
	withinDist := p.NodeTree.WithinRadius(p.NodeBTRange, curNode, []*NodeImpl{})
	numWithinDist := len(withinDist)

	curNode.GenerateHello(curNode.ComputeClusterScore(p,numWithinDist))

	//var buffer bytes.Buffer
	for j:=0; j<numWithinDist; j++ {
		curClusterP := withinDist[j].NodeClusterParams
		//if curClusterP.RecvMsgs != nil {
		//	if withinDist[j].Battery > withinDist[j].P.ThreshHoldBatteryToHave {
				//if curClusterP.ThisNodeHello.Sender != nil {
					curClusterP.RecvMsgs = append(curClusterP.RecvMsgs, curNode.NodeClusterParams.ThisNodeHello)
					//buffer.WriteString(fmt.Sprintf("SenderId=%v\tRecieverId=%v\tSenderCHS=%v\n",curNode.Id,withinDist[j].CurNode.Id,curNode.NodeClusterParams.ThisNodeHello.NodeCHScore))
					//curNode.DecrementPowerBT()
					adhoc.TotalMsgs++
				//}
			//}
		//}
	}
	//fmt.Fprintf(curNode.P.ClusterMessages,buffer.String())
}

func (node * NodeImpl) HasMaxNodeScore(p * Params)  *NodeImpl {
	maxNode := node //&(NodeImpl{})
	maxScore := node.NodeClusterParams.ThisNodeHello.NodeCHScore
	//print("cur node: ")
	//println(node.Id)
	for i:= 0; i<len(node.NodeClusterParams.RecvMsgs); i++{
		sender := node.NodeClusterParams.RecvMsgs[i].Sender
		//print("within: ")
		//print(node.NodeClusterParams.RecvMsgs[i].Sender.Id)
		//print("   score: ")
		//println(node.NodeClusterParams.RecvMsgs[i].Sender.NodeClusterParams.ThisNodeHello.NodeCHScore)
		//do not consider nodes already with a clusterhead
		//if received a message from a node who does not have a cluster head
		if !sender.IsClusterMember && len(sender.NodeClusterParams.CurrentCluster.ClusterMembers) < p.ClusterThreshold {
			//if their score higher than current node score
			if node.NodeClusterParams.RecvMsgs[i].NodeCHScore > maxScore {
				maxScore = node.NodeClusterParams.RecvMsgs[i].NodeCHScore
				maxNode = sender
			} else if math.Abs(node.NodeClusterParams.RecvMsgs[i].NodeCHScore - maxScore) < 0.1 {
				if sender.Id < maxNode.Id {
					maxNode = sender
				}
			}
		}
	}

	//print("max node: ")
	//println(maxNode.Id)
	//println()

	return maxNode
}

func (node * NodeImpl) PrintClusterNode(){
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
}

func (curNode * NodeImpl) ClearClusterParams(){
	//Reset Cluster Params (all but hello->sender since that will stay the same always)
	curNode.NodeClusterParams.CurrentCluster = &Cluster{}
	curNode.NodeClusterParams.CurrentCluster.ClusterMembers = []*NodeImpl{}

	curNode.NodeClusterParams.RecvMsgs = []*HelloMsg{}
	curNode.NodeClusterParams.ThisNodeHello = &HelloMsg{}

	curNode.IsClusterMember = false
	curNode.IsClusterHead = false
}

func (adhoc * AdHocNetwork) ResetClusters(p *Params){
	for i:=0; i<len(p.NodeList); i++{
		p.NodeList[i].ClearClusterParams()
	}
	adhoc.TotalHeads = 0
	adhoc.ClusterHeads = []*NodeImpl{}
	adhoc.SingularNodes = []*NodeImpl{}
	adhoc.TotalMsgs = 0
}

//sorts messages by distance to the node: 0th = closest, nth = farthest
func (node * NodeImpl) SortMessages(){

	distances := []float64{}

	msgs := node.NodeClusterParams.RecvMsgs

	for i :=0; i <len(msgs); i++{
		xDist := node.X - msgs[i].Sender.X
		yDist := node.Y - msgs[i].Sender.Y
		dist := math.Sqrt(float64(xDist*xDist) + float64(yDist*yDist))
		distances = append(distances, dist)
	}

	//TODO efficient sorting?
	for i:=0; i<len(msgs); i++{
		for j:=0; j<len(msgs)-i-1; j++{
			if distances[j]>distances[j+1] {
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
	for k<len(distances) && distances[k] <= 8 {
		k++
	}
	if k<len(distances) {
		node.NodeClusterParams.RecvMsgs  = node.NodeClusterParams.RecvMsgs[:k]
	}
}

func (adhoc * AdHocNetwork) ElectClusterHead(curNode * NodeImpl, p *Params){
	maxNode := curNode.HasMaxNodeScore(p)

	if !maxNode.IsClusterHead {
		maxNode.IsClusterHead = true
		maxNode.IsClusterMember = false
		adhoc.ClusterHeads = append(adhoc.ClusterHeads, maxNode)
		adhoc.TotalHeads = len(adhoc.ClusterHeads)
		maxNode.NodeClusterParams.CurrentCluster = &Cluster{maxNode,0,[]*NodeImpl{},adhoc, adhoc.NextClusterNum}
		adhoc.NextClusterNum++
	}
	maxNode.NodeClusterParams.CurrentCluster.ClusterMembers = append(maxNode.NodeClusterParams.CurrentCluster.ClusterMembers, curNode)
	curNode.NodeClusterParams.CurrentCluster = maxNode.NodeClusterParams.CurrentCluster
	curNode.IsClusterMember = true
}

//Assumed to be called by cluster heads
func (adhoc * AdHocNetwork) FormClusters(clusterHead * NodeImpl, p *Params){

	msgs := clusterHead.NodeClusterParams.RecvMsgs
	if clusterHead.NodeClusterParams.CurrentCluster==nil {
		clusterHead.NodeClusterParams.CurrentCluster = &Cluster{clusterHead, 0, []*NodeImpl{},adhoc, adhoc.NextClusterNum}
		adhoc.NextClusterNum++
		adhoc.ClusterHeads = append(adhoc.ClusterHeads, clusterHead)
		adhoc.TotalHeads++
	}

	for i:=0; i<len(msgs) && clusterHead.NodeClusterParams.CurrentCluster.Total < p.ClusterThreshold; i++{
		if !msgs[i].Sender.IsClusterHead && !msgs[i].Sender.IsClusterMember {
			//if(clusterHead.IsWithinRange(clusterHead.NodeClusterParams.RecvMsgs[i].Sender,8)){
			clusterHead.NodeClusterParams.CurrentCluster.ClusterMembers = append(clusterHead.NodeClusterParams.CurrentCluster.ClusterMembers, msgs[i].Sender)
			clusterHead.NodeClusterParams.CurrentCluster.Total = len(clusterHead.NodeClusterParams.CurrentCluster.ClusterMembers)

			msgs[i].Sender.IsClusterMember = true
			msgs[i].Sender.NodeClusterParams.CurrentCluster = clusterHead.NodeClusterParams.CurrentCluster

			//clusterHead.DecrementPowerBT()
			//clusterHead.NodeClusterParams.RecvMsgs[i].Sender.DecrementPowerBT()

			//}
		}
	}

	for i:=0; i<len(msgs); i++ {
		if msgs[i].Sender.IsClusterHead || msgs[i].Sender.IsClusterMember {
			//if i < len(msgs) {
			msgs = append(msgs[:i], msgs[i+1:]...)
			//}
		}
	}

	if len(msgs)>0 {
		for i:=0; i<len(msgs); i++{
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
					adhoc.SingularCount++
				}
			//}
		}
	}
}

//sorts clusterheads by distance to the current node
func (adhoc * AdHocNetwork) SortClusterHeads(curNode * NodeImpl) (viableOptions []*NodeImpl){

	distances := []float64{}
	viableOptions = []*NodeImpl{}

	for j:=0; j<len(adhoc.ClusterHeads);j++{
		xDist := curNode.X-adhoc.ClusterHeads[j].X
		yDist := curNode.Y-adhoc.ClusterHeads[j].Y
		dist := math.Sqrt(float64(xDist*xDist)+float64(yDist*yDist))
		distances = append(distances, dist)
		viableOptions = append(viableOptions, adhoc.ClusterHeads[j])
	}

	for i:=0; i<len(viableOptions);i++{
		for j:=0; j<len(viableOptions)-i-1; j++{
			if distances[j]>distances[j+1] {
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
	for k<len(distances) && distances[k] <= 8 {
		k++
	}
	if k<len(distances) {
		viableOptions  = viableOptions[:k]
	}

	//fmt.Println(distances)
	return viableOptions
}

func (adhoc * AdHocNetwork) FinalizeClusters(p * Params){
	//TODO clean this up. This code block should not be needed
	for i:=0; i<len(p.NodeList); i++ {
		//Nodes marked as members but not in a cluster added to SingularNodes
		if p.NodeList[i].IsClusterMember && !p.NodeList[i].IsClusterHead {
			if p.NodeList[i].NodeClusterParams.CurrentCluster==nil {
			} else{
				found := false
				j := 0
				for j<len(adhoc.ClusterHeads){
					if adhoc.ClusterHeads[j] ==  p.NodeList[i].NodeClusterParams.CurrentCluster.ClusterHead {
						break
					} else{
						j++
					}
				}
				if j<len(adhoc.ClusterHeads) {
					for k := 0; k<len(adhoc.ClusterHeads[j].NodeClusterParams.CurrentCluster.ClusterMembers) && !found; k++{
						if adhoc.ClusterHeads[j].NodeClusterParams.CurrentCluster.ClusterMembers[k]==p.NodeList[i] {
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
	for i:=0; i<len(adhoc.SingularNodes); i++{
		if adhoc.SingularNodes[i].IsClusterHead {
			adhoc.FormClusters(adhoc.SingularNodes[i], p)
		} else if !adhoc.SingularNodes[i].IsClusterMember {
			viableOptions := adhoc.SortClusterHeads(adhoc.SingularNodes[i])

			k := 0
			joined := false
			atj := []*NodeImpl{}
			for !joined && k<len(viableOptions){
				//fmt.Printf("\tViableOption Total: %d\n",viableOptions[0].NodeClusterParams.CurrentCluster.Total)
				if viableOptions[k].NodeClusterParams.CurrentCluster.Total < p.ClusterThreshold {
					clusterHead := viableOptions[k]

					clusterHead.NodeClusterParams.CurrentCluster.ClusterMembers = append(clusterHead.NodeClusterParams.CurrentCluster.ClusterMembers, adhoc.SingularNodes[i])
					clusterHead.NodeClusterParams.CurrentCluster.Total++

					adhoc.SingularNodes[i].IsClusterMember = true
					adhoc.SingularNodes[i].NodeClusterParams.CurrentCluster = clusterHead.NodeClusterParams.CurrentCluster
					joined = true

					//adhoc.SingularNodes[i].DecrementPowerBT()
					//clusterHead.DecrementPowerBT()
				} else{
					atj = append(atj, viableOptions[k])
				}
				k++
			}
			if k==len(viableOptions) {
				adhoc.SingularNodes[i].NodeClusterParams.CurrentCluster = &Cluster{}
				adhoc.SingularNodes[i].NodeClusterParams.CurrentCluster.ClusterHead = adhoc.SingularNodes[i]
				adhoc.SingularNodes[i].NodeClusterParams.CurrentCluster.Total = 0
				adhoc.SingularNodes[i].NodeClusterParams.CurrentCluster.ClusterMembers = []*NodeImpl{}
				adhoc.SingularNodes[i].NodeClusterParams.CurrentCluster.ClusterNetwork = adhoc
				adhoc.SingularNodes[i].IsClusterHead = true
				adhoc.SingularNodes[i].IsClusterMember = false

				adhoc.ClusterHeads = append(adhoc.ClusterHeads, adhoc.SingularNodes[i])
				adhoc.TotalHeads++
				adhoc.SingularNodes[i].NodeClusterParams.AttemptedToJoin = append(adhoc.SingularNodes[i].NodeClusterParams.AttemptedToJoin, atj...)

				//adhoc.SingularNodes[i].DecrementPowerBT()
			}
		}
	}

	//If no longer a singular node remove from array and decrease count
	for i:=0; i<len(adhoc.SingularNodes); i++{
		if adhoc.SingularNodes[i].IsClusterHead || adhoc.SingularNodes[i].IsClusterMember { // && adhoc.SingularNodes[i].NodeClusterParams.CurrentCluster!=nil && adhoc.SingularNodes[i].NodeClusterParams.CurrentCluster.Total>0){
			adhoc.SingularNodes = append(adhoc.SingularNodes[:i],adhoc.SingularNodes[i+1:]...)
			adhoc.SingularCount--
		}
	}

	//if a cluster head is part of a cluster remove from that cluster
	for i:=0; i<len(adhoc.ClusterHeads); i++ {
		for j:=0; j<len(adhoc.ClusterHeads); j++ {
			for k:=0; k<len(adhoc.ClusterHeads[j].NodeClusterParams.CurrentCluster.ClusterMembers); k++ {
				if adhoc.ClusterHeads[j].NodeClusterParams.CurrentCluster.ClusterMembers[k]==adhoc.ClusterHeads[i] {
					adhoc.ClusterHeads[j].NodeClusterParams.CurrentCluster.ClusterMembers = append(adhoc.ClusterHeads[j].NodeClusterParams.CurrentCluster.ClusterMembers[:k],adhoc.ClusterHeads[j].NodeClusterParams.CurrentCluster.ClusterMembers[k+1:]... )
					adhoc.ClusterHeads[j].NodeClusterParams.CurrentCluster.Total--
				}
			}
		}
	}
}

func (node * NodeImpl) IsWithinRange(node2 * NodeImpl, searchRange float64) (inRange bool){
	xDist := node.X - node2.X
	yDist := node.Y - node2.Y
	radDist := math.Sqrt(float64(xDist*xDist)+float64(yDist*yDist))

	inRange = radDist <= searchRange

	return inRange
}
