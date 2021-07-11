package main

import (
	"crypto/rand"
	"flag"
	"fmt"
	_ "mods/miya"
	mods "mods/modio"
	_ "mods/yysp"
	"strings"
	"sync"
	"time"

	"github.com/gookit/color"
)

type task struct {
	mod  *mods.ModuleIO
	tags []string
	dir  string
}

func getDownloadDir() string {
	randBytes := make([]byte, 8)
	rand.Reader.Read(randBytes)
	return fmt.Sprintf("v_auto_%x", randBytes)
}

func main() {
	fmt.Println("[sb-child/lsp]视频爬取工具 Golang版本")
	var (
		download_dir      string
		tags_string       string
		selected_mod_name string
		mods_get          bool
		list_get          bool
		tags_get          bool
	)
	flag.StringVar(&selected_mod_name, "mod", "", "指定要加载的模块")
	flag.BoolVar(&mods_get, "mods", false, "可选: 获取当前可选模块并终止")
	flag.StringVar(&download_dir, "dir", getDownloadDir(), "可选: 指定下载目录")
	flag.StringVar(&tags_string, "tag", "", "可选: 指定标签(编号)并终止, 用英文逗号分隔, 可指定多个, 否则为默认")
	flag.BoolVar(&tags_get, "tags", false, "可选: 获取当前模块中, 全部可用的标签并终止")
	flag.BoolVar(&list_get, "list", false, "可选: 仅拉取视频列表, 不下载")
	flag.Parse()
	if mods_get {
		fmt.Println("可用模块:")
		for k, v := range mods.GetAllModules() {
			fmt.Printf("模块名[%s] 模块描述[%s]\n", k, (*v).ModDesc())
		}
		return
	}
	if selected_mod_name == "" {
		fmt.Println("请指定一个模块")
		return
	}
	if tags_get {
		return
	}
	if list_get {
		return
	}
	mod := mods.GetModule(selected_mod_name)
	if mod == nil {
		fmt.Printf("找不到[%s]模块\n", selected_mod_name)
		return
	}
	_tags := strings.Split(tags_string, ",")
	tags := make([]string, 0)
	for _, v := range _tags {
		if v == "" {
			continue
		}
		tags = append(tags, v)
	}
	run(task{
		mod:  mod,
		tags: tags,
		dir:  download_dir,
	})
}
func run(t task) {
	mod := t.mod
	tags := t.tags
	上次调用时间 := -1.0
	输出锁 := sync.Mutex{}
	p_head := func(换行 bool) {
		输出锁.Lock()
		当前时间 := float64(time.Now().UnixNano()) / 1000000000
		if 上次调用时间 < 0 {
			color.Info.Print("-.  ")
			上次调用时间 = 当前时间
		} else {
			上次用时 := 当前时间 - 上次调用时间
			上次调用时间 = 当前时间
			color.Info.Printf("%.2f", 上次用时)
		}
		fmt.Print("^")
		color.Primary.Printf("%s", (*mod).ModName())
		if 换行 {
			color.Warn.Println("...")
			return
		}
		fmt.Print("> ")
	}
	需要换行 := func(s string) bool {
		a := strings.Index(s, "\n")
		// a == -1 : a里面没有换行符
		return a != -1
	}
	p_info := func(s string) {
		p_head(需要换行(s))
		color.Info.Println(s)
		输出锁.Unlock()
	}
	p_warn := func(s string) {
		p_head(需要换行(s))
		color.Warn.Println(s)
		输出锁.Unlock()
	}
	p_err := func(s string) {
		p_head(需要换行(s))
		color.Error.Println(s)
		输出锁.Unlock()
	}
	(*mod).OnInfo(p_info)
	(*mod).OnWarn(p_warn)
	(*mod).OnError(p_err)
	(*mod).Init()

	fmt.Printf("将下载到[%s]目录\n", getDownloadDir())
	fmt.Printf("%T %v\n", tags, tags)
	fmt.Printf("%v\n", *mod)
}
