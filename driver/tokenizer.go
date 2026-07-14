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
	Token   Token
	Content string
	Context *context.Context
}

func NewElement(token Token, content string, ctx *context.Context) *Element {
	e := &Element{
		Token:   token,
		Content: content,
		Context: ctx,
	}

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

func (t *Tokenizer) scanIdentifier(token Token) *Element {
	start := t.cursor.State()
	i := 0
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
	i := 0
	for {
		r, eol, eof := t.cursor.Peek(i)
		if eol || eof {
			break
		}

		switch {
		case '0' <= r && r <= '9':
		case r == '_':

		default:
			break
		}

		i++
	}

	finish := t.cursor.PeekState(i)
	t.cursor.SetState(finish)

	content, ctx := t.cursor.FinishWith(start, finish)
	return NewElement(TokenUint, content, ctx)
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
		result = t.scanIdentifier(TokenIdentifier)

	case '0' <= r && r <= '9':
		result = t.scanUnsignedNumber()

	}

	return result, nil
}
