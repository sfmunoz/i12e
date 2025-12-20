{{- define "__main__.py" -}}
{{- $k3s_url := "https://192.168.56.50:6443" -}}
#!/usr/bin/env python3
from sys import stderr
from logging import getLogger, basicConfig, INFO
import kopf
from kubernetes import client, config
from .install import GenesisInstall

basicConfig(format='%(asctime)s [%(relativeCreated)7.0f] [%(levelname).1s] %(message)s (%(module)s:%(lineno)d)',level=INFO,stream=stderr)
log = getLogger(__name__)
#config.load_kube_config()
config.load_incluster_config()

@kopf.timer('gdeployments', interval=5.0)
def on_timer(spec, **kwargs):
    try:
        log.info("on_timer()")
        api = client.CoreV1Api()
        pods = api.list_namespaced_pod("{{ .Release.Namespace }}")
        for i,pod in enumerate(pods.items):
            log.info("pod={0}: name='{1}', phase='{2}', ip='{3}'".format(i,pod.metadata.name,pod.status.phase,pod.status.pod_ip))
        GenesisInstall().run()
    except Exception as e:
        log.error("error: " + str(e))

def main():
    log.info("==== genesis begin ====")
    kopf.run(namespace="{{ .Release.Namespace }}",peering_name="{{ .Values.peering_name }}")
    log.info("---- genesis end ----")

if __name__ == "__main__":
    main()
{{ end }}
