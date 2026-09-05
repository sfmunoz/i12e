#!/bin/bash

set -e -o pipefail

DNAME="$(dirname "$0")"

[ "$I12E_SECRETS" = "" ] && I12E_SECRETS="${DNAME}/../../i12e-secrets"
[ "$I12E_ENV" = "" ] && I12E_ENV="dev"

TDIR="${XDG_RUNTIME_DIR}"

[ "$TDIR" = "" ] && TDIR="/dev/shm"

export RCLONE_CONFIG="$(mktemp -p "$TDIR" rclone-$(id -u).XXXXXXXXXX.conf)"
trap "rm -vf '${RCLONE_CONFIG}'" EXIT
chmod 600 "${RCLONE_CONFIG}"

sops decrypt "${I12E_SECRETS}/clusters/${I12E_ENV}/kube-system/rclone-conf.yaml" |
  yq -r .stringData.configData >"$RCLONE_CONFIG"

set -x
exec rclone "$@"
