package context

import (
	"testing"

	"bytes"
	"strings"
)

func TestLineContentBasicInfo(t *testing.T) {
	s := "lorem ipsum dolor sit amet"
	line := NewLineFromBytes(42, []byte(s), nil)
	if line.Line != 42 {
		t.Errorf("expected line number 42, got %d", line.Line)
	}

	if line.String() != s {
		t.Errorf("expected line content '%s', got '%s'", s, line.String())
	}

	if line.Length() != len(s) {
		t.Errorf("expected line length %d, got %d", len(s), line.Length())
	}
}

func TestLineContextMark(t *testing.T) {
	s := "lorem ipsum dolor sit amet"
	line := NewLineFromBytes(42, []byte(s), nil)
	ctx := line.markLine(6, 11)

	if ctx.StringContent() != s {
		t.Errorf("expected context content '%s', got '%s'", s, ctx.StringContent())
	}

	if ctx.Length() != len(s) {
		t.Errorf("expected context length %d, got %d", len(s), ctx.Length())
	}

	if len(ctx.Highlights) != 1 {
		t.Errorf("expected 1 highlight, got %d", len(ctx.Highlights))
	}
}

func TestReadFileData(t *testing.T) {
	filename := "example.txt"
	content := [][]byte{
		[]byte("lorem ipsum\n"),
		[]byte("dolor sit amet\r\n"),
		[]byte("consectetur adipiscing elit\n"),
		[]byte("\n"),
		[]byte("sed do eiusmod tempor incididunt\r\n"),
		[]byte("\r\n"),
		[]byte("ut labore et dolore magna aliqua\n"),
	}

	ctx := ReadFileData(filename, bytes.Join(content, nil))
	if ctx.Filename != filename {
		t.Errorf("expected filename %s, got %s", filename, ctx.Filename)
	}

	if ctx.Lines() != 8 {
		t.Errorf("expected 8 lines, got %d", ctx.Lines())
	}

	if ctx.Line(0).String() != "lorem ipsum" {
		t.Errorf("expected first line to be 'lorem ipsum', got '%s'", ctx.Line(0).String())
	}

	if !bytes.Equal(ctx.Line(0).EOLBytes(), EolLF) {
		t.Errorf("expected first line EOL to be LF, got %s", ctx.Line(0).EOLBytes())
	}

	if ctx.Line(1).String() != "dolor sit amet" {
		t.Errorf("expected second line to be 'dolor sit amet', got '%s'", ctx.Line(1).String())
	}

	if !bytes.Equal(ctx.Line(1).EOLBytes(), EolCRLF) {
		t.Errorf("expected second line EOL to be CRLF, got %s", ctx.Line(1).EOLBytes())
	}

	if ctx.Line(6).String() != "ut labore et dolore magna aliqua" {
		t.Errorf("expected last line to be 'ut labore et dolore magna aliqua', got '%s'", ctx.Line(6).String())
	}

	if c, eol := ctx.Rune(0, 0); c != 'l' || eol {
		t.Errorf("expected first rune to be 'l', got '%c', eof=%v", c, eol)
	}

	if c, eol := ctx.Rune(4, 5); c != 'o' || eol {
		t.Errorf("expected rune at (4, 5) to be 'o', got '%c', eof=%v", c, eol)
	}

	if c, eol := ctx.Rune(0, 100); c != 0 || !eol {
		t.Errorf("expected rune at (0, 100) to be EOF, got '%c', eof=%v", c, eol)
	}

	if c, eol := ctx.Rune(100, 0); c != 0 || !eol {
		t.Errorf("expected rune at (100, 0) to be EOF, got '%c', eof=%v", c, eol)
	}
}

func TestReadFileDataWithEmptyLine(t *testing.T) {
	filename := "example.txt"
	data1 := strings.Join([]string{
		"lorem ipsum",
		"",
	}, "\n")

	ctx1 := ReadFileString(filename, data1)
	if ctx1.Lines() != 2 {
		t.Errorf("expected 2 lines, got %d", ctx1.Lines())
	}

	data2 := strings.Join([]string{
		"lorem ipsum",
	}, "\n")

	ctx2 := ReadFileString(filename, data2)
	if ctx2.Lines() != 1 {
		t.Errorf("expected 1 line, got %d", ctx2.Lines())
	}
}

func TestReadFileWithoutEndOfLine(t *testing.T) {
	filename := "example.txt"
	data := []byte("the quick brown fox jumps over the lazy dog")

	ctx := ReadFileData(filename, data)
	if ctx.Lines() != 1 {
		t.Errorf("expected 1 line, got %d", ctx.Lines())
	}
}

func TestLineContentEOLString(t *testing.T) {
	cases := []struct {
		line     string
		eol      string
		expected string
	}{
		{line: "lorem ipsum", eol: "", expected: "<EOF>"},
		{line: "lorem ipsum", eol: "\n", expected: "<EOL LF>"},
		{line: "lorem ipsum", eol: "\r", expected: "<EOL CR>"},
		{line: "lorem ipsum", eol: "\r\n", expected: "<EOL CRLF>"},
		{line: "lorem ipsum", eol: "\x00", expected: "<EOL 0x00>"},
	}

	for _, c := range cases {
		line := NewLineFromString(0, c.line, c.eol)
		if line.EOLString() != c.expected {
			t.Errorf("expected EOL string '%s', got '%s'", c.expected, line.EOLString())
		}
	}
}
