package yazio

import "github.com/controlado/go-yazio/internal/application"

type grantType string

const (
	passwordGrant grantType = "password"
	googleGrant   grantType = "sign_in_with_google"
)

type usingPassword struct {
	username, password string
}

func NewPasswordCred(username, password string) application.Credentials {
	return &usingPassword{
		username: username,
		password: password,
	}
}

func (up *usingPassword) Body() map[string]any {
	return map[string]any{
		"grant_type":    passwordGrant,
		"username":      up.username,
		"password":      up.password,
		"client_id":     defaultClientID,
		"client_secret": defaultSecret,
	}
}
