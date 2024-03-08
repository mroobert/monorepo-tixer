package httpio

import (
	"log/slog"
	"net/http"

	"github.com/mroobert/monorepo-tixer/errs"
)

type errorBody struct {
	message string
	code    string
}

func (s *Server) logError(r *http.Request, err error) {
	slog.Error("encountered an error",
		slog.String("error", err.Error()),
		slog.String("request_method", r.Method),
		slog.String("request_url", r.URL.String()),
	)
}

func (s *Server) errorResponse(w http.ResponseWriter, r *http.Request, status int, code string, message string) {
	env := envelope{"error": errorBody{message, code}}

	err := s.writeJSON(w, status, env, nil)
	if err != nil {
		s.logError(r, err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (s *Server) internalServerErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	s.logError(r, err)

	message := "the server encountered a problem and could not process your request"
	s.errorResponse(w, r, http.StatusInternalServerError, errs.EINTERNAL, message)
}

func (s *Server) notFoundResponse(w http.ResponseWriter, r *http.Request) {
	message := "the requested resource could not be found"
	s.errorResponse(w, r, http.StatusNotFound, errs.ENOTFOUND, message)
}

func (s *Server) badRequestResponse(w http.ResponseWriter, r *http.Request, err error) {
	s.errorResponse(w, r, http.StatusBadRequest, errs.EINVALID, err.Error())
}
