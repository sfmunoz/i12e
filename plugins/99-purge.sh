#!/bin/bash
ENTRIES=(
  "/etc/i12e/namespaces"
  "/etc/i12e/traefik-config/"
  "/etc/i12e/csi-rclone"
  "/opt/libexec/i12e/plugins/02-namespaces.sh"
  "/opt/libexec/i12e/plugins/05-traefik-config.sh"
  "/opt/libexec/i12e/plugins/08-rclone-conf.sh"
  "/opt/libexec/i12e/plugins/10-csi-rclone.sh"
)
for E in "${ENTRIES[@]}"; do
  [ -e "$E" ] || continue
  set -x
  rm -rfv "$E"
  { set +x; } 2>/dev/null
done
