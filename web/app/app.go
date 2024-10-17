package app

import (
	"fmt"
	"html/template"
	"io/fs"
	"net/http"
	"path/filepath"

	rice "github.com/GeertJohan/go.rice"
	ag "github.com/sunshine69/automation-go/lib"
	"github.com/sunshine69/rental-management/configs"
)

var AllTemplate *template.Template

func loadAllTemplates() *template.Template {
	t := template.New("tmpl").Funcs(ag.GoTemplateFuncMap)
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
	// t, err := tloader.LoadTemplatesFromBinary(templatesBox, ag.GoTemplateFuncMap, true)
	// if err != nil {
	// 	log.Fatalf("can not parse templates %s", err.Error())
	// }
	return t
}

func Home(w http.ResponseWriter, r *http.Request) {
	AllTemplate.ExecuteTemplate(w, "index.html", nil)
}

func RouteRegisterApp(mux *http.ServeMux, Cfg *configs.Config) {
	AllTemplate = loadAllTemplates()
	// Web app part.
	assetBox := rice.MustFindBox("assets")
	mux.Handle(Cfg.PathBase+"/static/", http.StripPrefix(Cfg.PathBase+"/static/", http.FileServer(assetBox.HTTPBox())))
	mux.HandleFunc("GET "+Cfg.PathBase+"/home", Home)
}
