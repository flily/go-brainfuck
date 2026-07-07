package vm

type Configure uint32

const (
	StandardConfigureMask Configure = 0xff000000
	ConfigureStandard     Configure = iota
	ConfigureReadOnEOF
	ConfigureReadIgnoreOnEOF
)

func (c Configure) IsStandard() bool {
	return (c & StandardConfigureMask) == 0
}

type ConfigureContainer interface {
	GetBooleanConfigure(Configure) bool
	GetIntConfigure(Configure) int64
	GetUintConfigure(Configure) uint64
}
