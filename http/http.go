package http

import (
	"time"

	"net/http"

	"github.com/go-chi/chi/middleware"
	"go.uber.org/zap"

	"github.com/bobheadxi/zapx/httpctx"
)

type loggerMiddleware struct {
	l    *zap.Logger
	keys LogKeys
}

// LogKeys customizes the keys used for logging certain fields
type LogKeys struct {
	RequestID string
}

// NewMiddleware instantiates a middleware function that logs all requests
// using the provided logger
func NewMiddleware(l *zap.Logger, keys LogKeys) func(next http.Handler) http.Handler {
	return loggerMiddleware{
		// don't take stacktrace of wrapper class
		l.WithOptions(zap.AddCallerSkip(1)),
		keys,
	}.Handler
}

func (l loggerMiddleware) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

		next.ServeHTTP(ww, r)

		l.l.Info(r.Method+" "+r.URL.Path+": request completed",
			// request metadata
			zap.String("req.path", r.URL.Path),
			zap.String("req.query", r.URL.RawQuery),
			zap.String("req.method", r.Method),
			zap.String("req.ip", r.RemoteAddr),
			zap.String("req.user_agent", r.UserAgent()),

			// response metadata
			zap.Int("resp.status", ww.Status()),

			// additional metadata
			zap.Duration("duration", time.Since(start)),
			zap.String(l.keys.RequestID, httpctx.RequestID(r.Context())))
	})
}
