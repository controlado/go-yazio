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
					Body:        Payload[any]{"success": true},
					Headers:     Payload[string]{"fake-agent": "mock"},
					QueryParams: Payload[string]{"name": "maria"},
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

			srv, err := server.New(t,
				server.AssertMethod(tb.req.Method),
				server.AssertEndpoint(tb.req.Endpoint),
				server.AssertHeaders(tb.req.Headers),
				server.AssertBody(tb.req.Body),
				server.AssertQueryParams(tb.req.QueryParams),
				server.RespondStatus(tb.ServerStatus),
				server.RespondBodyString(tb.ServerBody),
			)
			assert.NotNil(t, srv)
			assert.NoError(t, err)

			c := New(
				WithBaseURL(srv.URL),
			)
			assert.NotNil(t, c)

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
