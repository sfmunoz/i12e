{ config, lib, pkgs, ... }:
let
  authorizedKeys = import ./authorized_keys.nix;   # rendered by 'install.sh' from 'i12e-secrets'
in
{
  users.users.i12e = {
    isNormalUser = true;
    extraGroups = [ "wheel" ]; # Enable sudo
    openssh.authorizedKeys.keys = authorizedKeys;
  #   packages = with pkgs; [
  #     tree
  #   ];
  };
  users.users.root.openssh.authorizedKeys.keys = authorizedKeys;
}
