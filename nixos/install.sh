#!/bin/sh

set -x -e

SOURCE="$(realpath "$(dirname "$0")")"
TARGET=""

mountpoint -q /mnt && TARGET="/mnt"

sudo mkdir -p "${TARGET}/etc/nixos"

for f in configuration.nix flake.nix users.nix; do
  sudo cp -v "${SOURCE}/etc/nixos/$f" "${TARGET}/etc/nixos/$f"
done
