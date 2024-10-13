package main

import (
	"flag"
	"fmt"
	"os"

	ag "github.com/sunshine69/automation-go/lib"
	"github.com/sunshine69/rental-management/utils"
)

func CodeGen(templateFile string) {
	if sqlb, err := os.ReadFile("db/schema.sql"); err == nil {
		sqls := ag.SplitTextByPattern(string(sqlb), `(?m)^CREATE TABLE IF NOT EXISTS .*`, true)
		// fmt.Printf("%s\n", u.JsonDump(sqls, "  "))
		for _, sqltext := range sqls {
			fmt.Println(sqltext)
			utils.GenerateClass(sqltext, templateFile)
		}
	} else {
		panic(err.Error())
	}
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
