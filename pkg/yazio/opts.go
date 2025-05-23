package yazio

import "github.com/controlado/go-yazio/internal/infra/client"

type Option func(a *API)

func WithRequester(r client.Requester) Option {
	return func(a *API) {
		a.client.Requester = r
	}
}

func WithBaseURL(bu string) Option {
	return func(a *API) {
		a.client.BaseURL = bu
	}
}
