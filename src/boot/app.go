package boot

import (
	"github.com/frostgrove/vv/app/module"

	"github.com/frostgrove/vv-template/src/mod/bounded_layered_example/example"
	"github.com/frostgrove/vv-template/src/mod/bounded_layered_example/product"
)

// Modules is the one list of bounded contexts this service is made of. A
// module says what it contributes and to which role (module.Definition); the
// profile in config decides which of those roles this process wires (D-106).
func Modules() module.Catalog {
	return module.MustCatalog(
		product.Module(),
		example.Module(),
	)
}
