package context

import (
	"testing"

	"fmt"
	"strings"
)

func createTestCursor1() *Cursor {
	parts := []string{
		"lorem\n",
		"ipsum\n",
	}

	return NewCursorFromString("example.txt", strings.Join(parts, ""))
}

func createTestCursor2() *Cursor {
	parts := []string{
		"lorem ipsum dolor sit amet\n",
		"consectetur adipiscing elit\n",
		"\n",
		"sed do eiusmod tempor incididunt\n",
		"ut labore et dolore magna aliqua\n",
		"ut enim ad minim veniam\r\n",
		"\r\n",
		"quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat\n",
		"duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur\n",
		"excepteur sint occaecat cupid atat non proident\n",
	}

	return NewCursorFromString("example.txt", strings.Join(parts, ""))
}

func TestCursorRuneBasic(t *testing.T) {
	cursor := createTestCursor1()

	firstLine := []rune{
		'l', 'o', 'r', 'e', 'm',
	}

	for i, expChar := range firstLine {
		r, eol, eof := cursor.Rune()
		if eol || eof {
			t.Fatalf("unexpected EOL/EOF at %s  (EOL=%v, EOF=%v)", cursor.PositionString(), eol, eof)
		}

		if r != expChar {
			t.Errorf("expected rune '%c' at %s, got '%c'", expChar, cursor.PositionString(), r)
		}

		expPosition := fmt.Sprintf("example.txt:1:%d", i+1)
		if cursor.PositionString() != expPosition {
			t.Errorf("expected position '%s', got '%s'", expPosition, cursor.PositionString())
		}

		_, eol = cursor.NextInLine()
		if i != len(firstLine)-1 && eol {
			t.Errorf("unexpected EOL at %s", cursor.PositionString())
		}
	}

	for range 5 {
		_, eol := cursor.NextInLine()
		if !eol {
			t.Errorf("expected EOL at %s", cursor.PositionString())
		}
	}
}

func TestCursorCurrentChar(t *testing.T) {
	cursor := createTestCursor1()

	r, c1 := cursor.CurrentChar()
	exp1 := strings.Join([]string{
		"    1 | lorem",
		"      | ^",
		"      | here",
	}, "\n")
	got1 := c1.HighlightText("here")
	if got1 != exp1 {
		t.Errorf("expected:\n%s\ngot:\n%s", exp1, got1)
	}

	if r != 'l' {
		t.Errorf("expected rune 'l' at %s, got '%c'", cursor.PositionString(), r)
	}

	_, _, _ = cursor.Next()
	r, c2 := cursor.CurrentChar()
	exp2 := strings.Join([]string{
		"    1 | lorem",
		"      |  ^",
		"      |  here",
	}, "\n")
	got2 := c2.HighlightText("here")
	if got2 != exp2 {
		t.Errorf("expected:\n%s\ngot:\n%s", exp2, got2)
	}

	if r != 'o' {
		t.Errorf("expected rune 'o' at %s, got '%c'", cursor.PositionString(), r)
	}

	_ = cursor.NextNonEmptyLine()
	r, c3 := cursor.CurrentChar()
	exp3 := strings.Join([]string{
		"    2 | ipsum",
		"      | ^",
		"      | here",
	}, "\n")
	got3 := c3.HighlightText("here")
	if got3 != exp3 {
		t.Errorf("expected:\n%s\ngot:\n%s", exp3, got3)
	}

	if r != 'i' {
		t.Errorf("expected rune 'i' at %s, got '%c'", cursor.PositionString(), r)
	}
}

func TestCursorCurrentCharAtEmptyLine(t *testing.T) {
	text := "\n"

	cursor := NewCursorFromString("example.txt", text)

	r, ctx := cursor.CurrentChar()
	exp := strings.Join([]string{
		"    1 | <EOL LF>",
		"      | ^^^^^^^^",
		"      | here",
	}, "\n")
	got := ctx.HighlightText("here")
	if got != exp {
		t.Errorf("expected:\n%s\ngot:\n%s", exp, got)
	}

	if r != 0 {
		t.Errorf("expected rune 0 at %s, got '%c'", cursor.PositionString(), r)
	}
}

