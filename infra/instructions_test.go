package infra

import (
	"testing"

	"slices"
)

type runeInst rune

func (r runeInst) Char() rune {
	return rune(r)
}

func (r runeInst) String() string {
	return string(r)
}

func TestConvertInstructionsFrom(t *testing.T) {
	input := []runeInst{
		'+',
		'-',
	}

	expected := []Instruction{
		runeInst('+'),
		runeInst('-'),
	}

	got := ConvertInstructionsFrom(input)
	if !slices.Equal(got, expected) {
		t.Fatalf("expected %v, got %v", expected, got)
	}
}
