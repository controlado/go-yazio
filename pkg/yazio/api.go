package yazio

import (
	"context"
	"fmt"
	"net/http"

	"github.com/controlado/go-yazio/internal/application"
	"github.com/controlado/go-yazio/internal/domain"
	"github.com/controlado/go-yazio/pkg/client"
)

type API struct {
	client *client.Client
}

func New(c *client.Client) (application.API, error) {
	if c == nil {
		return nil, ErrClientCannotBeNil
	}

	a := &API{
		client: c,
	}

	return a, nil
}

// Login implements application.API.
func (a *API) Login(ctx context.Context, c application.Credentials) (application.User, error) {
	var (
		dto LoginDTO
		req = client.Request{
			Method:   http.MethodPost,
			Endpoint: loginEndpoint,
			Body:     c.Body(),
		}
	)

	resp, err := a.client.Request(ctx, req)
	if err != nil {
		if resp.Response != nil {
			switch resp.StatusCode {
			case 400:
				return nil, domain.ErrInvalidCredentials
			}
		}
		return nil, fmt.Errorf("%s: %w", ErrRequestingToYazio, err)
	}

	if err := resp.BodyStruct(&dto); err != nil {
		return nil, fmt.Errorf("%s: %w", ErrDecodingToInternalDTO, err)
	}

	return dto.toUser(a.client)
}
