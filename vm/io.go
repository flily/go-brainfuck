package vm

import (
	"io"
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

func ConvertFromBE[T MemoryUnit](input []byte) (T, int) {
	l := GetLength[T]()
	if len(input) < l {
		return 0, 0
	}

	var result T
	for i := range l {
		result |= T(input[i]) << (8 * (l - 1 - i))
	}

	return result, l
}

func ConvertToBE[T MemoryUnit](value T, output []byte, offset int) int {
	l := GetLength[T]()
	if len(output) < offset+l {
		return 0
	}

	for i := range l {
		output[offset+i] = byte(value >> (8 * (l - 1 - i)))
	}

	return l
}

func ConvertFromLE[T MemoryUnit](input []byte) (T, int) {
	l := GetLength[T]()
	if len(input) < l {
		return 0, 0
	}

	var result T
	for i := range l {
		result |= T(input[i]) << (8 * i)
	}

	return result, l
}

func ConvertToLE[T MemoryUnit](value T, output []byte, offset int) int {
	l := GetLength[T]()
	if len(output) < offset+l {
		return 0
	}

	for i := range l {
		output[offset+i] = byte(value >> (8 * i))
	}

	return l
}

type EndianEncoder[T MemoryUnit] bool

func NewBigEndianEncoder[T MemoryUnit]() Encoder[T] {
	return EndianEncoder[T](true)
}

func NewEncoderBE[T MemoryUnit]() Encoder[T] {
	return NewBigEndianEncoder[T]()
}

func NewLittleEndianEncoder[T MemoryUnit]() Encoder[T] {
	return EndianEncoder[T](false)
}

func NewEncoderLE[T MemoryUnit]() Encoder[T] {
	return NewLittleEndianEncoder[T]()
}

func (e EndianEncoder[T]) Encode(value T, output []byte, offset int) int {
	if e {
		return ConvertToBE(value, output, offset)
	} else {
		return ConvertToLE(value, output, offset)
	}
}

func (e EndianEncoder[T]) Decode(input []byte) (T, int) {
	if e {
		return ConvertFromBE[T](input)
	} else {
		return ConvertFromLE[T](input)
	}
}

type Encoder[T MemoryUnit] interface {
	Encode(T, []byte, int) int
	Decode([]byte) (T, int)
}

type Reader[T MemoryUnit] interface {
	Read() (T, error)
}

type Writer[T MemoryUnit] interface {
	Write(T) error
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
