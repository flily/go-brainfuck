package driver

import (
	"testing"

	"strings"

	"github.com/flily/go-brainfuck/context"
)

const (
	TestDriverFilename = "test.bftest"
)

func loadDriver(t *testing.T, c string) *TestDriverItem {
	fctx := context.ReadFileString(TestDriverFilename, c)
	parser := NewParser(TestDriverFilename, fctx)
	driverItem, err := parser.Parse()
	if err != nil {
		t.Fatalf("failed to parse driver\n%s", err)
	}

	return driverItem
}

type testDriverCase struct {
	driver   []string
	code     []string
	expected []string
}

func (c testDriverCase) Ok(t *testing.T) {
	driverText := strings.Join(c.driver, "\n")
	driver := loadDriver(t, driverText)

	scriptName := driver.ScriptName.Value
	scriptText := strings.Join(c.code, "\n")
	scriptCtx := context.ReadFileString(scriptName, scriptText)

	err := genericRunWith(driver, scriptCtx)
	if err != nil {
		t.Fatalf("failed to run driver\n%s", err)
	}
}

func (c testDriverCase) Error(t *testing.T) {
	driverText := strings.Join(c.driver, "\n")
	driver := loadDriver(t, driverText)

	scriptName := driver.ScriptName.Value
	scriptText := strings.Join(c.code, "\n")
	scriptCtx := context.ReadFileString(scriptName, scriptText)

	err := genericRunWith(driver, scriptCtx)
	if err == nil {
		t.Fatalf("expected error but got none")
	}
	exp := strings.Join(c.expected, "\n")
	if exp != err.Error() {
		t.Fatalf("expected error:\n%s\nbut got:\n%s", exp, err.Error())
	}
}

func TestSimpleCase(t *testing.T) {
	testDriverCase{
		driver: []string{
			`script "test.bf"`,
			"init {",
			"    memory-size   1024",
			"    stack-size    128",
			"    word          uint8",
			"}",
			"",
			"case {",
			"    memory {",
			"        3 0 0 0",
			"    }",
			"}",
		},
		code: []string{
			"+++++--",
		},
	}.Ok(t)
}

func TestCaseSimpleIO(t *testing.T) {
	testDriverCase{
		driver: []string{
			`script "test.bf"`,
			"init {",
			"    memory-size   1024",
			"    stack-size    128",
			"    word          uint8",
			"}",
			"",
			"case {",
			"    input {",
			"        42",
			"    }",
			"    output {",
			"        45",
			"    }",
			"}",
		},
		code: []string{
			",+++.",
		},
	}.Ok(t)
}

func TestCaseErrorWithCompilationError(t *testing.T) {
	testDriverCase{
		driver: []string{
			`script "test.bf"`,
			"init {",
			"    memory-size   1024",
			"    stack-size    128",
			"    word          uint8",
			"}",
			"",
			"case {",
			"    input {",
			"        42",
			"    }",
			"    output {",
			"        45",
			"    }",
			"}",
		},
		code: []string{
			",+++[.",
		},
		expected: []string{
			"test.bf:1:5: error: unclosed loop bracket",
			"    1 | ,+++[.",
			"      |     ^",
			"      |     no matched ']' for this",
		},
	}.Error(t)
}

func TestCaseErrorWithInputNotRead(t *testing.T) {
	testDriverCase{
		driver: []string{
			`script "test.bf"`,
			"init {",
			"    memory-size   1024",
			"    stack-size    128",
			"    word          uint8",
			"}",
			"",
			"case {",
			"    input {",
			"        42 43",
			"    }",
			"    output {",
			"        45",
			"    }",
			"}",
		},
		code: []string{
			",+++.",
		},
		expected: []string{
			"test.bftest:10:12: error: input data not read by program",
			"   10 |         42 43",
			"      |            ^^",
			"      |            data not read by program",
		},
	}.Error(t)
}

