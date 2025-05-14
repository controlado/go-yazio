package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

type Payload map[string]any

type setter interface {
	Set(string, string)
}

func (p Payload) Set(s setter) {
	for key, value := range p {
		if valueString, ok := value.(string); ok {
			s.Set(key, valueString)
		}
	}
}

func (p Payload) Reader() (*bytes.Reader, error) {
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
	Body        Payload
	Headers     Payload
	QueryParams Payload
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
		r.Headers.Set(httpRequest.Header)
	}

	if r.QueryParams != nil {
		q := httpRequest.URL.Query()
		r.QueryParams.Set(q)
		httpRequest.URL.RawQuery = q.Encode()
	}

	return httpRequest, nil
}
