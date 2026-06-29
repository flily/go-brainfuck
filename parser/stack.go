package parser

type Stack[T any] struct {
	items []*T
	index int
}

const DefaultStackCapacity = 1024

func NewStackWithCapacity[T any](capacity int) *Stack[T] {
	s := &Stack[T]{
		items: make([]*T, capacity),
		index: 0,
	}

	return s
}

func NewStack[T any]() *Stack[T] {
	return NewStackWithCapacity[T](DefaultStackCapacity)
}

func (s *Stack[T]) Push(item *T) {
	s.items[s.index] = item
	s.index++

	if s.index >= len(s.items) {
		newStack := make([]*T, len(s.items)*2)
		copy(newStack, s.items)
		s.items = newStack
	}
}

func (s *Stack[T]) Pop() (*T, bool) {
	if s.index <= 0 {
		return nil, false
	}

	s.index--
	item := s.items[s.index]
	s.items[s.index] = nil
	return item, true
}
