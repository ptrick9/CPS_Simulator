package main

import(
	"./cps"
	"container/heap"
	"fmt"
)

func main() {
	pq := make(cps.PriorityQueue, 0)

	heap.Init(&pq)

	for i := 0; i < 10; i++ {
		item := &cps.Event{
			 "cat",
			 cps.RandomInt(0, 100),
			 i,
		}
		//heap.Push(&pq, item)
		pq.Push(item)
		//heap.Fix(&pq, i)
	}

	for pq.Len() > 0 {
		//item := pq.Pop().(*cps.Event)
		item := heap.Pop(&pq).(*cps.Event)
		fmt.Printf("%.2d:%s ", item.Time, item.Value)
	}

}
