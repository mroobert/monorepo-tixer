package tixer

import (
	"time"
)

// Ticket represents a ticket that can be purchased.
type Ticket struct {
	ID        int64
	PublicID  PublicID
	Title     string
	Price     int64
	Version   int32
	CreatedAt time.Time
	UpdatedAt time.Time
}

// Validate checks ticket's fields to ensure that the basic business rules are met.
// It returns a boolean indicating if the ticket is valid and a map of errors if it's not.
func (t Ticket) Validate() (bool, map[string]string) {
	errors := make(map[string]string)

	if len(t.Title) > 150 {
		errors["title"] = "must not be more than 150 characters long"
	}

	if t.Price <= 0 {
		errors["price"] = "must be greater than 0"
	}

	if t.Price > 50_000 {
		errors["price"] = "must be lower than 500 euros"
	}

	if len(errors) > 0 {
		return false, errors
	}

	return true, nil
}
