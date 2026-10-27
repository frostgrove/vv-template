package observability

import (
	"io"
	"log/slog"
	"os"

	vvotel "github.com/frostgrove/vv/otel"

	"github.com/frostgrove/vv-template/src/config"
)

// NewLogger is the process logger. vvotel.TraceHandler adds trace_id/span_id
// from the record's context, so a line logged with a request context is
// correlated with its span whenever tracing is on. It also becomes
// slog.Default, which is where framework code falls back when a context
// carries no logger (port.Logger, D-062).
func NewLogger(configuration *config.Config) *slog.Logger {
	logger := newLogger(os.Stdout, configuration)
	slog.SetDefault(logger)
	return logger
}

func newLogger(out io.Writer, configuration *config.Config) *slog.Logger {
	options := &slog.HandlerOptions{Level: level(configuration.Observability.Logger.Level)}
	var handler slog.Handler
	if configuration.Observability.Logger.Format == "json" {
		handler = slog.NewJSONHandler(out, options)
	} else {
		handler = slog.NewTextHandler(out, options)
	}
	return slog.New(vvotel.TraceHandler(handler)).With(slog.String("service", configuration.App.Name))
}

func level(name string) slog.Level {
	switch name {
	case "debug":
		return slog.LevelDebug
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	}
	return slog.LevelInfo
}