func TestCursorPeek(t *testing.T) {
	cursor := createTestCursor1()

	firstLine := []rune{
		'l', 'o', 'r', 'e', 'm',
	}

	r0, eol, eof := cursor.Rune()
	if eol || eof {
		t.Fatalf("unexpected EOL/EOF at %s  (EOL=%v, EOF=%v)", cursor.PositionString(), eol, eof)
	}

	if r0 != firstLine[0] {
		t.Errorf("expected rune '%c' at %s, got '%c'", firstLine[0], cursor.PositionString(), r0)
	}

	r1, eol, eof := cursor.Peek(1)
	if eol || eof {
		t.Fatalf("unexpected EOL/EOF at %s  (EOL=%v, EOF=%v)", cursor.PositionString(), eol, eof)
	}

	if r1 != firstLine[1] {
		t.Errorf("expected rune '%c' at %s, got '%c'", firstLine[1], cursor.PositionString(), r1)
	}

	r2, eol, eof := cursor.Peek(2)
	if eol || eof {
		t.Fatalf("unexpected EOL/EOF at %s  (EOL=%v, EOF=%v)", cursor.PositionString(), eol, eof)
	}

	if r2 != firstLine[2] {
		t.Errorf("expected rune '%c' at %s, got '%c'", firstLine[2], cursor.PositionString(), r2)
	}

	// move to end of line
	for range 4 {
		cursor.NextInLine()
	}

	re, eol, eof := cursor.Rune()
	if eol || eof {
		t.Fatalf("unexpected EOL/EOF at %s  (EOL=%v, EOF=%v)", cursor.PositionString(), eol, eof)
	}

	if re != firstLine[4] {
		t.Errorf("expected rune '%c' at %s, got '%c'", firstLine[4], cursor.PositionString(), re)
	}

	rp, eol, eof := cursor.Peek(1)
	if !eol || eof {
		t.Errorf("expected to be at end of line at %s, got rune '%c'", cursor.PositionString(), rp)
	}

	if rp != 0 {
		t.Errorf("expected rune to be 0 at end of line, got '%c'", rp)
	}
}

func TestCursorNextInLineBasic(t *testing.T) {
	cursor := createTestCursor1()

	firstLine := []rune{
		'i', 'p', 's', 'u', 'm',
	}

	eof := cursor.NextNonEmptyLine()
	if eof {
		t.Fatalf("unexpected EOF at %s", cursor.PositionString())
	}

	for i, expChar := range firstLine {
		r, eol, eof := cursor.Rune()
		if eol || eof {
			t.Fatalf("unexpected EOL/EOF at %s  (EOL=%v, EOF=%v)", cursor.PositionString(), eol, eof)
		}

		if r != expChar {
			t.Errorf("expected rune '%c' at %s, got '%c'", expChar, cursor.PositionString(), r)
		}

		_, eol = cursor.NextInLine()
		if i != len(firstLine)-1 && eol {
			t.Errorf("unexpected EOL at %s", cursor.PositionString())
		}
	}
}

func TestCursorNextInLineWithEmptyLine(t *testing.T) {
	cursor := createTestCursor2()

	content := []rune{
		's', 'e', 'd',
	}

	for range 2 {
		eof := cursor.NextNonEmptyLine()
		if eof {
			t.Fatalf("unexpected EOF at %s", cursor.PositionString())
		}
	}

	for _, expChar := range content {
		r, eol, eof := cursor.Rune()
		if eol || eof {
			t.Fatalf("unexpected EOL/EOF at %s  (EOL=%v, EOF=%v)", cursor.PositionString(), eol, eof)
		}

		if r != expChar {
			t.Errorf("expect rune '%c' at %s, got '%c'", expChar, cursor.PositionString(), r)
		}

		_, eol = cursor.NextInLine()
		if eol {
			t.Errorf("unexpected EOL at %s", cursor.PositionString())
		}
	}
}

