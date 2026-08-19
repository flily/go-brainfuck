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

func ctxNums(nums ...int64) ContextItem[[]ContextItem[int64]] {
	var result []ContextItem[int64]
	for _, n := range nums {
		result = append(result, NewContextItem(n, nil))
	}

	return NewContextItem(result, nil)
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

func TestParseErrorWithNoBlock(t *testing.T) {
	input := strings.Join([]string{
		`script "hello.bf"`,
		`init memory-size 1024`,
		`case {}`,
		``,
	}, "\n")

	expected := strings.Join([]string{
		"test.bft:2:6: error: unexpected token type",
		"    2 | init memory-size 1024",
		"      |      ^^^^^^^^^^^",
		"      |      expect BRACE-LEFT here, got IDENTIFIER",
	}, "\n")

	checkError(t, input, expected)
}

func TestParserErrorOnWrongTokenType(t *testing.T) {
	input := strings.Join([]string{
		`script 4ever.bf`,
		`init {}`,
		`case {}`,
		``,
	}, "\n")

	expected := strings.Join([]string{
		"test.bft:1:8: error: invalid number format '4ever.bf'",
		"    1 | script 4ever.bf",
		"      |        ^^^^^^^^",
		"      |        should be char [0-9] or underscore '_'",
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
		Tests: []ContextItem[TestCase]{
			NewContextItem(
				TestCase{Name: NewContextItem("test case 1", nil)},
				nil),
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
		Init: NewContextItem(InitParameters{
			MemorySize: NewContextItem(uint64(1024), nil),
			StackSize:  NewContextItem(uint64(256), nil),
			WordType:   NewContextItem(config.MemoryUnitTypeUint8, nil),
			EOFValue:   NewContextItem(int64(-1), nil),
			IgnoreEOF:  NewContextItem(true, nil),
			RaiseEOF:   NewContextItem(false, nil),
		}, nil),
		Tests: []ContextItem[TestCase]{
			NewContextItem(TestCase{Name: NewContextItem("test case 1", nil)}, nil),
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
		Init: NewContextItem(InitParameters{
			MemorySize: NewContextItem(uint64(1024), nil),
			StackSize:  NewContextItem(uint64(256), nil),
			WordType:   NewContextItem(config.MemoryUnitTypeUint8, nil),
		}, nil),
		Tests: []ContextItem[TestCase]{
			NewContextItem(TestCase{
				Name:   NewContextItem("example", nil),
				Input:  ctxNums(1, 2, 3, 4, 5, 6),
				Output: ctxNums(2, 3, 4, 5, 6, 7),
				Memory: ctxNums(3, 4, 5, 6, 7, 8),
			}, nil),
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
		Init: NewContextItem(InitParameters{
			MemorySize: NewContextItem(uint64(1024), nil),
			StackSize:  NewContextItem(uint64(256), nil),
			WordType:   NewContextItem(config.MemoryUnitTypeUint8, nil),
		}, nil),
		Tests: []ContextItem[TestCase]{
			NewContextItem(TestCase{
				Name:     NewContextItem("example", nil),
				Input:    ctxNums(1, 2, 3, 4, 5, 6),
				Output:   ctxNums(2, 3, 4, 5, 6, 7),
				Memory:   ctxNums(3, 4, 5, 6, 7, 8),
				MemoryAt: NewContextItem[uint64](42, nil),
			}, nil),
		},
	}

	checkOK(t, input, expected)
}

func TestTestParserOnSimpleCaseWithInit(t *testing.T) {
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
		`    init {`,
		`        3 4 5 6 7 8`,
		`    }`,
		`    memory {`,
		`        4 5 6 7 8 9`,
		`    }`,
		`}`,
	}, "\n")

	expected := &TestDriverItem{
		ScriptName: NewContextItem("path/to/script.bf", nil),
		Init: NewContextItem(InitParameters{
			MemorySize: NewContextItem(uint64(1024), nil),
			StackSize:  NewContextItem(uint64(256), nil),
			WordType:   NewContextItem(config.MemoryUnitTypeUint8, nil),
		}, nil),
		Tests: []ContextItem[TestCase]{
			NewContextItem(TestCase{
				Name:   NewContextItem("example", nil),
				Input:  ctxNums(1, 2, 3, 4, 5, 6),
				Output: ctxNums(2, 3, 4, 5, 6, 7),
				Init:   ctxNums(3, 4, 5, 6, 7, 8),
				Memory: ctxNums(4, 5, 6, 7, 8, 9),
			}, nil),
		},
	}

	checkOK(t, input, expected)
}

func TestTestParserOnSimpleCaseWithInitAt(t *testing.T) {
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
		`    init at 42 {`,
		`        3 4 5 6 7 8`,
		`    }`,
		`    memory at 53 {`,
		`        4 5 6 7 8 9`,
		`    }`,
		`}`,
	}, "\n")

	expected := &TestDriverItem{
		ScriptName: NewContextItem("path/to/script.bf", nil),
		Init: NewContextItem(InitParameters{
			MemorySize: NewContextItem(uint64(1024), nil),
			StackSize:  NewContextItem(uint64(256), nil),
			WordType:   NewContextItem(config.MemoryUnitTypeUint8, nil),
		}, nil),
		Tests: []ContextItem[TestCase]{
			NewContextItem(TestCase{
				Name:     NewContextItem("example", nil),
				Input:    ctxNums(1, 2, 3, 4, 5, 6),
				Output:   ctxNums(2, 3, 4, 5, 6, 7),
				Init:     ctxNums(3, 4, 5, 6, 7, 8),
				InitAt:   NewContextItem[uint64](42, nil),
				Memory:   ctxNums(4, 5, 6, 7, 8, 9),
				MemoryAt: NewContextItem[uint64](53, nil),
			}, nil),
		},
	}

	checkOK(t, input, expected)
}

func TestParserErrorOnWrongNumberList(t *testing.T) {
	input := strings.Join([]string{
		`script "path/to/script.bf"`,
		`init {`,
		`    memory-size  1024`,
		`    stack-size   256`,
		`    word         uint8`,
		`}`,
		`case example {`,
		`    input {}`,
		`    output 3 4 5 6 7 8`,
		`    memory {}`,
		`}`,
	}, "\n")

	expected := strings.Join([]string{
		"test.bft:9:12: error: unexpected token type",
		"    9 |     output 3 4 5 6 7 8",
		"      |            ^",
		"      |            expect BRACE-LEFT here, got INT",
	}, "\n")

	checkError(t, input, expected)
}

func TestParserErrorOnWrongFormatInNumberListInput(t *testing.T) {
	input := strings.Join([]string{
		`script "path/to/script.bf"`,
		`init {`,
		`    memory-size  1024`,
		`    stack-size   256`,
		`    word         uint8`,
		`}`,
		`case example {`,
		`    input {`,
		`        3 4x 5 6 7 8`,
		`    }`,
		`    output {}`,
		`    memory {}`,
		`}`,
	}, "\n")

	expected := strings.Join([]string{
		"test.bft:9:11: error: invalid number format '4x'",
		"    9 |         3 4x 5 6 7 8",
		"      |           ^^",
		"      |           should be char [0-9] or underscore '_'",
	}, "\n")

	checkError(t, input, expected)
}

func TestParserErrorOnWrongFormatInNumberListOutput(t *testing.T) {
	input := strings.Join([]string{
		`script "path/to/script.bf"`,
		`init {`,
		`    memory-size  1024`,
		`    stack-size   256`,
		`    word         uint8`,
		`}`,
		`case example {`,
		`    input {}`,
		`    output {`,
		`        3 4x 5 6 7 8`,
		`    }`,
		`    memory {}`,
		`}`,
	}, "\n")

	expected := strings.Join([]string{
		"test.bft:10:11: error: invalid number format '4x'",
		"   10 |         3 4x 5 6 7 8",
		"      |           ^^",
		"      |           should be char [0-9] or underscore '_'",
	}, "\n")

	checkError(t, input, expected)
}

func TestParserErrorOnWrongItemTypeInNumberList(t *testing.T) {
	input := strings.Join([]string{
		`script "path/to/script.bf"`,
		`init {`,
		`    memory-size  1024`,
		`    stack-size   256`,
		`    word         uint8`,
		`}`,
		`case example {`,
		`    input {}`,
		`    output {`,
		`        3 four 5 6 7 8`,
		`    }`,
		`    memory {}`,
		`}`,
	}, "\n")

	expected := strings.Join([]string{
		"test.bft:10:11: error: unexpected token type",
		"   10 |         3 four 5 6 7 8",
		"      |           ^^^^",
		"      |           expect integer or '}', got IDENTIFIER",
	}, "\n")

	checkError(t, input, expected)
}

func TestParserErrorOnUnclosedNumberList(t *testing.T) {
	input := strings.Join([]string{
		`script "path/to/script.bf"`,
		`init {`,
		`    memory-size  1024`,
		`    stack-size   256`,
		`    word         uint8`,
		`}`,
		`case example {`,
		`    input {}`,
		`    output {`,
		`        3 4 5 6 7 8`,
	}, "\n")

	expected := strings.Join([]string{
		"test.bft:10:20: error: unexpected EOF",
		"   10 |         3 4 5 6 7 8<EOF>",
		"      |                    ^^^^^",
		"      |                    expect '}' to close",
	}, "\n")

	checkError(t, input, expected)
}

func TestParserErrorOnMemoryBlockNoBrace(t *testing.T) {
	input := strings.Join([]string{
		`script "path/to/script.bf"`,
		`init {`,
		`    memory-size  1024`,
		`    stack-size   256`,
		`    word         uint8`,
		`}`,
		`case example {`,
		`    input {}`,
		`    output {}`,
		`    memory 3 4 5 6 7 8`,
		`}`,
	}, "\n")

	expected := strings.Join([]string{
		"test.bft:10:12: error: syntax error",
		"   10 |     memory 3 4 5 6 7 8",
		"      |            ^",
		"      |            expect '{' to start memory block",
	}, "\n")

	checkError(t, input, expected)
}

func TestParserErrorOnInitBlockNoBrace(t *testing.T) {
	input := strings.Join([]string{
		`script "path/to/script.bf"`,
		`init {`,
		`    memory-size  1024`,
		`    stack-size   256`,
		`    word         uint8`,
		`}`,
		`case example {`,
		`    input {}`,
		`    output {}`,
		`    init 3 4 5 6 7 8`,
		`}`,
	}, "\n")

	expected := strings.Join([]string{
		"test.bft:10:10: error: syntax error",
		"   10 |     init 3 4 5 6 7 8",
		"      |          ^",
		"      |          expect '{' to start init block",
	}, "\n")

	checkError(t, input, expected)
}

func TestParserErrorOnMemoryBlockWrongStartTokenFormat1(t *testing.T) {
	input := strings.Join([]string{
		`script "path/to/script.bf"`,
		`init {`,
		`    memory-size  1024`,
		`    stack-size   256`,
		`    word         uint8`,
		`}`,
		`case example {`,
		`    input {}`,
		`    output {}`,
		`    memory 0xqwer {}`,
		`}`,
	}, "\n")

	expected := strings.Join([]string{
		"test.bft:10:12: error: invalid number format '0xqwer'",
		"   10 |     memory 0xqwer {}",
		"      |            ^^^^^^",
		"      |            hexadecimal number should be 0x[0-9a-fA-F]+",
	}, "\n")

	checkError(t, input, expected)
}

func TestParserErrorOnMemoryBlockWrongStartTokenFormat2(t *testing.T) {
	input := strings.Join([]string{
		`script "path/to/script.bf"`,
		`init {`,
		`    memory-size  1024`,
		`    stack-size   256`,
		`    word         uint8`,
		`}`,
		`case example {`,
		`    input {}`,
		`    output {}`,
		`    memory at 42 0xqwer {}`,
		`}`,
	}, "\n")

	expected := strings.Join([]string{
		"test.bft:10:18: error: invalid number format '0xqwer'",
		"   10 |     memory at 42 0xqwer {}",
		"      |                  ^^^^^^",
		"      |                  hexadecimal number should be 0x[0-9a-fA-F]+",
	}, "\n")

	checkError(t, input, expected)
}

func TestParserErrorOnMemoryBlockWrongAtCommand(t *testing.T) {
	input := strings.Join([]string{
		`script "path/to/script.bf"`,
		`init {`,
		`    memory-size  1024`,
		`    stack-size   256`,
		`    word         uint8`,
		`}`,
		`case example {`,
		`    input {}`,
		`    output {}`,
		`    memory on 42 {}`,
		`}`,
	}, "\n")

	expected := strings.Join([]string{
		"test.bft:10:12: error: syntax error",
		"   10 |     memory on 42 {}",
		"      |            ^^",
		"      |            use 'at' to specify start address of memory block",
	}, "\n")

	checkError(t, input, expected)
}

func TestParserErrorOnMemoryBlockWrongAtCommandWrongAddressFormat(t *testing.T) {
	input := strings.Join([]string{
		`script "path/to/script.bf"`,
		`init {`,
		`    memory-size  1024`,
		`    stack-size   256`,
		`    word         uint8`,
		`}`,
		`case example {`,
		`    input {}`,
		`    output {}`,
		`    memory at 0xqwer {}`,
		`}`,
	}, "\n")

	expected := strings.Join([]string{
		"test.bft:10:15: error: invalid number format '0xqwer'",
		"   10 |     memory at 0xqwer {}",
		"      |               ^^^^^^",
		"      |               hexadecimal number should be 0x[0-9a-fA-F]+",
	}, "\n")

	checkError(t, input, expected)
}

func TestParserErrorOnMemoryBlockWrongAtCommandWrongAddressType(t *testing.T) {
	input := strings.Join([]string{
		`script "path/to/script.bf"`,
		`init {`,
		`    memory-size  1024`,
		`    stack-size   256`,
		`    word         uint8`,
		`}`,
		`case example {`,
		`    input {}`,
		`    output {}`,
		`    memory at zero {}`,
		`}`,
	}, "\n")

	expected := strings.Join([]string{
		"test.bft:10:15: error: syntax error",
		"   10 |     memory at zero {}",
		"      |               ^^^^",
		"      |               expect integer as address",
	}, "\n")

	checkError(t, input, expected)
}

func TestParserErrorOnMemoryBlockWrongAtCommandNegativeAddress(t *testing.T) {
	input := strings.Join([]string{
		`script "path/to/script.bf"`,
		`init {`,
		`    memory-size  1024`,
		`    stack-size   256`,
		`    word         uint8`,
		`}`,
		`case example {`,
		`    input {}`,
		`    output {}`,
		`    memory at -42 {}`,
		`}`,
	}, "\n")

	expected := strings.Join([]string{
		"test.bft:10:15: error: invalid address",
		"   10 |     memory at -42 {}",
		"      |               ^^^",
		"      |               address MUST BE positive",
	}, "\n")

	checkError(t, input, expected)
}

func TestParserErrorOnMemoryBlockWithWrongNumberList(t *testing.T) {
	input := strings.Join([]string{
		`script "path/to/script.bf"`,
		`init {`,
		`    memory-size  1024`,
		`    stack-size   256`,
		`    word         uint8`,
		`}`,
		`case example {`,
		`    input {}`,
		`    output {}`,
		`    memory at 42 {`,
		`        3 four 5 6 7 8`,
		`    }`,
		`}`,
	}, "\n")

	expected := strings.Join([]string{
		"test.bft:11:11: error: unexpected token type",
		"   11 |         3 four 5 6 7 8",
		"      |           ^^^^",
		"      |           expect integer or '}', got IDENTIFIER",
	}, "\n")

	checkError(t, input, expected)
}

func TestParserErrorOnMemoryBlockWithUnknownField(t *testing.T) {
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
		`    lorem ipsum`,
		`}`,
	}, "\n")

	expected := strings.Join([]string{
		"test.bft:17:5: error: unknown field name",
		"   17 |     lorem ipsum",
		"      |     ^^^^^",
		"      |     use one of 'input', 'output', 'memory', got 'lorem'",
	}, "\n")

	checkError(t, input, expected)
}

