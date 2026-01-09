#!/usr/bin/env python3
from os import environ
from logging import getLogger
from subprocess import Popen, PIPE
from json import loads as json_loads
log = getLogger(__name__)

# XXX: don't meant to be used directly: use 'Config.rclone_config()' instead

class Rclone(object):
    def __init__(self,cfg):
        self.__cfg = cfg
        self.__remote = self.__cfg["rclone_remote"]

    def __get_config(self):
        cmd = ['rclone','config','dump']
        env = environ.copy()
        env["RCLONE_CONFIG_PASS"] = self.__cfg["rclone_config_pass"]
        p = Popen(args=cmd,stdin=PIPE,stdout=PIPE,stderr=PIPE,env=env)
        (odata,edata) = p.communicate()
        if p.returncode != 0:
            raise Exception("'{0}' command failed: {1}".format(" ".join(cmd),edata.decode().strip()))
        return json_loads(odata.decode())

    def __remote_dump(self,cfg,name):
        if name not in cfg:
            raise Exception("cannot find '{0}' remote".format(name))
        c = cfg[name]
        lines = ["[{0}]".format("rem" if name == self.__remote else name)]
        for k,v in c.items():
            lines.append("{0} = {1}".format(k,v))
        if "remote" not in c.keys():
            return lines
        rem = c["remote"].split(":")[0]  # remote = name:path/to/subfolder
        return self.__remote_dump(cfg,rem) + [""] + lines

    def config(self):
        cfg = self.__get_config()
        lines = self.__remote_dump(cfg,self.__remote)
        return "\n".join(lines)
