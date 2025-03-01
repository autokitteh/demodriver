package driverdashboard

import (
	"encoding/json"
	"net/http"

	"github.com/tidwall/pretty"
	"go.uber.org/fx"

	"go.autokitteh.dev/demodriver/internal/app"
	"go.autokitteh.dev/demodriver/internal/driver"
	"go.autokitteh.dev/demodriver/internal/types"
)

func New() fx.Option {
	return app.Module[struct{}](
		"driverdashboard",
		fx.Invoke(func(driver driver.Driver, mux *http.ServeMux) {
			dashboard := &dashboard{driver: driver}

			mux.HandleFunc("GET /triggers", dashboard.triggers)
			mux.HandleFunc("GET /workflows", dashboard.workflows)
			mux.HandleFunc("GET /signals", dashboard.signals)
		}),
	)
}

type dashboard struct{ driver driver.Driver }

func (d *dashboard) triggers(w http.ResponseWriter, r *http.Request) {
	display(w, d.driver.ListTriggers(r.Context()))
}

func (d *dashboard) workflows(w http.ResponseWriter, r *http.Request) {
	tname := r.URL.Query().Get("tname")

	ws, err := d.driver.ListWorkflows(r.Context(), types.TriggerName(tname))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	display(w, ws)
}

func (d *dashboard) signals(w http.ResponseWriter, r *http.Request) {
	src := r.URL.Query().Get("src")
	wid := r.URL.Query().Get("wid")

	ws, err := d.driver.ListSignals(r.Context(), src, types.WorkflowID(wid))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	display(w, ws)
}

func display(w http.ResponseWriter, x any) {
	bs, err := json.Marshal(x)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, _ = w.Write(pretty.Pretty(bs))
}
