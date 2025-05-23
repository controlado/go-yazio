package yazio

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/controlado/go-yazio/internal/application"
	"github.com/controlado/go-yazio/internal/testutil/assert"
	"github.com/controlado/go-yazio/internal/testutil/server"
	"github.com/controlado/go-yazio/internal/testutil/times"
	"github.com/controlado/go-yazio/pkg/client"
	"github.com/google/uuid"
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

			var (
				handler = func(t *testing.T, w http.ResponseWriter, r *http.Request) {
					responseBody := map[string]any{
						"access_token":  newAccessToken,
						"expires_in":    172800,
						"refresh_token": newRefreshToken,
						"token_type":    "bearer",
					}
					bodyBytes, err := json.Marshal(responseBody)
					assert.NoError(t, err)

					if tb.ServerStatus != 0 {
						w.WriteHeader(tb.ServerStatus)
					}

					_, err = w.Write(bodyBytes)
					assert.NoError(t, err)
				}
				srv        = server.New(t, handler)
				httpClient = client.New(client.WithBaseURL(srv.URL))

				ctx  = context.Background()
				user = &User{client: httpClient, token: tb.token}
				a    = &API{client: httpClient}
			)

			err := a.Refresh(ctx, user)
			assert.WantErr(t, tb.wantErr, err)

			var ( // post update
				userAccess  = user.token.Access()
				userRefresh = user.token.Refresh()
			)

			if tb.wantUpdate {
				assert.Equal(t, userAccess, newAccessToken)
				assert.Equal(t, userRefresh, newRefreshToken)
			} else {
				assert.Equal(t, userAccess, userAccess)
				assert.Equal(t, userRefresh, userRefresh)
			}
		})
	}
}
