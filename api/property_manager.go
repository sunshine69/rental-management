// Generated by api-gen tool. Do not edit here but edit in api/api-template.go.tmpl 
package api

import (
	"fmt"
	"net/http"

	u "github.com/sunshine69/golang-tools/utils"

	"github.com/sunshine69/rental-management/model"
)

func CreateProperty_manager(w http.ResponseWriter, r *http.Request) {
	property_manager := ParseJSON[model.Property_manager](r)
	property_manager.Save()
}

func UpdateProperty_manager(w http.ResponseWriter, r *http.Request) {
	o := ParseJSON[model.Property_manager](r)
	if id := ParseID(r); id != 0 {
		o.Id = id
		o.Save()
		fmt.Fprint(w, `{"status": "OK", "msg": "Property_manager updated"}`)
		return
	} else {    
		if o.Email == ""  {
			fmt.Fprint(w, `{"status": "ERROR", "msg": "Property_manager no key value not provided"}`)
			return
		} else {
			o.Save()
			fmt.Fprint(w, `{"status": "OK", "msg": "Property_manager updated"}`)
			return
		}
	}
}

func DeleteProperty_manager(w http.ResponseWriter, r *http.Request) {
	if id := ParseID(r); id != 0 {
		model.DeleteProperty_managerByID(id)
		fmt.Fprint(w, `{"status": "OK", "msg": "Property_manager deleted"}`)
		return
	} else {
		property_manager := ParseJSON[model.Property_manager](r)
		if property_manager.Id != 0 {
			property_manager.Delete()
			fmt.Fprint(w, `{"status": "OK", "msg": "Property_manager deleted"}`)
		}
	}
}

func GetProperty_manager(w http.ResponseWriter, r *http.Request) {
	if id := ParseID(r); id == 0 {
		if where := r.PathValue("where"); where != "" {
			o := model.Property_manager{Where: where}
			os := o.Search()
			fmt.Fprint(w, u.JsonDump(os, ""))
			return
		} else {
			os := []model.Property_manager{}
			model.DB.Select(&os, `SELECT * FROM property_manager ORDER BY id LIMIT 200`)
			fmt.Fprint(w, u.JsonDump(os, ""))
			return
		}
	} else {
		o := model.GetProperty_managerByID(id)
		fmt.Fprint(w, u.JsonDump(*o, ""))
		return
	}
}

func Property_managerRouteRegister(mux *http.ServeMux, PathBase string) {
	mux.HandleFunc("POST "+PathBase+"/property_manager", CreateProperty_manager)
	mux.HandleFunc("PUT "+PathBase+"/property_manager/{id}{$}", UpdateProperty_manager)
	mux.HandleFunc("PUT "+PathBase+"/property_manager", UpdateProperty_manager)
	mux.HandleFunc("DELETE "+PathBase+"/property_manager/{id}{$}", DeleteProperty_manager)
	mux.HandleFunc("DELETE "+PathBase+"/property_manager", DeleteProperty_manager)
	mux.HandleFunc("GET "+PathBase+"/property_manager", GetProperty_manager)
	mux.HandleFunc("GET "+PathBase+"/property_manager/{id}{$}", GetProperty_manager)
	mux.HandleFunc("GET "+PathBase+"/property_manager/{where}{$}", GetProperty_manager)
}
