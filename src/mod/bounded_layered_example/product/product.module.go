package product

import (
	"github.com/frostgrove/vv/crud/decorators/specs"
	"github.com/frostgrove/vv/crud/http/crudfiber"
	"github.com/gofiber/fiber/v3"
	"go.uber.org/fx"
)

func Module() fx.Option {
	return fx.Options(
		fx.Provide(NewProductRepository),

		fx.Invoke(
			func(repo *ProductRepo, router *fiber.App) {
				var exec = specs.Executor(repo)
				router.Use("/products", crudfiber.New(exec).Routes())
			},
		),
	)
}
