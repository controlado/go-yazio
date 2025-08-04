package yazio

import (
	"context"
	"fmt"
	"net/http"

	"github.com/controlado/go-yazio/internal/application"
	"github.com/controlado/go-yazio/internal/infra/client"
)

// API is the main struct for interacting with the YAZIO API.
//
// It holds the HTTP client used for making requests.
type API struct {
	client *client.Client
}

// New creates a new instance of the [*API].
func New(opts ...Option) (*API, error) {
	defaultAPI := &API{
		client: client.New(
			client.WithBaseURL(baseURL),
		),
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
