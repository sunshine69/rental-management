package main

import (
	"fmt"
	"os"

	ag "github.com/sunshine69/automation-go/lib"
	"github.com/sunshine69/rental-management/utils"
)

func main() {
	if sqlb, err := os.ReadFile("db/schema.sql"); err == nil {
		sqls := ag.SplitTextByPattern(string(sqlb), `(?m)^CREATE TABLE IF NOT EXISTS .*`, true)
		// fmt.Printf("%s\n", u.JsonDump(sqls, "  "))
		for _, sqltext := range sqls {
			fmt.Println(sqltext)
			utils.GenerateClass(sqltext, "model/class-template.go.tmpl")
		}
	} else {
		panic(err.Error())
	}
}
