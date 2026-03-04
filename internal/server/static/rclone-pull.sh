#!/bin/bash
VERSION="v1.73.1"
SHA256SUM="10ffca7640f0483beba9a9c9f4cdb6a9ca30a13c4ab31f2329d09a832802694d"
RCLONE_BIN="/opt/bin/rclone"
set -x
[ -x "${RCLONE_BIN}" -a "$(sha256sum "${RCLONE_BIN}")" = "${SHA256SUM}  ${RCLONE_BIN}" ] && exit 0
set -e -o pipefail
rm -rf "${RCLONE_BIN}"
curl -sfL "https://github.com/rclone/rclone/releases/download/${VERSION}/rclone-${VERSION}-linux-amd64.zip" |
  bsdtar -xO "rclone-${VERSION}-linux-amd64/rclone" >"${RCLONE_BIN}"
chown 0:0 "${RCLONE_BIN}"
chmod 755 "${RCLONE_BIN}"
[ "$(sha256sum "${RCLONE_BIN}")" = "${SHA256SUM}  ${RCLONE_BIN}" ] && exit 0
rm -rf "${RCLONE_BIN}"
echo "error: '${RCLONE_BIN}' download failed" >&2
exit 1
