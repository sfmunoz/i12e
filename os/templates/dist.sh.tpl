{{- $profile := include "os.profile" . -}}
{{- if eq $profile "dist-sh" -}}
{{- $env := include "os.env" . -}}
script: |
  #!/bin/bash
  echo "env ....... {{ $env }}"
  echo "profile ... {{ $profile }}"
{{- end }}
