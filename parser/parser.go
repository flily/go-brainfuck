package parser

import (
	"slices"

	"github.com/flily/go-brainfuck/context"
	"github.com/flily/go-brainfuck/infra"
)

const DefaultStackCapacity = 1024

type (
	Instruction    = infra.Instruction
	InstructionSet = infra.InstructionSet
	CodeMap        = infra.CodeMap
)

type Parser struct {
	file            *context.FileContext
	instructionSets []InstructionSet
}

func NewParser(file *context.FileContext, instructionSets ...InstructionSet) *Parser {
	if len(instructionSets) <= 0 {
		instructionSets = []InstructionSet{infra.NewStandardInstructionSet()}
	}

	p := &Parser{
		file:            file,
		instructionSets: slices.Clone(instructionSets),
	}

	return p
}

func (p *Parser) checkSupportedInstruction(r rune, ctx *context.Context) *infra.Code {
	for _, set := range p.instructionSets {
		if code := set.CheckInstruction(r, ctx); code != nil {
			return code
		}
	}

	return nil
}

func (p *Parser) Parse() (*infra.CodeMap, error) {
	assembler := NewAssembler()
	cursor := context.NewCursor(p.file)

	for ; !cursor.EOF(); cursor.Next() {
		r, ctx := cursor.CurrentChar()
		code := p.checkSupportedInstruction(r, ctx)
		if code != nil {
			ok := assembler.AddCode(code)
			if !ok {
				err := context.NewError(code.Context, "unexpected closing loop bracket").
					With("no matched '[' for this")
				return nil, err
			}
		}
	}

	if code, index := assembler.Pop(); index >= 0 {
		err := context.NewError(code.Context, "unclosed loop bracket").
			With("no matched ']' for this")
		return nil, err
	}

	return assembler.Assemble(), nil
}
