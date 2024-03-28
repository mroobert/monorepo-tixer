package psql

import (
	"fmt"
	"strings"
)

// Sorter is used to sort records.
type Sorter struct {
	column string
}

// NewSorter creates a new Sorter.
// It returns an error if the sort column is not in the safe list.
// The safeList parameter is used to prevent SQL injection.
func NewSorter(column string, safeList []string) (*Sorter, error) {
	sorter := &Sorter{column: column}

	err := sorter.check(safeList)
	if err != nil {
		return nil, err
	}

	return sorter, nil
}

// check returns an error if the sort column is not in the safe list.
func (s *Sorter) check(safeList []string) error {
	for _, safeValue := range safeList {
		if s.column == safeValue {
			return nil
		}
	}

	return fmt.Errorf("unsafe sort parameter: %s", s.column)
}

// Column returns the column to sort by.
func (s *Sorter) Column() string {
	return strings.TrimPrefix(s.column, "-")
}

// sortDirection retuen the direction to sort by (ASC or DESC).
func (s *Sorter) SortDirection() string {
	if strings.HasPrefix(s.column, "-") {
		return "DESC"
	}

	return "ASC"
}
