package driver

import (
	"fmt"

	"github.com/flily/go-brainfuck/config"
	"github.com/flily/go-brainfuck/context"
)

type Parser struct {
	tokenizer *Tokenizer
}

func NewParser(filename string, data []byte) *Parser {
	file := context.ReadFileData(filename, data)
	tokenizer := NewTokenizer(file)

	p := &Parser{
		tokenizer: tokenizer,
	}

	return p
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
	switch name.Value {
	case FieldMemorySize:
		item.Init.MemorySize = value.UintValue()

	case FieldStackSize:
		item.Init.StackSize = value.UintValue()

	case FieldWord:
		var unitType config.MemoryUnitType
		err := unitType.Set(value.ValueString)
		if err != nil {
			return value.Errorf("invalid memory unit type '%s'", value.ValueString)
		}
		item.Init.WordType = NewContextItem(unitType, value.Context)
	}
	return nil
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

func (p *Parser) parseCaseParameter(item *TestCase) (bool, error) {
	token, err := p.nextToken()
	if err != nil {
		return true, err
	}

	stop := false
	switch token.Token {
	case TokenIdentifier:

	case TokenBraceRight:
		stop = true

	case TokenEOF:
		err = token.Errorf("unexpected end of file").
			With("expect '}' here, got %s", token.Token)

	default:
		err = token.Errorf("unexpected token type").
			With("expect identifier or '}', got %s", token.Token)
	}
	return stop, err
}

func (p *Parser) parseCase(keyword *Element, item *TestDriverItem) error {
	caseItem := &TestCase{}
	token, err := p.nextToken()
	if err != nil {
		return err
	}

	if token.Token == TokenIdentifier {
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

	for {
		stop, err := p.parseCaseParameter(caseItem)
		if err != nil {
			return err
		}

		if stop {
			break
		}
	}

	if caseItem.Name.Context == nil {
		name := fmt.Sprintf("test case %d", len(item.Tests)+1)
		caseItem.Name = NewContextItem(name, keyword.Context)
	}

	item.Tests = append(item.Tests, *caseItem)
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
		Tests: make([]TestCase, 0, 16),
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
			err = token.Errorf("expect identifier here, got %s", token.Token)
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
	parser := NewParser(filename, data)
	return parser.Parse()
}
