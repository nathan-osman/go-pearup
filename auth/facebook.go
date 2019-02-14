package auth

import (
	"context"
	"net/http"

	"github.com/dghubble/gologin"
	"github.com/dghubble/gologin/facebook"
	"github.com/dghubble/gologin/oauth2"
	"github.com/dghubble/sling"
	oauth2Facebook "golang.org/x/oauth2/facebook"
)

const facebookAPI = "https://graph.facebook.com/v3.0/"

// Facebook provides a Provider implementation for Facebook auth.
type Facebook struct {
	providerData
}

func NewFacebook(cfg *Config) *Facebook {
	f := &Facebook{}
	f.init(f, cfg, oauth2Facebook.Endpoint)
	return f
}

func (f *Facebook) Name() string {
	return "facebook"
}

func (f *Facebook) LoginHandler() http.Handler {
	return facebook.StateHandler(
		gologin.DebugOnlyCookieConfig,
		facebook.LoginHandler(
			f.config,
			f.errorHandler,
		),
	)
}

func (f *Facebook) CallbackHandler() http.Handler {
	return facebook.StateHandler(
		gologin.DebugOnlyCookieConfig,
		facebook.CallbackHandler(
			f.config,
			f.successHandler,
			f.errorHandler,
		),
	)
}

func (f *Facebook) User(ctx context.Context) (*User, error) {
	u, err := facebook.UserFromContext(ctx)
	if err != nil {
		return nil, err
	}
	t, err := oauth2.TokenFromContext(ctx)
	if err != nil {
		return nil, err
	}
	var (
		client = f.config.Client(ctx, t)
		p      = &struct {
			URL string `json:"url"`
		}{}
	)
	_, err = sling.
		New().
		Client(client).
		Base(facebookAPI).
		Set("Accept", "application/json").
		Get("me/picture?redirect=false").
		ReceiveSuccess(p)
	if err != nil {
		return nil, err
	}
	return &User{
		ID:      u.ID,
		Name:    u.Name,
		Email:   u.Email,
		Picture: p.URL,
	}, nil
}
