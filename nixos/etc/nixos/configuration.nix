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
      init.defaultBranch = "main";
      user.name = "sfmunoz";
      user.email = "46285520+sfmunoz@users.noreply.github.com";
    };
  };
  security.sudo.wheelNeedsPassword = false;
  services.openssh.enable = true;
  nix.settings.experimental-features = [ "nix-command" "flakes" ];
  # XXX Most users should NEVER change this value after the initial install
  # Refs: 'man 5 configuration.nix' and https://nixos.org/manual/nixos/stable/options#opt-system.stateVersion
  system.stateVersion = "26.05";
}
