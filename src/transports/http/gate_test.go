package httpserver

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/frostgrove/vv/app/http/appfiber"
	"github.com/frostgrove/vv/auth"
	"github.com/frostgrove/vv/auth/http/authhttp"
	"github.com/gofiber/fiber/v3"
	"go.uber.org/fx"

	"github.com/frostgrove/vv-template/src/config"
	coreauth "github.com/frostgrove/vv-template/src/core/auth"
	"github.com/frostgrove/vv-template/src/core/observability"
)

const permSecret auth.Permission = "secret:read"

func testConfig() *config.Config {
	configuration := &config.Config{}
	configuration.App = config.App{Name: "gate-test", Env: "test", Profile: "serving"}
	configuration.Transports.Fiber = config.FiberConfig{
		Addr:   "127.0.0.1:0",
		Prefix: "/api",
		APIKeys: []config.APIKey{
			{Key: "holder", Subject: "holder", Permissions: []string{string(permSecret)}},
			{Key: "nobody", Subject: "nobody"},
		},
	}
	configuration.Observability.Logger = config.LoggerConfig{Level: "error", Format: "text"}
	return configuration
}

// start builds the same HTTP surface the service runs — the Fiber app, the
// guard, the health and root routes, appfiber.Serving — plus the extra routes
// a test contributes, and starts it.
func start(t *testing.T, routes ...any) (*fiber.App, error) {
	t.Helper()
	configuration := testConfig()
	var app *fiber.App
	options := []fx.Option{
		fx.NopLogger,
		fx.Supply(configuration),
		observability.Register(),
		coreauth.Register(),
		Register(configuration),
		fx.Populate(&app),
	}
	for _, route := range routes {
		options = append(options, fx.Provide(appfiber.AsRoute(route)))
	}
	application := fx.New(options...)
	if err := application.Err(); err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := application.Start(ctx); err != nil {
		return nil, err
	}
	t.Cleanup(func() { _ = application.Stop(context.Background()) })
	return app, nil
}

type undeclared struct{}

func (undeclared) Mount(router fiber.Router) {
	router.Get("/backdoor", func(c fiber.Ctx) error { return c.SendString("open") })
}

func (undeclared) Access() []authhttp.Endpoint { return nil }

type unchecked struct{}

func (unchecked) Mount(router fiber.Router) {
	router.Get("/claimed", func(c fiber.Ctx) error { return c.SendString("open") })
}

func (unchecked) Access() []authhttp.Endpoint {
	return []authhttp.Endpoint{authhttp.Requires(http.MethodGet, "/claimed", permSecret)}
}

func TestTheServiceStartsWithEveryRouteDeclared(t *testing.T) {
	if _, err := start(t); err != nil {
		t.Fatalf("the plain HTTP surface did not start, so the refusals below prove nothing: %v", err)
	}
}

func TestARouteNobodyDeclaredRefusesTheStart(t *testing.T) {
	_, err := start(t, func() appfiber.Route { return undeclared{} })
	if err == nil {
		t.Fatal("a route mounted without an access declaration started; the gate let it through")
	}
}

func TestADeclarationNothingEnforcesRefusesTheStart(t *testing.T) {
	_, err := start(t, func() appfiber.Route { return unchecked{} })
	if !errors.Is(err, appfiber.ErrUnchecked) {
		t.Fatalf("a hand-written Requires with no check in front of it started (%v)", err)
	}
}

func TestADeclaredPermissionIsEnforced(t *testing.T) {
	app, err := start(t, func() (appfiber.Route, error) {
		return appfiber.Routes("/secret").
			GET("", appfiber.Requires(permSecret), func(c fiber.Ctx) error { return c.SendString("ok") }).
			Route()
	})
	if err != nil {
		t.Fatal(err)
	}

	for _, case_ := range []struct {
		key  string
		want int
	}{
		{"", http.StatusUnauthorized},
		{"wrong", http.StatusUnauthorized},
		{"nobody", http.StatusForbidden},
		{"holder", http.StatusOK},
	} {
		request := httptest.NewRequest(http.MethodGet, "/api/secret", nil)
		if case_.key != "" {
			request.Header.Set(coreauth.HeaderAPIKey, case_.key)
		}
		response, err := app.Test(request)
		if err != nil {
			t.Fatal(err)
		}
		if response.StatusCode != case_.want {
			t.Errorf("key %q got %d, want %d", case_.key, response.StatusCode, case_.want)
		}
	}
}

func TestHealthProbesArePublic(t *testing.T) {
	app, err := start(t)
	if err != nil {
		t.Fatal(err)
	}
	for _, path := range []string{"/api/health/live", "/api/health/ready"} {
		response, err := app.Test(httptest.NewRequest(http.MethodGet, path, nil))
		if err != nil {
			t.Fatal(err)
		}
		if response.StatusCode != http.StatusOK {
			t.Errorf("%s answered %d to an anonymous probe", path, response.StatusCode)
		}
	}
	response, err := app.Test(httptest.NewRequest(http.MethodGet, "/api/health/detail", nil))
	if err != nil {
		t.Fatal(err)
	}
	if response.StatusCode != http.StatusUnauthorized {
		t.Errorf("the operator health view answered %d to an anonymous caller", response.StatusCode)
	}
}
