package cps

import "fmt"

type AdHocNetwork struct {
	ClusterHeads	[]*NodeImpl
	TotalHeads		int
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
	return ((0.6*degree + 0.4*battery)*penalty)
}

//Generates Hello Message for node to form/maintain clusters. Returns message as a string
func (curNode * NodeImpl)GenerateHello(searchRange float64, score float64) {
	var option int

	if(curNode.IsClusterHead){
		option = curNode.NodeClusterParams.CurrentCluster.Total
	} else{
		option = 0
	}

	message := &HelloMsg{
		curNode,
		curNode.ClusterHead,
		score,
		option,
		0.2}
	curNode.NodeClusterParams.ThisNodeHello = message
}

func (adhoc * AdHocNetwork)GenerateClusters(transmitRange float64, curNode * NodeImpl){
	//assumes hello messages have already been generated

	//Assign clusterheads and form clusters
		//node already is a cluster head OR is already in a cluster
	if(curNode.IsClusterHead || curNode.ClusterHead != nil){
		return
	} else{
		//node is not a cluster head and is not in a cluster
		for j:=0; j<len(curNode.NodeClusterParams.RecvMsgs); j++{
			//if received a message from a cluster head
			if(curNode.NodeClusterParams.RecvMsgs[j].Sender == curNode.NodeClusterParams.RecvMsgs[j].ClusterHead){
				//join cluster
				curNode.NodeClusterParams.CurrentCluster.ClusterHead = curNode.NodeClusterParams.RecvMsgs[j].Sender
				curNode.NodeClusterParams.CurrentCluster.Total++
				break
			}
		}

		//No nodes within range are cluster heads
		//find node in range with max score, make it a cluster head
		if(curNode.HasMaxNodeScore()){
			//assign self as cluster head, and all in range to be in cluster
			curNode.IsClusterHead = true

			adhoc.ClusterHeads = append(adhoc.ClusterHeads, curNode)
			adhoc.TotalHeads++

			curNode.NodeClusterParams.CurrentCluster = &Cluster{curNode,0, []*NodeImpl{}, adhoc}
			//curNode.NodeClusterParams.CurrentCluster.ClusterHead = curNode
			for j:=0; j<len(curNode.NodeClusterParams.RecvMsgs) && j<adhoc.Threshold; j++{
				if(!curNode.NodeClusterParams.RecvMsgs[j].Sender.IsClusterMember){
					curNode.NodeClusterParams.RecvMsgs[j].Sender.NodeClusterParams.CurrentCluster = curNode.NodeClusterParams.CurrentCluster

					curNode.NodeClusterParams.CurrentCluster.ClusterMembers = append(curNode.NodeClusterParams.CurrentCluster.ClusterMembers, curNode.NodeClusterParams.RecvMsgs[j].Sender)

					curNode.NodeClusterParams.CurrentCluster.Total++
					if(!curNode.NodeClusterParams.RecvMsgs[j].Sender.IsClusterHead) {
						curNode.NodeClusterParams.RecvMsgs[j].Sender.IsClusterMember = true
					}

					//update hello / send second messsage saying sender joined a cluster
					curNode.NodeClusterParams.RecvMsgs[j].Sender.NodeClusterParams.ThisNodeHello.ClusterHead = curNode
				}
			}
		}
	}
}

func (curNode * NodeImpl) SendHelloMessage(transmitRange float64){
	withinDist := []*Bounds{}
	withinDist = curNode.P.NodeTree.WithinRadius(transmitRange,curNode.NodeBounds,curNode.NodeBounds.GetSearchBounds(transmitRange),withinDist)
	//withinDist = curNode.NodeTree.WithinRadius(transmitRange,curNode.NodeBounds,curNode.NodeBounds.GetSearchBounds(transmitRange),withinDist)
	numWithinDist := len(withinDist)

	curNode.GenerateHello(transmitRange, curNode.ComputeClusterScore( 1,numWithinDist))

	for j:=0; j<len(withinDist); j++ {
		if (withinDist[j].CurNode.NodeClusterParams.RecvMsgs != nil) {
			withinDist[j].CurNode.NodeClusterParams.RecvMsgs = append(withinDist[j].CurNode.NodeClusterParams.RecvMsgs, curNode.NodeClusterParams.ThisNodeHello)
		}
	}
}

func (curNode * NodeImpl) HasMaxNodeScore() (isMax bool){
	maxNode := curNode//&(NodeImpl{})
	maxScore := curNode.NodeClusterParams.ThisNodeHello.NodeCHScore
	for i:= 0; i<len(curNode.NodeClusterParams.RecvMsgs); i++{
		//do not consider nodes already with a clusterhead
		//if received a message from a node who has cluster head
		if(!curNode.NodeClusterParams.RecvMsgs[i].Sender.IsClusterMember){
			if(curNode.NodeClusterParams.RecvMsgs[i].NodeCHScore > maxScore){
				maxScore = curNode.NodeClusterParams.RecvMsgs[i].NodeCHScore
				maxNode = curNode.NodeClusterParams.RecvMsgs[i].Sender
			}
		}
	}

	return(maxNode == curNode)
}

func (curNode * NodeImpl) PrintClusterNode(){
	fmt.Print("{")
	fmt.Print(curNode.NodeBounds)
	fmt.Print(" ")
	fmt.Print(curNode.Battery)
	fmt.Print(" ")
	fmt.Print(curNode.ClusterHead)
	fmt.Print(" ")
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

	curNode.IsClusterMember = false
	curNode.IsClusterHead = false
	curNode.ClusterHead = nil
	adhoc.ClusterHeads = []*NodeImpl{}
	adhoc.TotalHeads = 0

}

