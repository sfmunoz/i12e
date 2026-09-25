#!/bin/bash

set -e -o pipefail

DNAME="$(realpath "$(dirname "$0")/..")"
SOPS="${DNAME}/scripts/sops.sh"

[ "$I12E_ENV" = "" ] && I12E_ENV="dev"
[ "$CLUSTER_NAME" = "" ] && CLUSTER_NAME="c${I12E_ENV}"
[ "$INSTALL_DISK" = "" ] && INSTALL_DISK="/dev/sda"
[ "$IP1" = "" ] && IP1="192.168.56.57"
[ "$IP2" = "" ] && IP2="192.168.56.58"
[ "$IP3" = "" ] && IP3="192.168.56.59"

IP_PUB=("----" "$IP1" "$IP2" "$IP3")
IP_PRIV=("----" "192.168.186.1" "192.168.186.2" "192.168.186.3")

function gen_config {
  set -e -o pipefail
  CFG_NAME="$1"
  case "$CFG_NAME" in
  node1)
    OUTPUT_TYPES="controlplane"
    ;;
  node2 | node3)
    OUTPUT_TYPES="worker"
    ;;
  talosconfig)
    OUTPUT_TYPES="talosconfig"
    ;;
  *)
    echo "error: unsupported '$1' argument"
    exit 1
    ;;
  esac
  { set +x; } 2>/dev/null
  SECRETS="$("${SOPS}" talos-secrets)"
  CONFIG_PATCH="$(
    { set +x; } 2>/dev/null
    set -e -o pipefail
    echo "---"
    cat "${DNAME}/talos/${I12E_ENV}/common.yaml"
    echo "---"
    "${SOPS}" talos-wge
    case "$CFG_NAME" in
    node1 | node2 | node3)
      echo "---"
      "${SOPS}" talos-wgm | "${DNAME}/scripts/talos-wgm.py" "${CFG_NAME#node}"
      ;;
    esac
  )"
  set -x
  talosctl gen config $CLUSTER_NAME https://${IP_PUB[1]}:6443 \
    --with-secrets <(
      { set +x; } 2>/dev/null
      echo "$SECRETS"
    ) \
    --install-disk "${INSTALL_DISK}" \
    --output - \
    --output-types "${OUTPUT_TYPES}" \
    --config-patch <(
      { set +x; } 2>/dev/null
      echo "$CONFIG_PATCH"
    ) \
    --config-patch-control-plane @"${DNAME}/talos/${I12E_ENV}/control-plane.yaml" \
    --config-patch-worker @"${DNAME}/talos/${I12E_ENV}/worker.yaml"
}

CMD="$1"

case "$CMD" in
talosconfig)
  set -x
  TFOLDER="${HOME}/.talos"
  mkdir -p "${TFOLDER}"
  chmod 700 "${TFOLDER}"
  gen_config talosconfig >"${TFOLDER}/config.${I12E_ENV}"
  ln -sf "config.${I12E_ENV}" "${TFOLDER}/config"
  talosctl config endpoint ${IP_PUB[1]}
  talosctl config node ${IP_PRIV[1]} ${IP_PRIV[2]} ${IP_PRIV[3]}
  chmod 600 "${TFOLDER}/config.${I12E_ENV}"
  ls -l "${TFOLDER}/config"
  ;;
debug-1 | debug-2 | debug-3)
  set -x
  N="${CMD#debug-}"
  NODE="node$N"
  gen_config $NODE
  ;;
install-1)
  NODE_CONFIG="$(gen_config node1)"
  set -x
  talosctl apply-config --nodes ${IP_PUB[1]} --file <(
    { set +x; } 2>/dev/null
    echo "$NODE_CONFIG"
  ) --insecure
  while true; do
    talosctl bootstrap --nodes ${IP_PUB[1]} && break
    sleep 10
  done
  ;;
install-2 | install-3 | update-1 | update-2 | update-3 | try-1 | try-2 | try-3)
  case "${CMD%%-*}" in
  install) FLAGS="--insecure" ;;
  try) FLAGS="--mode try" ;;
  *) FLAGS="" ;;
  esac
  N="${CMD#*-}"
  NODE="node$N"
  NODE_CONFIG="$(gen_config $NODE)"
  set -x
  talosctl apply-config --nodes ${IP_PUB[$N]} --file <(
    { set +x; } 2>/dev/null
    echo "$NODE_CONFIG"
  ) $FLAGS
  ;;
kubeconfig)
  set -x
  KFOLDER="${HOME}/.kube"
  mkdir -p "${KFOLDER}"
  chmod 700 "${KFOLDER}"
  talosctl kubeconfig - --nodes ${IP_PUB[1]} >"${KFOLDER}/config.${I12E_ENV}"
  ln -sf "config.${I12E_ENV}" "${KFOLDER}/config"
  chmod 600 "${KFOLDER}/config.${I12E_ENV}"
  ls -l "${KFOLDER}/config"
  ;;
*)
  BNAME="$(basename "$0")"
  echo
  echo "Usage (order matters):"
  echo
  echo "  \$ ${BNAME} talosconfig                    -- talosconfig gen"
  echo "  \$ ${BNAME} install-1                      -- control-plane node"
  echo "  \$ ${BNAME} kubeconfig                     -- kubeconfig gen"
  echo "  \$ ${BNAME} install-2/install-3            -- worker nodes"
  echo "  \$ ${BNAME} debug-1/debug-2/debug-3        -- debug config"
  echo "  \$ ${BNAME} try-1/try-2/try-3              -- try config"
  echo "  \$ ${BNAME} update-1/update-2/update-3     -- update config"
  echo
  ;;
esac
