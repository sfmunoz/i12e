#!/bin/bash

set -e -o pipefail

APP="$1"

case "$APP" in
anki | forgejo | trilium)
  NAMESPACE="${APP}"
  ;;
*)
  echo
  echo "usage:"
  echo
  echo "  \$ $(basename "$0") anki/forgejo/trilium [restic arguments...]"
  echo
  exit 1
  ;;
esac

shift

if [ "$I12E_SECRETS" = "" ]; then
  cd "$(dirname "$0")"
  I12E_SECRETS=../../i12e-secrets
fi

[ "$I12E_ENV" = "" ] && I12E_ENV="dev"

export RCLONE_CONFIG="$(mktemp -p /dev/shm rclone-$(id -u).XXXXXXXXXX.conf)"
trap "rm -vf '${RCLONE_CONFIG}'" EXIT
chmod 600 "${RCLONE_CONFIG}"

sops decrypt "${I12E_SECRETS}/clusters/${I12E_ENV}/kube-system/rclone-conf.yaml" |
  yq -r .stringData.configData >"$RCLONE_CONFIG"

export RESTIC_PASSWORD="$(sops decrypt ${I12E_SECRETS}/clusters/${I12E_ENV}/${APP}/restic-conf.yaml | yq -r .stringData.password)"
export RESTIC_REPOSITORY="rclone:rem:${NAMESPACE}/${APP}"
export RESTIC_CACHE_DIR=/dev/null

set -x
exec restic "$@"
