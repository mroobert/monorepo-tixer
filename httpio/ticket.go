package httpio

import (
	"fmt"
	"net/http"

	nanoid "github.com/matoous/go-nanoid/v2"
	tixer "github.com/mroobert/monorepo-tixer"
	"github.com/mroobert/monorepo-tixer/psql"
)

// registerTicketRoutes registers the ticket resource routes with the server.
func (s *Server) registerTicketRoutes(r *http.ServeMux) {
	r.HandleFunc("POST /v1/tickets", s.handleCreateTicket)
	r.HandleFunc("GET /v1/tickets", s.handleReadTickets)
	r.HandleFunc("GET /v1/tickets/{id}", s.handleReadTicket)
	r.HandleFunc("DELETE /v1/tickets/{id}", s.handleDeleteTicket)
	r.HandleFunc("PATCH /v1/tickets/{id}", s.handleUpdateTicket)
}

// ticketResponseBody represents the expected fields in the response body for a ticket resource.
type ticketResponseBody struct {
	PublicID string `json:"publicID"`
	Title    string `json:"title"`
	Price    int64  `json:"price"`
	Version  int32  `json:"version"`
}

// createTicketRequestBody represents the expected request body for creating a new ticket.
type createTicketRequestBody struct {
	Title string `json:"title"`
	Price int64  `json:"price"`
}

// handleCreateTicket handles the creation of a new ticket in the system.
func (s *Server) handleCreateTicket(w http.ResponseWriter, r *http.Request) {
	var body createTicketRequestBody

	err := s.readJSON(w, r, &body)
	if err != nil {
		s.badRequestResponse(w, r, err)
		return
	}

	validator := newValidator()
	if validator.validateCreateTicketRequestBody(body); !validator.valid() {
		s.failedValidationResponse(w, r, validator.errors)
		return
	}

	publicID, err := nanoid.Generate(tixer.PublicIDAlphabet, tixer.PublicIDLength)
	if err != nil {
		s.internalServerErrorResponse(w, r, err)
		return
	}

	ticket := tixer.Ticket{
		PublicID: tixer.PublicID(publicID),
		Title:    body.Title,
		Price:    body.Price,
	}
	if valid, errs := ticket.Validate(); !valid {
		s.failedValidationResponse(w, r, errs)
		return
	}

	ticketDB, err := s.TicketRepository.Insert(r.Context(), ticket)
	if err != nil {
		s.internalServerErrorResponse(w, r, err)
		return
	}

	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/tickets/%d", ticket.ID))

	err = s.writeJSON(w, http.StatusCreated, envelope{"ticket": ticketResponseBody{
		PublicID: string(ticketDB.PublicID),
		Title:    ticketDB.Title,
		Price:    ticketDB.Price,
		Version:  ticket.Version,
	}}, headers)
	if err != nil {
		s.internalServerErrorResponse(w, r, err)
	}
}

// urlQs represents the expected query string parameters for reading tickets.
type ticketUrlQs struct {
	title    string
	page     int
	pageSize int
	sort     string
}

// handleReadTickets handles reading tickets from the system.
func (s *Server) handleReadTickets(w http.ResponseWriter, r *http.Request) {
	sortSafeList := []string{"id", "title", "price", "-id", "-title", "-price"}

	validator := newValidator()
	qs := validator.validateTicketUrlValues(r.URL.Query(), sortSafeList)
	if !validator.valid() {
		s.failedValidationResponse(w, r, validator.errors)
		return
	}

	sorter, err := psql.NewSorter(qs.sort, sortSafeList)
	if err != nil {
		s.badRequestResponse(w, r, err)
		return
	}

	paginator := psql.NewPaginator(qs.page, qs.pageSize)

	ticketsDB, pagination, err := s.TicketRepository.SelectMultiple(r.Context(), psql.TicketFilter{
		Title:         qs.title,
		Limit:         paginator.Limit(),
		Offset:        paginator.Offset(),
		SortColumn:    sorter.Column(),
		SortDirection: sorter.SortDirection(),
	})
	if err != nil {
		s.internalServerErrorResponse(w, r, err)
		return
	}

	err = s.writeJSON(w, http.StatusOK, envelope{"tickets": toTicketResponseBody(ticketsDB), "pagination": pagination}, nil)
	if err != nil {
		s.internalServerErrorResponse(w, r, err)
	}
}

