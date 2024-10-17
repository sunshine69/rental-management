// Generated by api-gen tool. Do not edit here but edit in api/api-template.go.tmpl
package api

import (
	"fmt"
	"net/http"

	"github.com/R167/go-sets"
	u "github.com/sunshine69/golang-tools/utils"
	"github.com/sunshine69/rental-management/model"
)

func CreateTenant(w http.ResponseWriter, r *http.Request) {
	if tenant := ParseJSON[model.Tenant](r); tenant != nil {
		if err := tenant.Save(); err != nil {
			fmt.Fprintf(w, `{"status": "ERROR", "msg": "%s"}`, err.Error())
		}
	} else {
		fmt.Fprint(w, `{"status": "ERROR", "msg": "Can not parse json into tenant"}`)
	}
}

func UpdateTenant(w http.ResponseWriter, r *http.Request) {
	o := ParseJSONToMap(r)
	if id := ParseID(r); id != 0 {
		tn := model.GetTenantByID(id)
		if tn == nil {
			fmt.Fprint(w, `{"status": "ERROR", "msg": "Tenant not found with this id"}`)
			return
		}
		o["id"] = id
		if err := tn.Update(o); err == nil {
			fmt.Fprint(w, `{"status": "OK", "msg": "Tenant updated"}`)
		} else {
			fmt.Fprint(w, `{"status": "ERROR", "msg": "`+err.Error()+`"}`)
		}
		return
	} else {
		inputSet := sets.FromMap(o)
		keySet := sets.New("email")
		if !inputSet.Superset(keySet) {
			fmt.Fprint(w, `{"status": "ERROR", "msg": "Tenant no key value not provided"}`)
			return
		}
		tn := model.GetTenantByCompositeKeyOrNew(o)
		if err := tn.Update(o); err == nil {
			fmt.Fprint(w, `{"status": "OK", "msg": "Tenant updated"}`)
		} else {
			fmt.Fprint(w, `{"status": "ERROR", "msg": "`+err.Error()+`"}`)
		}
		return
	}
}

func DeleteTenant(w http.ResponseWriter, r *http.Request) {
	if id := ParseID(r); id != 0 {
		if err := model.DeleteTenantByID(id); err == nil {
			fmt.Fprint(w, `{"status": "OK", "msg": "Tenant deleted"}`)
		} else {
			fmt.Fprint(w, `{"status": "ERROR", "msg": `+err.Error())
		}
		return
	} else {
		tenant := ParseJSON[model.Tenant](r)
		if tenant != nil {
			if tenant.Id != 0 {
				model.DeleteTenantByID(tenant.Id)
				fmt.Fprint(w, `{"status": "OK", "msg": "Tenant deleted"}`)
				return
			} else if tenant.Email != "" {
				tenant.Delete()
				fmt.Fprint(w, `{"status": "OK", "msg": "Tenant deleted"}`)
				return
			} else {
				fmt.Fprint(w, `{"status": "ERROR", "msg": "No composite key found"}`)
				return
			}
		} else {
			fmt.Fprint(w, `{"status": "ERROR", "msg": "Cannot parse tenant"}`)
			return
		}
	}
}

func GetTenant(w http.ResponseWriter, r *http.Request) {
	if id := ParseID(r); id == 0 {
		if where := r.PathValue("where"); where != "" {
			o := model.Tenant{Where: where}
			os := o.Search()
			fmt.Fprint(w, u.JsonDump(os, ""))
			return
		} else {
			os := []model.Tenant{}
			model.DB.Select(&os, `SELECT * FROM tenant ORDER BY id`)
			fmt.Fprint(w, u.JsonDump(os, ""))
			return
		}
	} else {
		o := model.GetTenantByID(id)
		fmt.Fprint(w, u.JsonDump(*o, ""))
		return
	}
}

func RouteRegisterTenant(mux *http.ServeMux, PathBase string) {
	mux.HandleFunc("POST "+PathBase+"/api/tenant", CreateTenant)
	mux.HandleFunc("PUT "+PathBase+"/api/tenant/{id}", UpdateTenant)
	mux.HandleFunc("PUT "+PathBase+"/api/tenant", UpdateTenant)
	mux.HandleFunc("DELETE "+PathBase+"/api/tenant/{id}", DeleteTenant)
	mux.HandleFunc("DELETE "+PathBase+"/api/tenant", DeleteTenant)
	mux.HandleFunc("GET "+PathBase+"/api/tenant", GetTenant)
	mux.HandleFunc("GET "+PathBase+"/api/tenant/{id}", GetTenant)
	mux.HandleFunc("GET "+PathBase+"/api/tenant/q/{where}", GetTenant)
}