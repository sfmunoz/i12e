#!/usr/bin/env python3
from jinja2 import Template
import yaml
from subprocess import Popen,PIPE
from logging import getLogger
log = getLogger(__name__)

class Butane(object):
    def __init__(self,ssh_authorized_key):
        self.__tpl = """variant: flatcar
version: 1.0.0
passwd:
  users:
  - name: core
    ssh_authorized_keys:
    - "{{ ssh_authorized_key }}"
"""
        self.__ssh_authorized_key = ssh_authorized_key

    def ignition(self):
        buf = Template(self.__tpl) \
            .render(ssh_authorized_key=self.__ssh_authorized_key)
        _data = yaml.safe_load(buf)  # check it is valid
        cmd = ['butane']
        p = Popen(args=cmd,stdin=PIPE,stdout=PIPE,stderr=PIPE)
        (odata,edata) = p.communicate(buf.encode())
        if p.returncode != 0:
            raise Exception("'{0}' command failed: {1}".format(" ".join(cmd),edata.decode().strip()))
        return odata.decode().strip()
