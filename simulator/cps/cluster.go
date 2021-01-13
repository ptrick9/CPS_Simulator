package cps

import (
	"math"
	"sort"
)

type AdHocNetwork struct {
	ClusterHeads        []*NodeImpl
	FullReclusters      int //Counts the number of full reclusters that occur in a simulation
	LocalReclusters     int //Counts the number of local reclusters that occur in a simulation
	AverageClusterSize  int //The current average cluster size is added to this every iteration (when cluster print is enabled)
	AverageNumClusters  int //The current number of clusters is added to this every iteration
	CSJoins             int //Counts the number of times a new node hello leads to a node joining an existing cluster
	CSSolos             int //Counts the number of times a new node hello leads to the creation of a new cluster
	LRHeads             int //Counts the number of cluster heads created by local reclusters
	LostReadings        int //Counts the number of readings that were lost because the cluster head died before sending
	NextClusterNum      int //For testing, may remove later
	TotalWaits          int //Counts the number of times a node waited instead of initiating a cluster search
	ExpansiveExtras		int //Counts the number of clusters added to an expansive recluster
	ACSResets			int //Counts the number of time an alone node reset its wait threshold based on its movement speed

	TotalClustersFormed		int // total clusters formed in the simulation
	TotalClustersDissolved	int	// total clusters dissolved in the simulation
	MaxOriginalClusterSize	int // max initial cluster size for all formed clusters
	MaxMaxClusterSize		int // max max cluster size for all formed clusters
	MaxEndClusterSize		int // max end cluster size before a cluster is dissolved
	TotalOriginalClusterSize	int //sum of original cluster sizes for all original
	TotalMaxClusterSize			int //sum of max cluster sizes for all clusters formed
	TotalEndClusterSize			int //sum of end cluster sizes before dissolved
	// totals used to compute average at end of simulation

	MovingLocalReclusters	int //Counts the number of times a local recluster is triggered by a moving head node
	DyingLocalReclusters	int //Counts the number of times a local recluster is triggered by a head node dying
	TimeLocalReclusters	int //Counts the number of times a local recluster is triggered by a head node lasting too long
	LBatteryLocalReclusters	int //Counts the number of times a local recluster is triggered by a low battery head
	BTAloneNode				int //Counts the number of messages sent by alone nodes searching for a cluster head
}

type HelloMsg struct {
	Sender      *NodeImpl //pointer to Node sending the Hello Msg
	NodeCHScore float64 //score for determining how suitable a node is to be a clusterhead
}

type ClusterNode struct { //made for testing, only has parameters that the cluster needs to know from a NodeImpl
	NodeBounds        *Bounds
	Battery           float32
	ClusterHead       *NodeImpl
	IsClusterHead     bool
	IsClusterMember   bool
}

//Computes the cluster score (higher the score the better chance a node becomes a cluster head)
func (node *NodeImpl) ComputeClusterScore(p *Params, numWithinDist int, ) float64 {
	degree := math.Min(float64(numWithinDist), float64(p.ClusterMaxThreshold))
	//battery := node.GetBatteryPercentage() * float64(p.ClusterMaxThreshold) //Multiplying by threshold ensures that battery and degree have the same maximum value

	score := math.Pow(degree, p.DegreeWeight) * math.Pow(node.GetBatteryPercentage(), p.BatteryWeight)
	return score
}

//Generates Hello Message for node to form/maintain clusters
func (node *NodeImpl) GenerateHello(score float64) {
	message := &HelloMsg{ Sender: node, NodeCHScore: score}
	node.ThisNodeHello = message
}

