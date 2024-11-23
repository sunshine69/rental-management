package main

import (
	"net/http"

	"github.com/sunshine69/rental-management/configs"
	"github.com/sunshine69/rental-management/web/app"
)

var Cfg configs.Config

func main() {
	Cfg := configs.InitConfig()

	r := http.NewServeMux()

	// api.RouteRegisterAccount(r, Cfg.PathBase)
	// api.RouteRegisterContract(r, Cfg.PathBase)
	// api.RouteRegisterInvoice(r, Cfg.PathBase)
	// api.RouteRegisterMaintenance_request(r, Cfg.PathBase)
	// api.RouteRegisterPayment(r, Cfg.PathBase)
	// api.RouteRegisterProperty(r, Cfg.PathBase)
	// api.RouteRegisterProperty_manager(r, Cfg.PathBase)
	// api.RouteRegisterTenant(r, Cfg.PathBase)

	app.StartWebApp(r, &Cfg)

	wait_chan := make(chan int)
	r.HandleFunc("POST /quit", func(w http.ResponseWriter, r *http.Request) { wait_chan <- 1 })

	go http.ListenAndServe(":"+Cfg.Port, r)

	// go u.RunSystemCommandV2(fmt.Sprintf("%s http://localhost:%s/", "google-chrome ", ListenPort), true)

	<-wait_chan
}
