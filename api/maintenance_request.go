// Generated by api-gen tool. Do not edit here but edit in api/api-template.go.tmpl
package api

import (
	"fmt"
	"net/http"

	u "github.com/sunshine69/golang-tools/utils"
	"github.com/R167/go-sets"
	"github.com/sunshine69/rental-management/model"
)

func CreateMaintenance_request(w http.ResponseWriter, r *http.Request) {
	if maintenance_request := ParseJSON[model.Maintenance_request](r); maintenance_request != nil {
		if err := maintenance_request.Save(); err != nil {
			fmt.Fprintf(w, `{"status": "ERROR", "msg": "%s"}`, err.Error())
		}
	} else {
		fmt.Fprint(w, `{"status": "ERROR", "msg": "Can not parse json into maintenance_request"}`)
	}
}

func UpdateMaintenance_request(w http.ResponseWriter, r *http.Request) {
	o := ParseJSONToMap(r)
	if id := ParseID(r); id != 0 {
		tn := model.GetMaintenance_requestByID(id)
		if tn == nil {
			fmt.Fprint(w, `{"status": "ERROR", "msg": "Maintenance_request not found with this id"}`)
			return
		}
		o["id"] = id
		if err := tn.Update(o); err == nil {
			fmt.Fprint(w, `{"status": "OK", "msg": "Maintenance_request updated"}`)
		} else {
			fmt.Fprint(w, `{"status": "ERROR", "msg": "`+err.Error()+`"}`)
		}
		return
	} else {
		inputSet := sets.FromMap(o)
		keySet := sets.New("contract_id", "request_date")
		if !inputSet.Superset(keySet) {
			fmt.Fprint(w, `{"status": "ERROR", "msg": "Maintenance_request no key value not provided"}`)
			return
		}
		tn := model.GetMaintenance_requestByCompositeKey(o)
		if tn == nil {
			fmt.Fprint(w, `{"status": "ERROR", "msg": "Maintenance_request not found with this unique field"}`)
			return
		}
		if err := tn.Update(o); err == nil {
			fmt.Fprint(w, `{"status": "OK", "msg": "Maintenance_request updated"}`)
		} else {
			fmt.Fprint(w, `{"status": "ERROR", "msg": "`+err.Error()+`"}`)
		}
		return
	}
}

func DeleteMaintenance_request(w http.ResponseWriter, r *http.Request) {
	if id := ParseID(r); id != 0 {
		if err := model.DeleteMaintenance_requestByID(id); err == nil {
			fmt.Fprint(w, `{"status": "OK", "msg": "Maintenance_request deleted"}`)
		} else {
			fmt.Fprint(w, `{"status": "ERROR", "msg": `+err.Error())
		}
		return
	} else {
		maintenance_request := ParseJSON[model.Maintenance_request](r)
		if maintenance_request != nil {
			if maintenance_request.Id != 0 {
				model.DeleteMaintenance_requestByID(maintenance_request.Id)
				fmt.Fprint(w, `{"status": "OK", "msg": "Maintenance_request deleted"}`)
				return
			} else if maintenance_request.Contract_id != 0   && maintenance_request.Request_date != 0   {
				maintenance_request.Delete()
				fmt.Fprint(w, `{"status": "OK", "msg": "Maintenance_request deleted"}`)
				return
			} else {
				fmt.Fprint(w, `{"status": "ERROR", "msg": "No composite key found"}`)
				return
			}
		} else {
			fmt.Fprint(w, `{"status": "ERROR", "msg": "Cannot parse maintenance_request"}`)
			return
		}
	}
}

func GetMaintenance_request(w http.ResponseWriter, r *http.Request) {
	if id := ParseID(r); id == 0 {
		if where := r.PathValue("where"); where != "" {
			o := model.Maintenance_request{Where: where}
			os := o.Search()
			fmt.Fprint(w, u.JsonDump(os, ""))
			return
		} else {
			os := []model.Maintenance_request{}
			model.DB.Select(&os, `SELECT * FROM maintenance_request ORDER BY id`)
			fmt.Fprint(w, u.JsonDump(os, ""))
			return
		}
	} else {
		o := model.GetMaintenance_requestByID(id)
		fmt.Fprint(w, u.JsonDump(*o, ""))
		return
	}
}

func RouteRegisterMaintenance_request(mux *http.ServeMux, PathBase string) {
	mux.HandleFunc("POST "+PathBase+"/api/maintenance_request", CreateMaintenance_request)
	mux.HandleFunc("PUT "+PathBase+"/api/maintenance_request/{id}", UpdateMaintenance_request)
	mux.HandleFunc("PUT "+PathBase+"/api/maintenance_request", UpdateMaintenance_request)
	mux.HandleFunc("DELETE "+PathBase+"/api/maintenance_request/{id}", DeleteMaintenance_request)
	mux.HandleFunc("DELETE "+PathBase+"/api/maintenance_request", DeleteMaintenance_request)
	mux.HandleFunc("GET "+PathBase+"/api/maintenance_request", GetMaintenance_request)
	mux.HandleFunc("GET "+PathBase+"/api/maintenance_request/{id}", GetMaintenance_request)
	mux.HandleFunc("GET "+PathBase+"/api/maintenance_request/q/{where}", GetMaintenance_request)
}