/*
Used for sending hello messages in both local and global reclusters

curNode - The node sending the hello message
members - The nodes participating in the recluster. nil if global
p		- The parameters of the simulation
*/
func (adhoc *AdHocNetwork) SendHello(curNode *NodeImpl, members []*NodeImpl, counter *int, p *Params) {
	withinDist := p.NodeTree.WithinRadius(p.NodeBTRange, curNode, []*NodeImpl{})
	membersWithinDist := withinDist
	if members != nil {
		membersWithinDist = SimpleIntersect(withinDist, members)
	}

	curNode.DrainBatteryBluetooth(counter) //Broadcasting first hello message, no score

	for j := 0; j < len(withinDist); j++ {
		withinDist[j].DrainBatteryBluetooth(counter) //Every node in bluetooth range receives the first hello message, allowing them to count how many neighbors they have
	}

	if curNode.IsAlive() { //If the node survived sending the first message
		curNode.GenerateHello(curNode.ComputeClusterScore(p, len(membersWithinDist))) //Score is computed using members in range because only members respond to the initial hello
		curNode.DrainBatteryBluetooth(counter)                                               //Broadcasting second hello message with score
		for j := 0; j < len(withinDist); j++ {
			if withinDist[j].IsAlive() { //If the node survived receiving the first message
				withinDist[j].DrainBatteryBluetooth(counter) //Every node in bluetooth range then receives a second hello message, the one including curNode's score
			}
		}

		for j := 0; j < len(membersWithinDist); j++ {
			if withinDist[j].IsAlive() { //If the node survived receiving the first message
				//Of nodes in bluetooth range, only members add the hello message to their received messages
				membersWithinDist[j].RecvMsgs = append(membersWithinDist[j].RecvMsgs, curNode.ThisNodeHello)
			}
		}
	}
}

/*
Called after a node takes a measurement. Cluster members attempt to send the reading to their cluster head and will
initiate a cluster search if unsuccessful. Cluster heads will be unaffected unless alone cluster search is enabled and
they do not have enough members, in which case they will also initiate a cluster search.

node	- the node that took the reading
rd		- the reading
tp		- whether the reading was a true positive
p		- the Params object of the simulation
 */
func (adhoc *AdHocNetwork) UpdateClusterStatus(node *NodeImpl, rd *Reading, tp bool, p *Params) {
	//if node.Valid && node.IsAlive() {
		 if !node.IsClusterHead && p.CurrentTime >= p.InitClusterTime {
			if node.ClusterHead == nil {
				adhoc.ClusterSearch(node, rd, tp, p)
			} else {
				head := node.ClusterHead
				node.DrainBatteryBluetooth(&node.P.Server.ReadingBTCounter) //node sends reading to node, whether or not the head is in range
				if node.IsWithinRange(head, p.NodeBTRange) {
					node.SendToClusterHead(rd, tp , head)
				} else if node.IsAlive()  {//If the node survived sending the message to the out-of-range cluster head
					node.IsClusterMember = false
					adhoc.ClusterSearch(node, rd, tp, p)
				}
			}
		} else if node.IsClusterHead {
			if p.AloneNodeClusterSearch && len(node.ClusterMembers) <= p.ClusterMinThreshold{
			 if p.AdaptiveClusterSearch && p.ACSReset && node.OldX!=0 && node.OldY!=0{
				 distance := math.Sqrt((math.Pow(.5*float64(node.SX-node.OldX), 2)) + (math.Pow(.5*float64(node.SY-node.OldY), 2)))
				 //Distance is in half meters so have to divide by half to get meters
				 metersPerSecond := distance / (float64(node.SamplingPeriod)/1000)
				 if metersPerSecond > p.MaxMoveMeters{
					 adhoc.ACSResets++
					 node.WaitThresh = p.ClusterSearchThreshold
				 }
			 }
			 adhoc.ClusterSearch(node, rd, tp, p)
		 }
		 node.LostMostMembers = len(node.StoredReadings) < int(float64(node.LastNumStoredReadings) * p.LRMemberLostThreshold)
		 node.LastNumStoredReadings = len(node.StoredReadings)
		} else {
			node.Wait = 0
			node.WaitThresh = p.ClusterSearchThreshold
		}
	//}
}

