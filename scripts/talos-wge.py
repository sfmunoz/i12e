#!/usr/bin/env python3

import sys

import yaml


class TalosWge:
    def __line(self, template, value):
        return None if value is None else template.format(value)

    def __render_interface(self, document):
        lines = ["[Interface]"]
        lines.append(self.__line("PrivateKey = {0}", document.get("privateKey")))
        lines.append(self.__line("ListenPort = {0}", document.get("listenPort")))
        lines.append(self.__line("FwMark = {0}", document.get("firewallMark")))
        lines.append(self.__line("MTU = {0}", document.get("mtu")))
        addresses = document.get("addresses")
        if addresses:
            lines.append(
                "Address = " + ", ".join(entry["address"] for entry in addresses)
            )
        return "\n".join(line for line in lines if line is not None)

    def __render_peers(self, peers):
        sections = []
        for peer in peers:
            lines = ["[Peer]"]
            lines.append(self.__line("PublicKey = {0}", peer.get("publicKey")))
            lines.append(self.__line("PresharedKey = {0}", peer.get("presharedKey")))
            lines.append(self.__line("Endpoint = {0}", peer.get("endpoint")))
            keepalive = peer.get("persistentKeepaliveInterval")
            if keepalive is not None:
                lines.append(
                    "PersistentKeepalive = " + str(int(str(keepalive).rstrip("s")))
                )
            allowed = peer.get("allowedIPs")
            if allowed:
                lines.append("AllowedIPs = " + ", ".join(allowed))
            sections.append("\n".join(line for line in lines if line is not None))
        return sections

    def __render(self, document):
        sections = [self.__render_interface(document)]
        sections.extend(self.__render_peers(document.get("peers") or ()))
        return "\n\n".join(sections)

    def run(self):
        sys.stdout.write(self.__render(yaml.safe_load(sys.stdin)) + "\n")


if __name__ == "__main__":
    TalosWge().run()
