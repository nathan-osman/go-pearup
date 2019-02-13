package auth

import (
	"net/http"

	"github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
	facebookOAuth2 "golang.org/x/oauth2/facebook"
	googleOAuth2 "golang.org/x/oauth2/google"
	microsoftOAuth2 "golang.org/x/oauth2/microsoft"
)

var endpoints = map[Provider]oauth2.Endpoint{
	ProviderFacebook:  facebookOAuth2.Endpoint,
	ProviderGoogle:    googleOAuth2.Endpoint,
	ProviderMicrosoft: microsoftOAuth2.LiveConnectEndpoint,
}

// Auth provides an interface for authenticating with an account.
type Auth struct {
	oauth2Configs         map[Provider]*oauth2.Config
	loginSucceededHandler http.Handler
	loginFailedHandler    http.Handler
	log                   *logrus.Entry
}

// New creates a new Auth instance with the specified configuration.
func New(cfg *Config) *Auth {
	a := &Auth{
		oauth2Configs:         make(map[Provider]*oauth2.Config),
		loginSucceededHandler: cfg.LoginSucceededHandler,
		loginFailedHandler:    cfg.LoginFailedHandler,
		log:                   logrus.WithField("context", "auth"),
	}
	for k, v := range cfg.ProviderConfigs {
		a.oauth2Configs[k] = &oauth2.Config{
			ClientID:     v.ClientID,
			ClientSecret: v.ClientSecret,
			RedirectURL:  v.RedirectURL,
			Endpoint:     endpoints[k],
			Scopes:       []string{"email"},
		}
	}
	return a
}
