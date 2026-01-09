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
        self.__env = getenv("I12E_ENV")
        if self.__env is None:
            raise Exception("undefined 'I12E_ENV' env-var")
        if self.__env not in ["dev","prod"]:
            raise Exception("unknown I12E_ENV='{0}' value (valid: 'dev' or 'prod')".format(self.__env))

    def secrets_yaml(self):
        buf = getenv("I12E_SECRETS_YAML")
        if buf is None or len(buf) < 1:
            raise Exception("undefined 'I12E_SECRETS_YAML' env-var")
        return yaml.safe_load(decompress(b64decode(buf)))["env"][self.__env]

    def rclone_config(self):
        c = self.secrets_yaml()
        return Rclone(c["rclone_remote"],c["rclone_config_pass"]).config()
