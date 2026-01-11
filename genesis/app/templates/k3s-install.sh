#!/bin/bash
set -e -o pipefail
[ -f "/etc/systemd/system/k3s.service" ] && exit 0
echo "Installing k3s..."
set -x
export INSTALL_K3S_VERSION="v1.34.3+k3s1"
curl -sfL https://get.k3s.io | sh -s
