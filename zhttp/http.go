package zhttp

import (
	"context"
	"time"

	"net/http"

	"github.com/go-chi/chi/middleware"
	"go.uber.org/zap"
)

type loggerMiddleware struct {
	l *zap.Logger
	f LogFields
}

// LogFields customizes the fields logged
type LogFields map[string]func(context.Context) string

// NewMiddleware instantiates a middleware function that logs all requests
// using the provided logger
func NewMiddleware(l *zap.Logger, f LogFields) func(next http.Handler) http.Handler {
	return loggerMiddleware{
		// don't take stacktrace of wrapper class
		l.WithOptions(zap.AddCallerSkip(1)),
		f,
	}.Handler
}

func (l loggerMiddleware) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

		next.ServeHTTP(ww, r)

		fields := []zap.Field{
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
		}
		ctx := r.Context()
		for k, fn := range l.f {
			fields = append(fields, zap.String(k, fn(ctx)))
		}

		l.l.Info(r.Method+" "+r.URL.Path+": request completed", fields...)
	})
}
