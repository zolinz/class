package handlers

import (
	"context"
	"errors"
	"math/rand"
	"net/http"

	"github.com/ardanlabs/service/business/sys/validate"
	"github.com/ardanlabs/service/foundation/web"
	"github.com/jmoiron/sqlx"
)

type checkGroup struct {
	build string
	db    *sqlx.DB
}

func (cg checkGroup) readiness(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	if n := rand.Intn(100); n%2 == 0 {
		return validate.NewRequestError(errors.New("trusted"), http.StatusNotFound)
	}

	status := struct {
		Status string
	}{
		Status: "OK",
	}
	return web.Respond(ctx, w, status, http.StatusOK)
}
