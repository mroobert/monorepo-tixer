package httpio

import "net/http"

func (s *Server) registerTicketRoutes(r *http.ServeMux) {

	r.HandleFunc("POST /v1/tickets", s.handleCreateTicket)
	r.HandleFunc("GET /v1/tickets", s.handleGetAllTickets)
	r.HandleFunc("GET /v1/tickets/{id}", s.handleGetTicket)
	r.HandleFunc("DELETE /v1/tickets/{id}", s.handleDeleteTicket)
	r.HandleFunc("PATCH  /v1/tickets/{id}", s.handleUpdateTicket)

}

func (s *Server) handleCreateTicket(w http.ResponseWriter, r *http.Request) {
}

func (s *Server) handleGetAllTickets(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("get all tickets"))
}

func (s *Server) handleGetTicket(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("get ticket"))
}

func (s *Server) handleUpdateTicket(w http.ResponseWriter, r *http.Request) {
}

func (s *Server) handleDeleteTicket(w http.ResponseWriter, r *http.Request) {
}
