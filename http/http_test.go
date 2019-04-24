package http

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"

	"github.com/bobheadxi/zapx/test"
)

func Test_loggerMiddleware(t *testing.T) {
	type args struct {
		method      string
		path        string
		body        io.Reader
		middlewares []func(http.Handler) http.Handler
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			"GET with request ID",
			args{"GET", "/", nil, []func(http.Handler) http.Handler{middleware.RequestID}},
			[]string{"req.path", "req.id"},
		},
		{
			"GET with req.ip",
			args{"GET", "/", nil, []func(http.Handler) http.Handler{middleware.RealIP}},
			[]string{"req.path", "req.ip"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// create bootstrapped logger and middleware
			var l, out = test.NewObservable()
			var handler = NewMiddleware(l, LogKeys{
				RequestID: "req.id",
			})

			// set up mock router
			m := chi.NewRouter()
			m.Use(tt.args.middlewares...)
			m.Use(handler)
			m.Get("/", func(w http.ResponseWriter, r *http.Request) {
				render.JSON(w, r, map[string]string{"hi": "bye"})
			})

			// create a mock request to use
			req := httptest.NewRequest(tt.args.method, "http://testing"+tt.args.path,
				tt.args.body)

			// serve request
			m.ServeHTTP(httptest.NewRecorder(), req)

			// check for desired log fields
			for _, e := range out.All() {
				for _, f := range tt.want {
					// find field, cast as string, and check if empty
					if val, _ := e.ContextMap()[f].(string); val == "" {
						t.Errorf("field %s unexpectedly empty", f)
					}
				}
			}
		})
	}
}
