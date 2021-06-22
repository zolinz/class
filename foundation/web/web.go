package web

import (
	"context"
	"github.com/dimfeld/httptreemux/v5"
	"github.com/google/uuid"
	"net/http"
	"os"
	"syscall"
	"time"
)

// ctxKey represents the type of value for the context key.
type ctxKey int

// KeyValues is how request values are stored/retrieved.
const KeyValues ctxKey = 1


type Handler func(ctx context.Context, w http.ResponseWriter, r *http.Request) error



// Values represent state for each request.
type Values struct {
	TraceID    string
	Now        time.Time
	StatusCode int
}

type App struct{
	*httptreemux.ContextMux
	shutdown chan os.Signal
	mw []Middleware
}

func (a *App) SignalShutdown() {
	a.shutdown <- syscall.SIGTERM
}

func NewApp(shutdown chan os.Signal, mw ...Middleware ) *App{
	app := App{
		ContextMux:  httptreemux.NewContextMux(),
		shutdown:    shutdown,
		mw: 		 mw,
	}
	return &app
}

func (a *App) Handle(method string, path string, handler Handler, mw ...Middleware){

	// First wrap handler specific middleware around this handler.
	handler = wrapMiddleware(mw, handler)

	// Add the application's general middleware to the handler chain.
	handler = wrapMiddleware(a.mw, handler)

	h := func(w http.ResponseWriter, r *http.Request){

		v := Values{
			TraceID: uuid.New().String(),
			Now:     time.Now(),
		}
		ctx := context.WithValue(r.Context(), KeyValues, &v)
		// Call the wrapped handler functions.
		if err := handler(ctx, w, r); err != nil {
			a.SignalShutdown()
			return
		}
	}
	a.ContextMux.Handle(method, path, h)
}