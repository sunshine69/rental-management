package app

import (
	"bytes"
	"fmt"
	"html/template"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"reflect"
	"strings"

	rice "github.com/GeertJohan/go.rice"
	"github.com/go-playground/form/v4"
	"github.com/go-playground/validator/v10"
	ag "github.com/sunshine69/automation-go/lib"
	u "github.com/sunshine69/golang-tools/utils"
	"github.com/sunshine69/rental-management/configs"
	"github.com/sunshine69/rental-management/model"
)

var (
	validate *validator.Validate
	// use a single instance of Decoder, it caches struct info
	formDecoder *form.Decoder
	AllTemplate *template.Template
)

// Auto generate using app-validation.go.tmpl template
var (
	FormTenant              model.Tenant
	FormProperty            model.Property
	FormAccount             model.Account
	FormContract            model.Contract
	FormPayment             model.Payment
	FormMaintenance_request model.Maintenance_request
	FormProperty_manager    model.Property_manager
	FormInvoice             model.Invoice
)

func TenantStructLevelValidation(sl validator.StructLevel) {
	// Change it to suit
	/* form := sl.Current().Interface().(model.Tenant)
	  if FormTenant.XXX != "XXX" {
		sl.ReportError(form.XXX, "XXX", "XXX", "XXX_Should_Not_Be_Set_When_XXX", "")
	} */
}
func PropertyStructLevelValidation(sl validator.StructLevel) {
	// Change it to suit
	/* form := sl.Current().Interface().(model.Property)
	  if FormProperty.XXX != "XXX" {
		sl.ReportError(form.XXX, "XXX", "XXX", "XXX_Should_Not_Be_Set_When_XXX", "")
	} */
}
func AccountStructLevelValidation(sl validator.StructLevel) {
	// Change it to suit
	/* form := sl.Current().Interface().(model.Account)
	  if FormAccount.XXX != "XXX" {
		sl.ReportError(form.XXX, "XXX", "XXX", "XXX_Should_Not_Be_Set_When_XXX", "")
	} */
}
func ContractStructLevelValidation(sl validator.StructLevel) {
	// Change it to suit
	/* form := sl.Current().Interface().(model.Contract)
	  if FormContract.XXX != "XXX" {
		sl.ReportError(form.XXX, "XXX", "XXX", "XXX_Should_Not_Be_Set_When_XXX", "")
	} */
}
func PaymentStructLevelValidation(sl validator.StructLevel) {
	// Change it to suit
	/* form := sl.Current().Interface().(model.Payment)
	  if FormPayment.XXX != "XXX" {
		sl.ReportError(form.XXX, "XXX", "XXX", "XXX_Should_Not_Be_Set_When_XXX", "")
	} */
}
func Maintenance_requestStructLevelValidation(sl validator.StructLevel) {
	// Change it to suit
	/* form := sl.Current().Interface().(model.Maintenance_request)
	  if FormMaintenance_request.XXX != "XXX" {
		sl.ReportError(form.XXX, "XXX", "XXX", "XXX_Should_Not_Be_Set_When_XXX", "")
	} */
}
func Property_managerStructLevelValidation(sl validator.StructLevel) {
	// Change it to suit
	/* form := sl.Current().Interface().(model.Property_manager)
	  if FormProperty_manager.XXX != "XXX" {
		sl.ReportError(form.XXX, "XXX", "XXX", "XXX_Should_Not_Be_Set_When_XXX", "")
	} */
}

