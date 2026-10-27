package product

import (
	"github.com/frostgrove/vv/app/appfx"
	"github.com/frostgrove/vv/app/http/appfiber"
	"github.com/frostgrove/vv/app/module"
)

// Module is everything this bounded context contributes, filed by role. The
// profile decides what of it a process wires: an API replica mounts the routes,
// the seed command runs the seeder, every process gets the repositories.
func Module() module.Definition {
	return module.New("product").
		Order(100).
		Provide(NewProductRepository, NewProductImageRepository).
		Routes(appfiber.AsRoute(NewHttpHandler)).
		Seeders(appfx.AsSeeder(NewDemoSeeder)).
		MustBuild()
}
