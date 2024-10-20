package main

import (
	"flag"
	"fmt"
	"github.com/sunshine69/rental-management/model"
	"github.com/sunshine69/rental-management/utils"
)

func main() {
	gentype := flag.String("type", "", "Code gen type. Can be model | api")
	flag.Parse()

	switch *gentype {
	case "model":
		utils.CodeGen("model/class-template.go.tmpl")
	case "api":
		utils.CodeGen("api/api-template.go.tmpl")
	case "form":
		utils.FormGen(model.Tenant{}, "web/app/templates")
		utils.FormGen(model.Property{}, "web/app/templates")
		utils.FormGen(model.Account{}, "web/app/templates")
		utils.FormGen(model.Payment{}, "web/app/templates")
		utils.FormGen(model.Contract{}, "web/app/templates")
		utils.FormGen(model.Invoice{}, "web/app/templates")
		utils.FormGen(model.Maintenance_request{}, "web/app/templates")
		utils.FormGen(model.Property_manager{}, "web/app/templates")
	default:
		fmt.Println("Unknown type " + *gentype)
	}
}
