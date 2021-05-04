from getLinks import myReqGet


def preProcess(func):
    def checkLink(inp: str):
        lines = inp.split("\n")
        for i in range(inp.count("\n")):
            line = lines[i]
            if line.endswith(".m3u8"):
                # m3u8内嵌m3u8
                return line
        return ""

    def get(inp: str, source=""):
        r = checkLink(inp)
        if r == "":
            # 不包含嵌套, 返回.
            return inp, source
        # 补全主机名
        r = r if r.startswith("http") else f"{source}{r}"
        return get(myReqGet(r), source)

    def wrapped_f(inp: str, source=""):
        r = get(inp, source)
        return func(*r)

    return wrapped_f


@preProcess
def decoder(inp: str, source=""):
    lines = inp.split("\n")
    videos = []
    for i in range(inp.count("\n")):
        line = lines[i]
        if not line.endswith(".ts"):
            continue
        videos.append(line)
    return videos


@preProcess
def videoLen(inp: str, source=""):
    lines = inp.split("\n")
    # #EXTINF:5.080000,
    length = 0.0
    for i in lines:
        if i.startswith("#EXTINF:"):
            i = i.replace("#EXTINF:", "") \
                .replace(",", "")
            length += float(i)
    return length


# print(decoder(myReqGet("https://lbbf9.com/20200428/6s2AAeKd/index.m3u8"), source="https://lbbf9.com"))
