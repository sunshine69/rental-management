// Generated by model-gen tool. Do not edit but rather edit the template in utils/class-template.go.tmpl
package model

import (
	"database/sql"
	"errors"
	"fmt"
	u "github.com/sunshine69/golang-tools/utils"
	"os"
	"strings"
	"time"

	_ "github.com/mutecomm/go-sqlcipher/v4"
	ag "github.com/sunshine69/automation-go/lib"
)

type Invoice struct {
	Id          int64  `db:"id"`
	Date        string `db:"date"`
	Description string `db:"description"`
	Amount      int64  `db:"amount"`
	Number      string `db:"number,unique"`
	Issuer      string `db:"issuer,unique"`
	Payer       string `db:"payer"`
	Property    string `db:"property"`
	Due_date    string `db:"due_date"`
	Where       string `form:"-"`
}

func NewInvoice(number string, issuer string) Invoice {

	o := Invoice{}
	if err := DB.Get(&o, "SELECT * FROM invoice WHERE  number = ? AND  issuer = ?", number, issuer); errors.Is(err, sql.ErrNoRows) {
		o.Number = number
		o.Issuer = issuer
		if o.Date == "" {
			o.Date = time.Now().Format(u.TimeISO8601LayOut)
		}
		if o.Due_date == "" {
			o.Due_date = time.Now().Format(u.TimeISO8601LayOut)
		}
		o.Save()
	}
	// get one and test if exists return as it is
	return o
}

func GetInvoiceByCompositeKeyOrNew(data map[string]interface{}) *Invoice {
	data = ParseDatetimeFieldOfMapData(data)
	if rows, err := DB.NamedQuery(`SELECT * FROM invoice WHERE number=:number  AND issuer=:issuer `, data); err == nil {
		defer rows.Close()
		for rows.Next() {
			tn := Invoice{}
			if err = rows.StructScan(&tn); err == nil {
				return &tn
			} else {
				fmt.Fprintf(os.Stderr, "[ERROR] GetInvoiceByCompositeKey %s\n", err.Error())
				return nil
			}
		}
		// create new one
		tn := NewInvoice(data["number"].(string), data["issuer"].(string))
		tn.Update(data)
		return &tn
	} else {
		fmt.Fprintf(os.Stderr, "[ERROR] GetInvoiceByCompositeKey %s\n", err.Error())
	}
	return nil
}

func GetInvoice(number string, issuer string) *Invoice {
	o := Invoice{
		Number: number, Issuer: issuer,
		Where: "number=:number , issuer=:issuer "}
	if r := o.Search(); len(r) > 0 {
		return &r[0]
	} else {
		return nil
	}
}

func GetInvoiceByID(id int64) *Invoice {
	o := Invoice{
		Id:    id,
		Where: "id=:id"}
	if r := o.Search(); len(r) > 0 {
		return &r[0]
	} else {
		return nil
	}
}

// Search func
func (o *Invoice) Search() []Invoice {
	output := []Invoice{}
	if o.Where == "" {
		o.Where = "number LIKE '%" + o.Number + "%' AND issuer LIKE '%" + o.Issuer + "%'"
	}
	fmt.Println(o.Where)
	if rows, err := DB.NamedQuery(fmt.Sprintf(`SELECT * FROM invoice WHERE %s`, o.Where), o); err == nil {
		defer rows.Close()
		for rows.Next() {
			_t := Invoice{}
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
func (o *Invoice) Update(data map[string]interface{}) error {
	fields := ag.MapKeysToSlice(data)
	fieldsWithoutKey := ag.SliceMap(fields, func(s string) *string {
		if s != "id" && s != "number" && s != "issuer" {
			return &s
		}
		return nil
	})
	updateFields := ag.SliceMap(fieldsWithoutKey, func(s string) *string { s = s + " = :" + s; return &s })
	updateFieldsStr := strings.Join(updateFields, ",")

	if _, err := DB.NamedExec(`UPDATE invoice SET `+updateFieldsStr, data); err != nil {
		return err
	}
	return nil
}

// Save existing object which is saved it into db
func (o *Invoice) Save() error {
	if res, err := DB.NamedExec(`INSERT INTO invoice(date,description,amount,number,issuer,payer,property,due_date) VALUES(:date,:description,:amount,:number,:issuer,:payer,:property,:due_date) ON CONFLICT( number,issuer) DO UPDATE SET date=excluded.date,description=excluded.description,amount=excluded.amount,number=excluded.number,issuer=excluded.issuer,payer=excluded.payer,property=excluded.property,due_date=excluded.due_date`, o); err != nil {
		return err
	} else {
		o.Id, _ = res.LastInsertId()
	}
	return nil
}

// Delete one object
func (o *Invoice) Delete() error {
	if res, err := DB.NamedExec(`DELETE FROM invoice WHERE number=:number , issuer=:issuer `, o); err != nil {
		return err
	} else {
		r, err := res.RowsAffected()
		if err != nil {
			return err
		}
		if r == 0 {
			return fmt.Errorf("ERROR invoice not found")
		}
	}
	return nil
}

func DeleteInvoiceByID(id int64) error {
	// sqlx bug? If directly use Exec and sql is a pure string it never delete it but still return ok
	// looks like we always need to bind the named query with sqlx - can not parse pure string in
	if res, err := DB.NamedExec(`DELETE FROM invoice WHERE id = :id`, map[string]interface{}{"id": id}); err != nil {
		return err
	} else {
		r, err := res.RowsAffected()
		if err != nil {
			return err
		}
		if r == 0 {
			return fmt.Errorf("ERROR invoice not found")
		}
	}
	return nil
}
