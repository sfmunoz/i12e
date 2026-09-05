#!/bin/bash

set -e -o pipefail

function error_and_exit {
  echo "error: $1" >&2
  exit 1
}

DNAME="$(dirname "$0")"

[ "$I12E_SECRETS" = "" ] && I12E_SECRETS="${DNAME}/../../i12e-secrets"
[ "$I12E_ENV" = "" ] && I12E_ENV="dev"

export I12E_ENV

BLOCK="$1"

case "$BLOCK" in
rclone-conf)
  sops decrypt "${I12E_SECRETS}/clusters/${I12E_ENV}/kube-system/rclone-conf.yaml" | yq -r .stringData.configData
  exit $?
  ;;
restic-conf)
  APP="$2"
  case "$APP" in
  anki | forgejo | trilium)
    sops decrypt ${I12E_SECRETS}/clusters/${I12E_ENV}/${APP}/restic-conf.yaml | yq -r .stringData.password
    exit $?
    ;;
  *)
    error_and_exit "unknown '${APP}' app"
    ;;
  esac
  ;;
*)
  error_and_exit "unknown '${BLOCK}' block"
  ;;
esac
