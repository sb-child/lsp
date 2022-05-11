<div align='center'>
<img align='left' src='imgs/logo.svg' width='200px'>
<h1>
(lsp)老色批
</h1>
<h3>
颜色网站视频爬取下载工具
</h3>
<div align='center'>
<h4>已适配网站列表</h4>
<span>~~miya(蜜芽)~~ 蜜芽坏掉了!</span>
<span>yysp(夜夜视频)</span>
</div>
<br><br><br>
<h2>喜欢本项目? 不妨点个star~</h2>
</div>

## 已经适配了两个网站的 Python 版本的分支在下面

## [python 分支直达链接](../../tree/python)

## 这个分支是 Go 版本, 至少能用啦~

# 怎么用呢？

请事先安装好`go`和`ffmpeg`~

```bash
$ go get -v all
$ go run main.go [参数...]
```

下载完成后, 打开下载目录:

```bash
_lsp.db # 数据库文件，用于断点下载，全部下载完成后可以删除
1.mp4 # 原视频（可能无法正常观看），全部下载完成后可以删除
1.mp4.o.mp4 # 经过编码后的视频
2.mp4 # 第2个原视频
2.mp4.o.mp4 # 第2个经过编码后的视频
3.mp4 # 第3个原视频
3.mp4.o.mp4 # 第3个经过编码后的视频
...
```

就这么简单，\(溜了溜了~\)
