#!/bin/bash

FIN="/root/.config/rclone/rclone.conf"
SECRET_NAME="rclone-conf"
FOUT="/var/lib/rancher/k3s/server/manifests/${SECRET_NAME}.yaml"

[[ -f "$FIN" ]] || exit 0
[[ "$FIN" -nt "$FOUT" ]] || exit 0

set -x -e -o pipefail

(
  k3s kubectl create ns i12e --dry-run=client -o yaml

  echo "---"

  k3s kubectl create secret generic "${SECRET_NAME}" \
    --type Opaque \
    -n i12e \
    --from-file=configData="${FIN}" \
    --dry-run=client \
    -o yaml
) >"$FOUT"
