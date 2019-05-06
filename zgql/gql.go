package zgql

import (
	"context"
	"net/http"
	"time"

	"github.com/99designs/gqlgen/graphql"
	"go.uber.org/zap"

	"go.bobheadxi.dev/zapx/internal"
)

// httpCtx are context keys used for injected HTTP variables, mostly for the
// convenience of the GraphQL logger
type httpCtx int

const (
	httpCtxKeyUserAgent httpCtx = iota + 1
	httpCtxKeyRemoteAddr
)

// GraphCtxHandler injects request fields into context for use with the GraphQL
// request logger. Should wrap the mux used to handle the GraphQL resolver.
func GraphCtxHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), httpCtxKeyUserAgent, r.UserAgent())
		ctx = context.WithValue(ctx, httpCtxKeyRemoteAddr, r.RemoteAddr)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// LogFields customizes the fields logged
type LogFields map[string]func(context.Context) string

// NewMiddleware returns a logger for use with GraphQL queries
func NewMiddleware(l *zap.Logger, f LogFields) graphql.RequestMiddleware {
	// don't take stacktrace of wrapper class
	l = l.WithOptions(zap.AddCallerSkip(1))

	return func(ctx context.Context, next func(context.Context) []byte) []byte {
		// call handler
		start := time.Now()
		response := next(ctx)

		// log request
		// TODO: could implement more advanced tracing, hm, via RequestContext.Trace
		req := graphql.GetRequestContext(ctx)
		// TODO: log message not very informative
		// TODO: evaluate usefulness of logged fields
		fields := []zap.Field{
			// request metadata
			zap.Int("req.complexity", req.OperationComplexity),
			zap.Any("req.variables", req.Variables),
			zap.String("req.ip", internal.String(ctx, httpCtxKeyRemoteAddr)),
			zap.String("req.user_agent", internal.String(ctx, httpCtxKeyUserAgent)),

			// response metadata
			zap.Bool("resp.errored", len(req.Errors) > 0),
			zap.Int("resp.size", len(response)),

			// additional metadata
			zap.Duration("duration", time.Since(start)),
		}
		for k, fn := range f {
			fields = append(fields, zap.String(k, fn(ctx)))
		}

		l.Info("graph query completed", fields...)

		return response
	}
}
