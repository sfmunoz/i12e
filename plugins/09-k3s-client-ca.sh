#!/bin/bash

FKEY="/var/lib/rancher/k3s/server/tls/client-ca.key"
FCERT="/var/lib/rancher/k3s/server/tls/client-ca.crt"

SECRET_NAME="k3s-client-ca"
FOUT="/var/lib/rancher/k3s/server/manifests/${SECRET_NAME}.yaml"

[[ -f "$FKEY" ]] || exit 0
[[ -f "$FCERT" ]] || exit 0
[[ "$FKEY" -nt "$FOUT" || "$FCERT" -nt "$FOUT" ]] || exit 0

umask 077

set -x -e -o pipefail

k3s kubectl create secret tls "${SECRET_NAME}" \
  -n cert-manager \
  --key="$FKEY" \
  --cert="$FCERT" \
  --dry-run=client \
  -o yaml >"$FOUT"
