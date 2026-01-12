#!/usr/bin/env python3
from os import getenv
import yaml
from logging import getLogger
from base64 import b64decode
from gzip import decompress
from .rclone import Rclone
log = getLogger(__name__)

class Config(object):
    def __init__(self):
        buf = getenv("I12E_CONFIG")
        if buf is None or len(buf) < 1:
            raise Exception("undefined 'I12E_CONFIG' env-var")
        self.__cfg = yaml.safe_load(decompress(b64decode(buf)))

    def main_config(self):
        return self.__cfg

    def butane_config(self):
        buf = getenv("I12E_BUTANE")
        if buf is None or len(buf) < 1:
            return {}
        return yaml.safe_load(decompress(b64decode(buf)))

    def rclone_config(self):
        return Rclone(self.__cfg["rclone_remote"],self.__cfg["rclone_config_pass"]).config()
