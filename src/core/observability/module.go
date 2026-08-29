package observability

import (
	"log/slog"

	"github.com/frostgrove/vv-template/src/config"
	"go.uber.org/fx"
)

func Register() fx.Option {
	return fx.Options(
		fx.Provide(
			func(configuration *config.Config) *slog.Logger {
				return MustCreateStdLogger(configuration)
			},
		),
	)
}
