#!/bin/bash

set -e -o pipefail

DNAME="$(dirname "$0")"

[ "$I12E_SECRETS" = "" ] && I12E_SECRETS="${DNAME}/../../i12e-secrets"
[ "$I12E_ENV" = "" ] && I12E_ENV="dev"

TDIR="${XDG_RUNTIME_DIR}"

[ "$TDIR" = "" ] && TDIR="/dev/shm"

RCLONE_CONFIG="$(mktemp -p "$TDIR" rclone-$(id -u).XXXXXXXXXX.conf)"
chmod 600 "${RCLONE_CONFIG}"
sops decrypt "${I12E_SECRETS}/clusters/${I12E_ENV}/kube-system/rclone-conf.yaml" |
  yq -r .stringData.configData >"$RCLONE_CONFIG"

exec 3<"$RCLONE_CONFIG"
rm -f "$RCLONE_CONFIG"

export RCLONE_CONFIG="/dev/fd/3"

exec rclone "$@"
