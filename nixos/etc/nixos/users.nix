{ config, lib, pkgs, ... }:
let
  authorizedKeys = [ "@@AUTHORIZED_KEY@@" ];
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
