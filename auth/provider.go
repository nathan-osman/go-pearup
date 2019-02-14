package auth

import (
	"context"
	"net/http"

	"golang.org/x/oauth2"
)

// ProviderKey provides access to the Provider interface upon successful auth.
const ProviderKey = "provider"

// Provider implements authentication for a specific provider.
type Provider interface {

	// Name returns a unique name for the provider.
	Name() string

	// LoginHandler returns an HTTP handler for beginning auth.
	LoginHandler() http.Handler

	// CallbackHandler returns an HTTP handler for completing auth.
	CallbackHandler() http.Handler

	// User returns information about the user.
	User(context.Context) (*User, error)
}

type providerData struct {
	config         *oauth2.Config
	successHandler http.Handler
	errorHandler   http.Handler
}

func (p *providerData) init(provider Provider, cfg *Config, endpoint oauth2.Endpoint) {
	p.config = &oauth2.Config{
		ClientID:     cfg.ClientID,
		ClientSecret: cfg.ClientSecret,
		RedirectURL:  cfg.RedirectURL,
		Endpoint:     endpoint,
		Scopes:       []string{"email"},
	}
	p.successHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r = r.WithContext(context.WithValue(r.Context(), ProviderKey, provider))
		cfg.SuccessHandler.ServeHTTP(w, r)
	})
	p.errorHandler = cfg.ErrorHandler
}