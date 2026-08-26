package yazio

import (
	"time"

	"github.com/controlado/go-yazio/internal/infra/client"
)

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

func withNow(fn func() time.Time) Option {
	if fn == nil {
		fn = time.Now
	}
	return func(a *API) {
		a.now = fn
	}
}
