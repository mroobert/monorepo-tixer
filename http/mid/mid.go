// The package mid defines middlewares for the HTTP server.
package mid

import "net/http"

// MidHandler is a handler function designed to run code before and/or after
// another Handler. It is designed to remove boilerplate or other concerns not
// direct to any given app Handler.
type MidHandler func(next http.Handler) http.Handler
