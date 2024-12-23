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

func CreateProperty(w http.ResponseWriter, r *http.Request) {
	if property := u.ParseJsonReqBodyToStruct[model.Property](r); property != nil {
		if err := property.Save(); err != nil {
			fmt.Fprintf(w, `{"status": "ERROR", "msg": "%s"}`, err.Error())
		}
	} else {
		fmt.Fprint(w, `{"status": "ERROR", "msg": "Can not parse json into property"}`)
	}
}

func UpdateProperty(w http.ResponseWriter, r *http.Request) {
	o := u.ParseJsonReqBodyToMap(r)
	if id := utils.ParseID(r); id != 0 {
		tn := model.GetPropertyByID(id)
		if tn == nil {
			fmt.Fprint(w, `{"status": "ERROR", "msg": "Property not found with this id"}`)
			return
		}
		o["id"] = id
		if err := tn.Update(o); err == nil {
			fmt.Fprint(w, `{"status": "OK", "msg": "Property updated"}`)
		} else {
			fmt.Fprint(w, `{"status": "ERROR", "msg": "`+err.Error()+`"}`)
		}
		return
	} else {
		inputSet := sets.FromMap(o)
		keySet := sets.New("code")
		if !inputSet.Superset(keySet) {
			fmt.Fprint(w, `{"status": "ERROR", "msg": "Property no key value not provided"}`)
			return
		}
		tn := model.GetPropertyByCompositeKeyOrNew(o)
		if err := tn.Update(o); err == nil {
			fmt.Fprint(w, `{"status": "OK", "msg": "Property updated"}`)
		} else {
			fmt.Fprint(w, `{"status": "ERROR", "msg": "`+err.Error()+`"}`)
		}
		return
	}
}

func DeleteProperty(w http.ResponseWriter, r *http.Request) {
	if id := utils.ParseID(r); id != 0 {
		if err := model.DeletePropertyByID(id); err == nil {
			fmt.Fprint(w, `{"status": "OK", "msg": "Property deleted"}`)
		} else {
			fmt.Fprint(w, `{"status": "ERROR", "msg": `+err.Error())
		}
		return
	} else {
		property := u.ParseJsonReqBodyToStruct[model.Property](r)
		if property != nil {
			if property.Id != 0 {
				model.DeletePropertyByID(property.Id)
				fmt.Fprint(w, `{"status": "OK", "msg": "Property deleted"}`)
				return
			} else if property.Code != "" {
				property.Delete()
				fmt.Fprint(w, `{"status": "OK", "msg": "Property deleted"}`)
				return
			} else {
				fmt.Fprint(w, `{"status": "ERROR", "msg": "No composite key found"}`)
				return
			}
		} else {
			fmt.Fprint(w, `{"status": "ERROR", "msg": "Cannot parse property"}`)
			return
		}
	}
}

func GetProperty(w http.ResponseWriter, r *http.Request) {
	if id := utils.ParseID(r); id == 0 {
		if where := r.PathValue("where"); where != "" {
			o := model.Property{Where: where}
			os := o.Search()
			fmt.Fprint(w, u.JsonDump(os, ""))
			return
		} else {
			os := []model.Property{}
			model.DB.Select(&os, `SELECT * FROM property ORDER BY id`)
			fmt.Fprint(w, u.JsonDump(os, ""))
			return
		}
	} else {
		o := model.GetPropertyByID(id)
		fmt.Fprint(w, u.JsonDump(*o, ""))
		return
	}
}

func RouteRegisterProperty(mux *http.ServeMux, PathBase string) {
	mux.HandleFunc("POST "+PathBase+"/api/property", CreateProperty)
	mux.HandleFunc("PUT "+PathBase+"/api/property/{id}", UpdateProperty)
	mux.HandleFunc("PUT "+PathBase+"/api/property", UpdateProperty)
	mux.HandleFunc("DELETE "+PathBase+"/api/property/{id}", DeleteProperty)
	mux.HandleFunc("DELETE "+PathBase+"/api/property", DeleteProperty)
	mux.HandleFunc("GET "+PathBase+"/api/property", GetProperty)
	mux.HandleFunc("GET "+PathBase+"/api/property/{id}", GetProperty)
	mux.HandleFunc("GET "+PathBase+"/api/property/q/{where}", GetProperty)
}
