package mtools

import "github.com/gocolly/colly"

func CollyCollector() *colly.Collector {
	return colly.NewCollector(
		colly.UserAgent(
			"Mozilla/5.0 (X11; Linux x86_64; rv:89.0) Gecko/20100101 Firefox/89.0",
		),
		colly.AllowURLRevisit(),
		colly.Async(true),
	)
}
func CollyCollectorSlow() *colly.Collector {
	return colly.NewCollector(
		colly.UserAgent(
			"Mozilla/5.0 (X11; Linux x86_64; rv:89.0) Gecko/20100101 Firefox/89.0",
		),
		colly.AllowURLRevisit(),
		colly.Async(false),
	)
}

// func NewMyHttpClient() *httpclient.Client {
// 	return httpclient.NewClient(
// 		httpclient.WithHTTPTimeout(2*time.Second),
// 		httpclient.WithRetryCount(5),
// 	)
// }

// func HttpGet(cli *httpclient.Client, url string) (*http.Response, error) {
// 	h := http.Header{}
// 	h.Add("Accept-Language", "zh-CN,zh;q=0.8,zh-TW;q=0.7,zh-HK;q=0.5,en-US;q=0.3,en;q=0.2")
// 	h.Add("Cache-Control", "max-age=0")
// 	h.Add("User-Agent", "Mozilla/5.0 (X11; Linux x86_64; rv:89.0) Gecko/20100101 Firefox/89.0")
// 	return cli.Get(url, h)
// }

// func HttpString(r *http.Response) string {
// 	buf, err := ioutil.ReadAll(r.Body)
// 	if err != nil {
// 		return ""
// 	}
// 	return string(buf)
// }

// func Resp2HTML(r *http.Response) *html.Node {
// 	n, err := htmlquery.Parse(r.Body)
// 	if err != nil {
// 		return nil
// 	}
// 	return n
// }

func init() {
}
