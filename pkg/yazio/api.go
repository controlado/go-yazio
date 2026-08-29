package yazio

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/controlado/go-yazio/v3/internal/application"
	"github.com/controlado/go-yazio/v3/internal/infra/client"
)

// API is the main struct for interacting with the YAZIO API.
//
// It holds the HTTP client used for making requests.
type API struct {
	client *client.Client
	now    func() time.Time
}

// New creates a new instance of the [*API].
func New(opts ...Option) (*API, error) {
	defaultAPI := &API{
		client: client.New(
			client.WithBaseURL(baseURL),
		),
		now: time.Now,
	}

	for _, opt := range opts {
		opt(defaultAPI)
	}

	return defaultAPI, nil
}

func (a *API) Refresh(ctx context.Context, currentUser application.User) error {
	currentToken := currentUser.Token()
	if !currentToken.IsExpired() { // double-checking
		return nil
	}

	cred := newRefreshCred(currentToken)
	newUser, err := a.Login(ctx, cred)
	if err != nil {
		return fmt.Errorf("refreshing token: %w", err)
	}

	newToken := newUser.Token()
	currentToken.Update(newToken)

	return nil
}

// Login attempts to log in a user with the provided cred.
//
// It returns an [*User] containing the user's "connection"
// upon successful login, or an error if the login fails.
//
// On failure the error wraps either:
//   - [ErrInvalidCredentials]
//   - [ErrRequestingToYazio]
//   - [ErrDecodingResponse]
//   - Other: generic (DTO related)
func (a *API) Login(ctx context.Context, cred application.Credentials) (*User, error) {
	var (
		dto loginDTO
		req = client.Request{
			Method:   http.MethodPost,
			Endpoint: loginEndpoint,
			Headers:  defaultHeaders(nil),
			Body:     cred.Body(),
		}
	)

	resp, err := a.client.Request(ctx, req)
	if err != nil {
		if resp.Response != nil {
			switch resp.StatusCode {
			case http.StatusBadRequest:
				return nil, ErrInvalidCredentials
			}
		}
		return nil, fmt.Errorf("%s: %w", ErrRequestingToYazio, err)
	}

	if err := resp.BodyStruct(&dto); err != nil {
		return nil, fmt.Errorf("%s: %w", ErrDecodingResponse, err)
	}

	return dto.toUser(a.client)
}

// NewUserWithTokens creates a new [*User] instance using the provided
// access token, refresh token, and expiration time.
//
// This method allows you to create a User without going through
// the login flow, useful when you already have valid tokens.
//
// It returns an error if the access token or refresh token is empty.
func (a *API) NewUserWithTokens(accessToken, refreshToken string, expiresAt time.Time) (*User, error) {
	if accessToken == "" {
		return nil, fmt.Errorf("access token cannot be empty")
	}
	if refreshToken == "" {
		return nil, fmt.Errorf("refresh token cannot be empty")
	}

	return &User{
		client: a.client,
		token: &Token{
			expiresAt: expiresAt,
			access:    accessToken,
			refresh:   refreshToken,
		},
	}, nil
}

// NewUserWithTokensAndExpiresIn creates a new [*User] instance using the provided
// access token, refresh token, and expiration duration.
//
// This is a convenience method that calculates the expiration time
// by adding expiresIn to the current time.
//
// It returns an error if the access token or refresh token is empty.
func (a *API) NewUserWithTokensAndExpiresIn(accessToken, refreshToken string, expiresIn time.Duration) (*User, error) {
	expiresAt := a.now().Add(expiresIn)
	return a.NewUserWithTokens(accessToken, refreshToken, expiresAt)
}
