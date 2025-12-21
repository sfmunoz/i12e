#!/usr/bin/env python3
from os.path import isfile, islink
from os import unlink, symlink, readlink, fchmod
from subprocess import call
import json
from logging import getLogger
log = getLogger(__name__)
from .butane import Butane

class GenesisInstall(object):
    def __init__(self):
        self.__base = "/genesis"
        self.__genesis_reboot = "{0}/genesis_reboot".format(self.__base)
        self.__genesis_restart_update_engine = "{0}/genesis_restart_update_engine".format(self.__base)
        self.__flatcar_first_boot = "{0}/boot/flatcar/first_boot".format(self.__base)

    def __trigger(self,fname,enable=None):
        if enable is None:
            return isfile(fname)
        if enable:
            with open(fname,"w") as fp:
                fchmod(fp.fileno(),0o600)
            log.info("'{0}' created".format(fname))
            return True
        unlink(fname)
        log.info("'{0}' deleted".format(fname))
        return False

    def __flatcar_extensions(self):
        for entry in ["containerd","docker"]:
            fname = "{0}/etc/extensions/{1}-flatcar.raw".format(self.__base,entry)
            try:
                if not islink(fname):
                    log.info("skipping '{0}': it's not a symbolic link".format(fname))
                    continue
                lname = readlink(fname)
                if lname == "/dev/null":
                    continue
                log.info("(bef) {0}: {1}".format(fname,lname))
                unlink(fname)
                symlink("/dev/null",fname)
                log.info("(aft) {0}: {1}".format(fname,readlink(fname)))
                self.__trigger(self.__genesis_reboot,True)
            except FileNotFoundError as e:
                log.warning("skipping '{0}': {1}".format(fname,str(e)))

    def __flatcar_update_conf(self):
        fname = "{0}/etc/flatcar/update.conf".format(self.__base)
        if not isfile(fname):
            log.info("skipping '{0}': it's not a regular file".format(fname))
            return
        with open("/app/conf/flatcar-update.conf","r") as fp:
            buf_new = fp.read()
        with open(fname,"r") as fp:
            buf_old = fp.read()
        if buf_old == buf_new:
            log.info("nothing to do: '{0}' is up-to-date".format(fname))
            return
        with open(fname,"w") as fp:
            fp.write(buf_new)
            fchmod(fp.fileno(),0o644)
        log.info("'{0}' updated".format(fname))
        self.__trigger(self.__genesis_restart_update_engine,True)

    def __restart_update_engine(self):
        if not self.__trigger(self.__genesis_restart_update_engine):
            return
        cmd = ["chroot",self.__base,"systemctl","restart","update-engine"]
        ret = call(cmd)
        if ret != 0:
            raise Exception("'{0}' command failed: ret={1}".format(" ".join(cmd),ret))
        log.info("update-engine restarted")
        self.__trigger(self.__genesis_restart_update_engine,False)

    # def __k3s_config_yaml(self):
    #     fname = "{0}/etc/rancher/k3s/config.yaml".format(self.__base)
    #     buf = """{{ include "genesis.k3s.config.yaml" . }}"""
    #     with open(fname,"w") as fp:
    #         fp.write(buf.strip() + "\n")
    #         fchmod(fp.fileno(),0o600)
    #     log.info("'{0}' created".format(fname))

    # def __k3s_override_conf(self):
    #     dname = "{0}/etc/systemd/system/k3s.service.d".format(self.__base)
    #     fname = "{0}/override.conf".format(dname)
    #     buf = """{{ include "genesis.k3s.override.conf" . }}"""
    #     if not isdir(dname):
    #         mkdir(dname)
    #     if not isdir(dname):
    #         raise Exception("error: couldn't create '{0}' folder".format(dname))
    #     chmod(dname,0o755)
    #     with open(fname,"w") as fp:
    #         fp.write(buf.strip() + "\n")
    #         fchmod(fp.fileno(),0o644)
    #     log.info("'{0}' created".format(fname))

    # def __manifest_skip(self):
    #     # https://docs.k3s.io/installation/packaged-components
    #     # don't let genesis.yaml run on k3s(etcd)
    #     fname = "{0}/var/lib/rancher/k3s/server/manifests/genesis.yaml.skip".format(self.__base)
    #     with open(fname,"w") as fp:
    #         fchmod(fp.fileno(),0o600)
    #     log.info("'{0}' created".format(fname))

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

    def __reboot_required(self):
        for fname in [self.__genesis_reboot,self.__flatcar_first_boot]:
            if not self.__trigger(fname):
                continue
            log.warning("triggering reboot: '{0}' file exists...".format(fname))
            if fname == self.__genesis_reboot:
                self.__trigger(fname,False)
            return True
        return False

    def __reboot(self):
        if not self.__reboot_required():
            return
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
        self.__flatcar_extensions()
        self.__flatcar_update_conf()
        self.__restart_update_engine()
        #self.__butane()
        self.__reboot()
        log.info("---- genesis install end ----")
