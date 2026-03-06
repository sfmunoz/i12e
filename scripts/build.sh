#!/bin/bash

set -e -o pipefail

function error_and_exit {
  echo "error: $1" >&2
  exit 1
}

DIST="dist"
ROOTFS="${DIST}/rootfs"

[ "$GITHUB_REPOSITORY" = "" ] && error_and_exit "undefined 'GITHUB_REPOSITORY'"
REPO_NAME="${GITHUB_REPOSITORY#*/}"

echo "GITHUB_REPOSITORY ... '${GITHUB_REPOSITORY}'"
echo "REPO_NAME ........... '${REPO_NAME}'"

FLATCAR_EXT="${REPO_NAME}-flatcar"
FLATCAR_EXT_RAW="${FLATCAR_EXT}.raw"

set -x

rm -rf "${DIST}"
go test -v ./...
go build -trimpath -buildvcs=false -ldflags="-s -w" -o "${DIST}/${REPO_NAME}" main.go
mkdir -p "${ROOTFS}/usr/bin" "${ROOTFS}/usr/lib/extension-release.d" "${ROOTFS}/usr/lib64"
cp "${DIST}/${REPO_NAME}" "${ROOTFS}/usr/bin/${REPO_NAME}"
cat <<__EOF >"${ROOTFS}/usr/lib/extension-release.d/extension-release.${FLATCAR_EXT}"
ID=flatcar
SYSEXT_LEVEL=1.0
ARCHITECTURE=x86-64
__EOF
cp /usr/bin/tmux "${ROOTFS}/usr/bin/tmux"
cp -a /usr/lib/x86_64-linux-gnu/libutempter.so* "${ROOTFS}/usr/lib64"
mksquashfs "${ROOTFS}" "${DIST}/${FLATCAR_EXT_RAW}" -noappend -comp zstd -all-root

(cd "${DIST}" && sha256sum "${FLATCAR_EXT_RAW}") >"${DIST}/${FLATCAR_EXT_RAW}.sha256"
