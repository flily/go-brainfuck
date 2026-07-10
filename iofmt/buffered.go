package iofmt

import (
	"io"
)

type BufferedReader[T MemoryUnit] struct {
	Data   []T
	Offset int
}

func NewBufferedReader[T MemoryUnit](data []T) *BufferedReader[T] {
	r := &BufferedReader[T]{
		Data:   data,
		Offset: 0,
	}

	return r
}

func NewBufferedInput[T MemoryUnit](items ...T) *BufferedReader[T] {
	return NewBufferedReader(items)
}

func (r *BufferedReader[T]) Read() (T, error) {
	if r.Offset >= len(r.Data) {
		return 0, io.EOF
	}

	value := r.Data[r.Offset]
	r.Offset += 1

	return value, nil
}

type BufferedWriter[T MemoryUnit] struct {
	Data []T
}

func NewBufferedWriter[T MemoryUnit](size int) *BufferedWriter[T] {
	w := &BufferedWriter[T]{
		Data: make([]T, size),
	}

	return w
}

func (w *BufferedWriter[T]) Write(value T) error {
	w.Data = append(w.Data, value)
	return nil
}

func (w *BufferedWriter[T]) Dump() []T {
	return w.Data
}
