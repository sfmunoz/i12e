#!/usr/bin/env python3
from os import environ
from jinja2 import Environment, PackageLoader, select_autoescape, StrictUndefined
from logging import getLogger
from io import BytesIO
from tarfile import TarInfo, DIRTYPE, open as tar_open
from time import time
from subprocess import Popen, PIPE
from hashlib import sha256
from .config import Config
log = getLogger(__name__)

class Artifact(object):
    def __init__(self,modes):
        self.__modes = modes
        self.__cfg = Config().main_config()
        self.__env = Environment(
            loader = PackageLoader("genesis"),
            undefined = StrictUndefined,
            autoescape = select_autoescape(),
        )
        self.__tpl_critcl_yaml = self.__env.get_template("crictl.yaml")
        self.__tpl_opt_bin_e = self.__env.get_template("opt-bin-e.yaml")
        self.__tpl_flatcar_update_conf = self.__env.get_template("flatcar-update.conf")
        self.__tpl_k3s_config_yaml = self.__env.get_template("k3s-config.yaml")
        self.__tpl_k3s_override_conf = self.__env.get_template("k3s-override.conf")
        self.__tpl_systemd_genesis_conf = self.__env.get_template("systemd-genesis.conf")
        self.__tpl_i12e_k3s_install = self.__env.get_template("i12e-k3s-install.sh")
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

    def __flatcar_update_conf(self,tar):
        data = (self.__tpl_flatcar_update_conf.render() + "\n").encode()
        fname = "etc/flatcar/update.conf"
        finfo = self.__tarinfo(fname)
        finfo.size = len(data)
        tar.addfile(finfo,BytesIO(data))
        log.info("'{0}' added".format(fname))

    def __etc_i12e_k3s(self,tar):
        dname = "etc/i12e/k3s"
        dinfo = self.__tarinfo(dname)
        dinfo.mode = 0o700
        dinfo.type = DIRTYPE
        tar.addfile(dinfo)
        log.info("'{0}' added".format(dname))

    def __k3s_config_yaml(self,tar,mode):
        # https://docs.k3s.io/installation/configuration
        tls_san = self.__cfg["kube_vip"]["vip"]
        data = (self.__tpl_k3s_config_yaml.render(
            i12e_mode = mode,
            k3s_token = self.__cfg["k3s_token"],
            k3s_agent_token = self.__cfg["k3s_agent_token"],
            tls_san = tls_san,
            k3s_url = "https://{0}:6443".format(tls_san),
        ) + "\n").encode()
        fname = "etc/i12e/k3s/config-{0}.yaml".format(mode)
        finfo = self.__tarinfo(fname)
        finfo.size = len(data)
        finfo.mode = 0o600
        tar.addfile(finfo,BytesIO(data))
        log.info("'{0}' added".format(fname))

    def __k3s_service_d(self,tar):
        dname = "etc/systemd/system/k3s.service.d"
        dinfo = self.__tarinfo(dname)
        dinfo.mode = 0o755
        dinfo.type = DIRTYPE
        tar.addfile(dinfo)
        log.info("'{0}' added".format(dname))

    def __k3s_override_conf(self,tar):
        data = (self.__tpl_k3s_override_conf.render() + "\n").encode()
        fname = "etc/systemd/system/k3s.service.d/override.conf"
        finfo = self.__tarinfo(fname)
        finfo.size = len(data)
        tar.addfile(finfo,BytesIO(data))
        log.info("'{0}' added".format(fname))

    def __system_conf_d(self,tar):
        dname = "etc/systemd/system.conf.d"
        dinfo = self.__tarinfo(dname)
        dinfo.mode = 0o755
        dinfo.type = DIRTYPE
        tar.addfile(dinfo)
        log.info("'{0}' added".format(dname))

    def __systemd_genesis_conf(self,tar):
        data = (self.__tpl_systemd_genesis_conf.render() + "\n").encode()
        fname = "etc/systemd/system.conf.d/genesis.conf"
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

    def __opt_bin_e(self,tar):
        data = (self.__tpl_opt_bin_e.render() + "\n").encode()
        fname = "opt/bin/e"
        finfo = self.__tarinfo(fname)
        finfo.mode = 0o755
        finfo.size = len(data)
        tar.addfile(finfo,BytesIO(data))
        log.info("'{0}' added".format(fname))

    def __etc_i12e(self,tar):
        dname = "etc/i12e"
        dinfo = self.__tarinfo(dname)
        dinfo.mode = 0o700
        dinfo.type = DIRTYPE
        tar.addfile(dinfo)
        log.info("'{0}' added".format(dname))

    def __etc_i12e_iface_txt(self,tar):
        fname = "etc/i12e/iface.txt"
        iface = self.__cfg.get("flannel",{}).get("interface")
        if iface is None:
            log.info("'{0}' not added: flannel.interface not defined".format(fname))
            return
        iface = iface.strip()
        if len(iface) < 1:
            log.info("'{0}' not added: flannel.interface is empty".format(fname))
            return
        data = (iface + "\n").encode()
        finfo = self.__tarinfo(fname)
        finfo.mode = 0o600
        finfo.size = len(data)
        tar.addfile(finfo,BytesIO(data))
        log.info("'{0}' added".format(fname))

    def __opt_libexec(self,tar):
        dname = "opt/libexec"
        dinfo = self.__tarinfo(dname)
        dinfo.mode = 0o755
        dinfo.type = DIRTYPE
        tar.addfile(dinfo)
        log.info("'{0}' added".format(dname))

    def __i12e_k3s_install_sh(self,tar):
        data = (self.__tpl_i12e_k3s_install.render() + "\n").encode()
        fname = "opt/libexec/i12e-k3s-install.sh"
        finfo = self.__tarinfo(fname)
        finfo.mode = 0o700
        finfo.size = len(data)
        tar.addfile(finfo,BytesIO(data))
        log.info("'{0}' added".format(fname))

    def __etc_i12e_pull_done(self,tar):
        fname = "etc/i12e/pull-done"
        finfo = self.__tarinfo(fname)
        finfo.mode = 0o600
        tar.addfile(finfo)
        log.info("'{0}' added".format(fname))

    def __rclone_push(self,buf):
        rclone_config = "/root/rclone.conf"
        with open(rclone_config,"w") as fp:
            fp.write(Config().rclone_config())
        environ["RCLONE_CONFIG"] = rclone_config   # after encrypted one has been read
        # --------
        cmd = ['rclone','rcat','rem:artifact.tar.gz']
        log.info("$ {0}".format(" ".join(cmd)))
        p = Popen(args=cmd,stdin=PIPE)
        p.communicate(buf)
        if p.returncode != 0:
            raise Exception("'{0}' command failed".format(" ".join(cmd)))
        # --------
        cmd = ['rclone','cat','rem:artifact.tar.gz']
        log.info("$ {0}".format(" ".join(cmd)))
        p = Popen(args=cmd,stdout=PIPE)
        (odata,_) = p.communicate()
        if p.returncode != 0:
            raise Exception("'{0}' command failed".format(" ".join(cmd)))
        # --------
        sha256_1 = sha256(buf).hexdigest()
        sha256_2 = sha256(odata).hexdigest()
        log.info("sha256(bef): {0}".format(sha256_1))
        log.info("sha256(aft): {0}".format(sha256_2))
        if sha256_1 != sha256_2:
            raise Exception("sha256 checksum mismatch")
        # --------
        cmd = ['tar','tvz']
        log.info("$ {0}".format(" ".join(cmd)))
        p = Popen(args=cmd,stdin=PIPE)
        p.communicate(odata)
        if p.returncode != 0:
            raise Exception("'{0}' command failed".format(" ".join(cmd)))

    def run(self):
        log.info("==== genesis artifact begin ====")
        buf = BytesIO()
        with tar_open(fileobj=buf, mode="w:gz") as tar:
            self.__flatcar_update_conf(tar)
            self.__etc_i12e_k3s(tar)
            for mode in self.__modes:
                self.__k3s_config_yaml(tar,mode)
            self.__k3s_service_d(tar)
            self.__k3s_override_conf(tar)
            self.__system_conf_d(tar)
            self.__systemd_genesis_conf(tar)
            self.__etc_crictl_yaml(tar)
            self.__opt_bin_e(tar)
            self.__etc_i12e(tar)
            self.__etc_i12e_iface_txt(tar)
            self.__opt_libexec(tar)
            self.__i12e_k3s_install_sh(tar)
            self.__etc_i12e_pull_done(tar)
        self.__rclone_push(buf.getvalue())
        log.info("---- genesis artifact end ----")
