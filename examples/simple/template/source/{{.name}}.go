package source

import (
	"fmt"
)

func NewUserBio() string {
	return fmt.Sprintf(
		`%[1]s%[1]s
%[2]s
%[3]s
%[1]s%[1]s`,
		separator,
		newUserMsg,
		`
My name is {{ .name | stringsTitle }} and I am {{ .age }} years old. 
I like:
{{- define "interestsList" }}=> {{ . | stringsToUpper }}{{ end -}}
{{ range $i := .interests}}
  {{template "interestsList" $i}}
{{- end }}

{{/* this is a comment */}}
This is my favorite quote: "{{.quote | stringsToLower}}"
`,
	)
}
