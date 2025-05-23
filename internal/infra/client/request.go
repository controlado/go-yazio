package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

type (
	Payload[value any] map[string]value
)

type stringSetter interface {
	Set(string, string)
}

func setStrings(p Payload[string], s stringSetter) {
	for k, v := range p {
		s.Set(k, v)
	}
}

func (p Payload[T]) Reader() (*bytes.Reader, error) {
	if p == nil {
		return bytes.NewReader(nil), nil
	}

	payloadBytes, err := json.Marshal(p)
	if err != nil {
		return nil, fmt.Errorf("marshalling payload: %w", err)
	}

	return bytes.NewReader(payloadBytes), nil
}

type Request struct {
	BaseURL     string
	Method      string
	Endpoint    string
	Body        Payload[any]
	Headers     Payload[string]
	QueryParams Payload[string]
}

func (r *Request) HTTP(ctx context.Context) (*http.Request, error) {
	requestBody, err := r.Body.Reader()
	if err != nil {
		return nil, fmt.Errorf("getting body reader: %w", err)
	}

	rawURL, err := url.JoinPath(r.BaseURL, r.Endpoint)
	if err != nil {
		return nil, fmt.Errorf("preparing url path: %w", err)
	}

	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return nil, fmt.Errorf("parsing raw url: %w", err)
	}

	httpRequest, err := http.NewRequestWithContext(ctx, r.Method, parsedURL.String(), requestBody)
	if err != nil {
		return nil, fmt.Errorf("building request with context: %w", err)
	}

	if r.Headers != nil {
		setStrings(r.Headers, httpRequest.Header)
	}

	if r.QueryParams != nil {
		q := httpRequest.URL.Query()
		setStrings(r.QueryParams, q)
		httpRequest.URL.RawQuery = q.Encode()
	}

	return httpRequest, nil
}
