#!/bin/bash

SRC="/etc/i12e/wikijs/wikijs.yaml"
DST_DIR="/var/lib/rancher/k3s/server/manifests"
DST="${DST_DIR}/wikijs.yaml"

[ ! -f "${SRC}" ] && exit 0
[ ! -d "${DST_DIR}" ] && exit 0

if [ -f "${DST}" ]; then
  (cd /etc/i12e/wikijs && sha256sum wikijs.yaml) | (cd /var/lib/rancher/k3s/server/manifests/ && sha256sum -c) >/dev/null 2>&1 && exit 0
fi

set -x -e
install -o root -g root -m 0600 "${SRC}" "${DST}"
