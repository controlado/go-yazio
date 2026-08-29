package yazio

import (
	"net/http"
	"time"
)

type Option func(a *API)

type Requester interface {
	Do(*http.Request) (*http.Response, error)
}

func WithRequester(r Requester) Option {
	return func(a *API) {
		if r != nil {
			a.client.Requester = r
		}
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
