package mid

import (
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/mroobert/monorepo-tixer/httpio/rcontext"
)

// ContextInfo adds request and response information to the context.
func ContextInfo(next http.Handler) http.Handler {
	h := func(w http.ResponseWriter, r *http.Request) {
		ctx := rcontext.SetRequestInfo(r.Context(), &rcontext.RequestInfo{
			RequestID: uuid.NewString(),
			StartedAt: time.Now(),
		})
		ctx = rcontext.SetResponseInfo(ctx, &rcontext.ResponseInfo{})

		next.ServeHTTP(w, r.WithContext(ctx))
	}
	return http.HandlerFunc(h)
}
