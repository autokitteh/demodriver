package ddapp

import (
	"fmt"
	"log/slog"
	"net/http"

	"go.uber.org/fx"

	"go.autokitteh.dev/demodriver/internal/app"
	"go.autokitteh.dev/demodriver/internal/db"
	"go.autokitteh.dev/demodriver/internal/driver"
	"go.autokitteh.dev/demodriver/internal/driverapi"
	"go.autokitteh.dev/demodriver/internal/driverdashboard"
	"go.autokitteh.dev/demodriver/internal/httpsvc"
	"go.autokitteh.dev/demodriver/internal/sources/httpsource"
	"go.autokitteh.dev/demodriver/internal/sources/slacksource"
	"go.autokitteh.dev/demodriver/internal/temporal"
)

func New(l *slog.Logger, name string) *fx.App {
	return app.New(
		l,
		name,

		db.New(),
		httpsvc.New(),
		temporal.New(),
		driver.New(),
		driverapi.New(),
		driverdashboard.New(),
		httpsource.New(),
		slacksource.New(),

		fx.Invoke(func(mux *http.ServeMux) {
			mux.HandleFunc("/", func(w http.ResponseWriter, _ *http.Request) {
				fmt.Fprintln(w, "meow!")
			})
		}),

		fx.Invoke(func(lc fx.Lifecycle, l *slog.Logger) {
			lc.Append(fx.StartHook(func() {
				l.Info("ready")
			}))
		}),
	)
}
