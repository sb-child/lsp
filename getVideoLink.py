import bs4
import re
import requests

link_re = re.compile(r'"url":"(.*?)",')


def getLink(url: str):
    pg = requests.get(url,
                      headers={"User-Agent": "Mozilla/5.0 (X11; Linux x86_64; rv:81.0) Gecko/20100101 Firefox/83.0"}
                      ).text
    page = bs4.BeautifulSoup(pg, "lxml")
    for i in page.find_all("script"):
        i: bs4.element.Tag
        it = i.string
        if it is None:
            continue
        it1 = link_re.findall(it)
        if len(it1) == 0:
            continue
        if it1[0] == "MiYa188.com":
            continue
        res: str = it1[0]
        return res.replace("\\", "")
    return ""

