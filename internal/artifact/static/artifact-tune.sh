#!/bin/bash
set -e -o pipefail
function etc_wireguard {
  for i in $(seq 0 2); do
    f="/etc/wireguard/wg${i}.conf"
    s="wg-quick@wg${i}.service"
    if [ -f "$f" ]; then
      set -x
      systemctl enable $s
      systemctl start $s
      { set +x; } 2>/dev/null
    else
      set -x
      # stop won't work: '/etc/wireguard/wgN.conf' since doesn't exist
      # link deletion is best effort: 'PostUp', 'PreDown', ... won't be executed
      #systemctl stop $s
      ip link delete wg$i || true
      systemctl disable $s
      { set +x; } 2>/dev/null
    fi
  done
}
function artifact_tune {
  FLAG_FILE="/etc/i12e/flags/artifact-tuned"
  [ -f "$FLAG_FILE" ] && return 0
  set -x
  systemctl daemon-reload
  systemctl restart update-engine
  systemctl restart nftables
  systemctl enable nftables
  { set +x; } 2>/dev/null
  etc_wireguard
  set -x
  #systemctl restart locksmithd
  touch "$FLAG_FILE"
  { set +x; } 2>/dev/null
}
artifact_tune
