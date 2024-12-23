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

type Contract struct {
	Id                 int64  `db:"id"`
	Property           string `db:"property,unique"`
	Property_manager   string `db:"property_manager"`
	Tenant_main        string `db:"tenant_main,unique"`
	Tenants            string `db:"tenants"`
	Start_date         string `db:"start_date,unique"`
	End_date           string `db:"end_date"`
	Signed_date        string `db:"signed_date"`
	Term               string `db:"term"`
	Rent               int64  `db:"rent"`
	Rent_period        string `db:"rent_period"`
	Rent_paid_on       string `db:"rent_paid_on"`
	Water_charged      int64  `db:"water_charged"`
	Document_file_path string `db:"document_file_path"`
	Url                string `db:"url"`
	Note               string `db:"note" form:"Note,ele=textarea"`
	Where              string `form:"-"`
}

func NewContract(property string, start_date string, tenant_main string) Contract {

	o := Contract{}
	if err := DB.Get(&o, "SELECT * FROM contract WHERE  property = ? AND  start_date = ? AND  tenant_main = ?", property, start_date, tenant_main); errors.Is(err, sql.ErrNoRows) {
		o.Property = property
		o.Start_date = start_date
		o.Tenant_main = tenant_main
		if o.Start_date == "" {
			o.Start_date = time.Now().Format(u.TimeISO8601LayOut)
		}
		if o.End_date == "" {
			o.End_date = time.Now().Format(u.TimeISO8601LayOut)
		}
		if o.Signed_date == "" {
			o.Signed_date = time.Now().Format(u.TimeISO8601LayOut)
		}
		o.Save()
	}
	// get one and test if exists return as it is
	return o
}

func GetContractByCompositeKeyOrNew(data map[string]interface{}) *Contract {
	data = ParseDatetimeFieldOfMapData(data)
	if rows, err := DB.NamedQuery(`SELECT * FROM contract WHERE property=:property  AND start_date=:start_date  AND tenant_main=:tenant_main `, data); err == nil {
		defer rows.Close()
		for rows.Next() {
			tn := Contract{}
			if err = rows.StructScan(&tn); err == nil {
				return &tn
			} else {
				fmt.Fprintf(os.Stderr, "[ERROR] GetContractByCompositeKey %s\n", err.Error())
				return nil
			}
		}
		// create new one
		tn := NewContract(data["property"].(string), data["start_date"].(string), data["tenant_main"].(string))
		tn.Update(data)
		return &tn
	} else {
		fmt.Fprintf(os.Stderr, "[ERROR] GetContractByCompositeKey %s\n", err.Error())
	}
	return nil
}

func GetContract(property string, start_date string, tenant_main string) *Contract {
	o := Contract{
		Property: property, Start_date: start_date, Tenant_main: tenant_main,
		Where: "property=:property , start_date=:start_date , tenant_main=:tenant_main "}
	if r := o.Search(); len(r) > 0 {
		return &r[0]
	} else {
		return nil
	}
}

func GetContractByID(id int64) *Contract {
	o := Contract{
		Id:    id,
		Where: "id=:id"}
	if r := o.Search(); len(r) > 0 {
		return &r[0]
	} else {
		return nil
	}
}

// Search func
func (o *Contract) Search() []Contract {
	output := []Contract{}
	if o.Where == "" {
		o.Where = "property LIKE '%" + o.Property + "%' AND start_date LIKE '%" + o.Start_date + "%' AND tenant_main LIKE '%" + o.Tenant_main + "%'"
	}
	fmt.Println(o.Where)
	if rows, err := DB.NamedQuery(fmt.Sprintf(`SELECT * FROM contract WHERE %s`, o.Where), o); err == nil {
		defer rows.Close()
		for rows.Next() {
			_t := Contract{}
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
func (o *Contract) Update(data map[string]interface{}) error {
	fields := u.MapKeysToSlice(data)
	fieldsWithoutKey := u.SliceMap(fields, func(s string) *string {
		if s != "id" && s != "property" && s != "start_date" && s != "tenant_main" {
			return &s
		}
		return nil
	})
	updateFields := u.SliceMap(fieldsWithoutKey, func(s string) *string { s = s + " = :" + s; return &s })
	updateFieldsStr := strings.Join(updateFields, ",")

	if _, err := DB.NamedExec(`UPDATE contract SET `+updateFieldsStr, data); err != nil {
		return err
	}
	return nil
}

// Save existing object which is saved it into db
func (o *Contract) Save() error {
	if res, err := DB.NamedExec(`INSERT INTO contract(property,property_manager,tenant_main,tenants,start_date,end_date,signed_date,term,rent,rent_period,rent_paid_on,water_charged,document_file_path,url,note) VALUES(:property,:property_manager,:tenant_main,:tenants,:start_date,:end_date,:signed_date,:term,:rent,:rent_period,:rent_paid_on,:water_charged,:document_file_path,:url,:note) ON CONFLICT( property,start_date,tenant_main) DO UPDATE SET property=excluded.property,property_manager=excluded.property_manager,tenant_main=excluded.tenant_main,tenants=excluded.tenants,start_date=excluded.start_date,end_date=excluded.end_date,signed_date=excluded.signed_date,term=excluded.term,rent=excluded.rent,rent_period=excluded.rent_period,rent_paid_on=excluded.rent_paid_on,water_charged=excluded.water_charged,document_file_path=excluded.document_file_path,url=excluded.url,note=excluded.note`, o); err != nil {
		return err
	} else {
		o.Id, _ = res.LastInsertId()
	}
	return nil
}

// Delete one object
func (o *Contract) Delete() error {
	if res, err := DB.NamedExec(`DELETE FROM contract WHERE property=:property , start_date=:start_date , tenant_main=:tenant_main `, o); err != nil {
		return err
	} else {
		r, err := res.RowsAffected()
		if err != nil {
			return err
		}
		if r == 0 {
			return fmt.Errorf("ERROR contract not found")
		}
	}
	return nil
}

func DeleteContractByID(id int64) error {
	// sqlx bug? If directly use Exec and sql is a pure string it never delete it but still return ok
	// looks like we always need to bind the named query with sqlx - can not parse pure string in
	if res, err := DB.NamedExec(`DELETE FROM contract WHERE id = :id`, map[string]interface{}{"id": id}); err != nil {
		return err
	} else {
		r, err := res.RowsAffected()
		if err != nil {
			return err
		}
		if r == 0 {
			return fmt.Errorf("ERROR contract not found")
		}
	}
	return nil
}
