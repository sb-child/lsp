import urllib3


def urlGet(pool: urllib3.poolmanager.PoolManager, url: str):
    req: urllib3.response.HTTPResponse = pool.request("GET", url)
    return req.data


def urlGetToBinFile(pool, url, fn: str):
    with open(fn, "wb") as f:
        f.write(urlGet(pool, url))


def urlGetToStr(pool, url, encoding="utf-8"):
    return urlGet(pool, url).decode(encoding)
