package model

import (
	"fmt"
	"testing"

	u "github.com/sunshine69/golang-tools/utils"
)

func TestTenant(t *testing.T) {
	fmt.Println("start test")
	// os.Remove("test.sqlite3")
	SetupDBSchema("../db/schema.sql")
	m0 := map[string]any{"email": "k@k", "join_date": "12/02/2024"}
	m1 := ParseDatetimeFieldOfMapData(m0)
	fmt.Printf("m1: %s\n", u.JsonDump(m1, ""))
	at := GetTenantByCompositeKeyOrNew(map[string]any{"email": "k@k"})
	m1["address"] = "My address"
	u.CheckErr(at.Update(m1), "Update use map")
	at.Where = "email = 'k@k'"
	at1 := at.Search()[0]
	fmt.Printf("%s\n", u.JsonDump(at1, ""))
	at1.Address = "New address stevek"
	u.CheckErr(at1.Save(), "")
	at2 := at.Search()
	fmt.Printf("%s\n", u.JsonDump(at2, ""))
	// tn := Tenant{Address: "%moon%"}
}
