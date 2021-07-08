package main

import (
	"crypto/rand"
	"flag"
	"fmt"
	"mods"
)

func getDownloadDir() string {
	randBytes := make([]byte, 8)
	rand.Reader.Read(randBytes)
	return fmt.Sprintf("v_auto_%x", randBytes)
}

func main() {
	fmt.Println("[sb-child/lsp]视频爬取工具")
	var (
		mod_string string
		mods_get   bool
	)
	flag.StringVar(&mod_string, "mod", "", "指定要加载的模块")
	flag.BoolVar(&mods_get, "mods", false, "获取当前可选模块")
	flag.StringVar(&mod_string, "dir", getDownloadDir(), "可选: 指定下载目录")
	flag.StringVar(&mod_string, "tag", "", "可选: 指定标签(编号), 可指定多个, 否则为默认")
	flag.StringVar(&mod_string, "tags", "", "可选: 获取当前模块中, 全部可用的标签")
	flag.StringVar(&mod_string, "list", "", "可选: 仅拉取视频列表, 不下载")
	flag.Parse()

	var mod mods.ModuleIO

	_ = mod

	fmt.Printf("将下载到[%s]目录\n", getDownloadDir())
	fmt.Printf("将下载到[%s]目录\n", getDownloadDir())

}
