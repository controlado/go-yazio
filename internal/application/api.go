package application

import "context"

type Credentials interface {
	Body() map[string]any
}

type API interface {
	Refresh(context.Context, User) error
}
