{{- define "spark.py" -}}
{{- $k3s_url := "https://192.168.56.50:6443" -}}
#!/usr/bin/env python3
from os import fchmod,readlink,symlink,unlink
from os.path import islink,isfile
from sys import stderr
from logging import getLogger, basicConfig, INFO
basicConfig(format='%(asctime)s [%(relativeCreated)7.0f] [%(levelname).1s] %(message)s',level=INFO,stream=stderr)
log = getLogger(__name__)
SPARK_BASE = "/spark"
UPDATE_CONF = f"{SPARK_BASE}/etc/flatcar/update.conf"
UPDATE_CONF_BUF = """{{ include "spark.etc.flatcar.update.conf" . }}"""
SPARK_YAML = f"{SPARK_BASE}/var/lib/rancher/k3s/server/manifests/spark.yaml"
SPARK_SKIP = f"{SPARK_YAML}.skip"
CONFIG_YAML = f"{SPARK_BASE}/etc/rancher/k3s/config.yaml"
CONFIG_YAML_BUF = """
token: SOME_TOKEN
agent-token: ANOTHER_TOKEN
#token: ANOTHER_TOKEN
secrets-encryption: true
secrets-encryption-provider: secretbox
#tls-san: ???
flannel-backend: "wireguard-native"
cluster-init: true
#server: {{ $k3s_url | quote }}
node-ip: "192.168.56.51"
flannel-iface: "enp0s8"
"""
def main():
    log.info("==== spark begin ====")
    # https://docs.k3s.io/installation/packaged-components
    # don't let spark.yaml run on k3s(etcd)
    with open(SPARK_SKIP,"w") as fp:
        fchmod(fp.fileno(),0o600)
    log.info("empty '{0}' created".format(SPARK_SKIP))
    with open(CONFIG_YAML,"w") as fp:
        fp.write(CONFIG_YAML_BUF.strip() + "\n")
        fchmod(fp.fileno(),0o600)
    log.info("'{0}' created".format(CONFIG_YAML))
    for entry in ["containerd","docker"]:
        fname = "{0}/etc/extensions/{1}-flatcar.raw".format(SPARK_BASE,entry)
        try:
            if not islink(fname):
                log.warning("skipping '{0}': it's not a symlink".format(fname))
                continue
            log.info("(before) {0}: {1}".format(fname,readlink(fname)))
            unlink(fname)
            symlink("/dev/null",fname)
            log.info(" (after) {0}: {1}".format(fname,readlink(fname)))
        except FileNotFoundError as e:
            log.warning("skipping '{0}': {1}".format(fname,str(e)))
    if isfile(UPDATE_CONF):
        with open(UPDATE_CONF,"w") as fp:
            fp.write(UPDATE_CONF_BUF.strip() + "\n")
            fchmod(fp.fileno(),0o644)
        log.info("'{0}' updated".format(UPDATE_CONF))
    else:
        log.warning("skipping '{0}': it's not a regular file".format(UPDATE_CONF))
    log.info("---- spark end ----")
    # #chroot /spark systemd-run bash -c 'sleep 1 ; systemctl reboot'
    # # Failed to connect to system scope bus via local transport: No data available
    # chroot "${SPARK_BASE}" systemctl reboot
if __name__ == "__main__":
    main()
{{- end }}
