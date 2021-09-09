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
|     |     {{ if $l }}|     {{ end }}{{ if $s }}|     {{ end }}|
|-----|-----{{ if $l }}|:---:{{ end }}{{ if $s }}|----:{{ end }}|
| [{{- .NameWithOwner -}}]({{- .Url -}}) | {{ .Description }} {{ if .Archived }}(*archived*){{ end }} {{ if $l }}| {{ with .License}}{{ . }}{{ else }}-{{ end }}{{ end }} {{ if $s }}| ⭐️{{ .Stars }}{{ end }} |
    {{- end }}
{{- end }}
