import getLinks
import getVideoLink
import decryptLink


def main():
    lg = getLinks.Getter()
    dl = decryptLink.Decrypter()
    link_urls = lg.run()
    f1 = open("dld.sh", "w")
    count = 0
    link_urls.reverse()
    for link_url in link_urls:
        print(link_url)
        this_link = getVideoLink.getLink(link_url[0])
        print(this_link)
        this_url = dl.decrypt(this_link)
        print(this_url)
        f1.write(f"# {count}: {link_url[1]}\n")
        url_fix = this_url.replace('&', '\\&')
        f1.write(f"ffmpeg -i {url_fix} -c copy videos/{count}.mp4\n")
        # ffmpeg -i _url -c copy _file
        count += 1
    f1.close()


if __name__ == '__main__':
    main()
