#!/bin/bash

function error_and_exit {
  echo "error: $1" 1>&2
  exit 1
}

set -e -o pipefail

cd "$(dirname "$0")"

[ "$I12E_ENV" != "prod" ] && I12E_ENV="dev"

SFILE="../secrets-${I12E_ENV}.yaml"
[ -f "$SFILE" ] || error_and_exit "cannot find '${SFILE}' file"
I12E_SECRETS_YAML="$(sops decrypt "$SFILE" | gzip | base64 -w 0)"
[ "$I12E_SECRETS_YAML" = "" ] && error_and_exit "cannot get content from '${SFILE}' file"

[ "$RCLONE_CONFIG_FILE" = "" ] && RCLONE_CONFIG_FILE="${HOME}/.config/rclone/rclone.conf"
[ -f "$RCLONE_CONFIG_FILE" ] || error_and_exit "RCLONE_CONFIG_FILE='$RCLONE_CONFIG_FILE' file doesn't exist"

# "docker run -t" interferes with output capturing (e.g. jq weird indent behaviour)
case "$1" in
  sh|python3) T_OPT="t" ;;
  *) T_OPT="" ;;
esac

exec docker run -i$T_OPT --rm \
  -v "${RCLONE_CONFIG_FILE}:/root/.config/rclone/rclone.conf:ro" \
  -v ./app:/app/genesis:ro \
  -e "I12E_SECRETS_YAML=$I12E_SECRETS_YAML" \
  -e "PYTHONUNBUFFERED=1" \
  -e "PYTHONPATH=/app" \
  ghcr.io/sfmunoz/k8s-bulk:v1.8.0 \
  python3 -m genesis "$@"
