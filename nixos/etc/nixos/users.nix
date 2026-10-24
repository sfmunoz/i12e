{ config, lib, pkgs, ... }:
{
  users.users.i12e = {
    isNormalUser = true;
    extraGroups = [ "wheel" ]; # Enable sudo
  #   packages = with pkgs; [
  #     tree
  #   ];
  };
}
