package mtools

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"

	openssl "github.com/Luzifer/go-openssl/v4"
)

// regex
var (
	domainMatch,
	urlLinkMatch,
	m3u8UrlLinkMatch,
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
func TagLinkMatch() *regexp.Regexp {
	return tagLinkMatch
}
func init() {
	domainMatch, _ = regexp.Compile("(http[s]?://.*?)/")
	urlLinkMatch, _ = regexp.Compile("\"url\":\"(.*?)\"")
	m3u8UrlLinkMatch, _ = regexp.Compile("m3u8url = '(.*?)'")
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

func UrlGetToStr(url string) string {
	hc := http.Client{}
	r, err := hc.Get(url)
	if err != nil {
		return ""
	}
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return ""
	}
	return string(b)
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
func FindVideoSource(old string) (dir string, domain string) {
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
	result := UrlGetToStr(old)
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

type M3U8Decoder struct {
}

func (d *M3U8Decoder) Init() {
}