/*
Used when a node does not have a head, when is out of range of its head, or, if alone cluster search is enabled, when a
node is a head with few or no members. The node broadcasts a bluetooth message searching for nearby cluster heads. If
any respond, the node will join that cluster and send its reading if a reading has been included in the function call.
Otherwise, the node will form its own cluster. Also, if cluster search threshold is greater than 0, an actual cluster
search will only occur after this function has been called multiple times in a row by a node.

node	- the node searching for a cluster head
rd		- the reading taken by the node (nil if cluster search is not being called after a node took a reading)
tp		- whether the reading was a true positive (nil if cluster search is not being called after a node took a reading)
p		- the Params object of the simulation
*/
func (adhoc *AdHocNetwork) ClusterSearch(node *NodeImpl, rd *Reading, tp bool, p *Params) {
	if (!p.AdaptiveClusterSearch && node.Wait < p.ClusterSearchThreshold) || (p.AdaptiveClusterSearch && node.Wait < node.WaitThresh) {
		node.Wait++
		adhoc.TotalWaits++
	} else {
		node.Wait = 0
		toJoin := node.FindNearbyHeads(p, &p.Server.ClusterSearchBTCounter)

		if len(toJoin) > 0 {
			adhoc.ClearClusterParams(node)
			node.Join(toJoin[0])
			adhoc.CSJoins++
			if rd != nil {
				node.DrainBatteryBluetooth(&p.Server.ReadingBTCounter) //node sends reading to new head
				node.SendToClusterHead(rd, tp, toJoin[0])
			}
			node.WaitThresh = p.ClusterSearchThreshold
		} else {
			adhoc.ClearClusterParams(node)
			adhoc.FormCluster(node)
			adhoc.CSSolos++
			node.WaitThresh *= 2
		}
	}
}

/*
Used when a node joins a previously existing cluster.
 */
func (node *NodeImpl) Join(head *NodeImpl) {
	//node.IsClusterMember = true
	node.IsClusterHead = false
	node.ClusterHead = head
	//head.ClusterMembers[node] = node.P.CurrentTime
}

/*
Resets all cluster-related parameters for a node.

node - the node whose parameters will be reset
 */
func (adhoc *AdHocNetwork) ClearClusterParams(node *NodeImpl) {
	if node.IsClusterHead {
		node.IsClusterHead = false
		adhoc.ClusterHeads, _ = SearchRemove(adhoc.ClusterHeads, node)
		if node.P.LocalRecluster > 0 {
			adhoc.DissolveCluster(node)
			return
		}
	}
	//else if node.ClusterHead != nil {
	//	delete(node.ClusterHead.ClusterMembers, node)
	//}

	node.ClusterHead = nil
	node.LargestClusterSize = 0
	node.ClusterMembers = make(map[*NodeImpl]int)

	node.LastNumStoredReadings = 0
	node.LostMostMembers = false

	node.RecvMsgs = []*HelloMsg{}
	node.ThisNodeHello = &HelloMsg{Sender: node}

	node.IsClusterMember = false
}

/*
Calls ClearClusterParams on every member of a cluster. Assumes each member knows to do this because they are contacted
by the server.

node - the cluster head of the cluster
 */
func (adhoc *AdHocNetwork) DissolveCluster(node *NodeImpl) {
	//Assume node is cluster head
	adhoc.TotalClustersDissolved++
	adhoc.MaxMaxClusterSize = int(math.Max(float64(adhoc.MaxMaxClusterSize), float64(node.LargestClusterSize)))
	adhoc.MaxEndClusterSize = int(math.Max(float64(adhoc.MaxEndClusterSize), float64(len(node.ClusterMembers))))
	adhoc.TotalEndClusterSize += len(node.ClusterMembers)
	adhoc.TotalMaxClusterSize += node.LargestClusterSize
	for member := range node.ClusterMembers {
		adhoc.ClearClusterParams(member)
		/*	Assuming local reclustering is enabled, the members will know that this cluster has dissolved because
			they will be told so by the server */
		member.DrainBatteryWifi(1)
	}
	adhoc.ClearClusterParams(node)
}

/*
Calls ClearClusterParams on every alive node.

p - the Params object of the simulation
 */
func (adhoc *AdHocNetwork) ResetClusters(p *Params) {
	for node := range p.AliveNodes {
		node.IsClusterHead = false
		p.ClusterNetwork.ClearClusterParams(node)
	}
	adhoc.ClusterHeads = []*NodeImpl{}
}

/*
Called after a node has sent and received hello messages. From the messages received, the node will choose the node with
the highest score as its cluster head.

node - the node electing a cluster head
 */
