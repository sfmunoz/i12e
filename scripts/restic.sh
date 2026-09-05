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
SOPS="${DNAME}/sops.sh"
RCLONE="${DNAME}/rclone.sh"

export RESTIC_PASSWORD="$("${SOPS}" restic "${APP}")"
export RESTIC_REPOSITORY="rclone:rem:${NAMESPACE}/${APP}"
export RESTIC_CACHE_DIR=/dev/null

exec restic --option rclone.program="${RCLONE}" "$@"
