package client

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/controlado/go-yazio/internal/testutil/assert"
	"github.com/controlado/go-yazio/internal/testutil/server"
)

func TestNew(t *testing.T) {
	t.Parallel()

	var (
		blankRequester = &http.Client{}
		testBlocks     = []struct {
			name string
			want *Client
			args []Option
		}{
			{
				name: "no options",
				want: &Client{requester: http.DefaultClient},
			},
			{
				name: "custom requester",
				want: &Client{requester: blankRequester},
				args: []Option{
					WithRequester(blankRequester),
				},
			},
			{
				name: "nil requester",
				want: &Client{requester: http.DefaultClient},
				args: []Option{
					WithRequester(nil),
				},
			},
		}
	)

	for _, tb := range testBlocks {
		t.Run(tb.name, func(t *testing.T) {
			t.Parallel()
			got := New(tb.args...)
			assert.Equal(t, got.requester, tb.want.requester)
		})
	}
}

func TestClient_Request(t *testing.T) {
	t.Parallel()

	const (
		invalidBaseURL = " "
		validEndpoint  = "/test/this/is/valid"
	)

	var (
		testBlocks = []struct {
			name         string
			wantErr      bool
			ServerStatus int
			ServerBody   string
			ctx          context.Context
			req          Request
		}{
			{
				name: "with payloads",
				ctx:  context.Background(),
				req: Request{
					Method:      http.MethodPost,
					Endpoint:    validEndpoint,
					Body:        Payload{"success": true},
					Headers:     Payload{"fake-agent": "mock"},
					QueryParams: Payload{"name": "maria"},
				},
			},
			{
				name:    "with invalid base url",
				wantErr: true,
				ctx:     context.Background(),
				req: Request{
					Method:  http.MethodGet,
					BaseURL: invalidBaseURL,
				},
			},
			{
				name:    "with nil ctx",
				wantErr: true,
				ctx:     nil,
				req: Request{
					Method:   http.MethodGet,
					Endpoint: validEndpoint,
				},
			},
			{
				name:         "server responds bad",
				wantErr:      true,
				ctx:          context.Background(),
				ServerStatus: http.StatusInternalServerError,
				ServerBody:   "Please, try again later.",
				req: Request{
					Method:   http.MethodGet,
					Endpoint: validEndpoint,
				},
			},
		}
	)

	for _, tb := range testBlocks {
		t.Run(tb.name, func(t *testing.T) {
			t.Parallel()

			var (
				handler = func(t *testing.T, w http.ResponseWriter, r *http.Request) {
					assert.Equal(t, r.Method, tb.req.Method)
					assert.Equal(t, r.URL.Path, tb.req.Endpoint)

					if tb.req.Body != nil {
						var rBody Payload
						assert.DecodeDTO(t, r.Body, &rBody)
						assert.DeepEqual(t, rBody, tb.req.Body)
					}

					if tb.req.Headers != nil {
						for k, v := range tb.req.Headers {
							assert.DeepEqual(t, r.Header.Get(k), v)
						}
					}

					if tb.req.QueryParams != nil {
						q := r.URL.Query()
						for k, v := range tb.req.QueryParams {
							assert.DeepEqual(t, q.Get(k), v)
						}
					}

					if tb.ServerStatus != 0 {
						w.WriteHeader(tb.ServerStatus)
					}

					if tb.ServerBody != "" {
						bodyBytes := []byte(tb.ServerBody)
						w.Write(bodyBytes)
					}
				}
				srv = server.New(t, handler)
				c   = New(WithBaseURL(srv.URL))
			)

			resp, err := c.Request(tb.ctx, tb.req)
			if err == nil {
				if tb.wantErr {
					t.Fatal("want err, got nil")
				}
			} else {
				if !tb.wantErr {
					t.Fatalf("want nil error, got: %v", err)
				}

				if tb.ServerBody != "" {
					wantBody := fmt.Sprintf(
						"checking response: unexpected status %d (%s): %s",
						tb.ServerStatus, statusCategory(tb.ServerStatus/100), tb.ServerBody,
					)
					assert.Equal(t, err.Error(), wantBody)
				}
			}

			assert.NotNil(t, resp)
		})
	}
}
