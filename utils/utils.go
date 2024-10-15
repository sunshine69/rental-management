package utils

import (
	"path/filepath"
	"reflect"
	"regexp"
	"strings"

	ag "github.com/sunshine69/automation-go/lib"
	u "github.com/sunshine69/golang-tools/utils"
)

func Must[T any](o T, err error) T {
	u.CheckErr(err, "Must")
	return o
}

func GotypeLookup(sqltype string) string {
	switch {
	case strings.Contains(sqltype, "text") || strings.Contains(sqltype, "varchar"):
		return "string"
	case strings.Contains(sqltype, "int") || strings.Contains(sqltype, "integer"):
		return "int64"
	default:
		return ""
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
		if v[1] != "id" {
			fieldsList = append(fieldsList, v[1])
		}
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

func ReflectStruct(astruct any) map[string]map[string]interface{} {

	o := map[string]map[string]interface{}{}
	tagExtractPtn := regexp.MustCompile(`db:"([^"]+)"`)

	rf := reflect.TypeOf(astruct)
	if rf.Kind().String() != "struct" {
		panic("I need a struct")
	}
	rValue := reflect.ValueOf(astruct)
	for i := 0; i < rf.NumField(); i++ {
		f := rf.Field(i)
		o[f.Name] = map[string]interface{}{}
		dbcolum := tagExtractPtn.FindStringSubmatch(string(f.Tag))[1]
		fieldValue := rValue.Field(i)
		switch fieldValue.Type().String() {
		case "string":
			o[f.Name][dbcolum] = fieldValue.String()
		case "int64":
			o[f.Name][dbcolum] = fieldValue.Int()
		default:
			panic("Unsupported field type " + fieldValue.Type().String())
		}
	}
	return o
}
