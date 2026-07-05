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

func GetStandardInstructionSetHandlers[T MemoryUnit]() []InstructionHandlerEntry[T] {
	handlers := []InstructionHandlerEntry[T]{
		{Instruction: InstructionAdd, Handler: StandardHandlerAdd[T]},
		{Instruction: InstructionSub, Handler: StandardHandlerSub[T]},
		{Instruction: InstructionPointerDec, Handler: StandardHandlerPointerDec[T]},
		{Instruction: InstructionPointerInc, Handler: StandardHandlerPointerInc[T]},
		{Instruction: InstructionInput, Handler: StandardHandlerInput[T]},
		{Instruction: InstructionOutput, Handler: StandardHandlerOutput[T]},
		{Instruction: InstructionLoopBegin, Handler: StandardHandlerLoopBegin[T]},
		{Instruction: InstructionLoopEnd, Handler: StandardHandlerLoopEnd[T]},
	}

	return handlers
}

func StandardHandlerAdd[T MemoryUnit](vm *VM[T], conf ConfigureContainer) error {
	vm.Memory[vm.DP] += 1
	return nil
}

func StandardHandlerSub[T MemoryUnit](vm *VM[T], conf ConfigureContainer) error {
	vm.Memory[vm.DP] -= 1
	return nil
}

func StandardHandlerPointerDec[T MemoryUnit](vm *VM[T], conf ConfigureContainer) error {
	vm.DP -= 1
	return nil
}

func StandardHandlerPointerInc[T MemoryUnit](vm *VM[T], conf ConfigureContainer) error {
	vm.DP += 1
	return nil
}

func StandardHandlerInput[T MemoryUnit](vm *VM[T], conf ConfigureContainer) error {
	value, err := vm.Read()
	if err != nil {
		return err
	}

	vm.Memory[vm.DP] = value
	return nil
}

func StandardHandlerOutput[T MemoryUnit](vm *VM[T], conf ConfigureContainer) error {
	value := vm.Memory[vm.DP]
	err := vm.Write(value)
	if err != nil {
		return err
	}

	return nil
}

func StandardHandlerLoopBegin[T MemoryUnit](vm *VM[T], conf ConfigureContainer) error {
	value := vm.Memory[vm.DP]
	if value == 0 {
		next := vm.Code.GetNext(vm.IP)
		vm.IP = next

	} else {
		if err := vm.PushIP(); err != nil {
			return err
		}
	}

	return nil
}

func StandardHandlerLoopEnd[T MemoryUnit](vm *VM[T], conf ConfigureContainer) error {
	if vm.SP <= 0 {
		code := vm.GetCurrentCode()
		err := ReasonCallStackEmpty.
			OnFatal(code.Context, "call stack empty").
			With("SP=%d", vm.SP)

		return err
	}

	value := vm.Memory[vm.DP]
	if value == 0 {
		next := vm.Code.GetNext(vm.IP)
		vm.IP = next

	} else {

	}

	vm.SP -= 1
	vm.IP = int(vm.IPStack[vm.SP])

	return nil
}