func TestCursorNextInLineInTheLastLine(t *testing.T) {
	cursor := createTestCursor1()

	cursor.NextNonEmptyLine() // move to second line
	for range 4 {
		_, eof := cursor.NextInLine()
		if eof {
			t.Fatalf("unexpected EOF at %s", cursor.PositionString())
		}
	}

	// 5th char
	_, eof := cursor.NextInLine()
	if !eof {
		t.Errorf("expect EOL at %s", cursor.PositionString())
	}

	r, eof := cursor.NextInLine()
	if !eof {
		t.Errorf("expect EOL at %s, got rune '%c'", cursor.PositionString(), r)
	}

	eof = cursor.NextNonEmptyLine()
	if !eof {
		t.Errorf("expect EOL at %s, got rune '%c'", cursor.PositionString(), r)
	}

	r, eof = cursor.NextInLine()
	if !eof {
		t.Errorf("expect EOL at %s, got rune '%c'", cursor.PositionString(), r)
	}
}

func TestCursorNextLine1(t *testing.T) {
	text := []string{
		"lorem\n",
		"ipsum",
	}

	cursor := NewCursorFromString("example.txt", strings.Join(text, ""))

	eof1 := cursor.NextLine()
	if eof1 {
		t.Fatalf("unexpected EOF at %s", cursor.PositionString())
	}

	eof2 := cursor.NextLine()
	if !eof2 {
		t.Errorf("expected EOF at %s", cursor.PositionString())
	}
}

func TestCursorNextLine2(t *testing.T) {
	text := []string{
		"lorem\n",
		"ipsum\n",
	}

	cursor := NewCursorFromString("example.txt", strings.Join(text, ""))

	eof1 := cursor.NextLine()
	if eof1 {
		t.Fatalf("unexpected EOF at %s", cursor.PositionString())
	}

	eof2 := cursor.NextLine()
	if eof2 {
		t.Fatalf("unexpected EOF at %s", cursor.PositionString())
	}

	eof3 := cursor.NextLine()
	if !eof3 {
		t.Errorf("expected EOF at %s", cursor.PositionString())
	}
}

func TestCursorMark(t *testing.T) {
	cursor := createTestCursor2()

	s := cursor.State()
	for {
		r, _ := cursor.NextInLine()
		if r == ' ' {
			break
		}
	}

	content, ctx := cursor.Finish(s)
	got := ctx.HighlightText("here")
	expected := strings.Join([]string{
		"    1 | lorem ipsum dolor sit amet",
		"      | ^^^^^",
		"      | here",
	}, "\n")

	if got != expected {
		t.Errorf("expected:\n%s\ngot:\n%s", expected, got)
	}

	if content != "lorem" {
		t.Errorf("expected content 'lorem', got '%s'", content)
	}
}

func TestCursorPeekString(t *testing.T) {
	cursor := createTestCursor2()

	expPosition := "example.txt:1:1"
	if cursor.PositionString() != expPosition {
		t.Errorf("expected position '%s', got '%s'", expPosition, cursor.PositionString())
	}

	ctx1 := cursor.PeekString("lorem")
	if ctx1 == nil {
		t.Fatalf("expected to find 'lorem' at %s", cursor.PositionString())
	}

	expected1 := strings.Join([]string{
		"    1 | lorem ipsum dolor sit amet",
		"      | ^^^^^",
		"      | here",
	}, "\n")
	if expected1 != ctx1.HighlightText("here") {
		t.Errorf("expected:\n%s\ngot:\n%s", expected1, ctx1.HighlightText("here"))
	}

	if ctx := cursor.PeekString("ipsum"); ctx != nil {
		t.Errorf("expected to not find 'ipsum' at %s, got %s", cursor.PositionString(), ctx.HighlightText("here"))
	}
}

