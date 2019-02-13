package auth

import (
	"context"
	"net/http"
)

// withContext acts as a somewhat more concise way to inject context into a request.
func withContext(r *http.Request, key string, i interface{}) *http.Request {
	return r.WithContext(context.WithValue(r.Context(), key, i))
}
