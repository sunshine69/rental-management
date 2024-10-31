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

func CodeGen(templateFile string) {
	if sqlb, err := os.ReadFile("db/schema.sql"); err == nil {
		sqls := ag.SplitTextByPattern(string(sqlb), `(?m)^CREATE TABLE IF NOT EXISTS .*`, true)
		// fmt.Printf("%s\n", u.JsonDump(sqls, "  "))
		for _, sqltext := range sqls {
			fmt.Println(sqltext)
			GenerateClass(sqltext, templateFile)
		}
	} else {
		panic(err.Error())
	}
}

func GenerateClass(sqltext, classTemplateFile string) {
	typenamePtn := regexp.MustCompile(`CREATE TABLE .* ["]?([^\s"]+)["]? \(`)
	o := typenamePtn.FindStringSubmatch(sqltext)

	fieldPtn := regexp.MustCompile(`"([^\s"]+)" ([^\s\(]+)`)
	o1 := fieldPtn.FindAllStringSubmatch(sqltext, -1)
	fieldsList := []string{}
	fieldmap := map[string]interface{}{}
	for _, v := range o1 {
		fieldmap[v[1]] = GotypeLookup(v[2])
		fieldsList = append(fieldsList, v[1])
	}

	uniqueFieldPtn := regexp.MustCompile(`UNIQUE[ ]*\(([^\)]+)\)`)
	o2 := uniqueFieldPtn.FindStringSubmatch(sqltext)
	uniqueFields := strings.Split(o2[1], ",")
	uniqueFields = ag.SliceMap(uniqueFields, func(s string) *string { o := strings.TrimSpace(strings.ReplaceAll(s, `"`, ``)); return &o })
	uniqueFieldsMap := ag.SliceToMap(uniqueFields)
	for k := range uniqueFieldsMap {
		uniqueFieldsMap[k] = fieldmap[k]
	}
	query_new := "SELECT * FROM " + o[1] + " WHERE "
	for idx, _f := range uniqueFields {
		query_new = query_new + " " + _f + " = ?"
		if idx < len(uniqueFields)-1 {
			query_new = query_new + " AND "
		}
	}
	targetFile := filepath.Dir(classTemplateFile) + "/" + o[1] + ".go"
	ag.GoTemplateFile(classTemplateFile, targetFile, map[string]interface{}{
		"typename":        o[1],
		"fields":          fieldmap,
		"fieldsList":      fieldsList,
		"uniqueFieldsMap": uniqueFieldsMap,
		"uniqueFields":    uniqueFields,
		"query_new":       query_new,
	}, 0640)
	u.RunSystemCommandV2("go fmt "+targetFile, true)
}

func main() {
	gentype := flag.String("type", "", "Code gen type. Can be model | api")
	flag.Parse()

	switch *gentype {
	case "model":
		CodeGen("model/class-template.go.tmpl")
	case "api":
		CodeGen("api/api-template.go.tmpl")
	default:
		fmt.Println("Unknown type " + *gentype)
	}
}
