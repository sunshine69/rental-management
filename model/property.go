// Generated by model-gen tool. Do not edit but rather edit the template in utils/class-template.go.tmpl
package model

import (
	"database/sql"
	"errors"
	"fmt"
	u "github.com/sunshine69/golang-tools/utils"
	_ "modernc.org/sqlite"
	"os"
	"strings"
)

type Property struct {
	Id      int64  `db:"id"`
	Code    string `db:"code,unique"`
	Address string `db:"address"`
	Note    string `db:"note" form:"Note,ele=textarea"`
	Where   string `form:"-"`
}

func NewProperty(code string) Property {

	o := Property{}
	if err := DB.Get(&o, "SELECT * FROM property WHERE  code = ?", code); errors.Is(err, sql.ErrNoRows) {
		o.Code = code
		o.Save()
	}
	// get one and test if exists return as it is
	return o
}

func GetPropertyByCompositeKeyOrNew(data map[string]interface{}) *Property {
	data = ParseDatetimeFieldOfMapData(data)
	if rows, err := DB.NamedQuery(`SELECT * FROM property WHERE code=:code `, data); err == nil {
		defer rows.Close()
		for rows.Next() {
			tn := Property{}
			if err = rows.StructScan(&tn); err == nil {
				return &tn
			} else {
				fmt.Fprintf(os.Stderr, "[ERROR] GetPropertyByCompositeKey %s\n", err.Error())
				return nil
			}
		}
		// create new one
		tn := NewProperty(data["code"].(string))
		tn.Update(data)
		return &tn
	} else {
		fmt.Fprintf(os.Stderr, "[ERROR] GetPropertyByCompositeKey %s\n", err.Error())
	}
	return nil
}

func GetProperty(code string) *Property {
	o := Property{
		Code:  code,
		Where: "code=:code "}
	if r := o.Search(); len(r) > 0 {
		return &r[0]
	} else {
		return nil
	}
}

func GetPropertyByID(id int64) *Property {
	o := Property{
		Id:    id,
		Where: "id=:id"}
	if r := o.Search(); len(r) > 0 {
		return &r[0]
	} else {
		return nil
	}
}

// Search func
func (o *Property) Search() []Property {
	output := []Property{}
	if o.Where == "" {
		o.Where = "code LIKE '%" + o.Code + "%'"
	}
	fmt.Println(o.Where)
	if rows, err := DB.NamedQuery(fmt.Sprintf(`SELECT * FROM property WHERE %s`, o.Where), o); err == nil {
		defer rows.Close()
		for rows.Next() {
			_t := Property{}
			if er := rows.StructScan(&_t); er == nil {
				output = append(output, _t)
			} else {
				fmt.Printf("[ERROR] Scan %s\n", er.Error())
				continue
			}
		}
	} else {
		fmt.Printf("[ERROR] NamedQuery %s\n", err.Error())
	}
	return output
}

// Save new object which is saved it into db
func (o *Property) Update(data map[string]interface{}) error {
	fields := u.MapKeysToSlice(data)
	fieldsWithoutKey := u.SliceMap(fields, func(s string) *string {
		if s != "id" && s != "code" {
			return &s
		}
		return nil
	})
	updateFields := u.SliceMap(fieldsWithoutKey, func(s string) *string { s = s + " = :" + s; return &s })
	updateFieldsStr := strings.Join(updateFields, ",")

	if _, err := DB.NamedExec(`UPDATE property SET `+updateFieldsStr, data); err != nil {
		return err
	}
	return nil
}

// Save existing object which is saved it into db. Note that this will update all fields. If you only update some fields then better use the Update func above
func (o *Property) Save() error {
	if res, err := DB.NamedExec(`INSERT INTO property(code,address,note) VALUES(:code,:address,:note) ON CONFLICT( code) DO UPDATE SET code=excluded.code,address=excluded.address,note=excluded.note`, o); err != nil {
		return err
	} else {
		o.Id, _ = res.LastInsertId()
	}
	return nil
}

// Delete one object
func (o *Property) Delete() error {
	if res, err := DB.NamedExec(`DELETE FROM property WHERE code=:code `, o); err != nil {
		return err
	} else {
		r, err := res.RowsAffected()
		if err != nil {
			return err
		}
		if r == 0 {
			return fmt.Errorf("ERROR property not found")
		}
	}
	return nil
}

func DeletePropertyByID(id int64) error {
	// sqlx bug? If directly use Exec and sql is a pure string it never delete it but still return ok
	// looks like we always need to bind the named query with sqlx - can not parse pure string in
	if res, err := DB.NamedExec(`DELETE FROM property WHERE id = :id`, map[string]interface{}{"id": id}); err != nil {
		return err
	} else {
		r, err := res.RowsAffected()
		if err != nil {
			return err
		}
		if r == 0 {
			return fmt.Errorf("ERROR property not found")
		}
	}
	return nil
}
