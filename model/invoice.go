// Generated by model-gen tool. Do not edit but rather edit the template in utils/class-template.go.tmpl 
package model

import (
	"database/sql"
	"errors"
	"fmt"
	"time"
	
	_ "github.com/mutecomm/go-sqlcipher/v4"
)

type Invoice struct {
	Amount int64 `db:"amount"`
	Date int64 `db:"date"`
	Description string `db:"description"`
	Due_date int64 `db:"due_date"`
	Id int64 `db:"id"`
	Issuer string `db:"issuer"`
	Number string `db:"number"`
	Property_id int64 `db:"property_id"`
	To string `db:"to"`
	
	Where string 
}

func NewInvoice(number string ,issuer string ) Invoice {

	o := Invoice{}
	if err := DB.Get(&o, "SELECT * FROM invoice WHERE  number = ? AND  issuer = ?",number ,issuer ); errors.Is(err, sql.ErrNoRows) {		
		o.Number = number		
		o.Issuer = issuer
		if o.Date == 0 {
			o.Date = time.Now().Unix()
		}
		if o.Due_date == 0 {
			o.Due_date = time.Now().Unix()
		}
		o.Save()
	}
	// get one and test if exists return as it is
	return o	
}

func GetInvoice(number string, issuer string) *Invoice {
	o := Invoice{
		Number: number , Issuer: issuer , 
		Where: "number=:number , issuer=:issuer "}
	if r := o.Search(); r != nil {
		return &r[0]
	} else {
		return nil
	}
}

func GetInvoiceByID(id int64) *Invoice {
	o := Invoice{
		Id: id,
		Where: "id=:id"}
	if r := o.Search(); r != nil {
		return &r[0]
	} else {
		return nil
	}
}

// Search func
func (o *Invoice) Search() []Invoice {
	output := []Invoice{}
	if rows, err := DB.NamedQuery(fmt.Sprintf(`SELECT * FROM invoice WHERE %s`, o.Where), o); err == nil {
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

// Save existing object which is saved it into db 
func (o *Invoice) Save() {
	if res, err := DB.NamedExec(`INSERT INTO invoice(date,due_date,description,amount,number,issuer,to,property_id ) VALUES(:date,:due_date,:description,:amount,:number,:issuer,:to,:property_id ) ON CONFLICT(number,issuer) DO UPDATE SET date=excluded.date,due_date=excluded.due_date,description=excluded.description,amount=excluded.amount,number=excluded.number,issuer=excluded.issuer,to=excluded.to,property_id=excluded.property_id`, o); err != nil {
		fmt.Printf("[ERROR] %s\n", err.Error())
	} else {
		o.Id, _ = res.LastInsertId()
	}
}

// Delete one object
func (o *Invoice) Delete() {
	if _, err := DB.NamedExec(`DELETE FROM invoice WHERE number=:number AND issuer=:issuer`, o); err != nil {
		fmt.Printf("[ERROR] %s\n", err.Error())
	} else {
		o = nil
	}
}

func DeleteInvoiceByID(id int64) {
	if _, err := DB.NamedExec(`DELETE FROM invoice WHERE id=?`, id); err != nil {
		fmt.Printf("[ERROR] %s\n", err.Error())
	}
}