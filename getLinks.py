import bs4
import re
import requests


class Getter:
    def __init__(self, base_url="https://www.myvzao0ioh7xgjvdon7f8jqbqmef.xyz:59980"):
        self.base_url = base_url
        # pg = open("test1.html").read()
        self.play_url_re = re.compile("/index.php/vod/play/id/*")

    def run(self):
        pg = requests.get(self.base_url,
                          headers={"User-Agent": "Mozilla/5.0 (X11; Linux x86_64; rv:81.0) Gecko/20100101 Firefox/83.0"}
                          ).text
        page = bs4.BeautifulSoup(pg, "lxml")
        links_part: bs4.element.Tag = page.find("h3", attrs={"class": "title"}, text="最新更新")\
            .find_parent("div")\
            .find_parent("div")\
            .find_parent("div")
        links = links_part.find_all("a")
        link_urls = []
        for link in links:
            link: bs4.element.Tag
            link_url = link.attrs["href"]
            link_title = link.attrs["title"]
            if self.play_url_re.match(link_url):
                link_urls.append((self.base_url + link_url, link_title))
        link_urls = list(set(link_urls))
        return link_urls
