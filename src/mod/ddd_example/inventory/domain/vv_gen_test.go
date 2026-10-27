package domain

import "testing"

// Importing the package runs the generated init: every metamodel is checked
// against its model and every update DTO against its columns. A model changed
// without `make generate` fails here rather than at the first request.
func TestTheGeneratedMetamodelsMatchTheModels(t *testing.T) {
	if StockItem_.Name.Path().Name == "" {
		t.Fatal("the StockItem metamodel has an unnamed attribute")
	}
}
