package infra

import (
	"slices"

	"github.com/flily/go-brainfuck/context"
)

type Code struct {
	Instruction Instruction
	Context     *context.Context
}

func InstructionsToCodes(instructions []Instruction) []Code {
	codes := make([]Code, len(instructions))
	for i, ins := range instructions {
		codes[i].Instruction = ins
	}

	return codes
}

func CodesEqual(a []Code, b []Code) bool {
	if len(a) != len(b) {
		return false
	}

	for i, ac := range a {
		bc := b[i]
		if ac.Instruction != bc.Instruction {
			return false
		}
	}

	return true
}

type CodeMap struct {
	Files []*context.FileContext
	Codes []Code
	Next  []int
}

func NewCodeMap() *CodeMap {
	m := &CodeMap{
		Files: make([]*context.FileContext, 0, 2),
		Codes: nil,
		Next:  nil,
	}

	return m
}

func (m *CodeMap) GetNext(ip int) int {
	if ip < 0 || ip >= len(m.Next) {
		return -1
	}

	return m.Next[ip]
}

func (m *CodeMap) AddFile(file *context.FileContext) {
	m.Files = append(m.Files, file)
}

// CodeEquals compares instructions and next table of two CodeMap instances for equality.
// Context and file information are not compared.
func (m *CodeMap) CodeEquals(other *CodeMap) bool {
	if !CodesEqual(m.Codes, other.Codes) {
		return false
	}

	return slices.Equal(m.Next, other.Next)
}

func (m *CodeMap) Snapshot(ip int, codeBefore int, codeAfter int) []Code {
	codeStart := max(ip-codeBefore, 0)
	codeEnd := min(ip+codeAfter+1, len(m.Codes))
	return slices.Clone(m.Codes[codeStart:codeEnd])
}
