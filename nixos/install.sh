#!/usr/bin/env bash

set -x -e -o pipefail

[ "$HNAME" = "" ] && HNAME="$(hostname -s)"
SOURCE="$(realpath "$(dirname "$0")")"
SOURCE_SECRETS="$(realpath "${SOURCE}/../../i12e-secrets")"

TARGET=""
[ "$(basename "$0")" = "install.sh" ] && TARGET="/mnt"

sudo mkdir -p "${TARGET}/etc/nixos"
sudo chown 0:0 "${TARGET}/etc/nixos"
sudo chmod 0700 "${TARGET}/etc/nixos"

# Use "ssh-to-age < /etc/ssh/ssh_host_ed25519_key.pub" to get the AGE public key
env | grep -q "^SOPS_AGE_" || export SOPS_AGE_KEY_CMD="sudo ssh-to-age -private-key -i /etc/ssh/ssh_host_ed25519_key"

for f in network.nix; do
  sops decrypt "${SOURCE_SECRETS}/nixos/etc/nixos/vm/${HNAME}.yaml" |
    yq -r ".${f%.nix}_nix" |
    sudo tee "${TARGET}/etc/nixos/$f" >/dev/null
done

sed -e "s;@@HOSTNAME@@;${HNAME};g" "${SOURCE}/etc/nixos/flake.nix" |
  sudo tee "${TARGET}/etc/nixos/flake.nix" >/dev/null

for f in configuration.nix users.nix packages.nix; do
  sudo cp -v "${SOURCE}/etc/nixos/$f" "${TARGET}/etc/nixos/$f"
done

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