// handleReadTicket handles reading a single ticket from the system.
func (s *Server) handleReadTicket(w http.ResponseWriter, r *http.Request) {
	id, err := s.readIDParam(r)
	if err != nil {
		s.badRequestResponse(w, r, err)
		return
	}

	ticketDB, err := s.TicketRepository.SelectOne(r.Context(), id)
	if err != nil {
		switch err {
		case psql.ErrDbRecordNotFound:
			s.notFoundResponse(w, r)
		default:
			s.internalServerErrorResponse(w, r, err)
		}
		return
	}

	err = s.writeJSON(w, http.StatusOK, envelope{"ticket": ticketResponseBody{
		PublicID: string(ticketDB.PublicID),
		Title:    ticketDB.Title,
		Price:    ticketDB.Price,
		Version:  ticketDB.Version,
	}}, nil)
	if err != nil {
		s.internalServerErrorResponse(w, r, err)
	}
}

// updateTicketRequestBody represents the expected request body for updating an existing ticket.
type updateTicketRequestBody struct {
	Title *string `json:"title"`
	Price *int64  `json:"price"`
}

// handleUpdateTicket handles updating a ticket in the system.
func (s *Server) handleUpdateTicket(w http.ResponseWriter, r *http.Request) {
	var body updateTicketRequestBody

	id, err := s.readIDParam(r)
	if err != nil {
		s.badRequestResponse(w, r, err)
	}

	ticketDB, err := s.TicketRepository.SelectOne(r.Context(), id)
	if err != nil {
		switch err {
		case psql.ErrDbRecordNotFound:
			s.notFoundResponse(w, r)
		default:
			s.internalServerErrorResponse(w, r, err)
		}
		return
	}

	err = s.readJSON(w, r, &body)
	if err != nil {
		s.badRequestResponse(w, r, err)
	}

	if body.Title != nil {
		ticketDB.Title = *body.Title
	}
	if body.Price != nil {
		ticketDB.Price = *body.Price
	}

	valid, errs := ticketDB.Validate()
	if !valid {
		s.failedValidationResponse(w, r, errs)
	}

	err = s.TicketRepository.Update(r.Context(), &ticketDB)
	if err != nil {
		switch err {
		case psql.ErrDbEditConflict:
			s.editConflictResponse(w, r)
		default:
			s.internalServerErrorResponse(w, r, err)
		}
		return
	}

	err = s.writeJSON(w, http.StatusOK, envelope{"ticket": ticketResponseBody{
		PublicID: string(ticketDB.PublicID),
		Title:    ticketDB.Title,
		Price:    ticketDB.Price,
		Version:  ticketDB.Version,
	}}, nil)
	if err != nil {
		s.internalServerErrorResponse(w, r, err)
	}
}

// handleDeleteTicket handles the deletion of a ticket from the system.
func (s *Server) handleDeleteTicket(w http.ResponseWriter, r *http.Request) {
	id, err := s.readIDParam(r)
	if err != nil {
		s.badRequestResponse(w, r, err)
		return
	}

	err = s.TicketRepository.Delete(r.Context(), id)
	if err != nil {
		switch err {
		case psql.ErrDbRecordNotFound:
			s.notFoundResponse(w, r)
		default:
			s.internalServerErrorResponse(w, r, err)
		}
		return
	}

	err = s.writeJSON(w, http.StatusOK, envelope{"message": "ticket succesfully deleted"}, nil)
	if err != nil {
		s.internalServerErrorResponse(w, r, err)
	}
}

// toTicketResponseBody converts a slice of tickets that was read from DB
// to a slice of tickets that will be sent in the response body.
func toTicketResponseBody(ticketsDB []tixer.Ticket) []ticketResponseBody {
	tickets := make([]ticketResponseBody, len(ticketsDB))
	for i, ticketDB := range ticketsDB {
		tickets[i] = ticketResponseBody{
			PublicID: string(ticketDB.PublicID),
			Title:    ticketDB.Title,
			Price:    ticketDB.Price,
			Version:  ticketDB.Version,
		}
	}
	return tickets
}
