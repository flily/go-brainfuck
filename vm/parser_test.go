package vm

import (
	"testing"

	"slices"
	"strings"

	"github.com/flily/go-brainfuck/context"
)

const (
	testFilename = "test.bf"
)

func parseTestCodeOk(t *testing.T, code string, expCodes []Instruction, expNext []int) {
	t.Helper()

	file := context.ReadFileString(testFilename, code)
	parser := NewParser(file)

	codemap, err := parser.Parse()
	if err != nil {
		t.Fatalf("error parsing code:\n%s", err)
	}

	expected := buildCodeMap(expCodes, expNext)
	if !codemap.CodeEquals(expected) {
		t.Errorf("code expected: %v", expected.Codes)
		t.Errorf("code got:      %v", codemap.Codes)
		t.Errorf("next expected: %v", expected.Next)
		t.Errorf("next got:      %v", codemap.Next)
		t.Fatalf("got wrong code map")
	}
}

func parseTestCodeFailure(t *testing.T, code string, errMessage string) {
	t.Helper()

	file := context.ReadFileString(testFilename, code)
	parser := NewParser(file, NewStandardInstructionSet())
	codemap, err := parser.Parse()
	if codemap != nil {
		t.Fatalf("expected parse error, but got code map: %v", codemap)
	}

	if err == nil {
		t.Fatalf("expected parse error, but got nil")
	}

	if merr := err.Error(); merr != errMessage {
		t.Fatalf("error message mismatch, expected:\n%s\ngot:\n%s", errMessage, merr)
	}
}

func buildCodeMap(codes []Instruction, nexts []int) *CodeMap {
	codemap := NewCodeMap()

	codemap.Codes = InstructionsToCodes(codes)
	codemap.Next = slices.Clone(nexts)

	return codemap
}

func TestParseSimpleCode(t *testing.T) {
	code := "+++"

	codes := []Instruction{
		InstructionAdd,
		InstructionAdd,
		InstructionAdd,
	}
	nexts := []int{
		-1,
		-1,
		-1,
	}

	parseTestCodeOk(t, code, codes, nexts)
}

func TestParseSimpleLoopCode(t *testing.T) {
	code := "[-]"

	codes := []Instruction{
		InstructionLoopBegin,
		InstructionSub,
		InstructionLoopEnd,
	}
	nexts := []int{
		2,
		-1,
		0,
	}

	parseTestCodeOk(t, code, codes, nexts)
}

func TestParseWithNonInstructionCharacters(t *testing.T) {
	code := "++a--b"

	codes := []Instruction{
		InstructionAdd,
		InstructionAdd,
		InstructionSub,
		InstructionSub,
	}
	nexts := []int{
		-1,
		-1,
		-1,
		-1,
	}

	parseTestCodeOk(t, code, codes, nexts)
}

func TestParseErrorNoMatchedEndLoop(t *testing.T) {
	code := "[++--"

	expected := strings.Join([]string{
		"test.bf:1:1: error: unclosed loop bracket",
		"    1 | [++--",
		"      | ^",
		"      | no matched ']' for this",
	}, "\n")
	parseTestCodeFailure(t, code, expected)
}

func TestParseErrorNoMatchedBeginLoop(t *testing.T) {
	code := "++--]"

	expected := strings.Join([]string{
		"test.bf:1:5: error: unexpected closing loop bracket",
		"    1 | ++--]",
		"      |     ^",
		"      |     no matched '[' for this",
	}, "\n")
	parseTestCodeFailure(t, code, expected)
}
