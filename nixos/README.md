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

### 1. Create i12e user (as `root`)

`i12e` is granted with `sudo` privileges:

```
useradd -m -s /bin/bash -g users i12e
echo 'i12e ALL=(ALL:ALL) NOPASSWD:ALL' > /etc/sudoers.d/90-i12e
chmod 440 /etc/sudoers.d/90-i12e
```

### 2. Software install (as `i12e`)

Debian packages, `htop` is optional:

```
sudo apt update
sudo apt install -y xz-utils curl ca-certificates age git gh jq yq htop
```

Sops, not available as Debian package (ref: [https://github.com/getsops/sops/releases](https://github.com/getsops/sops/releases)):

```
curl -LO https://github.com/getsops/sops/releases/download/v3.13.3/sops-v3.13.3.linux.amd64
sudo mv sops-v3.13.3.linux.amd64 /usr/local/bin/sops
sudo chown root:root /usr/local/bin/sops
sudo chmod 0755 /usr/local/bin/sops
```

Nix, passwordless `sudo` required (granted by step 1):

```
curl -L https://nixos.org/nix/install | sh -s -- --daemon
```

### 3. Nix experimental features (as `root`)

```
printf 'experimental-features = nix-command flakes\n' >> /etc/nix/nix.conf
systemctl restart nix-daemon
```

### 4. NixOS install tools (as `root`)

```
nix profile add github:NixOS/nixpkgs/nixos-26.05#nixos-install-tools
```

Pinned to the `nixos-26.05` release branch on purpose: a bare `nixpkgs#` is
resolved through the flake registry, which points at `nixpkgs-unstable`, and
would install `26.11pre` tools (and with them a `26.11pre` system).

### 5. Target FS (as `root`)

Swap is not needed/wanted. Use `wipefs -a /dev/sda1`, `wipefs -a /dev/sda2`, ... if those partitions exist:

```
wipefs -a /dev/sda
parted /dev/sda -- mklabel msdos
parted /dev/sda -- mkpart primary 1MB 100%
parted /dev/sda -- set 1 boot on
mkfs.ext4 -L nixos /dev/sda1
mount /dev/disk/by-label/nixos /mnt
```

### 6. NixOS config generation (as `root`)

Create both `/mnt/etc/nixos/configuration.nix` and `/mnt/etc/nixos/hardware-configuration.nix`:

```
nixos-generate-config --root /mnt
```

### 7. Clone and install (as `i12e`)

Sops: `SOPS_AGE_SSH_PRIVATE_KEY_FILE`, `SOPS_AGE_SSH_PRIVATE_KEY_CMD`, `SOPS_AGE_KEY`, `SOPS_AGE_KEY_FILE` or `SOPS_AGE_KEY_CMD` may be used:

```
mkdir ~/src
git clone https://github.com/sfmunoz/i12e.git ~/src/i12e
gh auth login
git clone https://github.com/sfmunoz/i12e-secrets.git ~/src/i12e-secrets
export HNAME="<name-of-the-target-hostname>"
export SOPS_AGE_KEY="..."
~/src/i12e/nixos/install.sh
```

### 8. NixOS install (as `i12e`)

It may be required to run this twice if `error: NAR hash mismatch in input` aborts the process:

```
sudo -i nixos-install --flake /mnt/etc/nixos#${HNAME}
```

### 9. Password setup (as `root`)

```
nixos-enter --root /mnt
passwd i12e
exit
```

### 10. Reboot (as `root`)

```
umount /mnt
reboot
```
