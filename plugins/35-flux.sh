#!/bin/bash
export FLUX_VERSION="2.9.4" # exported: it's used by https://fluxcd.io/install.sh
SHA256SUM="0c91d4bbbc2aa9c84b42608184534437ef92e5f2d6b862e99a11c0bb24ad0941"
FLUX_BIN="/opt/bin/flux"
FLUX_CFG="/etc/i12e/flux/flux.cfg"
SOPS_AGE_YAML="/etc/i12e/flux/sops-age.yaml"
NS="flux-system"
[ "$KUBECONFIG" = "" ] && KUBECONFIG="/etc/rancher/k3s/k3s.yaml"
export KUBECONFIG
set -x -e -o pipefail
function flux_bin_install {
  [ -x "${FLUX_BIN}" ] && [ "$(sha256sum "${FLUX_BIN}")" = "${SHA256SUM}  ${FLUX_BIN}" ] && return 0
  rm -rf "${FLUX_BIN}"
  curl -s https://fluxcd.io/install.sh | bash -s - /opt/bin
  [ "$(sha256sum "${FLUX_BIN}")" = "${SHA256SUM}  ${FLUX_BIN}" ] && return 0
  rm -rf "${FLUX_BIN}"
  echo "error: '${FLUX_BIN}' download failed" >&2
  return 1
}
function flux_cfg_read {
  [ -f "$FLUX_CFG" ] || return 1
  source "$FLUX_CFG" 2>/dev/null #  GITHUB_TOKEN and AGE_KEY are not shown
  [ "${CLUSTER}" = "dev" -o "${CLUSTER}" = "prod" ] || return 1
  [ ${#GITHUB_TOKEN} -gt 0 ] || return 1 # length: GITHUB_TOKEN is not shown
  [ ${#AGE_KEY} -gt 0 ] || return 1      # length: AGE_KEY is not shown
  export CLUSTER GITHUB_TOKEN AGE_KEY
  return 0
}
function flux_bootstrap {
  flux bootstrap github \
    --token-auth \
    --owner=sfmunoz \
    --repository=i12e \
    --path=clusters/${CLUSTER} \
    --branch=main \
    --private=false \
    --personal=true \
    --author-name "flux-${CLUSTER}-bot" \
    --author-email "46285520+sfmunoz@users.noreply.github.com" \
    --components-extra=source-watcher
  return $?
}
function flux_create_namespace {
  k3s kubectl get ns "${NS}" >/dev/null 2>&1 && return 0
  { set +x; } 2>/dev/null
  echo "creating '${NS}' namespace (it doesn't exist)"
  set -x
  k3s kubectl create ns "${NS}" || return $?
  k3s kubectl get ns "${NS}" >/dev/null 2>&1 || return $?
  return 0
}
function flux_create_secret_sops_age {
  [ -f "$SOPS_AGE_YAML" ] || return 1
  k3s kubectl apply --server-side=true -f "${SOPS_AGE_YAML}"
  return $?
}
flux_bin_install || exit $?
flux check && exit 0
flux check --pre || exit $?
flux_create_namespace || exit $?
flux_create_secret_sops_age || exit $?
flux_cfg_read || exit $?
flux_bootstrap || exit $?
exit 0
