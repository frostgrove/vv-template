package auth

import (
	"github.com/frostgrove/vv/app/http/appfiber"
	"github.com/frostgrove/vv/auth"
	"github.com/frostgrove/vv/auth/apikey"
	"go.uber.org/fx"

	"github.com/frostgrove/vv-template/src/config"
)

// HeaderAPIKey carries the development API key: X-Api-Key: <key>.
const HeaderAPIKey = "X-Api-Key"

// Register authenticates every API request that presents a key and lets the
// others through anonymously. Authenticating is all it does: what a caller may
// reach is the policy each route declares (appfiber.Requires / Authenticated /
// Public, D-073 and D-100), and that check answers 401 to an anonymous caller.
//
// The static key list is for development. Swap the authenticator for a hashed
// key store (apikey.StoreFunc), auth/authjwt or auth/access in production; the
// routes do not change.
func Register() fx.Option {
	return fx.Provide(
		NewGuard,
		appfiber.AsMiddleware(func(guard *auth.Guard) appfiber.Middleware {
			return appfiber.Guarding("auth.apikey", guard)
		}),
	)
}

func NewGuard(configuration *config.Config) *auth.Guard {
	issued := make(map[string]auth.Principal, len(configuration.Transports.Fiber.APIKeys))
	for _, key := range configuration.Transports.Fiber.APIKeys {
		permissions := make([]auth.Permission, 0, len(key.Permissions))
		for _, permission := range key.Permissions {
			permissions = append(permissions, auth.Permission(permission))
		}
		issued[key.Key] = auth.Claims{Sub: key.Subject, Permissions: permissions}
	}
	return auth.NewGuard(apikey.New(apikey.Static(issued)), apikey.Header(HeaderAPIKey), auth.Optional())
}
