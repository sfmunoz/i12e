#!/bin/bash
#
# Refs:
#   https://github.com/rustfs/rustfs/
#   https://rustfs.com/download/?platform=linux
#     $ curl -O https://github.com/rustfs/rustfs/releases/download/1.0.0-alpha.82/rustfs-linux-x86_64-gnu-latest.zip
#     $ unzip rustfs-linux-x86_64-musl.zip
#     $ ./rustfs --version
#

set -e
cd "$(dirname "$0")"

if [ ! -f rustfs ]; then
  set -x
  curl -LO https://github.com/rustfs/rustfs/releases/download/1.0.0-alpha.83/rustfs-linux-x86_64-gnu-v1.0.0-alpha.83.zip
  echo "b3fbf4e0dbdede70fc774719509181229f747d987571815de1f7163d511b1d9f  rustfs-linux-x86_64-gnu-v1.0.0-alpha.83.zip" | sha256sum -c
  unzip rustfs-linux-x86_64-gnu-v1.0.0-alpha.83.zip
  rm rustfs-linux-x86_64-gnu-v1.0.0-alpha.83.zip
  { set +x; } 2>/dev/null
fi

sudo -u root mkdir -p data logs
sudo -u root chown -R 10001:10001 data logs
{ set +x; } 2>/dev/null
set +e

DATA_FOLDER="$(pwd)/data"
LOGS_FOLDER="$(pwd)/logs"

while true; do
  set -x
  docker run \
    -it \
    --rm \
    --name rustfs \
    -p 127.0.0.1:9000:9000 \
    -p 127.0.0.1:9001:9001 \
    -p 192.168.56.1:9000:9000 \
    -p 192.168.56.1:9001:9001 \
    -v ${DATA_FOLDER}:/data \
    -v ${LOGS_FOLDER}:/logs \
    rustfs/rustfs:latest
  { set +x; } 2>/dev/null
  echo "waiting 1s to try again..."
  sleep 1
done
