package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/flily/go-brainfuck/context"
	"github.com/flily/go-brainfuck/iofmt"
	"github.com/flily/go-brainfuck/parser"
	"github.com/flily/go-brainfuck/vm"
)

func GenericMain[T vm.MemoryUnit](conf *Configure, source string) {
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
	bfvm.Input = iofmt.NewReader(os.Stdin, iofmt.NewLittleEndianEncoder[T]())
	bfvm.Output = iofmt.NewWriter(os.Stdout, iofmt.NewLittleEndianEncoder[T]())
	bfvm.Configure = conf
	bfvm.LoadHandlers(vm.GetStandardInstructionSetHandlers[T]())
	bfvm.LoadCode(codemap)

	err = bfvm.Run()
	if err != nil {
		fmt.Printf("error executing code:\n%s\n", err)
	}
}

func makeArgumentParser() (*Configure, *flag.FlagSet) {
	conf := NewConfigure()

	set := flag.NewFlagSet("brainfuck", flag.ExitOnError)
	set.Var(&conf.MemoryUnitType, "word",
		"data type of memory unit cell, one of int[8/16/32/64] or uint[8/16/32/64]")
	set.Int64Var(&conf.ReadValueOnEOF, "eof-value", -1,
		"value returned when meets EOF")
	set.BoolVar(&conf.IgnoreReadEOF, "ignore-eof", false,
		"ignore EOF when reading input")
	set.BoolVar(&conf.RaiseReadEOF, "raise-eof", false,
		"raise error when meets EOF")

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
	case MemoryUnitTypeUint8:
		GenericMain[uint8](conf, filename)
	case MemoryUnitTypeUint16:
		GenericMain[uint16](conf, filename)
	case MemoryUnitTypeUint32:
		GenericMain[uint32](conf, filename)
	case MemoryUnitTypeUint64:
		GenericMain[uint64](conf, filename)
	case MemoryUnitTypeInt8:
		GenericMain[int8](conf, filename)
	case MemoryUnitTypeInt16:
		GenericMain[int16](conf, filename)
	case MemoryUnitTypeInt32:
		GenericMain[int32](conf, filename)
	case MemoryUnitTypeInt64:
		GenericMain[int64](conf, filename)
	default:
		fmt.Printf("invalid memory unit type '%s'\n", conf.MemoryUnitType.String())
	}
}
