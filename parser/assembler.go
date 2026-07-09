package parser

import (
	"github.com/flily/go-brainfuck/infra"
)

type StackItem struct {
	Code      *infra.Code
	CodeIndex int
}

type Assembler struct {
	stack *Stack[StackItem]
	code  []infra.Code
	next  []int
}

func NewAssemblerWithCapacity(capacity int) *Assembler {
	s := &Assembler{
		stack: NewStackWithCapacity[StackItem](capacity),
		code:  make([]infra.Code, 0, 32*capacity),
		next:  make([]int, 0, 32*capacity),
	}

	return s
}

func NewAssembler() *Assembler {
	return NewAssemblerWithCapacity(DefaultStackCapacity)
}

func (s *Assembler) Push(code *infra.Code, codeIndex int) {
	item := &StackItem{
		Code:      code,
		CodeIndex: codeIndex,
	}
	s.stack.Push(item)
}

func (s *Assembler) Pop() (*infra.Code, int) {
	item, ok := s.stack.Pop()
	if !ok {
		return nil, -1
	}

	return item.Code, item.CodeIndex
}

func (s *Assembler) AddCode(code *infra.Code) bool {
	next := -1

	result := true
	switch code.Instruction {
	case infra.InstructionLoopBegin:
		s.Push(code, len(s.code))

	case infra.InstructionLoopEnd:
		_, beginIndex := s.Pop()
		if beginIndex < 0 {
			// ']' without matching '['
			result = false
		} else {
			s.next[beginIndex] = len(s.code)
			next = beginIndex
		}
	}

	s.code = append(s.code, *code)
	s.next = append(s.next, next)

	return result
}

func (s *Assembler) Assemble() *infra.CodeMap {
	codemap := infra.NewCodeMap()
	codemap.Codes = s.code
	codemap.Next = s.next

	return codemap
}
