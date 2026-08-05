package main

import (
	"github.com/flily/go-brainfuck/config"
	"github.com/flily/go-brainfuck/context"
	"github.com/flily/go-brainfuck/driver"
	"github.com/flily/go-brainfuck/parser"
	"github.com/flily/go-brainfuck/vm"
)

func initDriver(filename string) (*driver.TestDriverItem, error) {
	fctx, err := context.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	parser := driver.NewParser(filename, fctx)
	driverItem, err := parser.Parse()
	if err != nil {
		return nil, err
	}

	return driverItem, nil
}

func genericRunCase[T driver.MemoryUnit](item *driver.TestDriverItem) error {
	scriptFilename := item.ScriptName.Value
	scriptCtx, err := context.ReadFile(scriptFilename)
	if err != nil {
		err = item.ScriptName.Context.Error("read file failed").
			With("%s", err)
		return err
	}

	parser := parser.NewParser(scriptCtx)
	codemap, err := parser.Parse()
	if err != nil {
		return err
	}

	bfvm := vm.New[T](int(item.Init.Value.MemorySize.Value), int(item.Init.Value.StackSize.Value))
	bfvm.LoadHandlers(vm.GetStandardInstructionSetHandlers[T]())
	bfvm.LoadCode(codemap)

	err = bfvm.Run()
	if err != nil {
		return err
	}

	return nil
}

func runCase(filename string) error {
	driverItem, err := initDriver(filename)
	if err != nil {
		return err
	}

	unitType := driverItem.Init.Value.WordType.Value
	switch unitType {
	case config.MemoryUnitTypeUint8:
		err = genericRunCase[uint8](driverItem)
	case config.MemoryUnitTypeInt8:
		err = genericRunCase[int8](driverItem)
	case config.MemoryUnitTypeUint16:
		err = genericRunCase[uint16](driverItem)
	case config.MemoryUnitTypeInt16:
		err = genericRunCase[int16](driverItem)
	case config.MemoryUnitTypeUint32:
		err = genericRunCase[uint32](driverItem)
	case config.MemoryUnitTypeInt32:
		err = genericRunCase[int32](driverItem)
	case config.MemoryUnitTypeUint64:
		err = genericRunCase[uint64](driverItem)
	case config.MemoryUnitTypeInt64:
		err = genericRunCase[int64](driverItem)

	default:
		ctx := driverItem.Init.Value.WordType.Context
		if ctx == nil {
			// word-type is not set
			kwCtx := driverItem.Init.Context
			err = kwCtx.Error("missing required field").
				With("no word-type specified")
		} else {
			err = driverItem.Init.Value.WordType.Context.
				Error("invalid memory unit type").
				With("invalid memory unit type '%s'", unitType.String())
		}
	}

	return err
}
