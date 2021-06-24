package handlers

import (
	"github.com/zolinz/class/business/mid"
	"github.com/zolinz/class/foundation/web"
	"log"
	"net/http"
	"os"
)

func API(build string, shutdown chan os.Signal, log *log.Logger) *web.App {
	app := web.NewApp(shutdown, mid.Logger(log), mid.Errors(log), mid.Metrics(), mid.Panics(log))


	check := check{
		log : log,
	}

	app.Handle(http.MethodGet, "/readiness", check.readiness)
	return app

}