package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	ag "github.com/sunshine69/automation-go/lib"
	u "github.com/sunshine69/golang-tools/utils"
)

func GotypeLookup(sqltype string) string {
	switch {
	case strings.Contains(sqltype, "text") || strings.Contains(sqltype, "varchar") || strings.Contains(sqltype, "jsonb") || strings.Contains(sqltype, "json"):
		return "string"
	case strings.Contains(sqltype, "int") || strings.Contains(sqltype, "integer"):
		return "int64"
	default:
		return ""
	}
}

var sqlCommentPtn *regexp.Regexp = regexp.MustCompile(`^[\s]*\-\-`)

var ObjectList []string

func CodeGen(templateFile string) {
	if sqlb, err := os.ReadFile("db/schema.sql"); err == nil {
		sqls := ag.SplitTextByPattern(string(sqlb), `(?m)^CREATE TABLE IF NOT EXISTS .*`, true)
		// fmt.Printf("%s\n", u.JsonDump(sqls, "  "))
		for _, sqltext := range sqls {
			GenerateClass(sqltext, templateFile)
		}
	} else {
		panic(err.Error())
	}
}

func GenerateClass(sqltext, classTemplateFile string) {
	fmt.Println("DEBUG " + sqltext)
	typenamePtn := regexp.MustCompile(`CREATE TABLE .* ["]?([^\s"]+)["]? \(`)
	var typeName string
	if typeNamematch := typenamePtn.FindStringSubmatch(sqltext); typeNamematch != nil {
		typeName = typeNamematch[1]
		ObjectList = append(ObjectList, typeName)
	} else {
		panic("[ERROR] can not parse typeName. Check your sql CREATE TABLE - need to match with my pattern here\n")
	}

	fieldPtn := regexp.MustCompile(`"([^\s"]+)" ([^\s\(]+)`)
	o1 := fieldPtn.FindAllStringSubmatch(sqltext, -1)
	fieldsList := []string{}
	//fieldmap[field name] => golang type as string
	fieldmap := map[string]interface{}{}
	for _, v := range o1 {
		fieldmap[v[1]] = GotypeLookup(v[2])
		fieldsList = append(fieldsList, v[1])
	}

	uniqueFieldPtn := regexp.MustCompile(`UNIQUE[ ]*\(([^\)]+)\)`)
	o2 := uniqueFieldPtn.FindStringSubmatch(sqltext)
	uniqueFields := strings.Split(o2[1], ",")
	uniqueFields = ag.SliceMap(uniqueFields, func(s string) *string { o := strings.TrimSpace(strings.ReplaceAll(s, `"`, ``)); return &o })
	//uniqueFieldsMap[field name] => golang type as string
	uniqueFieldsMap := ag.SliceToMap(uniqueFields)
	for k := range uniqueFieldsMap {
		uniqueFieldsMap[k] = fieldmap[k]
	}
	// uniqueStringFields used to generate search function thus we only take string
	uniqueStringFields := ag.SliceMap(uniqueFields, func(s string) *string {
		if uniqueFieldsMap[s].(string) == "string" {
			return &s
		} else {
			return nil
		}
	})
	query_new := "SELECT * FROM " + typeName + " WHERE "
	for idx, _f := range uniqueFields {
		query_new = query_new + " " + _f + " = ?"
		if idx < len(uniqueFields)-1 {
			query_new = query_new + " AND "
		}
	}
	targetFile := filepath.Dir(classTemplateFile) + "/" + typeName + ".go"
	ag.GoTemplateFile(classTemplateFile, targetFile, map[string]interface{}{
		"typename":           typeName,
		"fields":             fieldmap,
		"fieldsList":         fieldsList,
		"uniqueFieldsMap":    uniqueFieldsMap,
		"uniqueFields":       uniqueFields,
		"uniqueStringFields": uniqueStringFields,
		"query_new":          query_new,
	}, 0640)
	u.RunSystemCommandV2("go fmt "+targetFile, true)
}

func main() {
	gentype := flag.String("type", "", "Code gen type. Can be model | api")
	flag.Parse()

	switch *gentype {
	case "model":
		CodeGen("model/class-template.go.tmpl")
		tmpl := `var AllForms = map[string]any{
			{{ $g := .}}
			{{- range $typeName := $g }}
			"{{$typeName|title}}": {{$typeName|title}}{},
			{{- end}}
				}
		var AllModelObjects []any = []any{ {{range $idx, $typeName := $g}}{{$typeName|title}}{}{{if ne $idx (add (len $g) -1 ) }}, {{end}}{{end}} }`
		textrpl := ag.GoTemplateString(tmpl, ObjectList)
		ag.BlockInFile("model/global.go", []string{}, []string{`\/\/ End generate AllModelObjects`}, []string{`\/\/ Auto generate AllModelObjects`}, textrpl, true, false)
		u.RunSystemCommandV2("go fmt model/global.go", true)
	case "api":
		CodeGen("api/api-template.go.tmpl")
	default:
		fmt.Println("Unknown type " + *gentype)
	}
}
