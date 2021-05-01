def decoder(inp: str):
    lines = inp.split("\n")
    videos = []
    for i in range(inp.count("\n")):
        line = lines[i]
        if not line.endswith(".ts"):
            continue
        videos.append(line)
    return videos


def videoLen(inp: str):
    lines = inp.split("\n")
    # #EXTINF:5.080000,
    length = 0.0
    for i in lines:
        if i.startswith("#EXTINF:"):
            i = i.replace("#EXTINF:", "") \
                .replace(",", "")
            length += float(i)
    return length
