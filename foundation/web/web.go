package web

import (
	"context"
	"github.com/dimfeld/httptreemux/v5"
	"net/http"
	"os"
	"syscall"
)


type Handler func(ctx context.Context, w http.ResponseWriter, r *http.Request) error

type App struct{
	*httptreemux.ContextMux
	shutdown chan os.Signal
}

func (a *App) SignalShutdown() {
	a.shutdown <- syscall.SIGTERM
}

func NewApp(shutdown chan os.Signal) *App{
	app := App{
		ContextMux: httptreemux.NewContextMux(),
		shutdown: shutdown,
	}
	return &app
}

func (a *App) Handle(method string, path string, handler Handler){
	h := func(w http.ResponseWriter, r *http.Request){

		// Call the wrapped handler functions.
		if err := handler(r.Context(), w, r); err != nil {
			a.SignalShutdown()
			return
		}
		
	}
	a.ContextMux.Handle(method, path, h)
}