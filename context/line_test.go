package context

import (
	"testing"

	"strings"
)

func TestLineContextHighlightText(t *testing.T) {
	//    0         1         2         3         4
	//    0    5    0    5    0    5    0    5    0
	s := "the quick brown fox jumps over the lazy dog"
	line := NewLineFromBytes(42, []byte(s), EolLF)
	lctx := line.markLine(4, 9)

	got := lctx.HighlighText("lorem ipsum")
	expected := strings.Join([]string{
		"   43 | the quick brown fox jumps over the lazy dog",
		"      |     ^^^^^",
		"      |     lorem ipsum",
	}, "\n")
	if got != expected {
		t.Errorf("expected:\n%s\ngot:\n%s", expected, got)
	}

	if line, column := lctx.Position(); line != 42 || column != 4 {
		t.Errorf("expected position to be (42, 4), got (%d, %d)", line, column)
	}

	expLast := 9
	if lctx.Last() != expLast {
		t.Errorf("expected last highlight to be %d, got %d", expLast, lctx.Last())
	}
}

func TestLineContextHighlightTextWithoutMessage(t *testing.T) {
	//    0         1         2         3         4
	//    0    5    0    5    0    5    0    5    0
	s := "the quick brown fox jumps over the lazy dog"
	line := NewLineFromBytes(42, []byte(s), EolLF)
	lctx := line.markLine(4, 9)

	got := lctx.HighlighText(NoHighlightMessage)
	expected := strings.Join([]string{
		"   43 | the quick brown fox jumps over the lazy dog",
		"      |     ^^^^^",
	}, "\n")
	if got != expected {
		t.Errorf("expected:\n%s\ngot:\n%s", expected, got)
	}
}

func TestLineContextHighlightTextMultiParts(t *testing.T) {
	//    0         1         2         3         4
	//    0    5    0    5    0    5    0    5    0
	s := "the quick brown fox jumps over the lazy dog"
	line := NewLineFromBytes(42, []byte(s), EolLF)
	lctx := line.markLine(4, 9)
	lctx.MarkLine(16, 19)

	got := lctx.HighlighText("lorem ipsum")
	expected := strings.Join([]string{
		"   43 | the quick brown fox jumps over the lazy dog",
		"      |     ^^^^^       ^^^",
		"      |     lorem ipsum",
	}, "\n")
	if got != expected {
		t.Errorf("expected:\n%s\ngot:\n%s", expected, got)
	}
}

func TestLineContextHighlightTextWithTabCharacters(t *testing.T) {
	//    0         1         2         3         4
	//    0    5    0    5    0    5    0    5    0
	s := "the quick\tbrown fox jumps over the lazy dog"
	line := NewLineFromBytes(42, []byte(s), EolLF)
	lctx := line.markLine(10, 15)

	got := lctx.HighlighText("lorem ipsum")
	expected := strings.Join([]string{
		"   43 | the quick\tbrown fox jumps over the lazy dog",
		//  43 | the quick       brown fox jumps over the lazy dog
		"      |                 ^^^^^",
		"      |                 lorem ipsum",
	}, "\n")
	if got != expected {
		t.Errorf("expected:\n%s\ngot:\n%s", expected, got)
	}
}

func TestLineContextHighlightTextWithTabCharactersMultiParts(t *testing.T) {
	//    0         1         2         3         4
	//    0    5    0    5    0    5    0    5    0
	s := "the quick\tbrown fox jumps over the lazy dog"
	line := NewLineFromBytes(42, []byte(s), EolLF)
	lctx := line.markLine(10, 15)
	lctx.MarkLine(4, 9)

	got := lctx.HighlighText("lorem ipsum")
	expected := strings.Join([]string{
		"   43 | the quick\tbrown fox jumps over the lazy dog",
		//  43 | the quick       brown fox jumps over the lazy dog
		"      |     ^^^^^       ^^^^^",
		"      |     lorem ipsum",
	}, "\n")
	if got != expected {
		t.Errorf("expected:\n%s\ngot:\n%s", expected, got)
	}
}

func TestLineContextHighlightTextWithMixedLanguages1(t *testing.T) {
	//    0         1         2         3         4         5         6
	//    0    5    0    5    0    5    0    5    0    5    0    5    0
	s := "the quick brown fox 我能吞下玻璃而不伤身体 jumps over the lazy dog"
	line := NewLineFromBytes(42, []byte(s), EolLF)
	lctx := line.markLine(20, 24)

	got := lctx.HighlighText("lorem ipsum")
	expected := strings.Join([]string{
		"   43 | the quick brown fox 我能吞下玻璃而不伤身体 jumps over the lazy dog",
		"      |                     ^^^^^^^^",
		"      |                     lorem ipsum",
	}, "\n")
	if got != expected {
		t.Errorf("expected:\n%s\ngot:\n%s", expected, got)
	}
}

func TestLineContextHighlightTextWithMixedLanguages2(t *testing.T) {
	//    0         1         2         3         4         5         6
	//    0    5    0    5    0    5    0    5    0    5    0    5    0
	s := "the quick brown fox 我能吞下玻璃而不伤身体 jumps over the lazy dog"
	line := NewLineFromBytes(42, []byte(s), EolLF)
	lctx := line.markLine(32, 37)

	got := lctx.HighlighText("lorem ipsum")
	expected := strings.Join([]string{
		"   43 | the quick brown fox 我能吞下玻璃而不伤身体 jumps over the lazy dog",
		"      |                                            ^^^^^",
		"      |                                            lorem ipsum",
	}, "\n")
	if got != expected {
		t.Errorf("expected:\n%s\ngot:\n%s", expected, got)
	}
}
