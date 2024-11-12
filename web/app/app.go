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
	"strconv"
	"strings"

	// To bundle assets first build the binary - then get into this dir (where the go file has the rice findbox command) and run 'rice append --exec <path-to-bin>
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

func TenantStructLevelValidation(sl validator.StructLevel) {
	// Change it to suit
	/* form := sl.Current().Interface().(model.Tenant)
	  if formTenant.XXX != "XXX" {
		sl.ReportError(form.XXX, "XXX", "XXX", "XXX_Should_Not_Be_Set_When_XXX", "")
	} */
}
func Property_managerStructLevelValidation(sl validator.StructLevel) {
	// Change it to suit
	/* form := sl.Current().Interface().(model.Property_manager)
	  if formProperty_manager.XXX != "XXX" {
		sl.ReportError(form.XXX, "XXX", "XXX", "XXX_Should_Not_Be_Set_When_XXX", "")
	} */
}
func PropertyStructLevelValidation(sl validator.StructLevel) {
	// Change it to suit
	/* form := sl.Current().Interface().(model.Property)
	  if formProperty.XXX != "XXX" {
		sl.ReportError(form.XXX, "XXX", "XXX", "XXX_Should_Not_Be_Set_When_XXX", "")
	} */
}
func ContractStructLevelValidation(sl validator.StructLevel) {
	// Change it to suit
	/* form := sl.Current().Interface().(model.Contract)
	  if formContract.XXX != "XXX" {
		sl.ReportError(form.XXX, "XXX", "XXX", "XXX_Should_Not_Be_Set_When_XXX", "")
	} */
}
func AccountStructLevelValidation(sl validator.StructLevel) {
	// Change it to suit
	/* form := sl.Current().Interface().(model.Account)
	  if formAccount.XXX != "XXX" {
		sl.ReportError(form.XXX, "XXX", "XXX", "XXX_Should_Not_Be_Set_When_XXX", "")
	} */
}
func PaymentStructLevelValidation(sl validator.StructLevel) {
	// Change it to suit
	/* form := sl.Current().Interface().(model.Payment)
	  if formPayment.XXX != "XXX" {
		sl.ReportError(form.XXX, "XXX", "XXX", "XXX_Should_Not_Be_Set_When_XXX", "")
	} */
}
func InvoiceStructLevelValidation(sl validator.StructLevel) {
	// Change it to suit
	/* form := sl.Current().Interface().(model.Invoice)
	  if formInvoice.XXX != "XXX" {
		sl.ReportError(form.XXX, "XXX", "XXX", "XXX_Should_Not_Be_Set_When_XXX", "")
	} */
}
func Maintenance_requestStructLevelValidation(sl validator.StructLevel) {
	// Change it to suit
	/* form := sl.Current().Interface().(model.Maintenance_request)
	  if formMaintenance_request.XXX != "XXX" {
		sl.ReportError(form.XXX, "XXX", "XXX", "XXX_Should_Not_Be_Set_When_XXX", "")
	} */
}

// Common action for processing all forms.
func ProcessPreSteps[T any](w http.ResponseWriter, r *http.Request, currentFormType T) (T, error) {
	fmt.Fprintf(os.Stderr, "[DEBUG] Form of Type '%s'\n", reflect.TypeOf(currentFormType).Name())
	u.CheckErr(r.ParseForm(), "[ERROR] can not parse form")
	var err, err1 error
	var newT T
	err = u.CheckErrNonFatal(formDecoder.Decode(&newT, r.PostForm), "formDecoder.Decode")
	err1 = validate.Struct(newT)
	if err != nil {
		fmt.Fprintf(os.Stderr, "[ERROR] Form binding '%s'\n", err.Error())
		AllTemplate.ExecuteTemplate(w, "error.html", map[string]any{"output": "ERROR form binding please check server log"})
		return newT, err
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
		return newT, err
	}
	return newT, nil
}

func init () {
	validate = validator.New(validator.WithRequiredStructEnabled())
	formDecoder = form.NewDecoder()

	// Register validation
	validate.RegisterStructValidation(TenantStructLevelValidation, model.Tenant{})
	validate.RegisterStructValidation(Property_managerStructLevelValidation, model.Property_manager{})
	validate.RegisterStructValidation(PropertyStructLevelValidation, model.Property{})
	validate.RegisterStructValidation(ContractStructLevelValidation, model.Contract{})
	validate.RegisterStructValidation(AccountStructLevelValidation, model.Account{})
	validate.RegisterStructValidation(PaymentStructLevelValidation, model.Payment{})
	validate.RegisterStructValidation(InvoiceStructLevelValidation, model.Invoice{})
	validate.RegisterStructValidation(Maintenance_requestStructLevelValidation, model.Maintenance_request{})
}
// End app-validation.go.tmpl

