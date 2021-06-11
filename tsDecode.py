import re

from getLinks import myReqGet
from modBase import domain_re

uri_re = re.compile(r'URI="(.*?)"')


def preProcess(func):
    def checkLink(inp: str):
        lines = inp.split("\n")
        for i in range(inp.count("\n")):
            line = lines[i]
            if line.endswith(".m3u8"):
                # m3u8内嵌m3u8
                return line
        return ""

    def get(inp: str, source=""):
        r = checkLink(inp)
        if r == "":
            # 不包含嵌套, 返回.
            return inp, source
        # 补全主机名
        src_domain = domain_re.findall(source)[0]
        r = r if len(domain_re.findall(r)) != 0 else f"{src_domain}{r}"
        return get(myReqGet(r), source)

    def wrapped_f(inp: str, source=""):
        r = get(inp, source)
        return func(*r)

    return wrapped_f


def _mergeUrl(line: str, source: str):
    if len(domain_re.findall(line)) == 0 and line.find("/") == -1:
        # xx.ts
        return source + line
    elif len(domain_re.findall(line)) == 0:
        # xx/xx/xx.ts
        src_domain = domain_re.findall(source)[0]
        return src_domain + line
    else:
        # https://xx.xx/xx/xx.ts
        return line


@preProcess
def decoder(inp: str, source=""):
    # print(inp)
    lines = inp.split("\n")
    videos = []
    for i in range(inp.count("\n")):
        line = lines[i]
        if not line.endswith(".ts"):
            continue
        videos.append(_mergeUrl(line, source))
    # print(videos)
    return videos, source


@preProcess
def videoLen(inp: str, source=""):
    lines = inp.split("\n")
    # #EXTINF:5.080000,
    length = 0.0
    for i in lines:
        if i.startswith("#EXTINF:"):
            i = i.replace("#EXTINF:", "") \
                .replace(",", "")
            length += float(i)
    return length


@preProcess
def checkEncrypt(inp: str, source=""):
    # #EXT-X-KEY:METHOD=AES-128,URI="https://cdn1.baodelaike.com:4343/20210310/Cs5xTxiS/1000kb/hls/key.key"
    lines = inp.split("\n")
    for i in lines:
        if i.startswith("#EXT-X-KEY:"):
            uri: str = uri_re.findall(i)[0]
            return _mergeUrl(uri, source)
    return ""

# print(decoder(myReqGet("https://lbbf9.com/20200428/6s2AAeKd/index.m3u8"), source="https://lbbf9.com"))
