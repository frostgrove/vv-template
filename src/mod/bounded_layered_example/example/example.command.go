package example

import "github.com/google/uuid"

// Commands are the values a use case takes. They hold what the caller asked
// for, already decoded, and nothing about the transport it came over.

type createExampleCommand struct {
	Title string `json:"title"`
}

type completeExampleCommand struct {
	Id uuid.UUID
}
