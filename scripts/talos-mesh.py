#!/usr/bin/env python3

import argparse
import ipaddress
import subprocess
from pathlib import Path

import yaml

MAX_NODES = 254


class Mesh:
    def __init__(self, args):
        self.__args = args
        self.__network = ipaddress.ip_network(f"192.168.{args.network}.0/24")

    def __host(self, index):
        return str(self.__network.network_address + index)

    def __generate_key(self):
        result = subprocess.run(
            ["wg", "genkey"],
            check=True,
            capture_output=True,
            text=True,
        )
        return result.stdout.strip()

    def __gen_public_key(self, private_key):
        result = subprocess.run(
            ["wg", "pubkey"],
            input=private_key + "\n",
            check=True,
            capture_output=True,
            text=True,
        )
        return result.stdout.strip()

    def __parse_endpoint(self, value):
        ip, port = value.rsplit(":", 1)
        ipaddress.ip_address(ip)
        port = int(port)
        if not 1 <= port <= 65535:
            raise ValueError
        return ip, port

    def __create_wireguard_structure(self, node, nodes, listen_port):
        peers = []
        for peer in nodes:
            if peer["index"] == node["index"]:
                continue
            peers.append(
                {
                    "publicKey": peer["public_key"],
                    "endpoint": peer["endpoint"],
                    "allowedIPs": [f"{self.__host(peer['index'])}/32"],
                }
            )
        return {
            "apiVersion": "v1alpha1",
            "kind": "WireguardConfig",
            "name": "wgi",
            "privateKey": node["private_key"],
            "listenPort": listen_port,
            "mtu": 1420,
            "peers": peers,
            "up": True,
            "addresses": [{"address": f"{self.__host(node['index'])}/24"}],
        }

    def __render_yaml(self, structure):
        return yaml.safe_dump(
            structure,
            sort_keys=False,
            default_flow_style=False,
        )

    def __render_conf(self, structure):
        lines = ["[Interface]"]
        lines.append(f"PrivateKey = {structure['privateKey']}")
        for address in structure["addresses"]:
            lines.append(f"Address = {address['address']}")
        lines.append(f"ListenPort = {structure['listenPort']}")
        lines.append(f"MTU = {structure['mtu']}")
        for peer in structure["peers"]:
            lines.append("")
            lines.append("[Peer]")
            lines.append(f"PublicKey = {peer['publicKey']}")
            lines.append(f"Endpoint = {peer['endpoint']}")
            lines.append(f"AllowedIPs = {', '.join(peer['allowedIPs'])}")
        lines.append("")
        return "\n".join(lines)

    def __write_conf(self, content, filename):
        # wg-quick .conf files are intentionally left in plaintext: the user runs
        # wg-quick against them directly, so they are never sops-encrypted.
        Path(filename).write_text(content)

    def __encrypt_yaml(self, content, filename):
        # The plaintext only ever travels through the pipe: sops is handed the
        # rendered YAML on stdin and writes the encrypted result itself.
        subprocess.run(
            [
                "sops",
                "encrypt",
                "--filename-override",
                "secrets.yaml",
                "--output",
                str(filename),
            ],
            input=content,
            check=True,
            text=True,
        )

    def run(self):
        endpoints = []
        for value in self.__args.nodes:
            ip, port = self.__parse_endpoint(value)
            endpoints.append(
                {
                    "endpoint": f"{ip}:{port}",
                    "ip": ip,
                    "port": port,
                }
            )
        nodes = []
        for index, endpoint in enumerate(endpoints, start=1):
            private_key = self.__generate_key()
            public_key = self.__gen_public_key(private_key)
            nodes.append(
                {
                    "index": index,
                    "endpoint": endpoint["endpoint"],
                    "port": endpoint["port"],
                    "private_key": private_key,
                    "public_key": public_key,
                }
            )
        for node in nodes:
            structure = self.__create_wireguard_structure(
                node=node,
                nodes=nodes,
                listen_port=node["port"],
            )
            content = self.__render_yaml(structure)
            filename = Path(f"wg{node['index']}.yaml")
            self.__encrypt_yaml(content, filename)
            print(
                f"{filename}: "
                f"{self.__host(node['index'])} "
                f"endpoint={node['endpoint']} "
                f"public-key={node['public_key']}"
            )
            if self.__args.conf:
                self.__write_conf(
                    self.__render_conf(structure),
                    Path(f"wg{node['index']}.conf"),
                )
                print(f"wg{node['index']}.conf: written (plaintext wg-quick)")


if __name__ == "__main__":
    parser = argparse.ArgumentParser(
        description=(
            "Generate Talos WireguardConfig YAML files for a full WireGuard mesh."
        ),
        formatter_class=argparse.ArgumentDefaultsHelpFormatter,
    )

    parser.add_argument(
        "nodes",
        nargs="+",
        metavar="IP:PORT",
        help="WireGuard endpoint(s) for the nodes",
    )

    parser.add_argument(
        "-n",
        "--network",
        type=int,
        default=186,
        metavar="OCTET",
        help="third octet of the WireGuard subnet, enforced to 192.168.<OCTET>.0/24",
    )

    parser.add_argument(
        "-c",
        "--conf",
        action="store_true",
        help=(
            "also write wg-quick compatible wg<N>.conf files; "
            "these are plaintext and never sops-encrypted"
        ),
    )

    args = parser.parse_args()

    if len(args.nodes) > MAX_NODES:
        parser.error(f"too many nodes: maximum is {MAX_NODES} for a /24 subnet")

    Mesh(args).run()
