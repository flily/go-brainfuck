package vm

type Reader[T MemoryUnit] interface {
	Read() (T, error)
}

type Writer[T MemoryUnit] interface {
	Write(T) error
}

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
