package driver

import (
	"slices"

	"github.com/flily/go-brainfuck/context"
)

type Token int

const (
	TokenEOF Token = iota
	TokenIdentifier
	TokenUint
	TokenInt
	TokenBoolean
	TokenBracketLeft
	TokenBracketRight
)

var tokenNameMap = map[Token]string{
	TokenEOF:          "EOF",
	TokenIdentifier:   "IDENTIFIER",
	TokenUint:         "UINT",
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

func (e *Element) Uint(v uint64) *Element {
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
		case 'a' <= r && r <= 'z':
		case 'A' <= r && r <= 'Z':
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

func (t *Tokenizer) scanUnsignedNumber() *Element {
	start := t.cursor.State()

	value := uint64(0)
	i := 0
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
	return NewElement(TokenUint, content, ctx).Uint(value)
}

func (t *Tokenizer) scanHexadecimalNumber(startIndex int) *Element {
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
	return NewElement(TokenUint, content, ctx).Uint(value)
}

func (t *Tokenizer) scanOctalNumber(startIndex int) *Element {
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
	return NewElement(TokenUint, content, ctx).Uint(value)
}

func (t *Tokenizer) scanSignedInteger() *Element {
	return nil
}

func (t *Tokenizer) scanNumber() *Element {
	r0, _, _ := t.cursor.Rune()
	if r0 != '0' {
		return t.scanUnsignedNumber()
	}

	r1, eol, eof := t.cursor.Peek(1)
	if eol || eof {
		return t.scanUnsignedNumber()
	}

	switch r1 {
	case 'x', 'X':
		return t.scanHexadecimalNumber(2)
	case '1', '2', '3', '4', '5', '6', '7':
		return t.scanOctalNumber(1)
	}

	return nil
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
		result = t.scanNumber()

	}

	return result, nil
}
