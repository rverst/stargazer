# awesome stars

Total starred repositories: `{{ .Total }}`

## Contents

{{ range $key, $value := .Stars }}
* [{{- $key -}}](#{{- anchor $key -}}) ({{ len $value }})
{{ end }}

{{ range $key, $value := .Stars }}
## {{ $key }}
{{ range $value }}
- [{{- .Name -}}]({{- .Url -}}) - {{ .Description }} {{ with .License}} \[*{{ . }}*\]{{ end }} (⭐️{{ .Stars }}) {{ if .Archived }}*Archived!*{{ end }}
{{ end }}
{{ end }}
