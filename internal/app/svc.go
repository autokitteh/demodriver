package app

import (
	"log/slog"

	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
)

func New(l *slog.Logger, name string, opts ...fx.Option) *fx.App {
	opts = append([]fx.Option{
		fx.Supply(l),
		fx.Provide(func() *Config { return newConfig(name) }),
		fx.WithLogger(func() fxevent.Logger {
			sl := &fxevent.SlogLogger{Logger: l}
			sl.UseLogLevel(slog.LevelDebug)
			return sl
		}),
	}, opts...)

	return fx.New(opts...)
}
