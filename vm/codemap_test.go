package vm

import (
	"testing"

	"slices"
)

func TestCodeMapGetNext(t *testing.T) {
	codemap := NewCodeMap()
	codemap.AddFile(nil)

	codemap.Codes = InstructionsToCodes([]Instruction{
		InstructionAdd,
		InstructionSub,
		InstructionPointerDec,
	})
	codemap.Next = []int{
		1,
		2,
		3,
	}

	if r := codemap.GetNext(0); r != 1 {
		t.Fatalf("expected next for 0 to be 1, got %d", r)
	}

	if r := codemap.GetNext(1); r != 2 {
		t.Fatalf("expected next for 1 to be 2, got %d", r)
	}

	if r := codemap.GetNext(-1); r != -1 {
		t.Fatalf("expected next for -1 to be -1, got %d", r)
	}

	if r := codemap.GetNext(3); r != -1 {
		t.Fatalf("expected next for 3 to be -1, got %d", r)
	}
}

func TestCodeMapEqualsEqual(t *testing.T) {
	codemap1 := NewCodeMap()
	codemap1.Codes = InstructionsToCodes([]Instruction{
		InstructionAdd,
		InstructionSub,
	})
	codemap1.Next = []int{
		1,
		2,
	}

	codemap2 := NewCodeMap()
	codemap2.Codes = InstructionsToCodes([]Instruction{
		InstructionAdd,
		InstructionSub,
	})
	codemap2.Next = []int{
		1,
		2,
	}

	if !codemap1.CodeEquals(codemap2) {
		t.Fatalf("expected codemaps to be equal")
	}
}

func TestCodeMapEqualsNotEqualOnLength(t *testing.T) {
	codemap1 := NewCodeMap()
	codemap1.Codes = InstructionsToCodes([]Instruction{
		InstructionAdd,
		InstructionSub,
	})
	codemap1.Next = []int{
		1,
		2,
	}

	codemap2 := NewCodeMap()
	codemap2.Codes = InstructionsToCodes([]Instruction{
		InstructionAdd,
	})
	codemap2.Next = []int{
		1,
	}

	if codemap1.CodeEquals(codemap2) {
		t.Fatalf("expected codemaps to be not equal on length")
	}
}

func TestCodeMapEqualsNotEqualOnInstruction(t *testing.T) {
	codemap1 := NewCodeMap()
	codemap1.Codes = InstructionsToCodes([]Instruction{
		InstructionAdd,
		InstructionSub,
	})
	codemap1.Next = []int{
		1,
		2,
	}

	codemap2 := NewCodeMap()
	codemap2.Codes = InstructionsToCodes([]Instruction{
		InstructionAdd,
		InstructionPointerDec,
	})
	codemap2.Next = []int{
		1,
		2,
	}

	if codemap1.CodeEquals(codemap2) {
		t.Fatalf("expected codemaps to be not equal on instruction")
	}
}

func TestCodeMapEqualsNotEqualOnNext(t *testing.T) {
	codemap1 := NewCodeMap()
	codemap1.Codes = InstructionsToCodes([]Instruction{
		InstructionAdd,
		InstructionSub,
	})
	codemap1.Next = []int{
		1,
		2,
	}

	codemap2 := NewCodeMap()
	codemap2.Codes = InstructionsToCodes([]Instruction{
		InstructionAdd,
		InstructionSub,
	})
	codemap2.Next = []int{
		1,
		3,
	}

	if codemap1.CodeEquals(codemap2) {
		t.Fatalf("expected codemaps to be not equal on next")
	}
}

func TestCodeMapSnapshot(t *testing.T) {
	codemap := NewCodeMap()
	codemap.Codes = InstructionsToCodes([]Instruction{
		InstructionAdd,
		InstructionSub,
		InstructionAdd,
		InstructionPointerDec,
		InstructionPointerInc,
	})
	codemap.Next = []int{
		1,
		2,
		3,
		4,
	}

	{
		snapshot := codemap.Snapshot(1, 1, 2)
		expected := InstructionsToCodes([]Instruction{
			InstructionAdd,
			InstructionSub,
			InstructionAdd,
			InstructionPointerDec,
		})

		if !slices.Equal(snapshot, expected) {
			t.Fatalf("expected snapshot to be %v, got %v", expected, snapshot)
		}
	}

	{
		snapshot := codemap.Snapshot(0, 1, 2)
		expected := InstructionsToCodes([]Instruction{
			InstructionAdd,
			InstructionSub,
			InstructionAdd,
		})

		if !slices.Equal(snapshot, expected) {
			t.Fatalf("expected snapshot to be %v, got %v", expected, snapshot)
		}
	}
}
