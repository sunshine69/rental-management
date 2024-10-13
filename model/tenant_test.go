package model

import (
	"fmt"
	"testing"

	"github.com/sunshine69/rental-management/utils"

	u "github.com/sunshine69/golang-tools/utils"
)

func TestTenant(t *testing.T) {
	fmt.Println("start test")
	NewTenant("msh.computing@gmail.com")
	NewTenant("msh.computing1@gmail.com")
	k := NewTenant("kaycaonz@gmail.com")
	k.Address = "On the moon"
	k.Save()

	at := Tenant{Email: "andrew@bla"}
	at.Save()
	at1 := Tenant{Address: "Only address"}
	at1.Save()
	at2 := Tenant{Address: "Only address 2"}
	at2.Save()
	at3 := GetTenant("msh.computing@gmail.com")
	fmt.Printf("AT3 SEARCH BY ID %v\n", *at3)
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
