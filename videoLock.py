import hashlib
import sqlite3


class Lock:
    def __init__(self, lockfile=""):
        self.lockfile = lockfile
        self.sql = sqlite3.connect(lockfile)

    def addVideoLock(self, name: str):
        # todo
        self.sql.cursor()
        hashlib.sha256()
