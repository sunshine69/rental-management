package main

import (
	"fmt"
	"os"
	"strings"
	"testing"

	u "github.com/sunshine69/golang-tools/utils"

	"github.com/sunshine69/rental-management/model"
)

func TestFormGen(t *testing.T) {
	os.Chdir("../../")
	for _, it := range model.AllModelObjects {
		FormGen(it, "web/app/templates")
	}
}

func TestAppValidationGen(t *testing.T) {
	os.Chdir("../../")
	fmt.Println("Started tests")
	appGoB, _ := os.ReadFile("tools/gen-form/app-validation.go.tmpl")
	data := GetTemplateData(`form:"([^"]+)"`)

	code := u.GoTemplateString(string(appGoB), data)
	u.BlockInFile("web/app/app.go", []string{}, []string{`// End app-validation.go.tmpl`}, []string{`// Auto generate using app-validation.go.tmpl template`}, code, true, true)
}

func TestAppHanderGen(t *testing.T) {
	os.Chdir("../../")
	fmt.Println("Started tests")
	appGoB, _ := os.ReadFile("tools/gen-form/app-handler.go.tmpl")
	data := GetTemplateData(`db:"([^"]+)"`)
	for _, obj := range data["objs"].([]string) {
		sInfo := data[obj].(u.StructInfo)
		uniqFieldList := []string{}
		// We need to look up inverse, that is from db field => the struct Field
		dbFieldToStructField := map[string]string{}
		for _, f := range sInfo.FieldName {
			allTags := sInfo.TagCapture[f]
			for _, caps := range allTags {
				if len(caps) == 2 { // we have only once cap group as of now so shold only be 2 now
					cap := caps[1]
					capSplit := strings.Split(cap, ",")
					dbFieldName := capSplit[0] // First item is always the col name
					if strings.Contains(cap, "unique") {
						uniqFieldList = append(uniqFieldList, dbFieldName)
					}
					dbFieldToStructField[dbFieldName] = f
				}
			}
		}
		// Create new map field to store the uniq fields. We use it in the template
		data[obj+"_uniqFields"] = uniqFieldList
		data[obj+"_dbFieldToStructField"] = dbFieldToStructField
	}
	code := u.GoTemplateString(string(appGoB), data)
	u.BlockInFile("web/app/app.go", []string{}, []string{`// End app-handler.go.tmpl`}, []string{`// Auto generate using app-handler.go.tmpl template`}, code, true, true)
	u.RunSystemCommandV2("go fmt web/app/app.go", false)
}
