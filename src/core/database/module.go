package database

import (
	"database/sql"

	"github.com/frostgrove/vv/crud"
	"github.com/frostgrove/vv/crud/adapter/crudsql/crudsqlfx"
	"github.com/frostgrove/vv/health"
	"github.com/frostgrove/vv/health/healthfx"
	vvotel "github.com/frostgrove/vv/otel"
	_ "github.com/jackc/pgx/v5/stdlib" // the database/sql driver behind engine: postgres (D-057: the application picks it)
	"go.uber.org/fx"

	"github.com/frostgrove/vv-template/src/config"
	"github.com/frostgrove/vv-template/src/core/observability"
)

const tracingLayer = "otel.source"

// Register opens the application's database/sql pool once (crudsqlfx pings it
// at start and closes it on stop) and exposes it as the crud.Source every
// generated repository binds to. No repository depends on a driver or an ORM.
func Register(configuration *config.Config) fx.Option {
	options := []fx.Option{
		fx.Provide(healthfx.AsCheck(Check)),
	}
	var layers []crudsqlfx.Option
	if configuration.Observability.Otel.Enabled {
		layers = append(layers, crudsqlfx.Layers(tracingLayer))
		options = append(options, fx.Provide(crudsqlfx.AsWrapping(tracing)))
	}
	options = append(options, crudsqlfx.Module(&configuration.Infra.Db, layers...))
	return fx.Options(options...)
}

func Check(db *sql.DB) health.Contribution {
	return health.Contribution{
		Name:       "database.primary",
		Code:       "database",
		Importance: health.Required,
		Probe:      health.ProbeFunc(db.PingContext),
	}
}

func tracing(telemetry *observability.Telemetry) crudsqlfx.Wrapping {
	return crudsqlfx.Wrapping{
		Name: tracingLayer,
		Wrap: func(next crud.Source) crud.Source { return vvotel.Source(telemetry.Frost, next) },
	}
}
