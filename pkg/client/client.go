package client

import (
	"context"
	"fmt"
	"net/http"
)

type Requester interface {
	Do(*http.Request) (*http.Response, error)
}

// Client provides a mechanism to execute HTTP requests.
//
// It encapsulates a [Requester] for performing the actual HTTP calls
// and a base URL that can be used as a default for requests.
//
// Instances of Client should be created using [New].
type Client struct {
	requester Requester
	baseURL   string
}

// New creates and returns a new [Client] instance.
//
// It initializes the client with a default requester [http.DefaultClient].
//
// The behavior of the new client can be customized by passing one or more
// [Option] values. These options are applied sequentially to configure fields
// such as the base URL or a custom requester.
func New(opts ...Option) *Client {
	defaultClient := &Client{
		requester: http.DefaultClient,
	}
	defaultClient.apply(opts...)
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

func (c *Client) apply(opts ...Option) {
	for _, o := range opts {
		o(c)
	}
}