func TestCaseErrorWithOutputLessThanExpected(t *testing.T) {
	testDriverCase{
		driver: []string{
			`script "test.bf"`,
			"init {",
			"    memory-size   1024",
			"    stack-size    128",
			"    word          uint8",
			"}",
			"",
			"case {",
			"    output {",
			"        3 4 5 6",
			"    }",
			"}",
		},
		code: []string{
			"+++.+.",
		},
		expected: []string{
			"test.bftest:10:13: error: output data less than expected",
			"   10 |         3 4 5 6",
			"      |             ^",
			"      |             output 2 data items, expect 4",
		},
	}.Error(t)
}

func TestCaseErrorWithOutputMoreThanExpected(t *testing.T) {
	testDriverCase{
		driver: []string{
			`script "test.bf"`,
			"init {",
			"    memory-size   1024",
			"    stack-size    128",
			"    word          uint8",
			"}",
			"",
			"case {",
			"    output {",
			"        3 4 5 6",
			"    }",
			"}",
		},
		code: []string{
			"+++.+.+.+.+.",
		},
		expected: []string{
			"test.bftest:10:15: error: output data more than expected",
			"   10 |         3 4 5 6",
			"      |               ^",
			"      |               output 5 data items, expect 4",
		},
	}.Error(t)
}

func TestCaseErrorWithOutputValueMismatch(t *testing.T) {
	testDriverCase{
		driver: []string{
			`script "test.bf"`,
			"init {",
			"    memory-size   1024",
			"    stack-size    128",
			"    word          uint8",
			"}",
			"",
			"case {",
			"    output {",
			"        3 4 9 6",
			"    }",
			"}",
		},
		code: []string{
			"+++.+.+.+.",
		},
		expected: []string{
			"test.bftest:10:13: error: output value mismatch",
			"   10 |         3 4 9 6",
			"      |             ^",
			"      |             got 5",
		},
	}.Error(t)
}

func TestCaseErrorWithMemoryBaseOutOfRange(t *testing.T) {
	testDriverCase{
		driver: []string{
			`script "test.bf"`,
			"init {",
			"    memory-size   2",
			"    stack-size    128",
			"    word          uint8",
			"}",
			"",
			"case {",
			"    memory at 3 {",
			"        1 2 3 4",
			"    }",
			"}",
		},
		code: []string{
			"+++++--",
		},
		expected: []string{
			"test.bftest:9:15: error: memory base out of range",
			"    9 |     memory at 3 {",
			"      |               ^",
			"      |               vm memory size is set to 2",
		},
	}.Error(t)
}

func TestCaseErrorWithMemoryTooSmall(t *testing.T) {
	testDriverCase{
		driver: []string{
			`script "test.bf"`,
			"init {",
			"    memory-size   2",
			"    stack-size    128",
			"    word          uint8",
			"}",
			"",
			"case {",
			"    memory {",
			"        3 0 0 0",
			"    }",
			"}",
		},
		code: []string{
			"+++++--",
		},
		expected: []string{
			"test.bftest:10:13: error: memory assertion out of range",
			"   10 |         3 0 0 0",
			"      |             ^",
			"      |             vm memory size is set to 2",
		},
	}.Error(t)
}

func TestCaseErrorWithMemoryValueMismatch(t *testing.T) {
	testDriverCase{
		driver: []string{
			`script "test.bf"`,
			"init {",
			"    memory-size   1024",
			"    stack-size    128",
			"    word          uint8",
			"}",
			"",
			"case {",
			"    memory at 2 {",
			"        3 4 7 6",
			"    }",
			"}",
		},
		code: []string{
			">>",
			"+++>",
			"++++>",
			"+++++>",
			"++++++",
		},
		expected: []string{
			"test.bftest:10:13: error: memory value mismatch",
			"   10 |         3 4 7 6",
			"      |             ^",
			"      |             got 5",
		},
	}.Error(t)
}
