package boot

import (
	"github.com/frostgrove/vv-template/src/app/product"
	"go.uber.org/fx"
)

func appOptions() fx.Option {
	return fx.Options(
		product.Module(),
	)
}
