package application

import "time"

type Token interface {
	Update(Token)

	ExpiresAt() time.Time
	Refresh() string
	Access() string
	Bearer() string
	String() string

	IsExpired() bool
}
