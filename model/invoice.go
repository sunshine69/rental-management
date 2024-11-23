// Generated by model-gen tool. Do not edit but rather edit the template in utils/class-template.go.tmpl
package model

import (
	"context"
	"strings"
	"time"

	u "github.com/sunshine69/golang-tools/utils"
	"zombiezen.com/go/sqlite"
	"zombiezen.com/go/sqlite/sqlitex"
)

type Invoice struct {
	Id            int64          `db:"id"`
	Date          string         `db:"date"`
	Description   string         `db:"description"`
	Amount        int64          `db:"amount"`
	Number        string         `db:"number,unique"`
	Issuer        string         `db:"issuer,unique"`
	Payer         string         `db:"payer"`
	Property      string         `db:"property"`
	Due_date      string         `db:"due_date"`
	Where         string         `form:"-"`
	WhereNamedArg map[string]any `form:"-"`
}

func ParseInvoiceFromStmt(stmt *sqlite.Stmt) (o Invoice) {
	for idx := 0; idx < stmt.ColumnCount(); idx++ {
		col_name, col_val, _ := GetSqliteCol(stmt, idx)
		switch col_name {
		case "id":
			o.Id = col_val.(int64)
		case "date":
			o.Date = col_val.(string)
			if o.Date == "" {
				o.Date = time.Now().Format(u.TimeISO8601LayOut)
			}
		case "description":
			o.Description = col_val.(string)
		case "amount":
			o.Amount = col_val.(int64)
		case "number":
			o.Number = col_val.(string)
		case "issuer":
			o.Issuer = col_val.(string)
		case "payer":
			o.Payer = col_val.(string)
		case "property":
			o.Property = col_val.(string)
		case "due_date":
			o.Due_date = col_val.(string)
			if o.Due_date == "" {
				o.Due_date = time.Now().Format(u.TimeISO8601LayOut)
			}
		}
	}
	return
}

func NewInvoice(number string, issuer string) Invoice {
	o := Invoice{Where: ` number = :number AND  issuer = :issuer`, WhereNamedArg: map[string]any{":number": number, ":issuer": issuer}}
	output := o.Search()
	if len(output) == 0 {
		o.Number = number
		o.Issuer = issuer
		o.Save()
	} else {
		o = output[0]
	}
	return o
}

func GetInvoiceByCompositeKeyOrNew(data map[string]interface{}) *Invoice {
	DB := u.Must(DbPool.Take(context.TODO()))
	defer DbPool.Put(DB)
	t := Invoice{}
	data = ParseDatetimeFieldOfMapData(data)
	err := sqlitex.Execute(DB, `SELECT * FROM invoice WHERE  number = :number AND  issuer = :issuer`, &sqlitex.ExecOptions{
		Named: map[string]any{":number": data["number"], ":issuer": data["issuer"]},
		ResultFunc: func(stmt *sqlite.Stmt) error {
			t = ParseInvoiceFromStmt(stmt)
			return nil
		},
	})
	if err == nil && t.Id != 0 {
		return &t
	} else {
		// create new one
		tn := NewInvoice(data["number"].(string), data["issuer"].(string))
		tn.Update(data)
		return &tn
	}
}

func GetInvoice(number string, issuer string) *Invoice {
	o := Invoice{
		Number: number, Issuer: issuer,
		Where:         " number = :number AND  issuer = :issuer",
		WhereNamedArg: map[string]any{":number": number, ":issuer": issuer},
	}
	if r := o.Search(); len(r) > 0 {
		return &r[0]
	} else {
		return nil
	}
}

func GetInvoiceByID(id int64) *Invoice {
	o := Invoice{
		Id:            id,
		Where:         "id=:id",
		WhereNamedArg: map[string]any{":id": id},
	}
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
		o.Where = "true"
		if len(o.WhereNamedArg) == 0 {
			o.WhereNamedArg = map[string]any{}
		}
	}
	DB := u.Must(DbPool.Take(context.TODO()))
	defer DbPool.Put(DB)
	err := sqlitex.Execute(DB, "SELECT * FROM tenant WHERE "+o.Where, &sqlitex.ExecOptions{
		Named: o.WhereNamedArg,
		ResultFunc: func(stmt *sqlite.Stmt) error {
			t := ParseInvoiceFromStmt(stmt)
			output = append(output, t)
			return nil
		},
	})
	if err != nil {
		println("[ERROR] ", err.Error())
	}
	return output
}

// Save new object which is saved it into db
func (o *Invoice) Update(data map[string]interface{}) error {
	fields := u.MapKeysToSlice(data)
	fieldsWithoutKey := u.SliceMap(fields, func(s string) *string {
		if s != "id" && s != "number" && s != "issuer" {
			return &s
		}
		return nil
	})
	namedArgs := map[string]any{}
	updateFields := u.SliceMap(fieldsWithoutKey, func(s string) *string {
		s = s + " = :" + s
		namedArgs[":"+s] = data[s]
		return &s
	})
	updateFieldsStr := strings.Join(updateFields, ",")

	DB := u.Must(DbPool.Take(context.TODO()))
	defer DbPool.Put(DB)
	return sqlitex.Execute(DB, `UPDATE invoice SET `+updateFieldsStr, &sqlitex.ExecOptions{
		Named: namedArgs,
	})
}

// Save existing object which is saved it into db. Note that this will update all fields. If you only update some fields then better use the Update func above
func (o *Invoice) Save() error {
	DB := u.Must(DbPool.Take(context.TODO()))
	defer DbPool.Put(DB)
	sqlstr := `INSERT INTO invoice(date,description,amount,number,issuer,payer,property,due_date) VALUES(:date,:description,:amount,:number,:issuer,:payer,:property,:due_date) ON CONFLICT( number,issuer) DO UPDATE SET date=excluded.date,description=excluded.description,amount=excluded.amount,number=excluded.number,issuer=excluded.issuer,payer=excluded.payer,property=excluded.property,due_date=excluded.due_date`
	err := sqlitex.Execute(DB, sqlstr, &sqlitex.ExecOptions{
		Named: map[string]any{":id": o.Id, ":date": o.Date, ":description": o.Description, ":amount": o.Amount, ":number": o.Number, ":issuer": o.Issuer, ":payer": o.Payer, ":property": o.Property, ":due_date": o.Due_date},
	})
	if err != nil {
		return err
	}
	if DB.Changes() > 0 {
		o.Id = DB.LastInsertRowID()
	}
	return nil
}

// Delete one object
func (o *Invoice) Delete() error {
	DB := u.Must(DbPool.Take(context.TODO()))
	defer DbPool.Put(DB)
	sqlstr := `DELETE FROM invoice WHERE  number = :number AND  issuer = :issuer`
	return sqlitex.Execute(DB, sqlstr, &sqlitex.ExecOptions{
		Named: map[string]any{":number": o.Number, ":issuer": o.Issuer},
	})
}

func DeleteInvoiceByID(id int64) error {
	DB := u.Must(DbPool.Take(context.TODO()))
	defer DbPool.Put(DB)
	sqlstr := `DELETE FROM invoice WHERE id = :id`
	return sqlitex.Execute(DB, sqlstr, &sqlitex.ExecOptions{
		Named: map[string]any{
			":id": id,
		},
	})
}
