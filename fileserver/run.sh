#!/bin/bash
#
# Refs:
#   https://github.com/rustfs/rustfs/
#   https://rustfs.com/download/?platform=linux
#     $ curl -O https://github.com/rustfs/rustfs/releases/download/1.0.0-alpha.82/rustfs-linux-x86_64-gnu-latest.zip
#     $ unzip rustfs-linux-x86_64-musl.zip
#     $ ./rustfs --version
#

set -e -o pipefail
cd "$(dirname "$0")"

if [ ! -f rustfs ]; then
  set -x
  curl -LO https://github.com/rustfs/rustfs/releases/download/1.0.0-alpha.83/rustfs-linux-x86_64-gnu-v1.0.0-alpha.83.zip
  echo "b3fbf4e0dbdede70fc774719509181229f747d987571815de1f7163d511b1d9f  rustfs-linux-x86_64-gnu-v1.0.0-alpha.83.zip" | sha256sum -c
  unzip rustfs-linux-x86_64-gnu-v1.0.0-alpha.83.zip
  rm rustfs-linux-x86_64-gnu-v1.0.0-alpha.83.zip
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
