#!/bin/bash
#
# Refs:
#   https://github.com/rustfs/rustfs/
#   https://rustfs.com/download/?platform=linux
#     $ curl -O https://github.com/rustfs/rustfs/releases/download/1.0.0-alpha.82/rustfs-linux-x86_64-gnu-latest.zip
#     $ unzip rustfs-linux-x86_64-musl.zip
#     $ ./rustfs --version
#

VERSION="1.0.0-alpha.83"
SHA256SUM="b3fbf4e0dbdede70fc774719509181229f747d987571815de1f7163d511b1d9f"

set -e -o pipefail
cd "$(dirname "$0")"

if [ ! -f rustfs ]; then
  set -x
  curl -LO https://github.com/rustfs/rustfs/releases/download/${VERSION}/rustfs-linux-x86_64-gnu-v${VERSION}.zip
  echo "${SHA256SUM}  rustfs-linux-x86_64-gnu-v${VERSION}.zip" | sha256sum -c
  unzip rustfs-linux-x86_64-gnu-v${VERSION}.zip
  rm rustfs-linux-x86_64-gnu-v${VERSION}.zip
  { set +x; } 2>/dev/null
fi

set -x
mkdir -p data
{ set +x; } 2>/dev/null

set +e

while true; do
  set -x
  ./rustfs \
    --address 192.168.56.1:9000 \
    --console-address 192.168.56.1:9001 \
    data
  { set +x; } 2>/dev/null
  echo "waiting 1s to try again..."
  sleep 1
done
