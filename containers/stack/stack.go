package stack

import (
	"sync"
)

/* Implements a thread-safe stack using container/list */

type Stack interface {
	Len() int

	Push(interface{}) 
	Pop() (interface{}, bool) // second argument is true if stack is empty
	
	AddSlice([]interface{})
}

var _ Stack = &stack{} // typechecking

type stack struct {
	xs []interface{}
	m *sync.Mutex
}

func New() Stack {
	return &stack{xs: make([]interface{}, 0), m: new(sync.Mutex)}
}

func (s *stack)Len() int {
	s.m.Lock()
	defer s.m.Unlock()

	return len(s.xs)
}

func (s *stack)Push(v interface{}) {
	s.m.Lock()
	defer s.m.Unlock()

	s.xs = append(s.xs, v)
}

func (s *stack)Pop() (interface{}, bool) {
	s.m.Lock()
	defer s.m.Unlock()

	if len(s.xs) == 0 { return nil, true }

	v := s.xs[len(s.xs) - 1]
	s.xs = s.xs[0: len(s.xs) - 1]
	return v, false
}

func (s *stack)AddSlice(vs []interface{}) {
	s.m.Lock()
	defer s.m.Unlock()

	s.xs = append(s.xs, vs)
}