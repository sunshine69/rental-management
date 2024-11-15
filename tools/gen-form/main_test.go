package main

import (
	"fmt"
	"os"
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
	data := GetTemplateData()

	code := u.GoTemplateString(string(appGoB), data)
	u.BlockInFile("web/app/app.go", []string{}, []string{`// End app-validation.go.tmpl`}, []string{`// Auto generate using app-validation.go.tmpl template`}, code, true, true)
}

func TestAppHanderGen(t *testing.T) {
	os.Chdir("../../")
	fmt.Println("Started tests")
	appGoB, _ := os.ReadFile("tools/gen-form/app-handler.go.tmpl")
	data := GetTemplateData()

	code := u.GoTemplateString(string(appGoB), data)
	u.BlockInFile("web/app/app.go", []string{}, []string{`// End app-handler.go.tmpl`}, []string{`// Auto generate using app-handler.go.tmpl template`}, code, true, true)
}
