package vm

import (
	"github.com/flily/go-brainfuck/vm/utils"
)

type StackItem struct {
	Code      *Code
	CodeIndex int
}

type Assembler struct {
	stack *utils.Stack[StackItem]
	code  []Code
	next  []int
}

func NewAssemblerWithCapacity(capacity int) *Assembler {
	s := &Assembler{
		stack: utils.NewStackWithCapacity[StackItem](capacity),
		code:  make([]Code, 0, 32*capacity),
		next:  make([]int, 0, 32*capacity),
	}

	return s
}

func NewAssembler() *Assembler {
	return NewAssemblerWithCapacity(DefaultStackCapacity)
}

func (s *Assembler) Push(code *Code, codeIndex int) {
	item := &StackItem{
		Code:      code,
		CodeIndex: codeIndex,
	}
	s.stack.Push(item)
}

func (s *Assembler) Pop() (*Code, int) {
	item, ok := s.stack.Pop()
	if !ok {
		return nil, -1
	}

	return item.Code, item.CodeIndex
}

func (s *Assembler) AddCode(code *Code) bool {
	next := -1

	result := true
	switch code.Instruction {
	case InstructionLoopBegin:
		s.Push(code, len(s.code))

	case InstructionLoopEnd:
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

func (s *Assembler) Assemble() *CodeMap {
	codemap := NewCodeMap()
	codemap.Codes = s.code
	codemap.Next = s.next

	return codemap
}
