# yysp: 夜夜視頻資源站
# 视频拉取模块
import re
import time
import urllib3
import getLinks
import getVideoLink
import tsDecode
from downloader import urlGetToStr

domain_re = re.compile(r'https://(.*?)/')
link_re = re.compile(r'"url":"(.*?)"')


def _decoder(url: str):
    # 多视频来源解码
    if url.endswith(".m3u8"):
        # 不需要处理
        return url, ""
    result = getLinks.myReqGet(url)
    new_domain = domain_re.findall(url)[0]
    # https://vip4.ddyunbo.com
    new_dir = link_re.findall(result)[0]
    new_url = "https://" + new_domain + new_dir
    # var playlist = '[{"url":"/20210501/ItF819cE/600kb/hls/index.m3u8"}]';
    return new_url, new_domain


class Puller:
    def log(self, message: str):
        print(self.modName + ":", message)

    def __init__(self):
        self.modName = "[yysp]拉取模块"
        self.log("载入模块...")
        self._lg = getLinks.GetterYysp()
        self._dldPool = urllib3.PoolManager()
        self.lastLinks: list[tuple[str, str, str]] = []
        self.log("模块已准备就绪.")

    def fetch(self):
        self.log("拉取视频列表...")
        self.lastLinks = self._lg.run()
        self.log("视频列表拉取完成.")

    def getDownloadLink(self, index: int):
        self.log("获取下载链接...")
        link_url = self.lastLinks[index]
        print(f"视频{index}:\n" + getLinks.linkFormat(link_url))
        this_link = getVideoLink.getLink(link_url[0])
        this_url, domain = _decoder(this_link)
        print(f"* 下载链接:", this_url)
        video_list_str = urlGetToStr(self._dldPool, this_url)
        videos_list = tsDecode.decoder(video_list_str)
        for i in range(len(videos_list)):
            videos_list[i] = domain + videos_list[i]
        print(f"* 视频时长:", time.strftime("%H:%M:%S", time.gmtime(tsDecode.videoLen(video_list_str))))
        self.log("获取完成.")
        return {
            "list": videos_list,
            "links": link_url,
        }
