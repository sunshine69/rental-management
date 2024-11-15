// Generated by api-gen tool. Do not edit here but edit in api/api-template.go.tmpl
package api

import (
	"fmt"
	"net/http"

	"github.com/R167/go-sets"
	u "github.com/sunshine69/golang-tools/utils"
	"github.com/sunshine69/rental-management/model"
	"github.com/sunshine69/rental-management/utils"
)

func CreatePayment(w http.ResponseWriter, r *http.Request) {
	if payment := u.ParseJsonReqBodyToStruct[model.Payment](r); payment != nil {
		if err := payment.Save(); err != nil {
			fmt.Fprintf(w, `{"status": "ERROR", "msg": "%s"}`, err.Error())
		}
	} else {
		fmt.Fprint(w, `{"status": "ERROR", "msg": "Can not parse json into payment"}`)
	}
}

func UpdatePayment(w http.ResponseWriter, r *http.Request) {
	o := u.ParseJsonReqBodyToMap(r)
	if id := utils.ParseID(r); id != 0 {
		tn := model.GetPaymentByID(id)
		if tn == nil {
			fmt.Fprint(w, `{"status": "ERROR", "msg": "Payment not found with this id"}`)
			return
		}
		o["id"] = id
		if err := tn.Update(o); err == nil {
			fmt.Fprint(w, `{"status": "OK", "msg": "Payment updated"}`)
		} else {
			fmt.Fprint(w, `{"status": "ERROR", "msg": "`+err.Error()+`"}`)
		}
		return
	} else {
		inputSet := sets.FromMap(o)
		keySet := sets.New("account_id", "pay_date")
		if !inputSet.Superset(keySet) {
			fmt.Fprint(w, `{"status": "ERROR", "msg": "Payment no key value not provided"}`)
			return
		}
		tn := model.GetPaymentByCompositeKeyOrNew(o)
		if err := tn.Update(o); err == nil {
			fmt.Fprint(w, `{"status": "OK", "msg": "Payment updated"}`)
		} else {
			fmt.Fprint(w, `{"status": "ERROR", "msg": "`+err.Error()+`"}`)
		}
		return
	}
}

func DeletePayment(w http.ResponseWriter, r *http.Request) {
	if id := utils.ParseID(r); id != 0 {
		if err := model.DeletePaymentByID(id); err == nil {
			fmt.Fprint(w, `{"status": "OK", "msg": "Payment deleted"}`)
		} else {
			fmt.Fprint(w, `{"status": "ERROR", "msg": `+err.Error())
		}
		return
	} else {
		payment := u.ParseJsonReqBodyToStruct[model.Payment](r)
		if payment != nil {
			if payment.Id != 0 {
				model.DeletePaymentByID(payment.Id)
				fmt.Fprint(w, `{"status": "OK", "msg": "Payment deleted"}`)
				return
			} else if payment.Account_id != 0 && payment.Pay_date != "" {
				payment.Delete()
				fmt.Fprint(w, `{"status": "OK", "msg": "Payment deleted"}`)
				return
			} else {
				fmt.Fprint(w, `{"status": "ERROR", "msg": "No composite key found"}`)
				return
			}
		} else {
			fmt.Fprint(w, `{"status": "ERROR", "msg": "Cannot parse payment"}`)
			return
		}
	}
}

func GetPayment(w http.ResponseWriter, r *http.Request) {
	if id := utils.ParseID(r); id == 0 {
		if where := r.PathValue("where"); where != "" {
			o := model.Payment{Where: where}
			os := o.Search()
			fmt.Fprint(w, u.JsonDump(os, ""))
			return
		} else {
			os := []model.Payment{}
			model.DB.Select(&os, `SELECT * FROM payment ORDER BY id`)
			fmt.Fprint(w, u.JsonDump(os, ""))
			return
		}
	} else {
		o := model.GetPaymentByID(id)
		fmt.Fprint(w, u.JsonDump(*o, ""))
		return
	}
}

func RouteRegisterPayment(mux *http.ServeMux, PathBase string) {
	mux.HandleFunc("POST "+PathBase+"/api/payment", CreatePayment)
	mux.HandleFunc("PUT "+PathBase+"/api/payment/{id}", UpdatePayment)
	mux.HandleFunc("PUT "+PathBase+"/api/payment", UpdatePayment)
	mux.HandleFunc("DELETE "+PathBase+"/api/payment/{id}", DeletePayment)
	mux.HandleFunc("DELETE "+PathBase+"/api/payment", DeletePayment)
	mux.HandleFunc("GET "+PathBase+"/api/payment", GetPayment)
	mux.HandleFunc("GET "+PathBase+"/api/payment/{id}", GetPayment)
	mux.HandleFunc("GET "+PathBase+"/api/payment/q/{where}", GetPayment)
}
