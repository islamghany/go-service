package hackgrp

import (
	"context"
	"errors"
	"math/rand"
	"net/http"

	"github.com/islamghany/service/business/web/v1/respond"
	"github.com/islamghany/service/foundation/web"
)

func Hack(ctx context.Context, w http.ResponseWriter, r *http.Request) error {

	if n := rand.Intn(100) % 2; n%2 == 0 {
		return respond.NewError(errors.New("trusted error"), http.StatusUnprocessableEntity)
	}
	status := struct {
		Status string
	}{
		Status: "OK",
	}

	return web.Respond(ctx, w, status, http.StatusOK)
}
