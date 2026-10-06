#!/bin/bash
#
# Ref:
#   https://docs.siderolabs.com/talos/v1.14/platform-specific-installations/boot-assets
#

set -e -o pipefail

cd "$(dirname "$0")/.."

set -x

docker run --rm -t \
  -v ./_out:/out \
  ghcr.io/siderolabs/imager:v1.14.2 \
  metal
