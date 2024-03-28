package tixer

import (
	"errors"
	"fmt"
	"strings"
)

const (
	PublicIDAlphabet = "0123456789abcdefghijklmnopqrstuvwxyz"
	PublicIDLength   = 12
)

// PublicID represents a public ID value for human-readable use case.
type PublicID string

// ValidatePublicID checks if a given value is a valid PublicID.
func ValidatePublicID(value string) error {
	if value == "" {
		return errors.New("value cannot be blank")
	}

	if len(value) != PublicIDLength {
		return fmt.Errorf("value should be %d characters long", PublicIDLength)
	}

	if strings.Trim(value, PublicIDAlphabet) != "" {
		return errors.New("value has invalid characters")
	}

	return nil
}
