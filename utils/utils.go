package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"strings"

	ag "github.com/sunshine69/automation-go/lib"
	u "github.com/sunshine69/golang-tools/utils"
	"github.com/sunshine69/rental-management/model"
)

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

func AssertInt64ValueForMap(input map[string]interface{}) map[string]interface{} {
	for k, v := range input {
		if v, ok := v.(float64); ok {
			input[k] = int64(v)

		}
	}
	return input
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

type StructInfo struct {
	Name       string
	FieldName  []string
	FieldType  map[string]string
	FieldValue map[string]any
	TagCapture map[string][][]string
}

// Give it a struct and a tag pattern to capture the tag content - return a map of string which is the struct Field name, point to a map of
// string which is the capture in the pattern
func ReflectStruct(astruct any, tagPtn string) StructInfo {
	if tagPtn == "" {
		tagPtn = `db:"([^"]+)"`
	}
	o := StructInfo{}
	tagExtractPtn := regexp.MustCompile(tagPtn)

	rf := reflect.TypeOf(astruct)
	o.Name = rf.Name()
	if rf.Kind().String() != "struct" {
		panic("I need a struct")
	}
	rValue := reflect.ValueOf(astruct)
	o.FieldName = []string{}
	o.FieldType = map[string]string{}
	o.FieldValue = map[string]any{}
	o.TagCapture = map[string][][]string{}
	for i := 0; i < rf.NumField(); i++ {
		f := rf.Field(i)
		o.FieldName = append(o.FieldName, f.Name)
		fieldValue := rValue.Field(i)
		o.FieldType[f.Name] = fieldValue.Type().String()
		o.TagCapture[f.Name] = [][]string{}
		switch fieldValue.Type().String() {
		case "string":
			o.FieldValue[f.Name] = fieldValue.String()
		case "int64":
			o.FieldValue[f.Name] = fieldValue.Int()
		default:
			fmt.Printf("Unsupported field type " + fieldValue.Type().String())
		}
		if ext := tagExtractPtn.FindAllStringSubmatch(string(f.Tag), -1); ext != nil {
			o.TagCapture[f.Name] = append(o.TagCapture[f.Name], ext...)
		}
	}
	return o
}

var AllModelObjects []any = []any{model.Tenant{}, model.Property{}, model.Account{}, model.Contract{}, model.Payment{}, model.Maintenance_request{}, model.Property_manager{}, model.Invoice{}}

// Take all structs in model and generate golang template html form - write to target dir
func FormGen(structType any, writeDirectory string) {
	sInfo := ReflectStruct(structType, `([^\:\s]+):"([^"]+)"`)
	destFile := writeDirectory + "/" + sInfo.Name + ".html"

	fieldProp := map[string]map[string]any{}
	for _, v := range sInfo.FieldName {
		fieldProp[v] = map[string]any{"display": true, "ele": "<input", "close_ele": "/>", "type": "text", "label": ""}
		tags := sInfo.TagCapture[v]
		for _, tagset := range tags {
			if (tagset[1] == "form" && len(tagset) >= 3 && tagset[2] == "-") || v == "Id" {
				fieldProp[v]["display"] = false
			}
			if tagset[1] == "db" && len(tagset) >= 3 && strings.Contains(tagset[2], "unique") {
				fieldProp[v]["label"] = "*"
			}
			if tagset[1] == "form" && len(tagset) >= 3 && strings.Contains(tagset[2], "ele=textarea") {
				fieldProp[v]["ele"] = "<textarea"
				fieldProp[v]["close_ele"] = "></textarea>"
			}
		}
	}
	// fmt.Printf("%v\n", u.JsonDump(fieldProp, ""))
	data := map[string]any{
		"formName":   sInfo.Name,
		"formClass":  "form-group",
		"formAction": "/" + strings.ToLower(sInfo.Name),
		"formID":     strings.ToLower(sInfo.Name),
		"fInfo":      sInfo,
		"fieldProp":  fieldProp,
	}
	ag.GoTemplateFile("utils/form.go.tmpl", destFile, data, 0640)
	formNameList := []string{}
	for _, f := range AllModelObjects {
		sInfo := ReflectStruct(f, `form:"([^"]+)"`)
		formNameList = append(formNameList, sInfo.Name)
	}
	ag.GoTemplateFile("utils/form-header.go.tmpl", writeDirectory+"/form-header.html", map[string]any{"formNameList": formNameList}, 0640)
	// Generate some common func handler and form validation to copy/paste into the app.go

}
