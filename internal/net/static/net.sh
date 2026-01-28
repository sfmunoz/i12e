#!/bin/sh

function error_and_exit {
  echo "error: $1" >&2
  exit 1
}

set -e -o pipefail

case "$1" in
  priv-key)
    FNAME="$2"
    [ "$FNAME" = "" ] && error_and_exit "filename argument must be provided"
    if [ ! -f "$FNAME" ]
    then
      touch "$FNAME"
      chmod 600 "$FNAME"
      wg genkey > "$FNAME"
    fi
    cat "$FNAME"
  ;;
  pub-key)
    FNAME="$2"
    [ "$FNAME" = "" ] && error_and_exit "filename argument must be provided"
    wg pubkey < /etc/i12e/wg-priv-key
  ;;
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
    rclone touch "rem:mesh/${NODE_PATH}/${TS}/wg/${WG_PUBKEY}/${WG_ENDPOINT_IP}/${WG_ENDPOINT_PORT}"
    rclone touch "rem:mesh/${NODE_PATH}/${TS}/commit"
  ;;
esac

