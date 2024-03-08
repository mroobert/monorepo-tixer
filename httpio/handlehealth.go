package httpio

import (
	"net/http"
)

func (s *Server) handleHealthCheck(w http.ResponseWriter, r *http.Request) {
	env := envelope{
		"status": "available",
		"system_info": map[string]string{
			"environment": s.env,
		},
	}

	err := s.writeJSON(w, http.StatusOK, env, nil)
	if err != nil {
		s.internalServerErrorResponse(w, r, err)
	}
}
