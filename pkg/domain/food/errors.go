package food

import "errors"

var (
	ErrInvalidServing         = errors.New("given serving is invalid")
	ErrInvalidServingQuantity = errors.New("given serving quantity is invalid")
	ErrInvalidName            = errors.New("given food name is invalid")
	ErrAlreadyExists          = errors.New("given food already exists")
	ErrMissingNutrients       = errors.New("given food is missing some required nutrients")
)
