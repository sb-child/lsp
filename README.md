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
<span><del>miya(蜜芽)</del> 蜜芽坏掉了! </span>
<span>yysp(夜夜视频)</span>
</div>
<br><br><br>
<h2>喜欢本项目? 不妨点个star~</h2>
</div>

## 已经适配了两个网站的 Python 版本的分支在下面

## [python 分支直达链接](../../tree/python)

## 这个分支是 Go 版本, 至少能用啦~

# 怎么用呢？

## Windows 用户看过来~

选择你的阵营：

### 自己编译

1. 下载本 repo
2. 安装 `go 1.20+`
3. 把 [`ffmpeg`](https://github.com/BtbN/FFmpeg-Builds/releases/download/latest/ffmpeg-master-latest-win64-gpl.zip) 压缩包里 `ffmpeg-master-latest-win64-gpl/bin/` 的文件都放进 `C:\Windows\` 里面
4. 开始使用吧：

```bash
go get -v all
go run main.go -mod yysp [参数...]
```

### 懒得编译啦，下载预编译版

1. 在 ![GitHub release](https://img.shields.io/github/v/release/sb-child/lsp) 下载 `lsp_windows-amd64.exe`
2. 把 [`ffmpeg`](https://github.com/BtbN/FFmpeg-Builds/releases/download/latest/ffmpeg-master-latest-win64-gpl.zip) 压缩包里 `ffmpeg-master-latest-win64-gpl/bin/` 的文件都放进 `C:\Windows\` 里面
3. 在下载目录打开命令行，开始使用吧：

```bash
lsp_windows-amd64.exe -mod yysp [参数...]
```

## Linux 用户看过来~

选择你的阵营：

### 自己编译

1. 下载本 repo
2. 安装 `go 1.20+` 和 `ffmpeg`
3. 开始使用吧：

```bash
go get -v all
go run main.go -mod yysp [参数...]
```

### 懒得编译啦，下载预编译版

1. 在 ![GitHub release](https://img.shields.io/github/v/release/sb-child/lsp) 下载 `lsp_linux-amd64`，并赋予可执行权限 `chmod +x lsp_linux-amd64`
2. 安装 `ffmpeg`
3. 在下载目录打开命令行，开始使用吧：

```bash
./lsp_linux-amd64 -mod yysp [参数...]
```
## 能用了之后

### 查看可用参数

使用 `--help` 参数查看帮助

### 下载目录结构

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