// Common action for processing all forms.
func ProcessPreSteps(w http.ResponseWriter, r *http.Request, currentFormType any) {
	fmt.Fprintf(os.Stderr, "[DEBUG] Form of Type '%s'\n", reflect.TypeOf(currentFormType).Name())
	u.CheckErr(r.ParseForm(), "[ERROR] can not parse form")
	var err, err1 error
	// reflect.TypeOf will return a string representing the struct name, such as 'Form1'. Need to pass not using &
	switch reflect.TypeOf(currentFormType).Name() {
	case "Tenant":
		// Now we do type assertion based on string return by reflect to cast it to the original type (from any). This is needed as formDecoder
		// and validator needs the exact type to bind html for to struct because it needs to see the struct field and tags to collect html form data into it

		err = u.CheckErrNonFatal(formDecoder.Decode(&FormTenant, r.PostForm), "formDecoder.Decode")
		err1 = validate.Struct(FormTenant)

		if err != nil {
			fmt.Fprintf(os.Stderr, "[ERROR] Form binding '%s'\n", err.Error())
			AllTemplate.ExecuteTemplate(w, "error.html", map[string]any{"output": "ERROR form binding please check server log"})
			return
		}
		if err1 != nil { // Check validation so any errors will come out here. Currently it would just display the validation tag string
			fmt.Fprintf(os.Stderr, "[ERROR] Form validation '%s'\n", err1.Error())
			AllTemplate.ExecuteTemplate(w, "error.html", map[string]any{
				"err":    strings.ReplaceAll(err1.Error(), "\n", "<br/>"),
				"output": "",
				"action": fmt.Sprintf(`<p>Error</p>
					<a href="%s">click here to try again</a><br>
					Exit button to quit the program
					<div class="button-group">
						<form action="/quit" method="post">
						<input type="submit" name="submit" value="Exit">
						</form>
					</div>`, "/home"),
			})
			return
		}
	case "Property":
		// Now we do type assertion based on string return by reflect to cast it to the original type (from any). This is needed as formDecoder
		// and validator needs the exact type to bind html for to struct because it needs to see the struct field and tags to collect html form data into it

		err = u.CheckErrNonFatal(formDecoder.Decode(&FormProperty, r.PostForm), "formDecoder.Decode")
		err1 = validate.Struct(FormProperty)

		if err != nil {
			fmt.Fprintf(os.Stderr, "[ERROR] Form binding '%s'\n", err.Error())
			AllTemplate.ExecuteTemplate(w, "error.html", map[string]any{"output": "ERROR form binding please check server log"})
			return
		}
		if err1 != nil { // Check validation so any errors will come out here. Currently it would just display the validation tag string
			fmt.Fprintf(os.Stderr, "[ERROR] Form validation '%s'\n", err1.Error())
			AllTemplate.ExecuteTemplate(w, "error.html", map[string]any{
				"err":    strings.ReplaceAll(err1.Error(), "\n", "<br/>"),
				"output": "",
				"action": fmt.Sprintf(`<p>Error</p>
					<a href="%s">click here to try again</a><br>
					Exit button to quit the program
					<div class="button-group">
						<form action="/quit" method="post">
						<input type="submit" name="submit" value="Exit">
						</form>
					</div>`, "/home"),
			})
			return
		}
	case "Account":
		// Now we do type assertion based on string return by reflect to cast it to the original type (from any). This is needed as formDecoder
		// and validator needs the exact type to bind html for to struct because it needs to see the struct field and tags to collect html form data into it

		err = u.CheckErrNonFatal(formDecoder.Decode(&FormAccount, r.PostForm), "formDecoder.Decode")
		err1 = validate.Struct(FormAccount)

		if err != nil {
			fmt.Fprintf(os.Stderr, "[ERROR] Form binding '%s'\n", err.Error())
			AllTemplate.ExecuteTemplate(w, "error.html", map[string]any{"output": "ERROR form binding please check server log"})
			return
		}
		if err1 != nil { // Check validation so any errors will come out here. Currently it would just display the validation tag string
			fmt.Fprintf(os.Stderr, "[ERROR] Form validation '%s'\n", err1.Error())
			AllTemplate.ExecuteTemplate(w, "error.html", map[string]any{
				"err":    strings.ReplaceAll(err1.Error(), "\n", "<br/>"),
				"output": "",
				"action": fmt.Sprintf(`<p>Error</p>
					<a href="%s">click here to try again</a><br>
					Exit button to quit the program
					<div class="button-group">
						<form action="/quit" method="post">
						<input type="submit" name="submit" value="Exit">
						</form>
					</div>`, "/home"),
			})
			return
		}
	case "Contract":
		// Now we do type assertion based on string return by reflect to cast it to the original type (from any). This is needed as formDecoder
		// and validator needs the exact type to bind html for to struct because it needs to see the struct field and tags to collect html form data into it

		err = u.CheckErrNonFatal(formDecoder.Decode(&FormContract, r.PostForm), "formDecoder.Decode")
		err1 = validate.Struct(FormContract)

		if err != nil {
			fmt.Fprintf(os.Stderr, "[ERROR] Form binding '%s'\n", err.Error())
			AllTemplate.ExecuteTemplate(w, "error.html", map[string]any{"output": "ERROR form binding please check server log"})
			return
		}
		if err1 != nil { // Check validation so any errors will come out here. Currently it would just display the validation tag string
			fmt.Fprintf(os.Stderr, "[ERROR] Form validation '%s'\n", err1.Error())
			AllTemplate.ExecuteTemplate(w, "error.html", map[string]any{
				"err":    strings.ReplaceAll(err1.Error(), "\n", "<br/>"),
				"output": "",
				"action": fmt.Sprintf(`<p>Error</p>
					<a href="%s">click here to try again</a><br>
					Exit button to quit the program
					<div class="button-group">
						<form action="/quit" method="post">
						<input type="submit" name="submit" value="Exit">
						</form>
					</div>`, "/home"),
			})
			return
		}
	case "Payment":
		// Now we do type assertion based on string return by reflect to cast it to the original type (from any). This is needed as formDecoder
		// and validator needs the exact type to bind html for to struct because it needs to see the struct field and tags to collect html form data into it

		err = u.CheckErrNonFatal(formDecoder.Decode(&FormPayment, r.PostForm), "formDecoder.Decode")
		err1 = validate.Struct(FormPayment)

		if err != nil {
			fmt.Fprintf(os.Stderr, "[ERROR] Form binding '%s'\n", err.Error())
			AllTemplate.ExecuteTemplate(w, "error.html", map[string]any{"output": "ERROR form binding please check server log"})
			return
		}
		if err1 != nil { // Check validation so any errors will come out here. Currently it would just display the validation tag string
			fmt.Fprintf(os.Stderr, "[ERROR] Form validation '%s'\n", err1.Error())
			AllTemplate.ExecuteTemplate(w, "error.html", map[string]any{
				"err":    strings.ReplaceAll(err1.Error(), "\n", "<br/>"),
				"output": "",
				"action": fmt.Sprintf(`<p>Error</p>
					<a href="%s">click here to try again</a><br>
					Exit button to quit the program
					<div class="button-group">
						<form action="/quit" method="post">
						<input type="submit" name="submit" value="Exit">
						</form>
					</div>`, "/home"),
			})
			return
		}
	case "Maintenance_request":
		// Now we do type assertion based on string return by reflect to cast it to the original type (from any). This is needed as formDecoder
		// and validator needs the exact type to bind html for to struct because it needs to see the struct field and tags to collect html form data into it

		err = u.CheckErrNonFatal(formDecoder.Decode(&FormMaintenance_request, r.PostForm), "formDecoder.Decode")
		err1 = validate.Struct(FormMaintenance_request)

		if err != nil {
			fmt.Fprintf(os.Stderr, "[ERROR] Form binding '%s'\n", err.Error())
			AllTemplate.ExecuteTemplate(w, "error.html", map[string]any{"output": "ERROR form binding please check server log"})
			return
		}
		if err1 != nil { // Check validation so any errors will come out here. Currently it would just display the validation tag string
			fmt.Fprintf(os.Stderr, "[ERROR] Form validation '%s'\n", err1.Error())
			AllTemplate.ExecuteTemplate(w, "error.html", map[string]any{
				"err":    strings.ReplaceAll(err1.Error(), "\n", "<br/>"),
				"output": "",
				"action": fmt.Sprintf(`<p>Error</p>
					<a href="%s">click here to try again</a><br>
					Exit button to quit the program
					<div class="button-group">
						<form action="/quit" method="post">
						<input type="submit" name="submit" value="Exit">
						</form>
					</div>`, "/home"),
			})
			return
		}
	case "Property_manager":
		// Now we do type assertion based on string return by reflect to cast it to the original type (from any). This is needed as formDecoder
		// and validator needs the exact type to bind html for to struct because it needs to see the struct field and tags to collect html form data into it

		err = u.CheckErrNonFatal(formDecoder.Decode(&FormProperty_manager, r.PostForm), "formDecoder.Decode")
		err1 = validate.Struct(FormProperty_manager)

		if err != nil {
			fmt.Fprintf(os.Stderr, "[ERROR] Form binding '%s'\n", err.Error())
			AllTemplate.ExecuteTemplate(w, "error.html", map[string]any{"output": "ERROR form binding please check server log"})
			return
		}
		if err1 != nil { // Check validation so any errors will come out here. Currently it would just display the validation tag string
			fmt.Fprintf(os.Stderr, "[ERROR] Form validation '%s'\n", err1.Error())
			AllTemplate.ExecuteTemplate(w, "error.html", map[string]any{
				"err":    strings.ReplaceAll(err1.Error(), "\n", "<br/>"),
				"output": "",
				"action": fmt.Sprintf(`<p>Error</p>
					<a href="%s">click here to try again</a><br>
					Exit button to quit the program
					<div class="button-group">
						<form action="/quit" method="post">
						<input type="submit" name="submit" value="Exit">
						</form>
					</div>`, "/home"),
			})
			return
		}
	}
}

