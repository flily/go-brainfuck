package config

import (
	"testing"
)

func TestConfigureIsStandard(t *testing.T) {
	c1 := Configure(0x01020304)
	if c1.IsStandard() {
		t.Fatalf("expected c1 to not be standard")
	}

	c2 := Configure(0x00010203)
	if !c2.IsStandard() {
		t.Fatalf("expected c2 to be standard")
	}
}

func TestGenericConfigureInterface(t *testing.T) {
	gc := NewGenericConfigure()
	var _ ConfigureContainer = gc

	if len(gc) != 0 {
		t.Fatalf("expected generic configure to be empty, got %d", len(gc))
	}
}

func TestGenericConfigureBoolean(t *testing.T) {
	generic := GenericConfigure{
		1: true,
		2: int64(42),
		3: uint64(42),
	}

	v1, ok := generic.GetBoolean(1)
	if !ok {
		t.Fatalf("expected to find boolean configure for 1")
	}

	if v1 != true {
		t.Fatalf("expected boolean configure for 1 to be true, got %v", v1)
	}

	v2, ok := generic.GetInt(1)
	if ok {
		t.Fatalf("expected not to find int configure for 1")
	}

	if v2 != 0 {
		t.Fatalf("expected int configure for 1 to be 0, got %v", v2)
	}

	v3, ok := generic.GetUint(1)
	if ok {
		t.Fatalf("expected not to find uint configure for 1")
	}

	if v3 != 0 {
		t.Fatalf("expected uint configure for 1 to be 0, got %v", v3)
	}
}

func TestGenericConfigureInt(t *testing.T) {
	generic := GenericConfigure{
		1: true,
		2: int64(42),
		3: uint64(42),
	}

	v1, ok := generic.GetInt(2)
	if !ok {
		t.Fatalf("expected to find int configure for 2")
	}

	if v1 != 42 {
		t.Fatalf("expected int configure for 2 to be 42, got %v", v1)
	}

	v2, ok := generic.GetBoolean(2)
	if ok {
		t.Fatalf("expected not to find boolean configure for 2")
	}

	if v2 != false {
		t.Fatalf("expected boolean configure for 2 to be false, got %v", v2)
	}

	v3, ok := generic.GetUint(2)
	if ok {
		t.Fatalf("expected not to find uint configure for 2")
	}

	if v3 != 0 {
		t.Fatalf("expected uint configure for 2 to be 0, got %v", v3)
	}
}

func TestGenericConfigureUint(t *testing.T) {
	generic := GenericConfigure{
		1: true,
		2: int64(42),
		3: uint64(42),
	}

	v1, ok := generic.GetUint(3)
	if !ok {
		t.Fatalf("expected to find uint configure for 3")
	}

	if v1 != 42 {
		t.Fatalf("expected uint configure for 3 to be 42, got %v", v1)
	}

	v2, ok := generic.GetBoolean(3)
	if ok {
		t.Fatalf("expected not to find boolean configure for 3")
	}

	if v2 != false {
		t.Fatalf("expected boolean configure for 3 to be false, got %v", v2)
	}

	v3, ok := generic.GetInt(3)
	if ok {
		t.Fatalf("expected not to find int configure for 3")
	}

	if v3 != 0 {
		t.Fatalf("expected int configure for 3 to be 0, got %v", v3)
	}
}
