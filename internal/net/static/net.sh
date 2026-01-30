#!/bin/sh

function error_and_exit {
  echo "error: $1" >&2
  exit 1
}

set -e -o pipefail

case "$1" in
  node-push)
    NODE_PATH="$2"
    TS="$3"
    WG_PUBKEY="$4"
    WG_ENDPOINT_IP="$5"
    WG_ENDPOINT_PORT="$6"
    [ "$NODE_PATH" = "" ] && error_and_exit "NODE_PATH argument must be provided"
    [ "$TS" = "" ] && error_and_exit "TS argument must be provided"
    [ "$WG_PUBKEY" = "" ] && error_and_exit "WG_PUBKEY argument must be provided"
    [ "$WG_ENDPOINT_IP" = "" ] && error_and_exit "WG_ENDPOINT_IP argument must be provided"
    [ "$WG_ENDPOINT_PORT" = "" ] && error_and_exit "WG_ENDPOINT_PORT argument must be provided"
    rclone touch "rem:mesh/${NODE_PATH}/${TS}/${WG_PUBKEY}/${WG_ENDPOINT_IP}/${WG_ENDPOINT_PORT}"
    rclone lsf "rem:mesh/${NODE_PATH}" | sort -r | awk -v np="${NODE_PATH}" 'BEGIN { print "set -x -e -o pipefail" } { if (NR<2) next ; printf("rclone delete rem:mesh/%s/%s\n",np,$0) }' | bash
  ;;
  node-config)
    WG_IFACE="$2"
    WG_IPINT="$3"
    WG_PORT="$4"
    WG_FNAME="$5"
    [ "$WG_IFACE" = "" ] && error_and_exit "WG_IFACE argument must be provided"
    [ "$WG_IPINT" = "" ] && error_and_exit "WG_IPINT argument must be provided"
    [ "$WG_PORT" = "" ] && error_and_exit "WG_PORT argument must be provided"
    [ "$WG_FNAME" = "" ] && error_and_exit "WG_FNAME argument must be provided"
    ip link set $WG_IFACE down || true
    ip link del $WG_IFACE || true
    ip link add $WG_IFACE type wireguard
    ip addr add ${WG_IPINT}/16 dev $WG_IFACE
    wg set $WG_IFACE listen-port $WG_PORT private-key "${WG_FNAME}"
    ip link set $WG_IFACE up
  ;;
  *)
    error_and_exit "unknown command '$1'"
  ;;
esac

