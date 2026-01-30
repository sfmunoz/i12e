#!/bin/sh

function error_and_exit {
  echo "error: $1" >&2
  exit 1
}

set -e -o pipefail

case "$1" in
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

