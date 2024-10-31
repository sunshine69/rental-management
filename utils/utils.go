package utils

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
)

func ParseID(r *http.Request) int64 {
	idstr := r.PathValue("id")
	if idstr == "" {
		idstr = r.FormValue("id")
	}
	if idstr != "" {
		if id, err := strconv.ParseInt(idstr, 10, 64); err != nil {
			fmt.Fprintln(os.Stderr, `{"status": "ERROR", "msg": "[ERROR] parsing id"}`)
			return 0
		} else {
			return id
		}
	} else {
		fmt.Fprintln(os.Stderr, `{"status": "ERROR", "msg": "[ERROR] no id supplied"}`)
		return 0
	}
}
