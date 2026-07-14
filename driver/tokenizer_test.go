package driver

import (
	"testing"

	"strings"

	"github.com/flily/go-brainfuck/context"
)

type testTokenizerCase struct {
	t     *testing.T
	input string
}

func newTokenizerCase(t *testing.T, input string) *testTokenizerCase {
	c := &testTokenizerCase{
		t:     t,
		input: input,
	}

	return c
}

func (t *testTokenizerCase) scan() (*Element, error) {
	file := context.ReadFileString("test.txt", t.input)
	tokenizer := NewTokenizer(file)
	elem, err := tokenizer.Next()

	return elem, err
}

func (t *testTokenizerCase) Ok(token Token, content string, position string) {
	elem, err := t.scan()
	if err != nil {
		t.t.Fatalf("expect scan token without error, got:\n%s", err)
		return
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
}

func TestTokenizerScanIdentifier(t *testing.T) {
	input := "lorem ipsum"
	position := strings.Join([]string{
		"    1 | lorem ipsum",
		"      | ^^^^^",
		"      | here",
	}, "\n")

	newTokenizerCase(t, input).
		Ok(TokenIdentifier, "lorem", position)
}

func TestTokenizerScanBoolean(t *testing.T) {
	input := "false ipsum"
	position := strings.Join([]string{
		"    1 | false ipsum",
		"      | ^^^^^",
		"      | here",
	}, "\n")

	newTokenizerCase(t, input).
		Ok(TokenBoolean, "false", position)
}

func TestTokenizerScanUnsignedNumber(t *testing.T) {
	input := "42 ipsum"
	position := strings.Join([]string{
		"    1 | 42 ipsum",
		"      | ^^",
		"      | here",
	}, "\n")
	newTokenizerCase(t, input).
		Ok(TokenUint, "42", position)
}
