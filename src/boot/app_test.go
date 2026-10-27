package boot

import (
	"testing"

	"github.com/frostgrove/vv/app/module"
)

func TestEveryProfileOverTheCatalogIsUsable(t *testing.T) {
	for _, profile := range []module.Profile{module.Complete, module.Serving, module.Working, module.Seeding} {
		if diagnosis := module.Doctor(Modules(), profile); !diagnosis.OK() {
			t.Errorf("the %s profile would not start:\n%s", profile.Name, diagnosis)
		}
	}
}
