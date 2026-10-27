package boot

import (
	"log/slog"

	"github.com/frostgrove/vv/runtime/runtimefx"
	"github.com/frostgrove/vv/utils/vvcfg"
	"go.uber.org/fx"

	"github.com/frostgrove/vv-template/src/config"
	"github.com/frostgrove/vv-template/src/core/auth"
	"github.com/frostgrove/vv-template/src/core/database"
	"github.com/frostgrove/vv-template/src/core/observability"
)

// coreOptions is the infrastructure every process of this service shares,
// whatever its profile: configuration, logging and telemetry, the database
// and the supervisor of background runners (D-092).
func coreOptions(configuration *config.Config, report *vvcfg.Report) fx.Option {
	return fx.Options(
		fx.Supply(configuration),
		observability.Register(),
		fx.Invoke(func(logger *slog.Logger) {
			logger.Info("configuration loaded",
				slog.String("path", report.Path), slog.String("from", report.PathOrigin.String()))
		}),
		database.Register(configuration),
		auth.Register(),
		runtimefx.Auto(),
	)
}
