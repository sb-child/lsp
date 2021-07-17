package modio

type ModuleIO interface {
	ModName() string
	ModDesc() string
	Init() bool
	OnSucc(func(string))
	OnInfo(func(string))
	OnWarn(func(string))
	OnError(func(string))
	GetAllTags() map[string]string
}

var mods map[string]*ModuleIO

func init() {
	mods = make(map[string]*ModuleIO)
}

func AddModule(name string, mod *ModuleIO) {
	mods[name] = mod
}

func GetModule(name string) *ModuleIO {
	return mods[name]
}

func GetAllModules() map[string]*ModuleIO {
	return mods
}
