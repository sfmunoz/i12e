#!/bin/bash
ENTRIES=(
  "/etc/i12e/namespaces"
  "/etc/i12e/traefik-config"
  "/etc/i12e/cert-manager"
  "/etc/i12e/csi-rclone"
  "/etc/i12e/reflector"
  "/opt/libexec/i12e/plugins/02-namespaces.sh"
  "/opt/libexec/i12e/plugins/05-traefik-config.sh"
  "/opt/libexec/i12e/plugins/08-rclone-conf.sh"
  "/opt/libexec/i12e/plugins/09-k3s-client-ca.sh"
  "/opt/libexec/i12e/plugins/10-csi-rclone.sh"
  "/opt/libexec/i12e/plugins/20-reflector.sh"
  "/opt/libexec/i12e/plugins/30-cert-manager.sh"
)
for E in "${ENTRIES[@]}"; do
  [ -e "$E" ] || continue
  set -x
  rm -rfv "$E"
  { set +x; } 2>/dev/null
done
