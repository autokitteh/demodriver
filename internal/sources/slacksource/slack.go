package slacksource

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
	BotToken string `koanf:"bottoken"`
	AppToken string `koanf:"apptoken"`
	Debug    bool   `koanf:"debug"`
}

type slackSource struct {
	cfg           *config
	l             *slog.Logger
	driveCallback driver.DriveFunc
	client        *socketmode.Client
}

type logger struct{ *slog.Logger }

func (l logger) Output(_ int, msg string) error { l.Info(msg); return nil }

func New() fx.Option {
	return app.Module[config](
		"slacksource",
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

						source := &slackSource{
							cfg:           cfg,
							l:             l,
							driveCallback: drive,
							client:        client,
						}

						h := socketmode.NewSocketmodeHandler(client)
						h.HandleDefault(source.handleEvent)

						go func() {
							if err := h.RunEventLoop(); err != nil {
								l.Error("slack event loop error", "err", err)
								return
							}
						}()

						return nil
					},
				},
			)
		}),
	)
}

func (d *slackSource) drive(kind socketmode.EventType, data any) {
	_ = d.driveCallback(context.Background(), "slack", map[string]any{
		"type": kind,
		"data": data,
	})
}

func (d *slackSource) handleEvent(evt *socketmode.Event, client *socketmode.Client) {
	l := d.l

	switch evt.Type {
	case socketmode.EventTypeHello:
		l.Info("hello received from slack")

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
			l.Warn("unsupported Events API event received")
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
