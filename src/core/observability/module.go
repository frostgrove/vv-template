package observability

import (
	"log/slog"

	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
)

func Register() fx.Option {
	return fx.Options(
		fx.Provide(NewLogger, NewTelemetry),
		fx.WithLogger(func(logger *slog.Logger) fxevent.Logger {
			// The container's own wiring lines are debug; its errors stay errors.
			events := &fxevent.SlogLogger{Logger: logger.With(slog.String("component", "fx"))}
			events.UseLogLevel(slog.LevelDebug)
			return events
		}),
	)
}
