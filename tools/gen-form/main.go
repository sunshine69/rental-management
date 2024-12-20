package main

import (
	"fmt"
	"strings"

	u "github.com/sunshine69/golang-tools/utils"
	"github.com/sunshine69/rental-management/model"
)

func GetTemplateData(structTagPtn string) (data map[string]any) {
	data = map[string]any{}
	objs := []string{}

	for _, it := range model.AllModelObjects {
		sInfo := u.ReflectStruct(it, structTagPtn)
		objs = append(objs, sInfo.Name)
		data[sInfo.Name] = sInfo
	}
	data["objs"] = objs
	return
}

// Take all structs in model and generate golang template html form - write to target dir
func FormGen(structType any, writeDirectory string) {
	sInfo := u.ReflectStruct(structType, `([^\:\s]+):"([^"]+)"`)
	destFile := writeDirectory + "/" + sInfo.Name + ".html"
	fList := []string{}
	fieldProp := map[string]map[string]any{}
	for _, v := range sInfo.FieldName {
		fList = append(fList, v)
		fieldProp[v] = map[string]any{"display": true, "ele": "input", "close_ele": "/>", "type": "text", "label": ""}
		tags := sInfo.TagCapture[v]
		for _, tagset := range tags {
			if (tagset[1] == "form" && len(tagset) >= 3 && tagset[2] == "-") || v == "Id" {
				fieldProp[v]["display"] = false
			}
			if tagset[1] == "db" && len(tagset) >= 3 && strings.Contains(tagset[2], "unique") {
				fieldProp[v]["label"] = "*"
			}
			if tagset[1] == "form" && len(tagset) >= 3 && strings.Contains(tagset[2], "ele=textarea") {
				fieldProp[v]["ele"] = "textarea"
				fieldProp[v]["close_ele"] = "</textarea>"
			}
		}
	}
	data := map[string]any{
		"formName":   sInfo.Name,
		"formClass":  "form-group",
		"formAction": "/" + strings.ToLower(sInfo.Name),
		"formID":     strings.ToLower(sInfo.Name),
		"fList":      fList,
		"fieldProp":  fieldProp,
	}
	u.GoTemplateFile("tools/gen-form/form.go.tmpl", destFile, data, 0640)
	formNameList := []string{}
	for _, f := range model.AllModelObjects {
		sInfo := u.ReflectStruct(f, `form:"([^"]+)"`)
		formNameList = append(formNameList, sInfo.Name)
	}
	u.GoTemplateFile("tools/gen-form/form-header.go.tmpl", writeDirectory+"/form-header.html", map[string]any{"formNameList": formNameList}, 0640)
}

func main() {
	fmt.Println("Form gen")

}
