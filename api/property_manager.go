// Generated by api-gen tool. Do not edit here but edit in api/api-template.go.tmpl
package api

import (
	"fmt"
	"net/http"

	"github.com/R167/go-sets"
	u "github.com/sunshine69/golang-tools/utils"
	"github.com/sunshine69/rental-management/model"
)

func CreateProperty_manager(w http.ResponseWriter, r *http.Request) {
	if property_manager := ParseJSON[model.Property_manager](r); property_manager != nil {
		if err := property_manager.Save(); err != nil {
			fmt.Fprintf(w, `{"status": "ERROR", "msg": "%s"}`, err.Error())
		}
	} else {
		fmt.Fprint(w, `{"status": "ERROR", "msg": "Can not parse json into property_manager"}`)
	}
}

func UpdateProperty_manager(w http.ResponseWriter, r *http.Request) {
	o := ParseJSONToMap(r)
	if id := ParseID(r); id != 0 {
		tn := model.GetProperty_managerByID(id)
		if tn == nil {
			fmt.Fprint(w, `{"status": "ERROR", "msg": "Property_manager not found with this id"}`)
			return
		}
		o["id"] = id
		if err := tn.Update(o); err == nil {
			fmt.Fprint(w, `{"status": "OK", "msg": "Property_manager updated"}`)
		} else {
			fmt.Fprint(w, `{"status": "ERROR", "msg": "`+err.Error()+`"}`)
		}
		return
	} else {
		inputSet := sets.FromMap(o)
		keySet := sets.New("email")
		if !inputSet.Superset(keySet) {
			fmt.Fprint(w, `{"status": "ERROR", "msg": "Property_manager no key value not provided"}`)
			return
		}
		tn := model.GetProperty_managerByCompositeKeyOrNew(o)
		if err := tn.Update(o); err == nil {
			fmt.Fprint(w, `{"status": "OK", "msg": "Property_manager updated"}`)
		} else {
			fmt.Fprint(w, `{"status": "ERROR", "msg": "`+err.Error()+`"}`)
		}
		return
	}
}

func DeleteProperty_manager(w http.ResponseWriter, r *http.Request) {
	if id := ParseID(r); id != 0 {
		if err := model.DeleteProperty_managerByID(id); err == nil {
			fmt.Fprint(w, `{"status": "OK", "msg": "Property_manager deleted"}`)
		} else {
			fmt.Fprint(w, `{"status": "ERROR", "msg": `+err.Error())
		}
		return
	} else {
		property_manager := ParseJSON[model.Property_manager](r)
		if property_manager != nil {
			if property_manager.Id != 0 {
				model.DeleteProperty_managerByID(property_manager.Id)
				fmt.Fprint(w, `{"status": "OK", "msg": "Property_manager deleted"}`)
				return
			} else if property_manager.Email != "" {
				property_manager.Delete()
				fmt.Fprint(w, `{"status": "OK", "msg": "Property_manager deleted"}`)
				return
			} else {
				fmt.Fprint(w, `{"status": "ERROR", "msg": "No composite key found"}`)
				return
			}
		} else {
			fmt.Fprint(w, `{"status": "ERROR", "msg": "Cannot parse property_manager"}`)
			return
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
			model.DB.Select(&os, `SELECT * FROM property_manager ORDER BY id`)
			fmt.Fprint(w, u.JsonDump(os, ""))
			return
		}
	} else {
		o := model.GetProperty_managerByID(id)
		fmt.Fprint(w, u.JsonDump(*o, ""))
		return
	}
}

func RouteRegisterProperty_manager(mux *http.ServeMux, PathBase string) {
	mux.HandleFunc("POST "+PathBase+"/api/property_manager", CreateProperty_manager)
	mux.HandleFunc("PUT "+PathBase+"/api/property_manager/{id}", UpdateProperty_manager)
	mux.HandleFunc("PUT "+PathBase+"/api/property_manager", UpdateProperty_manager)
	mux.HandleFunc("DELETE "+PathBase+"/api/property_manager/{id}", DeleteProperty_manager)
	mux.HandleFunc("DELETE "+PathBase+"/api/property_manager", DeleteProperty_manager)
	mux.HandleFunc("GET "+PathBase+"/api/property_manager", GetProperty_manager)
	mux.HandleFunc("GET "+PathBase+"/api/property_manager/{id}", GetProperty_manager)
	mux.HandleFunc("GET "+PathBase+"/api/property_manager/q/{where}", GetProperty_manager)
}
