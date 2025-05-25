package server

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/controlado/go-yazio/internal/testutil/assert"
)

type TestBuilder struct {
	buildErrs []error

	receivedRequest bool

	respondHeaders map[string]string
	respondStatus  int
	respondBody    []byte

	assertRequest     bool
	assertEndpoint    string
	assertMethod      string
	assertBody        map[string]any
	assertHeaders     map[string]string
	assertQueryParams map[string]string
}

func New(t *testing.T, opts ...Option) (*httptest.Server, error) {
	t.Helper()

	tb := new(TestBuilder)

	for _, opt := range opts {
		opt(tb)
	}

	handler := func(w http.ResponseWriter, r *http.Request) {
		t.Helper()

		// to check with assertRequest
		tb.receivedRequest = true

		if tb.assertEndpoint != "" {
			assert.Equal(t, r.URL.Path, tb.assertEndpoint)
		}

		if tb.assertMethod != "" {
			assert.Equal(t, r.Method, tb.assertMethod)
		}

		if tb.assertBody != nil {
			rb := assert.ToJSON(t, r.Body)
			assert.DeepEqual(t, rb, tb.assertBody)
		}

		if tb.assertHeaders != nil {
			for k, v := range tb.assertHeaders {
				assert.Equal(t, r.Header.Get(k), v)
			}
		}

		if tb.assertQueryParams != nil {
			q := r.URL.Query()

			for k, v := range tb.assertQueryParams {
				assert.Equal(t, q.Get(k), v)
			}
		}

		if tb.respondHeaders != nil {
			h := w.Header()

			for k, v := range tb.respondHeaders {
				h.Set(k, v)
			}
		}

		if tb.respondStatus != 0 {
			w.WriteHeader(tb.respondStatus)
		}

		if tb.respondBody != nil {
			_, err := w.Write(tb.respondBody)
			assert.NoError(t, err)
		}
	}
	httpHandler := http.HandlerFunc(handler)
	srv := httptest.NewServer(httpHandler)

	t.Cleanup(
		func() {
			srv.Close()

			if tb.assertRequest && !tb.receivedRequest {
				t.Fatalf("want request, got none")
			}
		},
	)

	return srv, errors.Join(tb.buildErrs...)
}
