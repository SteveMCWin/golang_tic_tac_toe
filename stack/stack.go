package stack

import "log"

type Stack[T any] struct {
	elements []*T
	len int
}

func CreateStack[T any]() *Stack[T] {
	return &Stack[T] {
		elements: make([]*T, 0),
		len: 0,
	}
}

func (s *Stack[T]) Top() (*T, bool) {
	if s.len <= 0 {
		return nil, false
	}

	return s.elements[s.len-1], true
}

func (s *Stack[T]) Push(val T) {
	log.Println("Pushed to stack:", val)
	s.elements = append(s.elements, &val)
	s.len += 1
}

func (s *Stack[T]) Pop() bool {
	if s.len <= 0 {
		return false
	}

	s.len -= 1

	s.elements = s.elements[:s.len]
	return true
}

func (s *Stack[T]) Len() int {
	return s.len
}
