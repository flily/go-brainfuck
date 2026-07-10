package infra

import (
	"testing"
)

func TestBasicInstructionByte(t *testing.T) {
	tests := []struct {
		instruction StandardInstruction
		expected    rune
	}{
		{InstructionAdd, '+'},
		{InstructionSub, '-'},
		{InstructionPointerDec, '<'},
		{InstructionPointerInc, '>'},
		{InstructionInput, ','},
		{InstructionOutput, '.'},
		{InstructionLoopBegin, '['},
		{InstructionLoopEnd, ']'},
	}

	for _, test := range tests {
		raw := test.instruction
		var _ Instruction = raw

		if raw.Char() != test.expected {
			t.Fatalf("instruction char mismatch for '%c' (0x%x): expected 0x%x, got 0x%x",
				raw, raw, test.expected, raw)
		}
	}
}

func TestBasicInstructionString(t *testing.T) {
	tests := []struct {
		instruction StandardInstruction
		expected    string
	}{
		{InstructionAdd, "+"},
		{InstructionSub, "-"},
		{InstructionPointerDec, "<"},
		{InstructionPointerInc, ">"},
		{InstructionInput, ","},
		{InstructionOutput, "."},
		{InstructionLoopBegin, "["},
		{InstructionLoopEnd, "]"},
	}

	for _, test := range tests {
		raw := test.instruction
		var _ Instruction = raw

		if raw.String() != test.expected {
			t.Fatalf("instruction string mismatch for '%c' (0x%x): expected '%s', got '%c'",
				raw, raw, test.expected, raw)
		}
	}
}

func TestInstructionSetCheck(t *testing.T) {
	set := NewStandardInstructionSet()

	if code := set.CheckInstruction('+', nil); code == nil {
		t.Fatalf("instruction '+' not found in the standard instruction set")
	}

	if code := set.CheckInstruction('!', nil); code != nil {
		t.Fatalf("instruction '!' should not be found in the standard instruction set")
	}
}
