{{- define "os.env" -}}
{{ .Release.Name }}
{{- end }}

{{- define "os.cmd" -}}
{{ .Values.cmd | default "build-sh" }}
{{- end }}

{{- define "os.k3s_cmd" -}}
{{- $position := required "position missing" .position -}}
{{ le $position 3 | ternary "server" "agent" }}
{{- end }}
