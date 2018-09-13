package auth

import (
	"net/http"

	"github.com/dghubble/gologin"
	"github.com/dghubble/gologin/facebook"
	"github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
	facebookOAuth2 "golang.org/x/oauth2/facebook"
)

// Auth provides an interface for authenticating with a Facebook account.
type Auth struct {
	oauth2Config          *oauth2.Config
	log                   *logrus.Entry
	loginSucceededHandler http.Handler
	loginFailedHandler    http.Handler

	LoginHandler    http.Handler
	CallbackHandler http.Handler
}

// New creates a new Auth instance with the specified configuration.
func New(cfg *Config) *Auth {
	var (
		a = &Auth{
			oauth2Config: &oauth2.Config{
				ClientID:     cfg.ClientID,
				ClientSecret: cfg.ClientSecret,
				RedirectURL:  cfg.RedirectURL,
				Endpoint:     facebookOAuth2.Endpoint,
				Scopes:       []string{"email"},
			},
			log: logrus.WithField("context", "auth"),
			loginSucceededHandler: cfg.LoginSucceededHandler,
			loginFailedHandler:    cfg.LoginFailedHandler,
		}
	)
	a.LoginHandler = facebook.StateHandler(
		gologin.DebugOnlyCookieConfig,
		facebook.LoginHandler(
			a.oauth2Config,
			http.HandlerFunc(a.loginFailed),
		),
	)
	a.CallbackHandler = facebook.StateHandler(
		gologin.DebugOnlyCookieConfig,
		facebook.CallbackHandler(
			a.oauth2Config,
			http.HandlerFunc(a.loginSucceeded),
			http.HandlerFunc(a.loginFailed),
		),
	)
	return a
}
