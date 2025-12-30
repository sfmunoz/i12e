#!/bin/bash
#
# https://github.com/sfmunoz/i12e/issues/70
#
# Proof-of-concept:
#
#   - i12e-flatcar extension generation
#   - Only tmux is included together with required 'libutempter.so*'
#
# Tested on:
#
#   - Source: Linux Mint 22.2 Cinnamon
#   - Target: Flatcar 4459.2.2
#
# It works despite the following warning:
#
#   $ tmux
#   tmux: /usr/lib64/libtinfo.so.6: no version information available (required by tmux)
#

[ "$TARGET" = "" ] && TARGET="192.168.56.51"
D="build/i12e-flatcar"

set -e -o pipefail
cd "$(dirname "$0")/.."
set -x
rm -rf build
umask 022
mkdir -p $D/usr/bin $D/usr/lib/extension-release.d $D/usr/lib64
cp /usr/bin/tmux $D/usr/bin
cp -a /lib/x86_64-linux-gnu/libutempter.so* $D/usr/lib64
cat << __EOF > $D/usr/lib/extension-release.d/extension-release.i12e-flatcar
ID=flatcar
VERSION_ID=4459.2.2
ARCHITECTURE=x86-64
__EOF
mksquashfs $D $D.raw -noappend -comp zstd -all-root
ssh "core@${TARGET}" "sudo systemd-sysext status"
ssh "core@${TARGET}" "sudo systemd-sysext unmerge"
ssh "core@${TARGET}" "sudo rm -fv /etc/extensions/i12e-flatcar.raw"
ssh "core@${TARGET}" "sudo bash -c 'cat > /etc/extensions/i12e-flatcar.raw'" < $D.raw
ssh "core@${TARGET}" "sudo systemd-sysext refresh"
ssh "core@${TARGET}" "sudo systemd-sysext status"
