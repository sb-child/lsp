import bs4
import re

from getLinks import myReqGet

link_re = re.compile(r'"url":"(.*?)",')


def getLink(url: str):
    # print(url)
    pg = myReqGet(url)
    page = bs4.BeautifulSoup(pg, "lxml")
    for i in page.find_all("script"):
        i: bs4.element.Tag
        it = i.string
        if it is None:
            continue
        if "encrypt" not in it:
            continue
        it1 = link_re.findall(it)
        if len(it1) == 0:
            continue
        res: str = it1[0]
        return res.replace("\\", "")
    return ""
