package main

import (
	"github.com/sunshine69/rental-management/model"
)

func main() {
	model.SetupDBSchema("db/schema.sql")
}
