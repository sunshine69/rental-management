package main

import (
	"embed"
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/go-playground/form/v4"
	"github.com/go-playground/validator/v10"
	ag "github.com/sunshine69/automation-go/lib"
	u "github.com/sunshine69/golang-tools/utils"
)

var (
	TemplateFuncMap *template.FuncMap
	ListenPort      string
	validate        *validator.Validate
	// use a single instance of Decoder, it caches struct info
	formDecoder *form.Decoder
)

//go:embed assets
var staticFolder embed.FS

//go:embed templates
var tplFolder embed.FS // embeds the templates folder into variable tplFolder

func loadAllTemplates() *template.Template {
	t, err := template.New("templ").Funcs(ag.GoTemplateFuncMap).ParseFS(tplFolder, "templates/*.html")
	if err != nil {
		log.Fatalf("can not parse templates %s", err.Error())
	}
	return t
}

// As the content of vars-ansible.yaml is large creating one html form to hold data is big thus split them into several forms and collecting data
// the sequence steps. Some of the fields values depending the value of the previous field(s) thus we will write custom validation rules
type Form1 struct {
	ProjectName                  string `form:"project_name" validate:"required,alphanumunicode"`
	ProjectType                  string `form:"project_type"`
	Namespace_prefix             string `form:"namespace_prefix" validate:"required"`
	Tfs_project_name             string `form:"tfs_project_name" validate:"required"`
	Tfs_repo_name                string `form:"tfs_repo_name" validate:"required"`
	Build_operation_disabled     string `form:"build_operation_disabled"`
	Docker_registry_project_name string `form:"docker_registry_project_name" validate:"required"`
}

var (
	form1 Form1
)

func Home(w http.ResponseWriter, r *http.Request) {

}

// Step1 - show the GUI
func Step1(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":

	case "GET":
		AllTemplate.ExecuteTemplate(w, "step1.html", map[string]interface{}{})
	}
}

var AllTemplate *template.Template

func main() {
	flag.StringVar(&ListenPort, "p", "8080", "Local port to listen")
	flag.Parse()

	AllTemplate = loadAllTemplates()
	validate = validator.New(validator.WithRequiredStructEnabled())
	formDecoder = form.NewDecoder()

	fmt.Printf("[DEBUG] form initialization %s\n", u.JsonDump(form1, "  "))

	r := http.NewServeMux()

	r.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.FS(staticFolder))))

	// r.HandleFunc("/", Home)
	r.HandleFunc("/step1", Step1)

	wait_chan := make(chan int)
	r.HandleFunc("POST /quit", func(w http.ResponseWriter, r *http.Request) { wait_chan <- 1 })

	go http.ListenAndServe(":"+ListenPort, r)

	// go u.RunSystemCommandV2(fmt.Sprintf("%s http://localhost:%s/", "google-chrome ", ListenPort), true)

	<-wait_chan
}
