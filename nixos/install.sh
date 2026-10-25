#!/usr/bin/env bash

function error_and_exit {
  echo "error: $1" >&2
  exit 1
}

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

for FNAME in network.nix authorized_keys.nix; do
  sops decrypt "${SOURCE_SECRETS}/nixos/etc/nixos/vm/${HNAME}.yaml" |
    yq -r -e ".${FNAME%.nix}_nix" |
    sudo tee "${TARGET}/etc/nixos/${FNAME}" >/dev/null
done

for FNAME in flake.nix; do
  sed -e "s;@@HOSTNAME@@;${HNAME};g" "${SOURCE}/etc/nixos/${FNAME}" |
    sudo tee "${TARGET}/etc/nixos/${FNAME}" >/dev/null
done

for FNAME in configuration.nix users.nix packages.nix; do
  sudo cp -v "${SOURCE}/etc/nixos/${FNAME}" "${TARGET}/etc/nixos/${FNAME}"
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
