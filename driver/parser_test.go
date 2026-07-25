package driver

import (
	"testing"

	"strings"

	"github.com/flily/go-brainfuck/config"
)

const (
	testScriptFilename = "test.bft" // brainfuck test
)

func testParse(input string) (*TestDriverItem, error) {
	return Parse(testScriptFilename, []byte(input))
}

func checkError(t *testing.T, input string, expected string) {
	t.Helper()

	item, err := testParse(input)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	if item != nil {
		t.Fatalf("expected nil item, got: %+v", item)
	}

	if err.Error() != expected {
		t.Fatalf("expected error:\n%s\n\ngot:\n%s", expected, err.Error())
	}
}

func checkOK(t *testing.T, input string, expected *TestDriverItem) {
	t.Helper()

	item, err := testParse(input)
	if err != nil {
		t.Fatalf("expected no error, got: %s", err.Error())
	}

	if !expected.Equal(item) {
		t.Fatalf("expected:\n%+v\n\ngot:\n%+v", expected, item)
	}
}

func TestParserErrorEmptyContent(t *testing.T) {
	input := strings.Join([]string{
		"",
	}, "\n")

	expected := strings.Join([]string{
		"test.bft:1:1: error: missing required section",
		"    1 | <EOF>",
		"      | ^^^^^",
		"      | missing required section 'script'",
	}, "\n")

	checkError(t, input, expected)
}

func TestParserErrorStartWithInitSection(t *testing.T) {
	input := strings.Join([]string{
		"init {",
		"    memory-size: 1024",
		"}",
	}, "\n")

	expected := strings.Join([]string{
		"test.bft:1:1: error: wrong section layout",
		"    1 | init {",
		"      | ^^^^",
		"      | first section must be 'script', got 'init'",
	}, "\n")

	checkError(t, input, expected)
}

func TestParserErrorStartWithCaseSection(t *testing.T) {
	input := strings.Join([]string{
		"case example {",
		"}",
	}, "\n")

	expected := strings.Join([]string{
		"test.bft:1:1: error: wrong section layout",
		"    1 | case example {",
		"      | ^^^^",
		"      | first section must be 'script', got 'case'",
	}, "\n")

	checkError(t, input, expected)
}

func TestParserErrorWithUnknownSection(t *testing.T) {
	input := strings.Join([]string{
		`script "hello.bf"`,
		`lorem ipsum {`,
		`}`,
	}, "\n")

	expected := strings.Join([]string{
		"test.bft:2:1: error: unknown section",
		"    2 | lorem ipsum {",
		"      | ^^^^^",
		"      | unknown section 'lorem'",
	}, "\n")

	checkError(t, input, expected)
}

func TestParserErrorWithMissingRequiredSectionCase(t *testing.T) {
	input := strings.Join([]string{
		`script "hello.bf"`,
		`init {}`,
	}, "\n")

	expected := strings.Join([]string{
		"test.bft:2:8: error: missing required section",
		"    2 | init {}<EOF>",
		"      |        ^^^^^",
		"      |        missing required section 'case'",
	}, "\n")

	checkError(t, input, expected)
}

func TestParserErrorWithMissingRequiredSectionInit(t *testing.T) {
	input := strings.Join([]string{
		`script "hello.bf"`,
		`case {}`,
		``,
	}, "\n")

	expected := strings.Join([]string{
		"test.bft:3:1: error: missing required section",
		"    3 | <EOF>",
		"      | ^^^^^",
		"      | missing required section 'init'",
	}, "\n")

	checkError(t, input, expected)
}

func TestParserScriptNameOnly(t *testing.T) {
	input := strings.Join([]string{
		`script "path/to/script.bf"`,
		`init {}`,
		`case {}`,
	}, "\n")

	expected := &TestDriverItem{
		ScriptName: NewContextItem("path/to/script.bf", nil),
		Tests: []TestCase{
			{Name: NewContextItem("test case 1", nil)},
		},
	}

	checkOK(t, input, expected)
}

func TestParserErrorWrongFieldNameInInit(t *testing.T) {
	input := strings.Join([]string{
		`script "path/to/script.bf"`,
		`init {`,
		`    memory-size 1024`,
		`    lorem-ipsum 123`,
		`}`,
		`case {}`,
	}, "\n")

	expected := strings.Join([]string{
		"test.bft:4:5: error: unknown field name",
		"    4 |     lorem-ipsum 123",
		"      |     ^^^^^^^^^^^",
		"      |     unknown field name 'lorem-ipsum'",
	}, "\n")

	checkError(t, input, expected)
}

func TestParserWithInitFields(t *testing.T) {
	input := strings.Join([]string{
		`script "path/to/script.bf"`,
		`init {`,
		`    memory-size   1024`,
		`    stack-size    256`,
		`    word          uint8`,
		`    eof-value     -1`,
		`    ignore-eof    yes`,
		`    raise-eof     no`,
		`}`,
		`case {}`,
	}, "\n")

	expected := &TestDriverItem{
		ScriptName: NewContextItem("path/to/script.bf", nil),
		Init: InitParameters{
			MemorySize: NewContextItem(uint64(1024), nil),
			StackSize:  NewContextItem(uint64(256), nil),
			WordType:   NewContextItem(config.MemoryUnitTypeUint8, nil),
			EOFValue:   NewContextItem(int64(-1), nil),
			IgnoreEOF:  NewContextItem(true, nil),
			RaiseEOF:   NewContextItem(false, nil),
		},
		Tests: []TestCase{
			{Name: NewContextItem("test case 1", nil)},
		},
	}

	checkOK(t, input, expected)
}
