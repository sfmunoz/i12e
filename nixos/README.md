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

### 2. Debian software install (as `i12e`)

`htop` is optional:

```
sudo apt update
sudo apt install -y xz-utils curl ca-certificates age gh htop
```

### 3. Sops install (as `i12e`)

Ref: https://github.com/getsops/sops/releases

```
curl -LO https://github.com/getsops/sops/releases/download/v3.13.3/sops-v3.13.3.linux.amd64
sudo mv sops-v3.13.3.linux.amd64 /usr/local/bin/sops
sudo chown root:root /usr/local/bin/sops
sudo chmod 0755 /usr/local/bin/sops
```

### 4. Nix install (as `i12e`)

Passwordless `sudo` required (granted by step 1):

```
curl -L https://nixos.org/nix/install | sh -s -- --daemon
```

Notice: Nix won't work in active shell sessions until it is restarted.

### 5. Nix experimental features (as `root`)

```
printf 'experimental-features = nix-command flakes\n' >> /etc/nix/nix.conf
systemctl restart nix-daemon
```

### 6. NixOS install tools (as `root`)

```
nix profile add github:NixOS/nixpkgs/nixos-26.05#nixos-install-tools
```

Pinned to the `nixos-26.05` release branch on purpose: a bare `nixpkgs#` is
resolved through the flake registry, which points at `nixpkgs-unstable`, and
would install `26.11pre` tools (and with them a `26.11pre` system).

### 7. Target FS (as `root`)

Swap is not needed/wanted:

```
wipefs -a /dev/sda
parted /dev/sda -- mklabel msdos
parted /dev/sda -- mkpart primary 1MB 100%
parted /dev/sda -- set 1 boot on
mkfs.ext4 -L nixos /dev/sda1
mount /dev/disk/by-label/nixos /mnt
```

### 8. Network config (as `root`)

Create `/mnt/etc/nixos/network.nix` with the following content (must be adapted):

```
{ config, lib, pkgs, ... }:
{
  networking = {
    hostName = "vm";
    useDHCP = false;
    usePredictableInterfaceNames = false;   # enp0s3 -> eth0
    interfaces.eth0 = {
      ipv4.addresses = [ { address = "192.168.56.10"; prefixLength = 24; } ];
      ipv6.addresses = [ { address = "fd00::10"; prefixLength = 64; } ];
      tempAddress = "enabled";
    };
    defaultGateway = { address = "192.168.56.1"; interface = "eth0"; };
    defaultGateway6 = { address = "fd00::1"; interface = "eth0"; };
    nameservers = [ "192.168.56.1" "fd00::1" "1.1.1.1" "2606:4700:4700::1111" ];
    firewall.allowedTCPPorts = [ 22 ];
  };
}
```

### 9. NixOS config (as `root`)

```
nixos-generate-config --root /mnt
```

Changes to `/mnt/etc/nixos/configuration.nix` if diff format:

```diff
   imports =
     [ # Include the results of the hardware scan.
       ./hardware-configuration.nix
+      ./network.nix
     ];

   # Use the GRUB 2 boot loader.
```
```diff
-  # boot.loader.grub.device = "/dev/sda"; # or "nodev" for efi only
+  boot.loader.grub.device = "/dev/sda"; # or "nodev" for efi only
```
```diff
-  networking.networkmanager.enable = true;
+  #networking.networkmanager.enable = true;
```
```diff
-  # time.timeZone = "Europe/Amsterdam";
+  time.timeZone = "UTC";
```
```diff
+  security.sudo.wheelNeedsPassword = false;
```
```diff
-  # users.users.alice = {
-  #   isNormalUser = true;
-  #   extraGroups = [ "wheel" ]; # Enable ‘sudo’ for the user.
+  users.users.i12e = {
+    isNormalUser = true;
+    extraGroups = [ "wheel" ]; # Enable ‘sudo’ for the user.
   #   packages = with pkgs; [
   #     tree
   #   ];
-  # };
+  };
```
```diff
-  # services.openssh.enable = true;
+  services.openssh.enable = true;
```
The following keeps the installed system able to run `nixos-rebuild --flake` (step 11).
```diff
+  nix.settings.experimental-features = [ "nix-command" "flakes" ];
```

### 10. NixOS flake (as `root`)

Create `/mnt/etc/nixos/flake.nix` with the following content to pin the
`nixpkgs` revision the target system is built from. `nixos-install` will then
ignore `NIX_PATH` and the `nixpkgs` channels, and install exactly the release
declared here:

```
{
  inputs.nixpkgs.url = "github:NixOS/nixpkgs/nixos-26.05";

  outputs = { nixpkgs, ... }: {
    nixosConfigurations.vm = nixpkgs.lib.nixosSystem {
      system = "x86_64-linux";
      modules = [ ./configuration.nix ];
    };
  };
}
```

Notes:

- `vm` is the configuration name used in step 11 (`--flake .#vm`).
- The lock file `/mnt/etc/nixos/flake.lock` is generated during installation; keep it in the
repo so the machine can be rebuilt from the same inputs.

### 11. NixOS install (as `root`)

```
nixos-install --flake /mnt/etc/nixos#vm
```

### 12. Password setup (as `root`)

```
nixos-enter --root /mnt
passwd i12e
exit
```

### 13. Reboot (as `root`)

```
umount /mnt
reboot
```

## Day-two workflow

- `nixos-rebuild` implies `--flake /etc/nixos` when `/etc/nixos/flake.nix` exists
  (`--no-flake` disables it).
- The implied configuration name is the hostname: `nixosConfigurations.<hostName>`.
  Keep `networking.hostName = "vm"` in sync with `nixosConfigurations.vm`, or be
  explicit with `--flake /etc/nixos#vm`.
- Edit and rebuild:

```
sudoedit /etc/nixos/configuration.nix
sudo nixos-rebuild switch
```

- `/etc/nixos/flake.lock` is written by `nixos-install`; keep it in the repo so
  rebuilds use the same inputs. Regenerate with `sudo nix flake lock /etc/nixos`.
- `/etc/nixos` is not a git repo: the whole directory is copied into the store on
  every rebuild, so stray files count and `.gitignore` is irrelevant.
- `nix.settings.experimental-features` must be enabled, or `nixos-rebuild` fails
  with `experimental Nix feature 'nix-command' is disabled`.

