#!/usr/bin/env python3
from os.path import isfile
from subprocess import call
import json
from logging import getLogger
log = getLogger(__name__)
from .butane import Butane

class GenesisInstall(object):
    def __init__(self):
        self.__base = "/genesis"

    def __butane(self):
        fname = "/sec/ssh_authorized_key"
        with open(fname,"r") as fp:
            ssh_authorized_key = fp.read()
        config_ign_new = Butane(ssh_authorized_key).ignition()
        js_new = json.loads(config_ign_new)
        fname = "{0}/oem/config.ign".format(self.__base)
        with open(fname,"r") as fp:
            buf = fp.read()
        js_old = json.loads(buf)
        if js_new == js_old:
            log.info("nothing to do: '{0}' butane config did not change".format(fname))
            return
        log.info("upgrading flatcar...")
        fname = "{0}/root/config.ign".format(self.__base)
        with open(fname,"w") as fp:
            fp.write(config_ign_new)
        cmd = [
            "chroot",
            self.__base,
            "flatcar-reset",
            "--keep-machine-id",
            "--keep-paths",
            "/etc/ssh/ssh_host_.*",
            "/var/log",
            "/var/lib/rancher/k3s/agent/containerd",
            "-F",
            "/root/config.ign",
        ]
        ret = call(cmd)
        if ret != 0:
            raise Exception("'{0}' command failed: ret={1}".format(" ".join(cmd),ret))

    def __reboot(self):
        fname = "{0}/boot/flatcar/first_boot".format(self.__base)
        if not isfile(fname):
            log.info("reboot not required: '{0}' file doesn't exist".format(fname))
            return
        log.warning("triggering reboot: '{0}' file exists...".format(fname))
        cmd = [
            "chroot",
            self.__base,
            "systemd-run",
            "bash",
            "-c",
            "sleep 1 ; systemctl reboot",
        ]
        log.info("$ {0}".format(" ".join(cmd)))
        ret = call(cmd)
        if ret != 0:
            raise Exception("'{0}' command failed: ret={1}".format(" ".join(cmd),ret))

    def run(self):
        log.info("==== genesis install begin ====")
        self.__butane()
        self.__reboot()
        log.info("---- genesis install end ----")
