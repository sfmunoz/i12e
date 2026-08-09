#!/bin/bash
#
# https://github.com/sfmunoz/i12e/issues/276
#

function error_and_exit {
  echo "error: $1" >&2
  exit 1
}

set -e -o pipefail

[ "$1" != "" ] || error_and_exit "OVH credentials file must be provided as first argument with AK (16 chars), AS (32 chars) and CK (32 chars)"

source "$1"

[ "${#AK}" -eq 16 ] || error_and_exit "AK length is not 16"
[ "${#AS}" -eq 32 ] || error_and_exit "AS length is not 32"
[ "${#CK}" -eq 32 ] || error_and_exit "CK length is not 32"

# https://kubernetes.recipes/recipes/security/certmanager-ovh-dns01-wildcard-tls/

kubectl create secret generic ovh-credentials \
  --namespace cert-manager \
  --from-literal=applicationKey=$AK \
  --from-literal=applicationSecret=$AS \
  --from-literal=consumerKey=$CK
