#!/bin/bash

set -x -e -o pipefail
cd "$(dirname "$0")"
awk '!/^(#|$)/' udhcpd.conf
sudo busybox udhcpd -f udhcpd.conf
