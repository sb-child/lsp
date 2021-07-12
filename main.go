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
	flag.StringVar(&标签字符串, "tag", "", "可选: 指定标签(编号)并终止, 用英文逗号分隔, 可指定多个, 否则为默认")
	flag.BoolVar(&仅获取全部标签, "tags", false, "可选: 获取当前模块中, 全部可用的标签并终止")
	flag.BoolVar(&仅获取视频列表, "list", false, "可选: 仅拉取视频列表, 不下载")
	flag.Parse()
	if 下载目录 == "" && !仅获取视频列表 {
		fmt.Println("请指定一个下载目录")
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
		return
	}
	if 仅获取全部标签 {
		return
	}
	if 仅获取视频列表 {
		下载目录 = ""
	}
	模块实例 := mods.GetModule(选中的模块)
	if 模块实例 == nil {
		fmt.Printf("找不到[%s]模块\n", 选中的模块)
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
		mod:  模块实例,
		tags: 标签列表,
		dir:  下载目录,
	})
}
func run(t task) {
	mod := t.mod
	// tags := t.tags
	上次调用时间 := -1.0
	输出锁 := sync.Mutex{}
	p_head := func(换行 bool, 字符 rune) {
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
		color.LightMagenta.Printf("%c>", 字符)
	}
	需要换行 := func(s string) bool {
		a := strings.Index(s, "\n")
		// a == -1 : a里面没有换行符
		return a != -1
	}
	p_info := func(s string) {
		p_head(需要换行(s), 'i')
		color.Cyan.Println(s)
		输出锁.Unlock()
	}
	p_succ := func(s string) {
		p_head(需要换行(s), 'S')
		color.Success.Println(s)
		输出锁.Unlock()
	}
	p_warn := func(s string) {
		p_head(需要换行(s), 'W')
		color.Warn.Println(s)
		输出锁.Unlock()
	}
	p_err := func(s string) {
		p_head(需要换行(s), 'E')
		color.Error.Println(s)
		输出锁.Unlock()
	}
	(*mod).OnSucc(p_succ)
	(*mod).OnInfo(p_info)
	(*mod).OnWarn(p_warn)
	(*mod).OnError(p_err)
	(*mod).Init()
	fmt.Printf("初始化完成, 将下载到[%s]目录\n", getDownloadDir())
	// fmt.Printf("%T %v\n", tags, tags)
	// fmt.Printf("%v\n", *mod)
}
