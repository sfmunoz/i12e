{{- $profile := include "os.profile" . -}}
{{- if eq $profile "dist-sh" -}}
{{- $env := include "os.env" . -}}
{{- $e := index .Values.env $env -}}
{{- $tgt_file := now | date "20060102-150405" | printf "os-%s.json" -}}
script: |
  #!/bin/bash
  IGN_FILE="$1"
  set -e -o pipefail
  for TARGET in {{ $e.targets | join " " }}
  do
    echo "==== $TARGET ===="
    set -x
    scp "$IGN_FILE" "core@${TARGET}:{{ $tgt_file }}"
    ssh "core@${TARGET}" "sudo flatcar-reset --keep-machine-id --keep-paths '/etc/ssh/ssh_host_.*' /var/log -F {{ $tgt_file }} && sudo systemctl reboot"
    { set +x; } 2> /dev/null
  done
{{- end }}
