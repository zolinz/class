package mid

import (
	"context"
	"github.com/zolinz/class/foundation/web"
	"log"
	"net/http"
)

func Errors(log *log.Logger) web.Middleware{

	m := func(handler web.Handler) web.Handler{
		h := func(ctx context.Context, w http.ResponseWriter, r *http.Request) error{
			// If the context is missing this value, request the service
			// to be shutdown gracefully.
			v, ok := ctx.Value(web.KeyValues).(*web.Values)
			if !ok {
				return web.NewShutdownError("web value missing from context")
			}

			// Run the handler chain and catch any propagated error.
			if err := handler(ctx, w, r); err != nil {

				// Log the error.
				log.Printf("%s : ERROR     : %v", v.TraceID, err)

				// Respond to the error.
				if err := web.RespondError(ctx, w, err); err != nil {
					return err
				}

				// If we receive the shutdown err we need to return it
				// back to the base handler to shutdown the service.
				if ok := web.IsShutdown(err); ok {
					return err
				}
			}

		return nil
		}

		return h
	}

	return m
}
