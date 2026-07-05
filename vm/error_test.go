package vm

import (
	"io"
	"testing"

	"bytes"
	"errors"
	"strings"

	"github.com/flily/go-brainfuck/context"
)

func createTestFile() *context.FileContext {
	parts := [][]byte{
		[]byte("lorem ipsum dolor sit amet\n"),
		[]byte("consectetur adipiscing elit\n"),
		[]byte("\n"),
		[]byte("sed do eiusmod tempor incididunt\n"),
		[]byte("ut labore et dolore magna aliqua\n"),
		[]byte("ut enim ad minim veniam\r\n"),
		[]byte("\r\n"),
		[]byte("quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat\n"),
		[]byte("duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur\n"),
		[]byte("excepteur sint occaecat cupid atat non proident\n"),
	}

	content := bytes.Join(parts, []byte{})
	file := context.ReadFileData("example.txt", content)
	return file
}

func TestReasonString(t *testing.T) {
	cases := []struct {
		reason   Reason
		expected string
	}{
		{ReasonInvalid, "unknown"},
		{ReasonCallStackOverflow, "call stack overflow"},
	}

	for _, c := range cases {
		s1 := c.reason.String()
		if s1 != c.expected {
			t.Fatalf("reason string mismatch for reason %d: expected %s, got %s",
				c.reason, c.expected, s1)
		}

		s2 := c.reason.Error()
		if s2 != c.expected {
			t.Fatalf("reason error string mismatch for reason %d: expected %s, got %s",
				c.reason, c.expected, s2)
		}
	}
}

func TestReasonIsError(t *testing.T) {
	r := ReasonCallStackOverflow
	if !errors.Is(r, ReasonCallStackOverflow) {
		t.Fatalf("reason is not of type %s", ReasonCallStackOverflow)
	}

	derived := r.OnError(nil, "lorem ipsum")
	if !errors.Is(derived, r) {
		t.Fatalf("derived error is not of type %s", r)
	}
}

func TestRuntimeErrorWithoutContext(t *testing.T) {
	err := ReasonCallStackOverflow.OnError(nil, "lorem ipsum").
		With("dolor sit amet")
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	expected := strings.Join([]string{
		"[error]: lorem ipsum",
		"    dolor sit amet",
	}, "\n")
	if merr := err.Error(); merr != expected {
		t.Fatalf("error message mismatch, expected:\n%s\ngot:\n%s", expected, merr)
	}
}

func TestRuntimeErrorWithoutContextNote(t *testing.T) {
	err := ReasonCallStackOverflow.OnError(nil, "lorem ipsum")
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	expected := strings.Join([]string{
		"[error]: lorem ipsum",
	}, "\n")
	if merr := err.Error(); merr != expected {
		t.Fatalf("error message mismatch, expected:\n%s\ngot:\n%s", expected, merr)
	}
}

func TestRuntimeErrorWithoutContextMessage(t *testing.T) {
	err := ReasonCallStackOverflow.OnError(nil, "").
		With("dolor sit amet")
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	expected := strings.Join([]string{
		"[error]: call stack overflow",
		"    dolor sit amet",
	}, "\n")
	if merr := err.Error(); merr != expected {
		t.Fatalf("error message mismatch, expected:\n%s\ngot:\n%s", expected, merr)
	}
}

func TestRuntimeErrorWithoutContextMessageNote(t *testing.T) {
	err := ReasonCallStackOverflow.OnError(nil, "")
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	expected := strings.Join([]string{
		"[error]: call stack overflow",
	}, "\n")
	if merr := err.Error(); merr != expected {
		t.Fatalf("error message mismatch, expected:\n%s\ngot:\n%s", expected, merr)
	}
}

func TestRuntimeErrorWithContext(t *testing.T) {
	fd := createTestFile()
	line := fd.LineContext(3)
	ctx := line.Mark(7, 14)

	err := ReasonCallStackOverflow.OnError(ctx, "lorem ipsum").
		With("dolor sit amet")
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	expected := strings.Join([]string{
		"example.txt:4:8: error: lorem ipsum",
		"    4 | sed do eiusmod tempor incididunt",
		"      |        ^^^^^^^",
		"      |        dolor sit amet",
	}, "\n")
	if merr := err.Error(); merr != expected {
		t.Fatalf("error message mismatch, expected:\n%s\ngot:\n%s", expected, merr)
	}
}

func TestRuntimeErrorWithContextNoMessage(t *testing.T) {
	fd := createTestFile()
	line := fd.LineContext(3)
	ctx := line.Mark(7, 14)

	err := ReasonCallStackOverflow.OnError(ctx, "").
		With("dolor sit amet")
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	expected := strings.Join([]string{
		"example.txt:4:8: error: call stack overflow",
		"    4 | sed do eiusmod tempor incididunt",
		"      |        ^^^^^^^",
		"      |        dolor sit amet",
	}, "\n")
	if merr := err.Error(); merr != expected {
		t.Fatalf("error message mismatch, expected:\n%s\ngot:\n%s", expected, merr)
	}
}

func TestRuntimeErrorWithContextNoNote(t *testing.T) {
	fd := createTestFile()
	line := fd.LineContext(3)
	ctx := line.Mark(7, 14)

	err := ReasonCallStackOverflow.OnError(ctx, "lorem ipsum")
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	expected := strings.Join([]string{
		"example.txt:4:8: error: lorem ipsum",
		"    4 | sed do eiusmod tempor incididunt",
		"      |        ^^^^^^^",
		"      |        ",
	}, "\n")
	if merr := err.Error(); merr != expected {
		t.Fatalf("error message mismatch, expected:\n%s\ngot:\n%s", expected, merr)
	}
}

func TestRuntimeErrorWithContextNoMessageNoNote(t *testing.T) {
	fd := createTestFile()
	line := fd.LineContext(3)
	ctx := line.Mark(7, 14)

	err := ReasonCallStackOverflow.OnError(ctx, "")
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	expected := strings.Join([]string{
		"example.txt:4:8: error: call stack overflow",
		"    4 | sed do eiusmod tempor incididunt",
		"      |        ^^^^^^^",
		"      |        ",
	}, "\n")
	if merr := err.Error(); merr != expected {
		t.Fatalf("error message mismatch, expected:\n%s\ngot:\n%s", expected, merr)
	}
}

func TestRuntimeErrorOnDifferentLevels(t *testing.T) {
	base := ReasonCallStackOverflow
	cases := []struct {
		err *RuntimeError
		exp string
	}{
		{
			err: base.OnNote(nil, "lorem ipsum"),
			exp: "[note]: lorem ipsum",
		},
		{
			err: base.OnRemark(nil, "lorem ipsum"),
			exp: "[remark]: lorem ipsum",
		},
		{
			err: base.OnWarning(nil, "lorem ipsum"),
			exp: "[warning]: lorem ipsum",
		},
		{
			err: base.OnError(nil, "lorem ipsum"),
			exp: "[error]: lorem ipsum",
		},
		{
			err: base.OnFatal(nil, "lorem ipsum"),
			exp: "[fatal]: lorem ipsum",
		},
	}

	for _, c := range cases {
		if merr := c.err.Error(); merr != c.exp {
			t.Fatalf("error message mismatch for level %s, expected:\n%s\ngot:\n%s",
				c.err.Level, c.exp, merr)
		}

		if !errors.Is(c.err, base) {
			t.Fatalf("error is not of type %s, got %T", base, c.err)
		}

		if errors.Is(c.err, io.EOF) {
			t.Fatalf("error should not be of type %T", io.EOF)
		}
	}
}
