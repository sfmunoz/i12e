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
    def __init__(self,args,modes):
        self.__devel = args.devel
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
        self.__tpl_systemd_i12e_conf = self.__env.get_template("systemd-i12e.conf")
        self.__tpl_artifact_tune = self.__env.get_template("artifact-tune.sh")
        self.__tpl_nftables_conf = self.__env.get_template("nftables.conf")
        self.__tpl_nftables_service = self.__env.get_template("nftables.service")
        self.__time = time()

    def __tarinfo(self,fname):
        tinfo = TarInfo(name=fname)
        tinfo.mode = 0o644
        tinfo.mtime = self.__time
        tinfo.uname = "root"
        tinfo.gname = "root"
        return tinfo

    def __folders(self,tar):
        folders = [
            ("etc/i12e",0o700,"root","root"),
            ("etc/i12e/flags",0o700,"root","root"),
            ("etc/i12e/k3s",0o700,"root","root"),
            ("etc/systemd/system/k3s.service.d",0o755,"root","root"),
            ("etc/systemd/system.conf.d",0o755,"root","root"),
            ("opt/libexec",0o755,"root","root"),
            ("opt/libexec/i12e",0o755,"root","root"),
        ]
        for f in folders:
            ti = self.__tarinfo(f[0])
            ti.mode = f[1]
            ti.type = DIRTYPE
            ti.uname = f[2]
            ti.gname = f[3]
            tar.addfile(ti)
            log.info("'{0}' added".format(f[0]))

    def __etc_crictl_yaml(self,tar):
        data = (self.__tpl_critcl_yaml.render() + "\n").encode()
        fname = "etc/crictl.yaml"
        finfo = self.__tarinfo(fname)
        finfo.size = len(data)
        tar.addfile(finfo,BytesIO(data))
        log.info("'{0}' added".format(fname))

    def __flatcar_update_conf(self,tar):
        data = (self.__tpl_flatcar_update_conf.render() + "\n").encode()
        fname = "etc/flatcar/update.conf"
        finfo = self.__tarinfo(fname)
        finfo.size = len(data)
        tar.addfile(finfo,BytesIO(data))
        log.info("'{0}' added".format(fname))

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

    def __k3s_config_yaml(self,tar,mode=None):
        if mode is None:
            fname = "etc/rancher/k3s/config.yaml"
            finfo = self.__tarinfo(fname)
            finfo.mode = 0o600
            tar.addfile(finfo)
            log.info("'{0}' added".format(fname))
            return
        # https://docs.k3s.io/installation/configuration
        data = (self.__tpl_k3s_config_yaml.render(
            i12e_mode = mode,
            k3s_token = self.__cfg["k3s_token"],
            k3s_agent_token = self.__cfg["k3s_agent_token"],
            k3s_url = self.__cfg["k3s_url"],
            tls_san = self.__cfg.get("tls_san"),
        ) + "\n").encode()
        fname = "etc/i12e/k3s/config-{0}.yaml".format(mode)
        finfo = self.__tarinfo(fname)
        finfo.size = len(data)
        finfo.mode = 0o600
        tar.addfile(finfo,BytesIO(data))
        log.info("'{0}' added".format(fname))

    def __k3s_override_conf(self,tar,mode=None):
        if mode is None:
            fname = "etc/systemd/system/k3s.service.d/override.conf"
            finfo = self.__tarinfo(fname)
            finfo.mode = 0o644
            tar.addfile(finfo)
            log.info("'{0}' added".format(fname))
            return
        data = (self.__tpl_k3s_override_conf.render(
            i12e_mode = mode,
        ) + "\n").encode()
        fname = "etc/i12e/k3s/override-{0}.conf".format(mode)
        finfo = self.__tarinfo(fname)
        finfo.size = len(data)
        finfo.mode = 0o644
        tar.addfile(finfo,BytesIO(data))
        log.info("'{0}' added".format(fname))

    def __systemd_i12e_conf(self,tar):
        data = (self.__tpl_systemd_i12e_conf.render() + "\n").encode()
        fname = "etc/systemd/system.conf.d/i12e.conf"
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

    def __artifact_tune_sh(self,tar):
        data = (self.__tpl_artifact_tune.render() + "\n").encode()
        fname = "opt/libexec/i12e/artifact-tune.sh"
        finfo = self.__tarinfo(fname)
        finfo.mode = 0o755
        finfo.size = len(data)
        tar.addfile(finfo,BytesIO(data))
        log.info("'{0}' added".format(fname))

    def __etc_i12e_flags_artifact_pulled(self,tar):
        fname = "etc/i12e/flags/artifact-pulled"
        finfo = self.__tarinfo(fname)
        finfo.mode = 0o600
        tar.addfile(finfo)
        log.info("'{0}' added".format(fname))

    def __nftables_conf(self,tar):
        data = (self.__tpl_nftables_conf.render(
            port_knocking = self.__cfg["port_knocking"],
        ) + "\n").encode()
        fname = "etc/nftables.conf"
        finfo = self.__tarinfo(fname)
        finfo.mode = 0o600
        finfo.size = len(data)
        tar.addfile(finfo,BytesIO(data))
        log.info("'{0}' added".format(fname))

    def __nftables_service(self,tar):
        data = (self.__tpl_nftables_service.render() + "\n").encode()
        fname = "etc/systemd/system/nftables.service"
        finfo = self.__tarinfo(fname)
        finfo.size = len(data)
        tar.addfile(finfo,BytesIO(data))
        log.info("'{0}' added".format(fname))

    def __devel_output(self,buf):
        log.info("sha256: {0}".format(sha256(buf).hexdigest()))
        cmd = ['tar','tvz']
        log.info("$ {0}".format(" ".join(cmd)))
        p = Popen(args=cmd,stdin=PIPE)
        p.communicate(buf)
        if p.returncode != 0:
            raise Exception("'{0}' command failed".format(" ".join(cmd)))

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
            self.__folders(tar)
            self.__etc_crictl_yaml(tar)
            self.__flatcar_update_conf(tar)
            self.__etc_i12e_iface_txt(tar)
            for mode in self.__modes:
                self.__k3s_config_yaml(tar,mode)
            for mode in self.__modes:
                self.__k3s_override_conf(tar,mode)
            self.__nftables_conf(tar)
            self.__systemd_i12e_conf(tar)
            self.__k3s_override_conf(tar)
            self.__nftables_service(tar)
            self.__k3s_config_yaml(tar)
            self.__opt_bin_e(tar)
            self.__artifact_tune_sh(tar)
            self.__etc_i12e_flags_artifact_pulled(tar)
        if self.__devel:
            self.__devel_output(buf.getvalue())
            return
        self.__rclone_push(buf.getvalue())
        log.info("---- genesis artifact end ----")
