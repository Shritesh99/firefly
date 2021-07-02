{{/*
Expand the name of the chart.
*/}}
{{- define "firefly.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Create a default fully qualified app name.
We truncate at 63 chars because some Kubernetes name fields are limited to this (by the DNS naming spec).
If release name contains chart name it will be used as a full name.
*/}}
{{- define "firefly.fullname" -}}
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
{{- define "firefly.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Common labels
*/}}
{{- define "firefly.labels" -}}
helm.sh/chart: {{ include "firefly.chart" . }}
{{ include "firefly.selectorLabels" . }}
{{- if .Chart.AppVersion }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
{{- end }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- end }}

{{/*
Selector labels
*/}}
{{- define "firefly.selectorLabels" -}}
app.kubernetes.io/name: {{ include "firefly.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end }}

{{/*
Create the name of the service account to use
*/}}
{{- define "firefly.serviceAccountName" -}}
{{- if .Values.serviceAccount.create }}
{{- default (include "firefly.fullname" .) .Values.serviceAccount.name }}
{{- else }}
{{- default "default" .Values.serviceAccount.name }}
{{- end }}
{{- end }}

{{- define "firefly.dataexchangeP2PHost" -}}
{{- if .Values.dataexchange.ingress.enabled }}
{{- index .Values.dataexchange.ingress.hosts 0 }}
{{- else }}
{{- printf "%s-dx:%d" (include "firefly.fullname" .) (.Values.dataexchange.service.p2pPort | int64) }}
{{- end }}
{{- end }}

{{- define "firefly.core.config" -}}
{{- if .Values.config.debugEnabled }}
log:
  level: debug
debug:
  port: 6060
{{- end }}
http:
  port: 5000
  address: 0.0.0.0
admin:
  port: 5001
  address: 0.0.0.0
  enabled: {{ .Values.config.adminEnabled }}
  preinit: {{ and .Values.config.adminEnabled .Values.config.preInit }}
ui:
  path: ./frontend
node:
  name: {{ .Values.config.organizationName }}-{{ include "firefly.fullname" . }}
org:
  name: {{ .Values.config.organizationName }}
  identity: {{ .Values.config.organizationIdentity }}
{{- if .Values.config.blockchain }}
blockchain:
  {{- toYaml (tpl .Values.config.blockchain .) | nindent 2 }}
{{- else if .Values.config.ethconnectUrl }}
blockchain:
  type: ethereum
  ethereum:
    ethconnect:
      url: {{ tpl .Values.config.ethconnectUrl . }}
      instance: {{ .Values.config.fireflyContractAddress }}
      topic: "0" # TODO should likely be configurable
      {{- if and .Values.config.ethconnectUsername .Values.config.ethconnectPassword }}
      auth:
        username: {{ .Values.config.ethconnectUsername }}
        password: {{ .Values.config.ethconnectPassword }}
      {{- end }}
      {{- if .Values.config.ethconnectPrefixShort }}
      prefixShort: {{ .Values.config.ethconnectPrefixShort }}
      {{- end }}
      {{- if .Values.config.ethconnectPrefixLong }}
      prefixLong: {{ .Values.config.ethconnectPrefixLong }}
      {{- end }}
{{- end }}
{{- if .Values.config.database }}
database:
  {{- toYaml (tpl .Values.config.database .) | nindent 2 }}
{{- else if .Values.config.postgresUrl }}
database:
  type: postgres
  postgres:
    url: {{ tpl .Values.config.postgresUrl . }}
    migrations:
      auto: {{ .Values.config.postgresAutomigrate }}
{{- end }}
{{- if .Values.config.publicstorage }}
publicstorage:
  {{- toYaml (tpl .Values.config.publicstorage .) | nindent 2 }}
{{- else if and .Values.config.ipfsApiUrl .Values.config.ipfsGatewayUrl }}
publicstorage:
  type: ipfs
  ipfs:
    api:
      url: {{ tpl .Values.config.ipfsApiUrl . }}
    gateway:
      url: {{ tpl .Values.config.ipfsGatewayUrl . }}
{{- end }}
{{- if and .Values.config.dataexchange (not .Values.dataexchange.enabled) }}
dataexchange:
  {{- toYaml (tpl .Values.config.dataexchange .) | nindent 2 }}
{{- else }}
dataexchange:
  {{- if .Values.dataexchange.enabled }}
  https:
    url: http://{{ include "firefly.fullname" . }}-dx:{{ .Values.dataexchange.service.apiPort }}
  {{- else }}
  https:
    url: {{ tpl .Values.config.dataexchangeUrl . }}
  {{- end }}
{{- end }}
{{- end }}