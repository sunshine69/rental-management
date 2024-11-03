package main

import (
	"testing"

	ag "github.com/sunshine69/automation-go/lib"
	// u "github.com/sunshine69/golang-tools/utils"
)

func TestTemplate(t *testing.T) {
	tmpl := `var AllForms = map[string]any{
	{{ $g := .}}
	{{- range $typeName := $g }}
	"{{$typeName}}": model.{{$typeName}}{},
	{{- end}}
		}
var AllModelObjects []any = []any{ {{range $idx, $typeName := $g}}model.{{$typeName}}{}{{if ne $idx (add (len $g) -1 ) }}, {{end}}{{end}} }`
	textrpl := ag.GoTemplateString(tmpl, []string{`A`, `B`})
	println(textrpl)
}