// Auto generate using app-handler.go.tmpl template

func Tenant(w http.ResponseWriter, r *http.Request) {
    obj, err := ProcessPreSteps(w, r, model.Tenant{})
    if err != nil {
        fmt.Fprintf(w, "%s", err.Error())
        return
    }
    if err := obj.Save(); err != nil {
        fmt.Fprint(w, "[ERROR] saving object")
        fmt.Fprintf(os.Stderr, "%s", err.Error())
        return
    }
    fmt.Fprint(w, "Saved!")
}
func SearchTenant(w http.ResponseWriter, r *http.Request) {
    obj, err := ProcessPreSteps(w, r, model.Tenant{})
    if err != nil {
        fmt.Fprintf(w, "%s", err.Error())
        return
    }
    out := obj.Search()
    rows := []map[string]any{}
	fieldList := []string{}
	for _, v := range out {
		_fieldList, r := ag.ConvertStruct2Map(v)
		rows = append(rows, r)
		fieldList = _fieldList
	}
	AllTemplate.ExecuteTemplate(w, "search-result.html", map[string]any{
		"fieldList": fieldList,
		"rows":      rows,
        "tableName": "Tenant",
	})
}
func DeleteTenant(w http.ResponseWriter, r *http.Request) {
    if Id := ParseIdFromRequest(r); Id != -1 {
        if err := model.DeleteTenantByID(Id); err != nil {
            fmt.Fprint(w, "[ERROR]")
            fmt.Fprintf(os.Stderr, "[ERROR] deleting %s\n", err.Error())
            return
        }
        fmt.Fprint(w, "[OK] Deleted. Refresh the search button to get new rows to display")
    } else {
        fmt.Fprint(w, "[ERROR] ParseIdFromRequest Delete")
    }
}
func GetTenantById(w http.ResponseWriter, r *http.Request) {
    if Id := ParseIdFromRequest(r); Id != -1 {
        if obj := model.GetTenantByID(Id); obj != nil {
            AllTemplate.ExecuteTemplate(w, "Tenant.html", obj)
            return
        } else {
            fmt.Fprint(w, "[ERROR] Not found")
        }
    } else {
        fmt.Fprint(w, "[ERROR] ParseIdFromRequest GetTenantById")
    }
    return
}
func Property_manager(w http.ResponseWriter, r *http.Request) {
    obj, err := ProcessPreSteps(w, r, model.Property_manager{})
    if err != nil {
        fmt.Fprintf(w, "%s", err.Error())
        return
    }
    if err := obj.Save(); err != nil {
        fmt.Fprint(w, "[ERROR] saving object")
        fmt.Fprintf(os.Stderr, "%s", err.Error())
        return
    }
    fmt.Fprint(w, "Saved!")
}
func SearchProperty_manager(w http.ResponseWriter, r *http.Request) {
    obj, err := ProcessPreSteps(w, r, model.Property_manager{})
    if err != nil {
        fmt.Fprintf(w, "%s", err.Error())
        return
    }
    out := obj.Search()
    rows := []map[string]any{}
	fieldList := []string{}
	for _, v := range out {
		_fieldList, r := ag.ConvertStruct2Map(v)
		rows = append(rows, r)
		fieldList = _fieldList
	}
	AllTemplate.ExecuteTemplate(w, "search-result.html", map[string]any{
		"fieldList": fieldList,
		"rows":      rows,
        "tableName": "Property_manager",
	})
}
func DeleteProperty_manager(w http.ResponseWriter, r *http.Request) {
    if Id := ParseIdFromRequest(r); Id != -1 {
        if err := model.DeleteProperty_managerByID(Id); err != nil {
            fmt.Fprint(w, "[ERROR]")
            fmt.Fprintf(os.Stderr, "[ERROR] deleting %s\n", err.Error())
            return
        }
        fmt.Fprint(w, "[OK] Deleted. Refresh the search button to get new rows to display")
    } else {
        fmt.Fprint(w, "[ERROR] ParseIdFromRequest Delete")
    }
}
func GetProperty_managerById(w http.ResponseWriter, r *http.Request) {
    if Id := ParseIdFromRequest(r); Id != -1 {
        if obj := model.GetProperty_managerByID(Id); obj != nil {
            AllTemplate.ExecuteTemplate(w, "Property_manager.html", obj)
            return
        } else {
            fmt.Fprint(w, "[ERROR] Not found")
        }
    } else {
        fmt.Fprint(w, "[ERROR] ParseIdFromRequest GetProperty_managerById")
    }
    return
}
func Property(w http.ResponseWriter, r *http.Request) {
    obj, err := ProcessPreSteps(w, r, model.Property{})
    if err != nil {
        fmt.Fprintf(w, "%s", err.Error())
        return
    }
    if err := obj.Save(); err != nil {
        fmt.Fprint(w, "[ERROR] saving object")
        fmt.Fprintf(os.Stderr, "%s", err.Error())
        return
    }
    fmt.Fprint(w, "Saved!")
}
func SearchProperty(w http.ResponseWriter, r *http.Request) {
    obj, err := ProcessPreSteps(w, r, model.Property{})
    if err != nil {
        fmt.Fprintf(w, "%s", err.Error())
        return
    }
    out := obj.Search()
    rows := []map[string]any{}
	fieldList := []string{}
	for _, v := range out {
		_fieldList, r := ag.ConvertStruct2Map(v)
		rows = append(rows, r)
		fieldList = _fieldList
	}
	AllTemplate.ExecuteTemplate(w, "search-result.html", map[string]any{
		"fieldList": fieldList,
		"rows":      rows,
        "tableName": "Property",
	})
}
func DeleteProperty(w http.ResponseWriter, r *http.Request) {
    if Id := ParseIdFromRequest(r); Id != -1 {
        if err := model.DeletePropertyByID(Id); err != nil {
            fmt.Fprint(w, "[ERROR]")
            fmt.Fprintf(os.Stderr, "[ERROR] deleting %s\n", err.Error())
            return
        }
        fmt.Fprint(w, "[OK] Deleted. Refresh the search button to get new rows to display")
    } else {
        fmt.Fprint(w, "[ERROR] ParseIdFromRequest Delete")
    }
}
func GetPropertyById(w http.ResponseWriter, r *http.Request) {
    if Id := ParseIdFromRequest(r); Id != -1 {
        if obj := model.GetPropertyByID(Id); obj != nil {
            AllTemplate.ExecuteTemplate(w, "Property.html", obj)
            return
        } else {
            fmt.Fprint(w, "[ERROR] Not found")
        }
    } else {
        fmt.Fprint(w, "[ERROR] ParseIdFromRequest GetPropertyById")
    }
    return
}
func Contract(w http.ResponseWriter, r *http.Request) {
    obj, err := ProcessPreSteps(w, r, model.Contract{})
    if err != nil {
        fmt.Fprintf(w, "%s", err.Error())
        return
    }
    if err := obj.Save(); err != nil {
        fmt.Fprint(w, "[ERROR] saving object")
        fmt.Fprintf(os.Stderr, "%s", err.Error())
        return
    }
    fmt.Fprint(w, "Saved!")
}
func SearchContract(w http.ResponseWriter, r *http.Request) {
    obj, err := ProcessPreSteps(w, r, model.Contract{})
    if err != nil {
        fmt.Fprintf(w, "%s", err.Error())
        return
    }
    out := obj.Search()
    rows := []map[string]any{}
	fieldList := []string{}
	for _, v := range out {
		_fieldList, r := ag.ConvertStruct2Map(v)
		rows = append(rows, r)
		fieldList = _fieldList
	}
	AllTemplate.ExecuteTemplate(w, "search-result.html", map[string]any{
		"fieldList": fieldList,
		"rows":      rows,
        "tableName": "Contract",
	})
}
func DeleteContract(w http.ResponseWriter, r *http.Request) {
    if Id := ParseIdFromRequest(r); Id != -1 {
        if err := model.DeleteContractByID(Id); err != nil {
            fmt.Fprint(w, "[ERROR]")
            fmt.Fprintf(os.Stderr, "[ERROR] deleting %s\n", err.Error())
            return
        }
        fmt.Fprint(w, "[OK] Deleted. Refresh the search button to get new rows to display")
    } else {
        fmt.Fprint(w, "[ERROR] ParseIdFromRequest Delete")
    }
}
func GetContractById(w http.ResponseWriter, r *http.Request) {
    if Id := ParseIdFromRequest(r); Id != -1 {
        if obj := model.GetContractByID(Id); obj != nil {
            AllTemplate.ExecuteTemplate(w, "Contract.html", obj)
            return
        } else {
            fmt.Fprint(w, "[ERROR] Not found")
        }
    } else {
        fmt.Fprint(w, "[ERROR] ParseIdFromRequest GetContractById")
    }
    return
}
func Account(w http.ResponseWriter, r *http.Request) {
    obj, err := ProcessPreSteps(w, r, model.Account{})
    if err != nil {
        fmt.Fprintf(w, "%s", err.Error())
        return
    }
    if err := obj.Save(); err != nil {
        fmt.Fprint(w, "[ERROR] saving object")
        fmt.Fprintf(os.Stderr, "%s", err.Error())
        return
    }
    fmt.Fprint(w, "Saved!")
}
func SearchAccount(w http.ResponseWriter, r *http.Request) {
    obj, err := ProcessPreSteps(w, r, model.Account{})
    if err != nil {
        fmt.Fprintf(w, "%s", err.Error())
        return
    }
    out := obj.Search()
    rows := []map[string]any{}
	fieldList := []string{}
	for _, v := range out {
		_fieldList, r := ag.ConvertStruct2Map(v)
		rows = append(rows, r)
		fieldList = _fieldList
	}
	AllTemplate.ExecuteTemplate(w, "search-result.html", map[string]any{
		"fieldList": fieldList,
		"rows":      rows,
        "tableName": "Account",
	})
}
func DeleteAccount(w http.ResponseWriter, r *http.Request) {
    if Id := ParseIdFromRequest(r); Id != -1 {
        if err := model.DeleteAccountByID(Id); err != nil {
            fmt.Fprint(w, "[ERROR]")
            fmt.Fprintf(os.Stderr, "[ERROR] deleting %s\n", err.Error())
            return
        }
        fmt.Fprint(w, "[OK] Deleted. Refresh the search button to get new rows to display")
    } else {
        fmt.Fprint(w, "[ERROR] ParseIdFromRequest Delete")
    }
}
func GetAccountById(w http.ResponseWriter, r *http.Request) {
    if Id := ParseIdFromRequest(r); Id != -1 {
        if obj := model.GetAccountByID(Id); obj != nil {
            AllTemplate.ExecuteTemplate(w, "Account.html", obj)
            return
        } else {
            fmt.Fprint(w, "[ERROR] Not found")
        }
    } else {
        fmt.Fprint(w, "[ERROR] ParseIdFromRequest GetAccountById")
    }
    return
}
func Payment(w http.ResponseWriter, r *http.Request) {
    obj, err := ProcessPreSteps(w, r, model.Payment{})
    if err != nil {
        fmt.Fprintf(w, "%s", err.Error())
        return
    }
    if err := obj.Save(); err != nil {
        fmt.Fprint(w, "[ERROR] saving object")
        fmt.Fprintf(os.Stderr, "%s", err.Error())
        return
    }
    fmt.Fprint(w, "Saved!")
}
func SearchPayment(w http.ResponseWriter, r *http.Request) {
    obj, err := ProcessPreSteps(w, r, model.Payment{})
    if err != nil {
        fmt.Fprintf(w, "%s", err.Error())
        return
    }
    out := obj.Search()
    rows := []map[string]any{}
	fieldList := []string{}
	for _, v := range out {
		_fieldList, r := ag.ConvertStruct2Map(v)
		rows = append(rows, r)
		fieldList = _fieldList
	}
	AllTemplate.ExecuteTemplate(w, "search-result.html", map[string]any{
		"fieldList": fieldList,
		"rows":      rows,
        "tableName": "Payment",
	})
}
func DeletePayment(w http.ResponseWriter, r *http.Request) {
    if Id := ParseIdFromRequest(r); Id != -1 {
        if err := model.DeletePaymentByID(Id); err != nil {
            fmt.Fprint(w, "[ERROR]")
            fmt.Fprintf(os.Stderr, "[ERROR] deleting %s\n", err.Error())
            return
        }
        fmt.Fprint(w, "[OK] Deleted. Refresh the search button to get new rows to display")
    } else {
        fmt.Fprint(w, "[ERROR] ParseIdFromRequest Delete")
    }
}
func GetPaymentById(w http.ResponseWriter, r *http.Request) {
    if Id := ParseIdFromRequest(r); Id != -1 {
        if obj := model.GetPaymentByID(Id); obj != nil {
            AllTemplate.ExecuteTemplate(w, "Payment.html", obj)
            return
        } else {
            fmt.Fprint(w, "[ERROR] Not found")
        }
    } else {
        fmt.Fprint(w, "[ERROR] ParseIdFromRequest GetPaymentById")
    }
    return
}
func Invoice(w http.ResponseWriter, r *http.Request) {
    obj, err := ProcessPreSteps(w, r, model.Invoice{})
    if err != nil {
        fmt.Fprintf(w, "%s", err.Error())
        return
    }
    if err := obj.Save(); err != nil {
        fmt.Fprint(w, "[ERROR] saving object")
        fmt.Fprintf(os.Stderr, "%s", err.Error())
        return
    }
    fmt.Fprint(w, "Saved!")
}
func SearchInvoice(w http.ResponseWriter, r *http.Request) {
    obj, err := ProcessPreSteps(w, r, model.Invoice{})
    if err != nil {
        fmt.Fprintf(w, "%s", err.Error())
        return
    }
    out := obj.Search()
    rows := []map[string]any{}
	fieldList := []string{}
	for _, v := range out {
		_fieldList, r := ag.ConvertStruct2Map(v)
		rows = append(rows, r)
		fieldList = _fieldList
	}
	AllTemplate.ExecuteTemplate(w, "search-result.html", map[string]any{
		"fieldList": fieldList,
		"rows":      rows,
        "tableName": "Invoice",
	})
}
func DeleteInvoice(w http.ResponseWriter, r *http.Request) {
    if Id := ParseIdFromRequest(r); Id != -1 {
        if err := model.DeleteInvoiceByID(Id); err != nil {
            fmt.Fprint(w, "[ERROR]")
            fmt.Fprintf(os.Stderr, "[ERROR] deleting %s\n", err.Error())
            return
        }
        fmt.Fprint(w, "[OK] Deleted. Refresh the search button to get new rows to display")
    } else {
        fmt.Fprint(w, "[ERROR] ParseIdFromRequest Delete")
    }
}
func GetInvoiceById(w http.ResponseWriter, r *http.Request) {
    if Id := ParseIdFromRequest(r); Id != -1 {
        if obj := model.GetInvoiceByID(Id); obj != nil {
            AllTemplate.ExecuteTemplate(w, "Invoice.html", obj)
            return
        } else {
            fmt.Fprint(w, "[ERROR] Not found")
        }
    } else {
        fmt.Fprint(w, "[ERROR] ParseIdFromRequest GetInvoiceById")
    }
    return
}
func Maintenance_request(w http.ResponseWriter, r *http.Request) {
    obj, err := ProcessPreSteps(w, r, model.Maintenance_request{})
    if err != nil {
        fmt.Fprintf(w, "%s", err.Error())
        return
    }
    if err := obj.Save(); err != nil {
        fmt.Fprint(w, "[ERROR] saving object")
        fmt.Fprintf(os.Stderr, "%s", err.Error())
        return
    }
    fmt.Fprint(w, "Saved!")
}
func SearchMaintenance_request(w http.ResponseWriter, r *http.Request) {
    obj, err := ProcessPreSteps(w, r, model.Maintenance_request{})
    if err != nil {
        fmt.Fprintf(w, "%s", err.Error())
        return
    }
    out := obj.Search()
    rows := []map[string]any{}
	fieldList := []string{}
	for _, v := range out {
		_fieldList, r := ag.ConvertStruct2Map(v)
		rows = append(rows, r)
		fieldList = _fieldList
	}
	AllTemplate.ExecuteTemplate(w, "search-result.html", map[string]any{
		"fieldList": fieldList,
		"rows":      rows,
        "tableName": "Maintenance_request",
	})
}
func DeleteMaintenance_request(w http.ResponseWriter, r *http.Request) {
    if Id := ParseIdFromRequest(r); Id != -1 {
        if err := model.DeleteMaintenance_requestByID(Id); err != nil {
            fmt.Fprint(w, "[ERROR]")
            fmt.Fprintf(os.Stderr, "[ERROR] deleting %s\n", err.Error())
            return
        }
        fmt.Fprint(w, "[OK] Deleted. Refresh the search button to get new rows to display")
    } else {
        fmt.Fprint(w, "[ERROR] ParseIdFromRequest Delete")
    }
}
func GetMaintenance_requestById(w http.ResponseWriter, r *http.Request) {
    if Id := ParseIdFromRequest(r); Id != -1 {
        if obj := model.GetMaintenance_requestByID(Id); obj != nil {
            AllTemplate.ExecuteTemplate(w, "Maintenance_request.html", obj)
            return
        } else {
            fmt.Fprint(w, "[ERROR] Not found")
        }
    } else {
        fmt.Fprint(w, "[ERROR] ParseIdFromRequest GetMaintenance_requestById")
    }
    return
}

