package source

import (
	"fmt"
)

func GetDescription() string {
	return fmt.Sprintf("%s%s", newMsg, `

My name is {{.name | stringsTitle}} and I am {{.age}} years old. I like:
{{- define "interestsList"}}- {{. | stringsToUpper}}{{end}}{{range $i := .interests}}
    {{template "interestsList" $i}}
{{- end}}
This is my favorite quote: "{{.quote | stringsToLower}}"
`)
}
