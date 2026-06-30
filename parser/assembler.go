package parser

import (
	"github.com/flily/go-brainfuck/vm"
)

type StackItem struct {
	Code      *vm.Code
	CodeIndex int
}

type Assembler struct {
	code  []vm.Code
	stack *Stack[StackItem]
	next  []int
}

func NewAssemblerWithCapacity(capacity int) *Assembler {
	s := &Assembler{
		stack: NewStackWithCapacity[StackItem](capacity),
		code:  make([]vm.Code, 0, 32*capacity),
		next:  make([]int, 0, 32*capacity),
	}

	return s
}

func NewAssembler() *Assembler {
	return NewAssemblerWithCapacity(DefaultStackCapacity)
}

func (s *Assembler) Push(code *vm.Code, codeIndex int) {
	item := &StackItem{
		Code:      code,
		CodeIndex: codeIndex,
	}
	s.stack.Push(item)
}

func (s *Assembler) Pop() (*vm.Code, int) {
	item, ok := s.stack.Pop()
	if !ok {
		return nil, -1
	}

	return item.Code, item.CodeIndex
}

func (s *Assembler) AddCode(code *vm.Code) bool {
	s.code = append(s.code, *code)

	result := true
	switch code.Instruction {
	case vm.InstructionLoopBegin:
		s.Push(code, len(s.code)-1)

	case vm.InstructionLoopEnd:
		_, beginIndex := s.Pop()
		if beginIndex < 0 {
			// ']' without matching '['
			result = false
		}
	}

	return result
}

func (s *Assembler) Assemble() *vm.CodeMap {
	codemap := vm.NewCodeMap()
	codemap.Codes = s.code
	codemap.Next = s.next

	return codemap
}
