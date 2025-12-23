#!/bin/bash

function error_and_exit {
  echo "error: $1" 1>&2
  exit 1
}

[ "$GENESIS_OUTPUT" = "" ] && GENESIS_OUTPUT="bash"
[ "$GENESIS_TARGET" = "" ] && GENESIS_TARGET="192.168.56.51"

set -e -o pipefail

[ "$SSH_PUBKEY_FILE" = "" ] && SSH_PUBKEY_FILE="${HOME}/.ssh/id_rsa.pub"

[ -f "$SSH_PUBKEY_FILE" ] || error_and_exit "SSH_PUBKEY_FILE='$SSH_PUBKEY_FILE' file doesn't exist"

# "docker run -t" interferes with output capturing (e.g. jq weird indent behaviour)
case "$1" in
  sh|python3) T_OPT="t" ;;
  *) T_OPT="" ;;
esac

[[ "$@" = "" ]] && set -- python3 -m genesis

#set -x

cd "$(dirname "$0")"

exec docker run -i$T_OPT --rm \
  -v "${SSH_PUBKEY_FILE}:/ssh_authorized_key:ro" \
  -v ./app:/app/genesis:ro \
  -e PYTHONUNBUFFERED=1 \
  -e PYTHONPATH=/app \
  -e GENESIS_OUTPUT=$GENESIS_OUTPUT \
  -e GENESIS_TARGET=$GENESIS_TARGET \
  ghcr.io/sfmunoz/k8s-bulk:v1.6.0 \
  "$@"
