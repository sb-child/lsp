package miya

import (
	"fmt"
	"math/rand"
	mods "mods/modio"
	tools "mods/mtools"
	"strings"
	"sync"

	"github.com/gocolly/colly"
)

const (
	MOD_NAME = "yysp"
)

type Mod struct {
	s_成功    func(string)
	i_信息    func(string)
	w_警告    func(string)
	e_报错    func(string)
	_爬虫报错函数 func(*colly.Response, error)
	_爬虫请求函数 func(*colly.Request)
	_爬虫回应函数 func(*colly.Response)
	主站      string
}

func (m *Mod) ModDesc() string {
	return "夜夜视频(夜夜視頻)网站视频获取模块"
}
func (m *Mod) ModName() string {
	return MOD_NAME
}
func (m *Mod) makeSpider() *colly.Collector {
	m._爬虫报错函数 = func(r *colly.Response, err error) {
		m.e_报错(fmt.Sprintf("[%s]连接异常[%d]: %s", r.Request.URL.String(), r.StatusCode, err.Error()))
	}
	m._爬虫请求函数 = func(r *colly.Request) {
		m.i_信息(fmt.Sprintf("访问[%s]...", r.URL.String()))
	}
	m._爬虫回应函数 = func(r *colly.Response) {
		m.s_成功(fmt.Sprintf("[%s]回应[%d]", r.Request.URL.String(), r.StatusCode))
	}
	爬虫 := tools.CollyCollector()
	爬虫.OnError(m._爬虫报错函数)
	爬虫.OnRequest(m._爬虫请求函数)
	爬虫.OnResponse(m._爬虫回应函数)
	return 爬虫
}
func (m *Mod) Init() bool {
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
			m.主站 = r
			break
		}
	}
	if m.主站 == "" {
		m.e_报错("获取网址失败")
		return false
	}
	return true
}
func (m *Mod) AddTag(name string) {

}
func (m *Mod) ResetTags() {

}
func (m *Mod) GetAllTags() map[string]string {
	爬虫 := m.makeSpider()
	list := make(map[string]string, 0)
	爬虫.OnHTML(`a[class="1\=0"]`, func(e *colly.HTMLElement) {
		href := strings.TrimSpace(e.Attr("href"))
		// m.i_信息(href)
		f := tools.TagLinkMatch().FindStringSubmatch(href)
		if len(f) == 0 {
			return
		}
		r := strings.TrimSpace(f[1])
		rt := strings.TrimSpace(e.Text)
		list[r] = rt
	})
	爬虫.Visit(m.主站)
	爬虫.Wait()
	return list
}

func (m *Mod) tag2url(t string) string {
	return fmt.Sprintf(m.主站+"/index.php/vod/type/id/%s.html", t)
}

// GetVideos 获取指定分类或默认分类的视频网页链接
func (m *Mod) GetVideos(t []string) []mods.VideoContainer {
	爬虫 := m.makeSpider()
	r := make([]mods.VideoContainer, 0)
	rc := make(chan mods.VideoContainer, 20+10*len(t))
	goLock := sync.WaitGroup{}
	爬虫.OnHTML(`li>a[target="_blank"]`, func(e *colly.HTMLElement) {
		link := m.主站 + strings.TrimSpace(e.Attr("href"))
		title := strings.ReplaceAll(strings.TrimSpace(e.Attr("title")), "\n", " ")
		img := m.主站 + strings.TrimSpace(e.ChildAttr("img", "src"))
		// m.i_信息(fmt.Sprintf("l:%s t:%s i:%s", link, title, img))
		go func() {
			goLock.Add(1)
			rc <- mods.VideoContainer{Link: link, Title: title, Desc: "", Img: img}
		}()
	})
	go func() {
		for {
			select {
			case x, ok := <-rc:
				if !ok {
					continue
				}
				r = append(r, x)
				goLock.Done()
			}
		}
	}()
	// 分类转链接
	if len(t) == 0 {
		t = append(t, m.主站)
	} else {
		nt := make([]string, 0)
		for _, v := range t {
			nt = append(nt, m.tag2url(v))
		}
		t = nt
	}
	// 并行爬取
	for _, v := range t {
		爬虫.Visit(v)
	}
	爬虫.Wait()
	m.i_信息("全部获取完成, 等待汇总任务...")
	goLock.Wait()
	m.s_成功(fmt.Sprintf("汇总完成, 共[%d]个视频", len(r)))
	return r
}
func (m *Mod) GetVideoLink() {

}
func (m *Mod) t_网站测试(链接 string) string {
	结果 := ""
	爬虫 := m.makeSpider()
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
func (m *Mod) OnSucc(f func(s string)) {
	m.s_成功 = f
}
func (m *Mod) OnInfo(f func(s string)) {
	m.i_信息 = f
}
func (m *Mod) OnWarn(f func(s string)) {
	m.w_警告 = f
}
func (m *Mod) OnError(f func(s string)) {
	m.e_报错 = f
}

func init() {
	var m mods.ModuleIO = &Mod{}
	mods.AddModule(MOD_NAME, &m)
}
