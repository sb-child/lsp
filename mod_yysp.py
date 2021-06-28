# yysp: 夜夜視頻資源站
# 视频拉取模块

import time
import bs4
import getLinks
import getVideoLink
import tsDecode
import modBase

from typing import Union
from downloader import urlGetToStr

domain_re = modBase.domain_re
link_re = modBase.link_re
link2_re = modBase.link2_re
tag_re = modBase.tag_re


def _decoder(url: str):
    # 多视频来源解码
    new_domain: str = domain_re.findall(url)[0] + "/"
    if url.endswith(".m3u8"):
        # 不需要处理
        return url, new_domain
    result = getLinks.myReqGet(url)
    try:
        new_dir: str = link_re.findall(result)[0]
        new_dir = new_domain + new_dir
    except IndexError:
        new_dir: str = link2_re.findall(result)[0]
    return new_dir, new_domain


class Puller(modBase.Puller):
    def _init(self):
        self._name("yysp")
        self._lg = getLinks.GetterYysp()

    def _getTags(self):
        body = getLinks.myReqGet(self._lg.base_url)
        body_bs = bs4.BeautifulSoup(body, "lxml")
        tags = body_bs.find_all("a", {"class": "1=0"})
        tags_list = {}
        for i in range(len(tags)):
            tag: bs4.Tag = tags[i]
            # /index.php/vod/type/id/33.html
            tag_name = tag_re.findall(tag.attrs["href"])
            if len(tag_name) == 0:
                continue
            tags_list[tag_name[0]] = tag.text.strip()
            # {"33": "..."}
        self.lastTags = tags_list

    def _setTag(self, tag_name: list):
        self.selectedTag = tag_name

    def _fetch(self):
        if len(self.selectedTag) == 0:
            self.lastLinks = self._lg.run()
            return
        temp = []
        for i in self.selectedTag:
            temp.extend(self._lg.run(tag=str(i)))
        self.lastLinks = temp
        # self.lastLinks = self._lg.run(tag=self.selectedTag)

    def _getDownloadLink(self, index: int):
        link_url = self.lastLinks[index]
        print(f"视频{index}:\n" + getLinks.linkFormat(link_url))
        this_link = getVideoLink.getLink(link_url[0])
        try:
            this_url, domain = _decoder(this_link)
        except Exception as e:
            return {
                "error": str(e),
                "list": [],
                "links": ["-", "-", "-"],
                "encrypt": False,
            }
        this_url: str
        domain: str
        print(f"* 下载链接:", this_url)
        video_list_str = urlGetToStr(this_url)
        _decode_url = this_url[0:(this_url.rfind("/") + 1)]
        videos_list, _ = tsDecode.decoder(video_list_str, _decode_url)
        video_len: float = tsDecode.videoLen(video_list_str, _decode_url)
        print(f"* 视频时长:", time.strftime("%H:%M:%S", time.gmtime(video_len)))
        encrypt: str = tsDecode.checkEncrypt(video_list_str, _decode_url)
        print(f"* 密钥: [{encrypt if encrypt != '' else '无需解密'}]")
        # print(videos_list, encrypt)
        return {
            "list": videos_list,
            "links": link_url,
            "encrypt": encrypt,
        }

# a = Puller()
# a.getTags()
# print(a.lastTags)
# a.setTag("38")
# a.fetch()
# print(a.lastLinks)