func init() {
	validate = validator.New(validator.WithRequiredStructEnabled())
	formDecoder = form.NewDecoder()

	// Register validation
	validate.RegisterStructValidation(TenantStructLevelValidation, FormTenant)
	validate.RegisterStructValidation(PropertyStructLevelValidation, FormProperty)
	validate.RegisterStructValidation(AccountStructLevelValidation, FormAccount)
	validate.RegisterStructValidation(ContractStructLevelValidation, FormContract)
	validate.RegisterStructValidation(PaymentStructLevelValidation, FormPayment)
	validate.RegisterStructValidation(Maintenance_requestStructLevelValidation, FormMaintenance_request)
	validate.RegisterStructValidation(Property_managerStructLevelValidation, FormProperty_manager)
}

// End app-validation.go.tmpl

// Auto generate using app-handler.go.tmpl template

func Tenant(w http.ResponseWriter, r *http.Request) {
    ProcessPreSteps(w, r, model.Tenant{})
    if err := FormTenant.Save(); err != nil {
        fmt.Fprint(w, "[ERROR] while saving object. See the server log for details")
        fmt.Fprintf(os.Stderr, "%s", err.Error())
        return
    }
    fmt.Fprint(w, "Data saved")
}
func Property(w http.ResponseWriter, r *http.Request) {
    ProcessPreSteps(w, r, model.Property{})
    if err := FormProperty.Save(); err != nil {
        fmt.Fprint(w, "[ERROR] while saving object. See the server log for details")
        fmt.Fprintf(os.Stderr, "%s", err.Error())
        return
    }
    fmt.Fprint(w, "Data saved")
}
func Account(w http.ResponseWriter, r *http.Request) {
    ProcessPreSteps(w, r, model.Account{})
    if err := FormAccount.Save(); err != nil {
        fmt.Fprint(w, "[ERROR] while saving object. See the server log for details")
        fmt.Fprintf(os.Stderr, "%s", err.Error())
        return
    }
    fmt.Fprint(w, "Data saved")
}
func Contract(w http.ResponseWriter, r *http.Request) {
    ProcessPreSteps(w, r, model.Contract{})
    if err := FormContract.Save(); err != nil {
        fmt.Fprint(w, "[ERROR] while saving object. See the server log for details")
        fmt.Fprintf(os.Stderr, "%s", err.Error())
        return
    }
    fmt.Fprint(w, "Data saved")
}
func Payment(w http.ResponseWriter, r *http.Request) {
    ProcessPreSteps(w, r, model.Payment{})
    if err := FormPayment.Save(); err != nil {
        fmt.Fprint(w, "[ERROR] while saving object. See the server log for details")
        fmt.Fprintf(os.Stderr, "%s", err.Error())
        return
    }
    fmt.Fprint(w, "Data saved")
}
func Maintenance_request(w http.ResponseWriter, r *http.Request) {
    ProcessPreSteps(w, r, model.Maintenance_request{})
    if err := FormMaintenance_request.Save(); err != nil {
        fmt.Fprint(w, "[ERROR] while saving object. See the server log for details")
        fmt.Fprintf(os.Stderr, "%s", err.Error())
        return
    }
    fmt.Fprint(w, "Data saved")
}
func Property_manager(w http.ResponseWriter, r *http.Request) {
    ProcessPreSteps(w, r, model.Property_manager{})
    if err := FormProperty_manager.Save(); err != nil {
        fmt.Fprint(w, "[ERROR] while saving object. See the server log for details")
        fmt.Fprintf(os.Stderr, "%s", err.Error())
        return
    }
    fmt.Fprint(w, "Data saved")
}
func Invoice(w http.ResponseWriter, r *http.Request) {
    ProcessPreSteps(w, r, model.Invoice{})
    if err := FormInvoice.Save(); err != nil {
        fmt.Fprint(w, "[ERROR] while saving object. See the server log for details")
        fmt.Fprintf(os.Stderr, "%s", err.Error())
        return
    }
    fmt.Fprint(w, "Data saved")
}

func AddHandler(mux *http.ServeMux, Cfg *configs.Config) {
    mux.HandleFunc("POST "+Cfg.PathBase+"/tenant", Tenant)
    mux.HandleFunc("POST "+Cfg.PathBase+"/property", Property)
    mux.HandleFunc("POST "+Cfg.PathBase+"/account", Account)
    mux.HandleFunc("POST "+Cfg.PathBase+"/contract", Contract)
    mux.HandleFunc("POST "+Cfg.PathBase+"/payment", Payment)
    mux.HandleFunc("POST "+Cfg.PathBase+"/maintenance_request", Maintenance_request)
    mux.HandleFunc("POST "+Cfg.PathBase+"/property_manager", Property_manager)
    mux.HandleFunc("POST "+Cfg.PathBase+"/invoice", Invoice)
}
// End app-handler.go.tmpl

func loadAllTemplates() *template.Template {
	t := template.New("tmpl")
	myFuncmap := ag.GoTemplateFuncMap
	myFuncmap["CallTemplate"] = func(name string, data interface{}) (ret template.HTML, err error) {
		buf := bytes.NewBuffer([]byte{})
		// Need to use t (already created above) as when we call t parse ; will have the template name available
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
	AddHandler(mux, Cfg)
}
