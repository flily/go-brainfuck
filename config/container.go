package config

import (
	"fmt"
	"strings"
)

type (
	MemoryUnitType int
	Endian         int
)

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

	EndianLittle Endian = 0
	EndianBig    Endian = 1
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
	lower := strings.ToLower(value)
	for k, v := range memoryUnitTypeText {
		if v == lower {
			*t = k
			return nil
		}
	}

	return fmt.Errorf("invalid memory unit type '%s'", value)
}

func (e *Endian) String() string {
	switch *e {
	case EndianBig:
		return "big-endian"
	case EndianLittle:
		return "little-endian"
	default:
		return "unknown"
	}
}

func (e *Endian) Set(value string) error {
	v := strings.ToLower(value)
	switch v {
	case "big", "be", "big-endian", "bigendian", "b":
		*e = EndianBig
	case "little", "le", "little-endian", "littleendian", "l":
		*e = EndianLittle
	default:
		return fmt.Errorf("invalid endian type '%s'", value)
	}

	return nil
}

type Container struct {
	MemoryUnitType MemoryUnitType
	MemorySize     uint64
	StackSize      uint64
	Endian         Endian
	InputFilename  string
	OutputFilename string
	ReadValueOnEOF int64
	IgnoreReadEOF  bool
	RaiseReadEOF   bool
}

func NewContainer(memorySize uint64, stackSize uint64) *Container {
	conf := &Container{
		MemoryUnitType: MemoryUnitTypeUint8,
		MemorySize:     memorySize,
		StackSize:      stackSize,
		Endian:         EndianLittle,
		ReadValueOnEOF: -1,
		IgnoreReadEOF:  false,
		RaiseReadEOF:   false,
	}

	return conf
}

func (c *Container) GetBoolean(conf Configure) (bool, bool) {
	result := false
	found := true
	switch conf {
	case ConfigureReadValueIgnoreOnEOF:
		result = c.IgnoreReadEOF

	case ConfigureReadEOFRaiseError:
		result = c.RaiseReadEOF

	default:
		found = false
	}

	return result, found
}

func (c *Container) GetInt(conf Configure) (int64, bool) {
	result := int64(0)
	found := true
	switch conf {
	case ConfigureReadValueOnEOF:
		result = c.ReadValueOnEOF

	default:
		found = false
	}

	return result, found
}

func (c *Container) GetUint(conf Configure) (uint64, bool) {
	result := uint64(0)
	found := true
	// switch conf {
	// default:
	// 	found = false
	// }

	return result, found
}
