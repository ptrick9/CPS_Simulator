package cps

type Cluster struct {
	ClusterHead			*NodeImpl	//id of clusterhead
	Threshold			int //maximum # of nodes in a cluster
	Total				int //current # of nodes in a cluster
}

type ClusterMemberParams struct{
	CurrentCluster		*Cluster
	RecvMsgs			[]HelloMsg
}

type HelloMsg struct {
	Sender				*NodeImpl		//pointer to Node sending the Hello Msg
	ClusterHead			*NodeImpl		//pointer to of the cluster head
										//nil if not in a cluster
										//points to self if Node is a ClusterHead
	NodeCHScore			float64	//score for determining how suitible a node is to be a clusterhead
	Option				int		//0 for regular node, if a cluster head this is the # of nodes in the cluster
	BrodPeriod			float64	//broadcast period of the Sender
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
func (curNode * NodeImpl)GenerateHello(searchRange float64, score float64) (message HelloMsg){
	var option int

	if(curNode.IsClusterHead){
		option = curNode.NodeClusterParams.CurrentCluster.Total
	} else{
		option = 0
	}

	message = HelloMsg{
		curNode,
		curNode.ClusterHead,
		score,
		option,
		0.2}
	return message
}

func GenerateClusters(p * Params, searchRange float64){
	scores := make([]float64,len(p.NodeList))

	//Step 1: send hello messages from all nodes to all their neighbors
	for i:=0; i<len(p.NodeList); i++{
		withinDist := []*Bounds{}
		withinDist = p.NodeList[i].P.NodeTree.WithinRadius(searchRange,p.NodeList[i].NodeBounds,p.NodeList[i].NodeBounds.GetSearchBounds(searchRange),withinDist)
		numWithinDist := len(withinDist)

		scores[i] = p.NodeList[i].ComputeClusterScore( 1,numWithinDist)
		msg := p.NodeList[i].GenerateHello(searchRange, scores[i])

		for j:=0; j<len(withinDist); j++{
			if(withinDist[j].curNode.NodeClusterParams.RecvMsgs == nil){
				withinDist[j].curNode.NodeClusterParams.RecvMsgs = []HelloMsg{}
			}
			withinDist[j].curNode.NodeClusterParams.RecvMsgs = append(withinDist[j].curNode.NodeClusterParams.RecvMsgs, msg)
		}

	}


	//Step 2: assign clusterheads and form clusters
	for i:=0; i<len(p.NodeList); i++{
		//node already is a cluster head OR is already in a cluster
		if(p.NodeList[i].IsClusterHead || p.NodeList[i].ClusterHead != nil){
			continue
		}else{
			//node is not a cluster head and is not in a cluster
			for j:=0; j<len(p.NodeList[i].NodeClusterParams.RecvMsgs); j++{
				//if received a message from a cluster head
				if(p.NodeList[i].NodeClusterParams.RecvMsgs[j].Sender == p.NodeList[i].NodeClusterParams.RecvMsgs[j].ClusterHead){
					//join cluster
					p.NodeList[i].NodeClusterParams.CurrentCluster.ClusterHead = p.NodeList[i].NodeClusterParams.RecvMsgs[j].Sender
					p.NodeList[i].NodeClusterParams.CurrentCluster.Total++
					break
				}
			}

			thisNodeScore := scores[i]
			max := thisNodeScore
			maxScoreNode := &(NodeImpl{})
			//did not receive a message from a cluster head
			for j:=0; j<len(p.NodeList[i].NodeClusterParams.RecvMsgs); j++{
				//compare scores to all who sent message
				if(p.NodeList[i].NodeClusterParams.RecvMsgs[j].NodeCHScore > max){
					max = p.NodeList[i].NodeClusterParams.RecvMsgs[j].NodeCHScore
					maxScoreNode = p.NodeList[i]
				}
			}

			//if node score highest
			if(maxScoreNode == p.NodeList[i]){
				//assign self as cluster head, and all in range to be in cluster
				p.NodeList[i].NodeClusterParams.CurrentCluster.ClusterHead = p.NodeList[i]
				for j:=0; j<len(p.NodeList[i].NodeClusterParams.RecvMsgs); j++{
					p.NodeList[i].NodeClusterParams.RecvMsgs[i].Sender.NodeClusterParams.CurrentCluster.ClusterHead = p.NodeList[i]
				}
			}
		}
	}
}