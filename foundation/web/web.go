package web

import (
	"context"
	"net/http"
	"os"

	"github.com/dimfeld/httptreemux/v5"
)

type Handler func(ctx context.Context, w http.ResponseWriter, r *http.Request) error

type App struct {
	*httptreemux.ContextMux
	shutdown chan os.Signal
}

func NewApp(shutdown chan os.Signal) *App {
	return &App{
		ContextMux: httptreemux.NewContextMux(),
		shutdown:   shutdown,
	}
}

func (a *App) Handle(method, path string, handler Handler) {

	h := func(w http.ResponseWriter, r *http.Request) {
		// Add any Logging, Tracing, or Metrics here. For example, we could wrap the

		if err := handler(r.Context(), w, r); err != nil {
			// Log the error.
			return
		}
		// Add the handler to the ContextMux.

	}
	a.ContextMux.Handle(method, path, h)
}
