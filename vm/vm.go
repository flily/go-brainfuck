package vm

import (
	"errors"
	"io"
	"slices"

	"github.com/flily/go-brainfuck/config"
	"github.com/flily/go-brainfuck/infra"
)



type (
	MemoryUnit     = infra.MemoryUnit
	Instruction    = infra.Instruction
	InstructionSet = infra.InstructionSet
	Code           = infra.Code
	CodeMap        = infra.CodeMap
)

/*
 * Memory offset
 * +-----+-----+-----+-----+-----+-----+-----+
 * |  0  |  1  |  2  |  3  |  4  |  5  |  6  |  offset to memory slices
 * +-----+-----+-----+-----+-----+-----+-----+
 * |  -3 |  -2 |  -1 |  0  |  1  |  2  |  3  |  offset to current
 * +-----+-----+-----+-----+-----+-----+-----+
 *                      ^
 *                      DP
 *                    offset=3
 */
type Snapshot[T MemoryUnit] struct {
	Memory       []T
	MemoryOffset int
	Stack        []int
	Code         []Code
	CodeOffset   int
	IP           int
	DP           int
	SP           int
}

type InstructionHandler[T MemoryUnit] func(vm *VM[T], conf config.ConfigureContainer) error

type InstructionHandlerEntry[T MemoryUnit] struct {
	Instruction Instruction
	Handler     InstructionHandler[T]
}

type VM[T MemoryUnit] struct {
	Memory      []T
	Code        *CodeMap
	IPStack     []int
	Input       Reader[T]
	Output      Writer[T]
	Configure   config.ConfigureContainer
	MemorySize  int
	StackSize   int
	IP          int // Instruction Pointer
	DP          int // Data Pointer
	SP          int // Stack Pointer
	handlers    map[Instruction]InstructionHandler[T]
	currentCode *Code
}

func New[T MemoryUnit](memorySize int, stackSize int) *VM[T] {
	vm := &VM[T]{
		Memory:     make([]T, memorySize),
		Code:       nil,
		IPStack:    make([]int, stackSize),
		Input:      nil,
		Output:     nil,
		Configure:  config.NewGenericConfigure(),
		handlers:   make(map[Instruction]InstructionHandler[T]),
		MemorySize: memorySize,
		StackSize:  stackSize,
		IP:         0,
		DP:         0,
		SP:         0,
	}

	return vm
}

func (m *VM[T]) LoadHandlers(handlers []InstructionHandlerEntry[T]) *VM[T] {
	for _, entry := range handlers {
		ins := entry.Instruction
		m.handlers[ins] = entry.Handler
	}

	return m
}

func (m *VM[T]) LoadCode(code *CodeMap) *VM[T] {
	m.Code = code
	return m
}

func (m *VM[T]) LoadData(data []T) *VM[T] {
	for i := 0; i < min(len(data), len(m.Memory)); i++ {
		m.Memory[i] = data[i]
	}
	return m
}

func (m *VM[T]) Reset() {
	m.DP = 0
	m.SP = 0
	m.IP = 0
}

func (m *VM[T]) SetInput(input Reader[T]) *VM[T] {
	m.Input = input
	return m
}

func (m *VM[T]) SetOutput(output Writer[T]) *VM[T] {
	m.Output = output
	return m
}

func (m *VM[T]) fetchCode(ip int) *Code {
	if ip >= len(m.Code.Codes) {
		return nil
	}

	current := &(m.Code.Codes[ip])
	m.currentCode = current
	return current
}

func (m *VM[T]) GetCurrentCode() *Code {
	return m.currentCode
}

func (m *VM[T]) ExecuteInstruction(code *Code, conf config.ConfigureContainer) error {
	handler, ok := m.handlers[code.Instruction]
	if !ok {
		err := ReasonUnsupportedInstruction.
			OnFatal(code.Context, "").
			With("instruction=%c (%x)", code.Instruction, code.Instruction)
		return err
	}

	return handler(m, conf)
}

func (m *VM[T]) PushIP() error {
	value := m.IP
	if m.SP >= m.StackSize {
		code := m.GetCurrentCode()
		err := ReasonCallStackOverflow.
			OnFatal(code.Context, "call stack overflow, max stack size: %d", m.StackSize).
			With("SP=%d", m.SP)
		return err
	}

	m.IPStack[m.SP] = value
	m.SP += 1
	return nil
}

func (m *VM[T]) PopIP() (int, error) {
	if m.SP <= 0 {
		code := m.GetCurrentCode()
		err := ReasonCallStackEmpty.
			OnFatal(code.Context, "call stack is empty").
			With("SP=%d", m.SP)
		return -1, err
	}

	m.SP -= 1
	n := m.IPStack[m.SP]
	return n, nil
}

func (m *VM[T]) UseIP() error {
	if m.SP <= 0 {
		code := m.GetCurrentCode()
		err := ReasonCallStackEmpty.
			OnFatal(code.Context, "call stack is empty").
			With("SP=%d", m.SP)
		return err
	}

	m.IP = m.IPStack[m.SP-1]
	return nil
}

func (m *VM[T]) Read() (T, error) {
	if m.Input == nil {
		current := m.GetCurrentCode()
		err := ReasonNoInputDevice.
			OnFatal(current.Context, "no input device specified").
			With("read from input failed")
		return 0, err
	}

	value, err := m.Input.Read()
	if err != nil {
		current := m.GetCurrentCode()
		var e error

		if errors.Is(err, io.EOF) {
			e = ReasonReadEOF.
				OnFatal(current.Context, "read EOF").
				With("no more data to read")

		} else {
			e = ReasonReadError.
				OnFatal(current.Context, "read error: %s", err).
				With("read from input failed")
		}

		return 0, e
	}

	return value, nil
}

func (m *VM[T]) Write(value T) error {
	if m.Output == nil {
		current := m.GetCurrentCode()
		err := ReasonNoOutputDevice.
			OnFatal(current.Context, "no output device specified").
			With("write to output failed")
		return err
	}

	err := m.Output.Write(value)
	if err != nil {
		current := m.GetCurrentCode()
		e := ReasonWriteError.
			OnFatal(current.Context, "write error: %s", err).
			With("write to output failed")
		return e
	}

	return nil
}

func (m *VM[T]) Snapshot(dataBefore int, dataAfter int, codeBefore int, codeAfter int) *Snapshot[T] {
	snapshot := &Snapshot[T]{
		Memory:       nil,
		MemoryOffset: dataBefore,
		Stack:        nil,
		Code:         nil,
		CodeOffset:   codeBefore,
		IP:           m.IP,
		DP:           m.DP,
		SP:           m.SP,
	}

	dataStart := max(m.DP-dataBefore, 0)
	dataEnd := min(m.DP+dataAfter+1, m.MemorySize)
	snapshot.Memory = slices.Clone(m.Memory[dataStart:dataEnd])

	snapshot.Stack = slices.Clone(m.IPStack[:m.SP])
	snapshot.Code = m.Code.Snapshot(m.IP, codeBefore, codeAfter)

	return snapshot
}

func (m *VM[T]) Step() error {
	ins := m.fetchCode(m.IP)
	if ins == nil {
		return ReasonHalt
	}

	err := m.ExecuteInstruction(ins, m.Configure)
	m.IP += 1
	return err
}

func (m *VM[T]) Run() error {
	var err error
	for {
		err = m.Step()
		if err != nil {
			break
		}
	}

	if errors.Is(err, ReasonHalt) {
		return nil
	}

	return err
}
