#!/bin/bash

FKEY="/var/lib/rancher/k3s/server/tls/client-ca.key"
FCERT="/var/lib/rancher/k3s/server/tls/client-ca.crt"

SECRET_NAME="k3s-client-ca"
FOUT="/var/lib/rancher/k3s/server/manifests/${SECRET_NAME}.yaml"

[[ -f "$FKEY" ]] || exit 0
[[ -f "$FCERT" ]] || exit 0
[[ "$FKEY" -nt "$FOUT" || "$FCERT" -nt "$FOUT" ]] || exit 0

set -x -e -o pipefail

(
  k3s kubectl create ns i12e --dry-run=client -o yaml

  echo "---"

  k3s kubectl create secret tls "${SECRET_NAME}" -n i12e --key="$FKEY" --cert="$FCERT" --dry-run=client -o yaml
) >"$FOUT"
