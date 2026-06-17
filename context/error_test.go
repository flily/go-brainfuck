package context

import (
	"testing"

	"strings"
)

func TestErrorLevelNames(t *testing.T) {
	tests := []struct {
		level    ErrorLevel
		expected string
	}{
		{Note, "note"},
		{Warning, "warning"},
		{Error, "error"},
	}

	for _, test := range tests {
		if test.level.String() != test.expected {
			t.Fatalf("error level name mismatch for level %d: expected %s, got %s",
				test.level, test.expected, test.level.String())
		}
	}
}

func TestErrorMessage(t *testing.T) {
	fd := createTestFile1()

	line := fd.LineContext(3)
	ctx := line.Mark(7, 14)

	err := ctx.Error("the quick brown fox").With("jumps over the lazy dog")

	var _ DiagnosticInfo = err

	expected := strings.Join([]string{
		"example.txt:4:8: error: the quick brown fox",
		"    4 | sed do eiusmod tempor incididunt",
		"      |        ^^^^^^^",
		"      |        jumps over the lazy dog",
	}, "\n")

	if err.Error() != expected {
		t.Fatalf("error message mismatch, expected:\n%s\ngot:\n%s", expected, err.Error())
	}
}

func TestDiagnosticLevelString(t *testing.T) {
	fd := createTestFile1()

	line := fd.LineContext(3)
	ctx := line.Mark(7, 14)

	{
		err := ctx.Note("the quick brown fox").With("jumps over the lazy dog")
		expected := strings.Join([]string{
			"example.txt:4:8: note: the quick brown fox",
			"    4 | sed do eiusmod tempor incididunt",
			"      |        ^^^^^^^",
			"      |        jumps over the lazy dog",
		}, "\n")

		if err.Error() != expected {
			t.Fatalf("error message mismatch, expected:\n%s\ngot:\n%s", expected, err.Error())
		}
	}

	{
		err := ctx.Warning("the quick brown fox").With("jumps over the lazy dog")
		expected := strings.Join([]string{
			"example.txt:4:8: warning: the quick brown fox",
			"    4 | sed do eiusmod tempor incididunt",
			"      |        ^^^^^^^",
			"      |        jumps over the lazy dog",
		}, "\n")

		if err.Error() != expected {
			t.Fatalf("error message mismatch, expected:\n%s\ngot:\n%s", expected, err.Error())
		}
	}

	{
		err := ctx.Error("the quick brown fox").With("jumps over the lazy dog")
		expected := strings.Join([]string{
			"example.txt:4:8: error: the quick brown fox",
			"    4 | sed do eiusmod tempor incididunt",
			"      |        ^^^^^^^",
			"      |        jumps over the lazy dog",
		}, "\n")

		if err.Error() != expected {
			t.Fatalf("error message mismatch, expected:\n%s\ngot:\n%s", expected, err.Error())
		}
	}
}

func TestDiagnosticComboMessage(t *testing.T) {
	fd := createTestFile1()

	line1 := fd.LineContext(3)
	ctx1 := line1.Mark(7, 14)

	line2 := fd.LineContext(5)
	ctx2 := line2.Mark(11, 16)

	err := ctx1.Error("the quick brown fox").
		With("jumps over the lazy dog").
		For(ctx2.Note("lorem ipsum").With("dolor sit amet"))

	var _ DiagnosticInfo = err

	if lv := err.Level(); lv != Error {
		t.Fatalf("expected diagnostic combo level to be Error, got %s", lv.String())
	}

	expected := strings.Join([]string{
		"example.txt:4:8: error: the quick brown fox",
		"    4 | sed do eiusmod tempor incididunt",
		"      |        ^^^^^^^",
		"      |        jumps over the lazy dog",
		"example.txt:6:12: note: lorem ipsum",
		"    6 | ut enim ad minim veniam",
		"      |            ^^^^^",
		"      |            dolor sit amet",
	}, "\n")

	if err.Error() != expected {
		t.Fatalf("diagnostic combo message mismatch, expected:\n%s\ngot:\n%s", expected, err.Error())
	}
}

