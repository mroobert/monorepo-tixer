package mid

import (
	"log/slog"
	"net/http"
	"runtime/debug"
)

// Panics recovers from panics and returns a 500 Internal Server Error.
func Panics(next http.Handler) http.Handler {
	h := func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if r := recover(); r != nil {
				slog.Error("PANIC recovered", "error", r, "stack", debug.Stack())
				w.Header().Set("Connection", "close")
				http.Error(w, "something went wrong", http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(h)
}
