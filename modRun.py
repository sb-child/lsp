import os
import time
import mod_miya
import mod_yysp
import downloader
import argparse


def main(selected_mod: str, dld_dir: list, tags: bool, tag: str, no_dld: bool):
    if selected_mod == "miya":
        mod = mod_miya.Puller()
    elif selected_mod == "yysp":
        mod = mod_yysp.Puller()
    else:
        raise Exception("找不到模块")
    if tags:
        mod.getTags()
        for i in range(len(mod.lastTags)):
            j = mod.lastTags[i]
            print(f"标签编号[{j[0]}], 标签名[{j[1]}]")
        return 0
    if tag is not None:
        mod.setTag(tag)
    videos_dir = ""
    if not no_dld:
        videos_dir = "videos_auto_" + str(int(time.time())) if dld_dir is None else dld_dir[0]
        os.mkdir(videos_dir)
        print(f"视频将下载到[{videos_dir}]目录")
    mod.fetch()
    link_len = len(mod.lastLinks)
    print(f"获取到{link_len}个视频")
    for i in range(link_len):
        link = mod.getDownloadLink(i)
        if no_dld:
            continue
        downloader.downloadM3u8(link=link, out_dir=videos_dir, out_file=f"out_{i}")


if __name__ == '__main__':
    parser = argparse.ArgumentParser(description='[lsp]模块执行器')
    parser.add_argument('--mod', dest='mod', action='store', nargs=1,
                        help='指定要加载的模块', type=str, required=True)
    parser.add_argument('--dir', dest='dir', action='store', nargs=1,
                        help='指定下载目录', type=str)
    parser.add_argument('--tags', dest='tags', action='store_true',
                        help='获取标签列表')
    parser.add_argument('--tag', dest='tag', action='store', nargs=1,
                        help='指定标签(编号), 否则为默认', type=str)
    parser.add_argument('--not-download', dest='no_dld', action='store_true',
                        help='仅拉取视频列表, 不下载')
    args = parser.parse_args()
    # print(args)
    # exit()
    main(args.mod[0], args.dir, args.tags, args.tag[0], args.no_dld)
