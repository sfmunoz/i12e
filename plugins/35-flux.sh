#!/bin/bash
export FLUX_VERSION="2.9.4" # used by fluxcd.io/install.sh
SHA256SUM="0c91d4bbbc2aa9c84b42608184534437ef92e5f2d6b862e99a11c0bb24ad0941"
FLUX_BIN="/opt/bin/flux"
set -x -e -o pipefail
function flux_install {
  # if+if, no '-a' -> sha256sum would be executed otherwise
  if [ -x "${FLUX_BIN}" ]; then
    if [ "$(sha256sum "${FLUX_BIN}")" = "${SHA256SUM}  ${FLUX_BIN}" ]; then
      return 0
    fi
  fi
  rm -rf "${FLUX_BIN}"
  curl -s https://fluxcd.io/install.sh | bash -s - /opt/bin
  [ "$(sha256sum "${FLUX_BIN}")" = "${SHA256SUM}  ${FLUX_BIN}" ] && return 0
  rm -rf "${FLUX_BIN}"
  echo "error: '${FLUX_BIN}' download failed" >&2
  return 1
}
flux_install || exit $?
exit 0
