{{- $l := .WithLicense -}}
{{- $s := .WithStars -}}
# awesome stars

{{ .Credits.Text }}{{ .Credits.Link }}  
Total starred repositories: `{{ .Total }}`

{{ if .WithToc }}
## Contents
    {{ range $key, $value := .Stars }}
* [{{ $key }}](#{{ anchor $key }}) ({{ len $value }})
    {{- end }}
{{- end }}

{{ range $key, $value := .Stars }}
## {{ $key }}
    {{ range $value }}
- [{{- .NameWithOwner -}}]({{- .Url -}}) - {{ .Description }} 
{{ if $l }}{{ with .License}}\[*{{ . }}*\]{{ end }}{{ end }}
{{ if $s }} (⭐️{{ .Stars }}) {{ end }}
{{ if .Archived }}*Archived!*{{ end }}
    {{- end }}
{{ end }}
