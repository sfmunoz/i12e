#!/usr/bin/env python3
from jinja2 import Environment, PackageLoader, select_autoescape
import yaml, json
from subprocess import Popen, PIPE
from logging import getLogger
log = getLogger(__name__)

class Butane(object):
    def __init__(self):
        self.__env = Environment(
            loader = PackageLoader("genesis"),
            autoescape = select_autoescape(),
        )
        self.__tpl = self.__env.get_template("flatcar.yaml")
        self.__ssh_authorized_key = "<ssh-authorized-key>"  # TODO: set real file here

    def __buf_print(self,buf,prefix=""):
        for line in buf.strip().split("\n"):
            print("{0}{1}".format(prefix,line))

    def __ignition(self):
        buf = self.__tpl.render(
            ssh_authorized_key = self.__ssh_authorized_key,
        )
        self.__buf_print(buf,"<but> ")
        _data = yaml.safe_load(buf)  # check it is valid
        cmd = ['butane']
        p = Popen(args=cmd,stdin=PIPE,stdout=PIPE,stderr=PIPE)
        (odata,edata) = p.communicate(buf.encode())
        if p.returncode != 0:
            raise Exception("'{0}' command failed: {1}".format(" ".join(cmd),edata.decode().strip()))
        return odata.decode().strip()

    def __butane(self):
        buf1 = self.__ignition()
        js = json.loads(buf1)
        buf = json.dumps(js,indent=2,sort_keys=True)
        self.__buf_print(buf,"<ign> ")

    def run(self):
        self.__butane()
