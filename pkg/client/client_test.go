package client

import (
	"context"
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

	var (
		ctx        = context.Background()
		testBlocks = []struct {
			name           string
			req            Request
			wantErr        bool
			wantStatusCode int
			checkBody      func(t *testing.T, respBody Payload)
			serverHandle   server.Handler
		}{
			{
				name: "POST with body",
				req: Request{
					Method:   http.MethodPost,
					Endpoint: "/body",
					Body:     Payload{"user": "feminismo"},
				},
				wantStatusCode: http.StatusOK,
				checkBody: func(t *testing.T, respBody Payload) {
					assert.Equal(t, respBody["success"], true)
				},
				serverHandle: func(t *testing.T, w http.ResponseWriter, r *http.Request) {
					assert.Equal(t, r.Method, http.MethodPost)
					assert.Equal(t, r.URL.Path, "/body")

					gotBody := assert.JSON(t, r.Body)
					assert.Equal(t, gotBody["user"], "feminismo")

					respBody := Payload{"success": true}
					assert.Write(t, w, respBody)
				},
			},
			{
				name: "GET with invalid base url",
				req: Request{
					Method:  http.MethodGet,
					BaseURL: "@",
				},
				wantErr:      true,
				serverHandle: func(t *testing.T, w http.ResponseWriter, r *http.Request) {},
			},
		}
	)

	for _, tb := range testBlocks {
		t.Run(tb.name, func(t *testing.T) {
			t.Parallel()

			var (
				srv = server.New(t, tb.serverHandle)
				c   = New(
					WithBaseURL(srv.URL),
				)
			)

			resp, err := c.Request(ctx, tb.req)
			if resp.Response != nil && resp.Body != nil {
				defer resp.Body.Close()
			}

			if (err != nil) != tb.wantErr {
				t.Fatalf("want err (%v), got %v", tb.wantErr, err)
			}

			if tb.wantStatusCode != 0 {
				assert.Equal(t, tb.wantStatusCode, resp.StatusCode)
			}

			if tb.checkBody != nil {
				respBody := assert.JSON(t, resp.Body)
				tb.checkBody(t, respBody)
			}
		})
	}
}
