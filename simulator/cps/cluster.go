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
	NodeCHScore			float64	//score for determining how suitible a node is to be a clusterhead
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
	NodeTree			*Quadtree
}

//Computes the cluster score (higher the score the better chance a node beccomes a cluster head)
func (curNode * NodeImpl) ComputeClusterScore(penalty float64, numWithinDist int) (score float64){
	degree := float64(numWithinDist)
	battery := float64(curNode.Battery)

	//weighted sum, 60% from degree (# of nodes withinin distance), 40% from its battery life
	// penalty used to increase a nodes chance at staying a clusterhead
	return ((0.6*degree + 0.4*battery)*curNode.CHPenalty)
}

//Generates Hello Message for node to form/maintain clusters. Returns message as a string
func (curNode * NodeImpl)GenerateHello(searchRange float64, score float64) {
	var option int

	//if(curNode.IsClusterHead){
	//	option = curNode.NodeClusterParams.CurrentCluster.Total
	//} else{
		option = 0
	//}

	message := &HelloMsg{
		Sender: curNode,
		NodeCHScore: score,
		Option: option,
		BrodPeriod:	0.2}
	curNode.NodeClusterParams.ThisNodeHello = message
}

func (curNode * NodeImpl) SendHelloMessage(transmitRange float64){
	withinDist := []*Bounds{}
	withinDist = curNode.P.NodeTree.WithinRadius(transmitRange,curNode.NodeBounds,curNode.NodeBounds.GetSearchBounds(transmitRange),withinDist)
	//withinDist = curNode.NodeTree.WithinRadius(transmitRange,curNode.NodeBounds,curNode.NodeBounds.GetSearchBounds(transmitRange),withinDist)
	numWithinDist := len(withinDist)

	curNode.GenerateHello(transmitRange, curNode.ComputeClusterScore( 1,numWithinDist))

	//for j:=0; j<len(withinDist); {
	//	if(curNode.IsWithinRange(withinDist[j].CurNode,transmitRange)){
	//		j++
	//	} else{
	//			withinDist = append(withinDist[:j],withinDist[j+1:]...)
	//	}
	//}

	for j:=0; j<len(withinDist); j++ {
		if (withinDist[j].CurNode.NodeClusterParams.RecvMsgs != nil) {
			withinDist[j].CurNode.NodeClusterParams.RecvMsgs = append(withinDist[j].CurNode.NodeClusterParams.RecvMsgs, curNode.NodeClusterParams.ThisNodeHello)
			curNode.DecrementPowerBT()
		}
	}
}

func (curNode * NodeImpl) HasMaxNodeScore() (maxNode * NodeImpl){
	maxNode = curNode//&(NodeImpl{})
	maxScore := curNode.NodeClusterParams.ThisNodeHello.NodeCHScore
	for i:= 0; i<len(curNode.NodeClusterParams.RecvMsgs); i++{
		//do not consider nodes already with a clusterhead
		//if received a message from a node who does not have a cluster head
		if(!curNode.NodeClusterParams.RecvMsgs[i].Sender.IsClusterMember){
			//if their score higher than current node score
			if(curNode.NodeClusterParams.RecvMsgs[i].NodeCHScore > maxScore){
				maxScore = curNode.NodeClusterParams.RecvMsgs[i].NodeCHScore
				maxNode = curNode.NodeClusterParams.RecvMsgs[i].Sender
			}
		}
	}

	return maxNode
}

func (curNode * NodeImpl) PrintClusterNode(){
	fmt.Print("{")
	fmt.Print(curNode.NodeBounds)
	fmt.Print(" ")
	fmt.Print(curNode.Battery)
	fmt.Print(" ")
	//fmt.Print(curNode.ClusterHead)
	//fmt.Print(" ")
	fmt.Print(curNode.NodeClusterParams.CurrentCluster)
	fmt.Print(" ")
	fmt.Print(curNode.IsClusterHead)
	fmt.Print("}")
	fmt.Println()
}

