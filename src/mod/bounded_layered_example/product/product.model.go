package product

import (
	"time"

	"github.com/google/uuid"
)

type Product struct {
	Id          uuid.UUID
	Name        string
	Description *string

	CreatedAt time.Time
	UpdatedAt time.Time
}

type ProductImage struct {
	Id   uuid.UUID
	Path string
}
