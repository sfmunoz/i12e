#!/bin/bash
#
# Ref:
#   https://docs.siderolabs.com/talos/v1.14/platform-specific-installations/boot-assets
#

[ "$IMAGE_KIND" = "" ] && IMAGE_KIND="metal"

set -e -o pipefail

DNAME="$(realpath "$(dirname "$0")/..")"
TALOS_SETUP_SH="${DNAME}/scripts/talos-setup.sh"

TARGET="${DNAME}/_out"

set -x

sudo rm -rfv "$TARGET"

mkdir "$TARGET"

"${TALOS_SETUP_SH}" debug-1 >"${TARGET}/debug-1.yaml"

docker run --rm -t \
  -v "${TARGET}:/out" \
  ghcr.io/siderolabs/imager:v1.14.2 \
  "$IMAGE_KIND" \
  --embedded-config-path=/out/debug-1.yaml