func (adhoc * AdHocNetwork) ClearClusterParams(curNode * NodeImpl){
	//Reset Cluster Params (all but hello->sender since that will stay the same always)
	if(curNode.NodeClusterParams.CurrentCluster!=nil){
		curNode.NodeClusterParams.CurrentCluster.ClusterHead = nil
		curNode.NodeClusterParams.CurrentCluster.Total = 0
		curNode.NodeClusterParams.CurrentCluster.ClusterMembers = []*NodeImpl{}
	} else{
		curNode.NodeClusterParams.CurrentCluster = &Cluster{}
	}

	curNode.NodeClusterParams.RecvMsgs = []*HelloMsg{}
	if(curNode.NodeClusterParams.ThisNodeHello != nil){
		curNode.NodeClusterParams.ThisNodeHello.ClusterHead = nil
		curNode.NodeClusterParams.ThisNodeHello.NodeCHScore = 0
		curNode.NodeClusterParams.ThisNodeHello.BrodPeriod = 0
		curNode.NodeClusterParams.ThisNodeHello.Option = 0
	}else{
		curNode.NodeClusterParams.ThisNodeHello = &HelloMsg{}
	}

	if(curNode.IsClusterHead){
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
		if(adhoc.ClusterHeads[i].NodeClusterParams!=nil){
			adhoc.ClusterHeads[i].NodeClusterParams.AttemptedToJoin = []*NodeImpl{}
			if(adhoc.ClusterHeads[i].NodeClusterParams.CurrentCluster!=nil){
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
}

//sorts messages by distance to the node: 0th = closest, nth = farthest
func (curNode * NodeImpl) SortMessages(){

	distances := []float64{}

	for j:=0; j<len(curNode.NodeClusterParams.RecvMsgs);j++{
		xDist := curNode.X-curNode.NodeClusterParams.RecvMsgs[j].Sender.X
		yDist := curNode.Y-curNode.NodeClusterParams.RecvMsgs[j].Sender.Y
		dist := math.Sqrt(float64(xDist*xDist)+float64(yDist*yDist))
		distances = append(distances, dist)
	}

	for i:=0; i<len(curNode.NodeClusterParams.RecvMsgs);i++{
		for j:=0; j<len(curNode.NodeClusterParams.RecvMsgs)-i-1; j++{
			if(distances[j]>distances[j+1]){
				helloTemp := curNode.NodeClusterParams.RecvMsgs[j]
				curNode.NodeClusterParams.RecvMsgs[j] = curNode.NodeClusterParams.RecvMsgs[j+1]
				curNode.NodeClusterParams.RecvMsgs[j+1] = helloTemp

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
	if(k<len(distances)){
		curNode.NodeClusterParams.RecvMsgs  = curNode.NodeClusterParams.RecvMsgs[:k]
	}
}

func (adhoc * AdHocNetwork) ElectClusterHead(curNode * NodeImpl){
	maxNode := curNode.HasMaxNodeScore()

	if(maxNode.IsClusterHead == false){
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
	if(clusterHead.NodeClusterParams.CurrentCluster==nil){
		clusterHead.NodeClusterParams.CurrentCluster = &Cluster{clusterHead, 0, []*NodeImpl{},adhoc }
		adhoc.ClusterHeads = append(adhoc.ClusterHeads, clusterHead)
		adhoc.TotalHeads++
	}

	for i:=0; i<len(clusterHead.NodeClusterParams.RecvMsgs) && clusterHead.NodeClusterParams.CurrentCluster.Total<adhoc.Threshold; i++{
		if(!clusterHead.NodeClusterParams.RecvMsgs[i].Sender.IsClusterHead && !clusterHead.NodeClusterParams.RecvMsgs[i].Sender.IsClusterMember){
			//if(clusterHead.IsWithinRange(clusterHead.NodeClusterParams.RecvMsgs[i].Sender,8)){
			clusterHead.NodeClusterParams.CurrentCluster.ClusterMembers = append(clusterHead.NodeClusterParams.CurrentCluster.ClusterMembers, clusterHead.NodeClusterParams.RecvMsgs[i].Sender)
			clusterHead.NodeClusterParams.CurrentCluster.Total = len(clusterHead.NodeClusterParams.CurrentCluster.ClusterMembers)

			clusterHead.NodeClusterParams.RecvMsgs[i].Sender.IsClusterMember = true
			clusterHead.NodeClusterParams.RecvMsgs[i].Sender.NodeClusterParams.CurrentCluster = clusterHead.NodeClusterParams.CurrentCluster

			clusterHead.DecrementPowerBT()
			clusterHead.NodeClusterParams.RecvMsgs[i].Sender.DecrementPowerBT()

			//}
		}
	}

	for i:=0; i<len(msgs); i++ {
		if (msgs[i].Sender.IsClusterHead || msgs[i].Sender.IsClusterMember) {
			if (i < len(msgs)) {
				msgs = append(msgs[:i], msgs[i+1:]...)
			}
		}
	}

	if(len(msgs)>0){
		for i:=0; i<len(msgs); i++{
			adhoc.ElectClusterHead(msgs[i].Sender)

				//Add to SingularNodes if it has
				exists := false
				for j:=0; j<len(adhoc.SingularNodes); j++{
					if(adhoc.SingularNodes[j] == msgs[i].Sender){
						exists = true
						break
					}
				}

				if(!exists){
					adhoc.SingularNodes = append(adhoc.SingularNodes, msgs[i].Sender)
					adhoc.SingularCount++
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
			if(distances[j]>distances[j+1]){
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
	if(k<len(distances)){
		viableOptions  = viableOptions[:k]
	}

	//fmt.Println(distances)
	return viableOptions
}

func (adhoc * AdHocNetwork) FinalizeClusters(p * Params){
	for i:=0; i<len(p.NodeList); i++ {
		//Nodes marked as members but not in a cluster added to SingularNodes
		if(p.NodeList[i].IsClusterMember && !p.NodeList[i].IsClusterHead){
			if(p.NodeList[i].NodeClusterParams.CurrentCluster==nil){
			} else{
				found := false
				j := 0
				for j<len(p.ClusterNetwork.ClusterHeads){
					if(p.ClusterNetwork.ClusterHeads[j] ==  p.NodeList[i].NodeClusterParams.CurrentCluster.ClusterHead){
						break
					} else{
						j++
					}
				}
				if(j<len(p.ClusterNetwork.ClusterHeads)){
					for k := 0; k<len(p.ClusterNetwork.ClusterHeads[j].NodeClusterParams.CurrentCluster.ClusterMembers) && !found; k++{
						if(p.ClusterNetwork.ClusterHeads[j].NodeClusterParams.CurrentCluster.ClusterMembers[k]==p.NodeList[i]){
							found = true
						}
					}
				}
				if(!found){
					p.NodeList[i].IsClusterMember = false
					p.ClusterNetwork.SingularNodes = append(p.ClusterNetwork.SingularNodes, p.NodeList[i])
				}
			}
		}
	}

	//Finds viable options for singular nodes (join another cluster or form own)
	for i:=0; i<len(p.ClusterNetwork.SingularNodes); i++{
		if(p.ClusterNetwork.SingularNodes[i].IsClusterHead){
			p.ClusterNetwork.FormClusters(p.ClusterNetwork.SingularNodes[i])
		} else if(!p.ClusterNetwork.SingularNodes[i].IsClusterMember){
			viableOptions := p.ClusterNetwork.SortClusterHeads(p.ClusterNetwork.SingularNodes[i],p.NodeBTRange)

			k := 0
			joined := false
			atj := []*NodeImpl{}
			for !joined && k<len(viableOptions){
				//fmt.Printf("\tViableOption Total: %d\n",viableOptions[0].NodeClusterParams.CurrentCluster.Total)
				if(viableOptions[k].NodeClusterParams.CurrentCluster.Total < p.ClusterThreshold){
					clusterHead := viableOptions[k]

					clusterHead.NodeClusterParams.CurrentCluster.ClusterMembers = append(clusterHead.NodeClusterParams.CurrentCluster.ClusterMembers, p.ClusterNetwork.SingularNodes[i])
					clusterHead.NodeClusterParams.CurrentCluster.Total++

					p.ClusterNetwork.SingularNodes[i].IsClusterMember = true
					p.ClusterNetwork.SingularNodes[i].NodeClusterParams.CurrentCluster = clusterHead.NodeClusterParams.CurrentCluster
					joined = true

					p.ClusterNetwork.SingularNodes[i].DecrementPowerBT()
					clusterHead.DecrementPowerBT()
				} else{
					atj = append(atj, viableOptions[k])
				}
				k++
			}
			if(k==len(viableOptions)){
				p.ClusterNetwork.SingularNodes[i].NodeClusterParams.CurrentCluster = &Cluster{}
				p.ClusterNetwork.SingularNodes[i].NodeClusterParams.CurrentCluster.ClusterHead = p.ClusterNetwork.SingularNodes[i]
				p.ClusterNetwork.SingularNodes[i].NodeClusterParams.CurrentCluster.Total = 0
				p.ClusterNetwork.SingularNodes[i].NodeClusterParams.CurrentCluster.ClusterMembers = []*NodeImpl{}
				p.ClusterNetwork.SingularNodes[i].NodeClusterParams.CurrentCluster.ClusterNetwork = p.ClusterNetwork
				p.ClusterNetwork.SingularNodes[i].IsClusterHead = true
				p.ClusterNetwork.SingularNodes[i].IsClusterMember = false

				p.ClusterNetwork.ClusterHeads = append(p.ClusterNetwork.ClusterHeads, p.ClusterNetwork.SingularNodes[i])
				p.ClusterNetwork.TotalHeads++
				p.ClusterNetwork.SingularNodes[i].NodeClusterParams.AttemptedToJoin = append(p.ClusterNetwork.SingularNodes[i].NodeClusterParams.AttemptedToJoin, atj...)

				p.ClusterNetwork.SingularNodes[i].DecrementPowerBT()
			}
		}

	}

	//If no longer a singular node remove from array and decrease count
	for i:=0; i<len(p.ClusterNetwork.SingularNodes); i++{
		if(p.ClusterNetwork.SingularNodes[i].IsClusterHead || p.ClusterNetwork.SingularNodes[i].IsClusterMember){// && p.ClusterNetwork.SingularNodes[i].NodeClusterParams.CurrentCluster!=nil && p.ClusterNetwork.SingularNodes[i].NodeClusterParams.CurrentCluster.Total>0){
			p.ClusterNetwork.SingularNodes = append(p.ClusterNetwork.SingularNodes[:i],p.ClusterNetwork.SingularNodes[i+1:]...)
			p.ClusterNetwork.SingularCount--
		}
	}

	//if a cluster head is part of a cluster remove from that cluster
	for i:=0; i<len(p.ClusterNetwork.ClusterHeads); i++ {
		for j:=0; j<len(p.ClusterNetwork.ClusterHeads); j++ {
			for k:=0; k<len(p.ClusterNetwork.ClusterHeads[j].NodeClusterParams.CurrentCluster.ClusterMembers); k++ {
				if(p.ClusterNetwork.ClusterHeads[j].NodeClusterParams.CurrentCluster.ClusterMembers[k]==p.ClusterNetwork.ClusterHeads[i]){
					p.ClusterNetwork.ClusterHeads[j].NodeClusterParams.CurrentCluster.ClusterMembers = append(p.ClusterNetwork.ClusterHeads[j].NodeClusterParams.CurrentCluster.ClusterMembers[:k],p.ClusterNetwork.ClusterHeads[j].NodeClusterParams.CurrentCluster.ClusterMembers[k+1:]... )
					p.ClusterNetwork.ClusterHeads[j].NodeClusterParams.CurrentCluster.Total--
				}
			}
		}
	}
}

func (node1 * NodeImpl) IsWithinRange(node2 * NodeImpl, searchRange float64) (inRange bool){
	xDist := node1.X - node2.X
	yDist := node1.Y - node2.Y
	radDist := math.Sqrt(float64(xDist*xDist)+float64(yDist*yDist))

	inRange = radDist <= searchRange

	return inRange
}
