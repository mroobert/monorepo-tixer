package tixer

import (
	"github.com/google/uuid"
)

// Ticket represents a ticket that can be purchased.
type Ticket struct {
	ID      uuid.UUID
	OrderID uuid.UUID
	RefCode int64
	Title   string
	Price   int64
	UserId  string
	Version int32
}
