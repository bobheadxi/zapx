package zhttp

import (
	"context"
	"net/http"
	"time"

	"github.com/go-chi/chi/middleware"
	"go.uber.org/zap"

	"go.bobheadxi.dev/zapx/zapx"
)

// Middleware is the zapx/zhttp middleware client
type Middleware struct {
	l *zap.Logger
	f LogFields
}

// LogFields customizes the fields logged
type LogFields []func(context.Context) zap.Field

// NewMiddleware instantiates a middleware function that logs all requests
// using the provided logger
func NewMiddleware(l *zap.Logger, f LogFields) *Middleware {
	return &Middleware{
		// don't take stacktrace of wrapper class
		l.WithOptions(zap.AddCallerSkip(1)),
		f,
	}
}

// Logger is basic logging middleware
func (m *Middleware) Logger(next http.Handler) http.Handler {
	const defaultFields = 3
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

		next.ServeHTTP(ww, r)

		fields := make([]zap.Field, defaultFields+len(m.f))
		fields[0] = requestFields(r)
		fields[1] = responseFields(ww)
		fields[2] = zap.Duration("duration", time.Since(start))
		var ctx = r.Context()
		for i, fn := range m.f {
			fields[i+defaultFields] = fn(ctx)
		}

		m.l.Info(r.Method+" "+r.URL.Path+": request completed", fields...)
	})
}

// Recoverer is middleware for logging panics
func (m *Middleware) Recoverer(next http.Handler) http.Handler {
	const defaultFields = 2
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rvr := recover(); rvr != nil {
				fields := make([]zap.Field, defaultFields+len(m.f))
				fields[0] = zap.Any("panic", rvr)
				fields[1] = requestFields(r)
				var ctx = r.Context()
				for i, fn := range m.f {
					fields[i+defaultFields] = fn(ctx)
				}

				m.l.Error(r.Method+" "+r.URL.Path+": recovered from panic", fields...)
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			}
		}()

		next.ServeHTTP(w, r)
	})
}

func requestFields(r *http.Request) zap.Field {
	return zapx.FieldSet("req",
		zap.String("path", r.URL.Path),
		zap.String("query", r.URL.RawQuery),
		zap.String("method", r.Method),
		zap.String("ip", r.RemoteAddr),
		zap.String("user_agent", r.UserAgent()))
}

func responseFields(ww middleware.WrapResponseWriter) zap.Field {
	return zapx.FieldSet("resp",
		zap.Int("status", ww.Status()))
}
