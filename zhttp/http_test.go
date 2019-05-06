package zhttp

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bobheadxi/res"
	"go.bobheadxi.dev/zapx"
	"go.bobheadxi.dev/zapx/internal"
	"go.bobheadxi.dev/zapx/ztest"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest/observer"
)

func find(t *testing.T, out *observer.ObservedLogs, want map[string]bool) {
	for _, e := range out.All() {
		for _, f := range e.Context {
			if _, w := want[f.Key]; w {
				t.Log(f.Key, ":", f.String)
				assert.NotEmpty(t, f.String, "expected "+f.Key)
				want[f.Key] = true
			}
			if field, ok := f.Interface.(zapx.FieldSetMarshaller); ok {
				for _, sub := range field.Fields() {
					key := f.Key + "." + sub.Key
					if _, w := want[key]; w {
						t.Log(key, ":", sub.String)
						assert.NotEmpty(t, sub.String, "expected "+sub.Key)
						want[key] = true
					}
				}
			}
		}
	}
	for k, found := range want {
		assert.True(t, found, "should have found "+k)
	}
}

func TestMiddleware_Logger(t *testing.T) {
	type args struct {
		method      string
		path        string
		body        io.Reader
		middlewares []func(http.Handler) http.Handler
	}
	tests := []struct {
		name string
		args args
		want map[string]bool
	}{
		{
			"GET with request ID middleware",
			args{"GET", "/", nil, []func(http.Handler) http.Handler{middleware.RequestID}},
			map[string]bool{"req.id": false},
		},
		{
			"GET with real IP middleware",
			args{"GET", "/", nil, []func(http.Handler) http.Handler{middleware.RealIP}},
			map[string]bool{"req.ip": false},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// create bootstrapped logger and middleware
			var l, out = ztest.NewObservable()
			var md = NewMiddleware(l, LogFields{
				func(ctx context.Context) zap.Field {
					return zap.String("req.id", internal.String(ctx, middleware.RequestIDKey))
				},
			})

			// set up mock router
			m := chi.NewRouter()
			m.Use(tt.args.middlewares...)
			m.Use(md.Logger)
			m.Get("/", func(w http.ResponseWriter, r *http.Request) {
				res.R(w, r, res.MsgOK("hello world!"))
			})

			// create a mock request to use
			req := httptest.NewRequest(tt.args.method, "http://testing"+tt.args.path,
				tt.args.body)

			// serve request
			m.ServeHTTP(httptest.NewRecorder(), req)

			// check for desired log fields
			find(t, out, tt.want)
		})
	}
}

func TestMiddleware_Recoverer(t *testing.T) {
	type args struct {
		method      string
		path        string
		body        io.Reader
		middlewares []func(http.Handler) http.Handler
	}
	tests := []struct {
		name string
		args args
		want map[string]bool
	}{
		{
			"GET with request ID middleware and panic",
			args{"GET", "/", nil, []func(http.Handler) http.Handler{middleware.RequestID}},
			map[string]bool{"req.id": false, "panic": false, "req.ip": false},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// create bootstrapped logger and middleware
			var l, out = ztest.NewObservable()
			var md = NewMiddleware(l, LogFields{
				func(ctx context.Context) zap.Field {
					return zap.String("req.id", internal.String(ctx, middleware.RequestIDKey))
				},
			})

			// set up mock router
			m := chi.NewRouter()
			m.Use(tt.args.middlewares...)
			m.Use(md.Recoverer)
			m.Get("/", func(w http.ResponseWriter, r *http.Request) {
				panic("oh no")
			})

			// create a mock request to use
			req := httptest.NewRequest(tt.args.method, "http://testing"+tt.args.path,
				tt.args.body)

			// serve request
			m.ServeHTTP(httptest.NewRecorder(), req)

			// check for desired log fields
			find(t, out, tt.want)
		})
	}
}
