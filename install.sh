#!/bin/sh
# curl -sfL https://raw.githubusercontent.com/sfmunoz/i12e/refs/tags/${VERSION}/install.sh | sh -
VERSION="v0.0.3"
SHA256="3ec7822baf91b0a4a7dac3c58b3e3f8fbff9d5bd61943e461b89df85cb110fb4"
set -e -o pipefail
TFILE="$(mktemp)"
trap 'echo "exit code: $?" ; rm -f $TFILE' EXIT
set -x
curl -sfL -o $TFILE https://github.com/sfmunoz/i12e/releases/download/${VERSION}/i12e
echo "$SHA256  $TFILE" | sha256sum -c
mkdir -p /opt/bin
install -o root -g root -m 0755 $TFILE /opt/bin/i12e
