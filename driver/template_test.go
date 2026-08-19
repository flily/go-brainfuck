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
	expected.Value.Input = NewContextItem([]ContextItem[int64]{
		NewContextItem[int64](1, nil),
		NewContextItem[int64](2, nil),
	}, nil)
	expected.Value.Output = NewContextItem([]ContextItem[int64]{
		NewContextItem[int64](11, nil),
	}, nil)
	expected.Value.Memory = NewContextItem([]ContextItem[int64]{
		NewContextItem[int64](21, nil),
		NewContextItem[int64](22, nil),
	}, nil)
	expected.Value.MemoryAt = NewContextItem[uint64](1024, nil)
	expected.Value.Init = NewContextItem([]ContextItem[int64]{
		NewContextItem[int64](31, nil),
		NewContextItem[int64](32, nil),
	}, nil)
	expected.Value.InitAt = NewContextItem[uint64](2048, nil)

	o := NewTestCase("lorem ipsum", nil)
	if expected.Value.Equal(o.Value) {
		t.Fatalf("expected not equal, got equal")
	}

	o.Value.Name.Value = "example"
	if expected.Value.Equal(o.Value) {
		t.Fatalf("expected not equal, got equal")
	}

	o.Value.Input = NewContextItem([]ContextItem[int64]{
		NewContextItem[int64](1, nil),
		NewContextItem[int64](2, nil),
	}, nil)
	if expected.Value.Equal(o.Value) {
		t.Fatalf("expected not equal, got equal")
	}

	o.Value.Output = NewContextItem([]ContextItem[int64]{
		NewContextItem[int64](11, nil),
	}, nil)
	if expected.Value.Equal(o.Value) {
		t.Fatalf("expected not equal, got equal")
	}

	o.Value.Init = NewContextItem([]ContextItem[int64]{
		NewContextItem[int64](31, nil),
		NewContextItem[int64](32, nil),
	}, nil)
	if expected.Value.Equal(o.Value) {
		t.Fatalf("expected not equal, got equal")
	}

	o.Value.InitAt = NewContextItem[uint64](2048, nil)
	if expected.Value.Equal(o.Value) {
		t.Fatalf("expected not equal, got equal")
	}

	o.Value.Memory = NewContextItem([]ContextItem[int64]{
		NewContextItem[int64](21, nil),
		NewContextItem[int64](22, nil),
	}, nil)
	if expected.Value.Equal(o.Value) {
		t.Fatalf("expected not equal, got equal")
	}

	o.Value.MemoryAt = NewContextItem[uint64](1024, nil)

	if !expected.Value.Equal(o.Value) {
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
		Init: NewContextItem(InitParameters{
			MemorySize: NewContextItem[uint64](1024, nil),
			StackSize:  NewContextItem[uint64](512, nil),
			WordType:   NewContextItem(config.MemoryUnitTypeUint8, nil),
			EOFValue:   NewContextItem[int64](-1, nil),
			IgnoreEOF:  NewContextItem(true, nil),
			RaiseEOF:   NewContextItem(true, nil),
		}, nil),
		Tests: []ContextItem[TestCase]{
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

	o.Init = NewContextItem(InitParameters{
		MemorySize: NewContextItem[uint64](1024, nil),
		StackSize:  NewContextItem[uint64](512, nil),
		WordType:   NewContextItem(config.MemoryUnitTypeUint8, nil),
		EOFValue:   NewContextItem[int64](-1, nil),
		IgnoreEOF:  NewContextItem(true, nil),
		RaiseEOF:   NewContextItem(true, nil),
	}, nil)
	if expected.Equal(o) {
		t.Fatalf("expected not equal, got equal")
	}

	o.Tests = make([]ContextItem[TestCase], 0)
	if expected.Equal(o) {
		t.Fatalf("expected not equal, got equal")
	}

	o.Tests = []ContextItem[TestCase]{
		NewTestCase("test2", nil),
		NewTestCase("test1", nil),
	}
	if expected.Equal(o) {
		t.Fatalf("expected not equal, got equal")
	}

	o.Tests = []ContextItem[TestCase]{
		NewTestCase("test1", nil),
		NewTestCase("test2", nil),
	}
	if !expected.Equal(o) {
		t.Fatalf("expected equal, got not equal")
	}
}
