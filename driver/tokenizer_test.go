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

func (t *testTokenizerCase) CheckError(message string) *testTokenizerCase {
	t.t.Helper()

	elem, err := t.scan()
	if err == nil {
		t.t.Fatalf("expect scan token with error, got:\n%v", elem)
	}

	if err.Error() != message {
		t.t.Fatalf("wrong error message, expect\n%s\ngot\n%s", message, err.Error())
	}

	return t
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
	input := "lorem ipsum6 dolor sit-amet"
	position1 := strings.Join([]string{
		"    1 | lorem ipsum6 dolor sit-amet",
		"      | ^^^^^",
		"      | here",
	}, "\n")
	position2 := strings.Join([]string{
		"    1 | lorem ipsum6 dolor sit-amet",
		"      |       ^^^^^^",
		"      |       here",
	}, "\n")
	position3 := strings.Join([]string{
		"    1 | lorem ipsum6 dolor sit-amet",
		"      |              ^^^^^",
		"      |              here",
	}, "\n")
	position4 := strings.Join([]string{
		"    1 | lorem ipsum6 dolor sit-amet",
		"      |                    ^^^^^^^^",
		"      |                    here",
	}, "\n")
	position5 := strings.Join([]string{
		"    1 | lorem ipsum6 dolor sit-amet<EOF>",
		"      |                            ^^^^^",
		"      |                            here",
	}, "\n")

	newTokenizerCase(t, input).
		Check(TokenIdentifier, "lorem", position1).
		Check(TokenIdentifier, "ipsum6", position2).
		Check(TokenIdentifier, "dolor", position3).
		Check(TokenIdentifier, "sit-amet", position4).
		CheckEOF(position5)
}

func TestTokenizerScanQuotedIdentifier(t *testing.T) {
	input := `"lorem ipsum" "dolor sit amet"`
	position1 := strings.Join([]string{
		`    1 | "lorem ipsum" "dolor sit amet"`,
		"      |  ^^^^^^^^^^^",
		"      |  here",
	}, "\n")
	position2 := strings.Join([]string{
		`    1 | "lorem ipsum" "dolor sit amet"`,
		"      |                ^^^^^^^^^^^^^^",
		"      |                here",
	}, "\n")
	position3 := strings.Join([]string{
		`    1 | "lorem ipsum" "dolor sit amet"<EOF>`,
		"      |                               ^^^^^",
		"      |                               here",
	}, "\n")

	newTokenizerCase(t, input).
		Check(TokenIdentifier, "lorem ipsum", position1).
		Check(TokenIdentifier, "dolor sit amet", position2).
		CheckEOF(position3)
}

func TestTokenizerScanQutedIdentifierWithUnclosedQuoteAtEOF(t *testing.T) {
	input := `"lorem ipsum`
	message := strings.Join([]string{
		"test.txt:1:13: error: unclosed quoted identifier",
		`    1 | "lorem ipsum<EOF>`,
		"      |             ^^^^^",
		"      |             expect closing quote '\"' here",
	}, "\n")

	newTokenizerCase(t, input).
		CheckError(message)
}

func TestTokenizerScanQutedIdentifierWithUnclosedQuoteAtEOL(t *testing.T) {
	input := strings.Join([]string{
		`"lorem ipsum`,
		"",
	}, "\n")

	message := strings.Join([]string{
		"test.txt:1:13: error: unclosed quoted identifier",
		`    1 | "lorem ipsum<EOL LF>`,
		"      |             ^^^^^^^^",
		"      |             expect closing quote '\"' here",
	}, "\n")

	newTokenizerCase(t, input).
		CheckError(message)
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

func TestTokenizerScanErrorInvalidUnsignedNumberFormat(t *testing.T) {
	input := "4ever 42"
	message := strings.Join([]string{
		"test.txt:1:1: error: invalid number format '4ever'",
		"    1 | 4ever 42",
		"      | ^^^^^",
		"      | should be char [0-9] or underscore '_'",
	}, "\n")

	newTokenizerCase(t, input).
		CheckError(message)
}

func TestTokenizerScanErrorInvalidNumberFormatAtEnd(t *testing.T) {
	input := "0zero"
	message := strings.Join([]string{
		"test.txt:1:1: error: invalid number format '0zero'",
		"    1 | 0zero",
		"      | ^^^^^",
		"      | hexadecimal number should be 0x[0-9a-fA-F]+, octal number should be 0[0-7]+",
	}, "\n")

	newTokenizerCase(t, input).
		CheckError(message)
}

func TestTokenizerScanErrorInvalidNumberFormat(t *testing.T) {
	input := "0zero 42"
	message := strings.Join([]string{
		"test.txt:1:1: error: invalid number format '0zero'",
		"    1 | 0zero 42",
		"      | ^^^^^",
		"      | hexadecimal number should be 0x[0-9a-fA-F]+, octal number should be 0[0-7]+",
	}, "\n")

	newTokenizerCase(t, input).
		CheckError(message)
}

func TestTokenizerScanMultipleLines(t *testing.T) {
	input := strings.Join([]string{
		"lorem ipsum",
		"dolor sit amet",
	}, "\n")
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
		"    2 | dolor sit amet",
		"      | ^^^^^",
		"      | here",
	}, "\n")
	position4 := strings.Join([]string{
		"    2 | dolor sit amet",
		"      |       ^^^",
		"      |       here",
	}, "\n")
	position5 := strings.Join([]string{
		"    2 | dolor sit amet",
		"      |           ^^^^",
		"      |           here",
	}, "\n")
	position6 := strings.Join([]string{
		"    2 | dolor sit amet<EOF>",
		"      |               ^^^^^",
		"      |               here",
	}, "\n")

	newTokenizerCase(t, input).
		Check(TokenIdentifier, "lorem", position1).
		Check(TokenIdentifier, "ipsum", position2).
		Check(TokenIdentifier, "dolor", position3).
		Check(TokenIdentifier, "sit", position4).
		Check(TokenIdentifier, "amet", position5).
		CheckEOF(position6)
}

func TestTokenizerScanSymbols(t *testing.T) {
	input := "{lorem}"
	position1 := strings.Join([]string{
		"    1 | {lorem}",
		"      | ^",
		"      | here",
	}, "\n")
	position2 := strings.Join([]string{
		"    1 | {lorem}",
		"      |  ^^^^^",
		"      |  here",
	}, "\n")
	position3 := strings.Join([]string{
		"    1 | {lorem}",
		"      |       ^",
		"      |       here",
	}, "\n")
	position4 := strings.Join([]string{
		"    1 | {lorem}<EOF>",
		"      |        ^^^^^",
		"      |        here",
	}, "\n")

	newTokenizerCase(t, input).
		Check(TokenBracketLeft, "{", position1).
		Check(TokenIdentifier, "lorem", position2).
		Check(TokenBracketRight, "}", position3).
		CheckEOF(position4)
}

func TestTokenizerScanErrorUnknownCharacter(t *testing.T) {
	input := "#lorem"
	message := strings.Join([]string{
		"test.txt:1:1: error: unknown charactor found",
		"    1 | #lorem",
		"      | ^",
		"      | unknown charactor '#' (0x23)",
	}, "\n")

	newTokenizerCase(t, input).
		CheckError(message)
}
