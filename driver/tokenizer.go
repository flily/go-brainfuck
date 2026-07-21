package driver

import (
	"slices"
	"strings"

	"github.com/flily/go-brainfuck/context"
)

type Token int

const (
	TokenEOF Token = iota
	TokenIdentifier
	TokenInt
	TokenBoolean
	TokenBraceLeft
	TokenBraceRight
)

var tokenNameMap = map[Token]string{
	TokenEOF:        "EOF",
	TokenIdentifier: "IDENTIFIER",
	TokenInt:        "INT",
	TokenBoolean:    "BOOLEAN",
	TokenBraceLeft:  "BRACE-LEFT",
	TokenBraceRight: "BRACE-RIGHT",
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

func (e *Element) Errorf(format string, args ...any) *context.Diagnostic {
	err := context.NewError(e.Context, format, args...)
	return err
}

func (e *Element) StringValue() ContextItem[string] {
	return NewContextItem(e.ValueString, e.Context)
}

func (e *Element) IntValue() ContextItem[int64] {
	var v int64
	if e.ValueNegative {
		v = -int64(e.ValueUint)
	} else {
		v = int64(e.ValueUint)
	}

	return NewContextItem(v, e.Context)
}

func (e *Element) UintValue() ContextItem[uint64] {
	return NewContextItem(e.ValueUint, e.Context)
}

func (e *Element) BoolValue() ContextItem[bool] {
	var v bool
	switch strings.ToLower(e.ValueString) {
	case "true", "yes":
		v = true

	default:
	}

	return NewContextItem(v, e.Context)
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

func (t *Tokenizer) scanWord() (string, *context.Context) {
	start := t.cursor.State()
	i := 0
	for {
		r, eol, eof := t.cursor.Peek(i)
		if eol || eof {
			break
		}

		if r == ' ' || r == '\t' {
			break
		}

		i++
	}

	finish := t.cursor.PeekState(i)
	t.cursor.SetState(finish)

	content, ctx := t.cursor.FinishWith(start, finish)
	return content, ctx
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

func (t *Tokenizer) scanUnsignedNumber(startIndex int, negative bool) (*Element, error) {
	start := t.cursor.State()

	value := uint64(0)
	invalidFormat := false
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
			invalidFormat = true
		}

		if found {
			break
		}

		i++
	}

	finish := t.cursor.PeekState(i)
	t.cursor.SetState(finish)

	content, ctx := t.cursor.FinishWith(start, finish)
	if invalidFormat {
		err := context.NewError(ctx, "invalid number format '%s'", content).
			With("should be char [0-9] or underscore '_'")
		return nil, err
	}

	return NewElement(TokenInt, content, ctx).Int(value, negative), nil
}

func (t *Tokenizer) scanHexadecimalNumber(startIndex int, negative bool) (*Element, error) {
	start := t.cursor.State()

	value := uint64(0)
	invalidFormat := false
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

		case r == ' ' || r == '\t':
			found = true

		default:
			invalidFormat = true
		}

		if found {
			break
		}

		i++
	}

	finish := t.cursor.PeekState(i)
	t.cursor.SetState(finish)

	content, ctx := t.cursor.FinishWith(start, finish)
	if invalidFormat {
		err := context.NewError(ctx, "invalid number format '%s'", content).
			With("hexadecimal number should be 0x[0-9a-fA-F]+")
		return nil, err
	}

	return NewElement(TokenInt, content, ctx).Int(value, negative), nil
}

func (t *Tokenizer) scanOctalNumber(startIndex int, negative bool) (*Element, error) {
	start := t.cursor.State()

	value := uint64(0)
	invalidFormat := false
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

		case r == ' ' || r == '\t':
			found = true

		default:
			invalidFormat = true
		}

		if found {
			break
		}

		i++
	}

	finish := t.cursor.PeekState(i)
	t.cursor.SetState(finish)

	content, ctx := t.cursor.FinishWith(start, finish)
	elem := NewElement(TokenInt, content, ctx).Int(value, negative)

	if invalidFormat {
		err := elem.Errorf("invalid number format '%s'", content).
			With("octal number should be 0[0-7]+")
		return nil, err
	}

	return elem, nil
}

func (t *Tokenizer) scanPositiveNumber(startIndex int, negative bool) (*Element, error) {
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
		content, ctx := t.scanWord()
		err := context.NewError(ctx, "invalid number format '%s'", content).
			With("hexadecimal number should be 0x[0-9a-fA-F]+, octal number should be 0[0-7]+")
		return nil, err
	}
}

func (t *Tokenizer) scanQuotedIdentifier() (*Element, error) {
	start := t.cursor.PeekState(1)
	i := 1
	closed := false
	for {
		r, eol, eof := t.cursor.Peek(i)
		if eol || eof {
			break
		}

		if r == '"' {
			closed = true
			break
		}

		i++
	}

	finish := t.cursor.PeekState(i)
	t.cursor.SetState(finish)
	content, ctx := t.cursor.FinishWith(start, finish)
	if !closed {
		_, ctx := t.cursor.CurrentChar()
		err := context.NewError(ctx, "unclosed quoted identifier").
			With("expect closing quote '\"' here")
		return nil, err
	}

	t.cursor.Skip(1) // skip closing quote
	return NewElement(TokenIdentifier, content, ctx), nil
}

func (t *Tokenizer) scanNegativeInteger() (*Element, error) {
	elem, err := t.scanPositiveNumber(1, true)
	return elem, err
}

func (t *Tokenizer) scanSymbols(r rune, ctx *context.Context) (*Element, error) {
	var err error
	var elem *Element
	switch r {
	case '{':
		elem = NewElement(TokenBraceLeft, "{", ctx)

	case '}':
		elem = NewElement(TokenBraceRight, "}", ctx)

	default:
		err = context.NewError(ctx, "unknown charactor found").
			With("unknown charactor '%c' (0x%x)", r, r)
	}

	if err == nil {
		t.cursor.Skip(1)
	}

	return elem, err
}

func (t *Tokenizer) Next() (*Element, error) {
	if t.cursor.EOF() {
		eofCtx := t.cursor.EOFContext()
		return NewElement(TokenEOF, "", eofCtx), nil
	}

	t.cursor.SkipWhitespace()

	var err error
	var result *Element
	r, ctx := t.cursor.CurrentChar()
	switch {
	case ('a' <= r && r <= 'z') || ('A' <= r && r <= 'Z') || r == '_':
		result = t.scanIdentifier(TokenIdentifier, 0)

	case '0' <= r && r <= '9':
		result, err = t.scanPositiveNumber(0, false)

	case r == '"':
		result, err = t.scanQuotedIdentifier()

	case r == '-':
		result, err = t.scanNegativeInteger()

	default:
		result, err = t.scanSymbols(r, ctx)
	}

	return result, err
}
