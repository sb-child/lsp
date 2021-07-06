<div align='center'>
<img align='left' src='imgs/logo.svg' width='200px'>
<h1>
(lsp)老色批
</h1>
<h3>
颜色网站视频爬取下载工具
</h3>
<div align='center'>
<h4>已适配网站列表(<code>python</code>分支)</h4>
<span>miya(蜜芽)</span>
<span>yysp(夜夜视频)</span>
</div>
<br><br><br>
<h2>喜欢本项目? 不妨点个star~</h2>
</div>

## 已经适配了两个网站的Python版本的分支在下面
## [python分支直达链接](../python)
## 这个分支是Golang版本, 还在开发, 尚不成熟

<!-- # 快速入手
+ 安卓手机用户, 可使用`termux`等终端模拟器, 按照linux用户的步骤部署环境并运行.
  - 将在手机端串行解密ts文件, 速度较慢
+ 需要事先安装`python 3.8+`和`ffmpeg`, windows用户需要配置环境变量.  
  [windows下,ffmpeg安装教程](https://bbs.huaweicloud.com/blogs/243409)  
  [python阿里镜像](https://npm.taobao.org/mirrors/python/)  
  pypi清华镜像地址 https://opentuna.cn/pypi/web/simple/
+ 同时, windows用户需要安装 [vsBuildTools](https://visualstudio.microsoft.com/zh-hans/thank-you-downloading-visual-studio/?sku=Community&rel=15#) 以安装`pycrypto`库  
  - 注意: 在安装界面勾选 `c++ 生成工具` 组件
  - 若安装失败, 可考虑使用`虚拟机`, `git bash`或者`linux子系统`运行此项目

+ 安装 [jsdec-tiny](https://github.com/sb-child/jsdec-tiny) 插件  
  - 按照readme编译完成后, 将其`build`目录下的文件复制到本项目的`jsdec-tiny`目录下(若没有此目录, 可手动创建)

+ 安装python第三方库
  - linux用户:
    ```shell
    $# 如果你想使用python3运行脚本, 请执行下面的命令
    $ pip3 install requests pycrypto beautifulsoup4 lxml pycryptodome tqdm retrying colorama --user
    $ pip3 uninstall -y pycrypto
    $ pip3 uninstall -y pycryptodome
    $ pip3 install pycryptodome
    $# 如果你想使用pypy3运行脚本, 请执行下面的命令
    $ pypy3 -m pip install requests beautifulsoup4 lxml pycryptodome tqdm retrying colorama --user
    ```
  - windows用户: 运行 `init.cmd`
+ 现在你已经有了运行此脚本所需的环境, 可以[获取并下载视频](#获取并下载视频)啦!

# 获取并下载视频
## 通过命令行工具`modRun.py`下载视频
`modRun.py` 的命令行参数:
+ `--mod 模块名` 指定要加载的模块:
  - 可用模块:
    - `miya` \(蜜芽\)
    - `yysp` \(夜夜視頻資源站\)
  
+ `--dir 目录名` 可选: 指定下载目录
  - 当指定的目录上次下载时被中断或者崩溃, [将从上次进度下载](#断点下载)
+ `--tags` 可选: 获取当前模块中, 全部可用的标签
+ `--tag 标签编号 [标签编号...]` 可选: 指定标签(编号), 可指定多个, 否则为默认
+ `--not-download` 可选: 仅拉取视频列表, 不下载

### 运行脚本
+ linux用户:
  - `python3 modRun.py --mod 模块 [其他参数]`  
+ windows用户:
  - 运行 `modRunMiya.cmd` 或 `modRunYysp.cmd`, 也可自己在cmd中输入命令:
  - `python modRun.py --mod 模块 [其他参数]`

# 断点下载
下载中断后, 下次执行时只需在命令行中指定上次的下载目录即可  
例子:  
```shell
$ python3 modRun.py --mod miya
视频将下载到[xxx]目录
...
下载视频0...
$# 当脚本获取完链接或者下载时, 结束脚本进程
^C
$ python3 modRun.py --mod miya --dir xxx
上次的下载未完成, 将从上次的进度下载:
* 当前进度: 视频xxx / xxx
* 下载失败: 视频[xxx]
$# 此时, 脚本将自动从上次的进度开始下载
```

# 文件描述
+ `modRun.py` 模块启动器
+ `mod_*.py` 视频爬取模块
+ `downloader.py`  组件: 视频下载器
+ `getLinks.py` 组件: 获取网页链接
+ `getVideoLink.py` 组件: 获取视频链接
+ `decryptLink.py` 组件: 解密视频链接
+ `tsDecode.py` 组件: m3u8列表解码器
+ `videoLock.py` 组件: 锁文件操作

# 为此项目添砖加瓦
### 你可以
+ 适配一个新网站
+ 修复一个bug
+ 添加新功能
+ 完善文档 -->
