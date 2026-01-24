#!/bin/bash

function error_and_exit {
  echo "error: $1" >&2
  exit 1
}

set -e -o pipefail

WG_FNAME="/etc/i12e/wg-privkey"
TS="$(date +%Y%m%d_%H%M%S_%N)"
MACHINE_ID="$(cat /etc/machine-id)"
[ ${#MACHINE_ID} -eq 32 ] || error_and_exit "MACHINE_ID="${MACHINE_ID}" length is '${#MACHINE_ID}' (32 expected)"

if [ ! -f "$WG_FNAME" ]
then
  set -x
  touch "$WG_FNAME"
  chmod 600 "$WG_FNAME"
  wg genkey > "$WG_FNAME"
  { set +x; } 2> /dev/null
fi

# $ wg pubkey < private | base64 -d | xxd -p -c0 | xxd -p -r | base64
WG_PUBKEY="$(wg pubkey < "$WG_FNAME")"
WG_PUBKEY_HEX="$(base64 -d <<< "${WG_PUBKEY}"| xxd -p -c0)"

echo "TS ................ '${TS}'"
echo "MACHINE_ID ........ '${MACHINE_ID}'"
echo "WG_PUBKEY ......... '${WG_PUBKEY}'"
echo "WG_PUBKEY (HEX) ... '${WG_PUBKEY_HEX}'"

set -x

rclone touch "rem:mesh/${MACHINE_ID}/d/${TS}/wgpubkey/${WG_PUBKEY_HEX}"
rclone touch "rem:mesh/${MACHINE_ID}/c/${TS}"
rclone ls "rem:mesh/${MACHINE_ID}"

