#!/bin/bash
set -e -o pipefail
[ -f "/etc/systemd/system/k3s.service" ] && exit 0
echo "Installing k3s..."
MESH_CFG_FILE="/etc/i12e/k3s/mesh.cfg"
if [ ! -f "$MESH_CFG_FILE" ]; then
  echo "error: '${MESH_CFG_FILE}' file doesn't exit"
  set -x
  exit 1
fi
source "$MESH_CFG_FILE"
if [ "$IFACE" = "" ]; then
  echo "error: undefined 'IFACE' in '${MESH_CFG_FILE}' file"
  set -x
  rm -f "$MESH_CFG_FILE"
  exit 1
fi
if [ "$IP" = "" ]; then
  echo "error: undefined 'IP' in '${MESH_CFG_FILE}' file"
  set -x
  rm -f "$MESH_CFG_FILE"
  exit 1
fi
MODE_FILE="/etc/i12e/mode"
if [ ! -f "$MODE_FILE" ]; then
  echo "error: '${MODE_FILE}' file doesn't exit"
  set -x
  exit 1
fi
MODE="$(cat "$MODE_FILE")"
set -x
mkdir -p /etc/rancher/k3s /etc/systemd/system/k3s.service.d
chmod 755 /etc/rancher /etc/rancher/k3s /etc/systemd/system/k3s.service.d
cat /etc/i12e/k3s/override-${MODE}.conf >/etc/systemd/system/k3s.service.d/override.conf
chmod 644 /etc/systemd/system/k3s.service.d/override.conf
awk \
  -v IFACE="$IFACE" \
  -v IP="$IP" \
  '!/^(node-ip|flannel-iface):/ { print }
  END {
    if (IFACE != "" && IP != "") {
      printf("flannel-iface: \"%s\"\nnode-ip: \"%s\"\n",IFACE,IP)
    }
  }' \
  /etc/i12e/k3s/config-${MODE}.yaml >/etc/rancher/k3s/config.yaml
chmod 600 /etc/rancher/k3s/config.yaml
export INSTALL_K3S_VERSION="v1.36.3+k3s1"
export INSTALL_K3S_SKIP_SELINUX_RPM="true"
curl -sfL https://get.k3s.io | sh -s
{ set +x; } 2>/dev/null
