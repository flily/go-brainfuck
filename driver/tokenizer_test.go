package driver

import (
	"testing"

	"strings"

	"github.com/flily/go-brainfuck/context"
)

type testTokenizerCase struct {
	t         *testing.T
	input     string
	tokenizer *Tokenizer
}

func newTokenizerCase(t *testing.T, input string) *testTokenizerCase {
	file := context.ReadFileString("test.txt", input)
	tokenizer := NewTokenizer(file)

	c := &testTokenizerCase{
		t:         t,
		input:     input,
		tokenizer: tokenizer,
	}

	return c
}

func (t *testTokenizerCase) scan() (*Element, error) {
	elem, err := t.tokenizer.Next()
	return elem, err
}

func (t *testTokenizerCase) FindOk(token Token, content string, position string) *testTokenizerCase {
	elem, err := t.scan()
	if err != nil {
		t.t.Fatalf("expect scan token without error, got:\n%s", err)
		return t
	}

	ctx := elem.Context
	if elem.Token != token {
		message := ctx.HighlightText("expect token %s", token)
		t.t.Errorf("wrong token type, expect %s got %s", token, elem.Token)
		t.t.Fatalf("\n%s", message)
	}

	if elem.Content != content {
		message := ctx.HighlightText("expect '%s'", content)
		t.t.Errorf("wrong token, expect '%s' got '%s'", content, elem.Content)
		t.t.Fatalf("\n%s", message)
	}

	p := ctx.HighlightText("here")
	if p != position {
		t.t.Fatalf("got wrong position, expect\n%s\ngot\n%s", position, p)
	}

	return t
}

func TestTokenizerScanIdentifier(t *testing.T) {
	input := "lorem ipsum"
	position1 := strings.Join([]string{
		"    1 | lorem ipsum",
		"      | ^^^^^",
		"      | here",
	}, "\n")
	position2 := strings.Join([]string{
		"    1 | lorem ipsum",
		"      |       ^^^^^",
		"      |       here",
	}, "\n")
	position3 := strings.Join([]string{
		"    1 | lorem ipsum<EOF>",
		"      |            ^^^^^",
		"      |            here",
	}, "\n")

	newTokenizerCase(t, input).
		FindOk(TokenIdentifier, "lorem", position1).
		FindOk(TokenIdentifier, "ipsum", position2).
		FindOk(TokenEOF, "", position3)
}

func TestTokenizerScanBoolean(t *testing.T) {
	input := "false ipsum"
	position1 := strings.Join([]string{
		"    1 | false ipsum",
		"      | ^^^^^",
		"      | here",
	}, "\n")
	position2 := strings.Join([]string{
		"    1 | false ipsum",
		"      |       ^^^^^",
		"      |       here",
	}, "\n")
	position3 := strings.Join([]string{
		"    1 | false ipsum<EOF>",
		"      |            ^^^^^",
		"      |            here",
	}, "\n")

	newTokenizerCase(t, input).
		FindOk(TokenBoolean, "false", position1).
		FindOk(TokenIdentifier, "ipsum", position2).
		FindOk(TokenEOF, "", position3)
}

func TestTokenizerScanUnsignedNumber(t *testing.T) {
	input := "42 ipsum"
	position1 := strings.Join([]string{
		"    1 | 42 ipsum",
		"      | ^^",
		"      | here",
	}, "\n")
	position2 := strings.Join([]string{
		"    1 | 42 ipsum",
		"      |    ^^^^^",
		"      |    here",
	}, "\n")
	position3 := strings.Join([]string{
		"    1 | 42 ipsum<EOF>",
		"      |         ^^^^^",
		"      |         here",
	}, "\n")

	newTokenizerCase(t, input).
		FindOk(TokenUint, "42", position1).
		FindOk(TokenIdentifier, "ipsum", position2).
		FindOk(TokenEOF, "", position3)
}
