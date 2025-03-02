package driver

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"go.temporal.io/sdk/client"
	"go.uber.org/fx"
	"gopkg.in/yaml.v3"

	"go.autokitteh.dev/demodriver/internal/app"
	"go.autokitteh.dev/demodriver/internal/db"
	"go.autokitteh.dev/demodriver/internal/types"
)

type config struct {
	TriggersPath string `koanf:"triggerspath"`
}

func New() fx.Option {
	return app.Module[config](
		"driver",
		fx.Provide(func(
			cfg *config,
			l *slog.Logger,
			db db.DB,
			client client.Client,
			mux *http.ServeMux,
		) (Driver, DriveFunc, error) {
			var ts map[types.TriggerName]*trigger

			if cfg.TriggersPath != "" {
				var err error
				if ts, err = readTriggers(cfg.TriggersPath); err != nil {
					return nil, nil, fmt.Errorf("triggers: %w", err)
				}

				for _, t := range ts {
					l.Info("registering trigger", "t", t)
				}
			} else {
				l.Warn("no triggers path supplied, no triggers registered")
			}

			d := &driver{cfg: cfg, l: l, client: client, db: db, triggers: ts}

			return d, d.Drive, nil
		}),
	)
}

func readTriggers(path string) (map[types.TriggerName]*trigger, error) {
	bs, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}

		return nil, err
	}

	var tts map[types.TriggerName]types.Trigger

	if err := yaml.Unmarshal(bs, &tts); err != nil {
		return nil, fmt.Errorf("parse: %w", err)
	}

	ts := make(map[types.TriggerName]*trigger, len(tts))

	for name, tt := range tts {
		if tt.Name != "" {
			return nil, fmt.Errorf("trigger %q: already has a name", name)
		}

		f, err := parseFilter(tt.Filter)
		if err != nil {
			return nil, fmt.Errorf("trigger %q: filter: %w", name, err)
		}

		tt.Name = name

		ts[name] = &trigger{
			trigger: tt,
			filter:  f,
		}
	}

	return ts, nil
}
