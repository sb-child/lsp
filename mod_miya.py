# Miya: 蜜芽
# 视频拉取模块

import bs4

import colors
import getLinks
import getVideoLink
import decryptLink
import modBase
import tsDecode
import time

from typing import Union
from downloader import urlGetToStr

tag_re = modBase.tag_re


class Puller(modBase.Puller):
    def _init(self):
        self._name("Miya")
        self._lg = getLinks.GetterMiya()
        self._dl = decryptLink.Decrypter()
        self._tagName = ""

    def _getTags(self):
        body = getLinks.myReqGet(self._lg.base_url)
        body_bs = bs4.BeautifulSoup(body, "lxml")
        tags_root = body_bs.find("ul", {"class": "stui-header__menu type-slide"})
        tags = tags_root.find_all("a")
        tags_list = {}
        for i in range(len(tags)):
            tag: bs4.Tag = tags[i]
            tag_name = tag_re.findall(tag.attrs["href"])
            if len(tag_name) == 0:
                continue
            tags_list[tag_name[0]] = tag.text.strip()
        self.lastTags = tags_list

    def _setTag(self, tag_name: list):
        self.selectedTag = tag_name
        self.log("获取标签列表以指定标签...")
        self._getTags()
        self._tagName = [self.lastTags[str(i)] for i in tag_name]

    def _fetch(self):
        if len(self._tagName) == 0:
            self.lastLinks = self._lg.run()
            return
        temp = []
        for i in self._tagName:
            temp.extend(self._lg.run(tag=str(i)))
        self.lastLinks = temp

    def _getDownloadLink(self, index: int):
        link_url = self.lastLinks[index]
        print(f"视频[{colors.f_important(index)}]:\n" + getLinks.linkFormat(link_url))
        this_link = getVideoLink.getLink(link_url[0])
        this_url = self._dl.decrypt(this_link)
        print(f"* 下载链接:", colors.f_important(this_url))
        video_list_str = urlGetToStr(this_url)
        videos_list, _ = tsDecode.decoder(video_list_str)
        video_len: float = tsDecode.videoLen(video_list_str)
        print(f"* 视频时长:", colors.f_important(time.strftime("%H:%M:%S", time.gmtime(video_len))))
        return {
            "list": videos_list,
            "links": link_url,
            "encrypt": "",
        }


if __name__ == '__main__':
    pass
