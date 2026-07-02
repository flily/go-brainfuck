package vm

import (
	"testing"

	"slices"
)

func TestConvertInstructionsFrom(t *testing.T) {
	input := []StandardInstruction{
		InstructionAdd,
		InstructionSub,
	}

	expected := []Instruction{
		InstructionAdd,
		InstructionSub,
	}

	got := ConvertInstructionsFrom(input)
	if !slices.Equal(got, expected) {
		t.Fatalf("expected %v, got %v", expected, got)
	}
}
