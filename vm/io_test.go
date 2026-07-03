package vm

import (
	"testing"
	"testing/iotest"

	"bytes"
	"errors"
	"io"
	"slices"
)

func TestGetLengthOnUnsupportedType(t *testing.T) {
	if getLength(false) != 0 {
		t.Fatalf("expect length 0, got %d", getLength(false))
	}
}

func reverse(b []byte) []byte {
	nb := slices.Clone(b)
	slices.Reverse(nb)
	return nb
}

type testMemoryUnitEncodingCase[T MemoryUnit] struct {
	BinaryBE []byte
	Unit     T
}

type testMemoryUnitEncodingCases[T MemoryUnit] []testMemoryUnitEncodingCase[T]

func (c testMemoryUnitEncodingCases[T]) Run(t *testing.T) {
	t.Helper()

	buffer := make([]byte, 16)
	typeLength := GetLength[T]()

	encoderBE := NewEncoderBE[T]()
	encoderLE := NewEncoderLE[T]()

	for _, cc := range c {
		value, length := encoderBE.Decode(cc.BinaryBE)
		if value != cc.Unit {
			t.Fatalf("expect %v (%x), got %v (%x)",
				cc.Unit, cc.Unit, value, value)
		}

		if length != typeLength {
			t.Fatalf("expect length %d, got %d", typeLength, length)
		}

		got := encoderBE.Encode(cc.Unit, buffer, 0)
		if got != typeLength {
			t.Fatalf("expect length %d, got %d", typeLength, got)
		}

		if !slices.Equal(buffer[:typeLength], cc.BinaryBE) {
			t.Fatalf("expect %v (%x), got %v (%x)",
				cc.BinaryBE, cc.BinaryBE, buffer[:typeLength], buffer[:typeLength])
		}

		binaryLE := reverse(cc.BinaryBE)
		value, length = encoderLE.Decode(binaryLE)
		if value != cc.Unit {
			t.Fatalf("expect %v (%x), got %v (%x)",
				cc.Unit, cc.Unit, value, value)
		}

		if length != typeLength {
			t.Fatalf("expect length %d, got %d", typeLength, length)
		}

		got = encoderLE.Encode(cc.Unit, buffer, 0)
		if got != typeLength {
			t.Fatalf("expect length %d, got %d", typeLength, got)
		}

		if !slices.Equal(buffer[:typeLength], binaryLE) {
			t.Fatalf("expect %v (%x), got %v (%x)",
				binaryLE, binaryLE, buffer[:typeLength], buffer[:typeLength])
		}
	}
}

func TestConvertFromUint8(t *testing.T) {
	testMemoryUnitEncodingCases[uint8]{
		{BinaryBE: []byte{0x00}, Unit: 0},
		{BinaryBE: []byte{0x01}, Unit: 1},
		{BinaryBE: []byte{0x2a}, Unit: 42},
		{BinaryBE: []byte{0xff}, Unit: 255},
	}.Run(t)
}

func TestConvertFromInt8(t *testing.T) {
	testMemoryUnitEncodingCases[int8]{
		{BinaryBE: []byte{0x00}, Unit: 0},
		{BinaryBE: []byte{0x01}, Unit: 1},
		{BinaryBE: []byte{0x2a}, Unit: 42},
		{BinaryBE: []byte{0x7f}, Unit: 127},
		{BinaryBE: []byte{0x80}, Unit: -128},
		{BinaryBE: []byte{0xff}, Unit: -1},
	}.Run(t)
}

func TestConvertFromUint16(t *testing.T) {
	testMemoryUnitEncodingCases[uint16]{
		{BinaryBE: []byte{0x00, 0x00}, Unit: 0},
		{BinaryBE: []byte{0x00, 0x01}, Unit: 1},
		{BinaryBE: []byte{0x00, 0x2a}, Unit: 42},
		{BinaryBE: []byte{0x01, 0x00}, Unit: 256},
		{BinaryBE: []byte{0x01, 0x02}, Unit: 258},
		{BinaryBE: []byte{0xff, 0xff}, Unit: 65535},
	}.Run(t)
}

