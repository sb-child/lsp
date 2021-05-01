import urllib3
import getLinks
import getVideoLink
import decryptLink
import tsDecode
import time
import os
import tqdm
import subprocess
import retrying


def urlGet(pool: urllib3.poolmanager.PoolManager, url: str):
    req: urllib3.response.HTTPResponse = pool.request("GET", url)
    return req.data


def urlGetToBinFile(pool, url, fn: str):
    with open(fn, "wb") as f:
        f.write(urlGet(pool, url))


def urlGetToStr(pool, url, encoding="utf-8"):
    return urlGet(pool, url).decode(encoding)


def on_err(attempts, delay):
    print("下载失败, 稍后重试...")


@retrying.retry(stop_max_attempt_number=10, wait_random_min=5000,
                wait_random_max=10000, wait_incrementing_increment=0, stop_func=on_err)
def downloadVideoPart(dld_url: str, pool: urllib3.poolmanager.PoolManager, filename: str):
    urlGetToBinFile(pool, dld_url, filename)


def linkFormat(link: tuple[str, str, str]):
    return f"* 网页链接: {link[0]}\n" \
           f"* 标题: {link[1]}\n" \
           f"* 封面链接: {link[2]}"


def main():
    lg = getLinks.Getter()
    dl = decryptLink.Decrypter()
    link_urls = lg.run()
    dldPool = urllib3.PoolManager()
    progressbar = tqdm.tqdm
    videos_dir = "videos_auto_" + str(int(time.time()))
    print(f"脚本将自动创建 {videos_dir} 目录")
    os.mkdir(videos_dir)
    count = 0
    link_urls.reverse()
    for link_url in link_urls:
        print(f"视频{count}:\n" + linkFormat(link_url))
        this_link = getVideoLink.getLink(link_url[0])
        this_url = dl.decrypt(this_link)
        print(f"* 下载链接:", this_url)
        video_list_str = urlGetToStr(dldPool, this_url)
        videos_list = tsDecode.decoder(video_list_str)
        print(f"* 视频时长:", time.strftime("%H:%M:%S", time.gmtime(tsDecode.videoLen(video_list_str))))
        videos_list_len = len(videos_list)
        skip = False
        for i in progressbar(range(videos_list_len), desc="下载视频"):
            dld_url = videos_list[i]
            fn = os.path.join(videos_dir, f"temp_{count}_{i}.ts")
            try:
                downloadVideoPart(dld_url, dldPool, fn)
            except Exception:
                skip = True
                break
        if skip:
            print("多次下载失败, 跳过.")
            break
        fn = os.path.join(videos_dir, f"temp2_{count}.ts")
        for i in progressbar(range(videos_list_len), desc="合并视频"):
            fn2 = os.path.join(videos_dir, f"temp_{count}_{i}.ts")
            with open(fn2, "rb") as f:
                data = f.read()
            with open(fn, "ab+") as f:
                f.write(data)
            os.remove(fn2)
        print("转换为mp4格式...")
        fn3 = os.path.join(videos_dir, f"out_{count}.mp4")
        if subprocess.run(f"ffmpeg -v 0 -y -i {fn} -c copy {fn3}", shell=True).returncode != 0:
            print("格式转换时出错.")
            return 1
        os.remove(fn)
        print("下载封面...")
        fn_img = os.path.join(videos_dir, f"out_{count}_img.jpg")
        urlGetToBinFile(dldPool, link_url[2], fn_img)
        print("写出描述文件...")
        fn_desc = os.path.join(videos_dir, f"out_{count}.txt")
        with open(fn_desc, "w") as f:
            f.write("\n".join(link_url))
        print("完成.")
        count += 1
        # return


if __name__ == '__main__':
    main()
