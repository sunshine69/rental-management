// Generated by api-gen tool. Do not edit here but edit in api/api-template.go.tmpl
package api

import (
	"fmt"
	"net/http"

	"github.com/R167/go-sets"
	u "github.com/sunshine69/golang-tools/utils"
	"github.com/sunshine69/rental-management/model"
)

func CreateInvoice(w http.ResponseWriter, r *http.Request) {
	if invoice := ParseJSON[model.Invoice](r); invoice != nil {
		if err := invoice.Save(); err != nil {
			fmt.Fprintf(w, `{"status": "ERROR", "msg": "%s"}`, err.Error())
		}
	} else {
		fmt.Fprint(w, `{"status": "ERROR", "msg": "Can not parse json into invoice"}`)
	}
}

func UpdateInvoice(w http.ResponseWriter, r *http.Request) {
	o := ParseJSONToMap(r)
	if id := ParseID(r); id != 0 {
		tn := model.GetInvoiceByID(id)
		if tn == nil {
			fmt.Fprint(w, `{"status": "ERROR", "msg": "Invoice not found with this id"}`)
			return
		}
		o["id"] = id
		if err := tn.Update(o); err == nil {
			fmt.Fprint(w, `{"status": "OK", "msg": "Invoice updated"}`)
		} else {
			fmt.Fprint(w, `{"status": "ERROR", "msg": "`+err.Error()+`"}`)
		}
		return
	} else {
		inputSet := sets.FromMap(o)
		keySet := sets.New("number", "issuer")
		if !inputSet.Superset(keySet) {
			fmt.Fprint(w, `{"status": "ERROR", "msg": "Invoice no key value not provided"}`)
			return
		}
		tn := model.GetInvoiceByCompositeKey(o)
		if tn == nil {
			fmt.Fprint(w, `{"status": "ERROR", "msg": "Invoice not found with this unique field"}`)
			return
		}
		if err := tn.Update(o); err == nil {
			fmt.Fprint(w, `{"status": "OK", "msg": "Invoice updated"}`)
		} else {
			fmt.Fprint(w, `{"status": "ERROR", "msg": "`+err.Error()+`"}`)
		}
		return
	}
}

func DeleteInvoice(w http.ResponseWriter, r *http.Request) {
	if id := ParseID(r); id != 0 {
		if err := model.DeleteInvoiceByID(id); err == nil {
			fmt.Fprint(w, `{"status": "OK", "msg": "Invoice deleted"}`)
		} else {
			fmt.Fprint(w, `{"status": "ERROR", "msg": `+err.Error())
		}
		return
	} else {
		invoice := ParseJSON[model.Invoice](r)
		if invoice != nil {
			if invoice.Id != 0 {
				model.DeleteInvoiceByID(invoice.Id)
				fmt.Fprint(w, `{"status": "OK", "msg": "Invoice deleted"}`)
				return
			} else if invoice.Number != "" && invoice.Issuer != "" {
				invoice.Delete()
				fmt.Fprint(w, `{"status": "OK", "msg": "Invoice deleted"}`)
				return
			} else {
				fmt.Fprint(w, `{"status": "ERROR", "msg": "No composite key found"}`)
				return
			}
		} else {
			fmt.Fprint(w, `{"status": "ERROR", "msg": "Cannot parse invoice"}`)
			return
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
			model.DB.Select(&os, `SELECT * FROM invoice ORDER BY id`)
			fmt.Fprint(w, u.JsonDump(os, ""))
			return
		}
	} else {
		o := model.GetInvoiceByID(id)
		fmt.Fprint(w, u.JsonDump(*o, ""))
		return
	}
}

func RouteRegisterInvoice(mux *http.ServeMux, PathBase string) {
	mux.HandleFunc("POST "+PathBase+"/api/invoice", CreateInvoice)
	mux.HandleFunc("PUT "+PathBase+"/api/invoice/{id}", UpdateInvoice)
	mux.HandleFunc("PUT "+PathBase+"/api/invoice", UpdateInvoice)
	mux.HandleFunc("DELETE "+PathBase+"/api/invoice/{id}", DeleteInvoice)
	mux.HandleFunc("DELETE "+PathBase+"/api/invoice", DeleteInvoice)
	mux.HandleFunc("GET "+PathBase+"/api/invoice", GetInvoice)
	mux.HandleFunc("GET "+PathBase+"/api/invoice/{id}", GetInvoice)
	mux.HandleFunc("GET "+PathBase+"/api/invoice/q/{where}", GetInvoice)
}
