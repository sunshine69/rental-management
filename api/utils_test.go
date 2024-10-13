package api

import (
	"fmt"
	"testing"

	u "github.com/sunshine69/golang-tools/utils"
)

func TestUtils(t *testing.T) {
	instr := `{"first_name": "Steve","last_name": "Kieu", "email": "something@somewhere" }`
	fmt.Printf("%s\n", u.JsonDump(JsonToMap(instr), ""))
}
