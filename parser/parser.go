package parser

import (
	"slices"

	"github.com/flily/go-brainfuck/context"
	"github.com/flily/go-brainfuck/vm"
)

type Parser struct {
	file            *context.FileContext
	instructionSets []vm.InstructionSet
}

func NewParser(file *context.FileContext, instructionSets ...vm.InstructionSet) *Parser {
	if len(instructionSets) <= 0 {
		instructionSets = []vm.InstructionSet{vm.NewStandardInstructionSet()}
	}

	p := &Parser{
		file:            file,
		instructionSets: slices.Clone(instructionSets),
	}

	return p
}

func (p *Parser) checkSupportedInstruction(r rune, ctx *context.Context) *vm.Code {
	for _, set := range p.instructionSets {
		if code := set.CheckInstruction(r, ctx); code != nil {
			return code
		}
	}

	return nil
}

func (p *Parser) Parse() (*vm.CodeMap, error) {
	assembler := NewAssembler()
	cursor := context.NewCursor(p.file)

	for ; !cursor.EOF(); cursor.Next() {
		r, ctx := cursor.CurrentChar()
		code := p.checkSupportedInstruction(r, ctx)
		if code != nil {
			ok := assembler.AddCode(code)
			if !ok {
				err := context.NewError(code.Context, "no matched ']' found")
				return nil, err
			}
		}
	}

	if code, index := assembler.Pop(); index >= 0 {
		err := context.NewError(code.Context, "no matched '[' found")
		return nil, err
	}

	return assembler.Assemble(), nil
}
