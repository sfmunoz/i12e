#!/bin/bash
set -e -o pipefail
function artifact_tune {
  FLAG_FILE="/etc/i12e/flags/artifact-tuned"
  [ -f "$FLAG_FILE" ] && return 0
  MODE_FILE="/etc/i12e/mode"
  [ -f "$MODE_FILE" ] || return 0
  MODE="$(cat "$MODE_FILE")"
  set -x
  cat /etc/i12e/k3s/override-${MODE}.conf >/etc/systemd/system/k3s.service.d/override.conf
  awk \
    -v IFACE="$1" \
    -v IP="$2" \
    '!/^(node-ip|flannel-iface):/ { print }
    END {
      if (IFACE != "" && IP != "") {
        printf("flannel-iface: \"%s\"\nnode-ip: \"%s\"\n",IFACE,IP)
      }
    }' \
    /etc/i12e/k3s/config-${MODE}.yaml >/etc/rancher/k3s/config.yaml
  systemctl daemon-reload
  systemctl restart update-engine
  systemctl restart nftables
  systemctl enable nftables
  #systemctl restart locksmithd
  touch "$FLAG_FILE"
  { set +x; } 2>/dev/null
}
function k3s_install {
  [ -f "/etc/systemd/system/k3s.service" ] && return 0
  echo "Installing k3s..."
  set -x
  export INSTALL_K3S_VERSION="v1.34.3+k3s1"
  curl -sfL https://get.k3s.io | sh -s
  { set +x; } 2>/dev/null
}
artifact_tune "$1" "$2"
k3s_install
