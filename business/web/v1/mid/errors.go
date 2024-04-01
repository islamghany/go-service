package mid

import (
	"context"
	"net/http"

	"github.com/islamghany/service/business/web/v1/respond"
	"github.com/islamghany/service/foundation/logger"
	"github.com/islamghany/service/foundation/validate"
	"github.com/islamghany/service/foundation/web"
)

func Errors(log *logger.Logger) web.Middleware {
	m := func(handler web.Handler) web.Handler {
		h := func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
			if err := handler(ctx, w, r); err != nil {
				log.Error(ctx, "message", "msg", err.Error())

				var er respond.ErrorDocument
				var status int
				switch {
				// trusted  error
				case respond.IsError(err):
					reqErr := respond.GetError(err)
					if validate.IsFieldErrors(reqErr.Err) {
						fieldErrors := validate.GetFieldErrors(reqErr.Err)
						er = respond.ErrorDocument{
							Error:  "data validation error",
							Fields: fieldErrors.Fields(),
						}
						status = reqErr.Status
						break
					}

					er = respond.ErrorDocument{
						Error: reqErr.Error(),
					}
					status = reqErr.Status
					break
				// untrusted error
				default:
					er = respond.ErrorDocument{
						Error: http.StatusText(http.StatusInternalServerError),
					}
					status = http.StatusInternalServerError

				}
				if err := web.Respond(ctx, w, er, status); err != nil {
					return err
				}
				// If we receive the shutdown err we need to return it
				// back to the base handler to shut down the service.
				if web.IsShutdown(err) {
					return err
				}
			}
			return nil
		}
		return h
	}
	return m
}
