#!/usr/bin/env python3
from jinja2 import Environment, PackageLoader, select_autoescape, StrictUndefined
from logging import getLogger
from io import BytesIO
from base64 import b64encode
from tarfile import TarInfo, SYMTYPE, DIRTYPE, open as tar_open
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
        self.__time = time()

    def __tarinfo(self,fname):
        tinfo = TarInfo(name=fname)
        tinfo.mode = 0o644
        tinfo.mtime = self.__time
        tinfo.uid = 0
        tinfo.gid = 0
        tinfo.uname = "root"
        tinfo.gname = "root"
        return tinfo

    def __flatcar_extensions(self,tar):
        for entry in ["containerd","docker"]:
            fname = "etc/extensions/{0}-flatcar.raw".format(entry)
            finfo = self.__tarinfo(fname)
            finfo.mode = 0o777
            finfo.type = SYMTYPE
            finfo.linkname = "/dev/null"
            tar.addfile(finfo)
            log.info("'{0}' added".format(fname))

    def __flatcar_update_conf(self,tar):
        data = (self.__tpl_flatcar_update_conf.render() + "\n").encode()
        fname = "etc/flatcar/update.conf"
        finfo = self.__tarinfo(fname)
        finfo.size = len(data)
        tar.addfile(finfo,BytesIO(data))
        log.info("'{0}' added".format(fname))

    def __k3s_config_yaml(self,tar):
        # https://docs.k3s.io/installation/configuration
        tls_san = "192.168.56.50"
        data = (self.__tpl_k3s_config_yaml.render(
            position = 1,
            k3s_cmd = "server",
            k3s_token = "main-token",
            k3s_agent_token = "agent-token",
            tls_san = tls_san,
            k3s_url = "https://{0}:6443".format(tls_san),
            flannel_iface = "enp0s8",
            node_ip = "192.168.56.51",
        ) + "\n").encode()
        fname = "etc/rancher/k3s/config.yaml"
        finfo = self.__tarinfo(fname)
        finfo.size = len(data)
        finfo.mode = 0o600
        tar.addfile(finfo,BytesIO(data))
        log.info("'{0}' added".format(fname))

    def __k3s_override_conf(self,tar):
        data = (self.__tpl_k3s_override_conf.render() + "\n").encode()
        dname = "etc/systemd/system/k3s.service.d"
        dinfo = self.__tarinfo(dname)
        dinfo.mode = 0o755
        dinfo.type = DIRTYPE
        tar.addfile(dinfo)
        fname = "{0}/override.conf".format(dname)
        finfo = self.__tarinfo(fname)
        finfo.size = len(data)
        tar.addfile(finfo,BytesIO(data))
        log.info("'{0}' added".format(fname))

    def __systemd_genesis_conf(self,tar):
        data = (self.__tpl_systemd_genesis_conf.render() + "\n").encode()
        dname = "etc/systemd/system.conf.d"
        dinfo = self.__tarinfo(dname)
        dinfo.mode = 0o755
        dinfo.type = DIRTYPE
        tar.addfile(dinfo)
        fname = "{0}/genesis.conf".format(dname)
        finfo = self.__tarinfo(fname)
        finfo.size = len(data)
        tar.addfile(finfo,BytesIO(data))
        log.info("'{0}' added".format(fname))

    def __etc_crictl_yaml(self,tar):
        data = (self.__tpl_critcl_yaml.render() + "\n").encode()
        fname = "etc/crictl.yaml"
        finfo = self.__tarinfo(fname)
        finfo.size = len(data)
        tar.addfile(finfo,BytesIO(data))
        log.info("'{0}' added".format(fname))

    def run(self):
        buf = BytesIO()
        with tar_open(fileobj=buf, mode="w:gz") as tar:
            self.__flatcar_extensions(tar)
            self.__flatcar_update_conf(tar)
            self.__k3s_config_yaml(tar)
            self.__k3s_override_conf(tar)
            self.__systemd_genesis_conf(tar)
            self.__etc_crictl_yaml(tar)
        print(b64encode(buf.getvalue()).decode())
        log.info("---- genesis artifact end ----")
