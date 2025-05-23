package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type statusCategory int

const (
	Unknown       statusCategory = iota
	Informational                // 1xx
	Success                      // 2xx
	Redirection                  // 3xx
	ClientError                  // 4xx
	ServerError                  // 5xx
)

func (sc statusCategory) String() string {
	messages := [...]string{
		"Unknown",
		"Informational",
		"Success",
		"Redirection",
		"Client Error",
		"Server Error",
	}
	return messages[sc]
}

type Response struct {
	*http.Response
}

func (r *Response) check() error {
	statusCat := statusCategory(r.StatusCode / 100)

	switch statusCat {
	case Informational, Success, Redirection:
		return nil
	default:
		defer r.Body.Close()

		buffer, _ := io.ReadAll(r.Body)
		bufReader := bytes.NewReader(buffer)
		r.Body = io.NopCloser(bufReader)

		return fmt.Errorf(
			"unexpected status %d (%s): %s",
			r.StatusCode,
			statusCat,
			buffer,
		)
	}
}

func (r *Response) BodyString() (body string, err error) {
	defer func() { err = r.Body.Close() }()

	if r.ContentLength == 0 {
		return body, nil
	}

	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		return body, fmt.Errorf("reading body: %w", err)
	}

	return string(bodyBytes), nil
}

func (r *Response) BodyStruct(s any) (err error) {
	defer func() { err = r.Body.Close() }()

	if r.ContentLength == 0 {
		return nil
	}

	if err := json.NewDecoder(r.Body).Decode(s); err != nil {
		return fmt.Errorf("decoding body to struct: %w", err)
	}

	return nil
}
