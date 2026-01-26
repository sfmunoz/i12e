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
esac

