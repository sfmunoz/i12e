#!/bin/bash
set -e -o pipefail
FLAG_FILE="/etc/i12e/artifact-tuned"
[ -f "$FLAG_FILE" ] && exit 0
MODE_FILE="/etc/i12e/mode"
[ -f "$MODE_FILE" ] || exit 0
MODE="$(cat /etc/i12e/mode)"
set -x
cat /etc/i12e/k3s/override-${MODE}.conf > /etc/systemd/system/k3s.service.d/override.conf
awk \
  -v IFACE="$1" \
  -v IP="$2" \
  '!/^(node-ip|flannel-iface):/ { print }
  END {
    if (IFACE != "" && IP != "") {
      printf("flannel-iface: \"%s\"\nnode-ip: \"%s\"\n",IFACE,IP)
    }
  }' \
  /etc/i12e/k3s/config-${MODE}.yaml > /etc/rancher/k3s/config.yaml
touch "$FLAG_FILE"
