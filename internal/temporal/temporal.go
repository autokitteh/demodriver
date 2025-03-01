package temporal

import (
	"context"
	"fmt"
	"log/slog"

	"go.temporal.io/sdk/client"
	"go.uber.org/fx"

	"go.autokitteh.dev/demodriver/internal/app"
)

type config struct {
	HostPort  string
	Namespace string
}

func New() fx.Option {
	return app.Module[config](
		"temporal",
		fx.Decorate(func(in *config) *config {
			cfg := *in

			if cfg.HostPort == "" {
				cfg.HostPort = "localhost:7233"
			}

			if cfg.Namespace == "" {
				cfg.Namespace = "default"
			}

			return &cfg
		}),
		fx.Provide(func(lc fx.Lifecycle, cfg *config, l *slog.Logger) (client.Client, error) {
			cl, err := client.NewLazyClient(client.Options{
				HostPort:  cfg.HostPort,
				Namespace: cfg.Namespace,
				Logger:    l.WithGroup("client"),
			})
			if err != nil {
				return nil, err
			}

			lc.Append(fx.Hook{
				OnStart: func(ctx context.Context) error {
					_, err := cl.CheckHealth(ctx, &client.CheckHealthRequest{})
					if err != nil {
						return fmt.Errorf("health (hostport=%s): %w", cfg.HostPort, err)
					}

					l.Info("temporal connection is healthy", "hostport", cfg.HostPort, "namespace", cfg.Namespace)

					return nil
				},
				OnStop: func(context.Context) error {
					cl.Close()
					return nil
				},
			})

			return cl, nil
		}),
		fx.Invoke(func(client client.Client) {}),
	)
}
