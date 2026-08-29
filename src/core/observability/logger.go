package observability

import (
	"log/slog"

	"github.com/frostgrove/vv-template/src/config"
	"go.opentelemetry.io/contrib/bridges/otelslog"
	"go.opentelemetry.io/otel/exporters/stdout/stdoutlog"
	"go.opentelemetry.io/otel/sdk/log"
)

func MustCreateStdLogger(configuration *config.Config) *slog.Logger {
	logger, err := NewStdLogger(configuration)
	if err != nil {
		panic(err)
	}

	return logger
}

func NewStdLogger(configuration *config.Config) (*slog.Logger, error) {
	exporter, err := stdoutlog.New()
	if err != nil {
		return nil, err
	}

	provider := log.NewLoggerProvider(
		log.WithProcessor(log.NewBatchProcessor(exporter)),
	)

	logger := otelslog.NewLogger("test", otelslog.WithLoggerProvider(provider))
	return logger, nil
}
