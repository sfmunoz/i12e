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
i12e-conf)
  (
    sops decrypt "${I12E_SECRETS}/clusters/${I12E_ENV}/flux-system/flux-system.yaml" |
      yq -y '.stringData | { "github_token": .password }'
    sops decrypt "${I12E_SECRETS}/clusters/${I12E_ENV}/flux-system/sops-age.yaml" |
      yq -y '.stringData | { "age_key": ."age.agekey" }'
  ) |
    yq -y '{ "flux": . }'
  sops decrypt "${I12E_SECRETS}/clusters/${I12E_ENV}/i12e/butane.yaml" |
    yq -y '.stringData |
      to_entries |
      map(if .key == "ssh_authorized_keys" then .value |= (split("\n") | map(select(.!=""))) else . end) |
      from_entries |
      { "butane": . }'
  sops decrypt "${I12E_SECRETS}/clusters/${I12E_ENV}/i12e/artifact.yaml" |
    yq -y '.stringData |
      to_entries |
      map(if .key == "port_knocking" then .value |= (split(",") | map(tonumber)) else . end) |
      from_entries |
      { "artifact": . }'
  sops decrypt "${I12E_SECRETS}/clusters/${I12E_ENV}/i12e/mesh.yaml" |
    yq -y '.stringData | { "mesh": . }'
  sops decrypt "${I12E_SECRETS}/talos/${I12E_ENV}/wge.yaml" |
    awk '/^\-\-\-$/ { exit 0 } { print }' |
    "${DNAME}/talos-wge.py" |
    awk 'BEGIN { print "wg_conf: |" } { print($0=="" ? "" : "  "$0) }'
  exit $?
  ;;
talos-secrets | talos-wgm | talos-wge)
  sops decrypt "${I12E_SECRETS}/talos/${I12E_ENV}/${BLOCK#talos-}.yaml"
  exit $?
  ;;
sops-age)
  sops decrypt "${I12E_SECRETS}/clusters/${I12E_ENV}/flux-system/sops-age.yaml"
  exit $?
  ;;
*)
  error_and_exit "unknown '${BLOCK}' block"
  ;;
esac
