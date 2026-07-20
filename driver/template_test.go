package driver

import (
	"testing"

	"github.com/flily/go-brainfuck/config"
)

func TestContextItemUnpackValues(t *testing.T) {
	items := []ContextItem[int]{
		NewContextItem(1, nil),
		NewContextItem(2, nil),
		NewContextItem(3, nil),
	}

	values := UnpackValues(items)
	expected := []int{1, 2, 3}

	if len(values) != len(expected) {
		t.Fatalf("expected length %d, got %d", len(expected), len(values))
	}
}

func TestTestCaseEqual(t *testing.T) {
	expected := NewTestCase("example", nil)
	expected.Input = []ContextItem[int64]{
		NewContextItem[int64](1, nil),
		NewContextItem[int64](2, nil),
	}
	expected.Output = []ContextItem[int64]{
		NewContextItem[int64](11, nil),
	}
	expected.Memory = []ContextItem[int64]{
		NewContextItem[int64](21, nil),
		NewContextItem[int64](22, nil),
	}

	o := NewTestCase("lorem ipsum", nil)
	if expected.Equal(o) {
		t.Fatalf("expected not equal, got equal")
	}

	o.Name.Value = "example"
	if expected.Equal(o) {
		t.Fatalf("expected not equal, got equal")
	}

	o.Input = []ContextItem[int64]{
		NewContextItem[int64](1, nil),
		NewContextItem[int64](2, nil),
	}
	if expected.Equal(o) {
		t.Fatalf("expected not equal, got equal")
	}

	o.Output = []ContextItem[int64]{
		NewContextItem[int64](11, nil),
	}
	if expected.Equal(o) {
		t.Fatalf("expected not equal, got equal")
	}

	o.Memory = []ContextItem[int64]{
		NewContextItem[int64](21, nil),
		NewContextItem[int64](22, nil),
	}
	if !expected.Equal(o) {
		t.Fatalf("expected equal, got not equal")
	}
}

func TestInitParametersEqual(t *testing.T) {
	expected := InitParameters{
		MemorySize: NewContextItem[uint64](1024, nil),
		StackSize:  NewContextItem[uint64](512, nil),
		WordType:   NewContextItem(config.MemoryUnitTypeUint8, nil),
		EOFValue:   NewContextItem[int64](-1, nil),
		IgnoreEOF:  NewContextItem(true, nil),
		RaiseEOF:   NewContextItem(true, nil),
	}

	o := InitParameters{}
	if expected.Equal(o) {
		t.Fatalf("expected not equal, got equal")
	}

	o.MemorySize = NewContextItem[uint64](1024, nil)
	if expected.Equal(o) {
		t.Fatalf("expected not equal, got equal")
	}

	o.StackSize = NewContextItem[uint64](512, nil)
	if expected.Equal(o) {
		t.Fatalf("expected not equal, got equal")
	}

	o.WordType = NewContextItem(config.MemoryUnitTypeUint8, nil)
	if expected.Equal(o) {
		t.Fatalf("expected not equal, got equal")
	}

	o.EOFValue = NewContextItem[int64](-1, nil)
	if expected.Equal(o) {
		t.Fatalf("expected not equal, got equal")
	}

	o.IgnoreEOF = NewContextItem(true, nil)
	if expected.Equal(o) {
		t.Fatalf("expected not equal, got equal")
	}

	o.RaiseEOF = NewContextItem(true, nil)

	if !expected.Equal(o) {
		t.Fatalf("expected equal, got not equal")
	}
}

func TestTestDriverItemEqual(t *testing.T) {
	expected := &TestDriverItem{
		ScriptName: NewContextItem("script", nil),
		Init: InitParameters{
			MemorySize: NewContextItem[uint64](1024, nil),
			StackSize:  NewContextItem[uint64](512, nil),
			WordType:   NewContextItem(config.MemoryUnitTypeUint8, nil),
			EOFValue:   NewContextItem[int64](-1, nil),
			IgnoreEOF:  NewContextItem(true, nil),
			RaiseEOF:   NewContextItem(true, nil),
		},
		Tests: []TestCase{
			NewTestCase("test1", nil),
			NewTestCase("test2", nil),
		},
	}

	o := &TestDriverItem{}
	if expected.Equal(o) {
		t.Fatalf("expected not equal, got equal")
	}

	o.ScriptName = NewContextItem("script", nil)
	if expected.Equal(o) {
		t.Fatalf("expected not equal, got equal")
	}

	o.Init = InitParameters{
		MemorySize: NewContextItem[uint64](1024, nil),
		StackSize:  NewContextItem[uint64](512, nil),
		WordType:   NewContextItem(config.MemoryUnitTypeUint8, nil),
		EOFValue:   NewContextItem[int64](-1, nil),
		IgnoreEOF:  NewContextItem(true, nil),
		RaiseEOF:   NewContextItem(true, nil),
	}
	if expected.Equal(o) {
		t.Fatalf("expected not equal, got equal")
	}

	o.Tests = make([]TestCase, 0)
	if expected.Equal(o) {
		t.Fatalf("expected not equal, got equal")
	}

	o.Tests = []TestCase{
		NewTestCase("test2", nil),
		NewTestCase("test1", nil),
	}
	if expected.Equal(o) {
		t.Fatalf("expected not equal, got equal")
	}

	o.Tests = []TestCase{
		NewTestCase("test1", nil),
		NewTestCase("test2", nil),
	}
	if !expected.Equal(o) {
		t.Fatalf("expected equal, got not equal")
	}
}
