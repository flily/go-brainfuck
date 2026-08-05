package driver

import (
	"fmt"

	"github.com/flily/go-brainfuck/config"
	"github.com/flily/go-brainfuck/context"
)

type Parser struct {
	tokenizer *Tokenizer
}

func NewParser(filename string, file *context.FileContext) *Parser {
	p := &Parser{
		tokenizer: NewTokenizer(file),
	}

	return p
}

func NewParserWith(filename string, data []byte) *Parser {
	file := context.ReadFileData(filename, data)
	return NewParser(filename, file)
}

func (p *Parser) nextToken() (*Element, error) {
	return p.tokenizer.Next()
}

func (p *Parser) expectToken(expected Token) (*Element, error) {
	token, err := p.nextToken()
	if err != nil {
		return nil, err
	}

	if token.Token != expected {
		err = token.Errorf("unexpected token type").
			With("expect %s here, got %s", expected, token.Token)
		return nil, err
	}

	return token, nil
}

func (p *Parser) parseScript(item *TestDriverItem) error {
	name, err := p.expectToken(TokenIdentifier)
	if err != nil {
		return err
	}

	item.ScriptName = name.StringValue()

	return nil
}

func (p *Parser) setInitParameter(item *TestDriverItem, name ContextItem[string], value *Element) error {
	var err error
	switch name.Value {
	case FieldMemorySize:
		item.Init.Value.MemorySize = value.UintValue()

	case FieldStackSize:
		item.Init.Value.StackSize = value.UintValue()

	case FieldWord:
		var unitType config.MemoryUnitType
		vErr := unitType.Set(value.ValueString)
		if vErr != nil {
			err = value.Context.Error("invalid parameter value").
				With("invalid memory unit type '%s'", value.ValueString)
		}

		if err == nil {
			item.Init.Value.WordType = NewContextItem(unitType, value.Context)
		}

	case FieldEOFValue:
		item.Init.Value.EOFValue = value.IntValue()

	case FieldIgnoreEOF:
		item.Init.Value.IgnoreEOF, err = value.BoolValue()

	case FieldRaiseEOF:
		item.Init.Value.RaiseEOF, err = value.BoolValue()

	default:
		err = name.Context.Error("unknown field name").
			With("unknown field name '%s'", name.Value)
	}

	return err
}

func (p *Parser) parseInitParameters(item *TestDriverItem) (bool, error) {
	name, err := p.nextToken()
	if err != nil {
		return true, err
	}

	stop := false
	switch name.Token {
	case TokenIdentifier:
		var value *Element
		value, err = p.nextToken()
		if err != nil {
			return stop, err
		}

		err = p.setInitParameter(item, name.StringValue(), value)

	case TokenBraceRight:
		stop = true

	case TokenEOF:
		err = name.Errorf("unexpected EOF").
			With("expect '}' to close")

	default:
		err = name.Errorf("unexpected token type").
			With("expect identifier or '}', got %s", name.Token)
	}

	return stop, err
}

func (p *Parser) parseInit(item *TestDriverItem) error {
	if _, err := p.expectToken(TokenBraceLeft); err != nil {
		return err
	}

	for {
		stop, err := p.parseInitParameters(item)
		if err != nil {
			return err
		}

		if stop {
			break
		}
	}

	return nil
}

func (p *Parser) parseNumbersInList() ([]ContextItem[int64], error) {
	finish := false
	var result []ContextItem[int64]
	for !finish {
		token, err := p.nextToken()
		if err != nil {
			return nil, err
		}

		switch token.Token {
		case TokenInt:
			item := token.IntValue()
			result = append(result, item)

		case TokenBraceRight:
			finish = true

		case TokenEOF:
			err := token.Errorf("unexpected EOF").
				With("expect '}' to close")
			return nil, err

		default:
			err := token.Errorf("unexpected token type").
				With("expect integer or '}', got %s", token.Token)
			return nil, err
		}
	}

	return result, nil
}

func (p *Parser) parseNumberList() ([]ContextItem[int64], error) {
	if _, err := p.expectToken(TokenBraceLeft); err != nil {
		return nil, err
	}

	return p.parseNumbersInList()
}

func (p *Parser) parseCaseMemoryBlock(keyword *Element, item *TestCase) error {
	token, err := p.nextToken()
	if err != nil {
		return err
	}

	if token.Token == TokenIdentifier {
		if token.ValueString != KeywordAt {
			err := token.Errorf("syntax error").
				With("use 'at' to specify memory start address")
			return err
		}

		addrToken, err := p.nextToken()
		if err != nil {
			return err
		}

		if addrToken.Token != TokenInt {
			err := addrToken.Errorf("syntax error").
				With("expect integer as address")
			return err
		}

		if addrToken.ValueNegative {
			err := addrToken.Errorf("invalid address").
				With("address MUST BE positive")
			return err
		}

		item.MemoryAt = addrToken.UintValue()
		token, err = p.nextToken()
		if err != nil {
			return err
		}
	}

	if token.Token != TokenBraceLeft {
		err := token.Errorf("syntax error").
			With("expect '{' to start memory block")
		return err
	}

	memory, err := p.parseNumbersInList()
	if err != nil {
		return err
	}

	item.Memory = NewContextItem(memory, keyword.Context)
	return nil
}

