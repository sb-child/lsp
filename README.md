# lsp (强行老色批)
## 颜色网站视频爬取下载工具
## 完美适配:
+ miya(蜜芽)
+ yysp(夜夜视频)

# 安装依赖
+ 需要事先安装python和ffmpeg, windows用户需要配置环境变量.
+ 同时, windows用户需要安装 [vsBuildTools](https://visualstudio.microsoft.com/zh-hans/thank-you-downloading-visual-studio/?sku=Community&rel=15#) 以安装`pycrypto`库  
+ + 注意: 在安装界面勾选 `c++ 生成工具` 组件
+ + 若安装失败, 可考虑使用`虚拟机`, `git bash`或者`linux子系统`运行此项目

+ linux用户:
```shell
# python3
pip3 install requests pycrypto beautifulsoup4 lxml pycryptodome tqdm retrying --user
pip3 uninstall pycrypto
pip3 uninstall pycryptodome
pip3 install pycryptodome
# pypy3
pypy3 -m pip install requests beautifulsoup4 lxml pycryptodome tqdm retrying --user
```
+ windows用户: 运行 `init.cmd`

# 获取视频
## 通过命令行工具下载视频
+ `modRun.py` 命令行参数:
+ `--mod 模块名` 指定要加载的模块
+ + 可用模块:
+ + `miya` \(`蜜芽`\)
+ + `yysp` \(`夜夜視頻資源站`\)
+ `--dir 目录名` 指定下载目录(可选)  
+ `--tags` 获取当前模块中, 全部可用的标签
+ `--tag 标签编号` 指定标签(编号), 否则为默认
+ `--not-download` 仅拉取视频列表, 不下载
+
+ linux用户: `python3 modRun.py --mod 模块 [其他参数]`  
+ windows用户: 运行 `modRunMiya.cmd` 或 `modRunYysp.cmd`, 也可自己在cmd中输入命令来使用其他功能

# 文件描述
+ `modRun.py` 模块启动器
+ `mod_*.py` 视频爬取模块
+ `downloader.py`  组件: 视频下载器
+ `getLinks.py` 组件: 获取网页链接
+ `getVideoLink.py` 组件: 获取视频链接
+ `decryptLink.py` 组件: 解密视频链接
+ `tsDecode.py` 组件: m3u8列表解码器

# todo

[comment]: <> (+ `miya` 最近更换了域名..)
+ 尝试适配其他网站  
+ 链接数据库, 用于去重
