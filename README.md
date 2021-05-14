# LaoSePi - 人人都是老色批
## 颜色网站视频爬取工具

# 安装依赖
> 需要事先安装python和ffmpeg, windows用户需要配置环境变量.  
> 同时, windows用户需要安装 [vsBuildTools](https://visualstudio.microsoft.com/zh-hans/thank-you-downloading-visual-studio/?sku=Community&rel=15#) 以安装`pycrypto`库  
> > 注意: 在安装界面勾选 `c++ 生成工具` 组件  
> > 若安装失败, 可考虑使用`虚拟机`, `git bash`或者`linux子系统`运行此项目

> linux用户:
```shell
pip3 install requests pycrypto beautifulsoup4 lxml pycryptodome tqdm retrying --user
pip3 uninstall pycrypto
pip3 uninstall pycryptodome
pip3 install pycryptodome
```
> windows用户: 运行 `init.cmd`

# 获取视频
## \[推荐\] 模块化全自动下载
> `modRun.py` 命令行参数:  
> `--mod 模块名` 指定要加载的模块  
> > 可用模块:  
> `miya` \(`蜜芽`\)  
> `yysp` \(`夜夜視頻資源站`\)  

> `--dir 目录名` 指定下载目录(可选)  

> `--tags` 获取当前模块中, 全部可用的标签

> `--tag 标签编号` 指定标签(编号), 否则为默认

> `--not-download` 仅拉取视频列表, 不下载

> linux用户: `python3 modRun.py --mod 模块`  
> windows用户: 运行 `modRunMiya.cmd` 或 `modRunYysp.cmd`

[comment]: <> (## \[推荐\] 全自动下载)

[comment]: <> (> linux用户: `python3 test2.py`  )

[comment]: <> (> windows用户: 运行 `auto_download.cmd`)

## \[不推荐\]生成脚本方式
### 生成下载脚本
> linux用户: `python3 test.py`  
> windows用户: 运行 `gen_link.cmd`  
### 运行下载脚本
> linux 用户: `bash dld.sh`  
> windows 用户: 运行 `dld_windows.bat`

# 文件描述
> `test.py` 获取链接, 生成下载脚本  
> `test2.py` 全自动下载脚本  
> `getLinks.py` 模块: 获取网页链接  
> `getVideoLink.py` 模块: 获取视频链接  
> `decryptLink.py` 模块: 解密视频链接  
> `tsDecode.py` 模块: 解码m3u8列表

# todo
> 尝试适配其他网站  
> 链接数据库, 用于去重
