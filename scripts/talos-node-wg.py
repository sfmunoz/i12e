#!/usr/bin/env python3

# Input .... sops.sh talos-mesh
# Args ..... idx
# Output ... config for idx node

import subprocess
import sys

import yaml


def pubkey(privateKey):
    result = subprocess.run(
        ["wg", "pubkey"],
        input=privateKey + "\n",
        check=True,
        capture_output=True,
        text=True,
    )
    return result.stdout.strip()


if __name__ == "__main__":
    i1 = int(sys.argv[1])
    mesh = yaml.safe_load(sys.stdin)["mesh"]
    mesh_len = len(mesh)
    if i1 < 1:
        raise IndexError("index cannot be less than 1")
    if i1 > len(mesh):
        raise IndexError(f"index cannot be larger than {mesh_len}")
    i0 = i1 - 1
    node = mesh[i0]
    ret = {
        "apiVersion": "v1alpha1",
        "kind": "WireguardConfig",
        "name": "wgi",
        "mtu": 1420,
        "up": True,
        "privateKey": node["key"],
        "listenPort": node["port"],
        "addresses": [{"address": f"192.168.186.{i1}/24"}],
    }
    peers = []
    for i, n in enumerate(mesh, 1):
        if i == i1:
            continue
        peers.append(
            {
                "publicKey": pubkey(n["key"]),
                "endpoint": "{}:{}".format(n["ip"], n["port"]),
                "allowedIPs": [f"192.168.186.{i}/32"],
            },
        )
        if len(peers) > 0:
            ret["peers"] = peers
    print(
        yaml.safe_dump(
            ret,
            sort_keys=False,
            default_flow_style=False,
        )
    )
