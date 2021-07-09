package miya

import (
	mods "mods/modio"
	tools "mods/mtools"

	"github.com/gojek/heimdall/v7/httpclient"
)

const (
	MOD_NAME = "yysp"
)

type Mod struct {
	hc *httpclient.Client
}

func (m *Mod) ModDesc() string {
	return "夜夜视频(夜夜視頻)网站视频获取模块"
}
func (m *Mod) ModName() string {
	return MOD_NAME
}
func (m *Mod) Init() bool {
	m.hc = tools.NewMyHttpClient()
	return true
}
func (m *Mod) OnInfo(f func(...interface{})) {

}
func (m *Mod) OnWarn(f func(...interface{})) {

}
func (m *Mod) OnError(f func(...interface{})) {

}

func init() {
	var m mods.ModuleIO = &Mod{}
	mods.AddModule(MOD_NAME, &m)
}
