package ast

import (
	"github.com/flily/go-brainfuck/context"
)

type Instruction interface {
	Char() rune
	String() string
}

type InstructionSet interface {
	CheckInstruction(rune, *context.Context) *Code
}

func ConvertFrom[T Instruction](instructions []T) []Instruction {
	h := make([]Instruction, len(instructions))
	for i, instr := range instructions {
		h[i] = instr
	}

	return h
}
