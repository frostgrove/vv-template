package example

import "go.uber.org/fx"

func Module() fx.Option {
	return fx.Options(
		fx.Provide(
		// Register your modules here
		),

		fx.Invoke(
		// Register your commands here
		),
	)
}
