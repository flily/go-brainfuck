package main

import (
	"fmt"

	"github.com/flily/go-brainfuck/config"
)

const (
	DefaultMemorySize = 4 * 1024
	DefaultStackSize  = 128
)

type MemoryUnitType int

const (
	MemoryUnitTypeInvalid MemoryUnitType = iota
	MemoryUnitTypeUint8
	MemoryUnitTypeUint16
	MemoryUnitTypeUint32
	MemoryUnitTypeUint64
	MemoryUnitTypeInt8
	MemoryUnitTypeInt16
	MemoryUnitTypeInt32
	MemoryUnitTypeInt64
)

var memoryUnitTypeText = map[MemoryUnitType]string{
	MemoryUnitTypeUint8:  "uint8",
	MemoryUnitTypeUint16: "uint16",
	MemoryUnitTypeUint32: "uint32",
	MemoryUnitTypeUint64: "uint64",
	MemoryUnitTypeInt8:   "int8",
	MemoryUnitTypeInt16:  "int16",
	MemoryUnitTypeInt32:  "int32",
	MemoryUnitTypeInt64:  "int64",
}

func (t *MemoryUnitType) String() string {
	if text, ok := memoryUnitTypeText[*t]; ok {
		return text
	}

	return "unknown"
}

func (t *MemoryUnitType) Set(value string) error {
	for k, v := range memoryUnitTypeText {
		if v == value {
			*t = k
			return nil
		}
	}

	return fmt.Errorf("invalid memory unit type '%s'", value)
}

type Configure struct {
	MemoryUnitType MemoryUnitType
	MemorySize     uint64
	StackSize      uint64
	ReadValueOnEOF int64
	IgnoreReadEOF  bool
	RaiseReadEOF   bool
}

func NewConfigure() *Configure {
	conf := &Configure{
		MemoryUnitType: MemoryUnitTypeUint8,
		MemorySize:     DefaultMemorySize,
		StackSize:      DefaultStackSize,
		ReadValueOnEOF: -1,
		IgnoreReadEOF:  false,
		RaiseReadEOF:   false,
	}

	return conf
}

func (c *Configure) GetBoolean(conf config.Configure) (bool, bool) {
	result := false
	found := true
	switch conf {
	case config.ConfigureReadValueIgnoreOnEOF:
		result = c.IgnoreReadEOF

	case config.ConfigureReadEOFRaiseError:
		result = c.RaiseReadEOF

	default:
		found = false
	}

	return result, found
}

func (c *Configure) GetInt(conf config.Configure) (int64, bool) {
	result := int64(0)
	found := true
	switch conf {
	case config.ConfigureReadValueOnEOF:
		result = c.ReadValueOnEOF

	default:
		found = false
	}

	return result, found
}

func (c *Configure) GetUint(conf config.Configure) (uint64, bool) {
	result := uint64(0)
	found := true
	// switch conf {
	// default:
	// 	found = false
	// }

	return result, found
}