func TestConvertFromInt16(t *testing.T) {
	testMemoryUnitEncodingCases[int16]{
		{BinaryBE: []byte{0x00, 0x00}, Unit: 0},
		{BinaryBE: []byte{0x00, 0x01}, Unit: 1},
		{BinaryBE: []byte{0x00, 0x2a}, Unit: 42},
		{BinaryBE: []byte{0x01, 0x00}, Unit: 256},
		{BinaryBE: []byte{0x01, 0x02}, Unit: 258},
		{BinaryBE: []byte{0xff, 0xff}, Unit: -1},
	}.Run(t)
}

func TestConvertFromUint32(t *testing.T) {
	testMemoryUnitEncodingCases[uint32]{
		{BinaryBE: []byte{0x00, 0x00, 0x00, 0x00}, Unit: 0},
		{BinaryBE: []byte{0x00, 0x00, 0x00, 0x01}, Unit: 1},
		{BinaryBE: []byte{0x00, 0x00, 0x00, 0x2a}, Unit: 42},
		{BinaryBE: []byte{0x00, 0x00, 0x01, 0x00}, Unit: 256},
		{BinaryBE: []byte{0x00, 0x00, 0x01, 0x02}, Unit: 258},
		{BinaryBE: []byte{0x00, 0x00, 0xff, 0xff}, Unit: 65535},
		{BinaryBE: []byte{0x00, 0x01, 0x00, 0x00}, Unit: 65536},
		{BinaryBE: []byte{0x00, 0x01, 0x02, 0x03}, Unit: 66051},
		{BinaryBE: []byte{0x00, 0xff, 0xff, 0xff}, Unit: 16777215},
		{BinaryBE: []byte{0x01, 0x02, 0x03, 0x04}, Unit: 16909060},
		{BinaryBE: []byte{0xff, 0xff, 0xff, 0xff}, Unit: 4294967295},
	}.Run(t)
}

func TestConvertFromInt32(t *testing.T) {
	testMemoryUnitEncodingCases[int32]{
		{BinaryBE: []byte{0x00, 0x00, 0x00, 0x00}, Unit: 0},
		{BinaryBE: []byte{0x00, 0x00, 0x00, 0x01}, Unit: 1},
		{BinaryBE: []byte{0x00, 0x00, 0x00, 0x2a}, Unit: 42},
		{BinaryBE: []byte{0x00, 0x00, 0x01, 0x00}, Unit: 256},
		{BinaryBE: []byte{0x00, 0x00, 0x01, 0x02}, Unit: 258},
		{BinaryBE: []byte{0x00, 0x00, 0xff, 0xff}, Unit: 65535},
		{BinaryBE: []byte{0x00, 0x01, 0x00, 0x00}, Unit: 65536},
		{BinaryBE: []byte{0x00, 0x01, 0x02, 0x03}, Unit: 66051},
		{BinaryBE: []byte{0x00, 0xff, 0xff, 0xff}, Unit: 16777215},
		{BinaryBE: []byte{0x01, 0x02, 0x03, 0x04}, Unit: 16909060},
		{BinaryBE: []byte{0x7f, 0xff, 0xff, 0xff}, Unit: 2147483647},
		{BinaryBE: []byte{0x80, 0x00, 0x00, 0x00}, Unit: -2147483648},
		{BinaryBE: []byte{0xff, 0xff, 0xff, 0xff}, Unit: -1},
	}.Run(t)
}