func (p *Parser) parseCaseParameters(keyword *Element, item *TestCase) (bool, error) {
	switch keyword.ValueString {
	case FieldInput:
		input, err := p.parseNumberList()
		if err != nil {
			return false, err
		}

		item.Input = NewContextItem(input, keyword.Context)

	case FieldOutput:
		output, err := p.parseNumberList()
		if err != nil {
			return false, err
		}

		item.Output = NewContextItem(output, keyword.Context)

	case FieldMemory:
		err := p.parseCaseMemoryBlock(keyword, item)
		if err != nil {
			return false, err
		}

	default:
		err := keyword.Errorf("unknown field name").
			With("use one of 'input', 'output', 'memory', got '%s'", keyword.ValueString)
		return false, err
	}

	return false, nil
}

func (p *Parser) parseCase(keyword *Element, item *TestDriverItem) error {
	caseItem := &TestCase{}
	token, err := p.nextToken()
	if err != nil {
		return err
	}

	if token.Token == TokenIdentifier {
		// case with name
		caseItem.Name = token.StringValue()
		token, err = p.nextToken()
		if err != nil {
			return err
		}
	}

	if token.Token != TokenBraceLeft {
		err := token.Errorf("unexpected token found").
			With("expect '{' here, got %s", token.Token)
		return err
	}

	finish := false
	for !finish {
		token, err := p.nextToken()
		if err != nil {
			return err
		}

		switch token.Token {
		case TokenIdentifier:
			finish, err = p.parseCaseParameters(token, caseItem)
			if err != nil {
				return err
			}

		case TokenBraceRight:
			finish = true

		case TokenEOF:
			err := token.Errorf("unexpected EOF").
				With("expect '}' to close")
			return err

		default:
			err := token.Errorf("unexpected token type").
				With("expect identifier or '}', got %s", token.Token)
			return err
		}

	}

	if caseItem.Name.Context == nil {
		name := fmt.Sprintf("test case %d", len(item.Tests)+1)
		caseItem.Name = NewContextItem(name, keyword.Context)
	}

	item.Tests = append(item.Tests, NewContextItem(*caseItem, keyword.Context))
	return nil
}

func checkRequiredFirstSection(required map[string]bool, elem *Element) error {
	allFalse := true
	for _, v := range required {
		if v {
			allFalse = false
			break
		}
	}

	var err error
	if allFalse && elem.ValueString != SectionScript {
		err = elem.Errorf("wrong section layout").
			With("first section must be '%s', got '%s'", SectionScript, elem.ValueString)
	}

	return err
}

func (p *Parser) Parse() (*TestDriverItem, error) {
	item := &TestDriverItem{
		Tests: make([]ContextItem[TestCase], 0, 16),
	}

	sectionAppearances := map[string]bool{
		SectionScript: false,
		SectionInit:   false,
		SectionCase:   false,
	}

	var token *Element
	var err error

	for {
		token, err = p.nextToken()
		if err != nil {
			break
		}

		if token.Token == TokenEOF {
			break
		}

		if token.Token != TokenIdentifier {
			err = token.Errorf("invalid identifier").
				With("expect identifier here, got %s", token.Token)
			break
		}

		switch token.ValueString {
		case SectionScript:
			err = p.parseScript(item)
			sectionAppearances[SectionScript] = true

		case SectionInit:
			if err = checkRequiredFirstSection(sectionAppearances, token); err != nil {
				break
			}

			item.Init.Context = token.Context
			err = p.parseInit(item)
			sectionAppearances[SectionInit] = true

		case SectionCase:
			if err = checkRequiredFirstSection(sectionAppearances, token); err != nil {
				break
			}

			err = p.parseCase(token, item)
			sectionAppearances[SectionCase] = true

		default:
			err = token.Errorf("unknown section").
				With("unknown section '%s'", token.ValueString)
		}

		if err != nil {
			break
		}
	}

	if err != nil {
		return nil, err
	}

	if token.Token == TokenEOF {
		for _, section := range requiredSections {
			if !sectionAppearances[section] {
				err := context.NewError(token.Context, "missing required section").
					With("missing required section '%s'", section)
				return nil, err
			}
		}
	}

	return item, nil
}

func Parse(filename string, data []byte) (*TestDriverItem, error) {
	parser := NewParserWith(filename, data)
	return parser.Parse()
}
