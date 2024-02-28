package http

import "net/http"

func (s *Server) registerTicketRoutes(r *http.ServeMux) {

	r.HandleFunc("POST /tickets", s.handleCreateTicket)
	r.HandleFunc("GET /tickets", s.handleGetAllTickets)
	r.HandleFunc("GET /tickets/{id}", s.handleGetTicket)
	r.HandleFunc("DELETE /tickets/{id}", s.handleDeleteTicket)
	r.HandleFunc("PATCH /tickets/{id}", s.handleUpdateTicket)

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
