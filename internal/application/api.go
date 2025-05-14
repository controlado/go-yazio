package application

import "context"

type Credentials interface {
	Body() map[string]any
}

type API interface {
	Login(context.Context, Credentials) (User, error)
}
