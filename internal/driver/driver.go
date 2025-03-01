package driver

import (
	"go.uber.org/fx"

	"go.autokitteh.dev/demodriver/internal/app"
)

type config struct{}

type DriveFunc func(src string, data any) error

func New() fx.Option {
	return app.Module[config](
		"driver",
		fx.Provide(func() DriveFunc {
			return func(string, any) error { return nil }
		}),
	)
}