func TestConvertFromUint64(t *testing.T) {
	testMemoryUnitEncodingCases[uint64]{
		{BinaryBE: []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}, Unit: 0},
		{BinaryBE: []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01}, Unit: 1},
		{BinaryBE: []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x2a}, Unit: 42},
		{BinaryBE: []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01, 0x00}, Unit: 256},
		{BinaryBE: []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01, 0x02}, Unit: 258},
		{BinaryBE: []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0xff, 0xff}, Unit: 65535},
		{BinaryBE: []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00}, Unit: 65536},
		{BinaryBE: []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x01, 0x02, 0x03}, Unit: 66051},
		{BinaryBE: []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0xff, 0xff, 0xff}, Unit: 16777215},
		{BinaryBE: []byte{0x00, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00}, Unit: 16777216},
		{BinaryBE: []byte{0x00, 0x00, 0x00, 0x00, 0x01, 0x02, 0x03, 0x04}, Unit: 16909060},
		{BinaryBE: []byte{0x00, 0x00, 0x00, 0x00, 0xff, 0xff, 0xff, 0xff}, Unit: 4294967295},
		{BinaryBE: []byte{0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00}, Unit: 4294967296},
		{BinaryBE: []byte{0x00, 0x00, 0x00, 0x01, 0x02, 0x03, 0x04, 0x05}, Unit: 4328719365},
		{BinaryBE: []byte{0x00, 0x00, 0x00, 0xff, 0xff, 0xff, 0xff, 0xff}, Unit: 1099511627775},
		{BinaryBE: []byte{0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00}, Unit: 1099511627776},
		{BinaryBE: []byte{0x00, 0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06}, Unit: 1108152157446},
		{BinaryBE: []byte{0x00, 0x00, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff}, Unit: 281474976710655},
		{BinaryBE: []byte{0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}, Unit: 281474976710656},
		{BinaryBE: []byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07}, Unit: 283686952306183},
		{BinaryBE: []byte{0x00, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff}, Unit: 72057594037927935},
		{BinaryBE: []byte{0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}, Unit: 72057594037927936},
		{BinaryBE: []byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08}, Unit: 72623859790382856},
		{BinaryBE: []byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff}, Unit: 18446744073709551615},
	}.Run(t)
}

func TestConvertFromInt64(t *testing.T) {
	testMemoryUnitEncodingCases[int64]{
		{BinaryBE: []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}, Unit: 0},
		{BinaryBE: []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01}, Unit: 1},
		{BinaryBE: []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x2a}, Unit: 42},
		{BinaryBE: []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0xff}, Unit: 255},
		{BinaryBE: []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01, 0x00}, Unit: 256},
		{BinaryBE: []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01, 0x02}, Unit: 258},
		{BinaryBE: []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0xff, 0xff}, Unit: 65535},
		{BinaryBE: []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00}, Unit: 65536},
		{BinaryBE: []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x01, 0x02, 0x03}, Unit: 66051},
		{BinaryBE: []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0xff, 0xff, 0xff}, Unit: 16777215},
		{BinaryBE: []byte{0x00, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00}, Unit: 16777216},
		{BinaryBE: []byte{0x00, 0x00, 0x00, 0x00, 0x01, 0x02, 0x03, 0x04}, Unit: 16909060},
		{BinaryBE: []byte{0x00, 0x00, 0x00, 0x00, 0xff, 0xff, 0xff, 0xff}, Unit: 4294967295},
		{BinaryBE: []byte{0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00}, Unit: 4294967296},
		{BinaryBE: []byte{0x00, 0x00, 0x00, 0x01, 0x02, 0x03, 0x04, 0x05}, Unit: 4328719365},
		{BinaryBE: []byte{0x00, 0x00, 0x00, 0xff, 0xff, 0xff, 0xff, 0xff}, Unit: 1099511627775},
		{BinaryBE: []byte{0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00}, Unit: 1099511627776},
		{BinaryBE: []byte{0x00, 0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06}, Unit: 1108152157446},
		{BinaryBE: []byte{0x00, 0x00, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff}, Unit: 281474976710655},
		{BinaryBE: []byte{0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}, Unit: 281474976710656},
		{BinaryBE: []byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07}, Unit: 283686952306183},
		{BinaryBE: []byte{0x00, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff}, Unit: 72057594037927935},
		{BinaryBE: []byte{0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}, Unit: 72057594037927936},
		{BinaryBE: []byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08}, Unit: 72623859790382856},
		{BinaryBE: []byte{0x7f, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff}, Unit: 9223372036854775807},
		{BinaryBE: []byte{0x80, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}, Unit: -9223372036854775808},
		{BinaryBE: []byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff}, Unit: -1},
	}.Run(t)
}

func TestConvertFromFailure(t *testing.T) {
	b := []byte{0x01, 0x02}

	if n1, l := ConvertFromBE[uint32](b); n1 != 0 || l != 0 {
		t.Fatalf("expect 0, 0, got %d, %d", n1, l)
	}

	if n2, l := ConvertFromLE[uint32](b); n2 != 0 || l != 0 {
		t.Fatalf("expect 0, 0, got %d, %d", n2, l)
	}
}

