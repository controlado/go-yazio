package yazio

import "errors"

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrExpiredToken       = errors.New("used token is invalid")

	ErrClientCannotBeNil = errors.New("given client cannot be nil")
	ErrRequestingToYazio = errors.New("failed to request to yazio's api")
	ErrDecodingResponse  = errors.New("failed to decode response's body -> internal dto")
)
