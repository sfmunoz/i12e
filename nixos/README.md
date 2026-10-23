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
wipefs -a /dev/sda
parted /dev/sda -- mklabel msdos
parted /dev/sda -- mkpart primary 1MB 100%
parted /dev/sda -- set 1 boot on
mkfs.ext4 -L nixos /dev/sda1
```

Mount (as `root`):

```
mount /dev/disk/by-label/nixos /mnt
```

Generate NixOS config (as `root`):

```
nixos-generate-config --root /mnt
```

Configuration tuning (as `root`):

```
cp /mnt/etc/nixos/configuration.nix /mnt/etc/nixos/configuration.nix.orig
vi /mnt/etc/nixos/configuration.nix
```

Differences:

```diff
diff -u /mnt/etc/nixos/configuration.nix.orig /mnt/etc/nixos/configuration.nix
--- /mnt/etc/nixos/configuration.nix.orig       2026-10-23 17:06:34.348000000 +0000
+++ /mnt/etc/nixos/configuration.nix    2026-10-23 17:07:12.180000000 +0000
@@ -16,7 +16,7 @@
   # boot.loader.grub.efiInstallAsRemovable = true;
   # boot.loader.efi.efiSysMountPoint = "/boot/efi";
   # Define on which hard drive you want to install Grub.
-  # boot.loader.grub.device = "/dev/sda"; # or "nodev" for efi only
+  boot.loader.grub.device = "/dev/sda"; # or "nodev" for efi only

   # networking.hostName = "nixos"; # Define your hostname.

@@ -63,13 +63,13 @@
   # services.libinput.enable = true;

   # Define a user account. Don't forget to set a password with ‘passwd’.
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

   # programs.firefox.enable = true;

@@ -91,7 +91,7 @@
   # List services that you want to enable:

   # Enable the OpenSSH daemon.
-  # services.openssh.enable = true;
+  services.openssh.enable = true;

   # Open ports in the firewall.
   # networking.firewall.allowedTCPPorts = [ ... ];
```

Cleanup (as `root`):

```
rm /mnt/etc/nixos/configuration.nix.orig
```
