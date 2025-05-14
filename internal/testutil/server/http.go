package server

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

type Handler func(t *testing.T, w http.ResponseWriter, r *http.Request)

func New(t *testing.T, h Handler) *httptest.Server {
	t.Helper()
	srv := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				h(t, w, r)
			},
		),
	)
	t.Cleanup(srv.Close)
	return srv
}
