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


def on_err(attempts, delay):
    print("下载失败, 稍后重试...")


@retrying.retry(stop_max_attempt_number=10, wait_random_min=5000,
                wait_random_max=10000, wait_incrementing_increment=0, stop_func=on_err)
def downloadVideoPart(dld_url: str, pool: urllib3.poolmanager.PoolManager, filename: str):
    req: urllib3.response.HTTPResponse = pool.request("GET", dld_url)
    with open(filename, "wb") as f:
        f.write(req.data)
    req.close()


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
        print(f"视频{count}:", link_url)
        this_link = getVideoLink.getLink(link_url[0])
        this_url = dl.decrypt(this_link)
        print(f"视频{count}的链接:", this_url)
        req: urllib3.response.HTTPResponse = dldPool.request("GET", this_url)
        video_list_str = req.data.decode("utf-8")
        videos_list = tsDecode.decoder(video_list_str)  # [:10]
        videos_list_len = len(videos_list)
        req.close()
        skip = False
        for i in progressbar(range(videos_list_len), desc="下载视频"):
            dld_url = videos_list[i]
            fn = os.path.join(videos_dir, f"temp_{count}_{i}.ts")
            try:
                downloadVideoPart(dld_url, dldPool, fn)
            except Exception:
                skip = True
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
        fn3 = os.path.join(videos_dir, f"out_{count}.mp4")
        fn_desc = os.path.join(videos_dir, f"out_{count}.txt")
        print("格式转换...")
        if subprocess.run(f"ffmpeg -v 0 -y -i {fn} -c copy {fn3}", shell=True).returncode != 0:
            print("格式转换时出错.")
            return 1
        os.remove(fn)
        print("写出描述文件...")
        with open(fn_desc, "w") as f:
            f.write("\n".join(link_url))
        print("完成.")
        count += 1
        # return


if __name__ == '__main__':
    main()
