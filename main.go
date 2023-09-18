package main

import (
	"flag"
	"fmt"
	_ "lsp/services/miya"
	mods "lsp/services/modio"
	"lsp/services/mtools"
	_ "lsp/services/yysp"
	"math/rand"
	"os"
	"path"
	"strings"
	"sync"
	"time"

	"github.com/gookit/color"
)

const (
	RAND_LETTERS = "1qaz2wsx3edc4rfv5tgb6yhn7ujm8ik9ol0pQAZWSXEDCRFVTGBYHNUJMIKOLP"
)

type task struct {
	mod           *mods.ModuleIO
	tags          []string
	dir           string
	dbFile        string
	get_tags_only bool
}

func getDownloadDir() string {
	rands := ""
	for i := 0; i < 16; i++ {
		rands += string(RAND_LETTERS[rand.Intn(len(RAND_LETTERS))])
	}
	return fmt.Sprintf("v_auto_%s", rands)
}

func main() {
	fmt.Println("[sb-child/lsp]视频爬取工具 Go版本")
	var (
		downloadDir     string
		dbFile          string
		tagList         string
		selectedMod     string
		getModList      bool
		getVideoList    bool
		getTagList      bool
		writeToDatabase bool
	)
	flag.StringVar(&selectedMod, "mod", "", "指定要加载的模块")
	flag.StringVar(&downloadDir, "dir", getDownloadDir(), "可选: 指定下载目录")
	flag.StringVar(&dbFile, "db", "", "可选: 指定下载目录")
	flag.StringVar(&tagList, "tag", "", "可选: 指定分类(编号)并终止, 用英文逗号分隔, 否则为默认")
	flag.BoolVar(&getModList, "mods", false, "可选: 获取当前可选模块并终止")
	flag.BoolVar(&getTagList, "tags", false, "可选: 获取当前模块中, 全部可用的分类")
	flag.BoolVar(&getVideoList, "list", false, "可选: 仅获取视频列表, 不保存视频信息")
	flag.BoolVar(&writeToDatabase, "save", false, "可选: 仅将视频信息写入数据库")
	flag.Parse()
	if downloadDir == "" && !getVideoList {
		fmt.Println("请指定一个下载目录")
		os.Exit(4)
		return
	}
	if getVideoList && writeToDatabase {
		fmt.Println("请去掉 -list 参数以写入数据库")
		os.Exit(5)
		return
	}
	if getModList {
		fmt.Println("可用模块:")
		for k, v := range mods.GetAllModules() {
			fmt.Printf("模块名[%s] 模块描述[%s]\n", k, (*v).ModDesc())
		}
		return
	}
	if selectedMod == "" {
		fmt.Println("请指定一个模块")
		os.Exit(3)
		return
	}
	if getVideoList {
		downloadDir = ""
	}
	模块实例 := mods.GetModule(selectedMod)
	if 模块实例 == nil {
		fmt.Printf("找不到[%s]模块\n", selectedMod)
		os.Exit(2)
		return
	}
	_tags := strings.Split(tagList, ",")
	标签列表 := make([]string, 0)
	for _, v := range _tags {
		if v == "" {
			continue
		}
		标签列表 = append(标签列表, v)
	}
	fmt.Printf("载入[%s]模块...\n", selectedMod)
	run(task{
		mod:           模块实例,
		tags:          标签列表,
		dir:           downloadDir,
		get_tags_only: getTagList,
		dbFile:        dbFile,
	})
}

func printVideoList(video []mods.VideoContainer) {
	for k, v := range video {
		fmt.Print("[")
		color.Cyan.Printf("%d", k)
		fmt.Print("]: ")
		color.Yellow.Print(v.Title)
		fmt.Print("\n 链接 ")
		color.Blue.Println(v.Link)
		fmt.Print(" 封面 ")
		color.Red.Println(v.Img)
		fmt.Print(" 视频 ")
		color.Green.Println(v.VideoLink)
		fmt.Print(" 描述 ")
		color.Comment.Println(v.Desc)
		// md := mtools.M3U8Decoder{}
		// md.Init(v.VideoLink)
	}
}

