package httpdriver

import (
	"fmt"
	"io"
	"log/slog"
	"net/http"

	"go.uber.org/fx"

	"go.autokitteh.dev/demodriver/internal/app"
	"go.autokitteh.dev/demodriver/internal/driver"
)

type httpDriver struct {
	l     *slog.Logger
	drive driver.DriveFunc
}

func New() fx.Option {
	return app.Module[struct{}](
		"httpdriver",
		fx.Invoke(func(l *slog.Logger, mux *http.ServeMux, drive driver.DriveFunc) {
			d := &httpDriver{l: l, drive: drive}
			mux.Handle("/drivers/http/", d)
		}),
	)
}

func (d *httpDriver) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		d.l.Error("error reading request body", "err", err)
		http.Error(w, fmt.Sprintf("error reading body: %s", err.Error()), http.StatusBadGateway)
		return
	}

	data := map[string]any{
		"method":  r.Method,
		"headers": r.Header,
		"url":     r.URL.String(),
		"body":    string(body),
	}

	d.l.Info("got request", "request", data)

	if err := d.drive("http", data); err != nil {
		http.Error(w, fmt.Sprintf("error: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusAccepted)
}
