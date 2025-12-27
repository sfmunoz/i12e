{{/*
Expand the name of the chart.
*/}}
{{- define "genesis.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Create a default fully qualified app name.
We truncate at 63 chars because some Kubernetes name fields are limited to this (by the DNS naming spec).
If release name contains chart name it will be used as a full name.
*/}}
{{- define "genesis.fullname" -}}
{{- if .Values.fullnameOverride }}
{{- .Values.fullnameOverride | trunc 63 | trimSuffix "-" }}
{{- else }}
{{- $name := default .Chart.Name .Values.nameOverride }}
{{- if contains $name .Release.Name }}
{{- .Release.Name | trunc 63 | trimSuffix "-" }}
{{- else }}
{{- printf "%s-%s" .Release.Name $name | trunc 63 | trimSuffix "-" }}
{{- end }}
{{- end }}
{{- end }}

{{/*
Create chart name and version as used by the chart label.
*/}}
{{- define "genesis.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Common labels
*/}}
{{- define "genesis.labels" -}}
helm.sh/chart: {{ include "genesis.chart" . }}
{{ include "genesis.selectorLabels" . }}
{{- if .Chart.AppVersion }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
{{- end }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- end }}

{{/*
Selector labels
*/}}
{{- define "genesis.selectorLabels" -}}
app.kubernetes.io/name: {{ include "genesis.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end }}

{{/*
Create the name of the service account to use
*/}}
{{- define "genesis.serviceAccountName" -}}
{{- if .Values.serviceAccount.create }}
{{- default (include "genesis.fullname" .) .Values.serviceAccount.name }}
{{- else }}
{{- default "default" .Values.serviceAccount.name }}
{{- end }}
{{- end }}

{{- define "genesis.code_sha256sum" -}}
{{-
cat
  ( .Files.Get "app/__init__.py" )
  ( .Files.Get "app/__main__.py" )
  ( .Files.Get "app/butane.py" )
  ( .Files.Get "app/entrypoint.py" )
  ( .Files.Get "app/install.py" )
  ( .Files.Get "app/templates/crictl.yaml" )
  ( .Files.Get "app/templates/flatcar-update.conf" )
  ( .Files.Get "app/templates/flatcar.yaml" )
  ( .Files.Get "app/templates/k3s-config.yaml" )
  ( .Files.Get "app/templates/k3s-override.conf" )
  ( .Files.Get "app/templates/systemd-genesis.conf" )
  | sha256sum
-}}
{{- end -}}

{{- define "genesis.code_blob" -}}
{{-
list
  ( "app/__init__.py" | b64enc )
  ( .Files.Get "app/__init__.py" | b64enc )
  ( "app/__main__.py" | b64enc )
  ( .Files.Get "app/__main__.py" | b64enc )
  ( "app/butane.py" | b64enc )
  ( .Files.Get "app/butane.py" | b64enc )
  ( "app/entrypoint.py" | b64enc )
  ( .Files.Get "app/entrypoint.py" | b64enc )
  ( "app/install.py" | b64enc )
  ( .Files.Get "app/install.py" | b64enc )
  ( "app/templates/crictl.yaml" | b64enc )
  ( .Files.Get "app/templates/crictl.yaml" | b64enc )
  ( "app/templates/flatcar-update.conf" | b64enc )
  ( .Files.Get "app/templates/flatcar-update.conf" | b64enc )
  ( "app/templates/flatcar.yaml" | b64enc )
  ( .Files.Get "app/templates/flatcar.yaml" | b64enc )
  ( "app/templates/k3s-config.yaml" | b64enc )
  ( .Files.Get "app/templates/k3s-config.yaml" | b64enc )
  ( "app/templates/k3s-override.conf" | b64enc )
  ( .Files.Get "app/templates/k3s-override.conf" | b64enc )
  ( "app/templates/systemd-genesis.conf" | b64enc )
  ( .Files.Get "app/templates/systemd-genesis.conf" | b64enc )
| join "\n"
-}}
{{- end -}}
