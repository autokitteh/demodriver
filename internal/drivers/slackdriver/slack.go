package slackdriver

import (
	"context"
	"log/slog"

	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
	"github.com/slack-go/slack/socketmode"
	"go.uber.org/fx"

	"go.autokitteh.dev/demodriver/internal/app"
	"go.autokitteh.dev/demodriver/internal/driver"
)

type config struct {
	BotToken string `koanf:"bot_token"`
	AppToken string `koanf:"app_token"`
	Debug    bool   `koanf:"debug"`
}

type slackDriver struct {
	cfg           *config
	l             *slog.Logger
	driveCallback driver.DriveFunc
	client        *socketmode.Client
}

type logger struct{ *slog.Logger }

func (l logger) Output(_ int, msg string) error { l.Info(msg); return nil }

func New() fx.Option {
	return app.Module[config](
		"slackdriver",
		fx.Invoke(func(lc fx.Lifecycle, l *slog.Logger, cfg *config, drive driver.DriveFunc) {
			if cfg.BotToken == "" && cfg.AppToken == "" {
				l.Warn("no slack bot token or app token supplied, not connecting to slack")
				return
			}

			lc.Append(
				fx.Hook{
					OnStart: func(ctx context.Context) error {
						api := slack.New(
							cfg.BotToken,
							slack.OptionDebug(cfg.Debug),
							slack.OptionLog(logger{l}),
							slack.OptionAppLevelToken(cfg.AppToken),
						)

						resp, err := api.AuthTest()
						if err != nil {
							l.Error("slack auth test error", "err", err)
							return err
						}

						l.Info("connected to slack", "info", resp)

						client := socketmode.New(
							api,
							socketmode.OptionDebug(cfg.Debug),
							socketmode.OptionLog(logger{l}),
						)

						driver := &slackDriver{
							cfg:           cfg,
							l:             l,
							driveCallback: drive,
							client:        client,
						}

						go driver.run()

						go driver.serve()

						return nil
					},
				},
			)
		}),
	)
}

func (d *slackDriver) run() {
	if err := d.client.Run(); err != nil {
		d.l.Error("slack run error", "err", err)
	}
}

func (d *slackDriver) serve() {
	for evt := range d.client.Events {
		d.processEvent(evt)
	}
}

func (d *slackDriver) drive(kind socketmode.EventType, data any) {
	_ = d.driveCallback(context.Background(), "slack", map[string]any{
		"type": kind,
		"data": data,
	})
}

func (d *slackDriver) processEvent(evt socketmode.Event) {
	l, client := d.l, d.client

	switch evt.Type {
	case socketmode.EventTypeConnecting:
		l.Info("connecting to slack with socket mode...")
	case socketmode.EventTypeConnectionError:
		l.Info("connection failed. retrying later...")
	case socketmode.EventTypeConnected:
		l.Info("connected to slack with socket mode.")
	case socketmode.EventTypeEventsAPI:
		apiEvent, ok := evt.Data.(slackevents.EventsAPIEvent)
		if !ok {
			l.Debug("slack event ignored", "event", evt)
			return
		}

		switch apiEvent.Type {
		case slackevents.CallbackEvent:
			l.Info("slack callback event received", "event", apiEvent)
			d.drive(evt.Type, apiEvent.InnerEvent)
		default:
			l.Debug("unsupported Events API event received")
		}

		client.Ack(*evt.Request)
	case socketmode.EventTypeInteractive:
		callback, ok := evt.Data.(slack.InteractionCallback)
		if !ok {
			l.Debug("slack interactive event ignored", "event", evt)
			return
		}

		l.Info("slack interaction received", "callback", callback)

		d.drive(evt.Type, callback)

		client.Ack(*evt.Request)
	case socketmode.EventTypeSlashCommand:
		cmd, ok := evt.Data.(slack.SlashCommand)
		if !ok {
			l.Debug("slack slash command event ignored", "event", evt)
			return
		}

		l.Info("slash command received", "cmd", cmd)

		d.drive(evt.Type, cmd)

		client.Ack(*evt.Request)
	default:
		l.Error("unexpected event type received", "type", evt.Type)
	}
}
