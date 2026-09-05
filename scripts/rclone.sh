#!/bin/bash

set -e -o pipefail

DNAME="$(dirname "$0")"
SOPS="${DNAME}/sops.sh"
TDIR="${XDG_RUNTIME_DIR}"

[ "$TDIR" = "" ] && TDIR="/dev/shm"

RCLONE_CONFIG="$(mktemp -p "$TDIR" rclone-$(id -u).XXXXXXXXXX.conf)"
chmod 600 "${RCLONE_CONFIG}"
$SOPS rclone >"$RCLONE_CONFIG"

exec 3<"$RCLONE_CONFIG"
rm -f "$RCLONE_CONFIG"

export RCLONE_CONFIG="/dev/fd/3"

exec rclone "$@"
