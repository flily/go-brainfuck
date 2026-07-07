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
	GetBooleanConfigure(Configure) (bool, bool)
	GetIntConfigure(Configure) (int64, bool)
	GetUintConfigure(Configure) (uint64, bool)
}

type GenericConfigure map[Configure]any

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

func (c GenericConfigure) GetBooleanConfigure(conf Configure) (bool, bool) {
	return genericGetConfigure[bool](c, conf)
}

func (c GenericConfigure) GetIntConfigure(conf Configure) (int64, bool) {
	return genericGetConfigure[int64](c, conf)
}

func (c GenericConfigure) GetUintConfigure(conf Configure) (uint64, bool) {
	return genericGetConfigure[uint64](c, conf)
}
