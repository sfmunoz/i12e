#!/bin/sh
# curl -sfL https://raw.githubusercontent.com/sfmunoz/i12e/refs/tags/${VERSION}/scripts/install.sh | sh -
VERSION="v0.0.6"
SHA256="05ea4a7117b43e18ced82bcce5bf23b5a444585c0cd36193bd754d914852f5fc"
set -e -o pipefail
TFILE="$(mktemp)"
trap 'echo "exit code: $?" ; rm -f $TFILE' EXIT
set -x
curl -sfL -o $TFILE https://github.com/sfmunoz/i12e/releases/download/${VERSION}/i12e-flatcar.raw
echo "$SHA256  $TFILE" | sha256sum -c
[ -f /etc/systemd/system/i12e.service ] && systemctl stop i12e
mkdir -p /etc/extensions
[ -f /etc/extensions/i12e-flatcar.raw ] && rm -fv /etc/extensions/i12e-flatcar.raw
install -o root -g root -m 0644 $TFILE /etc/extensions/i12e-flatcar.raw
systemd-sysext refresh
systemd-sysext status
I12E_INSTALL=1 /usr/bin/i12e
