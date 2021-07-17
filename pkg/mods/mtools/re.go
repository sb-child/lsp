package mtools

import "regexp"

var domainMatch *regexp.Regexp
var urlLinkMatch *regexp.Regexp
var m3u8UrlLinkMatch *regexp.Regexp
var tagLinkMatch *regexp.Regexp

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
}
