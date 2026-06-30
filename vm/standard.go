package vm

import (
	"slices"

	"github.com/flily/go-brainfuck/context"
)

type StandardInstruction rune

const (
	InstructionAdd        StandardInstruction = '+'
	InstructionSub        StandardInstruction = '-'
	InstructionPointerDec StandardInstruction = '<'
	InstructionPointerInc StandardInstruction = '>'
	InstructionInput      StandardInstruction = ','
	InstructionOutput     StandardInstruction = '.'
	InstructionLoopBegin  StandardInstruction = '['
	InstructionLoopEnd    StandardInstruction = ']'
)

func (i StandardInstruction) Char() rune {
	return rune(i)
}

func (i StandardInstruction) String() string {
	return string(i)
}

type StandardInstructionSet struct {
	SupportedInstructions []rune
}

var standardInstructions = []rune{
	rune(InstructionAdd),
	rune(InstructionSub),
	rune(InstructionPointerDec),
	rune(InstructionPointerInc),
	rune(InstructionInput),
	rune(InstructionOutput),
	rune(InstructionLoopBegin),
	rune(InstructionLoopEnd),
}

func NewStandardInstructionSet() InstructionSet {
	s := &StandardInstructionSet{
		SupportedInstructions: standardInstructions,
	}

	return s
}

func (s *StandardInstructionSet) CheckInstruction(r rune, ctx *context.Context) *Code {
	var result *Code = nil
	if slices.Contains(s.SupportedInstructions, r) {
		result = &Code{
			Instruction: StandardInstruction(r),
			Context:     ctx,
		}
	}

	return result
}
