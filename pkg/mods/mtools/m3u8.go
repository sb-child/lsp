package mtools

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"path"
	"regexp"
	"strings"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	openssl "github.com/Luzifer/go-openssl/v4"
)

// regex
var (
	domainMatch,
	urlLinkMatch,
	m3u8UrlLinkMatch,
	m3u8ContentInfoMatch,
	m3u8TsInfoMatch,
	tagLinkMatch *regexp.Regexp
)

// private
var (
	ddyunboContentMatch,
	ddyunboKeyMatch,
	urlLink2Match,
	ddyunboPlaylistMatch,
	ddyunboMainMatch *regexp.Regexp
)

// end of private

func DomainMatch() *regexp.Regexp {
	return domainMatch
}
func UrlLinkMatch() *regexp.Regexp {
	return urlLinkMatch
}
func M3U8UrlLinkMatch() *regexp.Regexp {
	return m3u8UrlLinkMatch
}
func M3U8ContentInfoMatch() *regexp.Regexp {
	return m3u8ContentInfoMatch
}
func M3U8TsInfoMatch() *regexp.Regexp {
	return m3u8TsInfoMatch
}
func TagLinkMatch() *regexp.Regexp {
	return tagLinkMatch
}
func init() {
	domainMatch, _ = regexp.Compile("(http[s]?://.*?)/")
	urlLinkMatch, _ = regexp.Compile("\"url\":\"(.*?)\"")
	m3u8UrlLinkMatch, _ = regexp.Compile("m3u8url = '(.*?)'")
	m3u8ContentInfoMatch, _ = regexp.Compile(`(http[s]?://)?(.*?\.m3u8)(\?.*)*`)
	m3u8TsInfoMatch, _ = regexp.Compile(`(http[s]?://)?(.*?\.ts)(\?.*)*`)
	tagLinkMatch, _ = regexp.Compile(`/index\.php/vod/type/id/(.*?).html`)
	// private
	urlLink2Match, _ = regexp.Compile(`"url":"(.*?)","url_next"`)
	ddyunboContentMatch, _ = regexp.Compile(`var content = "(.*?)";`)
	ddyunboKeyMatch, _ = regexp.Compile(`CryptoJS.AES.decrypt\(content, '(.*?)'\);`)
	ddyunboPlaylistMatch, _ = regexp.Compile(`'\[\{"url":"(.*?)"\}\]'`)
	ddyunboMainMatch, _ = regexp.Compile(`var main = "(.*?)"`)
	// end of private
}

// end of regex

func UrlGetToStr(url string) (string, error) {
	hc := http.Client{}
	r, err := hc.Get(url)
	if err != nil {
		return "", err
	}
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return "", err
	}
	return string(b), nil
}
func UrlGetToStrMust(url string) string {
	r, err := UrlGetToStr(url)
	if err != nil {
		return ""
	}
	return r
}
func _AESDecryptForDdyunbo(s, k string) string {
	o := openssl.New()
	dec, err := o.DecryptBytes(k, []byte(s), openssl.BytesToKeyMD5)
	if err != nil {
		return ""
	}
	return string(dec)
}
func DecryptMethodForDdyunbo(html string) (link string) {
	defer func() {
		rec := recover()
		if rec != nil {
			fmt.Println("解密模块内部错误:")
			panic(rec)
		}
	}()
	content := ddyunboContentMatch.FindStringSubmatch(html)[1]
	key := ddyunboKeyMatch.FindStringSubmatch(html)[1]
	r := _AESDecryptForDdyunbo(content, key)
	return r
}
func FindVideoSource(old string) (dir string, domain string, e error) {
	// defer func() {
	// 	fmt.Printf("dir=%s domain=%s\n", dir, domain)
	// }()
	// https://xxx.xx/xxx/xxx.m3u8
	domain = DomainMatch().FindStringSubmatch(old)[1]
	dir = old
	if strings.HasSuffix(old, ".m3u8") {
		return
	}
	// https://xxx.xx/xxx/xxx
	result, err := UrlGetToStr(old)
	if err != nil {
		e = err
		return
	}
	re1 := urlLink2Match.FindStringSubmatch(result)
	if len(re1) == 0 {
		// var main = "http://xxx.xx/xxx/xxx.m3u8?xxx"
		re4 := ddyunboMainMatch.FindStringSubmatch(result)
		if len(re4) != 0 {
			dir = domain + re4[1]
			return
		}
		// m3u8url = 'https://xxx.xx/xxx/xxx.m3u8'
		re2 := ddyunboPlaylistMatch.FindStringSubmatch(result)
		// fallback to ddyunbo
		if len(re2) == 0 {
			content := DecryptMethodForDdyunbo(result)
			re3 := UrlLinkMatch().FindStringSubmatch(content)
			dir = domain + re3[1]
			return
		}
		dir = re2[1]
		return
	}
	// "url":"https://xxx.xx/xxx/xxx.m3u8"
	dir = domain + strings.ReplaceAll(re1[1], "\\", "")
	return
}

type VideoDatabase struct {
	dir string
	db  *gorm.DB
}

type M3U8Video struct {
	gorm.Model
	// modio.VideoContainer
	Link      string // The link to the video
	VideoLink string // The m3u8 of the video
	Title     string // The title of the video
	Img       string // The image of the video
	Desc      string // The description of the video
}

func (vdb *VideoDatabase) Init(dir string) error {
	vdb.dir = path.Join(dir, "_lsp.db")
	db, err := gorm.Open(sqlite.Open(vdb.dir), &gorm.Config{})
	if err != nil {
		fmt.Printf("打不开数据库: %s\n", err)
		return err
	}
	vdb.db = db
	db.AutoMigrate(&M3U8Video{})
	return nil
}

func (vdb *VideoDatabase) Add(video *M3U8Video) error {
	return vdb.db.Create(video).Error
}

type M3U8Decoder struct {
	content string
}

func (d *M3U8Decoder) Init(m3u8url string) error {
	// domain := DomainMatch().FindStringSubmatch(m3u8url)[1]
	m3u8Content, err := UrlGetToStr(m3u8url)
	if err != nil {
		return fmt.Errorf("内容获取失败: %s", err.Error())
	}
	d.content = m3u8Content
	buffer := make([][]string, 0)
	fmt.Println(d)
	for _, line := range strings.Split(d.content, "\n") {
		if strings.HasPrefix(line, "#") {
			continue
		}
		ln := M3U8ContentInfoMatch().FindStringSubmatch(line)
		if len(ln) != 0 {
			ln[0] = "m"
			buffer = append(buffer, ln)
		}
		ln = M3U8TsInfoMatch().FindStringSubmatch(line)
		if len(ln) != 0 {
			ln[0] = "t"
			buffer = append(buffer, ln)
		}
	}
	fmt.Println(buffer)
	return nil
}
