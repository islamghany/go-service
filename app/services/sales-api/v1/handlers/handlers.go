package handlers

import (
	"github.com/islamghany/service/app/services/sales-api/v1/handlers/hackgrp"
	v1 "github.com/islamghany/service/business/web/v1"
	"github.com/islamghany/service/foundation/web"
)

type Routes struct{}

func (Routes) Add(app *web.App, cfg v1.APIMuxConfig) {
	hackgrp.Routes(app)
}
