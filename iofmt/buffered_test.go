package iofmt

import (
	"testing"

	"io"
)

func TestBufferedReaderRead(t *testing.T) {
	reader := NewBufferedInput[int32](1, 2, 3)
	expected := []int32{1, 2, 3}
	for i, exp := range expected {
		value, err := reader.Read()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if value != exp {
			t.Fatalf("expected %v, got %v at index %d", exp, value, i)
		}
	}

	value, err := reader.Read()
	if err != io.EOF {
		t.Fatalf("expected EOF error, got %v", err)
	}

	if value != 0 {
		t.Fatalf("expected zero value, got %v", value)
	}
}

func TestBufferedWriterWrite(t *testing.T) {
	writer := NewBufferedWriter[int32](0)
	expected := []int32{1, 2, 3}
	for _, exp := range expected {
		err := writer.Write(exp)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	}

	output := writer.Dump()
	if len(output) != len(expected) {
		t.Fatalf("expected output length %d, got %d", len(expected), len(output))
	}

	for i, exp := range expected {
		if output[i] != exp {
			t.Fatalf("expected %v, got %v at index %d", exp, output[i], i)
		}
	}
}
