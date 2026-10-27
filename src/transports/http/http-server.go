package httpserver

import (
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/frostgrove/vv/app/http/appfiber"
	"github.com/frostgrove/vv/crud/http/crudfiber"
	"github.com/frostgrove/vv/port"
	fiberotel "github.com/gofiber/contrib/v3/otel"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/gofiber/fiber/v3/middleware/recover"

	"github.com/frostgrove/vv-template/src/config"
	coreauth "github.com/frostgrove/vv-template/src/core/auth"
	"github.com/frostgrove/vv-template/src/core/observability"
)

// NewFiberApp is the one *fiber.App of the process. What is mounted on it comes
// from the modules (appfiber.AsRoute); appfiber.Serving mounts it, compares the
// routing table with every declared access policy, listens, and shuts it down
// gracefully when fx stops.
func NewFiberApp(configuration *config.Config, logger *slog.Logger, telemetry *observability.Telemetry) *fiber.App {
	app := fiber.New(fiber.Config{
		AppName: configuration.App.Name,
		// Router refusals (404, 405, 413) and handler errors share the CRUD envelope.
		ErrorHandler: crudfiber.ErrorHandler(),
	})

	app.Use(recover.New())
	if telemetry.Enabled {
		app.Use(fiberotel.Middleware(
			fiberotel.WithTracerProvider(telemetry.TracerProvider),
			fiberotel.WithMeterProvider(telemetry.MeterProvider),
			fiberotel.WithPropagators(telemetry.Propagator),
			fiberotel.WithSpanNameFormatter(func(c fiber.Ctx) string {
				return c.Method() + " " + c.Route().Path
			}),
		))
	}
	app.Use(requestLogging(logger))
	if origins := configuration.Transports.Fiber.AllowOrigins; len(origins) > 0 {
		app.Use(cors.New(cors.Config{
			AllowOrigins: origins,
			AllowHeaders: []string{fiber.HeaderOrigin, fiber.HeaderContentType, coreauth.HeaderAPIKey},
			AllowMethods: []string{
				fiber.MethodGet, fiber.MethodPost, fiber.MethodHead,
				fiber.MethodPut, fiber.MethodDelete, fiber.MethodPatch,
			},
			AllowCredentials: true,
			MaxAge:           3600,
		}))
	}
	return app
}

// requestLogging puts the process logger into every request context, which is
// where framework code looks for it (port.Logger, D-062), and writes one line
// per request. The error handler runs here so the line carries the final status.
func requestLogging(logger *slog.Logger) fiber.Handler {
	return func(c fiber.Ctx) error {
		started := time.Now()
		c.SetContext(port.WithLogger(c.Context(), logger))

		if err := c.Next(); err != nil {
			if handled := c.App().Config().ErrorHandler(c, err); handled != nil {
				_ = c.SendStatus(http.StatusInternalServerError)
			}
		}

		level := slog.LevelInfo
		if strings.Contains(c.Path(), appfiber.DefaultHealthPath+"/") {
			level = slog.LevelDebug
		}
		logger.Log(c.Context(), level, "http request",
			slog.String("method", c.Method()),
			slog.String("route", c.Route().Path),
			slog.String("path", c.Path()),
			slog.Int("status", c.Response().StatusCode()),
			slog.Duration("duration", time.Since(started)),
		)
		return nil
	}
}

// NewRootRoutes answers outside the API prefix. Every route, here too, says
// what reaching it takes; Public needs the reason in writing (D-073, D-100).
func NewRootRoutes(app *fiber.App, configuration *config.Config) (appfiber.Route, error) {
	return appfiber.RootRoutes(app).
		GET("/", appfiber.Public("names the service to whoever opens its address"),
			func(c fiber.Ctx) error {
				return c.JSON(fiber.Map{"name": configuration.App.Name, "timestamp": time.Now()})
			}).
		GET("/favicon.ico", appfiber.Public("browsers ask for it unprompted and get nothing"),
			func(c fiber.Ctx) error { return c.SendStatus(http.StatusNoContent) }).
		Route()
}
