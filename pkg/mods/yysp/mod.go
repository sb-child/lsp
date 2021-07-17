package miya

import (
	"fmt"
	"math/rand"
	mods "mods/modio"
	tools "mods/mtools"
	"strings"

	"github.com/gocolly/colly"
)

const (
	MOD_NAME = "yysp"
)

type Mod struct {
	_成功函数   func(string)
	_信息函数   func(string)
	_警告函数   func(string)
	_报错函数   func(string)
	_爬虫报错函数 func(*colly.Response, error)
	_爬虫请求函数 func(*colly.Request)
	_爬虫回应函数 func(*colly.Response)
	获取到的网址  string
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
	m._爬虫请求函数 = func(r *colly.Request) {
		m.i_信息(fmt.Sprintf("访问[%s]...", r.URL.String()))
	}
	m._爬虫回应函数 = func(r *colly.Response) {
		m.s_成功(fmt.Sprintf("回应[%d]", r.StatusCode))
	}
	网址列表 := make([]string, 0)
	for a := 1; a < 10; a++ {
		网址列表 = append(网址列表, fmt.Sprintf("https://yyspzy%d.xyz", a))
	}
	rand.Shuffle(len(网址列表), func(i, j int) {
		网址列表[i], 网址列表[j] = 网址列表[j], 网址列表[i]
	})
	for _, v := range 网址列表 {
		r := m.t_网站测试(v)
		if r != "" {
			m.获取到的网址 = r
			break
		}
	}
	if m.获取到的网址 == "" {
		m.e_报错("获取网址失败")
		return false
	}
	return true
}
func (m *Mod) AddTag() {
	
}
func (m *Mod) ResetTags() {
	
}
func (m *Mod) GetAllTags() {
	
}
func (m *Mod) GetVideos() {
	
}
func (m *Mod) GetVideoLink() {
	
}
func (m *Mod) t_网站测试(链接 string) string {
	结果 := ""
	爬虫 := tools.CollyCollector()
	爬虫.OnError(m._爬虫报错函数)
	爬虫.OnRequest(m._爬虫请求函数)
	爬虫.OnResponse(m._爬虫回应函数)
	爬虫.OnHTML("meta[http-equiv=\"refresh\"]", func(e *colly.HTMLElement) {
		// 跳转
		主页 := e.Attr("content")[8:]
		主页 = strings.Replace(主页, "http://", "https://", 1)
		m.w_警告(fmt.Sprintf("跳转至[%s]", 主页))
		爬虫.Visit(主页)
	})
	爬虫.OnHTML("meta[name=\"renderer\"]", func(e *colly.HTMLElement) {
		// 主页
		结果 = e.Request.URL.String()
	})
	爬虫.Visit(链接)
	爬虫.Wait()
	return 结果
}
func (m *Mod) g_获取分类列表() {
	爬虫 := tools.CollyCollector()
	爬虫.OnError(m._爬虫报错函数)
	爬虫.OnRequest(m._爬虫请求函数)
	爬虫.OnResponse(m._爬虫回应函数)
	
}
func (m *Mod) OnSucc(f func(s string)) {
	m._成功函数 = f
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
func (m *Mod) s_成功(s string) {
	m._成功函数(s)
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
