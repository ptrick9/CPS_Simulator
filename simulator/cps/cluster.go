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
	Threshold		int //maximum # of nodes in a cluster
	TotalMsgs		int //used to counts total messages sent/received in one iteration
}

type Cluster struct {
	ClusterHead			*NodeImpl//*NodeImpl	//id of clusterhead
	Total				int //current # of nodes in a cluster
	ClusterMembers		[]*NodeImpl
	ClusterNetwork		*AdHocNetwork
}

type ClusterMemberParams struct{
	CurrentCluster		*Cluster
	RecvMsgs			[]*HelloMsg
	ThisNodeHello		*HelloMsg
	AttemptedToJoin		[]*NodeImpl
}

type HelloMsg struct {
	Sender				*NodeImpl//*NodeImpl		//pointer to Node sending the Hello Msg
	ClusterHead			*NodeImpl//NodeImpl		//pointer to of the cluster head
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

//Computes the cluster score (higher the score the better chance a node beccomes a cluster head)
func (node * NodeImpl) ComputeClusterScore(penalty float64, numWithinDist int) (score float64){
	degree := float64(numWithinDist)
	battery := float64(node.Battery)

	//weighted sum, 60% from degree (# of nodes within distance), 40% from its battery life
	// penalty used to increase a nodes chance at staying a clusterhead
	return (0.6*degree + 0.4*battery)* node.CHPenalty
}

//Generates Hello Message for node to form/maintain clusters. Returns message as a string
func (node * NodeImpl)GenerateHello(searchRange float64, score float64) {
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

func (adhoc * AdHocNetwork) SendHelloMessage(transmitRange float64, curNode * NodeImpl, nodeTree *Quadtree){
	withinDist := []*NodeImpl{}
	withinDist = nodeTree.WithinRadius(transmitRange, curNode, GetSearchBounds(curNode, transmitRange), withinDist)
	numWithinDist := len(withinDist)

	curNode.GenerateHello(transmitRange, curNode.ComputeClusterScore( 1,numWithinDist))

	//var buffer bytes.Buffer
	for j:=0; j<len(withinDist); j++ {
		curClusterP := withinDist[j].NodeClusterParams
		if curClusterP.RecvMsgs != nil {
			if withinDist[j].Battery > withinDist[j].P.ThreshHoldBatteryToHave {
				curClusterP.RecvMsgs = append(curClusterP.RecvMsgs, curClusterP.ThisNodeHello)
				//buffer.WriteString(fmt.Sprintf("SenderId=%v\tRecieverId=%v\tSenderCHS=%v\n",curNode.Id,withinDist[j].CurNode.Id,curNode.NodeClusterParams.ThisNodeHello.NodeCHScore))
				//curNode.DecrementPowerBT()
				adhoc.TotalMsgs++
			}
		}
	}
	//fmt.Fprintf(curNode.P.ClusterMessages,buffer.String())
}

func (node * NodeImpl) HasMaxNodeScore() (maxNode * NodeImpl){
	maxNode = node //&(NodeImpl{})
	maxScore := node.NodeClusterParams.ThisNodeHello.NodeCHScore
	for i:= 0; i<len(node.NodeClusterParams.RecvMsgs); i++{
		//do not consider nodes already with a clusterhead
		//if received a message from a node who does not have a cluster head
		if !node.NodeClusterParams.RecvMsgs[i].Sender.IsClusterMember {
			//if their score higher than current node score
			if node.NodeClusterParams.RecvMsgs[i].NodeCHScore > maxScore {
				maxScore = node.NodeClusterParams.RecvMsgs[i].NodeCHScore
				maxNode = node.NodeClusterParams.RecvMsgs[i].Sender
			}
		}
	}

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

func (adhoc * AdHocNetwork) ClearClusterParams(curNode * NodeImpl){
	//Reset Cluster Params (all but hello->sender since that will stay the same always)
	if curNode.NodeClusterParams.CurrentCluster!=nil {
		curNode.NodeClusterParams.CurrentCluster.ClusterHead = nil
		curNode.NodeClusterParams.CurrentCluster.Total = 0
		curNode.NodeClusterParams.CurrentCluster.ClusterMembers = []*NodeImpl{}
	} else{
		curNode.NodeClusterParams.CurrentCluster = &Cluster{}
	}

	curNode.NodeClusterParams.RecvMsgs = []*HelloMsg{}
	if curNode.NodeClusterParams.ThisNodeHello != nil {
		curNode.NodeClusterParams.ThisNodeHello.ClusterHead = nil
		curNode.NodeClusterParams.ThisNodeHello.NodeCHScore = 0
		curNode.NodeClusterParams.ThisNodeHello.BrodPeriod = 0
		curNode.NodeClusterParams.ThisNodeHello.Option = 0
	}else{
		curNode.NodeClusterParams.ThisNodeHello = &HelloMsg{}
	}

	if curNode.IsClusterHead {
		curNode.CHPenalty = 0.85
	}else{
		curNode.CHPenalty = 1.00
	}

	curNode.IsClusterMember = false
	curNode.IsClusterHead = false

	adhoc.ClusterHeads = []*NodeImpl{}
	adhoc.TotalHeads = 0
}

func (adhoc * AdHocNetwork) ResetClusters(){
	for i:=0; i<len(adhoc.ClusterHeads); i++{
		if adhoc.ClusterHeads[i].NodeClusterParams!=nil {
			adhoc.ClusterHeads[i].NodeClusterParams.AttemptedToJoin = []*NodeImpl{}
			if adhoc.ClusterHeads[i].NodeClusterParams.CurrentCluster!=nil {
				for j:=0; j<len(adhoc.ClusterHeads[i].NodeClusterParams.CurrentCluster.ClusterMembers); j++{
					adhoc.ClusterHeads[i].NodeClusterParams.CurrentCluster.ClusterMembers[j].IsClusterMember = false
					adhoc.ClusterHeads[i].NodeClusterParams.CurrentCluster.ClusterMembers[j].NodeClusterParams.CurrentCluster = nil
					adhoc.ClusterHeads[i].NodeClusterParams.CurrentCluster.ClusterMembers[j].NodeClusterParams.RecvMsgs = []*HelloMsg{}
					adhoc.ClusterHeads[i].NodeClusterParams.CurrentCluster.ClusterMembers[j].NodeClusterParams.AttemptedToJoin = []*NodeImpl{}
					adhoc.ClusterHeads[i].NodeClusterParams.CurrentCluster.ClusterMembers[j].NodeClusterParams.ThisNodeHello = nil
				}
				adhoc.ClusterHeads[i].NodeClusterParams.CurrentCluster = nil
				adhoc.ClusterHeads[i].IsClusterHead = false
			}
		}
	}
	adhoc.TotalHeads = 0
	adhoc.ClusterHeads = []*NodeImpl{}
	adhoc.SingularNodes = []*NodeImpl{}
	adhoc.TotalMsgs = 0
}

//sorts messages by distance to the node: 0th = closest, nth = farthest
func (node * NodeImpl) SortMessages(){

	distances := []float64{}

	for j:=0; j<len(node.NodeClusterParams.RecvMsgs);j++{
		xDist := node.X- node.NodeClusterParams.RecvMsgs[j].Sender.X
		yDist := node.Y- node.NodeClusterParams.RecvMsgs[j].Sender.Y
		dist := math.Sqrt(float64(xDist*xDist)+float64(yDist*yDist))
		distances = append(distances, dist)
	}

	//TODO efficient sorting?
	for i:=0; i<len(node.NodeClusterParams.RecvMsgs);i++{
		for j:=0; j<len(node.NodeClusterParams.RecvMsgs)-i-1; j++{
			if distances[j]>distances[j+1] {
				helloTemp := node.NodeClusterParams.RecvMsgs[j]
				node.NodeClusterParams.RecvMsgs[j] = node.NodeClusterParams.RecvMsgs[j+1]
				node.NodeClusterParams.RecvMsgs[j+1] = helloTemp

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

func (adhoc * AdHocNetwork) ElectClusterHead(curNode * NodeImpl){
	maxNode := curNode.HasMaxNodeScore()

	if maxNode.IsClusterHead == false {
		maxNode.IsClusterHead = true
		maxNode.IsClusterMember = false
		adhoc.ClusterHeads = append(adhoc.ClusterHeads, maxNode)
		adhoc.TotalHeads = len(adhoc.ClusterHeads)
		maxNode.NodeClusterParams.CurrentCluster = &Cluster{maxNode,0,[]*NodeImpl{},adhoc}
	}
}

//Assumed to be called by cluster heads
func (adhoc * AdHocNetwork) FormClusters(clusterHead * NodeImpl){

	msgs := clusterHead.NodeClusterParams.RecvMsgs
	if clusterHead.NodeClusterParams.CurrentCluster==nil {
		clusterHead.NodeClusterParams.CurrentCluster = &Cluster{clusterHead, 0, []*NodeImpl{},adhoc }
		adhoc.ClusterHeads = append(adhoc.ClusterHeads, clusterHead)
		adhoc.TotalHeads++
	}

	for i:=0; i<len(msgs) && clusterHead.NodeClusterParams.CurrentCluster.Total<adhoc.Threshold; i++{
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
			if msgs[i].Sender.Battery > msgs[i].Sender.P.ThreshHoldBatteryToHave {

				adhoc.ElectClusterHead(msgs[i].Sender)

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
			}
		}
	}
}

//sorts clusterheads by distance to the current node
func (adhoc * AdHocNetwork) SortClusterHeads(curNode * NodeImpl, searchRange float64) (viableOptions []*NodeImpl){

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
			adhoc.FormClusters(adhoc.SingularNodes[i])
		} else if !adhoc.SingularNodes[i].IsClusterMember {
			viableOptions := adhoc.SortClusterHeads(adhoc.SingularNodes[i],p.NodeBTRange)

			k := 0
			joined := false
			atj := []*NodeImpl{}
			for !joined && k<len(viableOptions){
				//fmt.Printf("\tViableOption Total: %d\n",viableOptions[0].NodeClusterParams.CurrentCluster.Total)
				if viableOptions[k].NodeClusterParams.CurrentCluster.Total < adhoc.Threshold {
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
