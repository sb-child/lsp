# LaoSePi - 人人都是老色批
## 颜色网站视频爬取工具

# 安装依赖
windows用户, 将`pip3`替换为`pip`
```
pip3 install requests pycrypto beautifulsoup4 lxml --user
```

# 文件描述
> `test.py` 获取链接, 生成下载脚本
> `getLinks.py` 模块: 获取网页链接  
> `getVideoLink.py` 模块: 获取视频链接  
> `decryptLink.py` 模块: 解密视频链接  

# todo
> 域名有时会更换, 尝试自动获取域名    
> 链接数据库, 用于去重    