func TestConvertToFailure(t *testing.T) {
	b := make([]byte, 2)

	value := uint64(42)

	if l := ConvertToBE(value, b, 0); l != 0 {
		t.Fatalf("expect 0, got %d", l)
	}

	if l := ConvertToLE(value, b, 0); l != 0 {
		t.Fatalf("expect 0, got %d", l)
	}
}

func TestReaderRead(t *testing.T) {
	data := []byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08}

	{
		buffer := bytes.NewBuffer(slices.Clone(data))
		reader := NewReader(buffer, NewEncoderBE[uint32]())
		v, err := reader.Read()
		if err != nil {
			t.Fatalf("expect no error, got %v", err)
		}

		if v != 0x01020304 {
			t.Fatalf("expect 0x01020304, got %x", v)
		}
	}

	{
		buffer := bytes.NewBuffer(slices.Clone(data))
		reader := NewReader(buffer, NewEncoderLE[uint32]())
		v, err := reader.Read()
		if err != nil {
			t.Fatalf("expect no error, got %v", err)
		}

		if v != 0x04030201 {
			t.Fatalf("expect 0x04030201, got %x", v)
		}
	}
}

func TestReaderReadFailure(t *testing.T) {
	exp := errors.New("lorem ipsum")

	reader := NewReader(iotest.ErrReader(exp), NewEncoderBE[uint32]())
	_, err := reader.Read()
	if err == nil {
		t.Fatalf("expect error, got nil")
	}

	if !errors.Is(err, exp) {
		t.Fatalf("expect %v, got %v", exp, err)
	}
}

func TestReaderReadShortData(t *testing.T) {
	data := []byte{0x01, 0x02}
	buffer := bytes.NewBuffer(data)
	reader := NewReader(iotest.DataErrReader(buffer), NewEncoderBE[uint32]())
	_, err := reader.Read()
	if err == nil {
		t.Fatalf("expect error, got nil")
	}

	if !errors.Is(err, io.ErrUnexpectedEOF) {
		t.Fatalf("expect %v, got %v", io.ErrUnexpectedEOF, err)
	}
}

func TestWriterWrite(t *testing.T) {
	{
		buffer := bytes.NewBuffer(nil)
		writer := NewWriter(buffer, NewEncoderBE[uint32]())
		err := writer.Write(0x01020304)
		if err != nil {
			t.Fatalf("expect no error, got %v", err)
		}

		expect := []byte{0x01, 0x02, 0x03, 0x04}
		if !bytes.Equal(buffer.Bytes(), expect) {
			t.Fatalf("expect %v, got %v", expect, buffer.Bytes())
		}
	}

	{
		buffer := bytes.NewBuffer(nil)
		writer := NewWriter(buffer, NewEncoderLE[uint32]())
		err := writer.Write(0x01020304)
		if err != nil {
			t.Fatalf("expect no error, got %v", err)
		}

		expect := []byte{0x04, 0x03, 0x02, 0x01}
		if !bytes.Equal(buffer.Bytes(), expect) {
			t.Fatalf("expect %v, got %v", expect, buffer.Bytes())
		}
	}
}

type errWriter struct {
	err error
}

func newErrWriter(err error) errWriter {
	return errWriter{err: err}
}

func (e errWriter) Write(p []byte) (int, error) {
	return -1, e.err
}

func TestWriterWriteFailure(t *testing.T) {
	exp := errors.New("lorem ipsum")

	writer := NewWriter(newErrWriter(exp), NewEncoderBE[uint32]())
	err := writer.Write(0x01020304)
	if err == nil {
		t.Fatalf("expect error, got nil")
	}

	if !errors.Is(err, exp) {
		t.Fatalf("expect %v, got %v", exp, err)
	}
}

type truncateWriter int

func newTruncateWriter(n int) truncateWriter {
	return truncateWriter(n)
}

func (t truncateWriter) Write(p []byte) (int, error) {
	return int(t), nil
}

func TestWriterWriteShortData(t *testing.T) {
	writer := NewWriter(newTruncateWriter(2), NewEncoderBE[uint32]())
	err := writer.Write(0x01020304)
	if err == nil {
		t.Fatalf("expect error, got nil")
	}

	if !errors.Is(err, io.ErrShortWrite) {
		t.Fatalf("expect %v, got %v", io.ErrShortWrite, err)
	}
}
