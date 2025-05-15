package yazio

import "github.com/controlado/go-yazio/internal/application"

type grantType string

const (
	passwordGrant grantType = "password"
	googleGrant   grantType = "sign_in_with_google"
	refreshGrant  grantType = "refresh_token"
)

type refreshToken struct {
	refreshTokenValue string
}

func newRefreshCred(tk application.Token) application.Credentials {
	return &refreshToken{
		refreshTokenValue: tk.Refresh(),
	}
}

func (rt *refreshToken) Body() map[string]any {
	return map[string]any{
		"grant_type":    refreshGrant,
		"refresh_token": rt.refreshTokenValue,
		"client_id":     defaultClientID,
		"client_secret": defaultSecret,
	}
}

type usingPassword struct {
	username, password string
}

// NewPasswordCred creates a new Credentials object
// for password-based authentication.
//
// It takes the username and password as input
// and returns an [application.Credentials] interface.
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
