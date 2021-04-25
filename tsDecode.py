def decoder(inp: str):
    lines = inp.split("\n")
    videos = []
    for i in range(inp.count("\n")):
        line = lines[i]
        if not line.endswith(".ts"):
            continue
        # if line is None:
        #     continue
        videos.append(line)
    return videos
