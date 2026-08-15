#!/bin/bash

function error_and_exit {
  echo "error: $1" >&2
  exit 1
}

CLUSTER="$1"

[ "${CLUSTER}" = "dev" -o "${CLUSTER}" = "prod" ] || error_and_exit "unknown cluster '${CLUSTER}' (valid: 'dev' or 'prod')"

set -x -e -o pipefail

flux bootstrap github \
  --token-auth \
  --owner=sfmunoz \
  --repository=i12e \
  --path=clusters/${CLUSTER} \
  --branch=main \
  --private=false \
  --personal=true \
  --author-name "flux-${CLUSTER}-bot" \
  --author-email "46285520+sfmunoz@users.noreply.github.com"
