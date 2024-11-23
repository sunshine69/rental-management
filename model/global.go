package model

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	u "github.com/sunshine69/golang-tools/utils"
	"zombiezen.com/go/sqlite"
	"zombiezen.com/go/sqlite/sqlitex"
)

var DbPool *sqlitex.Pool

func init() {
	DbPath := u.Getenv("DB_PATH", "test.sqlite3")
	DbPool = u.Must(sqlitex.NewPool("file:"+DbPath, sqlitex.PoolOptions{
		PoolSize: 10,
	}))
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

func SetupDBSchema(schemafile string) {
	DB := u.Must(DbPool.Take(context.TODO()))
	defer DbPool.Put(DB)
	fmt.Println("Setup db")
	sqlb, err := os.ReadFile(schemafile)
	u.CheckErr(err, "Read sql file")
	u.CheckErr(sqlitex.ExecuteScript(DB, string(sqlb), nil), "SetUpDB")
}

func GetSqliteCol(stmt *sqlite.Stmt, col_index int) (col_name string, col_val any, go_type string) {
	col_name = stmt.ColumnName(col_index)
	col_type_s := stmt.ColumnType(col_index)
	switch col_type_s {
	case sqlite.TypeInteger:
		col_val = stmt.ColumnInt64(col_index)
		go_type = "int64"
	case sqlite.TypeText:
		col_val = stmt.ColumnText(col_index)
		go_type = "string"
	case sqlite.TypeFloat:
		col_val = stmt.ColumnFloat(col_index)
		go_type = "float64"
	// BLOB and NUL not supported yet
	default:
		println("[WARN] GetSqliteCol un-supported type " + col_type_s.String())
	}
	return
}
