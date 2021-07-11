package miya

import (
	"fmt"
	mods "mods/modio"
	tools "mods/mtools"

	"github.com/gocolly/colly"
)

const (
	MOD_NAME = "yysp"
)

type Mod struct {
	_信息函数   func(string)
	_警告函数   func(string)
	_报错函数   func(string)
	_爬虫报错函数 func(*colly.Response, error)
}

func (m *Mod) ModDesc() string {
	return "夜夜视频(夜夜視頻)网站视频获取模块"
}
func (m *Mod) ModName() string {
	return MOD_NAME
}
func (m *Mod) Init() bool {
	m._爬虫报错函数 = func(r *colly.Response, err error) {
		m.e_报错(fmt.Sprintf("连接异常[%d]: %s", r.StatusCode, err.Error()))
	}
	r := m.t_网站测试("https://www.sbc-io.xyz:81")
	m.i_信息(r)
	r = m.t_网站测试("https://yyspzy1.xyz")
	m.i_信息(r)
	return true
}
func (m *Mod) t_网站测试(链接 string) string {
	结果 := ""
	爬虫 := tools.CollyCollector()
	m.i_信息(fmt.Sprintf("开始连接[%s]", 链接))
	爬虫.OnError(m._爬虫报错函数)
	爬虫.OnResponse(func(r *colly.Response) {
		m.i_信息(fmt.Sprintf("回应[%d]", r.StatusCode))
		// 结果 = string(r.Body)
	})

	爬虫.Visit(链接)
	爬虫.Wait()
	return 结果
}
func (m *Mod) OnInfo(f func(s string)) {
	m._信息函数 = f
}
func (m *Mod) OnWarn(f func(s string)) {
	m._警告函数 = f
}
func (m *Mod) OnError(f func(s string)) {
	m._报错函数 = f
}
func (m *Mod) i_信息(s string) {
	m._信息函数(s)
}
func (m *Mod) w_警告(s string) {
	m._警告函数(s)
}
func (m *Mod) e_报错(s string) {
	m._报错函数(s)
}

func init() {
	var m mods.ModuleIO = &Mod{}
	mods.AddModule(MOD_NAME, &m)
}
