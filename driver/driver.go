package driver

import (
	"fmt"

	"github.com/fatih/color"

	"github.com/flily/go-brainfuck/config"
	"github.com/flily/go-brainfuck/context"
	"github.com/flily/go-brainfuck/infra"
	"github.com/flily/go-brainfuck/iofmt"
	"github.com/flily/go-brainfuck/parser"
	"github.com/flily/go-brainfuck/vm"
)

func loadScriptWith[T MemoryUnit](file *context.FileContext) (*infra.CodeMap, error) {
	parser := parser.NewParser(file)
	codemap, err := parser.Parse()
	if err != nil {
		return nil, err
	}

	return codemap, nil
}

func initVM[T MemoryUnit](item *TestDriverItem, codemap *infra.CodeMap) *vm.VM[T] {
	bfvm := vm.New[T](int(item.Init.Value.MemorySize.Value), int(item.Init.Value.StackSize.Value))
	bfvm.LoadHandlers(vm.GetStandardInstructionSetHandlers[T]())
	bfvm.LoadCode(codemap)

	return bfvm
}

func convertList[T MemoryUnit, U MemoryUnit](origin []T) []U {
	h := make([]U, len(origin))
	for i, v := range origin {
		h[i] = U(v)
	}

	return h
}

func caseCheckInput[T MemoryUnit](kase *ContextItem[TestCase], input *iofmt.BufferedReader[T]) error {
	if !input.Empty() {
		offset := input.Offset
		remaining := kase.Value.Input.Value[offset]
		err := remaining.Context.Error("input data not read by program").
			With("data not read by program")
		return err
	}

	return nil
}

func caseCheckOutput[T MemoryUnit](kase *ContextItem[TestCase], got []T) error {
	expected := convertList[int64, T](UnpackValues(kase.Value.Output.Value))
	if len(got) != len(expected) {
		var err error
		if len(got) < len(expected) {
			ctx := kase.Value.Output.Value[len(got)].Context
			err = ctx.Error("output data less than expected").
				With("output %d data items, expect %d", len(got), len(expected))
		} else {
			ctx := kase.Value.Output.Value[len(expected)-1].Context
			err = ctx.Error("output data more than expected").
				With("output %d data items, expect %d", len(got), len(expected))
		}

		return err
	}

	for i, exp := range expected {
		if got[i] != exp {
			ctx := kase.Value.Output.Value[i].Context
			err := ctx.Error("output value mismatch").
				With("got %v", got[i])
			return err
		}
	}

	return nil
}

func caseCheckMemory[T MemoryUnit](bfvm *vm.VM[T], kase *ContextItem[TestCase]) error {
	memoryBase := 0
	if kase.Value.MemoryAt.Valid() {
		memoryBase = int(kase.Value.MemoryAt.Value)
		if memoryBase < 0 || memoryBase >= bfvm.MemorySize {
			err := kase.Value.MemoryAt.Context.Error("memory base out of range").
				With("vm memory size is set to %d", bfvm.MemorySize)
			return err
		}
	}

	memoryLength := len(kase.Value.Memory.Value)
	if memoryBase+memoryLength > bfvm.MemorySize {
		i := bfvm.MemorySize - memoryBase
		value := kase.Value.Memory.Value[i]
		err := value.Context.Error("memory assertion out of range").
			With("vm memory size is set to %d", bfvm.MemorySize)
		return err
	}

	for i, expected := range kase.Value.Memory.Value {
		got := bfvm.Memory[memoryBase+i]
		if T(expected.Value) != got {
			ctx := expected.Context
			err := ctx.Error("memory value mismatch").
				With("got %v", got)
			return err
		}
	}

	return nil
}

func RunCase[T MemoryUnit](item *TestDriverItem, codemap *infra.CodeMap, kase *ContextItem[TestCase]) error {
	bfvm := initVM[T](item, codemap)

	inputData := UnpackValues(kase.Value.Input.Value)
	input := iofmt.NewBufferedReader(convertList[int64, T](inputData))
	output := iofmt.NewBufferedWriter[T](0)

	bfvm.SetInput(input)
	bfvm.SetOutput(output)

	fmt.Printf("case %s on script '%s'\n",
		color.CyanString(kase.Value.Name.Value),
		color.YellowString(item.ScriptName.Value))
	err := bfvm.Run()
	if err != nil {
		fmt.Printf("  - run %s\n", color.RedString("failed"))
		return err
	}
	fmt.Printf("  - run %s\n", color.GreenString("success"))

	if err := caseCheckInput(kase, input); err != nil {
		return err
	}

	outputGot := output.Dump()
	if err := caseCheckOutput(kase, outputGot); err != nil {
		return err
	}

	if err := caseCheckMemory(bfvm, kase); err != nil {
		return err
	}

	return nil
}

func runWithScript[T MemoryUnit](item *TestDriverItem, script *context.FileContext) error {
	codemap, err := loadScriptWith[T](script)
	if err != nil {
		return err
	}

	for _, test := range item.Tests {
		err := RunCase[T](item, codemap, &test)
		if err != nil {
			return err
		}
	}

	return nil
}

func Run[T MemoryUnit](item *TestDriverItem) error {
	scriptFilename := item.ScriptName.Value
	scriptCtx, err := context.ReadFile(scriptFilename)
	if err != nil {
		err = item.ScriptName.Context.Error("read file failed").
			With("%s", err)
		return err
	}

	return runWithScript[T](item, scriptCtx)
}

func genericRunWith(item *TestDriverItem, script *context.FileContext) error {
	var err error
	ut := item.Init.Value.WordType.Value
	switch ut {
	case config.MemoryUnitTypeUint8:
		err = runWithScript[uint8](item, script)
	case config.MemoryUnitTypeUint16:
		err = runWithScript[uint16](item, script)
	case config.MemoryUnitTypeUint32:
		err = runWithScript[uint32](item, script)
	case config.MemoryUnitTypeUint64:
		err = runWithScript[uint64](item, script)
	case config.MemoryUnitTypeInt8:
		err = runWithScript[int8](item, script)
	case config.MemoryUnitTypeInt16:
		err = runWithScript[int16](item, script)
	case config.MemoryUnitTypeInt32:
		err = runWithScript[int32](item, script)
	case config.MemoryUnitTypeInt64:
		err = runWithScript[int64](item, script)
	default:
		ctx := item.Init.Value.WordType.Context
		if ctx == nil {
			// word-type is not set
			kwCtx := item.Init.Context
			err = kwCtx.Error("missing required field").
				With("no word-type specified")
		} else {
			err = item.Init.Value.WordType.Context.
				Error("invalid memory unit type").
				With("invalid memory unit type '%s'", ut.String())
		}
	}

	return err
}

func GenericRun(item *TestDriverItem) error {
	scriptFilename := item.ScriptName.Value
	scriptCtx, err := context.ReadFile(scriptFilename)
	if err != nil {
		err = item.ScriptName.Context.Error("read file failed").
			With("%s", err)
		return err
	}

	return genericRunWith(item, scriptCtx)
}
