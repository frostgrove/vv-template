package boot

import (
	"github.com/frostgrove/vv-template/src/config"
	"github.com/frostgrove/vv-template/src/core/database"
	"github.com/frostgrove/vv-template/src/core/observability"
	http_server_module "github.com/frostgrove/vv-template/src/transports/http"
	"github.com/frostgrove/vv/utils/vvcfg"
	"go.uber.org/fx"
)

func coreOptions() fx.Option {
	return fx.Options(
		fx.Provide(
			vvcfg.MustLoad[config.Config],
		),

		observability.Register(),

		database.Register(),

		http_server_module.Register(),
	)
}
