#!/bin/bash
set -e -o pipefail
FLAG_FILE="/etc/i12e/config-patched"
[ -f "$FLAG_FILE" ] && exit 0
MODE_FILE="/etc/i12e/mode"
[ -f "$MODE_FILE" ] || exit 0
MODE="$(cat /etc/i12e/mode)"
set -x
cat /etc/i12e/k3s/config-${MODE}.yaml > /etc/rancher/k3s/config.yaml
cat /etc/i12e/k3s/override-${MODE}.conf > /etc/systemd/system/k3s.service.d/override.conf
IFACE="$1"
IP="$2"
if [ "$IFACE" != "" -a "$IP" != "" ]
then
  awk \
    -v IFACE="$IFACE" \
    -v IP="$IP" \
    '!/^(node-ip|flannel-iface):/ {
      print
    }
    END {
      printf("flannel-iface: \"%s\"\nnode-ip: \"%s\"\n",IFACE,IP)
    }' \
    /etc/rancher/k3s/config.yaml > /etc/rancher/k3s/config.yaml.new
  cat /etc/rancher/k3s/config.yaml.new > /etc/rancher/k3s/config.yaml
  rm -f /etc/rancher/k3s/config.yaml.new
fi
touch "$FLAG_FILE"
