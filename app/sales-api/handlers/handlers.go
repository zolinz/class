package handlers

import (
	"github.com/zolinz/class/business/auth"
	"github.com/zolinz/class/business/mid"
	"github.com/zolinz/class/foundation/web"
	"log"
	"net/http"
	"os"
)

func API(build string, shutdown chan os.Signal, log *log.Logger, a *auth.Auth) *web.App {
	app := web.NewApp(shutdown, mid.Logger(log), mid.Errors(log), mid.Metrics(), mid.Panics(log))


	check := check{
		log : log,
		build: build,
	}

	app.Handle(http.MethodGet, "/readiness", check.readiness )
	//app.Handle(http.MethodGet, "/readiness", check.readiness, mid.Authenticate(a), mid.Authorize(auth.RoleAdmin))
	return app

}
