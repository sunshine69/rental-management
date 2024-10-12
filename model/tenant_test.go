package model

import (
	"fmt"
	"github.com/sunshine69/rental-management/utils"
	"testing"

	u "github.com/sunshine69/golang-tools/utils"
)

func TestTenant(t *testing.T) {
	fmt.Println("start test")
	NewTenant("msh.computing@gmail.com")
	NewTenant("msh.computing1@gmail.com")
	k := NewTenant("kaycaonz@gmail.com")
	k.Address = "On the moon"
	k.Save()
	ts := Tenant{Email: "%msh%", Where: "email LIKE :email"}
	o := ts.Search()
	// fmt.Printf("%v\n", k)
	// tn := Tenant{Address: "%moon%"}
	fmt.Printf("%s\n", u.JsonDump(o, "  "))

}

func TestReflect(t *testing.T) {
	o := utils.ReflectStruct(Tenant{Email: "msh@example.com"})

	fmt.Printf("%s\n", u.JsonDump(o, "  "))
}
