package app

import (
	"bytes"
	"fmt"
	"html/template"
	"io/fs"
	"net/http"
	"path/filepath"

	rice "github.com/GeertJohan/go.rice"
	"github.com/go-playground/form/v4"
	"github.com/go-playground/validator/v10"
	ag "github.com/sunshine69/automation-go/lib"
	u "github.com/sunshine69/golang-tools/utils"
	"github.com/sunshine69/rental-management/configs"
)

var (
	validate *validator.Validate
	// use a single instance of Decoder, it caches struct info
	formDecoder *form.Decoder
	AllTemplate *template.Template
)

func loadAllTemplates() *template.Template {
	t := template.New("tmpl")
	myFuncmap := ag.GoTemplateFuncMap
	myFuncmap["CallTemplate"] = func(name string, data interface{}) (ret template.HTML, err error) {
		buf := bytes.NewBuffer([]byte{})
		// fmt.Printf("[DEBUG] templating %s with data %v\n", name, data)
		err = t.ExecuteTemplate(buf, "Form"+name, data)
		u.CheckErr(err, "ERR")
		ret = template.HTML(buf.String())
		return
	}
	t.Funcs(myFuncmap)
	templatesBox := rice.MustFindBox("templates")
	templatesBox.Walk("/", func(path string, info fs.FileInfo, err error) error {
		fmt.Println(path)
		if info.IsDir() {
			return nil
		}
		fname := filepath.Base(path)
		t = template.Must(t.New(fname).Parse(templatesBox.MustString(fname)))
		return nil
	})
	return t
}

func Home(w http.ResponseWriter, r *http.Request) {
	formList := []string{"Tenant",
		"Property",
		"Account",
		"Payment",
		"Contract",
		"Invoice",
		"Maintenance_request",
		"Property_manager"}
	AllTemplate.ExecuteTemplate(w, "index.html", map[string]any{"formList": formList})
	// AllTemplate.ExecuteTemplate(w, "index.html", nil)
}

// Startup the app
func StartWebApp(mux *http.ServeMux, Cfg *configs.Config) {

	AllTemplate = loadAllTemplates()
	// Web app part.
	assetBox := rice.MustFindBox("assets")
	mux.Handle(Cfg.PathBase+"/static/", http.StripPrefix(Cfg.PathBase+"/static/", http.FileServer(assetBox.HTTPBox())))
	mux.HandleFunc("GET "+Cfg.PathBase+"/home", Home)
	mux.HandleFunc("POST "+Cfg.PathBase+"/tenant", Home)
}
