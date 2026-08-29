package client

import (
	"bytes"
	"encoding/json"
	"errors"
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

func (r *Response) check() (err error) {
	statusCat := statusCategory(r.StatusCode / 100)

	switch statusCat {
	case Informational, Success, Redirection:
		return nil
	default:
		defer closeAndJoin(r.Body, &err)

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
	defer closeAndJoin(r.Body, &err)

	if r.ContentLength == 0 {
		return body, nil
	}

	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		return body, err
	}

	return string(bodyBytes), nil
}

func (r *Response) BodyStruct(s any) (err error) {
	defer closeAndJoin(r.Body, &err)

	if r.ContentLength == 0 {
		return nil
	}

	return json.NewDecoder(r.Body).Decode(s)
}

func closeAndJoin(c io.Closer, err *error) {
	if closeErr := c.Close(); closeErr != nil {
		closeErr = fmt.Errorf("closing response body: %w", closeErr)
		*err = errors.Join(*err, closeErr)
	}
}
