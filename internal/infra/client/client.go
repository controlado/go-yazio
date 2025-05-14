package client

import (
	"context"
	"fmt"
	"net/http"
)

type Requester interface {
	Do(*http.Request) (*http.Response, error)
}

type Client struct {
	requester Requester
	baseURL   string
}

func New(opts ...Option) *Client {
	defaultClient := &Client{
		requester: http.DefaultClient,
	}

	for _, o := range opts {
		o(defaultClient)
	}

	return defaultClient
}

func (c *Client) Request(ctx context.Context, req Request) (resp Response, err error) {
	if req.BaseURL == "" {
		req.BaseURL = c.baseURL
	}

	httpRequest, err := req.HTTP(ctx)
	if err != nil {
		return resp, fmt.Errorf("parsing request to http.Request: %w", err)
	}

	resp.Response, err = c.requester.Do(httpRequest)
	if err != nil {
		return resp, fmt.Errorf("executing http.Request: %w", err)
	}

	if err := resp.check(); err != nil {
		return resp, fmt.Errorf("checking response: %w", err)
	}

	return
}
