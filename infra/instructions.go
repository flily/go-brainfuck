package infra

import (
	"github.com/flily/go-brainfuck/context"
)

type Instruction interface {
	Char() rune
	String() string
}

type MemoryUnit interface {
	uint8 | uint16 | uint32 | uint64 | int8 | int16 | int32 | int64
}

type InstructionSet interface {
	CheckInstruction(rune, *context.Context) *Code
}

func ConvertInstructionsFrom[T Instruction](instructions []T) []Instruction {
	h := make([]Instruction, len(instructions))
	for i, instr := range instructions {
		h[i] = instr
	}

	return h
}
