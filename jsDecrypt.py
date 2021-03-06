import pathlib
import os
import subprocess
import base64

jsdec_dir = "jsdec-tiny"


def checkJsDec():
    dirs = pathlib.Path(jsdec_dir)
    if not dirs.exists():
        print(f"jsDecrypt: 找不到 {jsdec_dir} 目录")
        return 1
    if not (dirs / "main").exists() and not (dirs / "main.exe").exists():
        print(f"jsDecrypt: 找不到可执行文件")
        return 2
    return 0


def dec(inp: str) -> str:
    file = f"./main" if os.name != "nt" else f"main.exe"
    r = subprocess.Popen([file], cwd=jsdec_dir,
                         stdout=subprocess.PIPE, stdin=subprocess.PIPE, stderr=subprocess.STDOUT)
    out = r.communicate(base64.b64encode(inp.encode()) + b"\n")[0].decode()
    if r.returncode != 0:
        print(f"jsDecrypt: 解密失败: {out}")
        return ""
    return out