func TestCursorNextLiteral(t *testing.T) {
	cursor := createTestCursor1()

	expPosition := "example.txt:1:1"
	if cursor.PositionString() != expPosition {
		t.Errorf("expected position '%s', got '%s'", expPosition, cursor.PositionString())
	}

	ctx1 := cursor.PeekString("lorem")
	if ctx1 == nil {
		t.Fatalf("expected to find 'lorem' at %s", cursor.PositionString())
	}

	expected1 := strings.Join([]string{
		"    1 | lorem",
		"      | ^^^^^",
		"      | here",
	}, "\n")
	if expected1 != ctx1.HighlightText("here") {
		t.Errorf("expected:\n%s\ngot:\n%s", expected1, ctx1.HighlightText("here"))
	}

	expPosition = "example.txt:1:1"
	if cursor.PositionString() != expPosition {
		t.Errorf("expected position '%s', got '%s'", expPosition, cursor.PositionString())
	}

	ctx2 := cursor.NextString("ipsum")
	if ctx2 != nil {
		t.Fatalf("expected to not find 'ipsum' at %s, got %s", cursor.PositionString(), ctx2.HighlightText("here"))
	}

	ctx3 := cursor.NextString("lorem")
	if expected1 != ctx3.HighlightText("here") {
		t.Errorf("expected:\n%s\ngot:\n%s", expected1, ctx3.HighlightText("here"))
	}

	expPosition = "example.txt:1:6"
	if cursor.PositionString() != expPosition {
		t.Errorf("expected position '%s', got '%s'", expPosition, cursor.PositionString())
	}
}

func TestCursorNextLiteralInEndOfLine(t *testing.T) {
	cursor := createTestCursor1()

	ctx1 := cursor.NextString("lorem")
	if ctx1 == nil {
		t.Fatalf("expected to find 'lorem' at %s", cursor.PositionString())
	}

	expected1 := strings.Join([]string{
		"    1 | lorem",
		"      | ^^^^^",
		"      | here",
	}, "\n")
	if expected1 != ctx1.HighlightText("here") {
		t.Errorf("expected:\n%s\ngot:\n%s", expected1, ctx1.HighlightText("here"))
	}

	cursor.NextNonEmptyLine()

	ctx2 := cursor.NextString("ipsum")
	if ctx2 == nil {
		t.Fatalf("expected to find 'ipsum' at %s", cursor.PositionString())
	}
	expected2 := strings.Join([]string{
		"    2 | ipsum",
		"      | ^^^^^",
		"      | here",
	}, "\n")
	if expected2 != ctx2.HighlightText("here") {
		t.Errorf("expected:\n%s\ngot:\n%s", expected2, ctx2.HighlightText("here"))
	}
}

func TestCursorNext(t *testing.T) {
	cursor := createTestCursor1()

	first, _, _ := cursor.Rune()
	result := make([]rune, 0, 100)
	result = append(result, first)
	for {
		r, _, eof := cursor.Next()
		if eof {
			break
		}

		result = append(result, r)
	}

	if len(result) != 12 {
		t.Errorf("expected 12 runes, got %d", len(result))
		t.Errorf("result: [%s]", string(result))
	}
}

func TestCursorIsFirstNonWhiteChar2(t *testing.T) {
	text := []string{
		"the quick brown fox\n",
		"   jumps over the lazy dog\n",
	}

	cursor := NewCursorFromString("example.txt", strings.Join(text, ""))

	if !cursor.IsFirstNonWhiteChar() {
		t.Errorf("expected to be first non-white char at %s", cursor.PositionString())
	}

	for {
		_, eol := cursor.NextInLine()
		if eol {
			break
		}

		if cursor.IsFirstNonWhiteChar() {
			t.Errorf("expected to not be first non-white char at %s", cursor.PositionString())
		}
	}

	cursor.NextNonEmptyLine()

	if !cursor.IsFirstNonWhiteChar() {
		t.Errorf("expected to be first non-white char at %s", cursor.PositionString())
	}

	for {
		r, eol := cursor.NextInLine()
		if eol {
			break
		}

		if !cursor.IsFirstNonWhiteChar() {
			t.Errorf("expected to be first non-white char at %s", cursor.PositionString())
		}

		if r != ' ' {
			break
		}
	}

	for {
		_, eol := cursor.NextInLine()
		if eol {
			break
		}

		if cursor.IsFirstNonWhiteChar() {
			t.Errorf("expected to not be first non-white char at %s", cursor.PositionString())
		}
	}
}

