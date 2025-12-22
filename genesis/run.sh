#!/bin/bash

function error_and_exit {
  echo "error: $1" 1>&2
  exit 1
}

set -e -o pipefail

[ "$SSH_PUBKEY_FILE" = "" ] && SSH_PUBKEY_FILE="${HOME}/.ssh/id_rsa.pub"

[ -f "$SSH_PUBKEY_FILE" ] || error_and_exit "SSH_PUBKEY_FILE='$SSH_PUBKEY_FILE' file doesn't exist"

[[ "$@" = "" ]] && set -- python3 -m genesis

set -x

cd "$(dirname "$0")"

exec docker run -it --rm \
  -v "${SSH_PUBKEY_FILE}:/ssh_authorized_key:ro" \
  -v ./app:/app/genesis:ro \
  -e PYTHONUNBUFFERED=1 \
  -e PYTHONPATH=/app \
  -e GENESIS_LOCAL=1 \
  ghcr.io/sfmunoz/k8s-bulk:v1.6.0 \
  "$@"
