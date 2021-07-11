package miya

import (
	mods "mods/modio"
	// tools "mods/mtools"

	"github.com/gocolly/colly"
	"github.com/gojek/heimdall/v7/httpclient"
)

const (
	MOD_NAME = "miya"
)

type Mod struct {
	hc *httpclient.Client
}

func (m *Mod) ModDesc() string {
	return "蜜芽网站视频获取模块"
}
func (m *Mod) ModName() string {
	return MOD_NAME
}
func (m *Mod) Init() bool {
	colly.NewCollector()
	// m.hc = tools.NewMyHttpClient()
	return true
}
func (m *Mod) OnInfo(f func(s string)) {

}
func (m *Mod) OnWarn(f func(s string)) {

}
func (m *Mod) OnError(f func(s string)) {

}

func init() {
	var m mods.ModuleIO = &Mod{}
	mods.AddModule(MOD_NAME, &m)
}
