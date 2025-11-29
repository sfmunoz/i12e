{{- $cmd := include "os.cmd" . -}}
{{- if eq $cmd "build-sh" -}}
{{- $env := include "os.env" . -}}
{{- $e := index .Values.env $env -}}
{{- $vip := $e | dig "kube_vip" "vip" "" -}}
{{- $tgt_prefix := now | date "20060102-150405" | printf "os-%s" -}}
script: |
  #!/bin/bash
  set -e -o pipefail
  cd "$(dirname "$0")/.."
  POS=0
  echo "================ BUILD ================"
  for TARGET in {{ $e.targets | join " " }}
  do
    POS="$((POS+1))"
    OS_YAML="build/os-${TARGET}.yaml"
    OS_JSON="build/os-${TARGET}.json"
    set -x
    $CMD --set "cmd=flatcar-yaml,target=${TARGET},position=${POS}" > "$OS_YAML"
    ls -l "$OS_YAML"
    docker run --rm -i quay.io/coreos/butane:latest < "$OS_YAML" > "$OS_JSON"
    ls -l "$OS_YAML"
    { set +x; } 2> /dev/null
  done
  [ "$I12E_DIST" = "1" ] || exit 0
  echo "================ DIST ================"
  for TARGET in {{ $e.targets | join " " }}
  do
    OS_JSON_LOC="build/os-${TARGET}.json"
    OS_JSON_REM="{{ $tgt_prefix }}-${TARGET}.json"
    set -x
    scp "$OS_JSON_LOC" "core@${TARGET}:${OS_JSON_REM}"
    # systemd-run usage: avoid "ssh connection closed by remote host" error
    ssh "core@${TARGET}" "sudo flatcar-reset --keep-machine-id --keep-paths '/etc/ssh/ssh_host_.*' /var/log -F $OS_JSON_REM && sudo systemd-run bash -c 'sleep 1 ; systemctl reboot'"
    { set +x; } 2> /dev/null
  done
{{- end }}
