// Package stack is a really simple implementation of the LIFO stack structure.
package stack

// Stack represents a stack that can be populated with values of eny one type determined on stack creation.
type Stack[T any] struct {
	elements []*T
	len int
}

// CreateStack returns an empty stack that can contain elements of any one type.
func CreateStack[T any]() *Stack[T] {
	return &Stack[T] {
		elements: make([]*T, 0),
		len: 0,
	}
}

// Top returns the pointer to the element that was last added to the stack, or nil if there are no elements in the stack.
func (s *Stack[T]) Top() *T {
	if s.len <= 0 {
		return nil
	}

	return s.elements[s.len-1]
}

// Push adds an element to the stack according to the LIFO structure rules.
func (s *Stack[T]) Push(val T) {
	s.elements = append(s.elements, &val)
	s.len += 1
}

// Pop removes an element from the stack according to the LIFO structure rules.
func (s *Stack[T]) Pop() bool {
	if s.len <= 0 {
		return false
	}

	s.len -= 1

	s.elements = s.elements[:s.len]
	return true
}

// Len returns the number of elements currently inside the stack.
func (s *Stack[T]) Len() int {
	return s.len
}
