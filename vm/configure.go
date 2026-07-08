package vm

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

type GenericConfigure map[Configure]any

func NewGenericConfigure() GenericConfigure {
	return make(GenericConfigure)
}

func genericGetConfigure[T any](g GenericConfigure, c Configure) (T, bool) {
	found := false
	var result T
	if v, ok := g[c]; ok {
		if t, ok := v.(T); ok {
			result = t
			found = true
		}
	}

	return result, found
}

func (c GenericConfigure) GetBoolean(conf Configure) (bool, bool) {
	return genericGetConfigure[bool](c, conf)
}

func (c GenericConfigure) GetInt(conf Configure) (int64, bool) {
	return genericGetConfigure[int64](c, conf)
}

func (c GenericConfigure) GetUint(conf Configure) (uint64, bool) {
	return genericGetConfigure[uint64](c, conf)
}
