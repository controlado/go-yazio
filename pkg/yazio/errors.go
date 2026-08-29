package yazio

import "errors"

var (
	ErrUserCannotBeNil        = errors.New("given user cannot be nil")
	ErrTokenCannotBeNil       = errors.New("given token cannot be nil")
	ErrClientCannotBeNil      = errors.New("given client cannot be nil")
	ErrCredentialsCannotBeNil = errors.New("given credentials cannot be nil")

	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrExpiredToken       = errors.New("used token is invalid")

	ErrRequestingToYazio = errors.New("failed to request to yazio's api")
	ErrDecodingResponse  = errors.New("failed to decode response's body -> internal dto")
)
