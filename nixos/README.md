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

Software install (as `root`, `htop` is optional):

```
apt update
apt install xz-utils htop
```

`i12e` user creation with `sudo` privileges (as `root`):

```
useradd -m -s /bin/bash -g users i12e
echo 'i12e ALL=(ALL:ALL) NOPASSWD:ALL' > /etc/sudoers.d/90-i12e
chmod 440 /etc/sudoers.d/90-i12e
```

Nix install (as `i12e` user, `sudo` used):

```
curl -L https://nixos.org/nix/install | sh -s -- --daemon
```

Nix `experimental-features` config (as `root`):

```
printf 'experimental-features = nix-command flakes\n' >> /etc/nix/nix.conf
systemctl restart nix-daemon
```

NixOS install tools (as `root`):

```
nix profile add nixpkgs#nixos-install-tools
```

Disk partition/formatting, swap not needed/wanted (as `root`):

```
parted /dev/sda -- mklabel msdos
parted /dev/sda -- mkpart primary 1MB 100%
parted /dev/sda -- set 1 boot on
mkfs.ext4 -L nixos /dev/sda1
```

Mount (as `root`):

```
mount /dev/disk/by-label/nixos /mnt
```
