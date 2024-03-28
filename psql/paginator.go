package psql

import (
	"math"
)

// Paginator is used to paginate records.
type Paginator struct {
	page     int
	pageSize int
}

func NewPaginator(page, pageSize int) *Paginator {
	return &Paginator{
		page:     page,
		pageSize: pageSize,
	}
}

// limit returns the number of records to return per page.
func (p Paginator) Limit() int {
	return p.pageSize
}

// offset returns the number of records to skip.
func (p Paginator) Offset() int {
	return (p.page - 1) * p.pageSize
}

// Pagination represents the pagination information.
type Pagination struct {
	CurrentPage  int `json:"currentPage,omitempty"`
	PageSize     int `json:"pageSize,omitempty"`
	FirstPage    int `json:"firstPage,omitempty"`
	LastPage     int `json:"lastPage,omitempty"`
	TotalRecords int `json:"totalRecords,omitempty"`
}

func calculatePagination(totalRecords, page, pageSize int) Pagination {
	if totalRecords == 0 {

		return Pagination{}
	}

	return Pagination{
		CurrentPage:  page,
		PageSize:     pageSize,
		FirstPage:    1,
		LastPage:     int(math.Ceil(float64(totalRecords) / float64(pageSize))),
		TotalRecords: totalRecords,
	}
}
