package hackgrp

import (
	"context"
	"net/http"

	"github.com/islamghany/service/foundation/web"
)

func Hack(ctx context.Context, w http.ResponseWriter, r *http.Request) error {

	status := struct {
		Status string
	}{
		Status: "OK",
	}

	return web.Respond(ctx, w, status, http.StatusOK)
}
