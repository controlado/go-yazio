package food

import "errors"

var (
	ErrInvalidName      = errors.New("given food name is invalid")
	ErrAlreadyExists    = errors.New("given food already exists")
	ErrMissingNutrients = errors.New("given food is missing some required nutrients")
)
