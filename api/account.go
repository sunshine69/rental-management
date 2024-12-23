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

func CreateAccount(w http.ResponseWriter, r *http.Request) {
	if account := u.ParseJsonReqBodyToStruct[model.Account](r); account != nil {
		if err := account.Save(); err != nil {
			fmt.Fprintf(w, `{"status": "ERROR", "msg": "%s"}`, err.Error())
		}
	} else {
		fmt.Fprint(w, `{"status": "ERROR", "msg": "Can not parse json into account"}`)
	}
}

func UpdateAccount(w http.ResponseWriter, r *http.Request) {
	o := u.ParseJsonReqBodyToMap(r)
	if id := utils.ParseID(r); id != 0 {
		tn := model.GetAccountByID(id)
		if tn == nil {
			fmt.Fprint(w, `{"status": "ERROR", "msg": "Account not found with this id"}`)
			return
		}
		o["id"] = id
		if err := tn.Update(o); err == nil {
			fmt.Fprint(w, `{"status": "OK", "msg": "Account updated"}`)
		} else {
			fmt.Fprint(w, `{"status": "ERROR", "msg": "`+err.Error()+`"}`)
		}
		return
	} else {
		inputSet := sets.FromMap(o)
		keySet := sets.New("contract_id")
		if !inputSet.Superset(keySet) {
			fmt.Fprint(w, `{"status": "ERROR", "msg": "Account no key value not provided"}`)
			return
		}
		tn := model.GetAccountByCompositeKeyOrNew(o)
		if err := tn.Update(o); err == nil {
			fmt.Fprint(w, `{"status": "OK", "msg": "Account updated"}`)
		} else {
			fmt.Fprint(w, `{"status": "ERROR", "msg": "`+err.Error()+`"}`)
		}
		return
	}
}

func DeleteAccount(w http.ResponseWriter, r *http.Request) {
	if id := utils.ParseID(r); id != 0 {
		if err := model.DeleteAccountByID(id); err == nil {
			fmt.Fprint(w, `{"status": "OK", "msg": "Account deleted"}`)
		} else {
			fmt.Fprint(w, `{"status": "ERROR", "msg": `+err.Error())
		}
		return
	} else {
		account := u.ParseJsonReqBodyToStruct[model.Account](r)
		if account != nil {
			if account.Id != 0 {
				model.DeleteAccountByID(account.Id)
				fmt.Fprint(w, `{"status": "OK", "msg": "Account deleted"}`)
				return
			} else if account.Contract_id != 0 {
				account.Delete()
				fmt.Fprint(w, `{"status": "OK", "msg": "Account deleted"}`)
				return
			} else {
				fmt.Fprint(w, `{"status": "ERROR", "msg": "No composite key found"}`)
				return
			}
		} else {
			fmt.Fprint(w, `{"status": "ERROR", "msg": "Cannot parse account"}`)
			return
		}
	}
}

func GetAccount(w http.ResponseWriter, r *http.Request) {
	if id := utils.ParseID(r); id == 0 {
		if where := r.PathValue("where"); where != "" {
			o := model.Account{Where: where}
			os := o.Search()
			fmt.Fprint(w, u.JsonDump(os, ""))
			return
		} else {
			os := []model.Account{}
			model.DB.Select(&os, `SELECT * FROM account ORDER BY id`)
			fmt.Fprint(w, u.JsonDump(os, ""))
			return
		}
	} else {
		o := model.GetAccountByID(id)
		fmt.Fprint(w, u.JsonDump(*o, ""))
		return
	}
}

func RouteRegisterAccount(mux *http.ServeMux, PathBase string) {
	mux.HandleFunc("POST "+PathBase+"/api/account", CreateAccount)
	mux.HandleFunc("PUT "+PathBase+"/api/account/{id}", UpdateAccount)
	mux.HandleFunc("PUT "+PathBase+"/api/account", UpdateAccount)
	mux.HandleFunc("DELETE "+PathBase+"/api/account/{id}", DeleteAccount)
	mux.HandleFunc("DELETE "+PathBase+"/api/account", DeleteAccount)
	mux.HandleFunc("GET "+PathBase+"/api/account", GetAccount)
	mux.HandleFunc("GET "+PathBase+"/api/account/{id}", GetAccount)
	mux.HandleFunc("GET "+PathBase+"/api/account/q/{where}", GetAccount)
}
