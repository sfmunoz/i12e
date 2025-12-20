{{- define "install.py" -}}
{{- $k3s_url := "https://192.168.56.50:6443" -}}
#!/usr/bin/env python3
from os import chmod,fchmod,readlink,symlink,unlink,mkdir,getenv
from os.path import islink,isfile,isdir
from subprocess import call
from logging import getLogger

log = getLogger(__name__)

class GenesisInstall(object):
    def __init__(self):
        self.__base = "/genesis"

    def __flatcar_extensions(self):
        for entry in ["containerd","docker"]:
            fname = "{0}/etc/extensions/{1}-flatcar.raw".format(self.__base,entry)
            try:
                if not islink(fname):
                    log.warning("skipping '{0}': it's not a symlink".format(fname))
                    continue
                log.info("(bef) {0}: {1}".format(fname,readlink(fname)))
                unlink(fname)
                symlink("/dev/null",fname)
                log.info("(aft) {0}: {1}".format(fname,readlink(fname)))
            except FileNotFoundError as e:
                log.warning("skipping '{0}': {1}".format(fname,str(e)))

    def __flatcar_update_conf(self):
        fname = "{0}/etc/flatcar/update.conf".format(self.__base)
        if not isfile(fname):
            log.warning("skipping '{0}': it's not a regular file".format(fname))
            return
        with open(fname,"r") as fp:
            buf_old = fp.read()
        buf_new = """{{ include "genesis.flatcar.update.conf" . }}"""
        if buf_old == buf_new:
            log.info("nothing to do: '{0}' is up-to-date".format(fname))
            return
        with open(fname,"w") as fp:
            fp.write(buf_new)
            fchmod(fp.fileno(),0o644)
        cmd = ["chroot",self.__base,"systemctl","restart","update-engine"]
        ret = call(cmd)
        if ret != 0:
            raise Exception("'{0}' command failed: ret={1}".format(" ".join(cmd),ret))
        log.info("'{0}' updated and update-engine restarted".format(fname))

    def __k3s_config_yaml(self):
        fname = "{0}/etc/rancher/k3s/config.yaml".format(self.__base)
        buf = """{{ include "genesis.k3s.config.yaml" . }}"""
        with open(fname,"w") as fp:
            fp.write(buf.strip() + "\n")
            fchmod(fp.fileno(),0o600)
        log.info("'{0}' created".format(fname))

    def __k3s_override_conf(self):
        dname = "{0}/etc/systemd/system/k3s.service.d".format(self.__base)
        fname = "{0}/override.conf".format(dname)
        buf = """{{ include "genesis.k3s.override.conf" . }}"""
        if not isdir(dname):
            mkdir(dname)
        if not isdir(dname):
            raise Exception("error: couldn't create '{0}' folder".format(dname))
        chmod(dname,0o755)
        with open(fname,"w") as fp:
            fp.write(buf.strip() + "\n")
            fchmod(fp.fileno(),0o644)
        log.info("'{0}' created".format(fname))

    def __manifest_skip(self):
        # https://docs.k3s.io/installation/packaged-components
        # don't let genesis.yaml run on k3s(etcd)
        fname = "{0}/var/lib/rancher/k3s/server/manifests/genesis.yaml.skip".format(self.__base)
        with open(fname,"w") as fp:
            fchmod(fp.fileno(),0o600)
        log.info("'{0}' created".format(fname))

    def __reboot(self):
        # #chroot /genesis systemd-run bash -c 'sleep 1 ; systemctl reboot'
        # # Failed to connect to system scope bus via local transport: No data available
        # chroot "${GENESIS_BASE}" systemctl reboot
        call(["chroot",self.__base,"systemctl","reboot"])

    def run(self):
        log.info("==== genesis install begin ====")
        #self.__flatcar_extensions()
        self.__flatcar_update_conf()
        #self.__k3s_config_yaml()
        #self.__k3s_override_conf()
        #self.__manifest_skip()
        #self.__reboot()
        log.info("---- genesis install end ----")  # never reached
{{ end }}