func TestParserErrorOnCaseName(t *testing.T) {
	input := strings.Join([]string{
		`script "path/to/script.bf"`,
		`init {`,
		`    memory-size  1024`,
		`    stack-size   256`,
		`    word         uint8`,
		`}`,
		`case 0xqwer {}`,
	}, "\n")

	expected := strings.Join([]string{
		"test.bft:7:6: error: invalid number format '0xqwer'",
		"    7 | case 0xqwer {}",
		"      |      ^^^^^^",
		"      |      hexadecimal number should be 0x[0-9a-fA-F]+",
	}, "\n")

	checkError(t, input, expected)
}

func TestParserErrorOnCaseWithoutWrongFormat(t *testing.T) {
	input := strings.Join([]string{
		`script "path/to/script.bf"`,
		`init {`,
		`    memory-size  1024`,
		`    stack-size   256`,
		`    word         uint8`,
		`}`,
		`case example 0xqwer`,
	}, "\n")

	expected := strings.Join([]string{
		"test.bft:7:14: error: invalid number format '0xqwer'",
		"    7 | case example 0xqwer",
		"      |              ^^^^^^",
		"      |              hexadecimal number should be 0x[0-9a-fA-F]+",
	}, "\n")

	checkError(t, input, expected)
}

func TestParserErrorOnCaseWithoutBlock(t *testing.T) {
	input := strings.Join([]string{
		`script "path/to/script.bf"`,
		`init {`,
		`    memory-size  1024`,
		`    stack-size   256`,
		`    word         uint8`,
		`}`,
		`case example input { 1 2 3 4 5 6 }`,
	}, "\n")

	expected := strings.Join([]string{
		"test.bft:7:14: error: unexpected token found",
		"    7 | case example input { 1 2 3 4 5 6 }",
		"      |              ^^^^^",
		"      |              expect '{' here, got IDENTIFIER",
	}, "\n")

	checkError(t, input, expected)
}

