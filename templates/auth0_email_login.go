package templates

type Auth0EmailLogin struct {
	Key        string
	EmailError bool
}

func NewAuth0EmailLogin(key string, emailError bool) *Auth0EmailLogin {
	return &Auth0EmailLogin{Key: key, EmailError: emailError}
}
