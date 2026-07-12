package config

import (
	"bytes"
	"testing"

	"flag"
	"strings"
)

func TestMemoryUnitType(t *testing.T) {
	cases := []struct {
		input    string
		expected MemoryUnitType
	}{
		{"uint8", MemoryUnitTypeUint8},
		{"UINT8", MemoryUnitTypeUint8},
		{"uint16", MemoryUnitTypeUint16},
		{"UINT16", MemoryUnitTypeUint16},
		{"uint32", MemoryUnitTypeUint32},
		{"UINT32", MemoryUnitTypeUint32},
		{"uint64", MemoryUnitTypeUint64},
		{"UINT64", MemoryUnitTypeUint64},
		{"int8", MemoryUnitTypeInt8},
		{"INT8", MemoryUnitTypeInt8},
		{"int16", MemoryUnitTypeInt16},
		{"INT16", MemoryUnitTypeInt16},
		{"int32", MemoryUnitTypeInt32},
		{"INT32", MemoryUnitTypeInt32},
		{"int64", MemoryUnitTypeInt64},
		{"INT64", MemoryUnitTypeInt64},
	}

	for _, c := range cases {
		set := flag.NewFlagSet("test", flag.ContinueOnError)

		var mt MemoryUnitType
		set.Var(&mt, "word", "data type of memory unit cell")
		err := set.Parse([]string{"-word", c.input})
		if err != nil {
			t.Fatalf("unexpected error parsing input '%s': %v", c.input, err)
		}

		if mt != c.expected {
			t.Errorf("unexpected MemoryUnitType for input '%s': got %v, want %v", c.input, mt, c.expected)
		}

		lower := strings.ToLower(c.input)
		if mt.String() != lower {
			t.Errorf("unexpected string representation for MemoryUnitType %v: got '%s', want '%s'", mt, mt.String(), c.input)
		}
	}
}

func TestMemoryUnitTypeInvalid(t *testing.T) {
	mt := MemoryUnitTypeInvalid

	if mt.String() != "unknown" {
		t.Errorf("unexpected string representation for MemoryUnitTypeInvalid: got '%s', want 'unknown'", mt.String())
	}

	out := bytes.NewBuffer(nil)

	set := flag.NewFlagSet("test", flag.ContinueOnError)
	set.SetOutput(out)

	set.Var(&mt, "word", "data type of memory unit cell")
	err := set.Parse([]string{"-word", "invalid"})
	if err == nil {
		t.Fatalf("expected error parsing invalid input, got nil")
	}

	if len(out.Bytes()) <= 0 {
		t.Errorf("expected error output, got empty")
	}
}

func TestEndian(t *testing.T) {
	cases := []struct {
		input    string
		expected Endian
	}{
		{"big", EndianBig},
		{"BIG", EndianBig},
		{"be", EndianBig},
		{"BE", EndianBig},
		{"big-endian", EndianBig},
		{"BIG-ENDIAN", EndianBig},
		{"little", EndianLittle},
		{"LITTLE", EndianLittle},
		{"le", EndianLittle},
		{"LE", EndianLittle},
		{"little-endian", EndianLittle},
		{"LITTLE-ENDIAN", EndianLittle},
	}

	for _, c := range cases {
		set := flag.NewFlagSet("test", flag.ContinueOnError)

		var e Endian
		set.Var(&e, "endian", "endian type for input/output")
		err := set.Parse([]string{"-endian", c.input})
		if err != nil {
			t.Fatalf("unexpected error parsing input '%s': %v", c.input, err)
		}

		if e != c.expected {
			t.Errorf("unexpected Endian for input '%s': got %v, want %v", c.input, e, c.expected)
		}

		expected := ""
		switch c.expected {
		case EndianBig:
			expected = "big-endian"
		case EndianLittle:
			expected = "little-endian"
		}

		if e.String() != expected {
			t.Errorf("unexpected string representation for Endian %v: got '%s', want '%s'", e, e.String(), expected)
		}
	}
}

func TestEndianInvalid(t *testing.T) {
	e := Endian(1000)

	if e.String() != "unknown" {
		t.Errorf("unexpected string representation for invalid Endian: got '%s', want 'unknown'", e.String())
	}

	out := bytes.NewBuffer(nil)

	set := flag.NewFlagSet("test", flag.ContinueOnError)
	set.SetOutput(out)

	set.Var(&e, "endian", "endian type for input/output")
	err := set.Parse([]string{"-endian", "invalid"})
	if err == nil {
		t.Fatalf("expected error parsing invalid input, got nil")
	}

	if len(out.Bytes()) <= 0 {
		t.Errorf("expected error output, got empty")
	}
}
