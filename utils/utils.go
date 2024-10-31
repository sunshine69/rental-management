package utils

import (
	"bytes"
	"fmt"
	json "github.com/json-iterator/go"
	ag "github.com/sunshine69/automation-go/lib"
	"io"
	"net/http"
	"os"
	"strconv"
)

func AssertInt64ValueForMap(input map[string]interface{}) map[string]interface{} {
	for k, v := range input {
		if v, ok := v.(float64); ok {
			input[k] = int64(v)

		}
	}
	return input
}

func JsonToMap(jsonStr string) map[string]interface{} {
	result := make(map[string]interface{})
	json.Unmarshal([]byte(jsonStr), &result)
	return AssertInt64ValueForMap(result)
}

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

func ParseJSONToMap(r *http.Request) map[string]interface{} {
	switch r.Method {
	case "POST", "PUT", "DELETE":
		jsonBytes := bytes.Buffer{}
		if _, err := io.Copy(&jsonBytes, r.Body); err != nil {
			fmt.Fprintf(os.Stderr, "[ERROR] ParseJSONToMap loading request body - %s\n", err.Error())
		}
		defer r.Body.Close()
		return JsonToMap(string(jsonBytes.Bytes()))
	default:
		fmt.Fprintf(os.Stderr, "[ERROR] ParseJSONToMap Do not call me with this method - %s\n", r.Method)
		return nil
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

// Convert any model type to a map - useful to render it in template
func Convert2Map[T any](t T) ([]string, map[string]any) {
	sInfo := ag.ReflectStruct(t, "")
	out := map[string]any{}
	for _, f := range sInfo.FieldName {
		out[f] = sInfo.FieldValue[f]
	}
	return sInfo.FieldName, out
}
