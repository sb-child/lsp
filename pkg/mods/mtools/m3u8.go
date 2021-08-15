package mtools

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
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
	ddyunboKeyMatch, _ = regexp.Compile(`CryptoJS.AES.decrypt\(content, '(.*?)'\);`)
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
func _PKCS7UnPadding(pt []byte) []byte {
	length := len(pt)
	unp := int(pt[length-1])
	return pt[:(length - unp)]
}
func _PaddingLeft(ori []byte, pad byte, length int) []byte {
	if len(ori) >= length {
		return ori[:length]
	}
	pads := bytes.Repeat([]byte{pad}, length-len(ori))
	return append(pads, ori...)
}
func _AESDecryptForDdyunbo(s, key string) string {
	out, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		return ""
	}
	// padding for key./
	bkey := _PaddingLeft([]byte(key), '0', 16)
	// decrypt
	akey, err := aes.NewCipher(bkey)
	if err != nil {
		return ""
	}
	decrypter := cipher.NewCBCDecrypter(akey, bkey)
	dec := make([]byte, len(out))
	decrypter.CryptBlocks(dec, out)
	dec = _PKCS7UnPadding(dec)
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
		re2 := M3U8UrlLinkMatch().FindStringSubmatch(result)
		if len(re2) == 0 {
			content := DecryptMethodForDdyunbo(result)
			fmt.Println("[:" + content)
			return
		}
		dir = re2[1]
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
