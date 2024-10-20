package model

import (
	"fmt"
	"testing"

	"github.com/sunshine69/rental-management/utils"

	u "github.com/sunshine69/golang-tools/utils"
)

func TestTenant(t *testing.T) {
	fmt.Println("start test")
	m0 := map[string]interface{}{"email": "k@k", "start_date": 12333, "end_date": "12/02/2023 00:00:00 +11"}
	m1 := ParseDatetimeFieldOfMapData(m0)
	fmt.Printf("m1: %s\n", u.JsonDump(m1, ""))
	GetTenantByCompositeKeyOrNew(map[string]interface{}{"email": "myf@ptcm"})
	// fmt.Printf("%v\n", k)
	// tn := Tenant{Address: "%moon%"}
}

func TestReflect(t *testing.T) {
	o := utils.ReflectStruct(Tenant{Email: "msh@example.com"}, `form:"([^"]+)"`)

	fmt.Printf("%s\n", u.JsonDump(o, ""))
}
