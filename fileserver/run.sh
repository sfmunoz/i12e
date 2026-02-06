#!/bin/bash
#
# Ref:
#   https://github.com/rustfs/rustfs/
#

set -e -o pipefail
cd "$(dirname "$0")"
set -x
sudo -u root mkdir -p data logs
sudo -u root chown -R 10001:10001 data logs
exec docker run \
  -it \
  --rm \
  --name rustfs \
  -p 9000:9000 \
  -p 9001:9001 \
  -v $(pwd)/data:/data \
  -v $(pwd)/logs:/logs \
  rustfs/rustfs:latest

