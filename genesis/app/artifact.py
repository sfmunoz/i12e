#!/usr/bin/env python3
from os.path import isfile, islink, isdir
from os import unlink, symlink, readlink, fchmod, mkdir, chmod, getenv
from jinja2 import Environment, PackageLoader, select_autoescape, StrictUndefined
from logging import getLogger
from io import BytesIO
from base64 import b64encode
from tarfile import TarInfo, open as tar_open
from time import time
log = getLogger(__name__)

class Artifact(object):
    def __init__(self):
        self.__env = Environment(
            loader = PackageLoader("genesis"),
            undefined = StrictUndefined,
            autoescape = select_autoescape(),
        )
        self.__tpl_critcl_yaml = self.__env.get_template("crictl.yaml")
        self.__tpl_flatcar_update_conf = self.__env.get_template("flatcar-update.conf")
        self.__tpl_k3s_config_yaml = self.__env.get_template("k3s-config.yaml")
        self.__tpl_k3s_override_conf = self.__env.get_template("k3s-override.conf")
        self.__tpl_systemd_genesis_conf = self.__env.get_template("systemd-genesis.conf")

    def __flatcar_extensions(self):
        for entry in ["containerd","docker"]:
            fname = "/etc/extensions/{0}-flatcar.raw".format(entry)
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
            except FileNotFoundError as e:
                log.warning("skipping '{0}': {1}".format(fname,str(e)))

    def __flatcar_update_conf(self):
        fname = "/etc/flatcar/update.conf"
        if not isfile(fname):
            log.info("skipping '{0}': it's not a regular file".format(fname))
            return
        buf_new = self.__tpl_flatcar_update_conf.render() + "\n"
        with open(fname,"r") as fp:
            buf_old = fp.read()
        if buf_old == buf_new:
            log.info("nothing to do: '{0}' is up-to-date".format(fname))
            return
        with open(fname,"w") as fp:
            fp.write(buf_new)
            fchmod(fp.fileno(),0o644)
        log.info("'{0}' updated".format(fname))

    def __k3s_config_yaml(self):
        # https://docs.k3s.io/installation/configuration
        tls_san = "192.168.56.50"
        buf_new = self.__tpl_k3s_config_yaml.render(
            position = 1,
            k3s_cmd = "server",
            k3s_token = "main-token",
            k3s_agent_token = "agent-token",
            tls_san = tls_san,
            k3s_url = "https://{0}:6443".format(tls_san),
            flannel_iface = "enp0s8",
            node_ip = "192.168.56.51",
        ) + "\n"
        fname = "/etc/rancher/k3s/config.yaml"
        buf_old = ""
        if isfile(fname):
            with open(fname,"r") as fp:
                buf_old = fp.read()
        if buf_old == buf_new:
            log.info("nothing to do: '{0}' is up to date".format(fname))
            return
        with open(fname,"w") as fp:
            fp.write(buf_new)
            fchmod(fp.fileno(),0o600)
        log.info("'{0}' created/updated".format(fname))

    def __k3s_override_conf(self):
        buf_new = self.__tpl_k3s_override_conf.render() + "\n"
        dname = "/etc/systemd/system/k3s.service.d"
        fname = "{0}/override.conf".format(dname)
        if not isdir(dname):
            mkdir(dname)
        if not isdir(dname):
            raise Exception("error: couldn't create '{0}' folder".format(dname))
        chmod(dname,0o755)  # TODO: avoid doing this on every iteration
        buf_old = ""
        if isfile(fname):
            with open(fname,"r") as fp:
                buf_old = fp.read()
        if buf_old == buf_new:
            log.info("nothing to do: '{0}' is up to date".format(fname))
            return
        with open(fname,"w") as fp:
            fp.write(buf_new)
            fchmod(fp.fileno(),0o644)
        log.info("'{0}' created/updated".format(fname))

    def __systemd_genesis_conf(self):
        buf_new = self.__tpl_systemd_genesis_conf.render() + "\n"
        dname = "/etc/systemd/system.conf.d"
        fname = "{0}/genesis.conf".format(dname)
        if not isdir(dname):
            mkdir(dname)
        if not isdir(dname):
            raise Exception("error: couldn't create '{0}' folder".format(dname))
        chmod(dname,0o755)  # TODO: avoid doing this on every iteration
        buf_old = ""
        if isfile(fname):
            with open(fname,"r") as fp:
                buf_old = fp.read()
        if buf_old == buf_new:
            log.info("nothing to do: '{0}' is up to date".format(fname))
            return
        with open(fname,"w") as fp:
            fp.write(buf_new)
            fchmod(fp.fileno(),0o644)
        log.info("'{0}' created/updated".format(fname))

    def __etc_crictl_yaml(self,tar):
        fname = "etc/crictl.yaml"
        data = (self.__tpl_critcl_yaml.render() + "\n").encode()
        tinfo = TarInfo(name=fname)
        tinfo.mode = 0o644
        tinfo.mtime = time()
        tinfo.uid = 0
        tinfo.gid = 0
        tinfo.uname = "root"
        tinfo.gname = "root"
        tinfo.size = len(data)
        tar.addfile(tinfo,BytesIO(data))
        log.info("'{0}' added".format(fname))

    def run(self):
        if getenv("I12E_LEGACY") == "1":
            self.__flatcar_extensions()
            self.__flatcar_update_conf()
            self.__k3s_config_yaml()
            self.__k3s_override_conf()
            self.__systemd_genesis_conf()
        buf = BytesIO()
        with tar_open(fileobj=buf, mode="w:gz") as tar:
            self.__etc_crictl_yaml(tar)
        print(b64encode(buf.getvalue()).decode())
        log.info("---- genesis artifact end ----")
