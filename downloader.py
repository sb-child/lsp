# 下载模块
import os
import subprocess
import uuid
import retrying
import requests
from typing import Union
from tqdm import tqdm
from Crypto.Cipher import AES


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


# def downloadM3u8(link: dict[str, Union[str, list[str], tuple[str, str, str], float]],
def downloadM3u8(link: dict,
                 out_dir: str, out_file: str):
    videos_list = link["list"]
    link_url = link["links"]
    video_encrypt: str = link["encrypt"]
    videos_list_len = len(videos_list)
    uid = uuid.uuid4().__str__()

    # debug: encrypt
    if video_encrypt == "":
        return 3

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
        if video_encrypt != "":
            keyStr = urlGetToStr(video_encrypt)
            keyBin = keyStr.encode()
            aesDec = AES.new(keyBin, AES.MODE_CBC, keyBin)
            data = aesDec.decrypt(data)
        with open(fn, "ab+") as f:
            f.write(data)
        os.remove(fn2)

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
