#!/bin/bash
ENTRIES=(
  "/etc/i12e/namespaces"
  "/etc/i12e/traefik-config/"
  "/opt/libexec/i12e/plugins/02-namespaces.sh"
  "/opt/libexec/i12e/plugins/05-traefik-config.sh"
  "/opt/libexec/i12e/plugins/08-rclone-conf.sh"
)
for E in "${ENTRIES[@]}"; do
  [ -e "$E" ] || continue
  set -x
  rm -rfv "$E"
  { set +x; } 2>/dev/null
done
