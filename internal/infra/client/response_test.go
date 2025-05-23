package client

import (
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/controlado/go-yazio/internal/testutil/assert"
)

func TestResponse_BodyString(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		httpResp *http.Response
		wantBody string
	}{
		{
			name: "valid body and content length",
			httpResp: &http.Response{
				ContentLength: 7,
				Body: io.NopCloser(
					strings.NewReader("success"),
				),
			},
			wantBody: "success",
		},
		{
			name: "no content length",
			httpResp: &http.Response{
				Body: io.NopCloser(
					strings.NewReader("success"),
				),
			},
		},
	}

	for _, tb := range tests {
		t.Run(tb.name, func(t *testing.T) {
			t.Parallel()

			resp := &Response{Response: tb.httpResp}
			got, err := resp.BodyString()
			assert.NoError(t, err)

			assert.Equal(t, got, tb.wantBody)
		})
	}
}

func TestResponse_BodyStruct(t *testing.T) {
	t.Parallel()

	type testDTO struct {
		APIMessage string `json:"api_message"`
	}

	testBlocks := []struct {
		name     string
		want     testDTO
		httpResp *http.Response
	}{
		{
			name: "valid struct",
			want: testDTO{
				APIMessage: "success",
			},
			httpResp: &http.Response{
				ContentLength: 26,
				Body: io.NopCloser(
					strings.NewReader(`{"api_message": "success"}`),
				),
			},
		},
	}

	for _, tb := range testBlocks {
		t.Run(tb.name, func(t *testing.T) {
			t.Parallel()

			var (
				got  testDTO
				resp = Response{Response: tb.httpResp}
			)

			if err := resp.BodyStruct(&got); err != nil {
				t.Fatalf("want err (nil), got %v", err)
			}

			assert.Equal(t, got, tb.want)
		})
	}
}
