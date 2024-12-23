// Generated by model-gen tool. Do not edit but rather edit the template in utils/class-template.go.tmpl
package model

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	_ "github.com/mutecomm/go-sqlcipher/v4"
	u "github.com/sunshine69/golang-tools/utils"
)

type Tenant struct {
	Id             int64  `db:"id"`
	First_name     string `db:"first_name"`
	Last_name      string `db:"last_name"`
	Address        string `db:"address"`
	Contact_number string `db:"contact_number"`
	Email          string `db:"email,unique"`
	Join_date      string `db:"join_date"`
	Note           string `db:"note" form:"Note,ele=textarea"`
	Where          string `form:"-"`
}

func NewTenant(email string) Tenant {

	o := Tenant{}
	if err := DB.Get(&o, "SELECT * FROM tenant WHERE  email = ?", email); errors.Is(err, sql.ErrNoRows) {
		o.Email = email
		if o.Join_date == "" {
			o.Join_date = time.Now().Format(u.TimeISO8601LayOut)
		}
		o.Save()
	}
	// get one and test if exists return as it is
	return o
}

func GetTenantByCompositeKeyOrNew(data map[string]interface{}) *Tenant {
	data = ParseDatetimeFieldOfMapData(data)
	if rows, err := DB.NamedQuery(`SELECT * FROM tenant WHERE email=:email `, data); err == nil {
		defer rows.Close()
		for rows.Next() {
			tn := Tenant{}
			if err = rows.StructScan(&tn); err == nil {
				return &tn
			} else {
				fmt.Fprintf(os.Stderr, "[ERROR] GetTenantByCompositeKey %s\n", err.Error())
				return nil
			}
		}
		// create new one
		tn := NewTenant(data["email"].(string))
		tn.Update(data)
		return &tn
	} else {
		fmt.Fprintf(os.Stderr, "[ERROR] GetTenantByCompositeKey %s\n", err.Error())
	}
	return nil
}

func GetTenant(email string) *Tenant {
	o := Tenant{
		Email: email,
		Where: "email=:email "}
	if r := o.Search(); len(r) > 0 {
		return &r[0]
	} else {
		return nil
	}
}

func GetTenantByID(id int64) *Tenant {
	o := Tenant{
		Id:    id,
		Where: "id=:id"}
	if r := o.Search(); len(r) > 0 {
		return &r[0]
	} else {
		return nil
	}
}

// Search func
func (o *Tenant) Search() []Tenant {
	output := []Tenant{}
	if o.Where == "" {
		o.Where = "email LIKE '%" + o.Email + "%'"
	}
	fmt.Println(o.Where)
	if rows, err := DB.NamedQuery(fmt.Sprintf(`SELECT * FROM tenant WHERE %s`, o.Where), o); err == nil {
		defer rows.Close()
		for rows.Next() {
			_t := Tenant{}
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
func (o *Tenant) Update(data map[string]interface{}) error {
	fields := u.MapKeysToSlice(data)
	fieldsWithoutKey := u.SliceMap(fields, func(s string) *string {
		if s != "id" && s != "email" {
			return &s
		}
		return nil
	})
	updateFields := u.SliceMap(fieldsWithoutKey, func(s string) *string { s = s + " = :" + s; return &s })
	updateFieldsStr := strings.Join(updateFields, ",")

	if _, err := DB.NamedExec(`UPDATE tenant SET `+updateFieldsStr, data); err != nil {
		return err
	}
	return nil
}

// Save existing object which is saved it into db
func (o *Tenant) Save() error {
	if res, err := DB.NamedExec(`INSERT INTO tenant(first_name,last_name,address,contact_number,email,join_date,note) VALUES(:first_name,:last_name,:address,:contact_number,:email,:join_date,:note) ON CONFLICT( email) DO UPDATE SET first_name=excluded.first_name,last_name=excluded.last_name,address=excluded.address,contact_number=excluded.contact_number,email=excluded.email,join_date=excluded.join_date,note=excluded.note`, o); err != nil {
		return err
	} else {
		o.Id, _ = res.LastInsertId()
	}
	return nil
}

// Delete one object
func (o *Tenant) Delete() error {
	if res, err := DB.NamedExec(`DELETE FROM tenant WHERE email=:email `, o); err != nil {
		return err
	} else {
		r, err := res.RowsAffected()
		if err != nil {
			return err
		}
		if r == 0 {
			return fmt.Errorf("ERROR tenant not found")
		}
	}
	return nil
}

func DeleteTenantByID(id int64) error {
	// sqlx bug? If directly use Exec and sql is a pure string it never delete it but still return ok
	// looks like we always need to bind the named query with sqlx - can not parse pure string in
	if res, err := DB.NamedExec(`DELETE FROM tenant WHERE id = :id`, map[string]interface{}{"id": id}); err != nil {
		return err
	} else {
		r, err := res.RowsAffected()
		if err != nil {
			return err
		}
		if r == 0 {
			return fmt.Errorf("ERROR tenant not found")
		}
	}
	return nil
}
