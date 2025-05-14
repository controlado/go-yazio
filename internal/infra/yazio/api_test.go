package yazio

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/controlado/go-yazio/internal/infra/client"
	"github.com/controlado/go-yazio/internal/testutil/assert"
	"github.com/controlado/go-yazio/internal/testutil/server"
)

func TestNew(t *testing.T) {
	t.Parallel()

	c := client.New(
		client.WithBaseURL(DefaultBaseURL),
	)
	api, err := New(c)

	switch {
	case err != nil:
		t.Fatalf("want Yazio, got err: %v", err)
	case api == nil:
		t.Fatalf("want Yazio, got nil")
	}
}

func TestYazio_Login(t *testing.T) {
	t.Parallel()

	const (
		username = "testingUsername"
		password = "testingPassword"
	)

	handler := func(t *testing.T, w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Method, http.MethodPost)
		assert.Equal(t, r.URL.Path, loginEndpoint)

		respBody := LoginDTO{
			ExpiresInSec: 172800,
			AccessToken:  "302af606a79142cb2ab862bf9488cfd4",
			RefreshToken: "302af606a79142cb2ab862bf9488cfd4",
		}

		err := json.NewEncoder(w).Encode(respBody)
		assert.NoError(t, err)
	}

	var (
		ctx = context.Background()
		srv = server.New(t, handler)
		c   = client.New(
			client.WithBaseURL(srv.URL),
		)
	)

	api, err := New(c)
	assert.NoError(t, err)

	cred := NewPasswordCred(username, password)
	_, err = api.Login(ctx, cred)
	assert.NoError(t, err)
}
