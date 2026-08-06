#!/bin/bash
#
# https://github.com/sfmunoz/i12e/issues/276
#

set -x -e -o pipefail

cd "$(dirname "$0")"

helm install \
  cert-manager oci://quay.io/jetstack/charts/cert-manager \
  --version v1.21.1 \
  --namespace cert-manager \
  --create-namespace \
  --set crds.enabled=true
