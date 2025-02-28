package app

import (
	"log/slog"

	"go.uber.org/fx"
)

func Module[Config any](name string, opts ...fx.Option) fx.Option {
	opts = append(
		[]fx.Option{
			fx.Decorate(func(l *slog.Logger) *slog.Logger {
				// TODO: This just names the argument, not the logger.
				return l.WithGroup(name)
			}),
			provideConfig[Config](name),
		},
		opts...,
	)

	return fx.Module(name, opts...)
}
