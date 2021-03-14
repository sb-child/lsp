import bs4
import re
import requests

base_url = "https://www.myvzao0ioh7xgjvdon7f8jqbqmef.xyz:59980"
# pg = open("test1.html").read()
pg = requests.get(base_url).text
play_url_re = re.compile("/index.php/vod/play/id/*")


def main():
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
        if play_url_re.match(link_url):
            link_urls.append((base_url + link_url, link_title))
    link_urls = list(set(link_urls))
    for link_url in link_urls:
        print(link_url)


if __name__ == '__main__':
    main()
