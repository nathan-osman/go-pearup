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
	pictureDir            string
	log                   *logrus.Entry
	loginSucceededHandler http.Handler
	loginFailedHandler    http.Handler

	LoginHandler    http.Handler
	CallbackHandler http.Handler
}

// New creates a new Auth instance with the specified configuration.
func New(cfg *Config) *Auth {
	var (
		oauth2Config = &oauth2.Config{
			ClientID:     cfg.ClientID,
			ClientSecret: cfg.ClientSecret,
			RedirectURL:  cfg.RedirectURL,
			Endpoint:     facebookOAuth2.Endpoint,
		}
		a = &Auth{
			pictureDir: cfg.PictureDir,
			log:        logrus.WithField("context", "auth"),
		}
	)
	a.LoginHandler = facebook.StateHandler(
		gologin.DebugOnlyCookieConfig,
		facebook.LoginHandler(
			oauth2Config,
			http.HandlerFunc(a.loginFailed),
		),
	)
	a.CallbackHandler = facebook.StateHandler(
		gologin.DebugOnlyCookieConfig,
		facebook.CallbackHandler(
			oauth2Config,
			http.HandlerFunc(a.loginSucceeded),
			http.HandlerFunc(a.loginFailed),
		),
	)
	return a
}
