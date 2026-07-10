package iofmt

import (
	"io"

	"github.com/flily/go-brainfuck/infra"
)

type (
	MemoryUnit = infra.MemoryUnit
)

func getLength(t any) int {
	switch t.(type) {
	case uint8, int8:
		return 1
	case uint16, int16:
		return 2
	case uint32, int32:
		return 4
	case uint64, int64:
		return 8

	default:
		return 0
	}
}

func GetLength[T any]() int {
	var t T
	return getLength(t)
}

type Reader[T MemoryUnit] interface {
	Read() (T, error)
}

type Writer[T MemoryUnit] interface {
	Write(T) error
}

type DumpableWriter[T MemoryUnit] interface {
	Writer[T]
	Dump() []T
}

type MemoryUnitReader[T MemoryUnit] struct {
	Reader  io.Reader
	Encoder Encoder[T]
}

func NewReader[T MemoryUnit](r io.Reader, encoder Encoder[T]) Reader[T] {
	e := &MemoryUnitReader[T]{
		Reader:  r,
		Encoder: encoder,
	}

	return e
}

func (r *MemoryUnitReader[T]) Read() (T, error) {
	var buffer [8]byte
	l := GetLength[T]()
	n, err := r.Reader.Read(buffer[:l])
	if 0 < n && n < l {
		return 0, io.ErrUnexpectedEOF
	}

	if err != nil {
		return 0, err
	}

	value, _ := r.Encoder.Decode(buffer[:l])
	return value, nil
}

type MemoryUnitWriter[T MemoryUnit] struct {
	Writer  io.Writer
	Encoder Encoder[T]
}

func NewWriter[T MemoryUnit](w io.Writer, encoder Encoder[T]) Writer[T] {
	e := &MemoryUnitWriter[T]{
		Writer:  w,
		Encoder: encoder,
	}

	return e
}

func (w *MemoryUnitWriter[T]) Write(value T) error {
	var buffer [8]byte
	l := w.Encoder.Encode(value, buffer[:], 0)
	n, err := w.Writer.Write(buffer[:l])
	if err != nil {
		return err
	}

	if n < l {
		return io.ErrShortWrite
	}

	return nil
}
