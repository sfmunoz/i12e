#!/bin/bash
#
# Refs:
# - https://www.flatcar.org/docs/latest/provisioning/sysext/
# - https://github.com/sfmunoz/i12e/issues/70
# - https://github.com/sfmunoz/i12e/issues/83
# - https://github.com/sfmunoz/i12e/issues/85
#
# Proof-of-concept:
# - i12e-flatcar extension generation
# - tmux (with required 'libutempter.so*') and rclone are included
#
# Tested on:
# - Source: Linux Mint 22.2 Cinnamon
# - Target: Flatcar 4459.2.2
#
# tmux works despite the following warning:
#   $ tmux
#   tmux: /usr/lib64/libtinfo.so.6: no version information available (required by tmux)
#
# Req on Linux Mint 22.2:
#   # apt install squashfs-tools erofs-utils tmux rclone
#
# FS=squashfs is the default for now since compression is better:
#   [FS=squashfs] ./scripts/i12e-ext.sh ---> du -s build/i12e-flatcar.raw = 17228
#   FS=erofs ./scripts/i12e-ext.sh --------> du -s build/i12e-flatcar.raw = 27448
#

[ "$TARGET" = "" ] && TARGET="192.168.56.51"
[ "$FS" = "" ] && FS="squashfs"

case "$FS" in
 squashfs|erofs) ;;
 *) echo "error: unsupported FS='$FS'; it must be 'squashfs' (default) or 'erofs'" ; exit 1 ;;
esac

D="build/i12e-flatcar"
DRAW="${D}.raw"

set -e -o pipefail
cd "$(dirname "$0")/.."
set -x
make clean
umask 022
mkdir -p $D/usr/bin $D/usr/lib/extension-release.d $D/usr/lib64
make
cp build/i12e $D/usr/bin/i12e
cp /usr/bin/{tmux,rclone} $D/usr/bin
cp -a /usr/lib/x86_64-linux-gnu/libutempter.so* $D/usr/lib64
cat << __EOF > $D/usr/lib/extension-release.d/extension-release.i12e-flatcar
ID=flatcar
SYSEXT_LEVEL=1.0
ARCHITECTURE=x86-64
__EOF
{ set +x; } 2>/dev/null
case "$FS" in
  squashfs)
    set -x
    mksquashfs $D $DRAW -noappend -comp zstd -all-root
    { set +x; } 2>/dev/null
  ;;
  erofs)
    set -x
    mkfs.erofs --all-root -z lz4hc $DRAW $D
    { set +x; } 2>/dev/null
  ;;
esac
set -x
ssh "core@${TARGET}" "sudo systemd-sysext status"
ssh "core@${TARGET}" "sudo systemd-sysext unmerge"
ssh "core@${TARGET}" "sudo systemd-sysext status"
ssh "core@${TARGET}" "sudo rm -fv /etc/extensions/i12e-flatcar.raw"
ssh "core@${TARGET}" "sudo bash -c 'cat > /etc/extensions/i12e-flatcar.raw'" < $DRAW
ssh "core@${TARGET}" "sudo systemd-sysext refresh"
ssh "core@${TARGET}" "sudo systemd-sysext status"
