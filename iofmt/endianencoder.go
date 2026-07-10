package iofmt

type Encoder[T MemoryUnit] interface {
	Encode(T, []byte, int) int
	Decode([]byte) (T, int)
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
