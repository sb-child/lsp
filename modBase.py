import re
from typing import Union

domain_re = re.compile(r'(http[s]?://.*?)/')
link_re = re.compile(r'"url":"(.*?)"')
link2_re = re.compile(r'm3u8url = \'(.*?)\'')
tag_re = re.compile(r'/index\.php/vod/type/id/(.*?).html')


class Puller:
    def log(self, message: str):
        print(self.modName + ":", message)

    def _name(self, mod_name=""):
        self.modName = f"[{mod_name}]拉取模块"

    def _init(self):
        pass

    def __init__(self):
        self.modName = ""
        self._init()
        # self.log("载入模块...")
        # self.lastLinks: list[tuple[str, str, str]] = []
        # self.lastTags: dict[str, str] = {}
        self.lastLinks = []
        self.lastTags = {}
        self.selectedTag = []
        self.log("模块已准备就绪.")

    def _getTags(self):
        pass

    def getTags(self):
        self.log("拉取分类列表...")
        self._getTags()
        self.log(f"分类列表拉取完成, 共[{len(self.lastTags)}]个.")

    def _setTag(self, tag_name: list):
        pass

    def setTag(self, tag_name: Union[list, int]):
        if isinstance(tag_name, int):
            tag_name = [tag_name]
        self._setTag(tag_name=tag_name)
        if len(tag_name) == 0:
            self.log(f"切换到默认标签.")
        self.log(f"指定标签[{','.join([str(i) for i in tag_name])}].")

    def _fetch(self):
        pass

    def fetch(self):
        self.log("拉取视频列表...")
        self._fetch()
        self.log(f"视频列表拉取完成, 共[{len(self.lastLinks)}]个.")

    # def _getDownloadLink(self, index: int) -> dict[str, tuple[str, str, str]]:
    def _getDownloadLink(self, index: int) -> dict:
        pass

    def getDownloadLink(self, index: int):
        self.log("获取下载链接...")
        r = self._getDownloadLink(index=index)
        if "error" in r:
            print("获取链接时出错:", r["error"])
            return r
        self.log("获取完成.")
        '''
        list:
        [v1.ts, v2.ts, ...]
        links:
        [网页链接, 标题, 封面链接]
        '''
        return r
