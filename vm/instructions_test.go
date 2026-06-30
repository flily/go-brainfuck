package vm

import (
	"testing"

	"slices"
)

func TestConvertFrom(t *testing.T) {
	input := []StandardInstruction{
		InstructionAdd,
		InstructionSub,
	}

	expected := []Instruction{
		InstructionAdd,
		InstructionSub,
	}

	got := ConvertFrom(input)
	if !slices.Equal(got, expected) {
		t.Fatalf("expected %v, got %v", expected, got)
	}
}
