package queue

import (
	"sync"
	"container/list"
)

/* Implements a thread-safe queue using container/list */

type Queue interface {
	Len() int

	Enq(interface{}) 
	Deq() (interface{}, bool) // second argument is true if queue is empty

	AddSlice([]interface{})
}

var _ Queue = &queue{} // typechecking

type queue struct {
	xs *list.List
	m *sync.Mutex
}

func New() Queue {
	return &queue{xs: list.New(), m: new(sync.Mutex)}
}

func (q *queue)Len() int {
	q.m.Lock()
	defer q.m.Unlock()

	return q.xs.Len()
}

func (q *queue)Enq(v interface{}) {
	q.m.Lock()
	defer q.m.Unlock()

	q.xs.PushFront(v)
}

func (q *queue) Deq() (interface{}, bool) {
	q.m.Lock()
	defer q.m.Unlock()

	if q.xs.Len() == 0 { return nil, true }

	elem := q.xs.Back()
	q.xs.Remove(elem)
	return elem.Value, false
}

func (q *queue)AddSlice(vs []interface{}) {
	q.m.Lock()
	defer q.m.Unlock()

	for _, v := range vs { q.xs.PushFront(v) }
}