#!/bin/bash
#
# https://github.com/sfmunoz/i12e/issues/276
#

set -x -e -o pipefail

cd "$(dirname "$0")"

# https://kubernetes.recipes/recipes/security/certmanager-ovh-dns01-wildcard-tls/

helm repo add cert-manager-webhook-ovh https://aureq.github.io/cert-manager-webhook-ovh

helm repo update

helm upgrade --install \
  cert-manager-webhook-ovh \
  cert-manager-webhook-ovh/cert-manager-webhook-ovh \
  --namespace cert-manager \
  --create-namespace \
  --set groupName=acme.example.com
