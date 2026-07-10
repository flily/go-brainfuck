package iofmt

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

func TestWriterWriteFailure(t *testing.T) {
	exp := errors.New("lorem ipsum")

	writer := NewWriter(NewErrWriter(exp), NewEncoderBE[uint32]())
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
