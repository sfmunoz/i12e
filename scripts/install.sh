#!/bin/sh
# curl -sfL https://raw.githubusercontent.com/sfmunoz/i12e/refs/tags/${VERSION}/scripts/install.sh | sh -
VERSION="v0.0.5"
SHA256="e8983394d46dc642e9474224e094dd5369f3dca0afd2ba0bc8d7abb42ac7f695"
set -e -o pipefail
TFILE="$(mktemp)"
trap 'echo "exit code: $?" ; rm -f $TFILE' EXIT
set -x
curl -sfL -o $TFILE https://github.com/sfmunoz/i12e/releases/download/${VERSION}/i12e
echo "$SHA256  $TFILE" | sha256sum -c
mkdir -p /opt/bin
install -o root -g root -m 0755 $TFILE /opt/bin/i12e
[ -f /etc/systemd/system/i12e.service ] && systemctl stop i12e
I12E_INSTALL=1 /opt/bin/i12e
