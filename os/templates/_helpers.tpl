{{- define "os.env" -}}
{{ regexSplit "-" .Release.Name -1 | first }}
{{- end }}

{{- define "os.profile" -}}
{{ regexSplit "-" .Release.Name -1 | rest | join "-" | default "main" }}
{{- end }}
