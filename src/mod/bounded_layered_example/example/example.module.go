package example

import (
	"github.com/frostgrove/vv/app/http/appfiber"
	"github.com/frostgrove/vv/app/module"
	"github.com/frostgrove/vv/runtime/runtimefx"
)

func Module() module.Definition {
	return module.New("example").
		Order(200).
		Provide(NewExampleRepository, newRepo, newService, NewAPI).
		Routes(appfiber.AsRoute(newHttpHandler)).
		Workers(runtimefx.AsRunner(newOpenReporter)).
		MustBuild()
}
