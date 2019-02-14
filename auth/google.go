package auth

import (
	"context"
	"net/http"

	"github.com/dghubble/gologin"
	"github.com/dghubble/gologin/google"
	oauth2Google "golang.org/x/oauth2/google"
)

// Google provides a Provider implementation for Google auth.
type Google struct {
	providerData
}

func NewGoogle(cfg *Config) *Google {
	g := &Google{}
	g.init(g, cfg, oauth2Google.Endpoint)
	return g
}

func (g *Google) Name() string {
	return "google"
}

func (g *Google) LoginHandler() http.Handler {
	return google.StateHandler(
		gologin.DebugOnlyCookieConfig,
		google.LoginHandler(
			g.config,
			g.errorHandler,
		),
	)
}

func (g *Google) CallbackHandler() http.Handler {
	return google.StateHandler(
		gologin.DebugOnlyCookieConfig,
		google.CallbackHandler(
			g.config,
			g.successHandler,
			g.errorHandler,
		),
	)
}

func (g *Google) User(ctx context.Context) (*User, error) {
	u, err := google.UserFromContext(ctx)
	if err != nil {
		return nil, err
	}
	return &User{
		ID:      u.Id,
		Name:    u.Name,
		Email:   u.Email,
		Picture: u.Picture,
	}, nil
}
