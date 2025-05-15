package yazio

import (
	"context"
	"fmt"
	"net/http"

	"github.com/controlado/go-yazio/internal/application"
	"github.com/controlado/go-yazio/pkg/client"
)

// API is the main struct for interacting with the YAZIO API.
//
// It holds the HTTP client used for making requests.
type API struct {
	client *client.Client
}

// New creates a new instance of the API.
//
// On failure the error wraps either:
//   - [ErrClientCannotBeNil]
func New(c *client.Client) (application.API, error) {
	if c == nil {
		return nil, ErrClientCannotBeNil
	}

	a := &API{
		client: c,
	}

	return a, nil
}

// Login implements the [application.API] interface.
// It attempts to log in a user with the provided cred.
//
// It returns an application.User containing the user's information
// upon successful login, or an error if the login fails.
//
// On failure the error wraps either:
//   - [ErrInvalidCredentials]
//   - [ErrRequestingToYazio]
//   - [ErrDecodingResponse]
//   - Other: generic (DTO related)
func (a *API) Login(ctx context.Context, cred application.Credentials) (application.User, error) {
	var (
		dto LoginDTO
		req = client.Request{
			Method:   http.MethodPost,
			Endpoint: loginEndpoint,
			Body:     cred.Body(),
		}
	)

	resp, err := a.client.Request(ctx, req)
	if err != nil {
		if resp.Response != nil {
			switch resp.StatusCode {
			case 400:
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
