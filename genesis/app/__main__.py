#!/usr/bin/env python3
import sys
from os import getenv
from logging import getLogger, basicConfig, INFO
import kopf
from kubernetes import client, config
from .install import GenesisInstall
from .butane import Butane

basicConfig(format='%(asctime)s [%(relativeCreated)7.0f] [%(levelname).1s] %(message)s (%(module)s:%(lineno)d)',level=INFO,stream=sys.stderr)
log = getLogger(__name__)

genesis_target = getenv("GENESIS_TARGET")

if genesis_target is not None and genesis_target != "":
    Butane(genesis_target).run()
    sys.exit(0)

#config.load_kube_config()
config.load_incluster_config()

class Namespace(object):
    __ns = None
    @classmethod
    def get(cls):
        fname = "/run/secrets/kubernetes.io/serviceaccount/namespace"
        if cls.__ns is None:
            with open(fname,"r") as fp:
                cls.__ns = fp.read()
        if cls.__ns is None:
            raise Exception(f"cannot get namespace from '{fname}'")
        return cls.__ns

@kopf.timer('gdeployments', interval=5.0)
def on_timer(spec, **kwargs):
    try:
        log.info("on_timer()")
        api = client.CoreV1Api()
        pods = api.list_namespaced_pod(Namespace.get())
        for i,pod in enumerate(pods.items):
            log.info("pod={0}: name='{1}', phase='{2}', ip='{3}'".format(i,pod.metadata.name,pod.status.phase,pod.status.pod_ip))
        GenesisInstall().run()
    except Exception as e:
        log.error("error: " + str(e))

def main():
    log.info("==== genesis begin ====")
    kopf.run(namespace=Namespace.get(),peering_name=getenv("VALUES_KOPF_PEERING_NAME"))
    log.info("---- genesis end ----")

if __name__ == "__main__":
    main()
