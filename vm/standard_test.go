package vm

import (
	"testing"
	"testing/iotest"

	"errors"
	"slices"
	"strings"

	"github.com/flily/go-brainfuck/config"
	"github.com/flily/go-brainfuck/context"
	"github.com/flily/go-brainfuck/iofmt"
	"github.com/flily/go-brainfuck/parser"
)

type nilWrite[T MemoryUnit] struct{}

func (w *nilWrite[T]) Write(p T) (err error) {
	return errors.New("write error")
}

func newNilWriter[T MemoryUnit]() iofmt.Writer[T] {
	return &nilWrite[T]{}
}

type testVMCase[T MemoryUnit] struct {
	code               string
	data               []T
	input              iofmt.Reader[T]
	output             iofmt.Writer[T]
	configure          ConfigureContainer
	memorySize         int
	stackSize          int
	expectedData       []T
	expectedDataOffset int
	expectedOutput     []T
	expectedError      string
}

const (
	minMemorySize = 32
	minStackSize  = 32
)

func parseCodeLiteral(code string, filename string) (*CodeMap, error) {
	file := context.ReadFileString(filename, code)
	parser := parser.NewParser(file)
	codemap, err := parser.Parse()
	if err != nil {
		return nil, err
	}

	return codemap, nil
}

func (c testVMCase[T]) Run(t *testing.T) {
	t.Helper()

	codemap, err := parseCodeLiteral(c.code, "test.bf")
	if err != nil {
		t.Fatalf("error parsing code:\n%s", err)
	}

	memorySize := max(c.memorySize, minMemorySize)
	stackSize := max(c.stackSize, minStackSize)
	vm := New[T](memorySize, stackSize)
	vm.LoadCode(codemap)
	vm.LoadData(c.data)
	if c.configure != nil {
		vm.Configure = c.configure
	}

	vm.SetInput(c.input)

	output := c.output
	if output == nil {
		if c.output == nil {
			output = iofmt.NewBufferedWriter[T](0)
		}
	} else if _, ok := c.output.(*nilWrite[T]); ok {
		output = nil
	}
	vm.SetOutput(output)

	vm.LoadHandlers(GetStandardInstructionSetHandlers[T]())

	err = vm.Run()
	if err != nil {
		if len(c.expectedError) > 0 {
			if err.Error() != c.expectedError {
				t.Fatalf("expected error:\n%s\n got:\n%s", c.expectedError, err)
			}
		} else {
			t.Fatalf("error running code:\n%s", err)
		}
	} else {
		if len(c.expectedError) > 0 {
			t.Fatalf("expected error but got no error")
		}
	}

	start := c.expectedDataOffset
	memorySlices := vm.Memory[start : start+len(c.expectedData)]
	if !slices.Equal(memorySlices, c.expectedData) {
		t.Errorf("expected: %v", c.expectedData)
		t.Errorf("got:      %v", memorySlices)
		t.Fatalf("memory mismatch")
	}

	dumpable, ok := output.(iofmt.DumpableWriter[T])
	if ok {
		outputData := dumpable.Dump()
		if !slices.Equal(outputData, c.expectedOutput) {
			t.Errorf("expected: %v", c.expectedOutput)
			t.Errorf("got:      %v", outputData)
			t.Fatalf("output mismatch")
		}
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
		input:          iofmt.NewBufferedInput[uint8](41),
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
		code: "[[-]>]",
		data: []uint8{
			1, 2, 3, 4, 5, 6,
		},
		expectedData: []uint8{
			0, 0, 0, 0, 0, 0,
		},
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

func TestVMReadEOFNoConfigure(t *testing.T) {
	testVMCase[uint8]{
		code:  ",,,,",
		input: iofmt.NewBufferedInput[uint8](1, 2),
		configure: config.GenericConfigure{
			config.ConfigureReadEOFRaiseError: true,
		},
		expectedData: []uint8{2},
		expectedError: strings.Join([]string{
			"test.bf:1:3: fatal: read EOF",
			"    1 | ,,,,",
			"      |   ^",
			"      |   no more data to read",
		}, "\n"),
	}.Run(t)
}

func TestVMReadEOFIgnore(t *testing.T) {
	testVMCase[uint8]{
		code:  ",,,,",
		input: iofmt.NewBufferedInput[uint8](1, 2),
		configure: config.GenericConfigure{
			config.ConfigureReadValueIgnoreOnEOF: true,
		},
		expectedData: []uint8{2},
	}.Run(t)
}

func TestVMReadEOFAsZero(t *testing.T) {
	testVMCase[uint8]{
		code:  ",,,,",
		input: iofmt.NewBufferedInput[uint8](1, 2),
		configure: config.GenericConfigure{
			config.ConfigureReadValueOnEOF: 0,
		},
		expectedData: []uint8{0},
	}.Run(t)
}

func TestVMReadEOFAsMinusOne(t *testing.T) {
	testVMCase[uint8]{
		code:  ",>,,,",
		input: iofmt.NewBufferedInput[uint8](1, 2),
		configure: config.GenericConfigure{
			config.ConfigureReadValueOnEOF: int64(-1),
		},
		expectedData: []uint8{
			1, 0xff, 0, 0,
		},
	}.Run(t)
}

func TestVMReadError(t *testing.T) {
	err := errors.New("read error")

	testVMCase[uint8]{
		code:  "++++,.,.",
		input: iofmt.NewReader(iotest.ErrReader(err), iofmt.NewEncoderLE[uint8]()),
		expectedError: strings.Join([]string{
			"test.bf:1:5: fatal: read error: read error",
			"    1 | ++++,.,.",
			"      |     ^",
			"      |     read from input failed",
		}, "\n"),
		expectedData: []uint8{
			4, 0, 0, 0,
		},
	}.Run(t)
}

func TestVMReadOnNilDevice(t *testing.T) {
	testVMCase[uint8]{
		code: "++++,.,.",
		expectedError: strings.Join([]string{
			"test.bf:1:5: fatal: no input device specified",
			"    1 | ++++,.,.",
			"      |     ^",
			"      |     read from input failed",
		}, "\n"),
		expectedData: []uint8{
			4, 0, 0, 0,
		},
	}.Run(t)
}

func TestVMWriteError(t *testing.T) {
	err := errors.New("write error")

	testVMCase[uint8]{
		code:   "++++..",
		output: iofmt.NewWriter(iofmt.NewErrWriter(err), iofmt.NewEncoderLE[uint8]()),
		expectedError: strings.Join([]string{
			"test.bf:1:5: fatal: write error: write error",
			"    1 | ++++..",
			"      |     ^",
			"      |     write to output failed",
		}, "\n"),
		expectedData: []uint8{
			4, 0, 0, 0,
		},
		expectedOutput: []uint8{},
	}.Run(t)
}

func TestVMWriteOnNilDevice(t *testing.T) {
	testVMCase[uint8]{
		code:   "++++..",
		output: newNilWriter[uint8](),
		expectedError: strings.Join([]string{
			"test.bf:1:5: fatal: no output device specified",
			"    1 | ++++..",
			"      |     ^",
			"      |     write to output failed",
		}, "\n"),

		expectedData: []uint8{
			4, 0, 0, 0,
		},
		expectedOutput: []uint8{},
	}.Run(t)
}
