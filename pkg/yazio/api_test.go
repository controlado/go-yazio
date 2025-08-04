package yazio

import (
	"context"
	"net/http"
	"testing"

	"github.com/controlado/go-yazio/internal/application"
	"github.com/controlado/go-yazio/internal/infra/client"
	"github.com/controlado/go-yazio/internal/testutil/assert"
	"github.com/controlado/go-yazio/internal/testutil/server"
	"github.com/controlado/go-yazio/internal/testutil/times"
	"github.com/google/uuid"
)

func TestNew(t *testing.T) {
	t.Parallel()

	api, err := New()
	assert.NotNil(t, api)
	assert.NoError(t, err)
}

func TestYazio_Login(t *testing.T) {
	t.Parallel()

	const (
		username = "testingUsername"
		password = "testingPassword"
	)

	respBody := loginDTO{
		ExpiresInSec: 172800,
		AccessToken:  "302af606a79142cb2ab862bf9488cfd4",
		RefreshToken: "302af606a79142cb2ab862bf9488cfd4",
	}
	srv, err := server.New(t,
		server.AssertMethod(http.MethodPost),
		server.AssertEndpoint(loginEndpoint),
		server.RespondBodyAny(respBody),
	)
	assert.NoError(t, err)
	assert.NotNil(t, srv)

	api, err := New(
		WithBaseURL(srv.URL),
	)
	assert.NotNil(t, api)
	assert.NoError(t, err)

	ctx := context.Background()
	cred := NewPasswordCred(username, password)
	user, err := api.Login(ctx, cred)
	assert.NotNil(t, user)
	assert.NoError(t, err)
}

func TestAPI_Refresh(t *testing.T) {
	var (
		randomToken = uuid.NewString()

		newAccessToken  = uuid.NewString()
		newRefreshToken = uuid.NewString()

		testBlocks = []struct {
			name         string
			wantErr      bool
			wantUpdate   bool
			ServerStatus int
			token        application.Token
		}{
			{
				name:       "using expired token should update",
				wantUpdate: true,
				token: &Token{
					expiresAt: times.Past(),
					access:    randomToken,
					refresh:   randomToken,
				},
			},
			{
				name: "using valid token should not update",
				token: &Token{
					expiresAt: times.Future(),
					access:    randomToken,
					refresh:   randomToken,
				},
			},
			{
				name:         "server responds invalid credentials",
				wantErr:      true,
				ServerStatus: http.StatusBadRequest,
				token: &Token{
					expiresAt: times.Past(),
					refresh:   randomToken,
				},
			},
		}
	)

	for _, tb := range testBlocks {
		t.Run(tb.name, func(t *testing.T) {
			t.Parallel()

			serverBody := map[string]any{
				"access_token":  newAccessToken,
				"expires_in":    172800,
				"refresh_token": newRefreshToken,
				"token_type":    "bearer",
			}
			srv, err := server.New(t,
				server.RespondBodyAny(serverBody),
				server.RespondStatus(tb.ServerStatus),
			)
			assert.NoError(t, err)
			assert.NotNil(t, srv)

			c := client.New(
				client.WithBaseURL(srv.URL),
			)
			assert.NotNil(t, c)

			var (
				ctx = context.Background()
				u   = &User{client: c, token: tb.token}
				a   = &API{client: c}
			)

			err = a.Refresh(ctx, u)
			assert.WantErr(t, tb.wantErr, err)

			var ( // post update
				userAccess  = u.token.Access()
				userRefresh = u.token.Refresh()
			)

			if tb.wantUpdate {
				assert.Equal(t, userAccess, newAccessToken)
				assert.Equal(t, userRefresh, newRefreshToken)
			} else {
				assert.Equal(t, userAccess, randomToken)
				assert.Equal(t, userRefresh, randomToken)
			}
		})
	}
}
