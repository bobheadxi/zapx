package httpctx

import (
	"context"

	"github.com/go-chi/chi/middleware"
)

// RequestID returns the request ID injected by the chi requestID middleware
func RequestID(ctx context.Context) string { return String(ctx, middleware.RequestIDKey) }

// String gets a string with the given key from the given context's values
func String(ctx context.Context, key interface{}) string {
	if v := ctx.Value(key); v != nil {
		return v.(string)
	}
	return ""
}
