package routeservice

import (
	"container/heap"
)

// An Item is something we manage in a priority queue.
type Item struct {
	node     int
	priority int
	index    int
}

// A PriorityQueue implements heap.Interface and holds Items.
type PriorityQueue []*Item

// NewPriorityQueue is a constructor for PriorityQueue.
func NewPriorityQueue() *PriorityQueue {
	pq := make(PriorityQueue, 0)
	heap.Init(&pq)

	return &pq
}

func (pq PriorityQueue) Len() int { return len(pq) }

func (pq PriorityQueue) Less(i, j int) bool {
	// We want Pop to give us the lowest priority so we use lesser than here.
	return pq[i].priority < pq[j].priority
}

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

// Push inserts a Item into the PriorityQueue.
func (pq *PriorityQueue) Push(x interface{}) {
	n := len(*pq)
	item := x.(*Item)
	item.index = n
	*pq = append(*pq, item)
}

// Pop removes an Item from the PriorityQueue and returns it.
func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil  // avoid memory leak
	item.index = -1 // for safety
	*pq = old[0 : n-1]
	return item
}

// Update reorganizes the Items indexes.
func (pq *PriorityQueue) Update(item *Item) {
	heap.Fix(pq, item.index)
}

func (pq *PriorityQueue) getItemReferenceByNode(node int) *Item {
	for _, item := range *pq {
		if item.node == node {
			return item
		}
	}

	return nil
}
