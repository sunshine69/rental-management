package main

import (
	"fmt"
	"strings"

	ag "github.com/sunshine69/automation-go/lib"
	"github.com/sunshine69/rental-management/model"
)

var AllModelObjects []any = []any{model.Tenant{}, model.Property{}, model.Account{}, model.Contract{}, model.Payment{}, model.Maintenance_request{}, model.Property_manager{}, model.Invoice{}}

func GetTemplateData() (data map[string]any) {
	data = map[string]any{}
	objs := []string{}

	for _, it := range AllModelObjects {
		sInfo := ag.ReflectStruct(it, `form:"([^"]+)"`)
		objs = append(objs, sInfo.Name)
		data[sInfo.Name] = sInfo
	}
	data["objs"] = objs
	return
}

// Take all structs in model and generate golang template html form - write to target dir
func FormGen(structType any, writeDirectory string) {
	sInfo := ag.ReflectStruct(structType, `([^\:\s]+):"([^"]+)"`)
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
	ag.GoTemplateFile("tools/gen-form/form.go.tmpl", destFile, data, 0640)
	formNameList := []string{}
	for _, f := range AllModelObjects {
		sInfo := ag.ReflectStruct(f, `form:"([^"]+)"`)
		formNameList = append(formNameList, sInfo.Name)
	}
	ag.GoTemplateFile("tools/gen-form/form-header.go.tmpl", writeDirectory+"/form-header.html", map[string]any{"formNameList": formNameList}, 0640)
	// Generate some common func handler and form validation to copy/paste into the app.go
}

func main() {
	fmt.Println("Form gen")

}
