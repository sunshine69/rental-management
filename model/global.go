package model

import (
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/mutecomm/go-sqlcipher/v4"
	u "github.com/sunshine69/golang-tools/utils"
)

// DbConn - Global DB connection
var DB *sqlx.DB

func init() {
	DB = sqlx.MustConnect("sqlite3", u.Getenv("DB_PATH", "test.sqlite3"))
}

// Take a map scan if the key contains the string `date` then detect if the value is a string but parsable into a datetime unix int value
// then modify that field into int value. This helps the model to creat/update rows allow parsing date string
func ParseDatetimeFieldOfMapData(input map[string]interface{}) map[string]interface{} {
	for k, v := range input {
		if strings.Contains(k, "date") {
			if _, ok := v.(int64); !ok {
				if v_as_str, ok := v.(string); ok {
					for _, layout := range []string{u.AUTimeLayout, u.CleanStringDateLayout, u.TimeISO8601LayOut} {
						if t, err := time.Parse(layout, v_as_str); err == nil {
							input[k] = t.Unix()
							break
						}
					}
				}
			}
		}
	}
	return input
}

// Auto generate AllModelObjects
var AllForms = map[string]any{

	"Tenant":              Tenant{},
	"Property_manager":    Property_manager{},
	"Property":            Property{},
	"Contract":            Contract{},
	"Account":             Account{},
	"Payment":             Payment{},
	"Invoice":             Invoice{},
	"Maintenance_request": Maintenance_request{},
}
var AllModelObjects []any = []any{Tenant{}, Property_manager{}, Property{}, Contract{}, Account{}, Payment{}, Invoice{}, Maintenance_request{}}

// End generate AllModelObjects
