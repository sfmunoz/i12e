#!/bin/bash

function error_and_exit {
  echo "error: $1" >&2
  exit 1
}

set -e -o pipefail

FNAME="$1"

[ "$FNAME" = "" ] && error_and_exit "filename argument must be provided"

if [ ! -f "$FNAME" ]
then
  set -x
  touch "$FNAME"
  chmod 600 "$FNAME"
  wg genkey > "$FNAME"
  { set +x; } 2> /dev/null
fi

cat "$FNAME"
