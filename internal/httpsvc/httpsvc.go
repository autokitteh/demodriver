package httpsvc

import (
	"context"
	"log/slog"
	"net"
	"net/http"

	"go.uber.org/fx"

	"go.autokitteh.dev/demodriver/internal/app"
)

type config struct {
	Addr string `koanf:"addr"`
}

func New() fx.Option {
	return app.Module[config](
		"http",
		fx.Decorate(func(cfg *config) *config {
			if cfg.Addr == "" {
				cfg.Addr = ":9001"
			}

			return cfg
		}),
		fx.Provide(func(lc fx.Lifecycle, l *slog.Logger, cfg *config) (*http.Server, *http.ServeMux) {
			srv := &http.Server{Addr: cfg.Addr}

			lc.Append(fx.Hook{
				OnStart: func(ctx context.Context) error {
					ln, err := net.Listen("tcp", srv.Addr)
					if err != nil {
						return err
					}
					l.Info("starting HTTP server", "addr", srv.Addr)
					go srv.Serve(ln)
					return nil
				},
				OnStop: func(ctx context.Context) error {
					return srv.Shutdown(ctx)
				},
			})

			mux := http.NewServeMux()
			srv.Handler = mux

			return srv, mux
		}),
	)
}