func (adhoc *AdHocNetwork) ElectClusterHead(node *NodeImpl) {
	//maxNode := curNode.HasMaxNodeScore(p)

	node.RecvMsgs = append(node.RecvMsgs, node.ThisNodeHello)
	//Sort the messages by Id. Lowest Id first
	sort.Slice(node.RecvMsgs, func(i int, j int) bool {
		return node.RecvMsgs[i].Sender.Id < node.RecvMsgs[j].Sender.Id
	})
	//Sort the messages by score. Highest score first
	sort.Slice(node.RecvMsgs, func(i int, j int) bool {
		return node.RecvMsgs[i].NodeCHScore > node.RecvMsgs[j].NodeCHScore
	})

	i := 0
	maxNode := node.RecvMsgs[i].Sender

	if node != maxNode {
		node.ClusterHead = maxNode
	} else {
		adhoc.FormCluster(node)
	}
}

/*
Called by a node to create a cluster with itself as cluster head.

node - the cluster head of the new cluster
 */
func (adhoc *AdHocNetwork) FormCluster(node *NodeImpl) {
	node.IsClusterHead = true
	node.TimeBecameClusterHead = node.P.CurrentTime
	node.BatteryBecameClusterHead = node.GetBatteryPercentage()
	node.IsClusterMember = false
	node.InitialClusterSize = len(node.RecvMsgs)
	node.LastNumStoredReadings = 0
	node.LostMostMembers = false
	adhoc.ClusterHeads = append(adhoc.ClusterHeads, node)
	node.ClusterHead = nil
	node.LargestClusterSize = 0
	node.ClusterMembers = make(map[*NodeImpl]int)
	adhoc.NextClusterNum++
	adhoc.TotalClustersFormed++
	adhoc.MaxOriginalClusterSize = int(math.Max(float64(adhoc.MaxOriginalClusterSize), float64(node.InitialClusterSize)))
	adhoc.TotalOriginalClusterSize += node.InitialClusterSize
}

/*
Checks if two nodes are in range of each other.

searchRange - the maximum distance between to two nodes to be considered in range
 */
func (node *NodeImpl) IsWithinRange(node2 *NodeImpl, searchRange float64) bool {
	xDist := node.X - node2.X
	yDist := node.Y - node2.Y
	radDist := math.Sqrt(float64(xDist*xDist) + float64(yDist*yDist))

	return radDist <= searchRange
}

/*
Performs a global recluster that affects the entire cluster network. All current cluster information is reset and all
alive nodes broadcast hello messages and choose cluster heads.
 */
func (adhoc *AdHocNetwork) FullRecluster(p *Params) {
	adhoc.FullReclusters++
	adhoc.ResetClusters(p)
	for node := range p.AliveNodes {
		node.DrainBatteryWifi(1)	//Server sends message to all nodes that reclustering is happening
		if node.Valid {
			adhoc.SendHello(node, nil, &p.Server.GlobalReclusterBTCounter, p)
		}
	}
	for node := range p.AliveNodes {
		if !node.IsClusterHead && node.Valid {
			adhoc.ElectClusterHead(node)
		}
	}
}

/*
Performs a local recluster by having every node is the members array broadcast hello messages and elect cluster heads.
If local recluster is set to minimal, nodes will perform a cluster search and will only participate in the recluster if
they do not find a cluster head to join.

head	- the head of the cluster that makes up the members of the recluster
members	- an array of the members that will participate in the recluster
p		- the Params object of the simulation
 */
