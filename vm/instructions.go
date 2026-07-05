package vm

import (
	"github.com/flily/go-brainfuck/context"
)

type Instruction interface {
	Char() rune
	String() string
}

type Configure uint32

const (
	StandardConfigureMask Configure = 0xff000000
)

func (c Configure) IsStandard() bool {
	return (c & StandardConfigureMask) == 0
}

type ConfigureContainer interface {
	GetConfigure(Configure) uint32
}

type MemoryUnit interface {
	uint8 | uint16 | uint32 | uint64 | int8 | int16 | int32 | int64
}

type InstructionHandler[T MemoryUnit] func(vm *VM[T], conf ConfigureContainer) error

type InstructionSet interface {
	CheckInstruction(rune, *context.Context) *Code
}

type InstructionHandlerEntry[T MemoryUnit] struct {
	Instruction Instruction
	Handler     InstructionHandler[T]
}

func ConvertInstructionsFrom[T Instruction](instructions []T) []Instruction {
	h := make([]Instruction, len(instructions))
	for i, instr := range instructions {
		h[i] = instr
	}

	return h
}
