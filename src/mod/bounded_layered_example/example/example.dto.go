package example

import (
	"time"

	"github.com/google/uuid"
)

// exampleDto is what leaves the bounded context: the model stays inside.
type exampleDto struct {
	Id        uuid.UUID `json:"id"`
	Title     string    `json:"title"`
	Done      bool      `json:"done"`
	CreatedAt time.Time `json:"createdAt"`
}

func toDto(model Example) exampleDto {
	return exampleDto{Id: model.Id, Title: model.Title, Done: model.Done, CreatedAt: model.CreatedAt}
}
