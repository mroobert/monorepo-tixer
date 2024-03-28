package mid

import (
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/mroobert/monorepo-tixer/httpio/rcontext"
)

// Logger writes information about the request start and end to the logs.
func Logger(next http.Handler) http.Handler {
	h := func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		if r.URL.RawQuery != "" {
			path = fmt.Sprintf("%s?%s", path, r.URL.RawQuery)
		}

		info := rcontext.GetRequestInfo(r.Context())

		slog.InfoContext(r.Context(), "request start",
			slog.String("method", r.Method),
			slog.String("path", path),
			slog.String("remoteaddr", r.RemoteAddr),
		)

		next.ServeHTTP(w, r)

		slog.InfoContext(r.Context(), "request end",
			slog.String("method", r.Method),
			slog.String("path", path),
			slog.String("remoteaddr", r.RemoteAddr),
			slog.String("time", time.Since(info.StartedAt).String()),
		)

	}

	return http.HandlerFunc(h)
}
