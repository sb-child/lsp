import pathlib
import time
import mod_miya
import mod_yysp
import downloader
import videoLock
import argparse
import random
import string


def main(selected_mod: str, dld_dir: list, tags: bool, tag: list, no_dld: bool):
    lockFile = "global_lsp.lock"
    lk = {"errors": []}
    restore = False
    if dld_dir is not None and not no_dld:
        lk: dict = videoLock.lockGet(dld_dir[0], fn=lockFile)
        if "videos" in lk:
            print("上次的下载未完成, 将从上次的进度下载:")
            print(f"* 当前进度: 视频{lk['progress']} / {len(lk['videos'])}")
            print(f"* 下载失败: 视频{lk['errors']}")
            restore = True
    if not restore and selected_mod == "":
        raise FileNotFoundError("不需要恢复上次进度, 所以需要指定模块")
    if restore:
        mod = None
    elif selected_mod == "miya":
        mod = mod_miya.Puller()
    elif selected_mod == "yysp":
        mod = mod_yysp.Puller()
    else:
        raise Exception("找不到模块")
    if tags:
        if restore:
            raise FileNotFoundError("当前模式下, 不能获取标签")
        mod.getTags()
        for i, j in mod.lastTags.items():
            print(f"标签编号[{i}], 标签名[{j}]")
        return 0
    if len(tag) != 0:
        if restore:
            raise FileNotFoundError("当前模式下, 不能指定标签")
        mod.setTag(tag)
    videos_dir = ""
    if not no_dld:
        videos_dir = "v_auto_" + \
                     "".join(random.choices(string.hexdigits, k=5)) + \
                     "_" + \
                     str(int(time.time())) \
            if dld_dir is None else dld_dir[0]
        pathlib.Path(videos_dir).mkdir(exist_ok=True)
        print(f"视频将下载到[{videos_dir}]目录")
    if mod is not None:
        mod.fetch()
        link_len = len(mod.lastLinks)
        links = []
        for i in range(link_len):
            link = mod.getDownloadLink(i)
            if "error" in link:
                lk["errors"].append(i)
            links.append(link)
        print("视频列表获取完成, 将下载视频")
    else:
        link_len = len(lk['videos'])
        links: list = lk['videos']
    if not no_dld and not restore:
        videoLock.lockSet(videos_dir, {"videos": links, "progress": 0, "errors": []}, fn=lockFile)
    for i in range(link_len):
        link = links[i]
        if no_dld:
            continue
        if i in lk['errors']:
            print(f"下载时出错过, 跳过视频{i}")
            continue
        if restore and i < lk["progress"]:
            print(f"之前已经下载过, 跳过视频{i}")
            continue
        print(f"下载视频{i}...")
        r = downloader.downloadM3u8(link=link, out_dir=videos_dir, out_file=f"out_{i}", restore=restore)
        lk = videoLock.lockGet(videos_dir, fn=lockFile)
        lk["progress"] = i + 1
        if r != 0:
            lk["errors"].append(i)
        videoLock.lockSet(videos_dir, lk, fn=lockFile)
        # 下次不用恢复进度, 否则会导致临时文件名冲突
        restore = False
    videoLock.lockSet(videos_dir, {}, fn=lockFile)


if __name__ == '__main__':
    parser = argparse.ArgumentParser(description='[lsp]模块执行器')
    parser.add_argument('--mod', dest='mod', action='store', nargs=1,
                        help='指定要加载的模块', type=str)
    parser.add_argument('--dir', dest='dir', action='store', nargs=1,
                        help='指定下载目录', type=str)
    parser.add_argument('--tags', dest='tags', action='store_true',
                        help='获取标签列表')
    parser.add_argument('--tag', dest='tag', action='store', nargs='+',
                        help='指定标签(编号), 否则为默认', type=int)
    parser.add_argument('--not-download', dest='no_dld', action='store_true',
                        help='仅拉取视频列表, 不下载')
    args = parser.parse_args()
    # print(args)
    # exit()
    if args.tag is None:
        args.tag = []
    if args.mod is None:
        args.mod = [""]
    main(args.mod[0], args.dir, args.tags, args.tag, args.no_dld)
