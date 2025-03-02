package db

import (
	"context"
	"database/sql"
	"log/slog"

	"go.uber.org/fx"
	_ "modernc.org/sqlite"

	"go.autokitteh.dev/demodriver/internal/app"
	"go.autokitteh.dev/demodriver/sqlc/sqlcgen"
)

type config struct {
	DSN string `koanf:"dsn"`
}

func New() fx.Option {
	return app.Module[config](
		"db",
		fx.Decorate(func(in *config) *config {
			cfg := *in

			if cfg.DSN == "" {
				cfg.DSN = "dd.sqlite"
			}

			return &cfg
		}),
		fx.Provide(func(lc fx.Lifecycle, l *slog.Logger, cfg *config) (DB, error) {
			sqldb, err := sql.Open("sqlite", cfg.DSN)
			if err != nil {
				return nil, err
			}

			l.Info("opened database", "dsn", cfg.DSN)

			db := &db{cfg: cfg, l: l, db: sqldb, q: sqlcgen.New(sqldb)}

			lc.Append(fx.Hook{
				OnStart: func(ctx context.Context) error {
					return db.Setup(ctx)
				},
				OnStop: func(ctx context.Context) error {
					return db.Close()
				},
			})

			return db, nil
		}),
	)
}
