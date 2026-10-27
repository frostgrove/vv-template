package httpserver

import (
	"github.com/frostgrove/vv/app/http/appfiber"
	"github.com/frostgrove/vv/auth"
	"github.com/frostgrove/vv/health"
	"github.com/frostgrove/vv/health/healthfx"
	"github.com/gofiber/fiber/v3"
	"go.uber.org/fx"

	"github.com/frostgrove/vv-template/src/config"
)

// PermHealthDetail opens GET <prefix>/health/detail, the operator view naming
// every dependency. Live and ready are public: a probe has no account.
const PermHealthDetail auth.Permission = "health:detail"

func Register(configuration *config.Config) fx.Option {
	return fx.Options(
		fx.Provide(
			NewFiberApp,
			appfiber.AsRoute(NewRootRoutes),
			appfiber.AsRoute(NewHealthRoutes),
			// Every Requires/Authenticated route in this template is mounted
			// through the registrar, which puts the check in front of the
			// handler. A hand-written declaration nothing enforces fails the
			// start instead of being logged.
			func() appfiber.UncheckedRule { return appfiber.RefusingUnchecked },
		),
		healthfx.Auto(),
		appfiber.Serving(appfiber.Spec{
			Prefix: configuration.Transports.Fiber.Prefix,
			Addr:   configuration.Transports.Fiber.Addr,
			// "the api is listening" is logged instead of Fiber's banner.
			Listen: fiber.ListenConfig{DisableStartupMessage: true},
		}),
	)
}

func NewHealthRoutes(registry *health.Registry) (appfiber.Route, error) {
	return appfiber.Health(appfiber.HealthSpec{
		Registry: registry,
		Operator: []auth.Permission{PermHealthDetail},
	})
}
