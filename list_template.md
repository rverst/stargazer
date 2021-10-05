{{- $wl := .WithLicense -}}
{{- $ws := .WithStars -}}
{{- $a := .Anchors -}}
{{- $s := .Stars -}}
# awesome stars

{{ .Credits.Text }}{{ .Credits.Link }}  
Total starred repositories: `{{ .Total }}`

{{- if .WithToc }}
## Contents
{{ range $key := .Keys }}
  - [{{ $key }}](#{{ with (index $a $key) }}{{ . }}{{ end }})
{{- end }}
{{- end }}


{{ range $key := .Keys }}
## {{ $key }}
{{ with (index $s $key) }}{{ range . }}
  - [{{- .NameWithOwner -}}]({{- .Url -}}) - {{ .Description }} 
{{- if $wl }}{{ with .License}} \[*{{ . }}*\]{{ end }}{{ end -}}
{{- if $ws }} (⭐️{{ .Stars }}){{ end -}}
{{- if .Archived }} *Archived!*{{ end -}}
{{- end }}
{{ end }}
{{ end }}
