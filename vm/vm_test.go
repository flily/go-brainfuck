package vm

import (
	"slices"
	"strings"
	"testing"

	"github.com/flily/go-brainfuck/context"
)

type testVMCase[T MemoryUnit] struct {
	code               string
	data               []T
	input              []T
	memorySize         int
	stackSize          int
	expectedData       []T
	expectedDataOffset int
	expectedOutput     []T
}

const (
	minMemorySize = 32
	minStackSize  = 32
)

func (c testVMCase[T]) Run(t *testing.T) {
	t.Helper()

	file := context.ReadFileString("test.bf", c.code)
	parser := NewParser(file)
	codemap, err := parser.Parse()
	if err != nil {
		t.Fatalf("error parsing code:\n%s", err)
	}

	memorySize := max(c.memorySize, minMemorySize)
	stackSize := max(c.stackSize, minStackSize)
	vm := New[T](memorySize, stackSize)
	vm.LoadCode(codemap)
	vm.LoadData(c.data)

	input := NewBufferedReader(c.input)
	vm.SetInput(input)
	output := NewBufferedWriter[T](0)
	vm.SetOutput(output)

	vm.LoadHandlers(GetStandardInstructionSetHandlers[T]())

	err = vm.Run()
	if err != nil {
		t.Fatalf("error running code:\n%s", err)
	}

	start := c.expectedDataOffset
	memorySlices := vm.Memory[start : start+len(c.expectedData)]
	if !slices.Equal(memorySlices, c.expectedData) {
		t.Errorf("expected: %v", c.expectedData)
		t.Errorf("got:      %v", memorySlices)
		t.Fatalf("memory mismatch")
	}

	if !slices.Equal(output.Data, c.expectedOutput) {
		t.Errorf("expected: %v", c.expectedOutput)
		t.Errorf("got:      %v", output.Data)
		t.Fatalf("output mismatch")
	}
}

func TestVMSimpleCodeAddSub(t *testing.T) {
	testVMCase[uint8]{
		code:         "+++--",
		expectedData: []uint8{1},
	}.Run(t)
}

func TestVMSimpleIO(t *testing.T) {
	testVMCase[uint8]{
		code:           ",+.",
		input:          []uint8{41},
		expectedOutput: []uint8{42},
	}.Run(t)
}

func TestVMSimpleLoop(t *testing.T) {
	testVMCase[uint8]{
		code:         "[-]",
		data:         []uint8{42},
		expectedData: []uint8{0},
	}.Run(t)
}

func TestVMCleanMemory(t *testing.T) {
	testVMCase[uint8]{
		code:         "[[-]>]",
		data:         []uint8{1, 2, 3, 4, 5, 6},
		expectedData: []uint8{0, 0, 0, 0, 0, 0},
	}.Run(t)
}

func TestVMHelloWorld(t *testing.T) {
	testVMCase[uint8]{
		code: strings.Join([]string{
			"++++++++++[>+++",
			"++++>+++++++++ ",
			"+>+++>+<<<<-]>+",
			"+.>+.+++++++..+",
			"++.>++.<<++++++",
			"+++++++++.>.++ ",
			"+.------.------",
			"--.>+.>.",
		}, ""),
		expectedOutput: []uint8("Hello World!\n"),
	}.Run(t)
}
