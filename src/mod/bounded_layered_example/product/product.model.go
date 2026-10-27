package product

import (
	"time"

	"github.com/google/uuid"
)

// Product is an ordinary Go struct: vv derives the "products" table, the
// snake_case columns and the Id primary key from it. Tags are only for the
// exceptions: the database generates the key and both timestamps
// (see resources/migrations), so a client can neither choose nor forge them.
// The json tags are the wire shape crudfiber reads and writes.
type Product struct {
	Id          uuid.UUID `db:"id,pk,auto" json:"id"`
	Name        string    `json:"name"`
	Description *string   `json:"description"`

	CreatedAt time.Time `db:"created_at,generated" json:"createdAt"`
	UpdatedAt time.Time `db:"updated_at,generated" json:"updatedAt"`
}

type ProductImage struct {
	Id   uuid.UUID `db:"id,pk,auto" json:"id"`
	Path string    `json:"path"`
}
