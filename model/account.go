// Generated by model-gen tool. Do not edit but rather edit the template in utils/class-template.go.tmpl
package model

import (
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/mutecomm/go-sqlcipher/v4"
	u "github.com/sunshine69/golang-tools/utils"
	"os"
	"strings"
)

type Account struct {
	Id          int64  `db:"id"`
	Balance     int64  `db:"balance"`
	Contract_id int64  `db:"contract_id,unique"`
	Tenant_main string `db:"tenant_main"`
	Note        string `db:"note" form:"Note,ele=textarea"`
	Where       string `form:"-"`
}

func NewAccount(contract_id int64) Account {

	o := Account{}
	if err := DB.Get(&o, "SELECT * FROM account WHERE  contract_id = ?", contract_id); errors.Is(err, sql.ErrNoRows) {
		o.Contract_id = contract_id
		o.Save()
	}
	// get one and test if exists return as it is
	return o
}

func GetAccountByCompositeKeyOrNew(data map[string]interface{}) *Account {
	data = ParseDatetimeFieldOfMapData(data)
	if rows, err := DB.NamedQuery(`SELECT * FROM account WHERE contract_id=:contract_id `, data); err == nil {
		defer rows.Close()
		for rows.Next() {
			tn := Account{}
			if err = rows.StructScan(&tn); err == nil {
				return &tn
			} else {
				fmt.Fprintf(os.Stderr, "[ERROR] GetAccountByCompositeKey %s\n", err.Error())
				return nil
			}
		}
		// create new one
		tn := NewAccount(data["contract_id"].(int64))
		tn.Update(data)
		return &tn
	} else {
		fmt.Fprintf(os.Stderr, "[ERROR] GetAccountByCompositeKey %s\n", err.Error())
	}
	return nil
}

func GetAccount(contract_id int64) *Account {
	o := Account{
		Contract_id: contract_id,
		Where:       "contract_id=:contract_id "}
	if r := o.Search(); len(r) > 0 {
		return &r[0]
	} else {
		return nil
	}
}

func GetAccountByID(id int64) *Account {
	o := Account{
		Id:    id,
		Where: "id=:id"}
	if r := o.Search(); len(r) > 0 {
		return &r[0]
	} else {
		return nil
	}
}

// Search func
func (o *Account) Search() []Account {
	output := []Account{}
	if o.Where == "" {
		o.Where = ""
	}
	fmt.Println(o.Where)
	if rows, err := DB.NamedQuery(fmt.Sprintf(`SELECT * FROM account WHERE %s`, o.Where), o); err == nil {
		defer rows.Close()
		for rows.Next() {
			_t := Account{}
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
func (o *Account) Update(data map[string]interface{}) error {
	fields := u.MapKeysToSlice(data)
	fieldsWithoutKey := u.SliceMap(fields, func(s string) *string {
		if s != "id" && s != "contract_id" {
			return &s
		}
		return nil
	})
	updateFields := u.SliceMap(fieldsWithoutKey, func(s string) *string { s = s + " = :" + s; return &s })
	updateFieldsStr := strings.Join(updateFields, ",")

	if _, err := DB.NamedExec(`UPDATE account SET `+updateFieldsStr, data); err != nil {
		return err
	}
	return nil
}

// Save existing object which is saved it into db
func (o *Account) Save() error {
	if res, err := DB.NamedExec(`INSERT INTO account(balance,contract_id,tenant_main,note) VALUES(:balance,:contract_id,:tenant_main,:note) ON CONFLICT( contract_id) DO UPDATE SET balance=excluded.balance,contract_id=excluded.contract_id,tenant_main=excluded.tenant_main,note=excluded.note`, o); err != nil {
		return err
	} else {
		o.Id, _ = res.LastInsertId()
	}
	return nil
}

// Delete one object
func (o *Account) Delete() error {
	if res, err := DB.NamedExec(`DELETE FROM account WHERE contract_id=:contract_id `, o); err != nil {
		return err
	} else {
		r, err := res.RowsAffected()
		if err != nil {
			return err
		}
		if r == 0 {
			return fmt.Errorf("ERROR account not found")
		}
	}
	return nil
}

func DeleteAccountByID(id int64) error {
	// sqlx bug? If directly use Exec and sql is a pure string it never delete it but still return ok
	// looks like we always need to bind the named query with sqlx - can not parse pure string in
	if res, err := DB.NamedExec(`DELETE FROM account WHERE id = :id`, map[string]interface{}{"id": id}); err != nil {
		return err
	} else {
		r, err := res.RowsAffected()
		if err != nil {
			return err
		}
		if r == 0 {
			return fmt.Errorf("ERROR account not found")
		}
	}
	return nil
}
