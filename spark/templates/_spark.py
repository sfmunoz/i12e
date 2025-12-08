{{- define "spark.py" -}}
{{- $k3s_url := "https://192.168.56.50:6443" -}}
#!/usr/bin/env python3
from os import chmod,fchmod,readlink,symlink,unlink,mkdir
from os.path import islink,isfile,isdir
from sys import stderr
from subprocess import call
from logging import getLogger, basicConfig, INFO
basicConfig(format='%(asctime)s [%(relativeCreated)7.0f] [%(levelname).1s] %(message)s',level=INFO,stream=stderr)
log = getLogger(__name__)
class Spark(object):
    def __init__(self):
        self.__base = "/spark"
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
        buf = """{{ include "spark.flatcar.update.conf" . }}"""
        with open(fname,"w") as fp:
            fp.write(buf.strip() + "\n")
            fchmod(fp.fileno(),0o644)
        log.info("'{0}' updated".format(fname))
    def __k3s_config_yaml(self):
        fname = "{0}/etc/rancher/k3s/config.yaml".format(self.__base)
        buf = """{{ include "spark.k3s.config.yaml" . }}"""
        with open(fname,"w") as fp:
            fp.write(buf.strip() + "\n")
            fchmod(fp.fileno(),0o600)
        log.info("'{0}' created".format(fname))
    def __k3s_override_conf(self):
        dname = "{0}/etc/systemd/system/k3s.service.d".format(self.__base)
        fname = "{0}/override.conf".format(dname)
        buf = """{{ include "spark.k3s.override.conf" . }}"""
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
        # don't let spark.yaml run on k3s(etcd)
        fname = "{0}/var/lib/rancher/k3s/server/manifests/spark.yaml.skip".format(self.__base)
        with open(fname,"w") as fp:
            fchmod(fp.fileno(),0o600)
        log.info("'{0}' created".format(fname))
    def __reboot(self):
        # #chroot /spark systemd-run bash -c 'sleep 1 ; systemctl reboot'
        # # Failed to connect to system scope bus via local transport: No data available
        # chroot "${SPARK_BASE}" systemctl reboot
        call(["chroot",self.__base,"systemctl","reboot"])
    def run(self):
        log.info("==== spark begin ====")
        self.__flatcar_extensions()
        self.__flatcar_update_conf()
        self.__k3s_config_yaml()
        self.__k3s_override_conf()
        self.__manifest_skip()
        #self.__reboot()
        log.info("---- spark end ----")  # never reached
if __name__ == "__main__":
    Spark().run()
{{- end }}
