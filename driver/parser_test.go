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

func ctxNums(nums ...int64) []ContextItem[int64] {
	var result []ContextItem[int64]
	for _, n := range nums {
		result = append(result, NewContextItem(n, nil))
	}

	return result
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

func TestParserErrorWithInvalidMemoryType(t *testing.T) {
	input := strings.Join([]string{
		`script "path/to/script.bf"`,
		`init {`,
		`    memory-size 1024`,
		`    stack-size  256`,
		`    word        byte`,
		`}`,
		`case {}`,
	}, "\n")

	expected := strings.Join([]string{
		"test.bft:5:17: error: invalid parameter value",
		"    5 |     word        byte",
		"      |                 ^^^^",
		"      |                 invalid memory unit type 'byte'",
	}, "\n")

	checkError(t, input, expected)
}

func TestParserErrorWithInvalidBooleanValue(t *testing.T) {
	input := strings.Join([]string{
		`script "path/to/script.bf"`,
		`init {`,
		`    memory-size 1024`,
		`    stack-size  256`,
		`    word        uint8`,
		`    ignore-eof  maybe`,
		`}`,
		`case {}`,
	}, "\n")

	expected := strings.Join([]string{
		"test.bft:6:17: error: invalid boolean value 'maybe'",
		"    6 |     ignore-eof  maybe",
		"      |                 ^^^^^",
		"      |                 use yes/no",
	}, "\n")

	checkError(t, input, expected)
}

func TestParserErrorWithUnclosedInitSection(t *testing.T) {
	input := strings.Join([]string{
		`script "path/to/script.bf"`,
		`init {`,
		`    memory-size 1024`,
		`    stack-size  256`,
	}, "\n")

	expected := strings.Join([]string{
		"test.bft:4:20: error: unexpected EOF",
		"    4 |     stack-size  256<EOF>",
		"      |                    ^^^^^",
		"      |                    expect '}' to close",
	}, "\n")

	checkError(t, input, expected)
}

func TestParserErrorWithWrongFieldNameInInit(t *testing.T) {
	input := strings.Join([]string{
		`script "path/to/script.bf"`,
		`init {`,
		`    memory-size 1024`,
		`    stack-size  256`,
		`    42          uint8`,
		`}`,
		`case {}`,
	}, "\n")

	expected := strings.Join([]string{
		"test.bft:5:5: error: unexpected token type",
		"    5 |     42          uint8",
		"      |     ^^",
		"      |     expect identifier or '}', got INT",
	}, "\n")

	checkError(t, input, expected)
}

func TestParserErrorWithWrongFieldNameFormatInInit(t *testing.T) {
	input := strings.Join([]string{
		`script "path/to/script.bf"`,
		`init {`,
		`    memory-size  1024`,
		`    stack-size   256`,
		`    0xword       uint8`,
		`}`,
		`case {}`,
	}, "\n")

	expected := strings.Join([]string{
		"test.bft:5:5: error: invalid number format '0xword'",
		"    5 |     0xword       uint8",
		"      |     ^^^^^^",
		"      |     hexadecimal number should be 0x[0-9a-fA-F]+",
	}, "\n")

	checkError(t, input, expected)
}

func TestParserErrorWithWrongFieldValueFormatInInit(t *testing.T) {
	input := strings.Join([]string{
		`script "path/to/script.bf"`,
		`init {`,
		`    memory-size  1024`,
		`    stack-size   256w`,
		`}`,
		`case {}`,
	}, "\n")

	expected := strings.Join([]string{
		"test.bft:4:18: error: invalid number format '256w'",
		"    4 |     stack-size   256w",
		"      |                  ^^^^",
		"      |                  should be char [0-9] or underscore '_'",
	}, "\n")

	checkError(t, input, expected)
}

func TestParserOnSimpleCase(t *testing.T) {
	input := strings.Join([]string{
		`script "path/to/script.bf"`,
		`init {`,
		`    memory-size  1024`,
		`    stack-size   256`,
		`    word         uint8`,
		`}`,
		`case example {`,
		`    input {`,
		`        1 2 3 4 5 6`,
		`    }`,
		`    output {`,
		`        2 3 4 5 6 7`,
		`    }`,
		`    memory {`,
		`        3 4 5 6 7 8`,
		`    }`,
		`}`,
	}, "\n")

	expected := &TestDriverItem{
		ScriptName: NewContextItem("path/to/script.bf", nil),
		Init: InitParameters{
			MemorySize: NewContextItem(uint64(1024), nil),
			StackSize:  NewContextItem(uint64(256), nil),
			WordType:   NewContextItem(config.MemoryUnitTypeUint8, nil),
		},
		Tests: []TestCase{
			{
				Name:   NewContextItem("example", nil),
				Input:  ctxNums(1, 2, 3, 4, 5, 6),
				Output: ctxNums(2, 3, 4, 5, 6, 7),
				Memory: ctxNums(3, 4, 5, 6, 7, 8),
			},
		},
	}

	checkOK(t, input, expected)
}

func TestParserOnSimpleCaseWithMemoryAt(t *testing.T) {
	input := strings.Join([]string{
		`script "path/to/script.bf"`,
		`init {`,
		`    memory-size  1024`,
		`    stack-size   256`,
		`    word         uint8`,
		`}`,
		`case example {`,
		`    input {`,
		`        1 2 3 4 5 6`,
		`    }`,
		`    output {`,
		`        2 3 4 5 6 7`,
		`    }`,
		`    memory at 42 {`,
		`        3 4 5 6 7 8`,
		`    }`,
		`}`,
	}, "\n")

	expected := &TestDriverItem{
		ScriptName: NewContextItem("path/to/script.bf", nil),
		Init: InitParameters{
			MemorySize: NewContextItem(uint64(1024), nil),
			StackSize:  NewContextItem(uint64(256), nil),
			WordType:   NewContextItem(config.MemoryUnitTypeUint8, nil),
		},
		Tests: []TestCase{
			{
				Name:     NewContextItem("example", nil),
				Input:    ctxNums(1, 2, 3, 4, 5, 6),
				Output:   ctxNums(2, 3, 4, 5, 6, 7),
				Memory:   ctxNums(3, 4, 5, 6, 7, 8),
				MemoryAt: NewContextItem[uint64](42, nil),
			},
		},
	}

	checkOK(t, input, expected)
}