func AddHandler(mux *http.ServeMux, Cfg *configs.Config) {
    mux.HandleFunc("POST "+Cfg.PathBase+"/tenant/search", SearchTenant)
    mux.HandleFunc("POST "+Cfg.PathBase+"/tenant/delete/{id}", DeleteTenant)
    mux.HandleFunc("GET "+Cfg.PathBase+"/tenant/get/{id}", GetTenantById)
    mux.HandleFunc("POST "+Cfg.PathBase+"/tenant", Tenant)
    mux.HandleFunc("POST "+Cfg.PathBase+"/property_manager/search", SearchProperty_manager)
    mux.HandleFunc("POST "+Cfg.PathBase+"/property_manager/delete/{id}", DeleteProperty_manager)
    mux.HandleFunc("GET "+Cfg.PathBase+"/property_manager/get/{id}", GetProperty_managerById)
    mux.HandleFunc("POST "+Cfg.PathBase+"/property_manager", Property_manager)
    mux.HandleFunc("POST "+Cfg.PathBase+"/property/search", SearchProperty)
    mux.HandleFunc("POST "+Cfg.PathBase+"/property/delete/{id}", DeleteProperty)
    mux.HandleFunc("GET "+Cfg.PathBase+"/property/get/{id}", GetPropertyById)
    mux.HandleFunc("POST "+Cfg.PathBase+"/property", Property)
    mux.HandleFunc("POST "+Cfg.PathBase+"/contract/search", SearchContract)
    mux.HandleFunc("POST "+Cfg.PathBase+"/contract/delete/{id}", DeleteContract)
    mux.HandleFunc("GET "+Cfg.PathBase+"/contract/get/{id}", GetContractById)
    mux.HandleFunc("POST "+Cfg.PathBase+"/contract", Contract)
    mux.HandleFunc("POST "+Cfg.PathBase+"/account/search", SearchAccount)
    mux.HandleFunc("POST "+Cfg.PathBase+"/account/delete/{id}", DeleteAccount)
    mux.HandleFunc("GET "+Cfg.PathBase+"/account/get/{id}", GetAccountById)
    mux.HandleFunc("POST "+Cfg.PathBase+"/account", Account)
    mux.HandleFunc("POST "+Cfg.PathBase+"/payment/search", SearchPayment)
    mux.HandleFunc("POST "+Cfg.PathBase+"/payment/delete/{id}", DeletePayment)
    mux.HandleFunc("GET "+Cfg.PathBase+"/payment/get/{id}", GetPaymentById)
    mux.HandleFunc("POST "+Cfg.PathBase+"/payment", Payment)
    mux.HandleFunc("POST "+Cfg.PathBase+"/invoice/search", SearchInvoice)
    mux.HandleFunc("POST "+Cfg.PathBase+"/invoice/delete/{id}", DeleteInvoice)
    mux.HandleFunc("GET "+Cfg.PathBase+"/invoice/get/{id}", GetInvoiceById)
    mux.HandleFunc("POST "+Cfg.PathBase+"/invoice", Invoice)
    mux.HandleFunc("POST "+Cfg.PathBase+"/maintenance_request/search", SearchMaintenance_request)
    mux.HandleFunc("POST "+Cfg.PathBase+"/maintenance_request/delete/{id}", DeleteMaintenance_request)
    mux.HandleFunc("GET "+Cfg.PathBase+"/maintenance_request/get/{id}", GetMaintenance_requestById)
    mux.HandleFunc("POST "+Cfg.PathBase+"/maintenance_request", Maintenance_request)
}

func ParseIdFromRequest(r *http.Request) int64 {
    id := r.PathValue("id")
    Id, err := strconv.ParseInt(id, 10, 64)
    if err != nil {
        fmt.Fprintf(os.Stderr, "[ERROR] Parse ID %s | %s\n", id, err.Error())
        return -1
    }
    return Id
}
// End app-handler.go.tmpl

func loadAllTemplates() *template.Template {
	t := template.New("tmpl")
	myFuncmap := ag.GoTemplateFuncMap
	// CallTemplate allows us to call a pre-defined template name (using block or define); the name can be dynamic (a variable with value at runtime).
	// Standard go `template` does not support this feature. See usage in the index.html
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
	AllTemplate.ExecuteTemplate(w, "index.html", model.AllForms)
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