func TestCursorSkipWhitespace(t *testing.T) {
	text := strings.Join([]string{
		"                lorem ipsum",
	}, "\n")

	cursor := NewCursorFromString("example.txt", text)

	_, ctx1 := cursor.CurrentChar()
	got1 := ctx1.HighlightText("here")
	exp1 := strings.Join([]string{
		"    1 |                 lorem ipsum",
		"      | ^",
		"      | here",
	}, "\n")

	if got1 != exp1 {
		t.Errorf("expected:\n%s\ngot:\n%s", exp1, got1)
	}

	cursor.SkipWhitespace()

	_, ctx2 := cursor.CurrentChar()
	got2 := ctx2.HighlightText("here")
	exp2 := strings.Join([]string{
		"    1 |                 lorem ipsum",
		"      |                 ^",
		"      |                 here",
	}, "\n")

	if got2 != exp2 {
		t.Errorf("expected:\n%s\ngot:\n%s", exp2, got2)
	}
}

func TestCursorSkipWhitespaceToNextLine(t *testing.T) {
	text := strings.Join([]string{
		"        lorem        ",
		"    ",
		"        \t\t\t        ",
		"      ipsum dolor sit amet",
	}, "\n")

	cursor := NewCursorFromString("example.txt", text)

	_, ctx1 := cursor.CurrentChar()
	got1 := ctx1.HighlightText("here")
	exp1 := strings.Join([]string{
		"    1 |         lorem        ",
		"      | ^",
		"      | here",
	}, "\n")

	if got1 != exp1 {
		t.Errorf("expected:\n%s\ngot:\n%s", exp1, got1)
	}

	cursor.SkipWhitespace()

	_, ctx2 := cursor.CurrentChar()
	got2 := ctx2.HighlightText("here")
	exp2 := strings.Join([]string{
		"    1 |         lorem        ",
		"      |         ^",
		"      |         here",
	}, "\n")

	if got2 != exp2 {
		t.Errorf("expected:\n%s\ngot:\n%s", exp2, got2)
	}

	cursor.Skip(5)
	cursor.SkipWhitespace()

	_, ctx3 := cursor.CurrentChar()
	got3 := ctx3.HighlightText("here")
	exp3 := strings.Join([]string{
		"    4 |       ipsum dolor sit amet",
		"      |       ^",
		"      |       here",
	}, "\n")

	if got3 != exp3 {
		t.Errorf("expected:\n%s\ngot:\n%s", exp3, got3)
	}
}

func TestCursorSkipWhitespaceInLine(t *testing.T) {
	text := strings.Join([]string{
		"        lorem        ",
		"    ",
		"        \t\t\t        ",
		"      ipsum dolor sit amet",
	}, "\n")

	cursor := NewCursorFromString("example.txt", text)

	_, ctx1 := cursor.CurrentChar()
	got1 := ctx1.HighlightText("here")
	exp1 := strings.Join([]string{
		"    1 |         lorem        ",
		"      | ^",
		"      | here",
	}, "\n")

	if got1 != exp1 {
		t.Errorf("expected:\n%s\ngot:\n%s", exp1, got1)
	}
	cursor.SkipWhitespaceInLine()

	_, ctx2 := cursor.CurrentChar()
	got2 := ctx2.HighlightText("here")
	exp2 := strings.Join([]string{
		"    1 |         lorem        ",
		"      |         ^",
		"      |         here",
	}, "\n")

	if got2 != exp2 {
		t.Errorf("expected:\n%s\ngot:\n%s", exp2, got2)
	}

	cursor.Skip(5)
	cursor.SkipWhitespaceInLine()

	_, ctx3 := cursor.CurrentChar()
	got3 := ctx3.HighlightText("here")
	exp3 := strings.Join([]string{
		"    1 |         lorem        <EOL LF>",
		"      |                      ^^^^^^^^",
		"      |                      here",
	}, "\n")

	if got3 != exp3 {
		t.Errorf("expected:\n%s\ngot:\n%s", exp3, got3)
	}
}

