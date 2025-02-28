package app

import (
	"fmt"
	"strings"

	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/v2"
	"go.uber.org/fx"
)

type Config struct {
	*koanf.Koanf
	prefix string
}

func newConfig(prefix string) *Config {
	return &Config{Koanf: koanf.New("."), prefix: strings.ToUpper(prefix) + "_"}
}

func (c *Config) load(path string, dst any) error {
	prefix := strings.ToUpper(fmt.Sprintf("%s_%s", c.prefix, path))

	if err := c.Koanf.Load(
		env.Provider(
			prefix,
			".",
			func(s string) string {
				return strings.Replace(
					strings.ToLower(strings.TrimPrefix(s, prefix)),
					"__",
					".",
					-1,
				)
			}),
		nil,
	); err != nil {
		return err
	}

	return c.Koanf.Unmarshal(path, dst)
}

func provideConfig[T any](name string, fs ...func(*T) error) fx.Option {
	return fx.Provide(func(c *Config) (*T, error) {
		t := new(T)

		if err := c.load(name, t); err != nil {
			return nil, err
		}

		for _, f := range fs {
			if err := f(t); err != nil {
				return nil, err
			}
		}

		return t, nil
	})
}
