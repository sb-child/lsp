import os
import time
import mod_miya
import mod_yysp
import downloader
import argparse


def main(selected_mod: str, dld_dir: list):
    if selected_mod == "miya":
        mod = mod_miya.Puller()
    elif selected_mod == "yysp":
        mod = mod_yysp.Puller()
    else:
        raise Exception("找不到模块")
    videos_dir = "videos_auto_" + str(int(time.time())) if dld_dir is None else dld_dir[0]
    os.mkdir(videos_dir)
    print(f"视频将下载到[{videos_dir}]目录")
    mod.fetch()
    for i in range(len(mod.lastLinks)):
        link = mod.getDownloadLink(i)
        downloader.downloadM3u8(link=link, out_dir=videos_dir, out_file=f"out_{i}")


if __name__ == '__main__':
    parser = argparse.ArgumentParser(description='[lsp]模块执行器')
    parser.add_argument('--mod', dest='mod', action='store', nargs=1,
                        help='指定要加载的模块', type=str)
    parser.add_argument('--dir', dest='dir', action='store', nargs=1,
                        help='指定下载目录', type=str)
    args = parser.parse_args()
    # print(args)
    main(str(args.mod[0]), args.dir)
