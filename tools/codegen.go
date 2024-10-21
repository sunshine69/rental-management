package main

import (
	"flag"
	"fmt"

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
		for it := range utils.AllModelObjects {
			utils.FormGen(it, "web/app/templates")
		}
	default:
		fmt.Println("Unknown type " + *gentype)
	}
}
