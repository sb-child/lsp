from Crypto.Cipher import AES
import base64


class Decrypter:
    def __init__(self):
        self.key = b"9q4h7kt7skwsc9af1qmwy14jkfq2biab"
        self.iv = b"6b3gslw69k6eazmw"
        self.base_url = "https://webplay.weilekangnet.com:59666"

    def decrypt(self, inp: str):
        inp = inp.replace("\\", "")
        mode = AES.MODE_CBC
        cryptos = AES.new(self.key, mode, self.iv)
        res = cryptos.decrypt(base64.b64decode(inp))
        res = res.strip()
        res_str = res.decode("utf8").split(",")
        return f"{self.base_url}/" \
               f"{res_str[0]}/" \
               f"{res_str[1]}/" \
               f"{res_str[2]}/play.m3u8?_KS={res_str[3]}&_KE={res_str[4]}"
