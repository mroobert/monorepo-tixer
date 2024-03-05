package mid

import "net/http"

// Panics recovers from panics and returns a 500 Internal Server Error.
func Panics(next http.Handler) http.Handler {
	h := func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if r := recover(); r != nil {
				w.Header().Set("Connection", "close")
				http.Error(w, "something went wrong", http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(h)
}
