#!/bin/bash

function error_and_exit {
  echo "error: $1" >&2
  exit 1
}

I12E_ENV="$1"

[ "${I12E_ENV}" = "dev" -o "${I12E_ENV}" = "prod" ] || error_and_exit "unknown I12E_ENV '${I12E_ENV}' (valid: 'dev' or 'prod')"

set -x -e -o pipefail

flux bootstrap github \
  --token-auth \
  --owner=sfmunoz \
  --repository=i12e \
  --path=clusters/${I12E_ENV} \
  --branch=main \
  --private=false \
  --personal=true \
  --author-name "flux-${I12E_ENV}-bot" \
  --author-email "46285520+sfmunoz@users.noreply.github.com" \
  --components-extra=source-watcher
