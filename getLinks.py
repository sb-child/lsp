import bs4
import re
import requests
import random


class Getter:
    def __init__(self, base_url=""):
        pg = requests.get(f"http://miyajump.xyz/?url=webjump",
                          headers={"User-Agent": "Mozilla/5.0 (X11; Linux x86_64; rv:81.0) Gecko/20100101 Firefox/83.0"}
                          ).text
        url1 = re.search("http://www.*.com/", pg)[0]
        pg = requests.get(f"{url1}?u={random.random()}&path=null",
                          headers={"User-Agent": "Mozilla/5.0 (X11; Linux x86_64; rv:81.0) Gecko/20100101 Firefox/83.0"}
                          ).text
        url_header_index = pg.find("'h'+'t'+'t'+'p'+'s'")
        url_footer_index = pg.find("'5'+'9'+'9'+'8'+'0'") + 19
        auto_url = pg[url_header_index:url_footer_index].replace("'", "").replace("+", "")

        self.base_url = auto_url if base_url == "" else base_url
        print(f"使用网址: {self.base_url}")
        self.play_url_re = re.compile("/index.php/vod/play/id/*")

    def run(self):
        pg = requests.get(self.base_url,
                          headers={"User-Agent": "Mozilla/5.0 (X11; Linux x86_64; rv:81.0) Gecko/20100101 Firefox/83.0"}
                          ).text
        page = bs4.BeautifulSoup(pg, "lxml")
        link_urls = []

        def _get_links(label: str):
            links_part: bs4.element.Tag = page.find("h3", attrs={"class": "title"}, text=label) \
                .find_parent("div") \
                .find_parent("div") \
                .find_parent("div")
            _links = links_part.find_all("a", attrs={"class": "stui-vodlist__thumb"})
            return _links

        links = _get_links("無碼")
        links.extend(_get_links("总热门"))
        for link in links:
            link: bs4.element.Tag
            link_url = link.attrs["href"]
            link_title = link.attrs["title"]
            link_img = link.attrs["data-original"]
            can_dld = True
            for i in link_urls:
                if i[0] == self.base_url + link_url:
                    print(f"重复: {link_title}")
                    can_dld = False
                    break
            if self.play_url_re.match(link_url) and can_dld:
                link_urls.append((self.base_url + link_url, link_title, link_img))
        link_urls = list(set(link_urls))
        return link_urls
