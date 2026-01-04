#!/usr/bin/env python3
from os import getenv
from logging import getLogger
from subprocess import Popen, PIPE
from json import loads as json_loads
log = getLogger(__name__)

class Rclone(object):
    def __init__(self):
        self.__remote = getenv("I12E_RCLONE_REMOTE")
        if self.__remote is None or len(self.__remote) < 1:
            raise Exception("cannot find 'I12E_RCLONE_REMOTE' value")

    def __get_config(self):
        cmd = ['rclone','config','dump']
        p = Popen(args=cmd,stdin=PIPE,stdout=PIPE,stderr=PIPE)
        (odata,edata) = p.communicate()
        if p.returncode != 0:
            raise Exception("'{0}' command failed: {1}".format(" ".join(cmd),edata.decode().strip()))
        return json_loads(odata.decode())

    def __remote_dump(self,cfg,name):
        if name not in cfg:
            raise Exception("cannot find '{0}' remote".format(name))
        c = cfg[name]
        lines = ["[{0}]".format(name)]
        for k,v in c.items():
            lines.append("{0} = {1}".format(k,v))
        if "remote" not in c.keys():
            return lines
        lines.append("")
        rem = c["remote"].split(":")[0]  # remote = name:path/to/subfolder
        lines.extend(self.__remote_dump(cfg,rem))
        return lines

    def run(self):
        cfg = self.__get_config()
        lines = self.__remote_dump(cfg,self.__remote)
        return "\n".join(lines)

