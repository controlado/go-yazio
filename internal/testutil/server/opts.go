package server

import (
	"encoding/json"
	"fmt"
)

type Option func(tb *TestBuilder)

func RespondHeaders(h map[string]string) Option {
	return func(tb *TestBuilder) {
		tb.respondHeaders = h
	}
}

func RespondStatus(s int) Option {
	return func(tb *TestBuilder) {
		tb.respondStatus = s
	}
}

func RespondBody(b []byte) Option {
	return func(tb *TestBuilder) {
		tb.respondBody = b
	}
}

func RespondBodyString(b string) Option {
	return func(tb *TestBuilder) {
		tb.respondBody = []byte(b)
	}
}

func RespondBodyAny(b any) Option {
	return func(tb *TestBuilder) {
		bodyBytes, err := json.Marshal(b)
		if err != nil {
			err = fmt.Errorf("setting (any) body response: %w", err)
			tb.buildErrs = append(tb.buildErrs, err)
		}

		tb.respondBody = bodyBytes
	}
}

func AssertRequest() Option {
	return func(tb *TestBuilder) {
		tb.assertRequest = true
	}
}

func AssertEndpoint(p string) Option {
	return func(tb *TestBuilder) {
		tb.assertEndpoint = p
	}
}

func AssertMethod(m string) Option {
	return func(tb *TestBuilder) {
		tb.assertMethod = m
	}
}

func AssertBody(b map[string]any) Option {
	return func(tb *TestBuilder) {
		tb.assertBody = b
	}
}

func AssertHeaders(h map[string]string) Option {
	return func(tb *TestBuilder) {
		tb.assertHeaders = h
	}
}

func AssertQueryParams(qp map[string]string) Option {
	return func(tb *TestBuilder) {
		tb.assertQueryParams = qp
	}
}
