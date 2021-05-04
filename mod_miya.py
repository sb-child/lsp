# Miya: 蜜芽
# 视频拉取模块
import urllib3
import getLinks
import getVideoLink
import decryptLink
import modBase
import tsDecode
import time
from downloader import urlGetToStr


class Puller:  # (modBase.Puller):
    def log(self, message: str):
        print(self.modName + ":", message)

    def __init__(self):
        self.modName = "[Miya]拉取模块"
        self.log("载入模块...")
        self._lg = getLinks.GetterMiya()
        self._dl = decryptLink.Decrypter()
        self._dldPool = urllib3.PoolManager()
        self.lastLinks: list[tuple[str, str, str]] = []
        self.lastTags: list[list] = []
        self.selectedTag = ""
        self.log("模块已准备就绪.")

    def getTags(self):
        self.log("拉取分类列表...")
        # body = getLinks.myReqGet(self._lg.base_url)
        # body_bs = bs4.BeautifulSoup(body, "lxml")
        # tags = body_bs.find_all("a", {"class": "1=0"})
        # tags_list = []
        # for i in range(len(tags)):
        #     tag: bs4.Tag = tags[i]
        #     # /index.php/vod/type/id/33.html
        #     tag_name = tag_re.findall(tag.attrs["href"])
        #     if len(tag_name) == 0:
        #         continue
        #     tags_list.append([tag_name[0], tag.text.strip()])
        #     # ["33", "..."]
        # self.lastTags = tags_list
        self.log(f"分类列表拉取完成, 共[{len(self.lastTags)}]个.")

    def setTag(self, tag_name: str):
        # "33"
        print("暂不支持..")
        exit()
        pass

    def fetch(self):
        self.log("拉取视频列表...")
        self.lastLinks = self._lg.run()
        self.log("视频列表拉取完成.")

    def getDownloadLink(self, index: int):
        self.log("获取下载链接...")
        link_url = self.lastLinks[index]
        print(f"视频{index}:\n" + getLinks.linkFormat(link_url))
        this_link = getVideoLink.getLink(link_url[0])
        this_url = self._dl.decrypt(this_link)
        print(f"* 下载链接:", this_url)
        video_list_str = urlGetToStr(self._dldPool, this_url)
        videos_list = tsDecode.decoder(video_list_str)
        print(f"* 视频时长:", time.strftime("%H:%M:%S", time.gmtime(tsDecode.videoLen(video_list_str))))
        self.log("获取完成.")
        '''
        list:
        [v1.ts, v2.ts, ...]
        links:
        [网页链接, 标题, 封面链接]
        '''
        return {
            "list": videos_list,
            "links": link_url,
        }


if __name__ == '__main__':
    pass
