{{- define "os.env" -}}
{{ .Release.Name }}
{{- end }}

{{- define "os.cmd" -}}
{{ .Values.cmd | default "build-sh" }}
{{- end }}
