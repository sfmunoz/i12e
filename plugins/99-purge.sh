#!/bin/bash
ENTRIES=(
  "/opt/libexec/i12e/plugins/02-namespaces.sh"
  "/etc/i12e/namespaces"
)
for E in "${ENTRIES[@]}"; do
  [ -e "$E" ] || continue
  set -x
  rm -rfv "$E"
  { set +x; } 2>/dev/null
done
