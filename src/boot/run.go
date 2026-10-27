package boot

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/frostgrove/vv/app"
	"github.com/frostgrove/vv/app/appfx"
	"github.com/frostgrove/vv/app/module"
	"go.uber.org/fx"

	"github.com/frostgrove/vv-template/src/config"
	httpserver "github.com/frostgrove/vv-template/src/transports/http"
)

// Run starts the service with the profile the configuration names and blocks
// until SIGINT/SIGTERM, then stops everything in reverse order: the HTTP
// server drains, runners stop, the pool closes, telemetry flushes.
func Run() {
	configuration, report := config.MustLoad()
	profile, _ := configuration.App.ModuleProfile()

	fx.New(
		coreOptions(configuration, report),
		appfx.Options(Modules(), profile),
		httpserver.Register(configuration),
	).Run()
	slog.Info("stopped")
}

// Seed runs every seeder the modules contribute for the configured env, then
// exits. Seeders are written to be safe to run twice.
func Seed() {
	configuration, report := config.MustLoad()

	var runner *app.Runner
	application := fx.New(
		coreOptions(configuration, report),
		appfx.Options(Modules(), module.Seeding),
		appfx.Seeding(app.Seeding{Env: configuration.App.Env, Known: config.Envs}),
		fx.Populate(&runner),
	)
	if err := application.Err(); err != nil {
		fail(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	if err := application.Start(ctx); err != nil {
		fail(err)
	}
	seeded := runner.Run(ctx)
	stopped := application.Stop(ctx)
	if seeded != nil {
		fail(seeded)
	}
	if stopped != nil {
		fail(stopped)
	}
}

func fail(err error) {
	fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}
