#!/bin/bash

set -e -o pipefail

if [ "$1" = "" ]; then
  echo
  echo "Usage:"
  echo
  echo "  $(basename "$0") target"
  echo
  echo "Example:"
  echo
  echo "  $(basename "$0") 192.168.56.51"
  echo
  exit 1
fi

TARGET="$1"

set -x
ssh "core@${TARGET}" "sudo systemd-sysext status"
ssh "core@${TARGET}" "sudo systemd-sysext unmerge"
ssh "core@${TARGET}" "sudo systemd-sysext status"
ssh "core@${TARGET}" "sudo rm -fv /etc/extensions/i12e-flatcar.raw"
ssh "core@${TARGET}" "sudo bash -c 'cat > /etc/extensions/i12e-flatcar.raw'" <"dist/i12e-flatcar.raw"
ssh "core@${TARGET}" "sudo systemd-sysext refresh"
ssh "core@${TARGET}" "sudo systemd-sysext status"
