package utils

import (
	"fmt"
	"os"
	"testing"

	// u "github.com/sunshine69/golang-tools/utils"
	ag "github.com/sunshine69/automation-go/lib"
)

func GetTemplateData() (data map[string]any) {
	data = map[string]any{}
	objs := []string{}

	for _, it := range AllModelObjects {
		sInfo := ReflectStruct(it, `form:"([^"]+)"`)
		objs = append(objs, sInfo.Name)
		data[sInfo.Name] = sInfo
	}

	data["objs"] = objs
	return
}

func TestAppValidationGen(t *testing.T) {
	os.Chdir("../")
	fmt.Println("Started tests")
	appGoB, _ := os.ReadFile("utils/app-validation.go.tmpl")
	data := GetTemplateData()

	code := ag.GoTemplateString(string(appGoB), data)
	ag.BlockInFile("web/app/app.go", []string{}, []string{`// End app-validation.go.tmpl`}, []string{`// Auto generate using app-validation.go.tmpl template`}, code, true, true)
}

func TestAppHanderGen(t *testing.T) {
	os.Chdir("../")
	fmt.Println("Started tests")
	appGoB, _ := os.ReadFile("utils/app-handler.go.tmpl")
	data := GetTemplateData()

	code := ag.GoTemplateString(string(appGoB), data)
	ag.BlockInFile("web/app/app.go", []string{}, []string{`// End app-handler.go.tmpl`}, []string{`// Auto generate using app-handler.go.tmpl template`}, code, true, true)
}
