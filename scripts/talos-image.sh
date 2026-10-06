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

{ set +x; } 2>/dev/null

cat <<__EOF

Grml box:

  # systemctl start ssh
  # passwd

Use the following command to push the image:

  $ ssh root@<VPS_IP_ADDRESS> "zstd -d | dd of=/dev/sdX bs=4M conv=fsync status=progress" < _out/metal-amd64.raw.zst

Once Talos has booted bootstrap is required:

  $ talos-setup.sh talosconfig
  $ talos-setup.sh kubeconfig   (optional but recommended)
  $ talosctl bootstrap --nodes 192.168.186.1

__EOF
