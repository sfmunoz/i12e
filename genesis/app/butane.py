#!/usr/bin/env python3
from os import getenv
from sys import stdout
from jinja2 import Environment, PackageLoader, select_autoescape
import yaml, json
from subprocess import Popen, PIPE
from logging import getLogger
from datetime import datetime, timedelta, UTC
from base64 import b64encode
from gzip import compress
log = getLogger(__name__)

O_BASH_B64 = 0
O_BASH_RAW = 1
O_IGNITION = 2
O_DEBUG = 3

class Butane(object):
    def __init__(self,target):
        self.__target = target
        e = getenv("GENESIS_OUTPUT")
        self.__output = O_DEBUG if e == "debug" else O_IGNITION if e == "ignition" else O_BASH_RAW if e == "bash_raw" else O_BASH_B64
        self.__env = Environment(
            loader = PackageLoader("genesis"),
            autoescape = select_autoescape(),
        )
        self.__tpl = self.__env.get_template("flatcar.yaml")
        with open("/ssh_authorized_key","r") as fp:
            self.__ssh_authorized_key = fp.read().strip()
        self.__fp = stdout

    def __buf_print(self,buf,prefix=""):
        if self.__output != O_DEBUG:
            return
        for line in buf.strip().split("\n"):
            self.__fp.write("{0}{1}\n".format(prefix,line))

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

    def __bash(self,buf):
        config_ign = "/oem/config.ign"
        cmd = " ; ".join([
            "set -x -e -o pipefail",
            "rm -fv {0}".format(config_ign),
            " | ".join([
                "base64 -d <<< \"" + b64encode(compress(buf.encode())).decode() + "\"",
                "gunzip",
                " ".join([
                    "flatcar-reset",
                    "--keep-machine-id",
                    "--keep-paths",
                    "'/etc/ssh/ssh_host_.*'",
                    "/var/log",
                    "/var/lib/rancher/k3s/agent/containerd",
                    "-F",
                    "/dev/stdin",
                ]),
            ]),
            "test -s {0}".format(config_ign),
            "jq < {0}".format(config_ign),
            "systemd-run bash -c 'sleep 1 ; systemctl reboot'",
        ])
        if self.__output == O_BASH_RAW:
            self.__fp.write(cmd)
        else:
            cmd2 = " | ".join([
                "base64 -d <<< \"" + b64encode(compress(cmd.encode())).decode() + "\"",
                "gunzip",
                "sudo bash",
            ])
            self.__fp.write(cmd2)
        self.__fp.flush()

    def __inject(self):
        buf1 = self.__ignition()
        if self.__output == O_IGNITION:
            self.__fp.write(buf1)
            self.__fp.flush()
            return
        if self.__output in [O_BASH_B64,O_BASH_RAW]:
            self.__bash(buf1)
            return
        js = json.loads(buf1)
        buf2 = json.dumps(js,indent=2,sort_keys=True)
        self.__buf_print(buf2,"<ign> ")
        tgt_user ="core@{0}".format(self.__target)
        tgt_file = datetime.now(UTC).strftime("os-%Y%m%d-%H%M%S.json")
        cmd = [
            "ssh",
            tgt_user,
            "cat > {0}".format(tgt_file),
        ]
        log.info("$ {0}".format(" ".join(cmd)))
        #p = Popen(args=cmd,stdin=PIPE,stdout=PIPE,stderr=PIPE)
        #(_,edata) = p.communicate(buf2.encode())
        #if p.returncode != 0:
        #    raise Exception("'{0}' command failed: {1}".format(" ".join(cmd),edata.decode().strip()))
        cmd = [
            "ssh",
            tgt_user,
            "sudo flatcar-reset --keep-machine-id --keep-paths '/etc/ssh/ssh_host_.*' /var/log /var/lib/rancher/k3s/agent/containerd -F {0} && sudo systemd-run bash -c 'sleep 1 ; systemctl reboot'".format(tgt_file),
        ]
        log.info("$ {0}".format(" ".join(cmd)))
        #p = Popen(args=cmd,stdin=PIPE,stdout=PIPE,stderr=PIPE)
        #(_,edata) = p.communicate(buf2.encode())
        #if p.returncode != 0:
        #    raise Exception("'{0}' command failed: {1}".format(" ".join(cmd),edata.decode().strip()))

    def run(self):
        self.__inject()
