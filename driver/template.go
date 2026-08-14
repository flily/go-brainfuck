package driver

import (
	"slices"

	"github.com/flily/go-brainfuck/config"
	"github.com/flily/go-brainfuck/context"
	"github.com/flily/go-brainfuck/infra"
)

type (
	MemoryUnit = infra.MemoryUnit
)

const (
	SectionScript   = "script"
	SectionInit     = "init"
	SectionCase     = "case"
	FieldMemorySize = "memory-size"
	FieldStackSize  = "stack-size"
	FieldWord       = "word"
	FieldEOFValue   = "eof-value"
	FieldIgnoreEOF  = "ignore-eof"
	FieldRaiseEOF   = "raise-eof"
	FieldInput      = "input"
	FieldOutput     = "output"
	FieldMemory     = "memory"
	KeywordAt       = "at"
)

var requiredSections = []string{
	SectionScript,
	SectionInit,
	SectionCase,
}

type ContextItem[T any] struct {
	Value   T
	Context *context.Context
}

func NewContextItem[T any](content T, ctx *context.Context) ContextItem[T] {
	item := ContextItem[T]{
		Value:   content,
		Context: ctx,
	}

	return item
}

func (i *ContextItem[T]) Valid() bool {
	return i.Context != nil
}

func UnpackValues[T any](items []ContextItem[T]) []T {
	values := make([]T, len(items))
	for i, item := range items {
		values[i] = item.Value
	}

	return values
}

type TestCase struct {
	Name     ContextItem[string]
	Input    ContextItem[[]ContextItem[int64]]
	Output   ContextItem[[]ContextItem[int64]]
	Memory   ContextItem[[]ContextItem[int64]]
	MemoryAt ContextItem[uint64]
}

func NewTestCase(name string, ctx *context.Context) ContextItem[TestCase] {
	c := TestCase{
		Name:     NewContextItem(name, ctx),
		Input:    NewContextItem(make([]ContextItem[int64], 0), ctx),
		Output:   NewContextItem(make([]ContextItem[int64], 0), ctx),
		Memory:   NewContextItem(make([]ContextItem[int64], 0), ctx),
		MemoryAt: NewContextItem[uint64](0, ctx),
	}

	return NewContextItem(c, ctx)
}

func (c *TestCase) Equal(o TestCase) bool {
	if c.Name.Value != o.Name.Value {
		return false
	}

	if !slices.Equal(UnpackValues(c.Input.Value), UnpackValues(o.Input.Value)) {
		return false
	}

	if !slices.Equal(UnpackValues(c.Output.Value), UnpackValues(o.Output.Value)) {
		return false
	}

	if !slices.Equal(UnpackValues(c.Memory.Value), UnpackValues(o.Memory.Value)) {
		return false
	}

	if c.MemoryAt.Value != o.MemoryAt.Value {
		return false
	}

	return true
}

type InitParameters struct {
	MemorySize ContextItem[uint64]
	StackSize  ContextItem[uint64]
	WordType   ContextItem[config.MemoryUnitType]
	EOFValue   ContextItem[int64]
	IgnoreEOF  ContextItem[bool]
	RaiseEOF   ContextItem[bool]
}

func (p *InitParameters) Equal(o InitParameters) bool {
	if p.MemorySize.Value != o.MemorySize.Value {
		return false
	}

	if p.StackSize.Value != o.StackSize.Value {
		return false
	}

	if p.WordType.Value != o.WordType.Value {
		return false
	}

	if p.EOFValue.Value != o.EOFValue.Value {
		return false
	}

	if p.IgnoreEOF.Value != o.IgnoreEOF.Value {
		return false
	}

	if p.RaiseEOF.Value != o.RaiseEOF.Value {
		return false
	}

	return true
}

type TestDriverItem struct {
	Filename   string
	ScriptName ContextItem[string]
	Init       ContextItem[InitParameters]
	Tests      []ContextItem[TestCase]
}

func (i *TestDriverItem) Equal(o *TestDriverItem) bool {
	if i.ScriptName.Value != o.ScriptName.Value {
		return false
	}

	if !i.Init.Value.Equal(o.Init.Value) {
		return false
	}

	if len(i.Tests) != len(o.Tests) {
		return false
	}

	for idx, test := range i.Tests {
		if !test.Value.Equal(o.Tests[idx].Value) {
			return false
		}
	}

	return true
}
