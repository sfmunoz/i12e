#!/bin/bash

set -e -o pipefail

DNAME="$(dirname "$0")"
SOPS="${DNAME}/sops.sh"
TALOS_MESH_PY="${DNAME}/talos-mesh.py"
[ "$I12E_ENV" = "" ] && I12E_ENV="dev"

[ "$CLUSTER_NAME" = "" ] && CLUSTER_NAME="cdev"
[ "$IP1" = "" ] && IP1="192.168.56.57"
[ "$IP2" = "" ] && IP2="192.168.56.58"
[ "$IP3" = "" ] && IP3="192.168.56.59"

IP_PUB=("----" "$IP1" "$IP2" "$IP3")
IP_PRIV=("----" "192.168.186.1" "192.168.186.2" "192.168.186.3")

export CLUSTER_NAME

export KUBECONFIG="./${CLUSTER_NAME}/kubeconfig"
export TALOSCONFIG="./${CLUSTER_NAME}/talosconfig"

function gen_config {
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
  set -x
  talosctl gen config $CLUSTER_NAME https://${IP_PUB[1]}:6443 \
    --with-secrets <(
      { set +x; } 2>/dev/null
      "${SOPS}" talos-secrets
    ) \
    --install-disk /dev/sda \
    --output - \
    --output-types "${OUTPUT_TYPES}" \
    --config-patch <(
      { set +x; } 2>/dev/null
      echo "---"
      cat "${DNAME}/../talos/${I12E_ENV}/common.yaml"
      case "$CFG_NAME" in
      node1 | node2 | node3)
        echo "---"
        sops decrypt "${CLUSTER_NAME}/wg${CFG_NAME#node}.yaml"
        ;;
      esac
      [ -f wg.yaml ] || exit 0
      echo "---"
      sops decrypt wg.yaml
    ) \
    --config-patch-control-plane <(
      { set +x; } 2>/dev/null
      echo "---"
      cat "${DNAME}/../talos/${I12E_ENV}/control-plane.yaml"
    ) \
    --config-patch-worker <(
      { set +x; } 2>/dev/null
      echo "---"
      cat "${DNAME}/../talos/${I12E_ENV}/worker.yaml"
    )
}

CMD="$1"

case "$CMD" in
mesh)
  set -x -e -o pipefail
  mkdir -p "${CLUSTER_NAME}"
  cd "${CLUSTER_NAME}"
  # "${TALOS_MESH_PY}" ${IP_PUB[1]}:51823 ${IP_PUB[2]}:51823 ${IP_PUB[3]}:51823
  # generate wg-quick files and host config to get into the mesh from the host
  "../${TALOS_MESH_PY}" -c ${IP_PUB[1]}:51823 ${IP_PUB[2]}:51823 ${IP_PUB[3]}:51823 192.168.56.51:51823
  cd ..
  ;;
talosconfig)
  set -x
  gen_config talosconfig >"${TALOSCONFIG}"
  talosctl config endpoint ${IP_PUB[1]}
  talosctl config node ${IP_PRIV[1]} ${IP_PRIV[2]} ${IP_PRIV[3]}
  ;;
debug-1 | debug-2 | debug-3)
  set -x
  N="${CMD#debug-}"
  NODE="node$N"
  gen_config $NODE
  ;;
install-1)
  set -x
  talosctl apply-config --nodes ${IP_PUB[1]} --file <(gen_config node1) --insecure
  while true; do
    talosctl bootstrap --nodes ${IP_PUB[1]} && break
    sleep 10
  done
  ;;
install-2 | install-3)
  set -x
  N="${CMD#install-}"
  NODE="node$N"
  talosctl apply-config --nodes ${IP_PUB[$N]} --file <(gen_config $NODE) --insecure
  ;;
update-1 | update-2 | update-3)
  set -x
  N="${CMD#update-}"
  NODE="node$N"
  talosctl apply-config --nodes ${IP_PRIV[$N]} --file <(gen_config $NODE)
  ;;
try-1 | try-2 | try-3)
  set -x
  N="${CMD#try-}"
  NODE="node$N"
  talosctl apply-config --nodes ${IP_PRIV[$N]} --file <(gen_config $NODE) --mode try
  ;;
kubeconfig)
  set -x
  talosctl kubeconfig --nodes ${IP_PUB[1]}
  ;;
source)
  cat <<__EOF
export KUBECONFIG="$KUBECONFIG"
export TALOSCONFIG="$TALOSCONFIG"
__EOF
  ;;
*)
  BNAME="$(basename "$0")"
  echo
  echo "Usage (order matters):"
  echo
  echo "  \$ ${BNAME} mesh                           -- mesh gen"
  echo "  \$ ${BNAME} talosconfig                    -- talosconfig gen"
  echo "  \$ ${BNAME} install-1                      -- control-plane node"
  echo "  \$ ${BNAME} kubeconfig                     -- kubeconfig gen"
  echo "  \$ ${BNAME} install-2/install-3            -- worker nodes"
  echo "  \$ ${BNAME} debug-1/debug-2/debug-3        -- debug config"
  echo "  \$ ${BNAME} try-1/try-2/try-3              -- try config"
  echo "  \$ ${BNAME} update-1/update-2/update-3     -- update config"
  echo "  \$ eval \$(${BNAME} source)                 -- set KUBECONFIG/TALOSCONFIG env vars"
  echo
  ;;
esac
