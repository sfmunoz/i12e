#!/bin/bash

IP="192.168.56.51"

set -e -o pipefail

cd "$(dirname "$0")"

set -x

FNAME="spark.yaml.$(date +%Y%m%d_%H%M%S)"

scp spark.yaml core@${IP}:$FNAME

ssh core@${IP} "set -e -x -o pipefail ; sudo chown -v 0:0 $FNAME ; sudo mv -v $FNAME /var/lib/rancher/k3s/server/manifests/spark.yaml"

