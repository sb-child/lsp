package mtools

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	openssl "github.com/Luzifer/go-openssl/v4"
	"github.com/gocolly/colly"
)

// regex
var (
	domainMatch,
	urlLinkMatch,
	m3u8UrlLinkMatch,
	m3u8ContentInfoMatch,
	m3u8TsInfoMatch,
	urlDirMatch,
	m3u8KeyUrlMatch,
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
func M3U8KeyUrlMatch() *regexp.Regexp {
	return m3u8KeyUrlMatch
}
func UrlDirMatch() *regexp.Regexp {
	return urlDirMatch
}
func TagLinkMatch() *regexp.Regexp {
	return tagLinkMatch
}
func init() {
	domainMatch, _ = regexp.Compile("(http[s]?://.*?)/")
	urlDirMatch, _ = regexp.Compile(`(http[s]?://.*/)`)
	urlLinkMatch, _ = regexp.Compile("\"url\":\"(.*?)\"")
	m3u8UrlLinkMatch, _ = regexp.Compile("m3u8url = '(.*?)'")
	m3u8KeyUrlMatch, _ = regexp.Compile(`URI="(.*?)"`)
	m3u8ContentInfoMatch, _ = regexp.Compile(`(http[s]?://.*?/)?(.*?\.m3u8)(\?.*)*`)
	m3u8TsInfoMatch, _ = regexp.Compile(`(http[s]?://.*?/)?(.*?\.ts)(\?.*)*`)
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
	dir  string
	db   *gorm.DB
	lock sync.Mutex
}

type M3U8Video struct {
	gorm.Model
	// modio.VideoContainer
	Link      string // The link to the video
	VideoLink string // The m3u8 of the video
	Title     string // The title of the video
	Img       string // The image of the video
	Desc      string // The description of the video
	// status
	Fetched bool // Whether the video has been fetched
}
type M3U8Content struct {
	gorm.Model
	VideoID    int
	Index      int
	Content    string
	Downloaded bool
}

func (vdb *VideoDatabase) Init(dir string) error {
	vdb.dir = path.Join(dir, "_lsp.db")
	db, err := gorm.Open(sqlite.Open(vdb.dir), &gorm.Config{
		PrepareStmt: true,
	})
	if err != nil {
		fmt.Printf("数据库错误: %s\n", err)
		return err
	}
	vdb.db = db
	db.AutoMigrate(&M3U8Video{})
	db.AutoMigrate(&M3U8Content{})
	vdb.lock = sync.Mutex{}
	return nil
}
func (vdb *VideoDatabase) VideoAdd(video *M3U8Video) error {
	vdb.lock.Lock()
	defer vdb.lock.Unlock()
	return vdb.db.Create(video).Error
}
func (vdb *VideoDatabase) VideoLen() (int64, error) {
	vdb.lock.Lock()
	defer vdb.lock.Unlock()
	var count int64
	err := vdb.db.Model(&M3U8Video{}).Count(&count).Error
	return count, err
}
func (vdb *VideoDatabase) VideoGet(id int) (*M3U8Video, error) {
	vdb.lock.Lock()
	defer vdb.lock.Unlock()
	var video M3U8Video
	err := vdb.db.First(&video, id).Error
	return &video, err
}
func (vdb *VideoDatabase) VideoSetFetched(id int, status bool) error {
	vdb.lock.Lock()
	defer vdb.lock.Unlock()
	var video M3U8Video
	err := vdb.db.First(&video, id).Error
	if err != nil {
		return err
	}
	video.Fetched = status
	return vdb.db.Save(&video).Error
}
func (vdb *VideoDatabase) M3U8ContentAdd(content *M3U8Content) error {
	vdb.lock.Lock()
	defer vdb.lock.Unlock()
	return vdb.db.Create(content).Error
}
func (vdb *VideoDatabase) M3U8ContentGet(videoID int, index int) (*M3U8Content, error) {
	vdb.lock.Lock()
	defer vdb.lock.Unlock()
	var content M3U8Content
	err := vdb.db.First(&content, "video_id = ? and index = ?", videoID, index).Error
	return &content, err
}
func (vdb *VideoDatabase) M3U8ContentGetAll(videoID int) ([]*M3U8Content, error) {
	vdb.lock.Lock()
	defer vdb.lock.Unlock()
	var contents []*M3U8Content
	err := vdb.db.Where("video_id = ?", videoID).Find(&contents).Error
	return contents, err
}
func (vdb *VideoDatabase) M3U8ContentLen(videoID int) (int64, error) {
	vdb.lock.Lock()
	defer vdb.lock.Unlock()
	var count int64
	err := vdb.db.Model(&M3U8Content{}).Where("video_id = ?", videoID).Count(&count).Error
	return count, err
}

type M3U8Decoder struct {
	content [][]string
}

func (d *M3U8Decoder) Init(m3u8url string) error {
	buffer := make([][]string, 0)
	ln := M3U8ContentInfoMatch().FindStringSubmatch(m3u8url)
	if len(ln) == 0 {
		return errors.New("m3u8 url error")
	}
	ln[0] = "m"
	buffer = append(buffer, ln)
	err := d.init(&buffer)
	if err != nil {
		return err
	}
	d.content = buffer
	return nil
}
func (d *M3U8Decoder) init(list *[][]string) error {
	ptr := -1
	for i, j := range *list {
		if j[0] == "m" {
			ptr = i
		}
	}
	if ptr == -1 {
		return nil
	}
	domain := (*list)[ptr][1]
	lastDir := UrlDirMatch().FindStringSubmatch(domain + (*list)[ptr][2])[1]
	m3u8url := domain + (*list)[ptr][2] + (*list)[ptr][3]
	fmt.Printf("正在获取 %s\n", m3u8url)
	m3u8Content, err := UrlGetToStr(m3u8url)
	if err != nil {
		return fmt.Errorf("获取失败: %s", err.Error())
	}
	buffer := make([][]string, 0)
	for _, line := range strings.Split(m3u8Content, "\n") {
		if strings.HasPrefix(line, "#EXT-X-KEY:") {
			// todo
			continue
		}
		if strings.HasPrefix(line, "#") {
			continue
		}
		ln := M3U8ContentInfoMatch().FindStringSubmatch(line)
		if ln != nil {
			ln[0] = "m"
			buffer = append(buffer, ln)
		}
		ln = M3U8TsInfoMatch().FindStringSubmatch(line)
		if ln != nil {
			ln[0] = "t"
			buffer = append(buffer, ln)
		}
	}
	for _, i := range buffer {
		if (i[1] == "") && (strings.HasPrefix(i[2], "/")) {
			i[1] = domain
			i[2] = strings.TrimPrefix(i[2], "/")
		} else if i[1] == "" {
			i[1] = lastDir
		}
	}
	buffer = append((*list)[0:ptr], buffer...)
	if len(*list) >= ptr+1 {
		buffer = append(buffer, (*list)[ptr+1:]...)
	}
	*list = buffer
	return d.init(list)
}
func (d *M3U8Decoder) Len() int {
	return len(d.content)
}
func (d *M3U8Decoder) Get(index int) ([]string, error) {
	if index < 0 || index >= d.Len() {
		return nil, errors.New("index out of range")
	}
	return d.content[index], nil
}

type M3U8Downloader struct {
	client *colly.Collector
	wg     *sync.WaitGroup
	buffer [][]byte
}

func (d *M3U8Downloader) Download(video []*M3U8Content, dir string, name string) error {
	d.buffer = make([][]byte, len(video))
	d.client = CollyCollectorSlow()
	d.client.SetRequestTimeout(time.Second * 3)
	d.wg = &sync.WaitGroup{}
	count := len(video)
	downloaded := 0
	t := time.NewTicker(time.Second)
	done := make(chan struct{})
	lock := sync.Mutex{}
	go func() {
		for {
			select {
			case <-done:
				return
			case <-t.C:
				lock.Lock()
				fmt.Printf("正在下载[%d]/[%d]: %d%%\n",
					downloaded, count, 100*downloaded/count,
				)
				lock.Unlock()
			}
		}
	}()
	d.client.OnRequest(func(r *colly.Request) {})
	d.client.OnError(func(r *colly.Response, err error) {
		fmt.Printf("下载 %s 时报错: %s, 正在重试...\n", r.Request.URL, err.Error())
		r.Request.Retry()
	})
	d.client.OnResponse(func(r *colly.Response) {
		if r.StatusCode != 200 {
			fmt.Printf("下载 %s 时, 非预期的状态码: %d, 正在重试...\n", r.Request.URL, r.StatusCode)
			r.Request.Retry()
			return
		}
		index, _ := strconv.ParseInt(r.Ctx.Get("video"), 10, 64)
		d.buffer[index] = r.Body
		lock.Lock()
		downloaded++
		lock.Unlock()
		d.wg.Done()
	})
	for index, v := range video {
		d.wg.Add(1)
		ctx := colly.NewContext()
		ctx.Put("video", strconv.FormatInt((int64)(index), 10))
		d.client.Request("GET", v.Content, nil, ctx, nil)
	}
	d.wg.Wait()
	t.Stop()
	close(done)
	fmt.Println("正在保存...")
	fileBytes := bytes.Buffer{}
	for _, v := range d.buffer {
		fileBytes.Write(v)
	}
	// write to file
	fileName := fmt.Sprintf("%s/%s.mp4", dir, name)
	f, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.Write(fileBytes.Bytes())
	if err != nil {
		return err
	}
	return nil
}
