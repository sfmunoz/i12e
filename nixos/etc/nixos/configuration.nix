#
# Refs:
#   - man 5 configuration.nix
#   - https://search.nixos.org/options
#   - nixos-help (= https://nixos.org/manual/nixos/stable/)
#

{ config, lib, pkgs, ... }:
{
  imports = [
    ./hardware-configuration.nix
    ./network.nix
    ./users.nix
    ./packages.nix
  ];
  boot.loader.grub.enable = true;
  boot.loader.grub.device = "/dev/sda";
  time.timeZone = "UTC";
  i18n.defaultLocale = "en_US.UTF-8";
  programs.git = {
    enable = true;
    config = {
      alias.ln = "log --graph --decorate -n 15 --pretty=tformat:'%ci %Cgreen%an%Creset %C(auto)%h%Creset%Cred%d%Creset %s'";
      color.ui = "auto";
      core.pager = "cat";
      init.defaultBranch = "main";
      user.name = "sfmunoz";
      user.email = "46285520+sfmunoz@users.noreply.github.com";
    };
  };
  programs.nix-ld.enable = true;
  security.sudo.wheelNeedsPassword = false;
  services.openssh = {
    enable = true;
    settings = {
      # Key-only auth: no account here has a password hash anyway (all locked),
      # so this removes the brute-force surface instead of relying on that luck.
      PasswordAuthentication = false;
      KbdInteractiveAuthentication = false;
      PermitRootLogin = "prohibit-password";
      # Built-in per-source penalties instead of fail2ban: the source *block* is
      # refused, not just one IP, which is what the rotating /24 botnets need.
      PerSourcePenalties = "crash:120s authfail:60s noauth:30s invaliduser:60s grace-exceeded:120s refuseconnection:120s max:3600s";
      PerSourceNetBlockSize = "24:64"; # IPv4 /24, IPv6 /64
      PerSourceMaxStartups = 4;        # max 4 unauthenticated conns per block
    };
  };
  nix.settings.experimental-features = [ "nix-command" "flakes" ];
  # XXX Most users should NEVER change this value after the initial install
  # Refs: 'man 5 configuration.nix' and https://nixos.org/manual/nixos/stable/options#opt-system.stateVersion
  system.stateVersion = "26.05";
}
