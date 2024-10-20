package utils

import (
	"fmt"
	"os"
	"testing"

	"github.com/sunshine69/rental-management/model"
	// u "github.com/sunshine69/golang-tools/utils"
)

func TestFormGen(t *testing.T) {
	os.Chdir("../")
	fmt.Println("Started tests")
	FormGen(model.Tenant{}, "web/app/templates")
}
