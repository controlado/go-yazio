package yazio

import "github.com/controlado/go-yazio/internal/application"

const (
	passwordGrantType = "password"
)

type UsingPassword struct {
	username string
	password string
}

func NewPasswordCred(username, password string) application.Credentials {
	return &UsingPassword{
		username: username,
		password: password,
	}
}

func (up *UsingPassword) Body() map[string]any {
	return map[string]any{
		"grant_type":    passwordGrantType,
		"username":      up.username,
		"password":      up.password,
		"client_id":     defaultClientID,
		"client_secret": defaultSecret,
	}
}
