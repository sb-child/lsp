package main

import (
	"crypto/rand"
	"flag"
	"fmt"
	_ "mods/miya"
	mods "mods/modio"
	"mods/mtools"
	_ "mods/yysp"
	"os"
	"strings"

	"github.com/gookit/color"
)

type task struct {
	mod           *mods.ModuleIO
	tags          []string
	dir           string
	get_tags_only bool
}

func getDownloadDir() string {
	randBytes := make([]byte, 16)
	rand.Reader.Read(randBytes)
	return fmt.Sprintf("v_auto_%x", randBytes)
}

func main() {
	fmt.Println("[sb-child/lsp]视频爬取工具 Golang版本")
	var (
		下载目录    string
		标签字符串   string
		选中的模块   string
		仅获取可选模块 bool
		仅获取视频列表 bool
		仅获取全部标签 bool
	)
	flag.StringVar(&选中的模块, "mod", "", "指定要加载的模块")
	flag.BoolVar(&仅获取可选模块, "mods", false, "可选: 获取当前可选模块并终止")
	flag.StringVar(&下载目录, "dir", getDownloadDir(), "可选: 指定下载目录")
	flag.StringVar(&标签字符串, "tag", "", "可选: 指定分类(编号)并终止, 用英文逗号分隔, 可指定多个, 否则为默认")
	flag.BoolVar(&仅获取全部标签, "tags", false, "可选: 获取当前模块中, 全部可用的分类并终止")
	flag.BoolVar(&仅获取视频列表, "list", false, "可选: 仅拉取视频列表, 不下载")
	flag.Parse()
	if 下载目录 == "" && !仅获取视频列表 {
		fmt.Println("请指定一个下载目录")
		os.Exit(4)
		return
	}
	if 仅获取可选模块 {
		fmt.Println("可用模块:")
		for k, v := range mods.GetAllModules() {
			fmt.Printf("模块名[%s] 模块描述[%s]\n", k, (*v).ModDesc())
		}
		return
	}
	if 选中的模块 == "" {
		fmt.Println("请指定一个模块")
		os.Exit(3)
		return
	}
	if 仅获取视频列表 {
		下载目录 = ""
	}
	模块实例 := mods.GetModule(选中的模块)
	if 模块实例 == nil {
		fmt.Printf("找不到[%s]模块\n", 选中的模块)
		os.Exit(2)
		return
	}
	_tags := strings.Split(标签字符串, ",")
	标签列表 := make([]string, 0)
	for _, v := range _tags {
		if v == "" {
			continue
		}
		标签列表 = append(标签列表, v)
	}
	fmt.Printf("载入[%s]模块...\n", 选中的模块)
	run(task{
		mod:           模块实例,
		tags:          标签列表,
		dir:           下载目录,
		get_tags_only: 仅获取全部标签,
	})
}
func run(t task) {
	mod := t.mod
	dld_dir := t.dir
	// tags := t.tags
	mc := mtools.NewMyColor((*mod).ModName())
	(*mod).OnSucc(mc.Succ)
	(*mod).OnInfo(mc.Info)
	(*mod).OnWarn(mc.Warn)
	(*mod).OnError(mc.Err)
	succ := (*mod).Init()
	if !succ {
		fmt.Println("初始化失败")
		os.Exit(1)
		return
	}
	fmt.Println("正在获取分类...")
	tags := (*mod).GetAllTags()
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
	for _, v := range t.tags {
		if !_checkTag(v) {
			fmt.Printf("%s 不属于任何一个分类\n", v)
			os.Exit(10)
			return
		}
	}
	fmt.Print("初始化完成, 获取视频列表...")
	r := (*mod).GetVideos(t.tags)
	for _, v := range r {
		fmt.Printf("%v\n", v)
	}
	if len(dld_dir) == 0 {
		fmt.Println("")
	} else {
		fmt.Printf("将下载到[%s]目录\n", dld_dir)
	}

}
