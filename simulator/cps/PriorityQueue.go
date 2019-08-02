package cps

import "container/heap"

// An Event is something we manage in a priority queue.

type Message int
const (
	_ = iota
	SENSE Message = iota + 1
	MOVE
	POSITION
	TIME
	SERVER
	ENERGYPRINT
	GRID
	GARBAGECOLLECT


	CLUSTERMSG
	CLUSTERHEADELECT
	CLUSTERFORM
	CLUSTERLESSFORM
	CLUSTERPRINT
	CLEANUPREADINGS
	ScheduleSensor
)

type Event struct {
	Node *NodeImpl
	Instruction Message // The Value of the item; arbitrary.
	Time  int    // The priority of the item in the queue.
	// The Index is needed by update and is maintained by the heap.Interface methods.
	Index int // The Index of the item in the heap.
}

// A PriorityQueue implements heap.Interface and holds Items.
type PriorityQueue []*Event

func (pq PriorityQueue) Len() int { return len(pq) }

func (pq PriorityQueue) Less(i, j int) bool {
	return pq[i].Time < pq[j].Time
}

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].Index = i
	pq[j].Index = j
}

func (pq *PriorityQueue) Push(x interface{}) {
	n := len(*pq)
	item := x.(*Event)
	item.Index = n
	*pq = append(*pq, item)
	heap.Fix(pq, item.Index)
}

func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	item.Index = -1 // for safety
	*pq = old[0 : n-1]
	return item
}

// update modifies the priority and Value of an Event in the queue.
func (pq *PriorityQueue) update(item *Event, value Message, priority int) {
	item.Instruction = value
	item.Time = priority
	heap.Fix(pq, item.Index)
}
