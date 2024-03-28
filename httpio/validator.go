package httpio

import (
	"net/url"
	"strconv"
)

// validator represents a data parser & validator for the http request payload.
type validator struct {
	errors map[string]string
}

// newValidator creates a new validator instance.
func newValidator() *validator {
	return &validator{errors: make(map[string]string)}
}

// valid returns true if the validator has no errors.
func (v *validator) valid() bool {
	return len(v.errors) == 0
}

// addError adds an error message to the validator.
func (v *validator) addError(key, message string) {
	if _, exists := v.errors[key]; !exists {
		v.errors[key] = message
	}
}

// check adds an error message to the validator if the condition is not met.
func (v *validator) check(ok bool, key, message string) {
	if !ok {
		v.addError(key, message)
	}
}

// readString reads a string value from the url query string.
func (v *validator) readString(qs url.Values, key string, defaultValue string) string {
	value := qs.Get(key)
	if value == "" {
		return defaultValue
	}
	return value
}

// readInt reads an integer value from the url query string.
func (v *validator) readInt(qs url.Values, key string, defaultValue int) int {
	value := qs.Get(key)
	if value == "" {
		return defaultValue
	}

	intValue, err := strconv.Atoi(value)
	if err != nil {
		v.addError(key, "must be an integer value")
		return defaultValue
	}
	return intValue
}

// validateCreateTicketBody validates the create ticket request body.
func (v *validator) validateCreateTicketRequestBody(body createTicketRequestBody) {
	v.check(body.Title != "", "title", "must be provided")
	v.check(body.Price > 0, "price", "must be provided")
}

// validateTicketUrlValues validates the url query string parameters used for reading multiple rows of tickets.
func (v *validator) validateTicketUrlValues(qs url.Values, sortSafeList []string) ticketUrlQs {
	title := v.readString(qs, "title", "")
	page := v.readInt(qs, "page", 1)
	pageSize := v.readInt(qs, "pageSize", 10)
	sort := v.readString(qs, "sort", "id")

	v.check(page <= 1000, "page", "must be a maximum of 1000")
	v.check(pageSize <= 25, "page_size", "must be a maximum of 25")

	v.check(permittedValue(sort, sortSafeList...), "sort", "invalid sort value")

	return ticketUrlQs{
		title:    title,
		page:     page,
		pageSize: pageSize,
		sort:     sort,
	}
}

func permittedValue[T comparable](value T, permittedValues ...T) bool {
	for _, val := range permittedValues {
		if val == value {
			return true
		}
	}
	return false
}
