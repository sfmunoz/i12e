#!/bin/bash
[ "$NODE" = "" ] && NODE="192.168.186.1"
[ "$FULL" = "true" ] || FULL="false"

set -e -o pipefail

X1=">"
X2=" "

[ "$FULL" = "true" ] && X1=" " && X2=">"

cat <<__EOF

$X1 FULL=false: only wipes EPHEMERAL partition (enough most of the times)
$X2 FULL=true:  wipes both EPHEMERAL and STATE partitions

__EOF

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
