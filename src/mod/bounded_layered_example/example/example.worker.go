package example

import (
	"context"
	"log/slog"
	"time"

	"github.com/frostgrove/vv/runtime"
)

// newOpenReporter is background work: a runtime.Runner the supervisor starts,
// drains and stops with the process, and whose failure stops the process too.
// Never a bare goroutine (D-092). Only a profile with the worker role runs it.
func newOpenReporter(api API, logger *slog.Logger) (runtime.Runner, error) {
	return runtime.NewPeriodic(runtime.PeriodicSpec{
		Name:      "example.open-reporter",
		Interval:  time.Minute,
		Timeout:   10 * time.Second,
		Immediate: true,
		Logger:    logger,
		Pass: func(ctx context.Context) error {
			open, err := api.CountOpen(ctx)
			if err != nil {
				return err
			}
			logger.InfoContext(ctx, "open examples", slog.Int64("count", open))
			return nil
		},
	})
}
