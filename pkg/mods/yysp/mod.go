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

// --- 模块基本功能

// ModDesc 模块介绍
func (m *Mod) ModDesc() string {
	return "夜夜视频(夜夜視頻)网站视频获取模块"
}

// ModName 模块名
func (m *Mod) ModName() string {
	return MOD_NAME
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

// --- 初始化方法

// Init 初始化
func (m *Mod) Init() bool {
	网址列表 := make([]string, 0)
	for a := 1; a < 3; a++ {
		网址列表 = append(网址列表, fmt.Sprintf("https://yyspzy%d.xyz", a))
	}
	for a := 1; a < 3; a++ {
		网址列表 = append(网址列表, fmt.Sprintf("https://yylu%d.com", a))
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

// --- 私有方法

// t_网站测试 返回一个直达的主页链接
func (m *Mod) t_网站测试(链接 string) string {
	结果 := ""
	爬虫 := m.makeSpider(true)
	爬虫.OnHTML("meta[http-equiv=\"refresh\"]", func(e *colly.HTMLElement) {
		// 跳转
		主页 := e.Attr("content")[8:]
		主页 = strings.Replace(主页, "http://", "https://", 1)
		m.w_警告(fmt.Sprintf("@[%s]", 主页))
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

// makeSpider 爬虫构造器
func (m *Mod) makeSpider(async bool) *colly.Collector {
	m._爬虫报错函数 = func(r *colly.Response, err error) {
		m.e_报错(fmt.Sprintf("*bot [%s]![%d]: %s", r.Request.URL.String(), r.StatusCode, err.Error()))
	}
	m._爬虫请求函数 = func(r *colly.Request) {
		m.i_信息(fmt.Sprintf("*bot >[%s]", r.URL.String()))
	}
	m._爬虫回应函数 = func(r *colly.Response) {
		m.s_成功(fmt.Sprintf("*bot [%s]>[%d]", r.Request.URL.String(), r.StatusCode))
	}
	var 爬虫 *colly.Collector
	if async {
		爬虫 = tools.CollyCollector()
	} else {
		爬虫 = tools.CollyCollectorSlow()
	}
	爬虫.OnError(m._爬虫报错函数)
	爬虫.OnRequest(m._爬虫请求函数)
	爬虫.OnResponse(m._爬虫回应函数)
	return 爬虫
}

// tag2url 分类转链接
func (m *Mod) tag2url(t string) string {
	return fmt.Sprintf(m.主站+"/index.php/vod/type/id/%s.html", t)
}

// getVideoM3U8 获取链接的m3u8地址
func (m *Mod) getVideoM3U8(links []string) (r map[string]string) {
	爬虫 := m.makeSpider(false)
	urlMap := sync.Map{}
	爬虫.OnHTML("script", func(e *colly.HTMLElement) {
		if strings.Count(e.Text, "encrypt") == 0 {
			return
		}
		finds := tools.UrlLinkMatch().FindStringSubmatch(e.Text)
		if len(finds) == 0 {
			urlMap.Store(e.Request.URL.String(), "")
			return
		}
		m3u8Url := finds[1]
		m3u8Url = strings.ReplaceAll(m3u8Url, "\\", "")
		m3u8Url, domain := tools.FindVideoSource(m3u8Url)
		if !strings.HasPrefix(m3u8Url, "https://") {
			m3u8Url = domain + "/" + m3u8Url
		}
		urlMap.Store(e.Request.URL.String(), m3u8Url)
	})
	for _, v := range links {
		爬虫.Visit(v)
	}
	爬虫.Wait()
	r = make(map[string]string, len(links))
	for _, v := range links {
		u, ok := urlMap.Load(v)
		if !ok {
			continue
		}
		r[v] = u.(string)
	}
	return
}

// --- 公共方法
// GetAllTags 获取所有分类
func (m *Mod) GetAllTags() map[string]string {
	爬虫 := m.makeSpider(true)
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

// GetVideos 获取指定分类或默认分类的视频网页链接
func (m *Mod) GetVideos(t []string) []mods.VideoContainer {
	爬虫 := m.makeSpider(true)
	r := make([]mods.VideoContainer, 0)
	rc := make(chan mods.VideoContainer, 20+10*len(t))
	goLock := sync.WaitGroup{}
	processTitle := func(ot string) (r string) {
		ot = strings.TrimSpace(ot)
		ot = strings.ReplaceAll(ot, "\n", " ")
		for {
			title := strings.ReplaceAll(
				ot,
				"  ",
				" ",
			)
			if ot == title {
				r = title
				break
			}
			ot = title
		}
		return
	}
	// 预处理
	processLink := func(ln string) string {
		ln = strings.ReplaceAll(ln, "detail", "play")
		ln = strings.ReplaceAll(ln, ".html", "/sid/1/nid/1.html")
		return ln
	}
	爬虫.OnHTML(`li>a[target="_blank"]`, func(e *colly.HTMLElement) {
		goLock.Add(1)
		link := m.主站 + strings.TrimSpace(processLink(e.Attr("href")))
		title := processTitle(e.Attr("title"))
		img := strings.TrimSpace(e.ChildAttr("img", "src"))
		rc <- mods.VideoContainer{Link: link, Title: title, Desc: "", Img: img}
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
		// 默认使用推荐视频
		t = append(t, m.主站)
	} else {
		// 使用自定义分类
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
	m.i_信息("获取视频M3U8...")
	links := make([]string, len(r))
	for _, v := range r {
		links = append(links, v.Link)
	}
	linksM3U8 := m.getVideoM3U8(links)
	for k, v := range linksM3U8 {
		for k2 := range r {
			if r[k2].Link != k {
				continue
			}
			r[k2].VideoLink = v
		}
	}
	m.s_成功(fmt.Sprintf("汇总完成, 共[%d]个视频", len(r)))
	return r
}

func init() {
	var m mods.ModuleIO = &Mod{}
	mods.AddModule(MOD_NAME, &m)
}
