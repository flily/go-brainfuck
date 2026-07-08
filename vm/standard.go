package vm

import (
	"errors"
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

// In beef (https://kiyuko.org/software/beef)
// data cell under the pointer is set to 0 when EOF is reached.
//
// In the original work by Urban Müller, (https://aminet.net/package/dev/lang/brainfuck-2.lha)
// behavior is not defined when EOF is reached.
// But in c implement code, use `int getchar()` to read input and return -1 when EOF is reached.
// The data cell under the pointer is set to 255 (0xff) when -1 is assigned.
func StandardHandlerInput[T MemoryUnit](vm *VM[T], conf ConfigureContainer) error {
	value, err := vm.Read()
	if err != nil {
		if !errors.Is(err, ReasonReadEOF) {
			return err

		} else {
			if raiseErr, found := conf.GetBoolean(ConfigureReadEOFRaiseError); found && raiseErr {
				return err
			}

			if ignoreEOF, found := conf.GetBoolean(ConfigureReadValueIgnoreOnEOF); found && ignoreEOF {
				return nil
			}

			if valueOnEOF, found := conf.GetInt(ConfigureReadValueOnEOF); found {
				value = T(valueOnEOF)
			}
		}
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
	if value != 0 {
		if err := vm.UseIP(); err != nil {
			return err
		}
	} else {
		if _, err := vm.PopIP(); err != nil {
			return err
		}
	}

	return nil
}
