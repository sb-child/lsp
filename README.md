# LaoSePi - 人人都是老色批
## 颜色网站视频爬取工具

# ~~链接失效问题~~
> ~~详见 `getLinks.py` 文件的第7行~~

# 安装依赖
> windows用户, 将`pip3`替换为`pip`
```
pip3 install requests pycrypto beautifulsoup4 lxml pycryptodome tqdm --user
```

# 获取视频
> 若需生成下载脚本, 需运行
> `python3 test.py`
> windows用户, 将`python3`替换为`python`

> 若需要全自动下载, 则
> `python3 test2.py`
> windows用户, 将`python3`替换为`python`

# 执行下载脚本
> 需要安装ffmpeg, windows用户需要配置环境变量.

> linux 用户:
```
bash dld.sh
```

> windows 用户:
> 运行 dld_windows.bat

# 文件描述
> `test.py` 获取链接, 生成下载脚本
> `test2.py` 全自动下载脚本
> `getLinks.py` 模块: 获取网页链接  
> `getVideoLink.py` 模块: 获取视频链接  
> `decryptLink.py` 模块: 解密视频链接  

# todo
> ~~域名有时会更换, 尝试自动获取域名~~ 目前可以自动获取域名    
> 链接数据库, 用于去重    
