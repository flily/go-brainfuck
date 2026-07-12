package main

import (
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/flily/go-brainfuck/config"
	"github.com/flily/go-brainfuck/context"
	"github.com/flily/go-brainfuck/iofmt"
	"github.com/flily/go-brainfuck/parser"
	"github.com/flily/go-brainfuck/vm"
)

const (
	DefaultMemorySize = 4 * 1024
	DefaultStackSize  = 128
)

func openReadFile(filename string, defaultReader io.Reader) (io.Reader, error) {
	if filename == "" {
		return defaultReader, nil
	}

	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	return file, nil
}

func openWriteFile(filename string, defaultWriter io.Writer) (io.Writer, error) {
	if filename == "" {
		return defaultWriter, nil
	}

	file, err := os.Create(filename)
	if err != nil {
		return nil, err
	}

	return file, nil
}

func GenericMain[T vm.MemoryUnit](conf *config.Container, source string) {
	fctx, err := context.ReadFile(source)
	if err != nil {
		fmt.Printf("error reading file: %v\n", err)
		return
	}

	parser := parser.NewParser(fctx)
	codemap, err := parser.Parse()
	if err != nil {
		fmt.Printf("error parsing file:\n%s\n", err)
		return
	}

	memorySize, stackSize := int(conf.MemorySize), int(conf.StackSize)
	if conf.MemorySize > uint64(memorySize) {
		fmt.Printf("warning: memory size too large, using %d\n", memorySize)
	}

	if conf.StackSize > uint64(stackSize) {
		fmt.Printf("warning: stack size too large, using %d\n", stackSize)
	}

	bfvm := vm.New[T](memorySize, stackSize)
	input, err := openReadFile(conf.InputFilename, os.Stdin)
	if err != nil {
		fmt.Printf("open input file '%s' error: %v\n", conf.InputFilename, err)
		return
	}

	output, err := openWriteFile(conf.OutputFilename, os.Stdout)
	if err != nil {
		fmt.Printf("open output file '%s' error: %v\n", conf.OutputFilename, err)
		return
	}

	bfvm.Input = iofmt.NewReader(input, iofmt.NewLittleEndianEncoder[T]())
	bfvm.Output = iofmt.NewWriter(output, iofmt.NewLittleEndianEncoder[T]())
	bfvm.Configure = conf
	bfvm.LoadHandlers(vm.GetStandardInstructionSetHandlers[T]())
	bfvm.LoadCode(codemap)

	err = bfvm.Run()
	if err != nil {
		fmt.Printf("error executing code:\n%s\n", err)
	}
}

func makeArgumentParser() (*config.Container, *flag.FlagSet) {
	conf := config.NewContainer(DefaultMemorySize, DefaultStackSize)

	set := flag.NewFlagSet("brainfuck", flag.ExitOnError)
	set.Var(&conf.MemoryUnitType, "word",
		"data type of memory unit cell, one of int[8/16/32/64] or uint[8/16/32/64]")

	set.Int64Var(&conf.ReadValueOnEOF, "eof-value", -1,
		"value returned when meets EOF")
	set.BoolVar(&conf.IgnoreReadEOF, "ignore-eof", false,
		"ignore EOF when reading input")
	set.BoolVar(&conf.RaiseReadEOF, "raise-eof", false,
		"raise error when meets EOF")
	set.StringVar(&conf.InputFilename, "in", "",
		"filename of input file, default to stdin if not specified")
	set.StringVar(&conf.OutputFilename, "out", "",
		"filename of output file, default to stdout if not specified")
	set.Var(&conf.Endian, "endian",
		"endian type for input/output, use 'be' for big-endian or 'le' for little-endian")

	return conf, set
}

func main() {
	rawArgs := os.Args[1:]
	conf, set := makeArgumentParser()
	_ = set.Parse(rawArgs)
	args := set.Args()

	if len(args) <= 0 || len(args) > 1 {
		fmt.Printf("brainfuck [ARGS...] source")
		set.Usage()
		return
	}

	filename := args[0]

	switch conf.MemoryUnitType {
	case config.MemoryUnitTypeUint8:
		GenericMain[uint8](conf, filename)
	case config.MemoryUnitTypeUint16:
		GenericMain[uint16](conf, filename)
	case config.MemoryUnitTypeUint32:
		GenericMain[uint32](conf, filename)
	case config.MemoryUnitTypeUint64:
		GenericMain[uint64](conf, filename)
	case config.MemoryUnitTypeInt8:
		GenericMain[int8](conf, filename)
	case config.MemoryUnitTypeInt16:
		GenericMain[int16](conf, filename)
	case config.MemoryUnitTypeInt32:
		GenericMain[int32](conf, filename)
	case config.MemoryUnitTypeInt64:
		GenericMain[int64](conf, filename)
	default:
		fmt.Printf("invalid memory unit type '%s'\n", conf.MemoryUnitType.String())
	}
}
