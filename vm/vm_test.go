package vm

import (
	"testing"

	"strings"

	"github.com/flily/go-brainfuck/infra"
)

func newSimpleCodeMap(insts ...Instruction) *CodeMap {
	codemap := infra.NewCodeMap()
	codemap.Codes = infra.InstructionsToCodes(insts)

	return codemap
}

func TestVMBasicMethods(t *testing.T) {
	m := New[uint8](32, 32)
	m.LoadCode(newSimpleCodeMap(InstructionAdd, InstructionAdd, InstructionAdd))

	ip1, dp1, sp1 := m.Registers()
	if ip1 != 0 || dp1 != 0 || sp1 != 0 {
		t.Fatalf("expected registers to be (0, 0, 0), got (%d, %d, %d)", ip1, dp1, sp1)
	}

	peekIP, err := m.PeekIP()
	if err == nil {
		t.Fatalf("expected error when peeking IP with empty stack, got nil")
	}

	if peekIP != -1 {
		t.Fatalf("expected peekIP to be -1 when stack is empty, got %d", peekIP)
	}

	m.IP = 2
	if err := m.PushIP(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if m.SP != 1 {
		t.Fatalf("expected SP to be 1 after PushIP, got %d", m.SP)
	}

	peekIP, err = m.PeekIP()
	if err != nil {
		t.Fatalf("unexpected error when peeking IP: %v", err)
	}

	if peekIP != 2 {
		t.Fatalf("expected peekIP to be 2, got %d", peekIP)
	}

	m.Reset()
	peekIP, err = m.PeekIP()
	if err == nil {
		t.Fatalf("expected error when peeking IP after reset, got nil")
	}

	if peekIP != -1 {
		t.Fatalf("expected peekIP to be -1 after reset, got %d", peekIP)
	}
}

func TestVMGetCurrentCodeOnErrorState(t *testing.T) {
	m := New[uint8](32, 32)
	m.LoadCode(newSimpleCodeMap(InstructionAdd, InstructionAdd, InstructionAdd))

	peekIP, err := m.PeekIP()
	if err == nil {
		t.Fatalf("expected error when peeking IP with empty stack, got nil")
	}

	if peekIP != -1 {
		t.Fatalf("expected peekIP to be -1 when stack is empty, got %d", peekIP)
	}

	m.IP = -100
	code := m.GetCurrentCode()
	if code != nil {
		t.Fatalf("expected GetCurrentCode to return nil for invalid IP, got %v", code)
	}

	m.IP = 100
	code = m.GetCurrentCode()
	if code != nil {
		t.Fatalf("expected GetCurrentCode to return nil for invalid IP, got %v", code)
	}
}

func TestVMExecuteInstructionNotInSet(t *testing.T) {
	code := "+++"
	codemap, err := parseCodeLiteral(code, "test.bf")
	if err != nil {
		t.Fatalf("error parsing code: %v", err)
	}

	vm := New[uint8](32, 32)
	vm.LoadCode(codemap)

	err = vm.Run()
	if err == nil {
		t.Fatalf("expected error when executing instruction not in set, got nil")
	}

	expected := strings.Join([]string{
		"test.bf:1:1: fatal: unsupported instruction",
		"    1 | +++",
		"      | ^",
		"      | instruction='+' (0x2b)",
	}, "\n")
	if err.Error() != expected {
		t.Fatalf("error message mismatch, expected:\n%s\ngot:\n%s", expected, err.Error())
	}
}
