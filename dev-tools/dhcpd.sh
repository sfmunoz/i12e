#!/bin/bash

cd "$(dirname "$0")"
while true; do
  set -x
  awk '!/^(#|$)/' dhcpd.conf
  sudo busybox udhcpd -f dhcpd.conf
  { set +x; } 2>/dev/null
  echo "waiting 1s to try again..."
  sleep 1
done
