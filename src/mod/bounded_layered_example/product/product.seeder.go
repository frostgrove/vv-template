package product

import (
	"context"

	"github.com/frostgrove/vv/app"
	"github.com/frostgrove/vv/crud"
	"github.com/frostgrove/vv/crud/decorators/specs"
)

// NewDemoSeeder fills an empty development database with a few products.
// It runs only in dev and is safe to run twice: it looks before it writes.
func NewDemoSeeder(repository *ProductRepo) app.Seeder {
	return app.Seeder{
		Name:  "product.demo",
		Order: 100,
		Envs:  []string{"dev"},
		Run: func(ctx context.Context) error {
			for _, name := range []string{"Hex bolt", "Hex nut", "Washer"} {
				exists, err := repository.Exists(ctx, crud.Where(specs.Predicate(Product_.Name.Eq(name))))
				if err != nil {
					return err
				}
				if exists {
					continue
				}
				if _, err := repository.Save(ctx, &Product{Name: name}); err != nil {
					return err
				}
			}
			return nil
		},
	}
}
