package example

import (
	"time"

	"github.com/google/uuid"
)

// Example is the persistence model. `make generate` writes ExampleUpdate,
// the Example_ metamodel and NewExampleRepository into vv_gen.go next to it.
type Example struct {
	Id        uuid.UUID `db:"id,pk,auto"`
	Title     string
	Done      bool
	CreatedAt time.Time `db:"created_at,generated"`
}
