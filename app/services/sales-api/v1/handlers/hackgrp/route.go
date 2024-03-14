package hackgrp

import (
	"net/http"

	"github.com/islamghany/service/foundation/web"
)

func Routes(app *web.App) {
	app.Handle(http.MethodGet, "/hack", Hack)
}
