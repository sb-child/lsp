package mtools

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"

	"github.com/wumansgy/goEncrypt"
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
	ddyunboKeyMatch *regexp.Regexp
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
	ddyunboContentMatch, _ = regexp.Compile(`var content = "(.*?)";`)
	ddyunboKeyMatch, _ = regexp.Compile(`var bytes =  CryptoJS.AES.decrypt(content, '(.*?)');`)
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
func _AESDecryptForDdyunbo(s, key string) string {
	// aes.NewCipher(key)
	out, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		return ""
	}
	dec, err := goEncrypt.AesCbcDecrypt(out, []byte(key), []byte(""))
	if err != nil {
		return ""
	}
	return string(dec)
}
func DecryptMethodForDdyunbo(html string) (link string) {
	content := ddyunboContentMatch.FindStringSubmatch(html)[1]
	key := ddyunboKeyMatch.FindStringSubmatch(html)[1]
	r := _AESDecryptForDdyunbo(content, key)
	return r
}
func FindVideoSource(old string) (dir string, domain string) {
	// https://xxx.xx/xxx/xxx.m3u8
	domain = DomainMatch().FindStringSubmatch(old)[1]
	dir = old
	if strings.HasSuffix(old, ".m3u8") {
		return
	}
	// https://xxx.xx/xxx/xxx
	result := UrlGetToStr(old)
	re1 := UrlLinkMatch().FindStringSubmatch(result)
	if len(re1) == 0 {
		// m3u8url = 'https://xxx.xx/xxx/xxx.m3u8'
		fmt.Println(old)
		fmt.Println(result)
		dir = M3U8UrlLinkMatch().FindStringSubmatch(result)[1]
		return
	}
	// "url":"https://xxx.xx/xxx/xxx.m3u8"
	dir = domain + re1[1]
	return
}

type M3U8Decoder struct {
}

func (d *M3U8Decoder) Init() {
}
