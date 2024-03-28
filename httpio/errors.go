package httpio

import (
	"log/slog"
	"net/http"

	tixer "github.com/mroobert/monorepo-tixer"
)

func (s *Server) logError(r *http.Request, err error) {
	slog.Error("encountered an error",
		slog.String("error", err.Error()),
		slog.String("request_method", r.Method),
		slog.String("request_url", r.URL.String()),
	)
}

func (s *Server) errorResponse(w http.ResponseWriter, r *http.Request, status int, code string, message string) {
	env := envelope{"error": map[string]string{
		"code":    code,
		"message": message,
	}}

	err := s.writeJSON(w, status, env, nil)
	if err != nil {
		s.logError(r, err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (s *Server) internalServerErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	s.logError(r, err)

	message := "the server encountered a problem and could not process your request"
	s.errorResponse(w, r, http.StatusInternalServerError, tixer.EINTERNAL, message)
}

func (s *Server) notFoundResponse(w http.ResponseWriter, r *http.Request) {
	message := "the requested resource could not be found"
	s.errorResponse(w, r, http.StatusNotFound, tixer.ENOTFOUND, message)
}

func (s *Server) badRequestResponse(w http.ResponseWriter, r *http.Request, err error) {
	s.errorResponse(w, r, http.StatusBadRequest, tixer.EINVALID, err.Error())
}

func (s *Server) editConflictResponse(w http.ResponseWriter, r *http.Request) {
	message := "unable to update the record due to an edit conflict, please try again"
	s.errorResponse(w, r, http.StatusConflict, tixer.ECONFLICT, message)
}

func (s *Server) failedValidationResponse(w http.ResponseWriter, r *http.Request, errors map[string]string) {
	env := envelope{"error": map[string]any{
		"code":        tixer.EUNPROCESSABLE,
		"validations": errors,
	}}

	err := s.writeJSON(w, http.StatusUnprocessableEntity, env, nil)
	if err != nil {
		s.logError(r, err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}
