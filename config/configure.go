package config

type Configure uint32

const (
	StandardConfigureMask Configure = 0xff000000
	ConfigureStandard     Configure = iota
	ConfigureReadValueOnEOF
	ConfigureReadValueIgnoreOnEOF
	ConfigureReadEOFRaiseError
)

func (c Configure) IsStandard() bool {
	return (c & StandardConfigureMask) == 0
}

type ConfigureContainer interface {
	GetBoolean(Configure) (bool, bool)
	GetInt(Configure) (int64, bool)
	GetUint(Configure) (uint64, bool)
}
