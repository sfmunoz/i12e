#!/bin/bash

function error_and_exit {
  echo "error: $1" >&2
  exit 1
}

I12E_ENV="$1"
I12E_K8S="$2"

[ "${I12E_ENV}" = "dev" -o "${I12E_ENV}" = "prod" ] || error_and_exit "unknown I12E_ENV '${I12E_ENV}' (valid: 'dev' or 'prod')"
[ "${I12E_K8S}" = "k3s" -o "${I12E_K8S}" = "talos" ] || error_and_exit "unknown I12E_K8S '${I12E_K8S}' (valid: 'k3s' or 'talos')"

set -e -o pipefail

DNAME="$(dirname "$0")"
SOPS="${DNAME}/sops.sh"

set -x

kubectl get ns flux-system || kubectl create ns flux-system
kubectl get ns flux-system
kubectl get secret -n flux-system sops-age || kubectl apply -f <("$SOPS" sops-age)
kubectl get secret -n flux-system sops-age

flux bootstrap github \
  --token-auth \
  --owner=sfmunoz \
  --repository=i12e \
  --path=clusters/${I12E_ENV}/${I12E_K8S} \
  --branch=main \
  --private=false \
  --personal=true \
  --author-name "flux-${I12E_ENV}-${I12E_K8S}-bot" \
  --author-email "46285520+sfmunoz@users.noreply.github.com" \
  --components-extra=source-watcher
