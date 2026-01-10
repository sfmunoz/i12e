#!/usr/bin/env python3
from sys import stdout
from jinja2 import Environment, PackageLoader, select_autoescape, StrictUndefined
import yaml, json
from subprocess import Popen, PIPE
from logging import getLogger
from base64 import b64encode
from gzip import compress
from .config import Config
log = getLogger(__name__)

O_BASH_B64 = "bash_b64"
O_BASH_RAW = "bash_raw"
O_IGNITION = "ignition"
O_DEBUG = "debug"

I12E_VERSION = "v0.0.14"
I12E_SHA256SUM = "bbeeb3aa306641f8b678dda388f3e86de4ab4411d9a9f22f79a94a22c3b9b10c"

class Butane(object):
    def __init__(self,args):
        self.__cfg = Config().main_config()
        self.__k3s_version = self.__cfg["k3s_version"]
        self.__mode = args.mode
        self.__output = args.output
        self.__env = Environment(
            loader = PackageLoader("genesis"),
            undefined = StrictUndefined,
            autoescape = select_autoescape(),
        )
        self.__tpl = self.__env.get_template("flatcar.yaml")
        self.__ssh_authorized_keys = self.__cfg["ssh_authorized_keys"]
        if len(self.__ssh_authorized_keys) < 1:
            raise Exception("'ssh_authorized_keys' list is empty")
        self.__fp = stdout

    def __buf_print(self,buf,prefix=""):
        if self.__output != O_DEBUG:
            return
        for line in buf.strip().split("\n"):
            self.__fp.write("{0}{1}\n".format(prefix,line))

    def __ignition(self):
        buf = self.__tpl.render(
            i12e_version = I12E_VERSION,
            i12e_sha256sum = I12E_SHA256SUM,
            mode = self.__mode,
            k3s_version = self.__k3s_version,
            ssh_authorized_keys = self.__ssh_authorized_keys,
            rclone_conf = Config().rclone_config(),
        )
        self.__buf_print(buf,"<but> ")
        yaml.safe_load(buf)  # return value ignored: check it is valid
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
            "sudo rm -fv {0}".format(config_ign),
            " | ".join([
                "base64 -d <<< \"" + b64encode(compress(buf.encode())).decode() + "\"",
                "gunzip",
                " ".join([
                    "sudo flatcar-reset",
                    "--keep-machine-id",
                    "--keep-paths",
                    "'/etc/ssh/ssh_host_.*'",
                    "/var/log",
                    "/var/lib/rancher/k3s/agent/containerd",
                    "-F",
                    "/dev/stdin",
                ]),
            ]),
            "sudo test -s {0}".format(config_ign),
            "sudo jq . {0}".format(config_ign),
            "sudo systemd-run bash -c 'sleep 1 ; systemctl reboot'",
        ])
        if self.__output == O_BASH_RAW:
            self.__fp.write(cmd)
        else:
            cmd2 = " | ".join([
                "base64 -d <<< \"" + b64encode(compress(cmd.encode())).decode() + "\"",
                "gunzip",
                "bash",
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

    def run(self):
        self.__inject()
