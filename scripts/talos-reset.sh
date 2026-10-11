#!/bin/bash
[ "$NODE" = "" ] && NODE="192.168.186.1"
[ "$FULL" = "true" ] || FULL="false"

set -e -o pipefail

OK="$((RANDOM + 10000))"

read -p "dangerous! type '${OK}' to reset NODE=${NODE} (FULL=${FULL}): "

if [ "$REPLY" != "$OK" ]; then
  echo "aborted!!"
  exit 1
fi

SYSTEMS_LABELS_TO_WIPE="--system-labels-to-wipe EPHEMERAL"

[ "$FULL" = "true" ] && SYSTEMS_LABELS_TO_WIPE="${SYSTEMS_LABELS_TO_WIPE} --system-labels-to-wipe STATE"

set -x
talosctl -n "$NODE" reset $SYSTEMS_LABELS_TO_WIPE --graceful=false --reboot
{ set +x; } 2>/dev/null

cat <<__EOF

Reset done. Next steps:

  \$ talosctl dashboard
  \$ talosctl bootstrap --nodes "$NODE"

__EOF
