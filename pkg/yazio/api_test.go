package yazio

import (
	"context"
	"net/http"
	"testing"

	"github.com/controlado/go-yazio/internal/testutil/assert"
	"github.com/controlado/go-yazio/internal/testutil/server"
	"github.com/controlado/go-yazio/pkg/client"
)

func TestNew(t *testing.T) {
	t.Parallel()

	c := client.New(
		client.WithBaseURL(BaseURL),
	)
	api, err := New(c)
	assert.NotNil(t, api)
	assert.NoError(t, err)
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

		assert.Write(t, w, respBody)
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
	user, err := api.Login(ctx, cred)
	assert.NotNil(t, user)
	assert.NoError(t, err)
}
