package client

import "net/http"

type Option func(c *Client)

func WithBaseURL(s string) Option {
	return func(c *Client) {
		c.BaseURL = s
	}
}

func WithRequester(r Requester) Option {
	if r == nil {
		r = http.DefaultClient
	}

	return func(c *Client) {
		c.Requester = r
	}
}
