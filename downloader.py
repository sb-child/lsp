# 下载模块
import os
import shutil
import subprocess
import uuid
import retrying
import requests
from typing import Union
from tqdm import tqdm


def urlGet(url: str):
    try:
        req = requests.get(url, timeout=5)
        if req.status_code != 200:
            raise ConnectionError("请求返回值 != 200")
    except Exception as e:
        print(f"* 下载[{url}]时抛出异常:", e)
        raise e
    return req.content


def urlGetToBinFile(url, fn: str):
    with open(fn, "wb") as f:
        f.write(urlGet(url))


def urlGetToStr(url, encoding="utf-8"):
    return urlGet(url).decode(encoding)


def on_err(attempts, delay):
    print("下载失败, 稍后重试...")


@retrying.retry(stop_max_attempt_number=10, wait_random_min=5000,
                wait_random_max=10000, wait_incrementing_increment=0, stop_func=on_err)
def downloadVideoPart(dld_url: str, filename: str):
    urlGetToBinFile(dld_url, filename)


def downloadM3u8(link: dict[str, Union[str, list[str], tuple[str, str, str], float]],
                 out_dir: str, out_file: str):
    videos_list = link["list"]
    link_url = link["links"]
    video_encrypt: str = link["encrypt"]
    video_len: float = link["len"]
    videos_list_len = len(videos_list)
    uid = uuid.uuid4().__str__()
    for i in tqdm(range(videos_list_len), desc="下载视频"):
        dld_url = videos_list[i]
        fn = os.path.join(out_dir, f"t_{uid}_{i}.ts")
        try:
            downloadVideoPart(dld_url=dld_url, filename=fn)
        except Exception:
            raise ConnectionError("多次下载失败, 放弃.")
    fn = os.path.join(out_dir, f"t2_{uid}.ts")
    for i in tqdm(range(videos_list_len), desc="合并视频"):
        fn2 = os.path.join(out_dir, f"t_{uid}_{i}.ts")
        with open(fn2, "rb") as f:
            data = f.read()
        with open(fn, "ab+") as f:
            f.write(data)
        os.remove(fn2)
    if video_encrypt != "":
        print("解密视频, 这通常需要很长时间...")
        fn_key = os.path.join(out_dir, f"t_key_{uid}.key")
        fn_dec_m3u8 = os.path.join(out_dir, f"t_dec_{uid}.m3u8")
        fn_dec_out = os.path.join(out_dir, f"t_dec_out_{uid}.ts")
        urlGetToBinFile(video_encrypt, fn_key)
        with open(fn_dec_m3u8, "w") as f:
            f.write(f"""
            #EXTM3U
            #EXT-X-VERSION:3
            #EXT-X-MEDIA-SEQUENCE:0
            #EXT-X-KEY:METHOD=AES-128,URI="{fn_key}"
            #EXTINF:{'{:.6f}'.format(video_len)},
            {fn}
            """.strip())
        if subprocess.run(f"ffmpeg -allowed_extensions ALL "
                          f"-v 0 -y -i {fn_dec_m3u8} -c copy "
                          f"{fn_dec_out}", shell=True).returncode != 0:
            print("解密时出错.")
            return 2
        os.remove(fn_key)
        os.remove(fn)
        shutil.move(fn_dec_out, fn)
        print("解密完成.")

    print("转换为mp4格式...")
    fn3 = os.path.join(out_dir, f"{out_file}.mp4")
    if subprocess.run(f"ffmpeg -v 0 -y -i {fn} -c copy {fn3}", shell=True).returncode != 0:
        print("格式转换时出错.")
        return 1
    os.remove(fn)
    print("下载封面...")
    fn_img = os.path.join(out_dir, f"{out_file}.jpg")
    urlGetToBinFile(link_url[2], fn_img)
    print("写出描述文件...")
    fn_desc = os.path.join(out_dir, f"{out_file}.txt")
    with open(fn_desc, "w") as f:
        f.write("\n".join(link_url))
    print("完成.")
    pass
