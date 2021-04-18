import getLinks
import getVideoLink
import decryptLink
import time

# 仅测试
ua = ""  # "-user_agent " + '"User-Agent: Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/88.0.4324.182 Safari/537.36"'
headers = ""  # "-headers " + '"sec-ch-ua: \'Chromium\';v=\'88\', \'Google Chrome\';v=\'88\', \';Not A Brand\';v=\'99\'"$\'\r\n\'"sec-ch-ua-mobile: ?0"$"Upgrade-Insecure-Requests: 1"'


def main():
    lg = getLinks.Getter()
    dl = decryptLink.Decrypter()
    link_urls = lg.run()
    f1 = open("dld.sh", "w")
    f2 = open("dld_windows.bat", "w")
    videos_dir = "videos_" + str(int(time.time()))
    print(f"生成的脚本将自动创建 {videos_dir} 目录")
    f1.write(f"mkdir {videos_dir}\n")
    f2.write(f"mkdir {videos_dir}\n")
    count = 0
    link_urls.reverse()
    f1.write(f"# all: {len(link_urls)}\n\n")
    f2.write(f":: all: {len(link_urls)}\n\n")
    for link_url in link_urls:
        print(f"视频{count}:", link_url)
        this_link = getVideoLink.getLink(link_url[0])
        # print(this_link)
        this_url = dl.decrypt(this_link)
        print(f"视频{count}的链接:", this_url)
        f1.write(f"# {count}: {link_url[1]}\n")
        f2.write(f":: {count}: {link_url[1]}\n")
        f1.write(f"# img: {link_url[2]}\n")
        f2.write(f":: img: {link_url[2]}\n")
        url_fix = this_url.replace('&', '\\&')
        f1.write(f"ffmpeg {ua} {headers} -i {url_fix} -c copy {videos_dir}/{count}.mp4\n")
        f2.write(f"ffmpeg {ua} {headers} -i {this_url} -c copy {videos_dir}/{count}.mp4\n")
        f1.write("\n")
        f2.write("\n")
        # ffmpeg -i _url -c copy _file
        count += 1
    f1.close()
    f2.close()


if __name__ == '__main__':
    main()
