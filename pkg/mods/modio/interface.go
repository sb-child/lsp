package modio

type VideoContainer struct {
	Link      string // The link to the video
	VideoLink string // The m3u8 of the video
	Title     string // The title of the video
	Img       string // The image of the video
	Desc      string // The description of the video
}

type ModuleIO interface {
	ModName() string // ModDesc 模块介绍
	ModDesc() string // ModName 模块名
	Init() bool      // Init 初始化
	OnSucc(func(string))
	OnInfo(func(string))
	OnWarn(func(string))
	OnError(func(string))
	GetAllTags() map[string]string       // GetAllTags 获取所有分类
	GetVideos([]string) []VideoContainer // GetVideos 获取指定分类或默认分类的视频网页链接
}

var mods map[string]*ModuleIO

func init() {
	mods = make(map[string]*ModuleIO)
}

func AddModule(name string, mod *ModuleIO) {
	mods[name] = mod
}

func GetModule(name string) *ModuleIO {
	return mods[name]
}

func GetAllModules() map[string]*ModuleIO {
	return mods
}
