#!/bin/bash
#
# Ref:
#   https://github.com/rustfs/rustfs/
#

set -e
cd "$(dirname "$0")"
set -x
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
