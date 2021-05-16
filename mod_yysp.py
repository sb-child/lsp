# yysp: 夜夜視頻資源站
# 视频拉取模块
import time
import bs4
import getLinks
import getVideoLink
import tsDecode
import modBase
from downloader import urlGetToStr

domain_re = modBase.domain_re
link_re = modBase.link_re
tag_re = modBase.tag_re


def _decoder(url: str):
    # 多视频来源解码
    new_domain: str = domain_re.findall(url)[0]
    if url.endswith(".m3u8"):
        # 不需要处理
        return url, new_domain
    result = getLinks.myReqGet(url)
    # https://vip4.ddyunbo.com
    new_dir = link_re.findall(result)[0]
    new_url = new_domain + new_dir
    # var playlist = '[{"url":"/20210501/ItF819cE/600kb/hls/index.m3u8"}]';
    return new_url, new_domain


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

    def _setTag(self, tag_name: str):
        self.selectedTag = tag_name

    def _fetch(self):
        self.lastLinks = self._lg.run(tag=self.selectedTag)

    def _getDownloadLink(self, index: int):
        link_url = self.lastLinks[index]
        print(f"视频{index}:\n" + getLinks.linkFormat(link_url))
        this_link = getVideoLink.getLink(link_url[0])
        this_url, domain = _decoder(this_link)
        this_url: str
        domain: str
        print(f"* 下载链接:", this_url)
        video_list_str = urlGetToStr(this_url)
        _decode_url = domain
        videos_list = tsDecode.decoder(video_list_str, _decode_url)
        for i in range(len(videos_list)):
            if len(domain_re.findall(videos_list[i])) == 0:
                videos_list[i] = domain + videos_list[i]
        video_len: float = tsDecode.videoLen(video_list_str, _decode_url)
        print(f"* 视频时长:", time.strftime("%H:%M:%S", time.gmtime(video_len)))
        encrypt = tsDecode.checkEncrypt(video_list_str, _decode_url)
        print(f"* 密钥: [{encrypt if encrypt != '' else '无需解密'}]")
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
