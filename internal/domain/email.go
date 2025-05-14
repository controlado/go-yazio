package domain

import "strings"

type Email struct {
	Value       string
	IsConfirmed bool
}

func (e *Email) String() string {
	return strings.ToLower(e.Value)
}
