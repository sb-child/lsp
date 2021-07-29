package mtools

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/gookit/color"
)

type MyColor struct {
	modName string
	输出锁     sync.Mutex
	上次调用时间  float64
	Info    func(string)
	Succ    func(string)
	Warn    func(string)
	Err     func(string)
}

func NewMyColor(mod string) *MyColor {
	mc := MyColor{
		modName: mod,
	}
	mc.Init()
	return &mc
}
func (mc *MyColor) Init() {
	mc.输出锁 = sync.Mutex{}
	mc.上次调用时间 = 0
	p_head := func(换行 bool, 字符 rune) {
		mc.输出锁.Lock()
		当前时间 := float64(time.Now().UnixNano()) / 1000000000
		if mc.上次调用时间 < 0 {
			color.Info.Print("-.  ")
			mc.上次调用时间 = 当前时间
		} else {
			上次用时 := 当前时间 - mc.上次调用时间
			mc.上次调用时间 = 当前时间
			color.Info.Printf("%.2f", 上次用时)
		}
		fmt.Print("^")
		color.Primary.Printf("%s", mc.modName)
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
	mc.Info = func(s string) {
		p_head(需要换行(s), 'i')
		color.Cyan.Println(s)
		mc.输出锁.Unlock()
	}
	mc.Succ = func(s string) {
		p_head(需要换行(s), 'S')
		color.Success.Println(s)
		mc.输出锁.Unlock()
	}
	mc.Warn = func(s string) {
		p_head(需要换行(s), 'W')
		color.Warn.Println(s)
		mc.输出锁.Unlock()
	}
	mc.Err = func(s string) {
		p_head(需要换行(s), 'E')
		color.Error.Println(s)
		mc.输出锁.Unlock()
	}
}
