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

### 1. Software install

As `root` (`htop` is optional):

```
apt update
apt install xz-utils htop
```

### 2. Create i12e user

With `sudo` privileges (as `root`):

```
useradd -m -s /bin/bash -g users i12e
echo 'i12e ALL=(ALL:ALL) NOPASSWD:ALL' > /etc/sudoers.d/90-i12e
chmod 440 /etc/sudoers.d/90-i12e
```

### 3. Nix install

As `i12e` user (passwordless `sudo` required):

```
curl -L https://nixos.org/nix/install | sh -s -- --daemon
```

Notice: Nix won't work in active shell sessions until it is restarted.

### 4. Nix experimental features

As `root`:

```
printf 'experimental-features = nix-command flakes\n' >> /etc/nix/nix.conf
systemctl restart nix-daemon
```

### 5. NixOS install tools

As `root`:

```
nix profile add github:NixOS/nixpkgs/nixos-26.05#nixos-install-tools
```

### 6. Target FS

Swap not needed/wanted, as `root`:

```
wipefs -a /dev/sda
parted /dev/sda -- mklabel msdos
parted /dev/sda -- mkpart primary 1MB 100%
parted /dev/sda -- set 1 boot on
mkfs.ext4 -L nixos /dev/sda1
mount /dev/disk/by-label/nixos /mnt
```

### 7. NixOS config

As `root`:

```
nixos-generate-config --root /mnt
cp /mnt/etc/nixos/configuration.nix /mnt/etc/nixos/configuration.nix.orig
vi /mnt/etc/nixos/configuration.nix
```

Differences using `diff -u /mnt/etc/nixos/configuration.nix.orig /mnt/etc/nixos/configuration.nix` (network configuration must be tuned):

```diff
--- /mnt/etc/nixos/configuration.nix.orig       2026-10-23 17:30:31.732000000 +0000
+++ /mnt/etc/nixos/configuration.nix    2026-10-23 17:48:02.276000000 +0000
@@ -16,15 +16,15 @@
   # boot.loader.grub.efiInstallAsRemovable = true;
   # boot.loader.efi.efiSysMountPoint = "/boot/efi";
   # Define on which hard drive you want to install Grub.
-  # boot.loader.grub.device = "/dev/sda"; # or "nodev" for efi only
+  boot.loader.grub.device = "/dev/sda"; # or "nodev" for efi only

   # networking.hostName = "nixos"; # Define your hostname.

   # Configure network connections interactively with nmcli or nmtui.
-  networking.networkmanager.enable = true;
+  #networking.networkmanager.enable = true;

   # Set your time zone.
-  # time.timeZone = "Europe/Amsterdam";
+  time.timeZone = "UTC";

   # Configure network proxy if necessary
   # networking.proxy.default = "http://user:password@proxy:port/";
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
@@ -99,6 +99,21 @@
   # Or disable the firewall altogether.
   # networking.firewall.enable = false;

+  networking = {
+    hostName = "vm";
+    useDHCP = false;
+    usePredictableInterfaceNames = false;   # enp0s3 -> eth0
+    interfaces.eth0 = {
+      ipv4.addresses = [ { address = "192.168.56.10"; prefixLength = 24; } ];
+      ipv6.addresses = [ { address = "fd00::10"; prefixLength = 64; } ];
+      tempAddress = "enabled";
+    };
+    defaultGateway = { address = "192.168.56.1"; interface = "eth0"; };
+    defaultGateway6 = { address = "fd00::1"; interface = "eth0"; };
+    nameservers = [ "192.168.56.1" "fd00::1" "1.1.1.1" "2606:4700:4700::1111" ];
+    firewall.allowedTCPPorts = [ 22 ];
+  };
+
   # Copy the NixOS configuration file and link it from the resulting system
   # (/run/current-system/configuration.nix). This is useful in case you
   # accidentally delete configuration.nix.
```

Cleanup (as `root`):

```
rm /mnt/etc/nixos/configuration.nix.orig
```

### 8. NixOS install

As `root`:

```
nixos-install
```

### 9. Password setup

As `root`:

```
nixos-enter --root /mnt
passwd i12e
exit
```

### 10. Reboot

As `root` (`umount /mnt` is optional):

```
umount /mnt
reboot
```
