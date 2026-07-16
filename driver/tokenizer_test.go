package driver

import (
	"testing"

	"strings"

	"github.com/flily/go-brainfuck/context"
)

func TestTokenString(t *testing.T) {
	cases := []struct {
		token    Token
		expected string
	}{
		{TokenEOF, "EOF"},
		{TokenIdentifier, "IDENTIFIER"},
		{TokenInt, "INT"},
		{TokenBoolean, "BOOLEAN"},
		{TokenBracketLeft, "BRACKET-LEFT"},
		{TokenBracketRight, "BRACKET-RIGHT"},
		{-1, "unknown"},
	}

	for _, c := range cases {
		got := c.token.String()
		if got != c.expected {
			t.Errorf("wrong token string, expect %s got %s", c.expected, got)
		}
	}
}

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
	t.t.Helper()

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
	t.t.Helper()

	t.check(token, content, position)
	return t
}

func (t *testTokenizerCase) CheckUint(content string, value uint64, position string) *testTokenizerCase {
	t.t.Helper()

	elem := t.check(TokenInt, content, position)
	if elem.ValueUint != value {
		t.t.Errorf("wrong token value, expect %d got %d", value, elem.ValueUint)
		t.t.Fatalf("\n%s", elem.Context.HighlightText("expect %d", value))
	}

	return t
}

func (t *testTokenizerCase) CheckInt(content string, value int64, position string) *testTokenizerCase {
	t.t.Helper()

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

func (t *testTokenizerCase) CheckEOF(position string) *testTokenizerCase {
	t.t.Helper()

	t.check(TokenEOF, "", position)
	return t
}

func TestTokenizerScanIdentifier(t *testing.T) {
	input := "lorem ipsum dolor sit-amet"
	position1 := strings.Join([]string{
		"    1 | lorem ipsum dolor sit-amet",
		"      | ^^^^^",
		"      | here",
	}, "\n")
	position2 := strings.Join([]string{
		"    1 | lorem ipsum dolor sit-amet",
		"      |       ^^^^^",
		"      |       here",
	}, "\n")
	position3 := strings.Join([]string{
		"    1 | lorem ipsum dolor sit-amet",
		"      |             ^^^^^",
		"      |             here",
	}, "\n")
	position4 := strings.Join([]string{
		"    1 | lorem ipsum dolor sit-amet",
		"      |                   ^^^^^^^^",
		"      |                   here",
	}, "\n")
	position5 := strings.Join([]string{
		"    1 | lorem ipsum dolor sit-amet<EOF>",
		"      |                           ^^^^^",
		"      |                           here",
	}, "\n")

	newTokenizerCase(t, input).
		Check(TokenIdentifier, "lorem", position1).
		Check(TokenIdentifier, "ipsum", position2).
		Check(TokenIdentifier, "dolor", position3).
		Check(TokenIdentifier, "sit-amet", position4).
		CheckEOF(position5)
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
		CheckEOF(position3)
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
		Check(TokenInt, "42", position1).
		Check(TokenInt, "030", position2).
		Check(TokenInt, "0x30", position3).
		CheckEOF(position4)
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
		CheckEOF(position4)
}

func TestTokenizerScanUnsignedNumberOnlyZero(t *testing.T) {
	input := "0"
	position1 := strings.Join([]string{
		"    1 | 0",
		"      | ^",
		"      | here",
	}, "\n")
	position2 := strings.Join([]string{
		"    1 | 0<EOF>",
		"      |  ^^^^^",
		"      |  here",
	}, "\n")

	newTokenizerCase(t, input).
		CheckUint("0", 0, position1).
		CheckEOF(position2)
}

func TestTokenizerScanNumberWithUnderscore(t *testing.T) {
	input := "1_000_000 0xdead_BEEF 0770_660"
	position1 := strings.Join([]string{
		"    1 | 1_000_000 0xdead_BEEF 0770_660",
		"      | ^^^^^^^^^",
		"      | here",
	}, "\n")
	position2 := strings.Join([]string{
		"    1 | 1_000_000 0xdead_BEEF 0770_660",
		"      |           ^^^^^^^^^^^",
		"      |           here",
	}, "\n")
	position3 := strings.Join([]string{
		"    1 | 1_000_000 0xdead_BEEF 0770_660",
		"      |                       ^^^^^^^^",
		"      |                       here",
	}, "\n")
	position4 := strings.Join([]string{
		"    1 | 1_000_000 0xdead_BEEF 0770_660<EOF>",
		"      |                               ^^^^^",
		"      |                               here",
	}, "\n")

	newTokenizerCase(t, input).
		CheckUint("1_000_000", 1000_000, position1).
		CheckUint("0xdead_BEEF", 0xdead_beef, position2).
		CheckUint("0770_660", 0770_660, position3).
		CheckEOF(position4)
}

func TestTokenizerScanNegativeNumber(t *testing.T) {
	input := "-42 -0x42 -077"
	position1 := strings.Join([]string{
		"    1 | -42 -0x42 -077",
		"      | ^^^",
		"      | here",
	}, "\n")
	position2 := strings.Join([]string{
		"    1 | -42 -0x42 -077",
		"      |     ^^^^^",
		"      |     here",
	}, "\n")
	position3 := strings.Join([]string{
		"    1 | -42 -0x42 -077",
		"      |           ^^^^",
		"      |           here",
	}, "\n")
	position4 := strings.Join([]string{
		"    1 | -42 -0x42 -077<EOF>",
		"      |               ^^^^^",
		"      |               here",
	}, "\n")

	newTokenizerCase(t, input).
		CheckInt("-42", -42, position1).
		CheckInt("-0x42", -0x42, position2).
		CheckInt("-077", -077, position3).
		CheckEOF(position4)
}

func TestTokenizerScanIdentifierStartsWithNumber(t *testing.T) {
	input := "4ever 0zero"
	position1 := strings.Join([]string{
		"    1 | 4ever 0zero",
		"      | ^^^^^",
		"      | here",
	}, "\n")
	position2 := strings.Join([]string{
		"    1 | 4ever 0zero",
		"      |       ^^^^^",
		"      |       here",
	}, "\n")
	position3 := strings.Join([]string{
		"    1 | 4ever 0zero<EOF>",
		"      |            ^^^^^",
		"      |            here",
	}, "\n")

	newTokenizerCase(t, input).
		Check(TokenIdentifier, "4ever", position1).
		Check(TokenIdentifier, "0zero", position2).
		CheckEOF(position3)
}
