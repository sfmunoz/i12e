#!/bin/bash
#
# Refs:
#   https://github.com/rustfs/rustfs/
#   https://rustfs.com/download/?platform=linux
#     $ curl -O https://github.com/rustfs/rustfs/releases/download/1.0.0-alpha.82/rustfs-linux-x86_64-gnu-latest.zip
#     $ unzip rustfs-linux-x86_64-musl.zip
#     $ ./rustfs --version
#

VERSION="1.0.0-beta.11"
SHA256SUM="ae77d173c520d523b7fd48fea93e8f8c4ad0c496b16cb47f14660677621424d5"

set -x -e -o pipefail
mkdir -p "${HOME}/.i12e/rustfs/data"
cd "${HOME}/.i12e/rustfs"
{ set +x; } 2>/dev/null

if [ ! -f rustfs ]; then
  set -x
  curl -LO https://github.com/rustfs/rustfs/releases/download/${VERSION}/rustfs-linux-x86_64-gnu-v${VERSION}.zip
  echo "${SHA256SUM}  rustfs-linux-x86_64-gnu-v${VERSION}.zip" | sha256sum -c
  unzip rustfs-linux-x86_64-gnu-v${VERSION}.zip
  rm rustfs-linux-x86_64-gnu-v${VERSION}.zip
  { set +x; } 2>/dev/null
fi

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
