import hashlib
import json
import os
import sqlite3


def lockInit(d: str, fn="lsp.lock"):
    lockFileName = os.path.join(d, fn)
    try:
        os.mkdir(d)
    except FileExistsError:
        pass
    if fn not in os.listdir(d):
        with open(lockFileName, "a+") as f:
            # f.seek(0)
            # f.truncate()
            f.write("{}")
    return lockFileName


def lockGet(d: str, fn="lsp.lock"):
    lockFileName = lockInit(d, fn)
    with open(lockFileName, "r") as f:
        r = json.load(f)
    return dict(r)


def lockSet(d: str, data: dict, fn="lsp.lock"):
    lockFileName = lockInit(d, fn)
    with open(lockFileName, "w") as f:
        json.dump(data, f)


class Lock:
    def __init__(self, lockfile=""):
        self.lockfile = lockfile
        self.sql = sqlite3.connect(lockfile)

    def addVideoLock(self, name: str):
        # todo
        self.sql.cursor()
        hashlib.sha256()
