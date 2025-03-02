package driverapi

import (
	"encoding/json"
	"net/http"

	"go.uber.org/fx"

	"go.autokitteh.dev/demodriver/internal/app"
	"go.autokitteh.dev/demodriver/internal/driver"
	"go.autokitteh.dev/demodriver/internal/types"
)

func New() fx.Option {
	return app.Module[struct{}](
		"driverapi",
		fx.Invoke(func(d driver.Driver, mux *http.ServeMux) {
			api := &api{driver: d}

			mux.HandleFunc("POST /api/signals", api.createSignal)
		}),
	)
}

type api struct{ driver driver.Driver }

func (a *api) createSignal(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Signal types.Signal `json:"signal"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := a.driver.RegisterSignal(r.Context(), req.Signal); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}
