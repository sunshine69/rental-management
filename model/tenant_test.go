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
	m0 := map[string]interface{}{"email": "k@k", "start_date": 12333, "end_date": "12/02/2023 00:00:00 +11"}
	m1 := ParseDatetimeFieldOfMapData(m0)
	fmt.Printf("m1: %s\n", u.JsonDump(m1, ""))
	GetTenantByCompositeKeyOrNew(map[string]interface{}{"email": "myf@ptcm"})
	// at := Tenant{Email: "msh@come"}
	// at.Save()
	at := Tenant{Where: "email like '%msh%'"}
	at1 := at.Search()[0]
	fmt.Printf("%s\n", u.JsonDump(at1, ""))
	at1.Address = "New address stevek "
	u.CheckErr(at1.Save(), "")
	fmt.Printf("%s\n", u.JsonDump(at1, ""))
	// tn := Tenant{Address: "%moon%"}
}
