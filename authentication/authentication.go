package authentication

import (
	"errors"
	"fmt"
	"github.com/go-resty/resty/v2"
	"log/slog"
)

// Auth0Config stores all data associated with the Auth0 rest API.
type Auth0Config struct {
	Auth0Domain  string
	ClientID     string
	ClientSecret string
}

// Auth stores the global logger object
type Auth struct {
	log    *slog.Logger
	config *Auth0Config
}

// New creates a new file
func New(log *slog.Logger, config *Auth0Config) *Auth {
	return &Auth{
		log:    log,
		config: config,
	}
}

// SendVerificationEmail sends the verification email to the specified user email address. It does this using the auth0
// rest APIs /passwordless/start function.
func (a *Auth) SendVerificationEmail(email string) error {
	client := resty.New()
	url := fmt.Sprintf("https://%s/passwordless/start", a.config.Auth0Domain)

	payload := map[string]string{
		"client_id":     a.config.ClientID,
		"client_secret": a.config.ClientSecret,
		"connection":    "email",
		"email":         email,
		"send":          "code",
	}

	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(payload).
		Post(url)
	if err != nil {
		a.log.Error("Failed to send verification email", "error", err)
		return err
	}

	if resp.StatusCode() != 200 {
		a.log.Error("Failed to send verification email with status code",
			"code", resp.StatusCode(), "status", resp.Status())
		return errors.New("failed to send verification email with status code")
	}

	a.log.Info("Verification email sent")
	return nil
}

func (a *Auth) ValidateOauthToken(email, code string) error {
	client := resty.New()
	url := fmt.Sprintf("https://%s/oauth/token", a.config.Auth0Domain)

	payload := map[string]string{
		"grant_type":    "http://auth0.com/oauth/grant-type/passwordless/otp",
		"client_id":     a.config.ClientID,
		"client_secret": a.config.ClientSecret,
		"username":      email,
		"otp":           code,
		"realm":         "email",
	}

	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(payload).
		Post(url)
	if err != nil {
		a.log.Error("Failed to validate oauth token", "error", err)
		return err
	}

	if resp.StatusCode() != 200 {
		a.log.Error("Failed to validate oauth token with status code",
			"code", resp.StatusCode(), "status", resp.Status())
		return errors.New("failed to validate oauth token with status code")
	}

	a.log.Info("Validated email successfully")
	return nil
}
