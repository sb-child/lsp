from typing import Union

import bs4
import re
import requests
import random


def myReqGet(url: str):
    return requests.get(url,
                        headers={"User-Agent": "Mozilla/5.0 (X11; Linux x86_64; rv:81.0) Gecko/20100101 Firefox/83.0"}
                        ).text


# def linkFormat(link: tuple[str, str, str]):
def linkFormat(link: Union[list, tuple]):
    return f"* 网页链接: {link[0]}\n" \
           f"* 标题: {link[1]}\n" \
           f"* 封面链接: {link[2]}"


class GetterMiya:
    def __init__(self):
        # print("-> 访问[miyajump.xyz]以获取链接...")
        # pg = myReqGet("http://miyajump.xyz/?url=webjump")
        # url1 = re.search("http://www.*.com/", pg)[0]
        url1 = "http://www.my1152.com/"
        print(f"-> 得到跳转页链接[{url1}].正在获取主页链接...")
        pg = myReqGet(f"{url1}?u={random.random()}&path=null")
        url_header_index = pg.find("'h'+'t'+'t'+'p'+'s'")
        url_footer_index = pg.find("'5'+'9'+'9'+'8'+'0'") + 19
        auto_url = pg[url_header_index:url_footer_index].replace("'", "").replace("+", "")
        print(f"-> 得到主页链接[{auto_url}].")
        self.base_url = auto_url
        self.play_url_re = re.compile("/index.php/vod/play/id/*")

    def run(self, tag=""):
        url = self.base_url
        # url = url if tag == "" else f"{url}/index.php/vod/type/id/{tag}.html"
        pg = myReqGet(url)
        page = bs4.BeautifulSoup(pg, "lxml")
        link_urls = []

        def _get_links(label: str):
            links_part: bs4.element.Tag = page.find("h3", attrs={"class": "title"}, text=label) \
                .find_parent("div") \
                .find_parent("div") \
                .find_parent("div")
            _links = links_part.find_all("a", attrs={"class": "stui-vodlist__thumb"})
            return _links

        # links = _get_links("無碼")
        # links.extend(_get_links("总热门"))
        links = _get_links("总热门" if tag == "" else tag)
        for link in links:
            link: bs4.element.Tag
            link_url = link.attrs["href"]
            link_title = link.attrs["title"]
            link_img = link.attrs["data-original"]
            can_dld = True
            for i in link_urls:
                if i[0] == self.base_url + link_url:
                    print(f"-> 重复: {link_title}")
                    can_dld = False
                    break
            if self.play_url_re.match(link_url) and can_dld:
                link_urls.append((self.base_url + link_url, link_title, link_img))
        link_urls = list(set(link_urls))
        return link_urls


class GetterYysp:
    def __init__(self):
        print("-> 遍历已知域名以确定主页...")
        domains = [
            "https://www.yyspzy5.xyz",
            "https://www.yyspzy4.xyz",
            "https://www.yyspzy1.xyz",
            "https://www.yyspzy2.xyz",
            "https://www.yyspzy3.xyz",
        ]
        selected_domain = ""
        for i in domains:
            print(f"-> 访问网址[{i}]...", end="")
            pg = myReqGet(i)
            page = bs4.BeautifulSoup(pg, "lxml")
            refresh_ele = page.find("meta", {"http-equiv": "refresh"})
            if refresh_ele is None:
                # 没有跳转, 即为主页
                selected_domain = i
                print("OK")
                break
            else:
                refresh_url: str = refresh_ele.attrs["content"]
                # 0.1;URL=http://www.yyspzy5.xyz
                selected_domain = refresh_url[8:]
                selected_domain = selected_domain.replace("http", "https")
                print(f"跳转到[{selected_domain}]")
                break
        if selected_domain == "":
            print("-> 找不到可用的域名.")
            raise ConnectionError("找不到可用的域名")
        print(f"-> 得到主页链接[{selected_domain}].")
        self.base_url = selected_domain

    def run(self, tag=""):
        url = self.base_url
        url = url if tag == "" else f"{url}/index.php/vod/type/id/{tag}.html"
        pg = myReqGet(url)
        page = bs4.BeautifulSoup(pg, "lxml")
        link_urls = []
        new_videos_root: bs4.element.Tag = page.find("div", attrs={"class": "pic"}).findChild("ul")
        new_videos: bs4.element.ResultSet = new_videos_root.find_all("a")
        for link in new_videos:
            # /index.php/vod/play/id/580061/sid/1/nid/1.html
            # /index.php/vod/detail/id/580061.html
            video_url: str = str(link.attrs["href"]) \
                .replace("detail", "play") \
                .replace(".html", "/sid/1/nid/1.html")
            title: str = link.attrs["title"]
            img: str = link.find("img").attrs["src"]
            link_urls.append((self.base_url + video_url, title, self.base_url + img))
        link_urls = list(set(link_urls))
        return link_urls
