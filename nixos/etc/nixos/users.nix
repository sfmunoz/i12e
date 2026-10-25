{ config, lib, pkgs, ... }:
let
  authorizedKeys = import ./authorized_keys.nix;   # rendered by 'install.sh' from 'i12e-secrets'
in
{
  programs.zsh.enable = true;

  users.users.i12e = {
    isNormalUser = true;
    extraGroups = [ "wheel" ]; # Enable sudo
    shell = pkgs.zsh;
    openssh.authorizedKeys.keys = authorizedKeys;
  #   packages = with pkgs; [
  #     tree
  #   ];
  };
  users.users.root = {
    shell = pkgs.zsh;
    openssh.authorizedKeys.keys = authorizedKeys;
  };
}