func TestDiagnosticContainerAdd(t *testing.T) {
	fd := createTestFile1()

	line1 := fd.LineContext(3)
	ctx1 := line1.Mark(7, 14)

	line2 := fd.LineContext(4)
	ctx2 := line2.Mark(13, 19)

	line3 := fd.LineContext(5)
	ctx3 := line3.Mark(11, 16)

	container := NewDiagnosticContainer(Warning)

	err1 := ctx1.Error("the quick brown fox").With("jumps over the lazy dog")
	err2 := ctx2.Warning("lorem ipsum").With("dolor sit amet")
	err3 := ctx3.Note("consectetur adipiscing elit").With("sed do eiusmod tempor")

	if lv := container.Level(); lv != Ignored {
		t.Fatalf("expected container level to be Ignored, got %s", lv.String())
	}

	if e := container.Add(err3); e != nil {
		t.Fatalf("did not expect err3 to raise error level")
	}

	if lv := container.Level(); lv != Note {
		t.Fatalf("expected container level to be Note, got %s", lv.String())
	}

	if e := container.Add(err2); e == nil {
		t.Fatalf("expected err2 to raise error level")
	}

	if lv := container.Level(); lv != Warning {
		t.Fatalf("expected container level to be Warning, got %s", lv.String())
	}

	if e := container.Add(err1); e == nil {
		t.Fatalf("expected err1 to raise error level")
	}

	if lv := container.Level(); lv != Error {
		t.Fatalf("expected container level to be Error, got %s", lv.String())
	}

	c0 := container.Count(Ignored)
	c1 := container.Count(Note)
	c2 := container.Count(Warning)
	c3 := container.Count(Error)

	if c0 != 3 || c1 != 3 || c2 != 2 || c3 != 1 {
		t.Fatalf("diagnostic count mismatch, expected (3,3,2,1), got (%d,%d,%d,%d)", c0, c1, c2, c3)
	}
}

func TestDiagnosticContaineMerge(t *testing.T) {
	fd := createTestFile1()

	line1 := fd.LineContext(3)
	ctx1 := line1.Mark(7, 14)

	line2 := fd.LineContext(4)
	ctx2 := line2.Mark(13, 19)

	line3 := fd.LineContext(5)
	ctx3 := line3.Mark(11, 16)

	container := NewDiagnosticContainer(Warning)

	err1 := ctx1.Error("the quick brown fox").With("jumps over the lazy dog")
	err2 := ctx2.Warning("lorem ipsum").With("dolor sit amet")
	err3 := ctx3.Note("consectetur adipiscing elit").With("sed do eiusmod tempor")

	c1 := err1.For(err2).ToContainer()
	c2 := err3.ToContainer()

	_ = container.Merge(c1)
	_ = container.Merge(c2)

	if container.Count(Ignored) != 2 {
		t.Fatalf("diagnostic count mismatch for Ignored, expected 2, got %d", container.Count(Ignored))
	}

	expected := strings.Join([]string{
		"example.txt:4:8: error: the quick brown fox",
		"    4 | sed do eiusmod tempor incididunt",
		"      |        ^^^^^^^",
		"      |        jumps over the lazy dog",
		"example.txt:5:14: warning: lorem ipsum",
		"    5 | ut labore et dolore magna aliqua",
		"      |              ^^^^^^",
		"      |              dolor sit amet",
		"example.txt:6:12: note: consectetur adipiscing elit",
		"    6 | ut enim ad minim veniam",
		"      |            ^^^^^",
		"      |            sed do eiusmod tempor",
	}, "\n")
	if container.Error() != expected {
		t.Fatalf("diagnostic container message mismatch, expected:\n%s\ngot:\n%s", expected, container.Error())
	}
}
