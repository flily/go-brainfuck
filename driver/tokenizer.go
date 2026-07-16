package driver

import (
	"slices"

	"github.com/flily/go-brainfuck/context"
)

type Token int

const (
	TokenEOF Token = iota
	TokenIdentifier
	TokenInt
	TokenBoolean
	TokenBracketLeft
	TokenBracketRight
)

var tokenNameMap = map[Token]string{
	TokenEOF:          "EOF",
	TokenIdentifier:   "IDENTIFIER",
	TokenInt:          "INT",
	TokenBoolean:      "BOOLEAN",
	TokenBracketLeft:  "BRACKET-LEFT",
	TokenBracketRight: "BRACKET-RIGHT",
}

func (t Token) String() string {
	if s, ok := tokenNameMap[t]; ok {
		return s
	}

	return "unknown"
}

var booleanWords = []string{
	"false",
	"true",
	"no",
	"yes",
}

type Element struct {
	Token         Token
	ValueNegative bool
	ValueUint     uint64
	ValueString   string
	Context       *context.Context
}

func NewElement(token Token, content string, ctx *context.Context) *Element {
	e := &Element{
		Token:       token,
		ValueString: content,
		Context:     ctx,
	}

	return e
}

func (e *Element) Int(v uint64, neg bool) *Element {
	e.ValueNegative = neg
	e.ValueUint = v
	return e
}

type Tokenizer struct {
	file   *context.FileContext
	cursor *context.Cursor
}

func NewTokenizer(file *context.FileContext) *Tokenizer {
	t := &Tokenizer{
		file:   file,
		cursor: context.NewCursor(file),
	}

	return t
}

func (t *Tokenizer) scanIdentifier(token Token, startIndex int) *Element {
	start := t.cursor.State()
	i := startIndex
	for {
		r, eol, eof := t.cursor.Peek(i)
		if eol || eof {
			break
		}

		found := false
		switch {
		case ('a' <= r && r <= 'z') || ('A' <= r && r <= 'Z'):
		case '0' <= r && r <= '9':
		case r == '_' || r == '-':

		default:
			found = true
		}

		if found {
			break
		}

		i++
	}

	finish := t.cursor.PeekState(i)
	t.cursor.SetState(finish)

	targetToken := token
	content, ctx := t.cursor.FinishWith(start, finish)
	if slices.Contains(booleanWords, content) {
		targetToken = TokenBoolean
	}

	return NewElement(targetToken, content, ctx)
}

func (t *Tokenizer) scanUnsignedNumber(startIndex int, negative bool) *Element {
	start := t.cursor.State()

	value := uint64(0)
	i := startIndex
	for {
		r, eol, eof := t.cursor.Peek(i)
		if eol || eof {
			break
		}

		found := false
		switch {
		case '0' <= r && r <= '9':
			value = value*10 + uint64(r-'0')

		case r == '_':

		case r == ' ' || r == '\t':
			found = true

		default:
			return t.scanIdentifier(TokenIdentifier, i)
		}

		if found {
			break
		}

		i++
	}

	finish := t.cursor.PeekState(i)
	t.cursor.SetState(finish)

	content, ctx := t.cursor.FinishWith(start, finish)
	return NewElement(TokenInt, content, ctx).Int(value, negative)
}

func (t *Tokenizer) scanHexadecimalNumber(startIndex int, negative bool) *Element {
	start := t.cursor.State()

	value := uint64(0)
	i := startIndex
	for {
		r, eol, eof := t.cursor.Peek(i)
		if eol || eof {
			break
		}

		found := false
		switch {
		case '0' <= r && r <= '9':
			value <<= 4
			value += uint64(r - '0')

		case 'a' <= r && r <= 'f':
			value <<= 4
			value += uint64(r - 'a' + 10)

		case 'A' <= r && r <= 'F':
			value <<= 4
			value += uint64(r - 'A' + 10)

		case r == '_':

		default:
			found = true
		}

		if found {
			break
		}

		i++
	}

	finish := t.cursor.PeekState(i)
	t.cursor.SetState(finish)

	content, ctx := t.cursor.FinishWith(start, finish)
	return NewElement(TokenInt, content, ctx).Int(value, negative)
}

func (t *Tokenizer) scanOctalNumber(startIndex int, negative bool) *Element {
	start := t.cursor.State()

	value := uint64(0)
	i := startIndex
	for {
		r, eol, eof := t.cursor.Peek(i)
		if eol || eof {
			break
		}

		found := false
		switch {
		case '0' <= r && r <= '7':
			value <<= 3
			value += uint64(r - '0')

		case r == '_':

		default:
			found = true
		}

		if found {
			break
		}

		i++
	}

	finish := t.cursor.PeekState(i)
	t.cursor.SetState(finish)

	content, ctx := t.cursor.FinishWith(start, finish)
	return NewElement(TokenInt, content, ctx).Int(value, negative)
}

func (t *Tokenizer) scanPositiveNumber(startIndex int, negative bool) *Element {
	r0, _, _ := t.cursor.Peek(startIndex)
	if r0 != '0' {
		return t.scanUnsignedNumber(startIndex, negative)
	}

	r1, eol, eof := t.cursor.Peek(startIndex + 1)
	if eol || eof {
		// '0'
		return t.scanUnsignedNumber(startIndex, negative)
	}

	switch r1 {
	case 'x', 'X':
		return t.scanHexadecimalNumber(startIndex+2, negative)
	case '1', '2', '3', '4', '5', '6', '7':
		return t.scanOctalNumber(startIndex+1, negative)

	default:
		return t.scanIdentifier(TokenIdentifier, 0)
	}
}

func (t *Tokenizer) scanNegativeInteger() *Element {
	elem := t.scanPositiveNumber(1, true)
	return elem
}

func (t *Tokenizer) Next() (*Element, error) {
	if t.cursor.EOF() {
		eofCtx := t.cursor.EOFContext()
		return NewElement(TokenEOF, "", eofCtx), nil
	}

	t.cursor.SkipWhitespace()
	var result *Element
	r, _ := t.cursor.CurrentChar()
	switch {
	case ('a' <= r && r <= 'z') || ('A' <= r && r <= 'Z') || r == '_':
		result = t.scanIdentifier(TokenIdentifier, 0)

	case '0' <= r && r <= '9':
		result = t.scanPositiveNumber(0, false)

	case r == '-':
		result = t.scanNegativeInteger()

	}

	return result, nil
}
