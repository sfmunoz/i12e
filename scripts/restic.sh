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

DNAME="$(dirname "$0")"

[ "$I12E_SECRETS" = "" ] && I12E_SECRETS="${DNAME}/../../i12e-secrets"

[ "$I12E_ENV" = "" ] && I12E_ENV="dev"

export RESTIC_PASSWORD="$(sops decrypt ${I12E_SECRETS}/clusters/${I12E_ENV}/${APP}/restic-conf.yaml | yq -r .stringData.password)"
export RESTIC_REPOSITORY="rclone:rem:${NAMESPACE}/${APP}"
export RESTIC_CACHE_DIR=/dev/null

set -x
exec restic --option rclone.program="${DNAME}/rclone.sh" "$@"
