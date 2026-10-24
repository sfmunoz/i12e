#!/bin/sh

set -x -e

[ "$HNAME" = "" ] && HNAME="$(hostname -s)"
SOURCE="$(realpath "$(dirname "$0")")"
TARGET=""

mountpoint -q /mnt && TARGET="/mnt"

sudo mkdir -p "${TARGET}/etc/nixos"

for f in configuration.nix users.nix packages.nix; do
  sudo cp -v "${SOURCE}/etc/nixos/$f" "${TARGET}/etc/nixos/$f"
done

sed -e "s;@@HOSTNAME@@;${HNAME};g" "${SOURCE}/etc/nixos/flake.nix" |
  sudo tee "${TARGET}/etc/nixos/flake.nix" >/dev/null

{ set +x; } 2>/dev/null

echo

if [ "$TARGET" = "" ]; then
  cat <<__EOF
Rebuild command:

  # sudo -i nixos-rebuild switch
  # sudo -i nixos-rebuild --flake /etc/nixos#${HNAME} switch          (default flake)
__EOF
else
  cat <<__EOF
Install command:

  # sudo -i nixos-install --flake ${TARGET}/etc/nixos#${HNAME}
__EOF
fi

echo