func (adhoc *AdHocNetwork) LocalRecluster(head *NodeImpl, members []*NodeImpl, p *Params) {
	//This adds any alone nodes with range of any members to be included in members. This is a simulation implementation
	//in reality this alone nodes would simply join the recluster once they received one of the hello messages.
	nearbyAlone := make(map[*NodeImpl]bool)
	for i := 0; i < len(members); i++ {
		nearby := p.NodeTree.WithinRadius(p.NodeBTRange, members[i], []*NodeImpl{})
		for _, node := range nearby {
			if node.IsClusterHead && len(node.ClusterMembers) <= 0 {
				nearbyAlone[node] = true
			}
		}
	}
	for alone := range nearbyAlone {
		members = append(members, alone)
	}
	//Local recluster starts in earnest
	adhoc.LocalReclusters++
	head.DrainBatteryWifi(1) //Head informs the server it is initiating local recluster
	for i := 0; i < len(members); i++ {
		adhoc.ClearClusterParams(members[i])
		if members[i].IsAlive() {
			if p.LocalRecluster == 1 {
				toJoin := members[i].FindNearbyHeads(p, &p.Server.LocalReclusterBTCounter)
				if len(toJoin) > 0 {
					members[i].Join(toJoin[0])
					members = append(members[:i], members[i+1:]...)
					i--
					adhoc.CSJoins++
					continue
				}
			}
			adhoc.SendHello(members[i], members, &p.Server.LocalReclusterBTCounter, p)
		}
	}
	for _, node := range members {
		if !node.IsClusterHead && node.IsAlive() {
			adhoc.ElectClusterHead(node)
		}
	}
	//For logging, not function
	for _, node := range members {
		if node.IsClusterHead {
			adhoc.LRHeads++
		}
	}
}

/*
Performs an expansive recluster by having the initiating cluster head broadcast a message asking nearby cluster heads to
join the recluster. Any nearby heads that have been a head long enough will join by contacting the server. The server
adds all the cluster heads' members to the members array and calls the LocalRecluster function.

head	- the head of the cluster that makes up the members of the recluster
members	- an array of the members that will participate in the recluster
p		- the Params object of the simulation
 */
func (adhoc *AdHocNetwork) ExpansiveLocalRecluster(head *NodeImpl, members []*NodeImpl, p *Params) {
	withinDist := p.NodeTree.WithinRadius(p.NodeBTRange, head, []*NodeImpl{})
	head.DrainBatteryBluetooth(&p.Server.LocalReclusterBTCounter) //Sends message to all in bluetooth range
	for i := 0; i < len(withinDist); i++ {
		withinDist[i].DrainBatteryBluetooth(&p.Server.LocalReclusterBTCounter) //All nodes in range receive message
		if withinDist[i].IsClusterHead && (p.CurrentTime - withinDist[i].TimeBecameClusterHead) / 1000 > int(float64(p.ClusterHeadTimeThreshold) * p.ExpansiveRatio) {
			adhoc.ExpansiveExtras++
			withinDist[i].DrainBatteryWifi(1) //This cluster head informs the server that it is within range of the dying cluster head
			adhoc.ClearClusterParams(withinDist[i])

			var memArr []*NodeImpl
			if _, ok := withinDist[i].P.Server.Clusters[withinDist[i]]; ok {
				memArr = make([]*NodeImpl, 0, len(withinDist[i].P.Server.Clusters[withinDist[i]].Members))
				for member := range withinDist[i].P.Server.Clusters[withinDist[i]].Members {
					memArr = append(memArr, member)
				}
			}
			members = append(members, memArr...)

			if withinDist[i].IsAlive() {
				members = append(members, withinDist[i])
			}
		}
	}
	adhoc.LocalRecluster(head, members, p)
}

/*
Used for logging. Allows the simulation to more accurately detect when a node is problematically out of range of its head.
 */
