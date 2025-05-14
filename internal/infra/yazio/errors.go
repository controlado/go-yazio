package yazio

import "errors"

var (
	ErrClientCannotBeNil     = errors.New("given client cannot be nil")
	ErrRequestingToYazio     = errors.New("failed to request to yazio's api")
	ErrDecodingToInternalDTO = errors.New("failed to decode response's body -> internal dto")
)
