#!/bin/bash

cd "$(dirname "$0")"
while true
do
  set -x
  awk '!/^(#|$)/' udhcpd.conf
  sudo busybox udhcpd -f udhcpd.conf
  { set +x; } 2>/dev/null
  echo "waiting 1s to try again..."
  sleep 1
done
