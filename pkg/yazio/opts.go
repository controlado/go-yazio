package yazio

import "github.com/controlado/go-yazio/internal/infra/client"

type Option func(a *API)

func WithRequester(r client.Requester) Option {
	return func(a *API) {
		a.client = client.New(
			client.WithRequester(r),
			client.WithBaseURL(baseURL),
		)
	}
}

func WithBaseURL(bu string) Option {
	return func(a *API) {
		a.client = client.New(
			client.WithBaseURL(bu),
		)
	}
}
