package http_server_module

import (
	"context"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/frostgrove/vv-template/src/config"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	fiberLogger "github.com/gofiber/fiber/v3/middleware/logger"
	"github.com/gofiber/fiber/v3/middleware/recover"
	"go.uber.org/fx"
)

func MustInit(
	configuration *config.Config,
	logger *slog.Logger,
) *fiber.App {
	app := fiber.New(fiber.Config{})

	app.Use(fiberLogger.New())
	app.Use(recover.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: configuration.Transports.Fiber.AllowOrigins,
		AllowHeaders: []string{"Origin", "Content-Type", "X-Access-Token"},
		AllowMethods: []string{
			fiber.MethodGet,
			fiber.MethodPost,
			fiber.MethodHead,
			fiber.MethodPut,
			fiber.MethodDelete,
			fiber.MethodPatch,
		},

		AllowCredentials: true,
		ExposeHeaders:    []string{"Set-Cookie"},
		MaxAge:           3600,
	}))

	app.Use(func(fiberContext fiber.Ctx) error {
		fiberContext.Response().Header.Del("Vary")

		return fiberContext.Next()
	})
	app.Get("/", func(fiberContext fiber.Ctx) error {
		return fiberContext.JSON(map[string]any{"timestamp": time.Now(), "name": configuration.App.Name})
	})

	app.Get("/favicon.ico", func(fiberContext fiber.Ctx) error {
		return fiberContext.SendStatus(http.StatusNoContent)
	})

	return app

}

func MustStart(lifecycle fx.Lifecycle, app *fiber.App, configuration *config.Config, logger *slog.Logger) {
	lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go func() {
				var err = app.Listen(":"+strconv.Itoa(configuration.Transports.Fiber.Port), fiber.ListenConfig{EnablePrefork: configuration.Transports.Fiber.Prefork})

				if err != nil && err != http.ErrServerClosed {
					logger.Error("Error when starting server", slog.Any("err", err.Error()))
				}
			}()
			return nil
		},
	})
}