func run(t task) {
	mod := t.mod
	dld_dir := t.dir
	mc := mtools.NewMyColor((*mod).ModName())
	(*mod).OnSucc(mc.Succ)
	(*mod).OnInfo(mc.Info)
	(*mod).OnWarn(mc.Warn)
	(*mod).OnError(mc.Err)
	// check if dld_dir exists
	if fi, err := os.Stat(dld_dir); (err == nil) && (fi.IsDir()) {
		if fi, err := os.Stat(t.dbFile); (err == nil) && (!fi.IsDir()) {
			goto skip
		} else if fi, err := os.Stat(path.Join(dld_dir, "_lsp.db")); (err == nil) && (!fi.IsDir()) {
			goto skip
		}
		fmt.Println("找不到数据库文件, 从网站拉取列表...")
	}
	// 初始化
	if succ := (*mod).Init(); !succ {
		fmt.Println("初始化失败")
		os.Exit(1)
		return
	}
	// 获取分类, 检查分类是否正确
	func() {
		fmt.Println("获取分类...")
		tags := (*mod).GetAllTags()
		fmt.Printf("共[%d]个\n", len(tags))
		if t.get_tags_only {
			for k, v := range tags {
				fmt.Print("分类[")
				color.Cyan.Print(v)
				fmt.Print("] 编号[")
				color.Cyan.Print(k)
				fmt.Print("]\n")
			}
			return
		}
		_checkTag := func(s string) bool {
			_, ok := tags[s]
			return ok
		}
		tags_temp := make(map[string]struct{})
		for _, v := range t.tags {
			if !_checkTag(v) {
				fmt.Printf("[%s]不属于任何一个分类\n", v)
				os.Exit(10)
				return
			}
			if _, ok := tags_temp[v]; ok {
				fmt.Printf("重复分类[%s]\n", v)
				os.Exit(10)
				return
			}
			tags_temp[v] = struct{}{}
		}
	}()
	// === 初始化完成 ===
	// 获取视频列表
	func() {
		fmt.Println("获取视频列表...")
		r := (*mod).GetVideos(t.tags)
		if len(dld_dir) == 0 {
			printVideoList(r)
			return
		}
		// 输出下载路径, 视频个数. 创建目录
		fmt.Printf("准备下载[%s]<-[%d]\n", dld_dir, len(r))
		os.Mkdir(dld_dir, os.ModePerm)
		// 保存解析结果
		fmt.Println("正在保存到数据库...")
		db := mtools.VideoDatabase{}
		if err := db.Init(dld_dir, t.dbFile); err != nil {
			os.Exit(1)
			return
		}
		for _, v := range r {
			mv := mtools.M3U8Video{
				Title:     v.Title,
				Link:      v.Link,
				Img:       v.Img,
				Desc:      v.Desc,
				VideoLink: v.VideoLink,
				Fetched:   false,
			}
			err := db.VideoAdd(&mv)
			if err != nil {
				fmt.Printf("保存时发生错误: %s\n", err.Error())
				os.Exit(1)
				return
			}
		}
	}()
skip:
	// 提取ts列表
	if err := fetchTs(dld_dir, t.dbFile); err != nil {
		fmt.Println(err)
		os.Exit(1)
		return
	}
	// 下载，合并，保存视频
	fmt.Println("开始下载...")
	if err := download(dld_dir, t.dbFile); err != nil {
		fmt.Println(err)
		os.Exit(1)
		return
	}
}
func fetchTs(dir, dbFile string) error {
	fmt.Println("读取数据库...")
	db := mtools.VideoDatabase{}
	if err := db.Init(dir, dbFile); err != nil {
		return err
	}
	fmt.Println("解析链接...")
	count, _ := db.VideoLen()
	t := time.NewTicker(time.Second)
	done := make(chan struct{})
	defer t.Stop()
	defer close(done)
	lock := sync.Mutex{}
	status := struct {
		VideoCount int
		VideoDone  int
		TsCount    int
		TsDone     int
	}{
		VideoCount: (int)(count),
		VideoDone:  0,
		TsCount:    1,
		TsDone:     1,
	}
	go func() {
		for {
			select {
			case <-done:
				return
			case <-t.C:
				lock.Lock()
				fmt.Printf("视频[%d](%d%%): 正在存储[%d]/[%d]: %d%%\n",
					status.VideoDone,
					(status.VideoDone*100)/status.VideoCount,
					status.TsDone,
					status.TsCount,
					(status.TsDone*100)/status.TsCount,
				)
				lock.Unlock()
			}
		}
	}()
	save := func(decoder mtools.M3U8Decoder, video *mtools.M3U8Video, tsIndex int, tsCount int, videoIndex int, videoCount int) error {
		ts, err := decoder.Get(tsIndex)
		if err != nil {
			return err
		}
		link := ts[1] + ts[2] + ts[3]
		key := ""
		if len(ts) == 5 {
			key = ts[4]
		}
		db.M3U8ContentAdd(&mtools.M3U8Content{
			VideoID:    int(video.ID),
			Index:      tsIndex,
			Content:    link,
			Downloaded: false,
			Key:        key,
		})
		lock.Lock()
		status.TsCount = tsCount
		status.TsDone = tsIndex + 1
		status.VideoCount = videoCount
		status.VideoDone = videoIndex
		lock.Unlock()
		return nil
	}
	for i := 1; i <= (int)(count); i++ {
		v, err := db.VideoGet(i)
		if v.Fetched {
			continue
		}
		if err != nil {
			return err
		}
		decoder := mtools.M3U8Decoder{}
		err = decoder.Init(v.VideoLink)
		if err != nil {
			return err
		}
		tsLen := decoder.Len()

		for j := 0; j < tsLen; j++ {
			save(decoder, v, j, tsLen, i, (int)(count))
		}
		db.VideoSetFetched(i, true)
	}
	return nil
}

func download(dir, dbFile string) error {
	fmt.Println("读取数据库...")
	db := mtools.VideoDatabase{}
	downloader := mtools.M3U8Downloader{}
	if err := db.Init(dir, dbFile); err != nil {
		return err
	}
	videoCount, _ := db.VideoLen()
	for i := 1; i <= (int)(videoCount); i++ {
		content, _ := db.M3U8ContentGetAll(i)
		videoDesc, _ := db.VideoGet(i)
		if videoDesc.Downloaded {
			fmt.Printf("跳过已下载的视频[%d]...\n", i)
			continue
		}
		videoTitleShort := videoDesc.Title
		if len([]rune(videoTitleShort)) > 10 {
			videoTitleShort = substr(videoTitleShort, 0, 10) + "..."
		}
		downloader.Download(content, dir, fmt.Sprintf("%d", videoDesc.ID), fmt.Sprintf("[%d/%d]%s", i, videoCount, videoTitleShort))
		db.VideoSetDownloaded(i, true)
	}
	return nil
}

func substr(input string, start int, length int) string {
	asRunes := []rune(input)
	if start >= len(asRunes) {
		return ""
	}
	if start+length > len(asRunes) {
		length = len(asRunes) - start
	}
	return string(asRunes[start : start+length])
}
