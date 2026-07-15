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

func (t *testTokenizerCase) check(token Token, content string, position string) *Element {
	elem, err := t.scan()
	if err != nil {
		t.t.Fatalf("expect scan token without error, got:\n%s", err)
	}

	ctx := elem.Context
	if elem.Token != token {
		message := ctx.HighlightText("expect token %s", token)
		t.t.Errorf("wrong token type, expect %s got %s", token, elem.Token)
		t.t.Fatalf("\n%s", message)
	}

	if elem.ValueString != content {
		message := ctx.HighlightText("expect '%s'", content)
		t.t.Errorf("wrong token, expect '%s' got '%s'", content, elem.ValueString)
		t.t.Fatalf("\n%s", message)
	}

	p := ctx.HighlightText("here")
	if p != position {
		t.t.Fatalf("got wrong position, expect\n%s\ngot\n%s", position, p)
	}

	return elem
}

func (t *testTokenizerCase) Check(token Token, content string, position string) *testTokenizerCase {
	t.check(token, content, position)
	return t
}

func (t *testTokenizerCase) CheckUint(content string, value uint64, position string) *testTokenizerCase {
	elem := t.check(TokenUint, content, position)
	if elem.ValueUint != value {
		t.t.Errorf("wrong token value, expect %d got %d", value, elem.ValueUint)
		t.t.Fatalf("\n%s", elem.Context.HighlightText("expect %d", value))
	}

	return t
}

func (t *testTokenizerCase) CheckInt(content string, value int64, position string) *testTokenizerCase {
	elem := t.check(TokenInt, content, position)

	exp := int64(elem.ValueUint)
	if elem.ValueNegative {
		exp = -exp
	}

	if exp != value {
		t.t.Errorf("wrong token value, expect %d got %d", value, exp)
		t.t.Fatalf("\n%s", elem.Context.HighlightText("expect %d", value))
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
		Check(TokenIdentifier, "lorem", position1).
		Check(TokenIdentifier, "ipsum", position2).
		Check(TokenEOF, "", position3)
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
		Check(TokenBoolean, "false", position1).
		Check(TokenIdentifier, "ipsum", position2).
		Check(TokenEOF, "", position3)
}

func TestTokenizerScanUnsignedNumber(t *testing.T) {
	input := "42 030 0x30"
	position1 := strings.Join([]string{
		"    1 | 42 030 0x30",
		"      | ^^",
		"      | here",
	}, "\n")
	position2 := strings.Join([]string{
		"    1 | 42 030 0x30",
		"      |    ^^^",
		"      |    here",
	}, "\n")
	position3 := strings.Join([]string{
		"    1 | 42 030 0x30",
		"      |        ^^^^",
		"      |        here",
	}, "\n")
	position4 := strings.Join([]string{
		"    1 | 42 030 0x30<EOF>",
		"      |            ^^^^^",
		"      |            here",
	}, "\n")

	newTokenizerCase(t, input).
		Check(TokenUint, "42", position1).
		Check(TokenUint, "030", position2).
		Check(TokenUint, "0x30", position3).
		Check(TokenEOF, "", position4)
}

func TestTokenizerScanUnsignedNumberWithValue(t *testing.T) {
	input := "42 030 0x30"
	position1 := strings.Join([]string{
		"    1 | 42 030 0x30",
		"      | ^^",
		"      | here",
	}, "\n")
	position2 := strings.Join([]string{
		"    1 | 42 030 0x30",
		"      |    ^^^",
		"      |    here",
	}, "\n")
	position3 := strings.Join([]string{
		"    1 | 42 030 0x30",
		"      |        ^^^^",
		"      |        here",
	}, "\n")
	position4 := strings.Join([]string{
		"    1 | 42 030 0x30<EOF>",
		"      |            ^^^^^",
		"      |            here",
	}, "\n")

	newTokenizerCase(t, input).
		CheckUint("42", 42, position1).
		CheckUint("030", 24, position2).
		CheckUint("0x30", 48, position3).
		Check(TokenEOF, "", position4)
}