func TestCursorEndWithLastEmptyLine(t *testing.T) {
	cursor := createTestCursor1()

	for range 5 {
		eol, eof := cursor.End()
		if eol || eof {
			t.Fatalf("unexpected EOL/EOF at %s  (EOL=%v, EOF=%v)", cursor.PositionString(), eol, eof)
		}

		cursor.NextInLine()
	}

	{
		eol, eof := cursor.End()
		if !eol {
			t.Errorf("expect EOL at %s", cursor.PositionString())
		}

		if eof {
			t.Fatalf("unexpected EOF at %s", cursor.PositionString())
		}
	}

	cursor.NextLine()
	for range 5 {
		eol, eof := cursor.End()
		if eol || eof {
			t.Fatalf("unexpected EOL/EOF at %s  (EOL=%v, EOF=%v)", cursor.PositionString(), eol, eof)
		}

		cursor.NextInLine()
	}

	{
		eol, eof := cursor.End()
		if !eol {
			t.Errorf("expect EOL at %s", cursor.PositionString())
		}

		if eof {
			t.Fatalf("unexpected EOF at %s", cursor.PositionString())
		}
	}

	cursor.NextLine()
	{
		eol, eof := cursor.End()
		if !eol {
			t.Errorf("expect EOL at %s", cursor.PositionString())
		}

		if !eof {
			t.Fatalf("expect EOF at %s", cursor.PositionString())
		}
	}

}

func TestCursorCurrentLineAndEOL(t *testing.T) {
	cursor := createTestCursor2()

	line1, ctx1 := cursor.CurrentLine()
	content1 := "lorem ipsum dolor sit amet"
	if line1 != content1 {
		t.Errorf("expected line content:\n%s\ngot:\n%s", content1, line1)
	}

	exp1 := strings.Join([]string{
		"    1 | lorem ipsum dolor sit amet",
		"      | ^^^^^^^^^^^^^^^^^^^^^^^^^^",
		"      | here",
	}, "\n")
	got1 := ctx1.HighlightText("here")
	if got1 != exp1 {
		t.Errorf("expected:\n%s\ngot:\n%s", exp1, got1)
	}

	eol1 := cursor.CurrentEOL()
	if string(eol1) != "\n" {
		t.Errorf("expected EOL '\\n', got '%s'", string(eol1))
	}

	cursor.NextLine()

	line2, ctx2 := cursor.CurrentLine()
	content2 := "consectetur adipiscing elit"
	if line2 != content2 {
		t.Errorf("expected line content:\n%s\ngot:\n%s", content2, line2)
	}

	exp2 := strings.Join([]string{
		"    2 | consectetur adipiscing elit",
		"      | ^^^^^^^^^^^^^^^^^^^^^^^^^^^",
		"      | here",
	}, "\n")
	got2 := ctx2.HighlightText("here")
	if got2 != exp2 {
		t.Errorf("expected:\n%s\ngot:\n%s", exp2, got2)
	}

	eol2 := cursor.CurrentEOL()
	if string(eol2) != "\n" {
		t.Errorf("expected EOL '\\n', got '%s'", string(eol2))
	}

	cursor.NextLine()
	cursor.NextLine()
	cursor.NextLine()
	cursor.NextLine()

	line6, ctx6 := cursor.CurrentLine()
	content6 := "ut enim ad minim veniam"
	if line6 != content6 {
		t.Errorf("expected line content:\n%s\ngot:\n%s", content6, line6)
	}

	exp6 := strings.Join([]string{
		"    6 | ut enim ad minim veniam",
		"      | ^^^^^^^^^^^^^^^^^^^^^^^",
		"      | here",
	}, "\n")
	got6 := ctx6.HighlightText("here")
	if got6 != exp6 {
		t.Errorf("expected:\n%s\ngot:\n%s", exp6, got6)
	}

	eol6 := cursor.CurrentEOL()
	if string(eol6) != "\r\n" {
		t.Errorf("expected EOL '\\r\\n', got '%s'", string(eol6))
	}
}
