package product

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"go.uber.org/fx"
)

func Module() fx.Option {
	return fx.Options(
		fx.Provide(NewProductRepository),

		fx.Invoke(
			func(repo *ProductRepo) {
				var saved = repo.Save(context.Background(), &Product{
					Id:   uuid.Must(uuid.NewV7()),
					Name: "test",
				})

				fmt.Printf("%v", saved)

				got, err := repo.Get(context.Background())
				fmt.Printf("%v, %v", got, err)
			},
		),
	)
}
