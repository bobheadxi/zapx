package benchmarks

import (
	"net/http"
	"net/http/httptest"

	"go.bobheadxi.dev/res"
)

type handler struct{}

func (h handler) ServeHTTP(w http.ResponseWriter, r *http.Request) { res.R(w, r, res.MsgOK("ok!")) }

func newTestServer() *httptest.Server { return httptest.NewServer(handler{}) }
