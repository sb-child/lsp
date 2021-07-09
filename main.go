package main

import (
	"crypto/rand"
	"flag"
	"fmt"
	_ "mods/miya"
	mods "mods/modio"
	_ "mods/yysp"
)

func getDownloadDir() string {
	randBytes := make([]byte, 8)
	rand.Reader.Read(randBytes)
	return fmt.Sprintf("v_auto_%x", randBytes)
}

func main() {
	fmt.Println("[sb-child/lsp]视频爬取工具")
	var (
		mod_string        string
		tags              string
		selected_mod_name string
		mods_get          bool
		list_get          bool
		tags_get          bool
	)
	flag.StringVar(&selected_mod_name, "mod", "", "指定要加载的模块")
	flag.BoolVar(&mods_get, "mods", false, "可选: 获取当前可选模块并终止")
	flag.StringVar(&mod_string, "dir", getDownloadDir(), "可选: 指定下载目录")
	flag.StringVar(&tags, "tag", "", "可选: 指定标签(编号)并终止, 用英文逗号分隔, 可指定多个, 否则为默认")
	flag.BoolVar(&tags_get, "tags", false, "可选: 获取当前模块中, 全部可用的标签并终止")
	flag.BoolVar(&list_get, "list", false, "可选: 仅拉取视频列表, 不下载")
	flag.Parse()
	// miya.MOD_NAME
	if mods_get {
		fmt.Println("可用模块:")
		for k, v := range mods.GetAllModules() {
			fmt.Printf("模块名[%s] 模块描述[%s]\n", k, (*v).ModDesc())
		}
		return
	}

	if tags_get {
		return
	}
	if list_get {
		return
	}

	fmt.Printf("将下载到[%s]目录\n", getDownloadDir())
	fmt.Printf("%s\n", tags)
	fmt.Printf("%v\n", mods.GetModule("miya"))
}
