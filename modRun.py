import os
import time
import mod_miya
import downloader


def main():
    videos_dir = "videos_auto_" + str(int(time.time()))
    os.mkdir(videos_dir)
    mod = mod_miya.Puller()
    mod.fetch()
    for i in range(len(mod.lastLinks)):
        link = mod.getDownloadLink(i)
        downloader.downloadM3u8(link=link, out_dir=videos_dir, out_file=f"out_{i}")


if __name__ == '__main__':
    main()
