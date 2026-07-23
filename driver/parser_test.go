package driver

import (
	"testing"

	"strings"
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

func TestParserEmptyContent(t *testing.T) {
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

func TestParserStartWithInitSection(t *testing.T) {
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

func TestParserStartWithCaseSection(t *testing.T) {
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
