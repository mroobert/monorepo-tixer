package httpio

import (
	"encoding/json"
	"net/http"
)

// envelope is the response envelope.
type envelope map[string]any

// writeJSON writes the data to the http response writer as JSON.
func (s *Server) writeJSON(w http.ResponseWriter, status int, data envelope, headers http.Header) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	for key, value := range headers {
		w.Header()[key] = value
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(jsonData)

	return nil
}
