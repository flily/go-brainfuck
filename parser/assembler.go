package parser

import (
	"github.com/flily/go-brainfuck/ast"
)

type StackItem struct {
	Code      *ast.Code
	CodeIndex int
}

type Assembler struct {
	code  []ast.Code
	stack *Stack[StackItem]
	next  []int
}

func NewAssemblerWithCapacity(capacity int) *Assembler {
	s := &Assembler{
		stack: NewStackWithCapacity[StackItem](capacity),
		code:  make([]ast.Code, 0, 32*capacity),
		next:  make([]int, 0, 32*capacity),
	}

	return s
}

func NewAssembler() *Assembler {
	return NewAssemblerWithCapacity(DefaultStackCapacity)
}

func (s *Assembler) Push(code *ast.Code, codeIndex int) {
	item := &StackItem{
		Code:      code,
		CodeIndex: codeIndex,
	}
	s.stack.Push(item)
}

func (s *Assembler) Pop() (*ast.Code, int) {
	item, ok := s.stack.Pop()
	if !ok {
		return nil, -1
	}

	return item.Code, item.CodeIndex
}

func (s *Assembler) AddCode(code *ast.Code) bool {
	s.code = append(s.code, *code)

	result := true
	switch code.Instruction {
	case ast.InstructionLoopBegin:
		s.Push(code, len(s.code)-1)

	case ast.InstructionLoopEnd:
		_, beginIndex := s.Pop()
		if beginIndex < 0 {
			// ']' without matching '['
			result = false
		}
	}

	return result
}

func (s *Assembler) Assemble() *ast.CodeMap {
	codemap := ast.NewCodeMap()
	codemap.Codes = s.code
	codemap.Next = s.next

	return codemap
}
