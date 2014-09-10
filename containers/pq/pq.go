package pq

import (
	"container/heap"
	"sync"
)

type Elem struct {
	Priority int64
	Value interface{} 

	index int
}

type PriorityQueue interface {
	Len() int

	Push(priority int64, value interface{})
	Pop() (*Elem, bool) // Pops the LOWEST priority
}

type priorityQueue struct {
	ph *priorityHeap
	m sync.Mutex
}

var _ PriorityQueue = new(priorityQueue) // typechecking

func (pq *priorityQueue) Len() int { 
	pq.m.Lock()
	defer pq.m.Unlock()

	return len(*pq.ph) 
}

func (pq *priorityQueue) Push(priority int64, value interface{}) {
	pq.m.Lock()
	defer pq.m.Unlock()

	heap.Push(pq.ph, &Elem{priority, value, -1})
}

func (pq *priorityQueue) Pop() (*Elem, bool) {
	pq.m.Lock()
	defer pq.m.Unlock()

	if pq.ph.Len() == 0 { return nil, false }
	elem, ok := heap.Pop(pq.ph).(*Elem)
	if !ok { panic("What are you doing?") }
	return elem, true
}

func New() PriorityQueue {
	pq := &priorityQueue{&priorityHeap{}, sync.Mutex{}}
	heap.Init(pq.ph)
	return pq
}

// This implementation is almost entirely borrowed from Go's container.heap
// documentation:

type priorityHeap []*Elem

func (ph priorityHeap) Len() int { return len(ph) }

func (ph priorityHeap) Less(i, j int) bool {
	return ph[i].Priority < ph[j].Priority
}

func (pq priorityHeap) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

func (pq *priorityHeap) Push(x interface{}) {
	n := len(*pq)
	elem := x.(*Elem)
	elem.index = n
	*pq = append(*pq, elem)
}

func (pq *priorityHeap) Pop() interface{} {
	old := *pq
	n := len(old)
	elem := old[n-1]
	elem.index = -1
	*pq = old[0 : n-1]
	return elem
}

