package gql

import (
	"context"
	"net/http"
	"time"

	"github.com/99designs/gqlgen/graphql"
	"go.uber.org/zap"

	"github.com/bobheadxi/zapx/httpctx"
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

// LogKeys customizes the keys used for logging certain fields
type LogKeys struct {
	RequestID string
}

// NewMiddleware returns a logger for use with GraphQL queries
func NewMiddleware(l *zap.Logger, keys LogKeys) graphql.RequestMiddleware {
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
		l.Info("graph query completed",
			// request metadata
			zap.Int("req.complexity", req.OperationComplexity),
			zap.Any("req.variables", req.Variables),
			zap.String("req.ip", httpctx.String(ctx, httpCtxKeyRemoteAddr)),
			zap.String("req.user_agent", httpctx.String(ctx, httpCtxKeyUserAgent)),

			// response metadata
			zap.Bool("resp.errored", len(req.Errors) > 0),
			zap.Int("resp.size", len(response)),

			// additional metadata
			zap.Duration("duration", time.Since(start)),
			zap.String(keys.RequestID, httpctx.RequestID(ctx)))

		return response
	}
}
