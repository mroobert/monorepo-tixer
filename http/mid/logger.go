package mid

import (
	"fmt"
	"log/slog"
	"net/http"
)

// Logger writes information about the request start and end to the logs.
func Logger(next http.Handler) http.Handler {
	h := func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		if r.URL.RawQuery != "" {
			path = fmt.Sprintf("%s?%s", path, r.URL.RawQuery)
		}

		slog.Info("request started",
			slog.String("method", r.Method),
			slog.String("path", path),
			slog.String("remoteaddr", r.RemoteAddr),
		)

		next.ServeHTTP(w, r)

		slog.Info("request completed",
			slog.String("method", r.Method),
			slog.String("path", path),
			slog.String("remoteaddr", r.RemoteAddr),
		)

	}

	return http.HandlerFunc(h)
}
