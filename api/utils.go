package api

import (
	"fmt"
	json "github.com/json-iterator/go"
	"net/http"
	"os"
	"strconv"

	"github.com/sunshine69/rental-management/model"
)

func ParseID(r *http.Request) int64 {
	idstr := r.PathValue("id")
	if idstr == "" {
		idstr = r.FormValue("id")
	}
	if idstr != "" {
		if id, err := strconv.ParseInt(idstr, 10, 64); err != nil {
			fmt.Fprintf(os.Stderr, `{"status": "ERROR", "msg": "[ERROR] parsing id"}`)
			return 0
		} else {
			return id
		}
	} else {
		fmt.Fprintf(os.Stderr, `{"status": "ERROR", "msg": "[ERROR] no id supplied"}`)
		return 0
	}
}

// ParseJSON parses the raw JSON body from an HTTP request into the specified struct.
func ParseJSON[T any](r *http.Request) *T {
	switch r.Method {
	case "POST", "PUT", "DELETE":
		decoder := json.NewDecoder(r.Body)
		defer r.Body.Close()
		var data T

		if err := decoder.Decode(&data); err != nil {
			fmt.Fprintf(os.Stderr, "[ERROR] parsing JSON - %s\n", err.Error())
			return nil
		}
		return &data
	default:
		fmt.Fprintf(os.Stderr, "[ERROR] ParseJSON Do not call me with this method - %s\n", r.Method)
		return nil
	}
}
