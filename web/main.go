package main

import (
	"fmt"
	"net/http"

	"github.com/go-playground/form/v4"
	"github.com/go-playground/validator/v10"
	u "github.com/sunshine69/golang-tools/utils"
	"github.com/sunshine69/rental-management/api"
	"github.com/sunshine69/rental-management/configs"
	"github.com/sunshine69/rental-management/web/app"
)

var (
	validate *validator.Validate
	// use a single instance of Decoder, it caches struct info
	formDecoder *form.Decoder
)

// As the content of vars-ansible.yaml is large creating one html form to hold data is big thus split them into several forms and collecting data
// the sequence steps. Some of the fields values depending the value of the previous field(s) thus we will write custom validation rules
type Form1 struct {
}

var (
	form1 Form1
)

var Cfg configs.Config

func main() {
	Cfg := configs.InitConfig()

	validate = validator.New(validator.WithRequiredStructEnabled())
	formDecoder = form.NewDecoder()

	fmt.Printf("[DEBUG] form initialization %s\n", u.JsonDump(form1, "  "))

	r := http.NewServeMux()

	api.RouteRegisterAccount(r, Cfg.PathBase)
	api.RouteRegisterContract(r, Cfg.PathBase)
	api.RouteRegisterInvoice(r, Cfg.PathBase)
	api.RouteRegisterMaintenance_request(r, Cfg.PathBase)
	api.RouteRegisterPayment(r, Cfg.PathBase)
	api.RouteRegisterProperty(r, Cfg.PathBase)
	api.RouteRegisterProperty_manager(r, Cfg.PathBase)
	api.RouteRegisterTenant(r, Cfg.PathBase)

	app.RouteRegisterApp(r, &Cfg)

	wait_chan := make(chan int)
	r.HandleFunc("POST /quit", func(w http.ResponseWriter, r *http.Request) { wait_chan <- 1 })

	go http.ListenAndServe(":"+Cfg.Port, r)

	// go u.RunSystemCommandV2(fmt.Sprintf("%s http://localhost:%s/", "google-chrome ", ListenPort), true)

	<-wait_chan
}