func TestParserErrorOnCaseWithWrongFieldNameFormat(t *testing.T) {
	input := strings.Join([]string{
		`script "path/to/script.bf"`,
		`init {`,
		`    memory-size  1024`,
		`    stack-size   256`,
		`    word         uint8`,
		`}`,
		`case example {`,
		`    0xqwer {}`,
		`    input {}`,
		`    output {}`,
		`    memory {}`,
		`}`,
	}, "\n")

	expected := strings.Join([]string{
		"test.bft:8:5: error: invalid number format '0xqwer'",
		"    8 |     0xqwer {}",
		"      |     ^^^^^^",
		"      |     hexadecimal number should be 0x[0-9a-fA-F]+",
	}, "\n")

	checkError(t, input, expected)
}

func TestParserErrorOnCaseWithUnclosedBlock(t *testing.T) {
	input := strings.Join([]string{
		`script "path/to/script.bf"`,
		`init {`,
		`    memory-size  1024`,
		`    stack-size   256`,
		`    word         uint8`,
		`}`,
		`case example {`,
		`    input {}`,
		`    output {}`,
	}, "\n")

	expected := strings.Join([]string{
		"test.bft:9:14: error: unexpected EOF",
		"    9 |     output {}<EOF>",
		"      |              ^^^^^",
		"      |              expect '}' to close",
	}, "\n")

	checkError(t, input, expected)
}

