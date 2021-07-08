package mods

type ModuleIO interface {
	ModName() string
	ModDesc() string
	Init() bool
	OnInfo(func(...interface{}))
	OnWarn(func(...interface{}))
	OnError(func(...interface{}))
}
