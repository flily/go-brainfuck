package utils

import (
	"testing"
)

func TestStack(t *testing.T) {
	items := []int{11, 13, 17, 19}
	expected := []int{19, 17, 13, 11}

	stack := NewStackWithCapacity[int](10)
	if !stack.IsEmpty() {
		t.Fatalf("expected stack to be empty")
	}

	if size := stack.Size(); size != 0 {
		t.Fatalf("expected stack size to be 0, got %d", size)
	}

	for _, item := range items {
		stack.Push(&item)
	}

	if size := stack.Size(); size != len(items) {
		t.Fatalf("expected stack size to be %d, got %d", len(items), size)
	}

	if stack.IsEmpty() {
		t.Fatalf("expected stack to not be empty")
	}

	for _, exp := range expected {
		item, ok := stack.Pop()
		if !ok {
			t.Fatalf("expected to pop item %d, but stack was empty", exp)
		}

		if *item != exp {
			t.Fatalf("expected popped item to be %d, got %d", exp, *item)
		}
	}

	if size := stack.Size(); size != 0 {
		t.Fatalf("expected stack size to be 0 after popping all items, got %d", size)
	}

	if !stack.IsEmpty() {
		t.Fatalf("expected stack to be empty after popping all items")
	}
}

func TestStackPushManyItems(t *testing.T) {
	stack := NewStackWithCapacity[int](2)
	for i := 0; i < 100; i++ {
		stack.Push(&i)
	}

	if size := stack.Size(); size != 100 {
		t.Fatalf("expected stack size to be 100, got %d", size)
	}

	for i := 99; i >= 0; i-- {
		item, ok := stack.Pop()
		if !ok {
			t.Fatalf("expected to pop item %d, but stack was empty", i)
		}

		if *item != i {
			t.Fatalf("expected popped item to be %d, got %d", i, *item)
		}
	}

	if size := stack.Size(); size != 0 {
		t.Fatalf("expected stack size to be 0 after popping all items, got %d", size)
	}

	if !stack.IsEmpty() {
		t.Fatalf("expected stack to be empty after popping all items")
	}
}
