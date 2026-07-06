package utils

type Stack[T any] struct {
	items []*T
	index int
}



func NewStackWithCapacity[T any](capacity int) *Stack[T] {
	s := &Stack[T]{
		items: make([]*T, capacity),
		index: 0,
	}

	return s
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

func (s *Stack[T]) IsEmpty() bool {
	return s.index <= 0
}

func (s *Stack[T]) Size() int {
	return s.index
}
