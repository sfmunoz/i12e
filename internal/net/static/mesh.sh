#!/bin/bash

function error_and_exit {
  echo "error: $1" >&2
  exit 1
}

set -e -o pipefail

WG_IFACE="wg0"
WG_PORT="51820"
IFACE="enp0s8"  # TODO: parametrize this
WG_IPV4="$(ip -j -4 addr show dev $IFACE | jq -r '.[0].addr_info.[0].local')"
[ "$WG_IPV4" = "" ] && error_and_exit "cannot get IPv4 address from IFACE='${IFACE}'"
WG_FNAME="/etc/i12e/wg-privkey"
TS="$(date +%Y%m%d_%H%M%S_%N)"
MACHINE_ID="$(cat /etc/machine-id)"
[ ${#MACHINE_ID} -eq 32 ] || error_and_exit "MACHINE_ID="${MACHINE_ID}" length is '${#MACHINE_ID}' (32 expected)"
WG_IPINT="10.56.$(awk '{ printf("%s.%s\n",strtonum("0x"substr($0,1,2)),strtonum("0x"substr($0,3,2))) }' <<< "${MACHINE_ID}")"

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
echo "WG_IPINT .......... '${WG_IPINT}'"
echo "IFACE ............. '${IFACE}'"
echo "WG_IPV4 ........... '${WG_IPV4}'"
echo "WG_PORT ........... '${WG_PORT}'"
echo "WG_PUBKEY ......... '${WG_PUBKEY}'"
echo "WG_PUBKEY (HEX) ... '${WG_PUBKEY_HEX}'"

set -x
rclone touch "rem:mesh/${MACHINE_ID}/d/${TS}/wgpubkey/${WG_PUBKEY_HEX}"
rclone touch "rem:mesh/${MACHINE_ID}/d/${TS}/wgipv4/${WG_IPV4}"
rclone touch "rem:mesh/${MACHINE_ID}/d/${TS}/wgport/${WG_PORT}"
rclone touch "rem:mesh/${MACHINE_ID}/d/${TS}/wgipint/${WG_IPINT}"
rclone touch "rem:mesh/${MACHINE_ID}/c/${TS}"
{ set +x; } 2> /dev/null

rclone lsf "rem:mesh/${MACHINE_ID}/c" | sort -r | awk -v m_id="${MACHINE_ID}" 'BEGIN { print "set -x -e -o pipefail" } { if (NR<2) next ; printf("rclone delete rem:mesh/%s/c/%s\nrclone delete rem:mesh/%s/d/%s\n",m_id,$0,m_id,$0) }' | bash

set -x
rclone ls "rem:mesh/${MACHINE_ID}" | sort -k2

ip link set $WG_IFACE down || true
ip link del $WG_IFACE || true
ip link add $WG_IFACE type wireguard
ip addr add ${WG_IPINT}/16 dev $WG_IFACE
wg set $WG_IFACE private-key "${WG_FNAME}"
ip link set $WG_IFACE up
