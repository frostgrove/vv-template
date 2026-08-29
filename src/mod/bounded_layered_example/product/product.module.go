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
			func(productRepository *ProductRepo, router *fiber.App) {
				var executor = specs.Executor(productRepository)
				router.Use("/products", crudfiber.New(executor).Routes())
			},
		),
	)
}
