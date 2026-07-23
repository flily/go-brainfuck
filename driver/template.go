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

func UnpackValues[T any](items []ContextItem[T]) []T {
	values := make([]T, len(items))
	for i, item := range items {
		values[i] = item.Value
	}

	return values
}

type TestCase struct {
	Name   ContextItem[string]
	Input  []ContextItem[int64]
	Output []ContextItem[int64]
	Memory []ContextItem[int64]
}

func NewTestCase(name string, ctx *context.Context) TestCase {
	c := TestCase{
		Name:   NewContextItem(name, ctx),
		Input:  make([]ContextItem[int64], 0),
		Output: make([]ContextItem[int64], 0),
		Memory: make([]ContextItem[int64], 0),
	}

	return c
}

func (c *TestCase) Equal(o TestCase) bool {
	if c.Name.Value != o.Name.Value {
		return false
	}

	if !slices.Equal(UnpackValues(c.Input), UnpackValues(o.Input)) {
		return false
	}

	if !slices.Equal(UnpackValues(c.Output), UnpackValues(o.Output)) {
		return false
	}

	if !slices.Equal(UnpackValues(c.Memory), UnpackValues(o.Memory)) {
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
	ScriptName ContextItem[string]
	Init       InitParameters
	Tests      []TestCase
}

func (i *TestDriverItem) Equal(o *TestDriverItem) bool {
	if i.ScriptName.Value != o.ScriptName.Value {
		return false
	}

	if !i.Init.Equal(o.Init) {
		return false
	}

	if len(i.Tests) != len(o.Tests) {
		return false
	}

	for idx, test := range i.Tests {
		if !test.Equal(o.Tests[idx]) {
			return false
		}
	}

	return true
}
