# NixOS

[NixOS](https://nixos.org/) configuration

## References

- https://nixos.org/
  - https://nixos.org/manual/nixos/stable/
  - https://search.nixos.org/packages
- https://nix.dev/

## Install

Ref: https://nixos.org/manual/nixos/stable/#sec-installation-manual

Initial status:

- Debian 12.13 on memory
- Empty /dev/sda HD (+100GB)

Std user creation (as `root`):

```
useradd -m -s /bin/bash -g users i12e
echo 'i12e ALL=(ALL:ALL) NOPASSWD:ALL' > /etc/sudoers.d/90-i12e
chmod 440 /etc/sudoers.d/90-i12e
```

Nix install (as `i12e` user, `htop` is optional):

```
sudo apt update
sudo apt install xz-utils htop
curl -L https://nixos.org/nix/install | sh -s -- --daemon
```

Disk partition/formatting (swap not needed/wanted):

```
sudo parted /dev/sda -- mklabel msdos
sudo parted /dev/sda -- mkpart primary 1MB 100%
sudo parted /dev/sda -- set 1 boot on
sudo mkfs.ext4 -L nixos /dev/sda1
```

Mount:

```
sudo mount /dev/disk/by-label/nixos /mnt
```
