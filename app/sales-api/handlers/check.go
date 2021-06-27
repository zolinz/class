package handlers

import (
	"context"
	"github.com/pkg/errors"
	"github.com/zolinz/class/foundation/web"
	"log"
	"math/rand"
	"net/http"
)

type check struct {
	log *log.Logger
}

func (c check) readiness(ctx context.Context, w http.ResponseWriter, r *http.Request) error{

	if n := rand.Intn(100) ; n%2 == 0 {
		return web.NewRequestError(	errors.New("trusted error "), http.StatusBadRequest)
		//panic("oh my god")
	}
	status := struct {
		Status string
	}{
		Status: "OK",
	}


	return web.Respond(ctx, w, status, http.StatusOK)
	//return json.NewEncoder(w).Encode(status)
}