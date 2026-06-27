package ast

import (
	"github.com/flily/go-brainfuck/context"
)

type Code struct {
	Instruction Instruction
	Context     *context.Context
}

type CodeMap struct {
	Files []*context.FileContext
	Codes []Code
	Next  []int
}

func NewCodeMap() *CodeMap {
	m := &CodeMap{
		Files: make([]*context.FileContext, 0),
		Codes: make([]Code, 0),
		Next:  make([]int, 0),
	}

	return m
}

func (m *CodeMap) AddFile(file *context.FileContext) {
	m.Files = append(m.Files, file)
}