func TestParserErrorOnCaseWithUnexpectedFieldName(t *testing.T) {
	input := strings.Join([]string{
		`script "path/to/script.bf"`,
		`init {`,
		`    memory-size  1024`,
		`    stack-size   256`,
		`    word         uint8`,
		`}`,
		`case example {`,
		`    input {}`,
		`    output {}`,
		`    42 {}`,
		`}`,
	}, "\n")

	expected := strings.Join([]string{
		"test.bft:10:5: error: unexpected token type",
		"   10 |     42 {}",
		"      |     ^^",
		"      |     expect identifier or '}', got INT",
	}, "\n")

	checkError(t, input, expected)
}

func TestParserErrorWithInvalidKeyword1(t *testing.T) {
	input := strings.Join([]string{
		`0xqwer`,
	}, "\n")

	expected := strings.Join([]string{
		"test.bft:1:1: error: invalid number format '0xqwer'",
		"    1 | 0xqwer",
		"      | ^^^^^^",
		"      | hexadecimal number should be 0x[0-9a-fA-F]+",
	}, "\n")

	checkError(t, input, expected)
}

func TestParserErrorWithInvalidKeyword2(t *testing.T) {
	input := strings.Join([]string{
		`42`,
	}, "\n")

	expected := strings.Join([]string{
		"test.bft:1:1: error: invalid identifier",
		"    1 | 42",
		"      | ^^",
		"      | expect identifier here, got INT",
	}, "\n")

	checkError(t, input, expected)
}