func (node *NodeImpl) UpdateOutOfRange(p *Params) {
	if !node.IsClusterHead {
		if node.ClusterHead != nil {
			if !node.IsWithinRange(node.ClusterHead, p.NodeBTRange) {
				if !node.OutOfRange {
					node.OutOfRange = true
					node.TimeMovedOutOfRange = p.CurrentTime
				}
			} else {
				node.OutOfRange = false
			}
		}
	} else {
		for member := range node.ClusterMembers {
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

/*
Returns the intersection of two node arrays.
 */
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

/*
Broadcasts a bluetooth message looking for nearby cluster heads and returns a list of these heads sorted by battery level.
 */
func (node *NodeImpl) FindNearbyHeads(p *Params, counter *int) []*NodeImpl {
	node.DrainBatteryBluetooth(counter)	//Broadcasting hello message
	node.P.ClusterNetwork.BTAloneNode++// stat tracking cluster search overhead
	withinDist := p.NodeTree.WithinRadius(p.NodeBTRange, node, []*NodeImpl{})
	headsInDist := []*NodeImpl{}

	for i := 0; i < len(withinDist); i++ {
		withinDist[i].DrainBatteryBluetooth(counter)	//Every node in bluetooth range receives the hello message
		if withinDist[i].IsClusterHead && len(withinDist[i].ClusterMembers) < p.ClusterMaxThreshold && withinDist[i].IsAlive() {
			headsInDist = append(headsInDist, withinDist[i])
			withinDist[i].DrainBatteryBluetooth(counter) //Every cluster head with room for this node sends a reply
			node.P.ClusterNetwork.BTAloneNode++// stat tracking cluster search overhead
			if node.IsAlive() {
				node.DrainBatteryBluetooth(counter) //The node receives messages from every cluster head with room (unless it died)
				node.P.ClusterNetwork.BTAloneNode++// stat tracking cluster search overhead
			}
		}
	}
	sort.Slice(headsInDist, func(i, j int) bool {
		return headsInDist[i].GetBatteryPercentage() > headsInDist[j].GetBatteryPercentage()
	})
	return headsInDist
}

/*
Determines if a cluster head should initiate a local recluster based on how long the node has been a head, how much
battery it has lost as a head, and its current battery level. This is only called when battery is drained.
See
 */
func (node *NodeImpl) ShouldLocalRecluster() bool {
	if (node.BatteryBecameClusterHead - node.GetBatteryPercentage()) > node.P.ClusterHeadBatteryDropThreshold {
		node.P.ClusterNetwork.LBatteryLocalReclusters++
		return true
	} else if (node.P.CurrentTime - node.TimeBecameClusterHead)/1000 > node.P.ClusterHeadTimeThreshold {
		node.P.ClusterNetwork.TimeLocalReclusters++
		return true
	} else if !node.IsAlive() {
		node.P.ClusterNetwork.DyingLocalReclusters++
		return true
	} else if node.LostMostMembers {
		node.P.ClusterNetwork.MovingLocalReclusters++
		return true
	} else {
		return false
	}
}

/*
Determines if a node has lost more than {LRMemberLostThreshold} of its members. This check is only done
for the cluster head
 */
//func (node *NodeImpl) LostMostMembers() bool {
//	if node.IsClusterHead {
//		return
//		//return 1 - float64(len(node.ClusterMembers))/float64(node.LargestClusterSize) >= node.P.LRMemberLostThreshold
//		//lost := 0.0
//		//for member := range node.ClusterMembers {
//		//	if member.ClusterHead == node && member.Wait >= node.P.ClusterSearchThreshold {
//		//		lost += 1
//		//		if lost >= node.P.LRMemberLostThreshold * float64(node.InitialClusterSize) {
//		//			return true
//		//		}
//		//	}
//		//}
//
//	}
//
//	return false
//}

/*
Initiates a local recluster by creating the initial members array and either calling the LocalRecluster or
ExpansiveLocalRecluster method.
 */
func (node *NodeImpl) InitLocalRecluster() {
	//node.IsClusterHead = false
	var members []*NodeImpl

	if _, ok := node.P.Server.Clusters[node]; ok {
		members = make([]*NodeImpl, 0, len(node.P.Server.Clusters[node].Members))
		for member := range node.P.Server.Clusters[node].Members {
			members = append(members, member)
		}
	}

	node.P.ClusterNetwork.ClearClusterParams(node)
	if node.P.LocalRecluster < 3 {
		node.P.ClusterNetwork.LocalRecluster(node, members, node.P)
	} else {
		node.P.ClusterNetwork.ExpansiveLocalRecluster(node, members, node.P)
	}
}

/*
Used to update a cluster head's knowledge about its members, based on the server's knowledge.

server - the FusionCenter object of the simulation
 */
func (node *NodeImpl) UpdateClusterInfo(server *FusionCenter) {
	node.ClusterMembers = make(map[*NodeImpl]int)
	if _, ok := server.Clusters[node]; ok {
		for k, v := range server.Clusters[node].Members {
			node.ClusterMembers[k] = v
		}
	}
	node.LargestClusterSize = max(node.LargestClusterSize, len(node.ClusterMembers))
}