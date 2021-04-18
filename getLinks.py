import bs4
import re
import requests
import random


class Getter:
    # 若脚本出错, 请更新此处的base_url. 在实现自动获取链接的功能之前, 欢迎贡献最新的链接
    def __init__(self, base_url=""):
        # http://www.my1132.com/?u=0.844066510415847&path=null
        # var _0x7e7501='h'+'t'+'t'+'p'+'s' ... '5'+'9'+'9'+'8'+'0';
        pg = requests.get(f"http://www.my1132.com/?u={random.random()}&path=null",
                          headers={"User-Agent": "Mozilla/5.0 (X11; Linux x86_64; rv:81.0) Gecko/20100101 Firefox/83.0"}
                          ).text
        url_header_index = pg.find("'h'+'t'+'t'+'p'+'s'")
        url_footer_index = pg.find("'5'+'9'+'9'+'8'+'0'") + 19
        auto_url = pg[url_header_index:url_footer_index].replace("'", "").replace("+", "")  # + "/"

        self.base_url = auto_url if base_url == "" else base_url
        print(f"使用网址: {self.base_url}")
        # pg = open("test1.html").read()
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

        # 最新更新
        links = _get_links("無碼")
        links.extend(_get_links("总热门"))
        for link in links:
            link: bs4.element.Tag
            # print(link)
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
