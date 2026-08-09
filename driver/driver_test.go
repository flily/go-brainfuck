package driver

import (
	"strings"
	"testing"

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

func TestCaseWithMemoryTooSmall(t *testing.T) {
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
