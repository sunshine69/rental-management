// Generated by api-gen tool. Do not edit here but edit in api/api-template.go.tmpl 
package api

import (
	"fmt"
	"net/http"

	u "github.com/sunshine69/golang-tools/utils"

	"github.com/sunshine69/rental-management/model"
)

func CreateInvoice(w http.ResponseWriter, r *http.Request) {
	invoice := ParseJSON[model.Invoice](r)
	invoice.Save()
}

func UpdateInvoice(w http.ResponseWriter, r *http.Request) {
	o := ParseJSON[model.Invoice](r)
	if id := ParseID(r); id != 0 {
		o.Id = id
		o.Save()
		fmt.Fprint(w, `{"status": "OK", "msg": "Invoice updated"}`)
		return
	} else {    
		if o.Number == ""  || o.Issuer == ""  {
			fmt.Fprint(w, `{"status": "ERROR", "msg": "Invoice no key value not provided"}`)
			return
		} else {
			o.Save()
			fmt.Fprint(w, `{"status": "OK", "msg": "Invoice updated"}`)
			return
		}
	}
}

func DeleteInvoice(w http.ResponseWriter, r *http.Request) {
	if id := ParseID(r); id != 0 {
		model.DeleteInvoiceByID(id)
		fmt.Fprint(w, `{"status": "OK", "msg": "Invoice deleted"}`)
		return
	} else {
		invoice := ParseJSON[model.Invoice](r)
		if invoice.Id != 0 {
			invoice.Delete()
			fmt.Fprint(w, `{"status": "OK", "msg": "Invoice deleted"}`)
		}
	}
}

func GetInvoice(w http.ResponseWriter, r *http.Request) {
	if id := ParseID(r); id == 0 {
		if where := r.PathValue("where"); where != "" {
			o := model.Invoice{Where: where}
			os := o.Search()
			fmt.Fprint(w, u.JsonDump(os, ""))
			return
		} else {
			os := []model.Invoice{}
			model.DB.Select(&os, `SELECT * FROM invoice ORDER BY id LIMIT 200`)
			fmt.Fprint(w, u.JsonDump(os, ""))
			return
		}
	} else {
		o := model.GetInvoiceByID(id)
		fmt.Fprint(w, u.JsonDump(*o, ""))
		return
	}
}

func InvoiceRouteRegister(mux *http.ServeMux, PathBase string) {
	mux.HandleFunc("POST "+PathBase+"/invoice", CreateInvoice)
	mux.HandleFunc("PUT "+PathBase+"/invoice/{id}{$}", UpdateInvoice)
	mux.HandleFunc("PUT "+PathBase+"/invoice", UpdateInvoice)
	mux.HandleFunc("DELETE "+PathBase+"/invoice/{id}{$}", DeleteInvoice)
	mux.HandleFunc("DELETE "+PathBase+"/invoice", DeleteInvoice)
	mux.HandleFunc("GET "+PathBase+"/invoice", GetInvoice)
	mux.HandleFunc("GET "+PathBase+"/invoice/{id}{$}", GetInvoice)
	mux.HandleFunc("GET "+PathBase+"/invoice/{where}{$}", GetInvoice)
}
